package forms

import (
	"errors"
	"net/http"
	"time"

	"github.com/ElvisReis2K/Form-Builder/backend/internal/auth"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/httpx"
)

type Handler struct {
	authService *auth.Service
	service     *Service
}

type formRequest struct {
	Title           string         `json:"title"`
	Description     *string        `json:"description"`
	ControllerEmail *string        `json:"controllerEmail"`
	PrivacyPurpose  *string        `json:"privacyPurpose"`
	RetentionPolicy *string        `json:"retentionPolicy"`
	Fields          []fieldRequest `json:"fields"`
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
	ID              string          `json:"id"`
	Title           string          `json:"title"`
	Description     *string         `json:"description"`
	ControllerEmail *string         `json:"controllerEmail"`
	PrivacyPurpose  *string         `json:"privacyPurpose"`
	RetentionPolicy *string         `json:"retentionPolicy"`
	Status          FormStatus      `json:"status"`
	PublicSlug      *string         `json:"publicSlug"`
	PublicURL       *string         `json:"publicUrl"`
	PublishedAt     *string         `json:"publishedAt"`
	Fields          []fieldResponse `json:"fields"`
	CreatedAt       string          `json:"createdAt"`
	UpdatedAt       string          `json:"updatedAt"`
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
	mux.HandleFunc("GET /api/forms", handler.listForms)
	mux.HandleFunc("POST /api/forms", handler.createForm)
	mux.HandleFunc("GET /api/forms/{formID}", handler.getForm)
	mux.HandleFunc("PUT /api/forms/{formID}", handler.updateForm)
	mux.HandleFunc("DELETE /api/forms/{formID}", handler.deleteForm)
	mux.HandleFunc("POST /api/forms/{formID}/publish", handler.publishForm)
	mux.HandleFunc("POST /api/forms/{formID}/unpublish", handler.unpublishForm)
	mux.HandleFunc("GET /api/public/forms/{slug}", handler.publicForm)
}

func (handler *Handler) listForms(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	forms, err := handler.service.ListForms(r.Context(), user.ID)
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, formListResponse{Forms: toFormResponses(forms)})
}

func (handler *Handler) createForm(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

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
}

func (handler *Handler) getForm(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	form, err := handler.service.GetForm(r.Context(), user.ID, r.PathValue("formID"))
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toFormResponse(form))
}

func (handler *Handler) updateForm(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	var payload formRequest
	if !httpx.DecodeJSON(w, r, &payload) {
		return
	}

	form, err := handler.service.UpdateForm(r.Context(), user.ID, r.PathValue("formID"), payload.toInput())
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toFormResponse(form))
}

func (handler *Handler) deleteForm(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	if err := handler.service.DeleteForm(r.Context(), user.ID, r.PathValue("formID")); err != nil {
		handler.writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (handler *Handler) publishForm(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	form, err := handler.service.PublishForm(r.Context(), user.ID, r.PathValue("formID"))
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toFormResponse(form))
}

func (handler *Handler) unpublishForm(w http.ResponseWriter, r *http.Request) {
	user, ok := handler.authenticate(w, r)
	if !ok {
		return
	}

	form, err := handler.service.UnpublishForm(r.Context(), user.ID, r.PathValue("formID"))
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toFormResponse(form))
}

func (handler *Handler) publicForm(w http.ResponseWriter, r *http.Request) {
	form, err := handler.service.GetPublishedForm(r.Context(), r.PathValue("slug"))
	if err != nil {
		handler.writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toFormResponse(form))
}

func (handler *Handler) authenticate(w http.ResponseWriter, r *http.Request) (auth.User, bool) {
	token, ok := auth.SessionToken(r)
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "autenticacao obrigatoria")
		return auth.User{}, false
	}

	user, err := handler.authService.Authenticate(r.Context(), token)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated", "autenticacao obrigatoria")
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
		httpx.WriteError(w, http.StatusNotFound, "not_found", "formulario nao encontrado")
	default:
		httpx.WriteError(w, http.StatusInternalServerError, "internal_error", "ocorreu um erro inesperado")
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
		Title:           request.Title,
		Description:     request.Description,
		ControllerEmail: request.ControllerEmail,
		PrivacyPurpose:  request.PrivacyPurpose,
		RetentionPolicy: request.RetentionPolicy,
		Fields:          fields,
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
		ID:              form.ID,
		Title:           form.Title,
		Description:     form.Description,
		ControllerEmail: form.ControllerEmail,
		PrivacyPurpose:  form.PrivacyPurpose,
		RetentionPolicy: form.RetentionPolicy,
		Status:          form.Status,
		PublicSlug:      form.PublicSlug,
		PublicURL:       publicURL(form.PublicSlug),
		PublishedAt:     optionalTime(form.PublishedAt),
		Fields:          fields,
		CreatedAt:       form.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       form.UpdatedAt.UTC().Format(time.RFC3339),
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
