package downloader

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const chromeUA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36"

func (b *base) httpGet(ctx context.Context, rawURL string, headers map[string]string) (string, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", nil, err
	}
	applyHeaders(req, headers)
	body, hdr, status, err := b.http.Do(req)
	if err != nil {
		return "", hdr, err
	}
	if status < 200 || status >= 400 {
		return "", hdr, fmt.Errorf("status upstream %d", status)
	}
	return string(body), hdr, nil
}

func (b *base) httpPostForm(ctx context.Context, rawURL string, form url.Values, headers map[string]string) (string, http.Header, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	applyHeaders(req, headers)
	body, hdr, status, err := b.http.Do(req)
	if err != nil {
		return "", hdr, err
	}
	if status < 200 || status >= 400 {
		return "", hdr, fmt.Errorf("status upstream %d", status)
	}
	return string(body), hdr, nil
}

func applyHeaders(req *http.Request, headers map[string]string) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", chromeUA)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
