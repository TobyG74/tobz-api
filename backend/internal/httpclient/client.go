package httpclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

var ErrBlockedHost = errors.New("host tujuan diblokir (kemungkinan SSRF)")

const maxResponseBytes = 10 << 20 // 10 MB

// Safe is an HTTP client with built-in SSRF protection.
type Safe struct {
	client *http.Client
}

// NewSafe builds a Safe client that blocks SSRF targets at dial time.
func NewSafe(timeout time.Duration) *Safe {
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, err
			}
			for _, ip := range ips {
				if isBlocked(ip.IP) {
					return nil, ErrBlockedHost
				}
			}

			return dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].IP.String(), port))
		},
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &Safe{
		client: &http.Client{
			Timeout:   timeout,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return errors.New("terlalu banyak redirect")
				}
				return nil
			},
		},
	}
}

func (s *Safe) Get(ctx context.Context, rawURL string) ([]byte, error) {
	if err := validateScheme(rawURL); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "tobz-api/2.0")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status upstream %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
}

func (s *Safe) Do(req *http.Request) (body []byte, header http.Header, status int, err error) {
	if err = validateScheme(req.URL.String()); err != nil {
		return nil, nil, 0, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes))
	if err != nil {
		return nil, resp.Header, resp.StatusCode, err
	}
	return b, resp.Header, resp.StatusCode, nil
}

func validateScheme(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrBlockedHost
	}
	return nil
}

// isBlocked rejects private, loopback, and cloud-metadata IPs (SSRF defense).
func isBlocked(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
		return true
	}

	if ip.Equal(net.ParseIP("169.254.169.254")) {
		return true
	}
	return false
}
