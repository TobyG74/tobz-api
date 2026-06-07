package search

import (
	"context"
	"fmt"
	"net/url"
)

// pexelsSource queries the Pexels search endpoint (port of pexels.scraper —
// mirrors the reference, which calls the endpoint without an API key).
type pexelsSource struct{ *base }

func (s *pexelsSource) ID() string { return "pexels" }

func (s *pexelsSource) Search(ctx context.Context, query string, limit int) ([]ImageResult, error) {
	perPage := clampLimit(limit, 60, 80)
	q := url.Values{"query": {query}, "per_page": {fmt.Sprint(perPage)}, "page": {"1"}}

	var resp struct {
		Photos []struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			URL    string `json:"url"`
			Alt    string `json:"alt"`
			Src    struct {
				Original string `json:"original"`
				Large    string `json:"large"`
				Medium   string `json:"medium"`
			} `json:"src"`
		} `json:"photos"`
	}
	if err := s.httpGetJSON(ctx, "https://api.pexels.com/v1/search?"+q.Encode(), nil, &resp); err != nil {
		return nil, err
	}
	if len(resp.Photos) == 0 {
		return nil, ErrNoResults
	}

	out := make([]ImageResult, 0, len(resp.Photos))
	for _, p := range resp.Photos {
		full := p.Src.Large
		if full == "" {
			full = p.Src.Original
		}
		title := p.Alt
		if title == "" {
			title = query
		}
		out = append(out, ImageResult{
			URL: full, Thumbnail: p.Src.Medium, Width: p.Width, Height: p.Height,
			Title: title, Source: "pexels", Type: "image",
		})
	}
	return out, nil
}
