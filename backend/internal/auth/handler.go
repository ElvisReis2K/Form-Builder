package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const sessionCookieName = "form_builder_session"

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

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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
	if !requireMethod(w, r, http.MethodPost) {
		return
	}

	var payload registerRequest
	if !decodeJSON(w, r, &payload) {
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
	writeJSON(w, http.StatusCreated, authResponse{User: toUserResponse(result.User)})
}

func (handler *Handler) login(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}

	var payload loginRequest
	if !decodeJSON(w, r, &payload) {
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
	writeJSON(w, http.StatusOK, authResponse{User: toUserResponse(result.User)})
}

func (handler *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}

	token, ok := sessionToken(r)
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
	if !requireMethod(w, r, http.MethodGet) {
		return
	}

	token, ok := sessionToken(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthenticated", "authentication is required")
		return
	}

	user, err := handler.service.Authenticate(r.Context(), token)
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, authResponse{User: toUserResponse(user)})
}

func (handler *Handler) setSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
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
		Name:     sessionCookieName,
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
		writeError(w, http.StatusBadRequest, "validation_error", validationError.Message)
	case errors.Is(err, ErrEmailTaken):
		writeError(w, http.StatusConflict, "email_taken", "email is already registered")
	case errors.Is(err, ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "email or password is invalid")
	case errors.Is(err, ErrUnauthenticated):
		writeError(w, http.StatusUnauthorized, "unauthenticated", "authentication is required")
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
	}
}

func sessionToken(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return "", false
	}

	return cookie.Value, true
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}

	w.Header().Set("Allow", method)
	writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method is not allowed")
	return false
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer func() {
		_ = r.Body.Close()
	}()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must be valid JSON")
		return false
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return false
	}

	return true
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, errorResponse{
		Error: apiError{
			Code:    code,
			Message: message,
		},
	})
}

func toUserResponse(user User) userResponse {
	return userResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
	}
}
