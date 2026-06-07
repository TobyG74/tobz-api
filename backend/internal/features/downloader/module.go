package downloader

import (
	"github.com/gofiber/fiber/v2"

	"github.com/tobz/tobz-api/internal/httpclient"
)

// Feature name & base path — used by the server for logging/registration.
const (
	FeatureName = "downloader"
	BasePath    = "/download"
)

// Register mounts the downloader routes. This is the feature-module contract:
// each feature exposes one Register(router, deps) so the server enables it in one line.
func Register(router fiber.Router, client *httpclient.Safe) {
	h := &handler{registry: NewRegistry(client)}

	g := router.Group(BasePath)
	g.Get("/", h.download)
	g.Get("/platforms", h.platforms)
}
