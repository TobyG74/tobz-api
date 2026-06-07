package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/tobz/tobz-api/internal/apikey"
	"github.com/tobz/tobz-api/internal/auth"
	"github.com/tobz/tobz-api/internal/response"
	"github.com/tobz/tobz-api/internal/whitelist"
)

const (
	ctxUserID = "user_id"
	ctxRole   = "role"
	ctxAPIKey = "api_key"
)

// RequireAuth validates the Bearer JWT in the Authorization header.
func RequireAuth(tokens *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Get("Authorization")
		if !strings.HasPrefix(raw, "Bearer ") {
			return response.Unauthorized(c, "Authorization header (Bearer) wajib")
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(raw, "Bearer "))

		claims, err := tokens.ParseAccessToken(tokenStr)
		if err != nil {
			if errors.Is(err, auth.ErrExpiredToken) {
				return response.Unauthorized(c, "Token kadaluarsa")
			}
			return response.Unauthorized(c, "Token tidak valid")
		}

		uid, err := uuid.Parse(claims.Subject)
		if err != nil {
			return response.Unauthorized(c, "Token tidak valid")
		}
		c.Locals(ctxUserID, uid)
		c.Locals(ctxRole, claims.Role)
		return c.Next()
	}
}

// RequireAdmin requires an admin role.
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if role, _ := c.Locals(ctxRole).(string); role != "admin" {
			return response.Forbidden(c, "Akses khusus admin")
		}
		return c.Next()
	}
}

// RequireAPIKey validates the API key from the X-API-Key header or apikey query
// param, and enforces the key owner's IP whitelist.
func RequireAPIKey(keys *apikey.Service, wl *whitelist.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Get("X-API-Key")
		if raw == "" {
			raw = c.Query("apikey")
		}
		if raw == "" {
			return response.Unauthorized(c, "API key wajib (header X-API-Key)")
		}

		key, err := keys.VerifyAndConsume(raw)
		if err != nil {
			switch {
			case errors.Is(err, apikey.ErrQuotaExceed):
				return response.TooManyRequests(c, "Kuota harian API key habis")
			case errors.Is(err, apikey.ErrExpired):
				return response.Forbidden(c, "API key kadaluarsa")
			case errors.Is(err, apikey.ErrRevoked):
				return response.Forbidden(c, "API key dicabut")
			default:
				return response.Unauthorized(c, "API key tidak valid")
			}
		}

		// Enforce the key owner's IP whitelist (empty list = unrestricted).
		if ok, _ := wl.Allowed(key.UserID, c.IP()); !ok {
			return response.Forbidden(c, "IP ini tidak ada di whitelist akun")
		}

		c.Locals(ctxAPIKey, key)
		c.Set("X-Quota-Limit", itoa(key.DailyQuota))
		c.Set("X-Quota-Remaining", itoa(key.DailyQuota-key.QuotaUsed))
		return c.Next()
	}
}

// UserID returns the authenticated user's ID, or uuid.Nil if unset.
func UserID(c *fiber.Ctx) uuid.UUID {
	if v, ok := c.Locals(ctxUserID).(uuid.UUID); ok {
		return v
	}
	return uuid.Nil
}

func itoa(n int) string {
	if n < 0 {
		n = 0
	}
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
