package downloader

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/tobz/tobz-api/internal/response"
)

// handler holds the registry and exposes this feature's HTTP endpoints.
type handler struct {
	registry *Registry
}

// download is the universal media downloader endpoint: GET /download?url=...
func (h *handler) download(c *fiber.Ctx) error {
	rawURL := c.Query("url")
	if rawURL == "" {
		return response.BadRequest(c, "Parameter 'url' wajib")
	}

	dl, err := h.registry.Resolve(rawURL)
	if err != nil {
		return response.BadRequest(c, "URL tidak didukung. Lihat /download/platforms untuk daftar platform.")
	}

	result, err := dl.Scrape(c.Context(), rawURL)
	if err != nil {
		if errors.Is(err, ErrNoMedia) {
			return response.NotFound(c, "Tidak ada media yang dapat diunduh dari URL tersebut")
		}
		// Don't leak upstream details to the client.
		return response.Fail(c, fiber.StatusBadGateway, "upstream_error",
			"Gagal mengambil media dari "+dl.PlatformName())
	}
	return response.OK(c, result)
}

// platforms lists supported platforms: GET /download/platforms
func (h *handler) platforms(c *fiber.Ctx) error {
	return response.OK(c, fiber.Map{"platforms": h.registry.Platforms()})
}
