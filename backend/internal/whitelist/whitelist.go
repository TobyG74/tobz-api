// Package whitelist manages per-user IP allow-lists used to gate login and
// API-key usage. Each user may store up to MaxIPs entries; the sentinel
// "0.0.0.0" means public (any IP allowed). An empty list means unrestricted.
package whitelist

import (
	"errors"
	"net"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/tobz/tobz-api/internal/models"
)

const (
	MaxIPs    = 5
	PublicIP  = "0.0.0.0"
)

var (
	ErrLimitReached = errors.New("maksimal 5 IP")
	ErrInvalidIP    = errors.New("alamat IP tidak valid")
	ErrDuplicate    = errors.New("IP sudah ada")
	ErrNotFound     = errors.New("IP tidak ditemukan")
)

type Service struct{ db *gorm.DB }

func NewService(db *gorm.DB) *Service { return &Service{db: db} }

// List returns a user's whitelisted IPs (newest first).
func (s *Service) List(userID uuid.UUID) ([]models.WhitelistIP, error) {
	var ips []models.WhitelistIP
	err := s.db.Where("user_id = ?", userID).Order("created_at asc").Find(&ips).Error
	return ips, err
}

// Add stores a new IP after validation and the 5-entry cap.
func (s *Service) Add(userID uuid.UUID, ip, label string) (*models.WhitelistIP, error) {
	if net.ParseIP(ip) == nil {
		return nil, ErrInvalidIP
	}

	var count int64
	if err := s.db.Model(&models.WhitelistIP{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count >= MaxIPs {
		return nil, ErrLimitReached
	}

	rec := &models.WhitelistIP{UserID: userID, IP: ip, Label: trim(label, 40)}
	if err := s.db.Create(rec).Error; err != nil {
		// Unique index (user_id, ip) -> duplicate.
		return nil, ErrDuplicate
	}
	return rec, nil
}

// Remove deletes a user's IP (ownership-checked to prevent IDOR).
func (s *Service) Remove(userID, id uuid.UUID) error {
	res := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.WhitelistIP{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Allowed reports whether clientIP may act for userID.
// Empty list -> unrestricted. "0.0.0.0" present -> public. Otherwise exact match.
func (s *Service) Allowed(userID uuid.UUID, clientIP string) (bool, error) {
	var ips []models.WhitelistIP
	if err := s.db.Where("user_id = ?", userID).Find(&ips).Error; err != nil {
		return false, err
	}
	if len(ips) == 0 {
		return true, nil
	}
	norm := normalize(clientIP)
	for _, w := range ips {
		if w.IP == PublicIP || normalize(w.IP) == norm {
			return true, nil
		}
	}
	return false, nil
}

// normalize canonicalizes an IP string (handles IPv4-in-IPv6, etc.).
func normalize(ip string) string {
	if parsed := net.ParseIP(ip); parsed != nil {
		if v4 := parsed.To4(); v4 != nil {
			return v4.String()
		}
		return parsed.String()
	}
	return ip
}

func trim(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}
