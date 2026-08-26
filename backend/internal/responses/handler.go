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
	Answers             map[string]any `json:"answers"`
	PrivacyAcknowledged bool           `json:"privacyAcknowledged"`
}

type submissionResponse struct {
	ID                    string         `json:"id"`
	FormID                string         `json:"formId"`
	Answers               map[string]any `json:"answers"`
	PrivacyAcknowledgedAt string         `json:"privacyAcknowledgedAt"`
	SubmittedAt           string         `json:"submittedAt"`
}

type submissionListResponse struct {
	Form      formSummaryResponse  `json:"form"`
	Responses []submissionResponse `json:"responses"`
}

type submissionExportResponse struct {
	Form       formSummaryResponse  `json:"form"`
	Responses  []submissionResponse `json:"responses"`
	ExportedAt string               `json:"exportedAt"`
}

type formSummaryResponse struct {
	ID              string          `json:"id"`
	Title           string          `json:"title"`
	ControllerEmail *string         `json:"controllerEmail"`
	PrivacyPurpose  *string         `json:"privacyPurpose"`
	RetentionPolicy *string         `json:"retentionPolicy"`
	Fields          []fieldResponse `json:"fields"`
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
	mux.HandleFunc("GET /api/forms/{formID}/responses/export", handler.exportResponses)
	mux.HandleFunc("DELETE /api/forms/{formID}/responses/{responseID}", handler.deleteResponse)
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

func (handler *Handler) exportResponses(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	form, responses, err := handler.service.ListResponses(r.Context(), user.ID, r.PathValue("formID"))
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, submissionExportResponse{
		Form:       toFormSummaryResponse(form),
		Responses:  toSubmissionResponses(responses),
		ExportedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

func (handler *Handler) deleteResponse(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	if err := handler.service.DeleteResponse(r.Context(), user.ID, r.PathValue("formID"), r.PathValue("responseID")); err != nil {
		handler.writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) submitResponse(w http.ResponseWriter, r *http.Request) {
	var payload submissionRequest
	if !httpx.DecodeJSON(w, r, &payload) {
		return
	}

	response, err := handler.service.SubmitResponse(r.Context(), r.PathValue("slug"), SubmissionInput{
		Answers:             payload.Answers,
		PrivacyAcknowledged: payload.PrivacyAcknowledged,
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
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "autenticação obrigatória")
		return auth.User{}, false
	}

	user, err := handler.authService.Authenticate(r.Context(), token)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "autenticação obrigatória")
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
		httpx.WriteError(w, http.StatusNotFound, "not_found", "recurso de resposta não encontrado")
	default:
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "ocorreu um erro inesperado")
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
		ID:                    response.ID,
		FormID:                response.FormID,
		Answers:               response.Answers,
		PrivacyAcknowledgedAt: response.PrivacyAcknowledgedAt.UTC().Format(time.RFC3339),
		SubmittedAt:           response.SubmittedAt.UTC().Format(time.RFC3339),
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
		ID:              form.ID,
		Title:           form.Title,
		ControllerEmail: form.ControllerEmail,
		PrivacyPurpose:  form.PrivacyPurpose,
		RetentionPolicy: form.RetentionPolicy,
		Fields:          fields,
	}
}
