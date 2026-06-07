package handlers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/tobz/tobz-api/internal/middleware"
	"github.com/tobz/tobz-api/internal/response"
	"github.com/tobz/tobz-api/internal/whitelist"
)

// ListWhitelist returns the user's whitelisted IPs.
func (h *Handlers) ListWhitelist(c *fiber.Ctx) error {
	ips, err := h.whitelist.List(middleware.UserID(c))
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, fiber.Map{"ips": ips, "max": whitelist.MaxIPs})
}

type addIPReq struct {
	IP    string `json:"ip"`
	Label string `json:"label"`
}

// AddWhitelist adds an IP (max 5; "0.0.0.0" means public).
func (h *Handlers) AddWhitelist(c *fiber.Ctx) error {
	var req addIPReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Body JSON tidak valid")
	}
	rec, err := h.whitelist.Add(middleware.UserID(c), strings.TrimSpace(req.IP), req.Label)
	if err != nil {
		switch {
		case errors.Is(err, whitelist.ErrLimitReached):
			return response.Fail(c, fiber.StatusConflict, "limit_reached", "Maksimal 5 IP per akun")
		case errors.Is(err, whitelist.ErrInvalidIP):
			return response.BadRequest(c, "Alamat IP tidak valid")
		case errors.Is(err, whitelist.ErrDuplicate):
			return response.Fail(c, fiber.StatusConflict, "duplicate", "IP itu sudah ada")
		default:
			return response.Internal(c)
		}
	}
	return response.Created(c, rec)
}

// RemoveWhitelist deletes an IP (ownership-checked).
func (h *Handlers) RemoveWhitelist(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "ID tidak valid")
	}
	if err := h.whitelist.Remove(middleware.UserID(c), id); err != nil {
		return response.NotFound(c, "IP tidak ditemukan")
	}
	return response.OK(c, fiber.Map{"message": "IP dihapus"})
}
