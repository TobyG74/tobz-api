package search

import (
	"context"
	"net/url"
	"regexp"
)

// duckduckgoSource scrapes DuckDuckGo images
type duckduckgoSource struct{ *base }

func (s *duckduckgoSource) ID() string { return "duckduckgo" }

const (
	ddgPageURL = "https://duckduckgo.com/"
	ddgImgURL  = "https://duckduckgo.com/i.js"
)

// vqd tokens may contain digits, dashes, letters and underscores, and appear in
// several forms across DDG responses. 
var vqdPatterns = []*regexp.Regexp{
	regexp.MustCompile(`vqd="([^"]+)"`),
	regexp.MustCompile(`vqd='([^']+)'`),
	regexp.MustCompile(`vqd=([\w-]+)&`),
	regexp.MustCompile(`"vqd":"([^"]+)"`),
	regexp.MustCompile(`vqd=([\w-]+)`),
}

func (s *duckduckgoSource) getVqd(ctx context.Context, query string) (string, error) {
	q := url.Values{"q": {query}, "ia": {"images"}, "iax": {"images"}}
	html, err := s.httpGet(ctx, ddgPageURL+"?"+q.Encode(), map[string]string{
		"Accept-Language": "en-US,en;q=0.9",
	})
	if err != nil {
		return "", err
	}
	for _, re := range vqdPatterns {
		if m := re.FindStringSubmatch(html); m != nil && m[1] != "" {
			return m[1], nil
		}
	}
	return "", ErrNoResults
}

func (s *duckduckgoSource) Search(ctx context.Context, query string, limit int) ([]ImageResult, error) {
	count := clampLimit(limit, 60, 150)
	vqd, err := s.getVqd(ctx, query)
	if err != nil {
		return nil, err
	}

	q := url.Values{
		"q":   {query},
		"vqd": {vqd},
		"o":   {"json"},
		"p":   {"-2"}, // safe search off
		"s":   {"0"},
		"u":   {"bing"},
		"f":   {",,,,"},
		"l":   {"us-en"},
	}
	headers := map[string]string{
		"Referer":          ddgPageURL + "?q=" + url.QueryEscape(query) + "&ia=images&iax=images",
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":  "en-US,en;q=0.9",
		"X-Requested-With": "XMLHttpRequest",
	}

	var resp struct {
		Results []struct {
			Image     string `json:"image"`
			Thumbnail string `json:"thumbnail"`
			Width     int    `json:"width"`
			Height    int    `json:"height"`
			Title     string `json:"title"`
			URL       string `json:"url"`
		} `json:"results"`
	}
	if err := s.httpGetJSON(ctx, ddgImgURL+"?"+q.Encode(), headers, &resp); err != nil {
		return nil, err
	}
	if len(resp.Results) == 0 {
		return nil, ErrNoResults
	}
	if count > len(resp.Results) {
		count = len(resp.Results)
	}
	out := make([]ImageResult, 0, count)
	for _, it := range resp.Results[:count] {
		thumb := it.Thumbnail
		if thumb == "" {
			thumb = it.Image
		}
		out = append(out, ImageResult{
			URL: it.Image, Thumbnail: thumb, Width: it.Width, Height: it.Height,
			Title: it.Title, Source: "duckduckgo", Type: "image",
		})
	}
	return out, nil
}
