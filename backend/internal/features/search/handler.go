package search

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/tobz/tobz-api/internal/response"
)

type handler struct {
	registry *Registry
}

// images searches images from a chosen source.
//
//	GET /search/images?q=cat&source=bing&limit=30
func (h *handler) images(c *fiber.Ctx) error {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		return response.BadRequest(c, "Parameter 'q' wajib")
	}
	src, err := h.registry.ImageSource(c.Query("source"))
	if err != nil {
		return response.BadRequest(c, "Sumber tidak didukung. Lihat /search/sources.")
	}
	limit, _ := strconv.Atoi(c.Query("limit"))

	results, err := src.Search(c.Context(), q, limit)
	if err != nil {
		return mapSearchError(c, err, src.ID())
	}
	return response.OK(c, fiber.Map{
		"source":  src.ID(),
		"query":   q,
		"count":   len(results),
		"results": results,
	})
}

// sources lists available image sources.
func (h *handler) sources(c *fiber.Ctx) error {
	return response.OK(c, fiber.Map{"sources": h.registry.Sources()})
}

// pixiv searches Pixiv.
//
//	GET /search/pixiv?q=miku&type=artworks&page=1
func (h *handler) pixiv(c *fiber.Ctx) error {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		return response.BadRequest(c, "Parameter 'q' wajib")
	}
	kind := c.Query("type", "artworks")
	page, _ := strconv.Atoi(c.Query("page"))

	results, err := h.registry.Pixiv().Search(c.Context(), kind, q, page)
	if err != nil {
		return mapSearchError(c, err, "pixiv")
	}
	return response.OK(c, fiber.Map{
		"query":   q,
		"type":    kind,
		"count":   len(results),
		"results": results,
	})
}

func mapSearchError(c *fiber.Ctx, err error, source string) error {
	switch {
	case errors.Is(err, ErrNoResults):
		return response.NotFound(c, "Tidak ada hasil ditemukan")
	case errors.Is(err, ErrNotConfigured):
		return response.Fail(c, fiber.StatusServiceUnavailable, "not_configured",
			"Sumber '"+source+"' belum dikonfigurasi di server")
	default:
		// Log the real cause server-side (not leaked to the client).
		log.Printf("[search] %s upstream error: %v", source, err)
		return response.Fail(c, fiber.StatusBadGateway, "upstream_error",
			"Gagal mengambil hasil dari "+source)
	}
}
