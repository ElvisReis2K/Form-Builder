package auth

import (
	"net/url"
	"strings"
	"testing"
)

func TestGoogleOAuthAuthCodeURL(t *testing.T) {
	oauth := NewGoogleOAuth(GoogleOAuthConfig{
		ClientID:     "123456789-clientid.apps.googleusercontent.com",
		ClientSecret: "client-secret",
		RedirectURL:  "http://localhost:8080/api/auth/google/callback",
	})

	authURL := oauth.AuthCodeURL("state-token")
	if !strings.HasPrefix(authURL, googleAuthorizationEndpoint+"?") {
		t.Fatalf("expected google authorization endpoint, got %q", authURL)
	}

	parsed, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("parse auth url: %v", err)
	}

	query := parsed.Query()
	assertQueryValue(t, query, "client_id", "123456789-clientid.apps.googleusercontent.com")
	assertQueryValue(t, query, "redirect_uri", "http://localhost:8080/api/auth/google/callback")
	assertQueryValue(t, query, "response_type", "code")
	assertQueryValue(t, query, "scope", googleOAuthScopes)
	assertQueryValue(t, query, "state", "state-token")
	assertQueryValue(t, query, "prompt", "select_account")
}

func TestGoogleOAuthConfiguredRequiresValidClientID(t *testing.T) {
	oauth := NewGoogleOAuth(GoogleOAuthConfig{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "http://localhost:8080/api/auth/google/callback",
	})

	if oauth.Configured() {
		t.Fatal("expected oauth config to be invalid")
	}

	if !strings.Contains(oauth.ConfigurationError(), "GOOGLE_CLIENT_ID inválido") {
		t.Fatalf("expected invalid client id error, got %q", oauth.ConfigurationError())
	}
}

func TestGoogleOAuthConfiguredRequiresClientSecret(t *testing.T) {
	oauth := NewGoogleOAuth(GoogleOAuthConfig{
		ClientID:    "123456789-clientid.apps.googleusercontent.com",
		RedirectURL: "http://localhost:8080/api/auth/google/callback",
	})

	if oauth.Configured() {
		t.Fatal("expected oauth config to be invalid")
	}

	if !strings.Contains(oauth.ConfigurationError(), "GOOGLE_CLIENT_SECRET não configurado") {
		t.Fatalf("expected missing client secret error, got %q", oauth.ConfigurationError())
	}
}

func assertQueryValue(t *testing.T, query url.Values, key string, expected string) {
	t.Helper()

	if got := query.Get(key); got != expected {
		t.Fatalf("expected %s=%q, got %q", key, expected, got)
	}
}
