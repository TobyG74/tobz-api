package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tobz/tobz-api/internal/models"
)

const (
	keyPrefix    = "tobz_"
	rawKeyBytes  = 24
	displayChars = 8
)

var (
	ErrInvalidKey  = errors.New("api key tidak valid")
	ErrRevoked     = errors.New("api key dicabut")
	ErrExpired     = errors.New("api key kadaluarsa")
	ErrQuotaExceed = errors.New("kuota harian habis")
)

var tierQuota = map[string]int{
	"free": 800,
	"pro":  50000,
	"vvip": 1_000_000,
}

type Service struct{ db *gorm.DB }

func NewService(db *gorm.DB) *Service { return &Service{db: db} }

// hashKey stores only the SHA-256 hash, never the raw key.
func hashKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func nextMidnightWIB() time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
	return next
}

func (s *Service) Create(userID uuid.UUID, name, tier string) (rawKey string, rec *models.APIKey, err error) {
	quota, ok := tierQuota[tier]
	if !ok {
		tier = "free"
		quota = tierQuota["free"]
	}

	b := make([]byte, rawKeyBytes)
	if _, err = rand.Read(b); err != nil {
		return "", nil, err
	}
	rawKey = keyPrefix + hex.EncodeToString(b)

	rec = &models.APIKey{
		UserID:       userID,
		Name:         name,
		Prefix:       rawKey[:displayChars],
		KeyHash:      hashKey(rawKey),
		Tier:         tier,
		DailyQuota:   quota,
		QuotaUsed:    0,
		QuotaResetAt: nextMidnightWIB(),
	}
	if err = s.db.Create(rec).Error; err != nil {
		return "", nil, err
	}
	return rawKey, rec, nil
}

func (s *Service) VerifyAndConsume(rawKey string) (*models.APIKey, error) {
	if !strings.HasPrefix(rawKey, keyPrefix) {
		return nil, ErrInvalidKey
	}

	var key models.APIKey
	if err := s.db.Where("key_hash = ?", hashKey(rawKey)).First(&key).Error; err != nil {
		return nil, ErrInvalidKey
	}
	if key.Revoked {
		return nil, ErrRevoked
	}
	if key.ExpiresAt != nil && time.Now().After(*key.ExpiresAt) {
		return nil, ErrExpired
	}

	now := time.Now()
	// reset daily quota once the window has passed
	if now.After(key.QuotaResetAt) {
		s.db.Model(&key).Updates(map[string]interface{}{
			"quota_used":     0,
			"quota_reset_at": nextMidnightWIB(),
		})
		key.QuotaUsed = 0
	}

	// atomic quota update: increments only if under limit
	res := s.db.Model(&models.APIKey{}).
		Where("id = ? AND quota_used < daily_quota", key.ID).
		Updates(map[string]interface{}{
			"quota_used":   gorm.Expr("quota_used + 1"),
			"last_used_at": now,
		})
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrQuotaExceed
	}
	key.QuotaUsed++
	return &key, nil
}

func (s *Service) List(userID uuid.UUID) ([]models.APIKey, error) {
	var keys []models.APIKey
	err := s.db.Where("user_id = ?", userID).Order("created_at desc").Find(&keys).Error
	return keys, err
}

func (s *Service) Revoke(userID, keyID uuid.UUID) error {
	res := s.db.Model(&models.APIKey{}).
		Where("id = ? AND user_id = ?", keyID, userID).
		Update("revoked", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrInvalidKey
	}
	return nil
}
