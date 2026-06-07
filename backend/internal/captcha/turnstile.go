package captcha

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"encoding/json"
)

type Verifier struct {
	secret  string
	enabled bool
	client  *http.Client
}

func New(secret string, enabled bool) *Verifier {
	return &Verifier{
		secret:  secret,
		enabled: enabled,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

const verifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type siteVerifyResp struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

func (v *Verifier) Verify(ctx context.Context, token, remoteIP string) error {
	if !v.enabled {
		return nil
	}
	if strings.TrimSpace(token) == "" {
		return ErrCaptchaRequired
	}

	form := url.Values{}
	form.Set("secret", v.secret)
	form.Set("response", token)
	if ip := sanitizeIP(remoteIP); ip != "" {
		form.Set("remoteip", ip)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return ErrCaptchaFailed
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.client.Do(req)
	if err != nil {
		return ErrCaptchaFailed
	}
	defer resp.Body.Close()

	var out siteVerifyResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ErrCaptchaFailed
	}
	if !out.Success {
		return ErrCaptchaFailed
	}
	return nil
}

func sanitizeIP(s string) string {
	if ip := net.ParseIP(strings.TrimSpace(s)); ip != nil {
		return ip.String()
	}
	return ""
}
