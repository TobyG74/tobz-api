package search

import (
	"context"
	"encoding/json"
	"net/url"
	"sort"
	"strings"
	"time"
)

// pinterestSource searches Pinterest pins 
type pinterestSource struct{ *base }

func (s *pinterestSource) ID() string { return "pinterest" }

func (s *pinterestSource) Search(ctx context.Context, query string, limit int) ([]ImageResult, error) {
	count := clampLimit(limit, 25, 50)

	data, _ := json.Marshal(map[string]interface{}{
		"options": map[string]interface{}{
			"isPrefetch":                   false,
			"query":                        query,
			"scope":                        "pins",
			"no_fetch_context_on_resource": false,
			"context":                      map[string]interface{}{},
			"bookmarks":                    []string{},
			"page_size":                    count,
		},
	})

	params := url.Values{
		"source_url": {"/search/pins/q=" + query},
		"data":       {string(data)},
		"_":          {timestamp()},
	}
	endpoint := "https://pinterest.com/resource/BaseSearchResource/get/?" + params.Encode()

	body, err := s.httpGet(ctx, endpoint, map[string]string{
		"X-Pinterest-PWS-Handler": "www/search/[scope].js",
		"Accept":                  "application/json, text/javascript, */*, q=0.01",
	})
	if err != nil {
		return nil, err
	}

	var root map[string]interface{}
	if err := json.Unmarshal([]byte(body), &root); err != nil {
		return nil, err
	}
	// resource_response.data.results
	results := arrField(objField(objField(root, "resource_response"), "data"), "results")
	if len(results) == 0 {
		return nil, ErrNoResults
	}

	out := make([]ImageResult, 0, len(results))
	for _, raw := range results {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		mediaURL, kind := resolvePinMedia(item)
		if mediaURL == "" {
			continue
		}
		pinner := objField(item, "pinner")
		caption := strField(item, "grid_title")
		if caption == "" {
			caption = query
		}
		out = append(out, ImageResult{
			URL:    mediaURL,
			Title:  caption,
			Author: strField(pinner, "username"),
			Source: "pinterest",
			Type:   kind,
		})
		if len(out) >= count {
			break
		}
	}
	if len(out) == 0 {
		return nil, ErrNoResults
	}
	return out, nil
}

// resolvePinMedia returns the best media URL and its kind (video|gif|image).
func resolvePinMedia(item map[string]interface{}) (string, string) {
	// Prefer an .mp4 from the video list.
	if videos := objField(item, "videos"); videos != nil {
		if list := objField(videos, "video_list"); list != nil {
			for _, v := range list {
				if vo, ok := v.(map[string]interface{}); ok {
					if u := strField(vo, "url"); strings.HasSuffix(u, ".mp4") {
						return u, "video"
					}
				}
			}
		}
	}

	// Otherwise pick the largest image variant.
	if u := largestImage(objField(item, "images")); u != "" {
		if strings.HasSuffix(u, ".gif") {
			return u, "gif"
		}
		return u, "image"
	}
	return "", ""
}

func largestImage(images map[string]interface{}) string {
	if images == nil {
		return ""
	}
	if orig := objField(images, "orig"); orig != nil {
		if u := strField(orig, "url"); u != "" {
			return u
		}
	}
	type variant struct {
		url  string
		area int
	}
	var vs []variant
	for _, v := range images {
		vo, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		u := strField(vo, "url")
		if !strings.HasPrefix(u, "http") {
			continue
		}
		vs = append(vs, variant{u, numField(vo, "width") * numField(vo, "height")})
	}
	if len(vs) == 0 {
		return ""
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].area > vs[j].area })
	return vs[0].url
}

func timestamp() string {
	return time.Now().UTC().Format("20060102150405")
}
