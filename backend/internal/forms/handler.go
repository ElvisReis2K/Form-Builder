package forms

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ElvisReis2K/Form-Builder/backend/internal/auth"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/httpx"
)

type Handler struct {
	authService *auth.Service
	service     *Service
}

type formRequest struct {
	Title       string         `json:"title"`
	Description *string        `json:"description"`
	Fields      []fieldRequest `json:"fields"`
}

type fieldRequest struct {
	Type        FieldType      `json:"type"`
	Label       string         `json:"label"`
	Required    bool           `json:"required"`
	Placeholder *string        `json:"placeholder"`
	Options     []string       `json:"options"`
	Config      map[string]any `json:"config"`
}

type formListResponse struct {
	Forms []formResponse `json:"forms"`
}

type formResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description *string         `json:"description"`
	Status      FormStatus      `json:"status"`
	PublicSlug  *string         `json:"publicSlug"`
	PublicURL   *string         `json:"publicUrl"`
	PublishedAt *string         `json:"publishedAt"`
	Fields      []fieldResponse `json:"fields"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

type fieldResponse struct {
	ID          string         `json:"id"`
	Position    int            `json:"position"`
	Type        FieldType      `json:"type"`
	Label       string         `json:"label"`
	Required    bool           `json:"required"`
	Placeholder *string        `json:"placeholder"`
	Options     []string       `json:"options"`
	Config      map[string]any `json:"config"`
}

func NewHandler(authService *auth.Service, service *Service) *Handler {
	return &Handler{
		authService: authService,
		service:     service,
	}
}

func (handler *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/forms", handler.forms)
	mux.HandleFunc("/api/forms/", handler.form)
	mux.HandleFunc("/api/public/forms/", handler.publicForm)
}

func (handler *Handler) forms(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		forms, err := handler.service.ListForms(r.Context(), user.ID)
		if err != nil {
			handler.writeServiceError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, formListResponse{Forms: toFormResponses(forms)})
	case http.MethodPost:
		var payload formRequest
		if !httpx.DecodeJSON(w, r, &payload) {
			return
		}

		form, err := handler.service.CreateForm(r.Context(), user.ID, payload.toInput())
		if err != nil {
			handler.writeServiceError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusCreated, toFormResponse(form))
	default:
		w.Header().Set("Allow", "GET, POST")
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method is not allowed")
	}
}

func (handler *Handler) form(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	formID, action, ok := parseFormPath(r.URL.Path)
	if !ok {
		httpx.WriteError(w, http.StatusNotFound, "not_found", "resource not found")
		return
	}

	if action != "" {
		handler.formAction(w, r, user.ID, formID, action)
		return
	}

	switch r.Method {
	case http.MethodGet:
		form, err := handler.service.GetForm(r.Context(), user.ID, formID)
		if err != nil {
			handler.writeServiceError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, toFormResponse(form))
	case http.MethodPut:
		var payload formRequest
		if !httpx.DecodeJSON(w, r, &payload) {
			return
		}

		form, err := handler.service.UpdateForm(r.Context(), user.ID, formID, payload.toInput())
		if err != nil {
			handler.writeServiceError(w, err)
			return
		}

		httpx.WriteJSON(w, http.StatusOK, toFormResponse(form))
	case http.MethodDelete:
		if err := handler.service.DeleteForm(r.Context(), user.ID, formID); err != nil {
			handler.writeServiceError(w, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, PUT, DELETE")
		httpx.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method is not allowed")
	}
}

func (handler *Handler) formAction(w http.ResponseWriter, r *http.Request, ownerID string, formID string, action string) {
	if action != "publish" && action != "unpublish" {
		httpx.WriteError(w, http.StatusNotFound, "not_found", "resource not found")
		return
	}

	if !httpx.RequireMethod(w, r, http.MethodPost) {
		return
	}

	var form Form
	var err error
	switch action {
	case "publish":
		form, err = handler.service.PublishForm(r.Context(), ownerID, formID)
	case "unpublish":
		form, err = handler.service.UnpublishForm(r.Context(), ownerID, formID)
	}

	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toFormResponse(form))
}

func (handler *Handler) publicForm(w http.ResponseWriter, r *http.Request) {
	if !httpx.RequireMethod(w, r, http.MethodGet) {
		return
	}

	slug := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/public/forms/"), "/")
	if slug == "" || strings.Contains(slug, "/") {
		httpx.WriteError(w, http.StatusNotFound, "not_found", "resource not found")
		return
	}

	form, err := handler.service.GetPublishedForm(r.Context(), slug)
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toFormResponse(form))
}

func (handler *Handler) authenticate(w http.ResponseWriter, r *http.Request) (auth.User, bool) {
	token, ok := auth.SessionToken(r)
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "authentication is required")
		return auth.User{}, false
	}

	user, err := handler.authService.Authenticate(r.Context(), token)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "authentication is required")
		return auth.User{}, false
	}

	return user, true
}

func (handler *Handler) writeServiceError(w http.ResponseWriter, err error) {
	var validationError ValidationError
	switch {
	case errors.As(err, &validationError):
		httpx.WriteError(w, http.StatusBadRequest, "validation_error", validationError.Message)
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, "not_found", "form not found")
	default:
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
	}
}

func parseFormPath(path string) (string, string, bool) {
	value := strings.Trim(strings.TrimPrefix(path, "/api/forms/"), "/")
	if value == "" {
		return "", "", false
	}

	parts := strings.Split(value, "/")
	switch len(parts) {
	case 1:
		return parts[0], "", true
	case 2:
		return parts[0], parts[1], true
	default:
		return "", "", false
	}
}

func (request formRequest) toInput() FormInput {
	fields := make([]FieldInput, 0, len(request.Fields))
	for _, field := range request.Fields {
		fields = append(fields, FieldInput{
			Type:        field.Type,
			Label:       field.Label,
			Required:    field.Required,
			Placeholder: field.Placeholder,
			Options:     field.Options,
			Config:      field.Config,
		})
	}

	return FormInput{
		Title:       request.Title,
		Description: request.Description,
		Fields:      fields,
	}
}

func toFormResponses(forms []Form) []formResponse {
	responses := make([]formResponse, 0, len(forms))
	for _, form := range forms {
		responses = append(responses, toFormResponse(form))
	}

	return responses
}

func toFormResponse(form Form) formResponse {
	fields := make([]fieldResponse, 0, len(form.Fields))
	for _, field := range form.Fields {
		fields = append(fields, fieldResponse{
			ID:          field.ID,
			Position:    field.Position,
			Type:        field.Type,
			Label:       field.Label,
			Required:    field.Required,
			Placeholder: field.Placeholder,
			Options:     field.Options,
			Config:      field.Config,
		})
	}

	return formResponse{
		ID:          form.ID,
		Title:       form.Title,
		Description: form.Description,
		Status:      form.Status,
		PublicSlug:  form.PublicSlug,
		PublicURL:   publicURL(form.PublicSlug),
		PublishedAt: optionalTime(form.PublishedAt),
		Fields:      fields,
		CreatedAt:   form.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   form.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func publicURL(slug *string) *string {
	if slug == nil {
		return nil
	}

	value := "/f/" + *slug
	return &value
}

func optionalTime(value *time.Time) *string {
	if value == nil {
		return nil
	}

	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}
