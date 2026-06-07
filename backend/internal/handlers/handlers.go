// Package handlers contains all HTTP handlers (Fiber).
package handlers

import (
	"gorm.io/gorm"

	"github.com/tobz/tobz-api/internal/apikey"
	"github.com/tobz/tobz-api/internal/auth"
	"github.com/tobz/tobz-api/internal/captcha"
	"github.com/tobz/tobz-api/internal/config"
	"github.com/tobz/tobz-api/internal/httpclient"
)

// Handlers holds shared dependencies for the core auth/account handlers.
type Handlers struct {
	cfg     *config.Config
	db      *gorm.DB
	tokens  *auth.Service
	keys    *apikey.Service
	captcha *captcha.Verifier
	oauth   *auth.OAuthRegistry
	http    *httpclient.Safe
}

func New(
	cfg *config.Config,
	db *gorm.DB,
	tokens *auth.Service,
	keys *apikey.Service,
	cap *captcha.Verifier,
	oauth *auth.OAuthRegistry,
	safe *httpclient.Safe,
) *Handlers {
	return &Handlers{
		cfg:     cfg,
		db:      db,
		tokens:  tokens,
		keys:    keys,
		captcha: cap,
		oauth:   oauth,
		http:    safe,
	}
}
