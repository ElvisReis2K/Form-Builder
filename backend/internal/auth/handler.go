package auth

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ElvisReis2K/Form-Builder/backend/internal/httpx"
)

const SessionCookieName = "form_builder_session"
const OAuthStateCookieName = "form_builder_oauth_state"

const oauthStateMaxAgeSeconds = 600

type Handler struct {
	service      *Service
	cookieSecure bool
	frontendURL  string
	googleOAuth  *GoogleOAuth
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type authResponse struct {
	User userResponse `json:"user"`
}

func NewHandler(service *Service, cookieSecure bool, frontendURL string, googleOAuth *GoogleOAuth) *Handler {
	return &Handler{
		service:      service,
		cookieSecure: cookieSecure,
		frontendURL:  frontendURL,
		googleOAuth:  googleOAuth,
	}
}

func (handler *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/register", handler.register)
	mux.HandleFunc("/api/auth/login", handler.login)
	mux.HandleFunc("/api/auth/logout", handler.logout)
	mux.HandleFunc("/api/auth/me", handler.me)
	mux.HandleFunc("GET /api/auth/google", handler.googleStart)
	mux.HandleFunc("GET /api/auth/google/callback", handler.googleCallback)
}

func (handler *Handler) register(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var payload registerRequest
	if !httpx.DecodeJSON(w, r, &payload) {
		return
	}

	result, err := handler.service.Register(r.Context(), RegisterInput{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	handler.setSessionCookie(w, result.Token, result.ExpiresAt)
	httpx.WriteJSON(w, http.StatusCreated, authResponse{User: toUserResponse(result.User)})
}

func (handler *Handler) login(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var payload loginRequest
	if !httpx.DecodeJSON(w, r, &payload) {
		return
	}

	result, err := handler.service.Login(r.Context(), LoginInput{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	handler.setSessionCookie(w, result.Token, result.ExpiresAt)
	httpx.WriteJSON(w, http.StatusOK, authResponse{User: toUserResponse(result.User)})
}

func (handler *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	token, ok := SessionToken(r)
	if !ok {
		handler.clearSessionCookie(w)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := handler.service.Logout(r.Context(), token); err != nil {
		handler.writeServiceError(w, err)
		return
	}

	handler.clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) me(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodGet) {
		return
	}

	token, ok := SessionToken(r)
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "autenticação obrigatória")
		return
	}

	user, err := handler.service.Authenticate(r.Context(), token)
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, authResponse{User: toUserResponse(user)})
}

func (handler *Handler) googleStart(w http.ResponseWriter, r *http.Request) {
	if reason := handler.googleOAuthConfigurationError(); reason != "" {
		http.Redirect(w, r, handler.frontendRedirect("/", "google_oauth_not_configured"), http.StatusFound)
		return
	}

	state, err := newSessionToken()
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "ocorreu um erro inesperado")
		return
	}

	handler.setOAuthStateCookie(w, state)
	http.Redirect(w, r, handler.googleOAuth.AuthCodeURL(state), http.StatusFound)
}

func (handler *Handler) googleCallback(w http.ResponseWriter, r *http.Request) {
	if !handler.googleOAuthConfigured(w) {
		return
	}

	expectedState, ok := OAuthState(r)
	handler.clearOAuthStateCookie(w)
	if !ok || expectedState != r.URL.Query().Get("state") {
		httpx.WriteError(w, http.StatusBadRequest, "invalid_oauth_state", "estado do OAuth inválido")
		return
	}

	if r.URL.Query().Get("error") != "" {
		http.Redirect(w, r, handler.frontendRedirect("/", "google_oauth_denied"), http.StatusFound)
		return
	}

	profile, err := handler.googleOAuth.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Redirect(w, r, handler.frontendRedirect("/", "google_oauth_failed"), http.StatusFound)
		return
	}

	result, err := handler.service.LoginWithGoogle(r.Context(), GoogleIdentityInput{
		Subject: profile.Subject,
		Email:   profile.Email,
		Name:    profile.Name,
	})
	if err != nil {
		http.Redirect(w, r, handler.frontendRedirect("/", "google_oauth_failed"), http.StatusFound)
		return
	}

	handler.setSessionCookie(w, result.Token, result.ExpiresAt)
	http.Redirect(w, r, handler.frontendRedirect("/admin", ""), http.StatusFound)
}

func (handler *Handler) googleOAuthConfigured(w http.ResponseWriter) bool {
	if reason := handler.googleOAuthConfigurationError(); reason != "" {
		httpx.WriteError(w, http.StatusServiceUnavailable, "google_oauth_not_configured", reason)
		return false
	}

	return true
}

func (handler *Handler) googleOAuthConfigurationError() string {
	if handler.googleOAuth == nil {
		return "login com Google não configurado"
	}
	if reason := handler.googleOAuth.ConfigurationError(); reason != "" {
		return reason
	}

	return ""
}

func (handler *Handler) setSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   handler.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (handler *Handler) setOAuthStateCookie(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     OAuthStateCookieName,
		Value:    state,
		Path:     "/api/auth/google",
		MaxAge:   oauthStateMaxAgeSeconds,
		HttpOnly: true,
		Secure:   handler.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (handler *Handler) clearOAuthStateCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     OAuthStateCookieName,
		Value:    "",
		Path:     "/api/auth/google",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   handler.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (handler *Handler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   handler.cookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (handler *Handler) writeServiceError(w http.ResponseWriter, err error) {
	var validationError ValidationError
	switch {
	case errors.As(err, &validationError):
		httpx.WriteError(w, http.StatusBadRequest, "validation_error", validationError.Message)
	case errors.Is(err, ErrEmailTaken):
		httpx.WriteError(w, http.StatusConflict, "email_taken", "e-mail já cadastrado")
	case errors.Is(err, ErrInvalidCredentials):
		httpx.WriteError(w, http.StatusUnauthorized, "invalid_credentials", "e-mail ou senha inválidos")
	case errors.Is(err, ErrUnauthenticated):
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "autenticação obrigatória")
	default:
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "ocorreu um erro inesperado")
	}
}

func SessionToken(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return "", false
	}

	return cookie.Value, true
}

func OAuthState(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(OAuthStateCookieName)
	if err != nil || cookie.Value == "" {
		return "", false
	}

	return cookie.Value, true
}

func (handler *Handler) frontendRedirect(path string, authError string) string {
	base := strings.TrimSpace(handler.frontendURL)
	if base == "" {
		base = "/"
	}

	redirectURL, err := url.Parse(base)
	if err != nil {
		return "/"
	}

	redirectURL.Path = path
	query := redirectURL.Query()
	if authError != "" {
		query.Set("authError", authError)
	}
	redirectURL.RawQuery = query.Encode()

	return redirectURL.String()
}

func toUserResponse(user User) userResponse {
	return userResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
	}
}
