package downloader

import (
	"testing"
	"time"

	"github.com/tobz/tobz-api/internal/httpclient"
)

func newTestRegistry() *Registry {
	return NewRegistry(httpclient.NewSafe(5 * time.Second))
}

func TestResolvePlatform(t *testing.T) {
	reg := newTestRegistry()
	cases := map[string]string{
		"https://www.tiktok.com/@user/video/123":      "tiktok",
		"https://vt.tiktok.com/abc":                   "tiktok",
		"https://youtu.be/dQw4w9WgXcQ":                "youtube",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ": "youtube",
		"https://www.instagram.com/p/abc/":            "instagram",
		"https://www.facebook.com/reel/123":           "facebook",
		"https://x.com/user/status/123":               "twitter",
		"https://twitter.com/user/status/123":         "twitter",
		"https://v.douyin.com/abc/":                   "douyin",
	}
	for url, wantPlatform := range cases {
		dl, err := reg.Resolve(url)
		if err != nil {
			t.Errorf("%s: tidak ada downloader (%v)", url, err)
			continue
		}
		if dl.Platform() != wantPlatform {
			t.Errorf("%s: platform = %s, mau %s", url, dl.Platform(), wantPlatform)
		}
	}
}

func TestResolveUnsupported(t *testing.T) {
	reg := newTestRegistry()
	if _, err := reg.Resolve("https://example.com/foo"); err == nil {
		t.Error("URL tak dikenal harus mengembalikan error")
	}
	if _, err := reg.Resolve(""); err == nil {
		t.Error("URL kosong harus error")
	}
}

func TestPlatformsList(t *testing.T) {
	reg := newTestRegistry()
	got := reg.Platforms()
	if len(got) < 6 {
		t.Errorf("harus ada >=6 platform, dapat %d", len(got))
	}
}
