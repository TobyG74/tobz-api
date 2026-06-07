package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/tobz/tobz-api/internal/auth"
	"github.com/tobz/tobz-api/internal/captcha"
	"github.com/tobz/tobz-api/internal/middleware"
	"github.com/tobz/tobz-api/internal/models"
	"github.com/tobz/tobz-api/internal/response"
)

const refreshCookieName = "tobz_refresh"

const (
	maxFailedLogins = 5
	lockoutDuration = 15 * time.Minute
)

type registerReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	DisplayName  string `json:"display_name"`
	CaptchaToken string `json:"captcha_token"`
}

type loginReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	CaptchaToken string `json:"captcha_token"`
}

type authResult struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int          `json:"expires_in"`
	User        *models.User `json:"user"`
}

func (h *Handlers) Register(c *fiber.Ctx) error {
	var req registerReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Body JSON tidak valid")
	}

	if err := h.captcha.Verify(c.Context(), req.CaptchaToken, c.IP()); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "captcha_failed", "Verifikasi captcha gagal")
	}

	email, ok := validateEmail(req.Email)
	if !ok {
		return response.BadRequest(c, "Format email tidak valid")
	}
	if !validatePassword(req.Password) {
		return response.BadRequest(c, "Password minimal 8 karakter dan mengandung huruf serta angka")
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return response.Internal(c)
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hash,
		DisplayName:  sanitizeName(req.DisplayName, email),
		Role:         "user",
	}
	if err := h.db.Create(user).Error; err != nil {
		return response.Fail(c, fiber.StatusConflict, "registration_failed", "Tidak dapat mendaftarkan akun ini")
	}

	return h.issueTokens(c, user, fiber.StatusCreated)
}

func (h *Handlers) Login(c *fiber.Ctx) error {
	var req loginReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Body JSON tidak valid")
	}

	if err := h.captcha.Verify(c.Context(), req.CaptchaToken, c.IP()); err != nil {
		if errors.Is(err, captcha.ErrCaptchaRequired) {
			return response.Fail(c, fiber.StatusBadRequest, "captcha_required", "Captcha wajib diisi")
		}
		return response.Fail(c, fiber.StatusBadRequest, "captcha_failed", "Verifikasi captcha gagal")
	}

	email, ok := validateEmail(req.Email)
	if !ok {
		return response.Unauthorized(c, "Email atau password salah")
	}

	var user models.User
	err := h.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_, _ = auth.VerifyPassword(req.Password, dummyHash)
		return response.Unauthorized(c, "Email atau password salah")
	}
	if err != nil {
		return response.Internal(c)
	}

	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return response.Fail(c, fiber.StatusTooManyRequests, "account_locked",
			"Akun terkunci sementara karena terlalu banyak percobaan gagal")
	}

	if user.PasswordHash == "" {
		return response.Unauthorized(c, "Akun ini menggunakan login sosial. Silakan login via provider.")
	}

	valid, _ := auth.VerifyPassword(req.Password, user.PasswordHash)
	if !valid {
		h.recordFailedLogin(&user)
		return response.Unauthorized(c, "Email atau password salah")
	}

	if user.FailedLoginCount != 0 || user.LockedUntil != nil {
		h.db.Model(&user).Updates(map[string]interface{}{"failed_login_count": 0, "locked_until": nil})
	}

	return h.issueTokens(c, &user, fiber.StatusOK)
}

func (h *Handlers) Refresh(c *fiber.Ctx) error {
	raw := c.Cookies(refreshCookieName)
	if raw == "" {
		return response.Unauthorized(c, "Refresh token tidak ada")
	}

	newRaw, user, err := h.tokens.RotateRefreshToken(raw)
	if err != nil {
		h.clearRefreshCookie(c)
		return response.Unauthorized(c, "Refresh token tidak valid")
	}

	access, err := h.tokens.GenerateAccessToken(user)
	if err != nil {
		return response.Internal(c)
	}
	h.setRefreshCookie(c, newRaw)
	return response.OK(c, authResult{
		AccessToken: access,
		TokenType:   "Bearer",
		ExpiresIn:   int(h.cfg.AccessTokenTTL.Seconds()),
		User:        user,
	})
}

func (h *Handlers) Logout(c *fiber.Ctx) error {
	if raw := c.Cookies(refreshCookieName); raw != "" {
		_ = h.tokens.RevokeRefreshToken(raw)
	}
	h.clearRefreshCookie(c)
	return response.OK(c, fiber.Map{"message": "Berhasil logout"})
}

func (h *Handlers) Me(c *fiber.Ctx) error {
	uid := middleware.UserID(c)
	var user models.User
	if err := h.db.First(&user, "id = ?", uid).Error; err != nil {
		return response.NotFound(c, "User tidak ditemukan")
	}
	return response.OK(c, user)
}

func (h *Handlers) recordFailedLogin(user *models.User) {
	updates := map[string]interface{}{"failed_login_count": user.FailedLoginCount + 1}
	if user.FailedLoginCount+1 >= maxFailedLogins {
		until := time.Now().Add(lockoutDuration)
		updates["locked_until"] = &until
		updates["failed_login_count"] = 0
	}
	h.db.Model(user).Updates(updates)
}

func (h *Handlers) issueTokens(c *fiber.Ctx, user *models.User, status int) error {
	access, err := h.tokens.GenerateAccessToken(user)
	if err != nil {
		return response.Internal(c)
	}
	refresh, err := h.tokens.IssueRefreshToken(user.ID)
	if err != nil {
		return response.Internal(c)
	}
	h.setRefreshCookie(c, refresh)
	return c.Status(status).JSON(response.Envelope{
		Success: true,
		Data: authResult{
			AccessToken: access,
			TokenType:   "Bearer",
			ExpiresIn:   int(h.cfg.AccessTokenTTL.Seconds()),
			User:        user,
		},
	})
}

func (h *Handlers) setRefreshCookie(c *fiber.Ctx, value string) {
	c.Cookie(&fiber.Cookie{
		Name:     refreshCookieName,
		Value:    value,
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   h.cfg.IsProduction(), // HTTPS-only in production
		SameSite: "Lax",                // CSRF defense
		Expires:  time.Now().Add(h.cfg.RefreshTokenTTL),
	})
}

func (h *Handlers) clearRefreshCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/api/v1/auth",
		HTTPOnly: true,
		Secure:   h.cfg.IsProduction(),
		SameSite: "Lax",
		Expires:  time.Now().Add(-time.Hour),
	})
}

func sanitizeName(name, fallbackEmail string) string {
	name = trimToLen(name, 60)
	if name == "" {
		for i, r := range fallbackEmail {
			if r == '@' {
				return fallbackEmail[:i]
			}
		}
	}
	return name
}

func trimToLen(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

var dummyHash = "$argon2id$v=19$m=65536,t=3,p=2$YWFhYWFhYWFhYWFhYWFhYQ$Yk5xZ2J1c2FtcGxlZHVtbXloYXNodmFsdWUwMDAwMA"
