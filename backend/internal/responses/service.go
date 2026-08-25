package responses

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/mail"
	"strconv"
	"strings"

	"github.com/ElvisReis2K/Form-Builder/backend/internal/forms"
)

const maxTextAnswerLength = 4000

type Service struct {
	formsRepo *forms.Repository
	repo      *Repository
}

func NewService(formsRepo *forms.Repository, repo *Repository) *Service {
	return &Service{
		formsRepo: formsRepo,
		repo:      repo,
	}
}

func (service *Service) SubmitResponse(ctx context.Context, slug string, input SubmissionInput) (Response, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return Response{}, ValidationError{Message: "form slug is required"}
	}

	form, err := service.formsRepo.GetPublishedBySlug(ctx, slug)
	if errors.Is(err, forms.ErrNotFound) {
		return Response{}, ErrNotFound
	}
	var formValidationError forms.ValidationError
	if errors.As(err, &formValidationError) {
		return Response{}, ValidationError{Message: formValidationError.Message}
	}
	if err != nil {
		return Response{}, err
	}

	answers, err := validateAnswers(form, input.Answers)
	if err != nil {
		return Response{}, err
	}

	return service.repo.Create(ctx, form.ID, answers)
}

func (service *Service) ListResponses(ctx context.Context, ownerID string, formID string) (forms.Form, []Response, error) {
	form, err := service.formsRepo.GetByOwner(ctx, ownerID, formID)
	if errors.Is(err, forms.ErrNotFound) {
		return forms.Form{}, nil, ErrNotFound
	}
	var formValidationError forms.ValidationError
	if errors.As(err, &formValidationError) {
		return forms.Form{}, nil, ValidationError{Message: formValidationError.Message}
	}
	if err != nil {
		return forms.Form{}, nil, err
	}

	responses, err := service.repo.ListByForm(ctx, ownerID, form.ID)
	if err != nil {
		return forms.Form{}, nil, err
	}

	return form, responses, nil
}

func validateAnswers(form forms.Form, answers map[string]any) (map[string]any, error) {
	if answers == nil {
		answers = map[string]any{}
	}

	fieldsByID := map[string]forms.Field{}
	for _, field := range form.Fields {
		fieldsByID[field.ID] = field
	}

	for fieldID := range answers {
		if _, ok := fieldsByID[fieldID]; !ok {
			return nil, ValidationError{Message: "answer field is invalid"}
		}
	}

	normalized := map[string]any{}
	for _, field := range form.Fields {
		value, hasValue := answers[field.ID]
		if !hasValue || isEmptyAnswer(field, value) {
			if field.Required {
				return nil, ValidationError{Message: fmt.Sprintf("%s is required", field.Label)}
			}

			continue
		}

		normalizedValue, err := normalizeAnswer(field, value)
		if err != nil {
			return nil, err
		}

		normalized[field.ID] = normalizedValue
	}

	return normalized, nil
}

func normalizeAnswer(field forms.Field, value any) (any, error) {
	switch field.Type {
	case forms.FieldTypeText, forms.FieldTypeTextarea:
		return normalizeTextAnswer(field, value)
	case forms.FieldTypeEmail:
		answer, err := normalizeTextAnswer(field, value)
		if err != nil {
			return nil, err
		}

		text := answer.(string)
		parsed, err := mail.ParseAddress(text)
		if err != nil || parsed.Address != text {
			return nil, ValidationError{Message: fmt.Sprintf("%s must be a valid email", field.Label)}
		}

		return text, nil
	case forms.FieldTypeNumber:
		return normalizeNumberAnswer(field, value)
	case forms.FieldTypeSelect:
		answer, err := normalizeTextAnswer(field, value)
		if err != nil {
			return nil, err
		}

		text := answer.(string)
		for _, option := range field.Options {
			if option == text {
				return text, nil
			}
		}

		return nil, ValidationError{Message: fmt.Sprintf("%s must be one of the available options", field.Label)}
	case forms.FieldTypeCheckbox:
		answer, ok := value.(bool)
		if !ok {
			return nil, ValidationError{Message: fmt.Sprintf("%s must be true or false", field.Label)}
		}

		if field.Required && !answer {
			return nil, ValidationError{Message: fmt.Sprintf("%s is required", field.Label)}
		}

		return answer, nil
	default:
		return nil, ValidationError{Message: "field type is invalid"}
	}
}

func normalizeTextAnswer(field forms.Field, value any) (any, error) {
	text, ok := value.(string)
	if !ok {
		return nil, ValidationError{Message: fmt.Sprintf("%s must be text", field.Label)}
	}

	text = strings.TrimSpace(text)
	if len(text) > maxTextAnswerLength {
		return nil, ValidationError{Message: fmt.Sprintf("%s is too long", field.Label)}
	}

	return text, nil
}

func normalizeNumberAnswer(field forms.Field, value any) (any, error) {
	switch number := value.(type) {
	case float64:
		if math.IsNaN(number) || math.IsInf(number, 0) {
			return nil, ValidationError{Message: fmt.Sprintf("%s must be a valid number", field.Label)}
		}

		return number, nil
	case string:
		number = strings.TrimSpace(number)
		if number == "" {
			return nil, ValidationError{Message: fmt.Sprintf("%s is required", field.Label)}
		}

		parsed, err := parseNumber(number)
		if err != nil {
			return nil, ValidationError{Message: fmt.Sprintf("%s must be a valid number", field.Label)}
		}
		if math.IsNaN(parsed) || math.IsInf(parsed, 0) {
			return nil, ValidationError{Message: fmt.Sprintf("%s must be a valid number", field.Label)}
		}

		return parsed, nil
	default:
		return nil, ValidationError{Message: fmt.Sprintf("%s must be a valid number", field.Label)}
	}
}

func parseNumber(value string) (float64, error) {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func isEmptyAnswer(field forms.Field, value any) bool {
	if value == nil {
		return true
	}

	switch answer := value.(type) {
	case string:
		return strings.TrimSpace(answer) == ""
	case bool:
		return field.Type == forms.FieldTypeCheckbox && !answer
	default:
		return false
	}
}
