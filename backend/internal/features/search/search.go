package search

import (
	"context"
	"errors"
	"strings"

	"github.com/tobz/tobz-api/internal/httpclient"
)

var (
	ErrUnknownSource = errors.New("unknown search source")
	ErrNoResults     = errors.New("no results found")
	ErrNotConfigured = errors.New("source not configured")
)

// ImageResult is a unified image-search hit across sources.
type ImageResult struct {
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Title     string `json:"title,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Source    string `json:"source,omitempty"`
	// Pinterest extras.
	Author string `json:"author,omitempty"`
	Type   string `json:"type,omitempty"` // image | video | gif
}

// ImageSource searches images for a query.
type ImageSource interface {
	ID() string
	Search(ctx context.Context, query string, limit int) ([]ImageResult, error)
}

// Registry holds the active image sources plus the Pixiv client.
type Registry struct {
	sources map[string]ImageSource
	order   []string
	pixiv   *pixivClient
}

// NewRegistry wires every source, sharing one anti-SSRF HTTP client.
func NewRegistry(client *httpclient.Safe) *Registry {
	b := &base{http: client}
	r := &Registry{sources: map[string]ImageSource{}, pixiv: &pixivClient{b}}

	for _, s := range []ImageSource{
		&bingSource{b},
		&duckduckgoSource{b},
		&pinterestSource{b},
		&pexelsSource{b},
	} {
		r.sources[s.ID()] = s
		r.order = append(r.order, s.ID())
	}
	return r
}

// ImageSource returns a source by id (defaults to bing if empty).
func (r *Registry) ImageSource(id string) (ImageSource, error) {
	if id == "" {
		id = "bing"
	}
	s, ok := r.sources[strings.ToLower(id)]
	if !ok {
		return nil, ErrUnknownSource
	}
	return s, nil
}

// Sources lists the available image-source ids.
func (r *Registry) Sources() []string { return r.order }

// Pixiv exposes the Pixiv client.
func (r *Registry) Pixiv() *pixivClient { return r.pixiv }

// base carries shared dependencies for all sources.
type base struct {
	http *httpclient.Safe
}

const chromeUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

func clampLimit(n, def, max int) int {
	if n <= 0 {
		return def
	}
	if n > max {
		return max
	}
	return n
}
