package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/tobz/tobz-api/internal/middleware"
	"github.com/tobz/tobz-api/internal/response"
)

type createKeyReq struct {
	Name string `json:"name"`
	Tier string `json:"tier"` // free | pro | vvip (non-free tiers are admin only)
}

func (h *Handlers) CreateAPIKey(c *fiber.Ctx) error {
	uid := middleware.UserID(c)
	var req createKeyReq
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Body JSON tidak valid")
	}

	tier := req.Tier
	if tier == "" {
		tier = "free"
	}
	if tier != "free" {
		if role, _ := c.Locals("role").(string); role != "admin" {
			return response.Forbidden(c, "Tier ini hanya untuk admin")
		}
	}

	raw, rec, err := h.keys.Create(uid, trimToLen(req.Name, 60), tier)
	if err != nil {
		return response.Internal(c)
	}

	return response.Created(c, fiber.Map{
		"id":          rec.ID,
		"name":        rec.Name,
		"tier":        rec.Tier,
		"daily_quota": rec.DailyQuota,
		"prefix":      rec.Prefix,
		"api_key":     raw, // shown only once
		"warning":     "Simpan API key ini sekarang. Nilai penuh tidak akan ditampilkan lagi.",
	})
}

func (h *Handlers) ListAPIKeys(c *fiber.Ctx) error {
	uid := middleware.UserID(c)
	keys, err := h.keys.List(uid)
	if err != nil {
		return response.Internal(c)
	}
	return response.OK(c, keys)
}

func (h *Handlers) RevokeAPIKey(c *fiber.Ctx) error {
	uid := middleware.UserID(c)
	keyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "ID key tidak valid")
	}
	if err := h.keys.Revoke(uid, keyID); err != nil {
		return response.NotFound(c, "Key tidak ditemukan")
	}
	return response.OK(c, fiber.Map{"message": "API key dicabut"})
}
