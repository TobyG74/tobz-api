package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/tobz/tobz-api/internal/auth"
	"github.com/tobz/tobz-api/internal/middleware"
	"github.com/tobz/tobz-api/internal/models"
	"github.com/tobz/tobz-api/internal/response"
)

type changePasswordReq struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ChangePassword updates the user's password. If the account already has a
// password, the current one must be verified first. On success all existing
// refresh tokens are revoked (logout everywhere) and a fresh session is issued
// for the current device.
func (h *Handlers) ChangePassword(c *fiber.Ctx) error {
	var req changePasswordReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Body JSON tidak valid")
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", middleware.UserID(c)).Error; err != nil {
		return response.NotFound(c, "User tidak ditemukan")
	}

	// Verify current password only if the account already has one
	// (OAuth-only accounts may set a password without it).
	if user.PasswordHash != "" {
		valid, _ := auth.VerifyPassword(req.CurrentPassword, user.PasswordHash)
		if !valid {
			return response.Unauthorized(c, "Password saat ini salah")
		}
	}

	if !validatePassword(req.NewPassword) {
		return response.BadRequest(c, "Password baru minimal 8 karakter dan mengandung huruf serta angka")
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return response.Internal(c)
	}
	if err := h.db.Model(&user).Update("password_hash", hash).Error; err != nil {
		return response.Internal(c)
	}

	// Invalidate every existing session, then re-issue one for this device.
	_ = h.tokens.RevokeAllForUser(user.ID)
	user.PasswordHash = hash
	return h.issueTokens(c, &user, fiber.StatusOK)
}
