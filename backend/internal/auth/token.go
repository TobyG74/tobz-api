package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tobz/tobz-api/internal/config"
	"github.com/tobz/tobz-api/internal/models"
)

var (
	ErrInvalidToken = errors.New("token tidak valid")
	ErrExpiredToken = errors.New("token kadaluarsa")
)

type Claims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type Service struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewService(cfg *config.Config, db *gorm.DB) *Service {
	return &Service{cfg: cfg, db: db}
}

func (s *Service) GenerateAccessToken(user *models.User) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.AccessTokenTTL)),
			Issuer:    "tobz-api",
			ID:        uuid.NewString(),
		},
		Role: user.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.cfg.JWTSecret)
}

func (s *Service) ParseAccessToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.cfg.JWTSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (s *Service) IssueRefreshToken(userID uuid.UUID) (string, error) {
	raw, err := randomToken(32)
	if err != nil {
		return "", err
	}
	rt := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(raw),
		ExpiresAt: time.Now().Add(s.cfg.RefreshTokenTTL),
	}
	if err := s.db.Create(rt).Error; err != nil {
		return "", err
	}
	return raw, nil
}

// RotateRefreshToken revokes the old token and issues a new one atomically.
func (s *Service) RotateRefreshToken(raw string) (newRaw string, user *models.User, err error) {
	var rt models.RefreshToken
	if err = s.db.Where("token_hash = ?", hashToken(raw)).First(&rt).Error; err != nil {
		return "", nil, ErrInvalidToken
	}
	if rt.Revoked || time.Now().After(rt.ExpiresAt) {
		return "", nil, ErrInvalidToken
	}

	var u models.User
	if err = s.db.First(&u, "id = ?", rt.UserID).Error; err != nil {
		return "", nil, ErrInvalidToken
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Model(&rt).Update("revoked", true).Error; e != nil {
			return e
		}
		nrt := &models.RefreshToken{
			UserID:    u.ID,
			TokenHash: hashToken(mustNewRaw(&newRaw)),
			ExpiresAt: time.Now().Add(s.cfg.RefreshTokenTTL),
		}
		return tx.Create(nrt).Error
	})
	if err != nil {
		return "", nil, err
	}
	return newRaw, &u, nil
}

func (s *Service) RevokeRefreshToken(raw string) error {
	return s.db.Model(&models.RefreshToken{}).
		Where("token_hash = ?", hashToken(raw)).
		Update("revoked", true).Error
}

func (s *Service) RevokeAllForUser(userID uuid.UUID) error {
	return s.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// hashToken stores only the SHA-256 hash, never the raw token.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func mustNewRaw(out *string) string {
	r, err := randomToken(32)
	if err != nil {
		panic(err)
	}
	*out = r
	return r
}
