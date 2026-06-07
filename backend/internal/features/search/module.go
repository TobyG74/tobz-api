package search

import (
	"github.com/gofiber/fiber/v2"

	"github.com/tobz/tobz-api/internal/httpclient"
)

const (
	FeatureName = "search"
	BasePath    = "/search"
)

// Register mounts the search feature's routes (feature-module contract).
func Register(router fiber.Router, client *httpclient.Safe) {
	h := &handler{registry: NewRegistry(client)}

	g := router.Group(BasePath)
	g.Get("/images", h.images)   // ?q=&source=&limit=
	g.Get("/sources", h.sources) // available image sources
	g.Get("/pixiv", h.pixiv)     // ?q=&type=&page=
}
