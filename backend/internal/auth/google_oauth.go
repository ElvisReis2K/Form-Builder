package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	googleAuthorizationEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenEndpoint         = "https://oauth2.googleapis.com/token"
	googleUserInfoEndpoint      = "https://openidconnect.googleapis.com/v1/userinfo"
	googleOAuthScopes           = "openid email profile"
)

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type GoogleOAuth struct {
	config GoogleOAuthConfig
	client *http.Client
}

type GoogleProfile struct {
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

type googleUserInfoResponse struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
}

func NewGoogleOAuth(config GoogleOAuthConfig) *GoogleOAuth {
	return &GoogleOAuth{
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (oauth *GoogleOAuth) Configured() bool {
	return strings.TrimSpace(oauth.config.ClientID) != "" &&
		strings.TrimSpace(oauth.config.ClientSecret) != "" &&
		strings.TrimSpace(oauth.config.RedirectURL) != ""
}

func (oauth *GoogleOAuth) AuthCodeURL(state string) string {
	values := url.Values{}
	values.Set("client_id", oauth.config.ClientID)
	values.Set("redirect_uri", oauth.config.RedirectURL)
	values.Set("response_type", "code")
	values.Set("scope", googleOAuthScopes)
	values.Set("state", state)
	values.Set("include_granted_scopes", "true")
	values.Set("prompt", "select_account")

	return googleAuthorizationEndpoint + "?" + values.Encode()
}

func (oauth *GoogleOAuth) Exchange(ctx context.Context, code string) (GoogleProfile, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return GoogleProfile{}, ValidationError{Message: "authorization code is required"}
	}

	token, err := oauth.exchangeCode(ctx, code)
	if err != nil {
		return GoogleProfile{}, err
	}

	return oauth.fetchProfile(ctx, token.AccessToken)
}

func (oauth *GoogleOAuth) exchangeCode(ctx context.Context, code string) (googleTokenResponse, error) {
	body := url.Values{}
	body.Set("client_id", oauth.config.ClientID)
	body.Set("client_secret", oauth.config.ClientSecret)
	body.Set("code", code)
	body.Set("grant_type", "authorization_code")
	body.Set("redirect_uri", oauth.config.RedirectURL)

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, googleTokenEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return googleTokenResponse{}, fmt.Errorf("create google token request: %w", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := oauth.client.Do(request)
	if err != nil {
		return googleTokenResponse{}, fmt.Errorf("exchange google authorization code: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	payload, err := decodeGoogleJSON[googleTokenResponse](response.Body)
	if err != nil {
		return googleTokenResponse{}, err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return googleTokenResponse{}, fmt.Errorf("google token endpoint failed: %s", payload.Error)
	}
	if strings.TrimSpace(payload.AccessToken) == "" {
		return googleTokenResponse{}, fmt.Errorf("google token response did not include access token")
	}

	return payload, nil
}

func (oauth *GoogleOAuth) fetchProfile(ctx context.Context, accessToken string) (GoogleProfile, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserInfoEndpoint, nil)
	if err != nil {
		return GoogleProfile{}, fmt.Errorf("create google userinfo request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := oauth.client.Do(request)
	if err != nil {
		return GoogleProfile{}, fmt.Errorf("fetch google profile: %w", err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	payload, err := decodeGoogleJSON[googleUserInfoResponse](response.Body)
	if err != nil {
		return GoogleProfile{}, err
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return GoogleProfile{}, fmt.Errorf("google userinfo endpoint failed")
	}

	profile := GoogleProfile{
		Subject:       strings.TrimSpace(payload.Subject),
		Email:         normalizeEmail(payload.Email),
		EmailVerified: payload.EmailVerified,
		Name:          strings.TrimSpace(payload.Name),
	}
	if profile.Subject == "" {
		return GoogleProfile{}, ValidationError{Message: "google profile subject is required"}
	}
	if err := validateEmail(profile.Email); err != nil {
		return GoogleProfile{}, ValidationError{Message: "google profile email is required"}
	}
	if !profile.EmailVerified {
		return GoogleProfile{}, ValidationError{Message: "google profile email must be verified"}
	}

	return profile, nil
}

func decodeGoogleJSON[T any](body io.Reader) (T, error) {
	var payload T
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return payload, fmt.Errorf("decode google response: %w", err)
	}

	return payload, nil
}
