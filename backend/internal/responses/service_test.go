package responses

import (
	"testing"

	"github.com/ElvisReis2K/Form-Builder/backend/internal/forms"
)

func TestValidateAnswersNormalizesValidSubmission(t *testing.T) {
	form := testForm()

	answers, err := validateAnswers(form, map[string]any{
		"text-field":   "  Ada  ",
		"email-field":  "ada@example.com",
		"number-field": "42.5",
		"phone-field":  "(123) 999999999",
		"select-field": "Pro",
		"check-field":  true,
	})
	if err != nil {
		t.Fatalf("validate answers: %v", err)
	}

	if answers["text-field"] != "Ada" {
		t.Fatalf("expected trimmed text answer, got %#v", answers["text-field"])
	}
	if answers["number-field"] != 42.5 {
		t.Fatalf("expected parsed number answer, got %#v", answers["number-field"])
	}
	if answers["phone-field"] != "123999999999" {
		t.Fatalf("expected normalized phone answer, got %#v", answers["phone-field"])
	}
}

func TestValidateAnswersRejectsUnknownField(t *testing.T) {
	_, err := validateAnswers(testForm(), map[string]any{
		"unknown-field": "value",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateAnswersRequiresMandatoryFields(t *testing.T) {
	_, err := validateAnswers(testForm(), map[string]any{
		"text-field": "",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateAnswersRejectsInvalidSelectOption(t *testing.T) {
	_, err := validateAnswers(testForm(), map[string]any{
		"text-field":   "Ada",
		"email-field":  "ada@example.com",
		"number-field": float64(42),
		"select-field": "Enterprise",
		"check-field":  true,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateAnswersRejectsInvalidPhone(t *testing.T) {
	_, err := validateAnswers(testForm(), map[string]any{
		"text-field":   "Ada",
		"email-field":  "ada@example.com",
		"number-field": float64(42),
		"phone-field":  "11999999999",
		"select-field": "Pro",
		"check-field":  true,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidatePrivacyAcknowledgementRequiresConfirmation(t *testing.T) {
	if err := validatePrivacyAcknowledgement(false); err == nil {
		t.Fatal("expected validation error")
	}

	if err := validatePrivacyAcknowledgement(true); err != nil {
		t.Fatalf("expected valid acknowledgement, got %v", err)
	}
}

func testForm() forms.Form {
	return forms.Form{
		ID: "form-id",
		Fields: []forms.Field{
			{ID: "text-field", Type: forms.FieldTypeText, Label: "Name", Required: true},
			{ID: "email-field", Type: forms.FieldTypeEmail, Label: "Email", Required: true},
			{ID: "number-field", Type: forms.FieldTypeNumber, Label: "Score", Required: true},
			{ID: "phone-field", Type: forms.FieldTypePhone, Label: "Phone", Required: true},
			{ID: "select-field", Type: forms.FieldTypeSelect, Label: "Plan", Required: true, Options: []string{"Basic", "Pro"}},
			{ID: "check-field", Type: forms.FieldTypeCheckbox, Label: "Accept terms", Required: true},
		},
	}
}
