package search

import (
	"testing"
	"time"

	"github.com/tobz/tobz-api/internal/httpclient"
)

func newReg() *Registry {
	return NewRegistry(httpclient.NewSafe(5 * time.Second))
}

func TestImageSourceResolution(t *testing.T) {
	reg := newReg()
	for _, id := range []string{"bing", "duckduckgo", "pexels", "pinterest"} {
		s, err := reg.ImageSource(id)
		if err != nil {
			t.Errorf("source %s: %v", id, err)
			continue
		}
		if s.ID() != id {
			t.Errorf("got id %s, want %s", s.ID(), id)
		}
	}
}

func TestImageSourceDefaultAndUnknown(t *testing.T) {
	reg := newReg()
	if s, err := reg.ImageSource(""); err != nil || s.ID() != "bing" {
		t.Errorf("empty source should default to bing, got %v / %v", s, err)
	}
	if _, err := reg.ImageSource("nope"); err == nil {
		t.Error("unknown source must error")
	}
}

func TestSourcesList(t *testing.T) {
	if len(newReg().Sources()) != 4 {
		t.Errorf("expected 4 image sources")
	}
}
