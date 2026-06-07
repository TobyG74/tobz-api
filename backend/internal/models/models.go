package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email         string    `gorm:"uniqueIndex;not null" json:"email"`
	EmailVerified bool      `gorm:"not null;default:false" json:"email_verified"`
	DisplayName   string    `gorm:"not null;default:''" json:"display_name"`
	AvatarURL     string    `gorm:"not null;default:''" json:"avatar_url,omitempty"`
	PasswordHash  string    `gorm:"not null;default:''" json:"-"`
	Role          string    `gorm:"not null;default:'user'" json:"role"` // user | admin

	FailedLoginCount int        `gorm:"not null;default:0" json:"-"`
	LockedUntil      *time.Time `json:"-"`

	OAuthAccounts []OAuthAccount `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	APIKeys       []APIKey       `gorm:"constraint:OnDelete:CASCADE" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type OAuthAccount struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Provider       string    `gorm:"not null;uniqueIndex:idx_provider_user" json:"provider"`
	ProviderUserID string    `gorm:"not null;uniqueIndex:idx_provider_user" json:"provider_user_id"`
	Email          string    `gorm:"not null;default:''" json:"email"`
	CreatedAt      time.Time `json:"created_at"`
}

func (o *OAuthAccount) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string    `gorm:"uniqueIndex;not null" json:"-"`
	Revoked   bool      `gorm:"not null;default:false" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type APIKey struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID  uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Name    string    `gorm:"not null;default:''" json:"name"`
	Prefix  string    `gorm:"not null;index" json:"prefix"` // 8 char
	KeyHash string    `gorm:"uniqueIndex;not null" json:"-"`
	Tier    string    `gorm:"not null;default:'free'" json:"tier"` // free | pro | vvip

	DailyQuota   int       `gorm:"not null;default:800" json:"daily_quota"`
	QuotaUsed    int       `gorm:"not null;default:0" json:"quota_used"`
	QuotaResetAt time.Time `gorm:"not null" json:"quota_reset_at"`

	Revoked    bool       `gorm:"not null;default:false" json:"revoked"`
	ExpiresAt  *time.Time `json:"expires_at"` // nil = no expiration
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

func (k *APIKey) BeforeCreate(tx *gorm.DB) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return nil
}

// WhitelistIP is an IP a user allows for login and API usage. Max 5 per user.
// The sentinel "0.0.0.0" means public (any IP allowed).
type WhitelistIP struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_ip" json:"user_id"`
	IP        string    `gorm:"not null;uniqueIndex:idx_user_ip" json:"ip"`
	Label     string    `gorm:"not null;default:''" json:"label"`
	CreatedAt time.Time `json:"created_at"`
}

func (w *WhitelistIP) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&OAuthAccount{},
		&RefreshToken{},
		&APIKey{},
		&WhitelistIP{},
	)
}
