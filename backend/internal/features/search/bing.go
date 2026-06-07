package search

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// bingSource scrapes Bing's async image endpoint
type bingSource struct{ *base }

func (s *bingSource) ID() string { return "bing" }

const bingAsyncURL = "https://www.bing.com/images/async"

var reBingImg = regexp.MustCompile(`"murl":"(https?:[^"]+)","turl":"(https?:[^"]+)","md5":"[^"]*","shkey":"[^"]*","t":"([^"]*)"`)

func (s *bingSource) Search(ctx context.Context, query string, limit int) ([]ImageResult, error) {
	count := clampLimit(limit, 60, 150)
	q := url.Values{
		"q":     {query},
		"first": {"1"},
		"count": {fmt.Sprint(count)},
		"mkt":   {"en-US"},
		"adlt":  {"Off"},
		"qft":   {""},
	}
	headers := map[string]string{
		"Referer":         "https://www.bing.com/images/search?q=" + url.QueryEscape(query),
		"Accept-Language": "en-US,en;q=0.9",
	}

	html, err := s.httpGet(ctx, bingAsyncURL+"?"+q.Encode(), headers)
	if err != nil {
		return nil, err
	}
	decoded := strings.NewReplacer("&quot;", `"`, "&amp;", "&").Replace(html)

	matches := reBingImg.FindAllStringSubmatch(decoded, count)
	if len(matches) == 0 {
		return nil, ErrNoResults
	}
	out := make([]ImageResult, 0, len(matches))
	for _, m := range matches {
		out = append(out, ImageResult{
			URL:       m[1],
			Thumbnail: m[2],
			Title:     m[3],
			Source:    "bing",
			Type:      "image",
		})
	}
	return out, nil
}
