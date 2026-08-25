package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/ElvisReis2K/Form-Builder/backend/internal/httpx"
)

const SessionCookieName = "form_builder_session"

type Handler struct {
	service      *Service
	cookieSecure bool
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

func NewHandler(service *Service, cookieSecure bool) *Handler {
	return &Handler{
		service:      service,
		cookieSecure: cookieSecure,
	}
}

func (handler *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/register", handler.register)
	mux.HandleFunc("/api/auth/login", handler.login)
	mux.HandleFunc("/api/auth/logout", handler.logout)
	mux.HandleFunc("/api/auth/me", handler.me)
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
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "authentication is required")
		return
	}

	user, err := handler.service.Authenticate(r.Context(), token)
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, authResponse{User: toUserResponse(user)})
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
		httpx.WriteError(w, http.StatusConflict, "email_taken", "email is already registered")
	case errors.Is(err, ErrInvalidCredentials):
		httpx.WriteError(w, http.StatusUnauthorized, "invalid_credentials", "email or password is invalid")
	case errors.Is(err, ErrUnauthenticated):
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "authentication is required")
	default:
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
	}
}

func SessionToken(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return "", false
	}

	return cookie.Value, true
}

func toUserResponse(user User) userResponse {
	return userResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
	}
}
