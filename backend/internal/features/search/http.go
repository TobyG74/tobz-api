package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// httpGet performs a GET with custom headers and returns the body as a string.
func (b *base) httpGet(ctx context.Context, url string, headers map[string]string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	applyHeaders(req, headers)
	body, _, status, err := b.http.Do(req)
	if err != nil {
		return "", err
	}
	if status < 200 || status >= 400 {
		return "", fmt.Errorf("upstream status %d", status)
	}
	return string(body), nil
}

// httpGetJSON performs a GET and decodes the JSON body into out.
func (b *base) httpGetJSON(ctx context.Context, url string, headers map[string]string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	applyHeaders(req, headers)
	body, _, status, err := b.http.Do(req)
	if err != nil {
		return err
	}
	if status < 200 || status >= 400 {
		return fmt.Errorf("upstream status %d", status)
	}
	return json.Unmarshal(body, out)
}

func applyHeaders(req *http.Request, headers map[string]string) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", chromeUA)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
