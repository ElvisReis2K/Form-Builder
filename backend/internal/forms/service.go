package forms

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/mail"
	"regexp"
	"strings"
)

const (
	maxTitleLength           = 160
	maxDescriptionLength     = 1000
	maxControllerEmailLength = 160
	maxPrivacyPurposeLength  = 1000
	maxRetentionPolicyLength = 1000
	maxFieldsPerForm         = 50
	maxLabelLength           = 160
	maxPlaceholderLength     = 160
	maxOptionLength          = 120
)

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (service *Service) ListForms(ctx context.Context, ownerID string) ([]Form, error) {
	return service.repo.ListByOwner(ctx, ownerID)
}

func (service *Service) CreateForm(ctx context.Context, ownerID string, input FormInput) (Form, error) {
	normalized, err := normalizeFormInput(input)
	if err != nil {
		return Form{}, err
	}

	return service.repo.Create(ctx, ownerID, normalized)
}

func (service *Service) GetForm(ctx context.Context, ownerID string, formID string) (Form, error) {
	if err := validateID(formID); err != nil {
		return Form{}, err
	}

	return service.repo.GetByOwner(ctx, ownerID, formID)
}

func (service *Service) UpdateForm(ctx context.Context, ownerID string, formID string, input FormInput) (Form, error) {
	if err := validateID(formID); err != nil {
		return Form{}, err
	}

	normalized, err := normalizeFormInput(input)
	if err != nil {
		return Form{}, err
	}

	existing, err := service.repo.GetByOwner(ctx, ownerID, formID)
	if err != nil {
		return Form{}, err
	}
	if existing.Status == FormStatusPublished {
		if err := validatePublishableInput(normalized); err != nil {
			return Form{}, err
		}
	}

	return service.repo.Update(ctx, ownerID, formID, normalized)
}

func (service *Service) DeleteForm(ctx context.Context, ownerID string, formID string) error {
	if err := validateID(formID); err != nil {
		return err
	}

	return service.repo.Delete(ctx, ownerID, formID)
}

func (service *Service) PublishForm(ctx context.Context, ownerID string, formID string) (Form, error) {
	form, err := service.GetForm(ctx, ownerID, formID)
	if err != nil {
		return Form{}, err
	}

	if len(form.Fields) == 0 {
		return Form{}, ValidationError{Message: "o formulario deve ter pelo menos um campo antes da publicacao"}
	}
	if err := validatePrivacyNotice(form.ControllerEmail, form.PrivacyPurpose, form.RetentionPolicy); err != nil {
		return Form{}, err
	}

	for attempt := 0; attempt < 5; attempt++ {
		slug, err := newPublicSlug()
		if err != nil {
			return Form{}, err
		}

		published, err := service.repo.Publish(ctx, ownerID, formID, slug)
		if errors.Is(err, ErrSlugTaken) {
			continue
		}
		if err != nil {
			return Form{}, err
		}

		return published, nil
	}

	return Form{}, ErrSlugTaken
}

func (service *Service) UnpublishForm(ctx context.Context, ownerID string, formID string) (Form, error) {
	if err := validateID(formID); err != nil {
		return Form{}, err
	}

	return service.repo.Unpublish(ctx, ownerID, formID)
}

func (service *Service) GetPublishedForm(ctx context.Context, slug string) (Form, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return Form{}, ValidationError{Message: "slug e obrigatorio"}
	}

	return service.repo.GetPublishedBySlug(ctx, slug)
}

func normalizeFormInput(input FormInput) (FormInput, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return FormInput{}, ValidationError{Message: "titulo e obrigatorio"}
	}
	if len(title) > maxTitleLength {
		return FormInput{}, ValidationError{Message: "titulo muito longo"}
	}

	description := normalizeOptionalString(input.Description)
	if description != nil && len(*description) > maxDescriptionLength {
		return FormInput{}, ValidationError{Message: "descricao muito longa"}
	}

	controllerEmail, err := normalizeControllerEmail(input.ControllerEmail)
	if err != nil {
		return FormInput{}, err
	}

	privacyPurpose := normalizeOptionalString(input.PrivacyPurpose)
	if privacyPurpose != nil && len(*privacyPurpose) > maxPrivacyPurposeLength {
		return FormInput{}, ValidationError{Message: "finalidade do tratamento muito longa"}
	}

	retentionPolicy := normalizeOptionalString(input.RetentionPolicy)
	if retentionPolicy != nil && len(*retentionPolicy) > maxRetentionPolicyLength {
		return FormInput{}, ValidationError{Message: "politica de retencao muito longa"}
	}

	if len(input.Fields) > maxFieldsPerForm {
		return FormInput{}, ValidationError{Message: "formulario tem campos demais"}
	}

	fields := make([]FieldInput, 0, len(input.Fields))
	for _, field := range input.Fields {
		normalized, err := normalizeFieldInput(field)
		if err != nil {
			return FormInput{}, err
		}

		fields = append(fields, normalized)
	}

	return FormInput{
		Title:           title,
		Description:     description,
		ControllerEmail: controllerEmail,
		PrivacyPurpose:  privacyPurpose,
		RetentionPolicy: retentionPolicy,
		Fields:          fields,
	}, nil
}

func validatePublishableInput(input FormInput) error {
	if len(input.Fields) == 0 {
		return ValidationError{Message: "o formulario deve ter pelo menos um campo antes da publicacao"}
	}

	return validatePrivacyNotice(input.ControllerEmail, input.PrivacyPurpose, input.RetentionPolicy)
}

func validatePrivacyNotice(controllerEmail *string, privacyPurpose *string, retentionPolicy *string) error {
	if controllerEmail == nil {
		return ValidationError{Message: "e-mail de contato do controlador e obrigatorio antes da publicacao"}
	}
	if privacyPurpose == nil {
		return ValidationError{Message: "finalidade do tratamento e obrigatoria antes da publicacao"}
	}
	if retentionPolicy == nil {
		return ValidationError{Message: "politica de retencao e obrigatoria antes da publicacao"}
	}

	return nil
}

func normalizeFieldInput(input FieldInput) (FieldInput, error) {
	if !allowedFieldType(input.Type) {
		return FieldInput{}, ValidationError{Message: "tipo de campo invalido"}
	}

	label := strings.TrimSpace(input.Label)
	if label == "" {
		return FieldInput{}, ValidationError{Message: "rotulo do campo e obrigatorio"}
	}
	if len(label) > maxLabelLength {
		return FieldInput{}, ValidationError{Message: "rotulo do campo muito longo"}
	}

	placeholder := normalizeOptionalString(input.Placeholder)
	if placeholder != nil && len(*placeholder) > maxPlaceholderLength {
		return FieldInput{}, ValidationError{Message: "texto de ajuda do campo muito longo"}
	}

	options := normalizeOptions(input.Options)
	if input.Type == FieldTypeSelect && len(options) == 0 {
		return FieldInput{}, ValidationError{Message: "campos de selecao devem ter pelo menos uma opcao"}
	}

	config := input.Config
	if config == nil {
		config = map[string]any{}
	}

	return FieldInput{
		Type:        input.Type,
		Label:       label,
		Required:    input.Required,
		Placeholder: placeholder,
		Options:     options,
		Config:      config,
	}, nil
}

func allowedFieldType(fieldType FieldType) bool {
	switch fieldType {
	case FieldTypeText, FieldTypeTextarea, FieldTypeEmail, FieldTypeNumber, FieldTypeSelect, FieldTypeCheckbox:
		return true
	default:
		return false
	}
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func normalizeControllerEmail(value *string) (*string, error) {
	email := normalizeOptionalString(value)
	if email == nil {
		return nil, nil
	}
	if len(*email) > maxControllerEmailLength {
		return nil, ValidationError{Message: "e-mail de contato do controlador muito longo"}
	}

	parsed, err := mail.ParseAddress(*email)
	if err != nil || parsed.Address != *email {
		return nil, ValidationError{Message: "e-mail de contato do controlador deve ser valido"}
	}

	return email, nil
}

func normalizeOptions(options []string) []string {
	normalized := []string{}
	seen := map[string]bool{}

	for _, option := range options {
		option = strings.TrimSpace(option)
		if option == "" || len(option) > maxOptionLength || seen[option] {
			continue
		}

		seen[option] = true
		normalized = append(normalized, option)
	}

	return normalized
}

func validateID(id string) error {
	if !uuidPattern.MatchString(strings.TrimSpace(id)) {
		return ValidationError{Message: "id do formulario invalido"}
	}

	return nil
}

func newPublicSlug() (string, error) {
	bytes := make([]byte, 9)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
