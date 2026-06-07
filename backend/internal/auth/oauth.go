package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tobz/tobz-api/internal/config"
)

var ErrOAuthUnsupported = errors.New("provider OAuth tidak didukung")

type OAuthUserInfo struct {
	Provider       string
	ProviderUserID string
	Email          string
	EmailVerified  bool
	DisplayName    string
	AvatarURL      string
}

type OAuthProvider struct {
	Name         string
	ClientID     string
	clientSecret string
	authURL      string
	tokenURL     string
	scopes       []string
	redirectURL  string
	fetchUser    func(ctx context.Context, accessToken string) (*OAuthUserInfo, error)
	client       *http.Client
}

type OAuthRegistry struct {
	providers map[string]*OAuthProvider
}

func NewOAuthRegistry(cfg *config.Config) *OAuthRegistry {
	r := &OAuthRegistry{providers: map[string]*OAuthProvider{}}
	client := &http.Client{Timeout: 10 * time.Second}

	if cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
		r.providers["google"] = &OAuthProvider{
			Name:         "google",
			ClientID:     cfg.GoogleClientID,
			clientSecret: cfg.GoogleClientSecret,
			authURL:      "https://accounts.google.com/o/oauth2/v2/auth",
			tokenURL:     "https://oauth2.googleapis.com/token",
			scopes:       []string{"openid", "email", "profile"},
			redirectURL:  cfg.BaseURL + "/api/v1/auth/oauth/google/callback",
			fetchUser:    fetchGoogleUser,
			client:       client,
		}
	}
	if cfg.GitHubClientID != "" && cfg.GitHubClientSecret != "" {
		r.providers["github"] = &OAuthProvider{
			Name:         "github",
			ClientID:     cfg.GitHubClientID,
			clientSecret: cfg.GitHubClientSecret,
			authURL:      "https://github.com/login/oauth/authorize",
			tokenURL:     "https://github.com/login/oauth/access_token",
			scopes:       []string{"read:user", "user:email"},
			redirectURL:  cfg.BaseURL + "/api/v1/auth/oauth/github/callback",
			fetchUser:    fetchGitHubUser,
			client:       client,
		}
	}
	return r
}

func (r *OAuthRegistry) Get(name string) (*OAuthProvider, error) {
	p, ok := r.providers[name]
	if !ok {
		return nil, ErrOAuthUnsupported
	}
	return p, nil
}

func (p *OAuthProvider) AuthCodeURL(state string) string {
	q := url.Values{}
	q.Set("client_id", p.ClientID)
	q.Set("redirect_uri", p.redirectURL)
	q.Set("response_type", "code")
	q.Set("scope", strings.Join(p.scopes, " "))
	q.Set("state", state)
	return p.authURL + "?" + q.Encode()
}

func (p *OAuthProvider) Exchange(ctx context.Context, code string) (*OAuthUserInfo, error) {
	form := url.Values{}
	form.Set("client_id", p.ClientID)
	form.Set("client_secret", p.clientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", p.redirectURL)
	form.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange gagal: status %d", resp.StatusCode)
	}

	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return nil, err
	}
	if tok.AccessToken == "" {
		return nil, errors.New("access token kosong dari provider")
	}
	return p.fetchUser(ctx, tok.AccessToken)
}

func fetchGoogleUser(ctx context.Context, accessToken string) (*OAuthUserInfo, error) {
	var body struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := getJSON(ctx, "https://openidconnect.googleapis.com/v1/userinfo", accessToken, &body); err != nil {
		return nil, err
	}
	if body.Sub == "" || body.Email == "" {
		return nil, errors.New("profil google tidak lengkap")
	}
	return &OAuthUserInfo{
		Provider:       "google",
		ProviderUserID: body.Sub,
		Email:          strings.ToLower(body.Email),
		EmailVerified:  body.EmailVerified,
		DisplayName:    body.Name,
		AvatarURL:      body.Picture,
	}, nil
}

func fetchGitHubUser(ctx context.Context, accessToken string) (*OAuthUserInfo, error) {
	var profile struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	}
	if err := getJSON(ctx, "https://api.github.com/user", accessToken, &profile); err != nil {
		return nil, err
	}
	if profile.ID == 0 {
		return nil, errors.New("profil github tidak lengkap")
	}

	email := profile.Email
	verified := false
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := getJSON(ctx, "https://api.github.com/user/emails", accessToken, &emails); err == nil {
		for _, e := range emails {
			if e.Primary {
				email = e.Email
				verified = e.Verified
				break
			}
		}
	}
	name := profile.Name
	if name == "" {
		name = profile.Login
	}
	return &OAuthUserInfo{
		Provider:       "github",
		ProviderUserID: strconv.FormatInt(profile.ID, 10),
		Email:          strings.ToLower(email),
		EmailVerified:  verified,
		DisplayName:    name,
		AvatarURL:      profile.AvatarURL,
	}, nil
}

func getJSON(ctx context.Context, urlStr, accessToken string, out interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "tobz-api")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gagal ambil profil: status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
