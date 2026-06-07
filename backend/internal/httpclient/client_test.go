package httpclient

import (
	"net"
	"testing"
)

func TestIsBlocked(t *testing.T) {
	blocked := []string{
		"127.0.0.1",       // loopback
		"10.0.0.5",        // private
		"192.168.1.1",     // private
		"172.16.0.1",      // private
		"169.254.169.254", // cloud metadata
		"0.0.0.0",         // unspecified
		"::1",             // loopback v6
	}
	for _, s := range blocked {
		if !isBlocked(net.ParseIP(s)) {
			t.Errorf("%s seharusnya diblokir (SSRF)", s)
		}
	}

	allowed := []string{"8.8.8.8", "1.1.1.1", "93.184.216.34"}
	for _, s := range allowed {
		if isBlocked(net.ParseIP(s)) {
			t.Errorf("%s seharusnya diizinkan", s)
		}
	}
}

func TestValidateScheme(t *testing.T) {
	if err := validateScheme("http://example.com"); err != nil {
		t.Errorf("http harus valid: %v", err)
	}
	if err := validateScheme("https://example.com"); err != nil {
		t.Errorf("https harus valid: %v", err)
	}
	if err := validateScheme("file:///etc/passwd"); err == nil {
		t.Error("skema file harus ditolak")
	}
	if err := validateScheme("gopher://x"); err == nil {
		t.Error("skema gopher harus ditolak")
	}
}
