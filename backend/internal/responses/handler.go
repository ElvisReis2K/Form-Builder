package responses

import (
	"errors"
	"net/http"
	"time"

	"github.com/ElvisReis2K/Form-Builder/backend/internal/auth"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/forms"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/httpx"
)

type Handler struct {
	authService *auth.Service
	service     *Service
}

type submissionRequest struct {
	Answers map[string]any `json:"answers"`
}

type submissionResponse struct {
	ID          string         `json:"id"`
	FormID      string         `json:"formId"`
	Answers     map[string]any `json:"answers"`
	SubmittedAt string         `json:"submittedAt"`
}

type submissionListResponse struct {
	Form      formSummaryResponse  `json:"form"`
	Responses []submissionResponse `json:"responses"`
}

type formSummaryResponse struct {
	ID     string          `json:"id"`
	Title  string          `json:"title"`
	Fields []fieldResponse `json:"fields"`
}

type fieldResponse struct {
	ID       string          `json:"id"`
	Position int             `json:"position"`
	Type     forms.FieldType `json:"type"`
	Label    string          `json:"label"`
	Required bool            `json:"required"`
}

func NewHandler(authService *auth.Service, service *Service) *Handler {
	return &Handler{
		authService: authService,
		service:     service,
	}
}

func (handler *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/forms/{formID}/responses", handler.listResponses)
	mux.HandleFunc("POST /api/public/forms/{slug}/responses", handler.submitResponse)
}

func (handler *Handler) listResponses(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	form, responses, err := handler.service.ListResponses(r.Context(), user.ID, r.PathValue("formID"))
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, submissionListResponse{
		Form:      toFormSummaryResponse(form),
		Responses: toSubmissionResponses(responses),
	})
}

func (handler *Handler) submitResponse(w http.ResponseWriter, r *http.Request) {
	var payload submissionRequest
	if !httpx.DecodeJSON(w, r, &payload) {
		return
	}

	response, err := handler.service.SubmitResponse(r.Context(), r.PathValue("slug"), SubmissionInput{
		Answers: payload.Answers,
	})
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, toSubmissionResponse(response))
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
		httpx.WriteError(w, http.StatusNotFound, "not_found", "response resource not found")
	default:
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
	}
}

func toSubmissionResponses(responses []Response) []submissionResponse {
	items := make([]submissionResponse, 0, len(responses))
	for _, response := range responses {
		items = append(items, toSubmissionResponse(response))
	}

	return items
}

func toSubmissionResponse(response Response) submissionResponse {
	return submissionResponse{
		ID:          response.ID,
		FormID:      response.FormID,
		Answers:     response.Answers,
		SubmittedAt: response.SubmittedAt.UTC().Format(time.RFC3339),
	}
}

func toFormSummaryResponse(form forms.Form) formSummaryResponse {
	fields := make([]fieldResponse, 0, len(form.Fields))
	for _, field := range form.Fields {
		fields = append(fields, fieldResponse{
			ID:       field.ID,
			Position: field.Position,
			Type:     field.Type,
			Label:    field.Label,
			Required: field.Required,
		})
	}

	return formSummaryResponse{
		ID:     form.ID,
		Title:  form.Title,
		Fields: fields,
	}
}
