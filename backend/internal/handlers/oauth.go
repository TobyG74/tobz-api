package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/tobz/tobz-api/internal/auth"
	"github.com/tobz/tobz-api/internal/models"
	"github.com/tobz/tobz-api/internal/response"
)

const oauthStateCookie = "tobz_oauth_state"

// OAuthLogin starts the flow; the state cookie guards the callback against CSRF.
func (h *Handlers) OAuthLogin(c *fiber.Ctx) error {
	providerName := c.Params("provider")
	provider, err := h.oauth.Get(providerName)
	if err != nil {
		return response.NotFound(c, "Provider tidak didukung")
	}

	state, err := randomState()
	if err != nil {
		return response.Internal(c)
	}
	c.Cookie(&fiber.Cookie{
		Name:     oauthStateCookie,
		Value:    state,
		Path:     "/api/v1/auth/oauth",
		HTTPOnly: true,
		Secure:   h.cfg.IsProduction(),
		SameSite: "Lax",
		Expires:  time.Now().Add(10 * time.Minute),
	})

	return c.Redirect(provider.AuthCodeURL(state), fiber.StatusTemporaryRedirect)
}

// OAuthCallback validates state, exchanges the code, then finds-or-creates the user.
func (h *Handlers) OAuthCallback(c *fiber.Ctx) error {
	providerName := c.Params("provider")
	provider, err := h.oauth.Get(providerName)
	if err != nil {
		return h.redirectAuthError(c, "provider_unsupported")
	}

	// Anti-CSRF: cookie state must match query state.
	stateCookie := c.Cookies(oauthStateCookie)
	stateQuery := c.Query("state")
	if stateCookie == "" || stateQuery == "" || stateCookie != stateQuery {
		return h.redirectAuthError(c, "invalid_state")
	}
	// Single-use state.
	c.Cookie(&fiber.Cookie{Name: oauthStateCookie, Value: "", Path: "/api/v1/auth/oauth",
		HTTPOnly: true, Secure: h.cfg.IsProduction(), SameSite: "Lax", Expires: time.Now().Add(-time.Hour)})

	if errParam := c.Query("error"); errParam != "" {
		return h.redirectAuthError(c, "cancelled")
	}
	code := c.Query("code")
	if code == "" {
		return h.redirectAuthError(c, "no_code")
	}

	info, err := provider.Exchange(c.Context(), code)
	if err != nil {
		return h.redirectAuthError(c, "exchange_failed")
	}

	user, created, err := h.findOrCreateOAuthUser(info)
	if err != nil {
		return h.redirectAuthError(c, "server_error")
	}

	// On first OAuth login, seed the whitelist with the current IP.
	if created {
		_, _ = h.whitelist.Add(user.ID, c.IP(), "Login pertama")
	}

	// Enforce the user's IP whitelist (empty list = unrestricted).
	if ok, _ := h.whitelist.Allowed(user.ID, c.IP()); !ok {
		return h.redirectAuthError(c, "ip_not_allowed")
	}

	// Set the refresh cookie and redirect back to the SPA, which restores the
	// session via /auth/refresh on load. No token is placed in the URL.
	refresh, err := h.tokens.IssueRefreshToken(user.ID)
	if err != nil {
		return h.redirectAuthError(c, "server_error")
	}
	h.setRefreshCookie(c, refresh)
	return c.Redirect(h.cfg.FrontendURL+"/?login=success", fiber.StatusTemporaryRedirect)
}

// redirectAuthError sends the browser back to the SPA with an error code the UI can surface.
func (h *Handlers) redirectAuthError(c *fiber.Ctx, code string) error {
	return c.Redirect(h.cfg.FrontendURL+"/?auth_error="+code, fiber.StatusTemporaryRedirect)
}

// findOrCreateOAuthUser links by existing account, then by verified email, else
// creates a new user. The bool reports whether a brand-new user was created.
func (h *Handlers) findOrCreateOAuthUser(info *auth.OAuthUserInfo) (*models.User, bool, error) {
	var user models.User
	created := false

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var acct models.OAuthAccount
		err := tx.Where("provider = ? AND provider_user_id = ?", info.Provider, info.ProviderUserID).First(&acct).Error
		if err == nil {
			return tx.First(&user, "id = ?", acct.UserID).Error
		}
		if err != gorm.ErrRecordNotFound {
			return err
		}

		// Link via verified email.
		if info.Email != "" && info.EmailVerified {
			if e := tx.Where("email = ?", info.Email).First(&user).Error; e == nil {
				return tx.Create(&models.OAuthAccount{
					UserID: user.ID, Provider: info.Provider,
					ProviderUserID: info.ProviderUserID, Email: info.Email,
				}).Error
			}
		}

		user = models.User{
			Email:         info.Email,
			EmailVerified: info.EmailVerified,
			DisplayName:   sanitizeName(info.DisplayName, info.Email),
			AvatarURL:     info.AvatarURL,
			Role:          "user",
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		created = true
		return tx.Create(&models.OAuthAccount{
			UserID: user.ID, Provider: info.Provider,
			ProviderUserID: info.ProviderUserID, Email: info.Email,
		}).Error
	})
	if err != nil {
		return nil, false, err
	}
	return &user, created, nil
}

func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
