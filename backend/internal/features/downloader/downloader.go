// Package downloader is a multi-platform media downloader: each platform
// implements Downloader and registers with the Registry.
package downloader

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tobz/tobz-api/internal/httpclient"
)

var (
	ErrUnsupportedURL = errors.New("URL tidak didukung oleh downloader manapun")
	ErrNoMedia        = errors.New("tidak ada media yang dapat diunduh")
)

// DownloadItem is a single download option (e.g. HD video, audio, etc.).
type DownloadItem struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	MimeType string `json:"mime_type,omitempty"`
	Quality  string `json:"quality,omitempty"`
}

// MediaResult is the scrape result, uniform across platforms.
type MediaResult struct {
	Platform      string         `json:"platform"`
	PlatformName  string         `json:"platform_name"`
	Downloader    string         `json:"downloader"`
	Title         string         `json:"title,omitempty"`
	AuthorName    string         `json:"author_name,omitempty"`
	Duration      int            `json:"duration,omitempty"`
	Thumbnail     string         `json:"thumbnail,omitempty"`
	DownloadItems []DownloadItem `json:"download_items"`
	Images        []string       `json:"images,omitempty"`
}

// Downloader is the contract for each platform.
type Downloader interface {
	ID() string
	Platform() string
	PlatformName() string
	CanHandle(url string) bool
	Scrape(ctx context.Context, rawURL string) (*MediaResult, error)
}

// PlatformInfo is used by the endpoint that lists supported platforms.
type PlatformInfo struct {
	ID           string `json:"id"`
	Platform     string `json:"platform"`
	PlatformName string `json:"platform_name"`
	Downloader   string `json:"downloader"`
}

// Registry holds all active downloaders.
type Registry struct {
	downloaders []Downloader
}

// NewRegistry builds a registry of all downloaders sharing one anti-SSRF client.
func NewRegistry(client *httpclient.Safe) *Registry {
	base := newBase(client)
	return &Registry{
		downloaders: []Downloader{
			&tiktokDownloader{base},
			&youtubeDownloader{base},
			newSnapsave(base, "instagram", "Instagram", []string{"instagram.com", "instagr.am"}, snapsaveIGEndpoint),
			newSnapsave(base, "facebook", "Facebook", []string{"facebook.com", "fb.watch", "fb.com"}, snapsaveFBEndpoint),
			newSnapsave(base, "twitter", "Twitter/X", []string{"twitter.com", "x.com"}, snapsaveTWEndpoint),
			&douyinDownloader{base},
		},
	}
}

// Resolve picks the first downloader that can handle the URL.
func (r *Registry) Resolve(rawURL string) (Downloader, error) {
	u := strings.TrimSpace(rawURL)
	if u == "" {
		return nil, ErrUnsupportedURL
	}
	for _, d := range r.downloaders {
		if d.CanHandle(u) {
			return d, nil
		}
	}
	return nil, ErrUnsupportedURL
}

// Platforms returns the list of supported platforms.
func (r *Registry) Platforms() []PlatformInfo {
	out := make([]PlatformInfo, 0, len(r.downloaders))
	for _, d := range r.downloaders {
		out = append(out, PlatformInfo{
			ID:           d.ID(),
			Platform:     d.Platform(),
			PlatformName: d.PlatformName(),
			Downloader:   d.ID(),
		})
	}
	return out
}

// base holds the shared dependencies for all downloaders.
type base struct {
	http *httpclient.Safe
}

func newBase(c *httpclient.Safe) *base { return &base{http: c} }

const defaultTimeout = 30 * time.Second
