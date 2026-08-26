package forms

import "testing"

func TestNormalizeFormInput(t *testing.T) {
	description := "  Customer feedback  "
	controllerEmail := "  privacy@example.com  "
	privacyPurpose := "  Answer customer feedback  "
	retentionPolicy := "  Responses are retained for 90 days  "
	placeholder := "  Your name  "

	input, err := normalizeFormInput(FormInput{
		Title:           "  Feedback form  ",
		Description:     &description,
		ControllerEmail: &controllerEmail,
		PrivacyPurpose:  &privacyPurpose,
		RetentionPolicy: &retentionPolicy,
		Fields: []FieldInput{
			{
				Type:        FieldTypeText,
				Label:       "  Name  ",
				Required:    true,
				Placeholder: &placeholder,
			},
			{
				Type:    FieldTypeSelect,
				Label:   "Plan",
				Options: []string{" Basic ", "Pro", "Basic", ""},
			},
		},
	})
	if err != nil {
		t.Fatalf("normalize form input: %v", err)
	}

	if input.Title != "Feedback form" {
		t.Fatalf("expected title to be trimmed, got %q", input.Title)
	}
	if input.Description == nil || *input.Description != "Customer feedback" {
		t.Fatalf("expected description to be normalized, got %#v", input.Description)
	}
	if input.ControllerEmail == nil || *input.ControllerEmail != "privacy@example.com" {
		t.Fatalf("expected controller email to be normalized, got %#v", input.ControllerEmail)
	}
	if input.PrivacyPurpose == nil || *input.PrivacyPurpose != "Answer customer feedback" {
		t.Fatalf("expected privacy purpose to be normalized, got %#v", input.PrivacyPurpose)
	}
	if input.RetentionPolicy == nil || *input.RetentionPolicy != "Responses are retained for 90 days" {
		t.Fatalf("expected retention policy to be normalized, got %#v", input.RetentionPolicy)
	}
	if input.Fields[0].Placeholder == nil || *input.Fields[0].Placeholder != "Your name" {
		t.Fatalf("expected placeholder to be normalized, got %#v", input.Fields[0].Placeholder)
	}
	if got := input.Fields[1].Options; len(got) != 2 || got[0] != "Basic" || got[1] != "Pro" {
		t.Fatalf("expected duplicate and blank options to be removed, got %#v", got)
	}
}

func TestNormalizeFormInputRequiresSelectOptions(t *testing.T) {
	_, err := normalizeFormInput(FormInput{
		Title: "Feedback form",
		Fields: []FieldInput{
			{
				Type:  FieldTypeSelect,
				Label: "Plan",
			},
		},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNormalizeFormInputNormalizesFieldID(t *testing.T) {
	fieldID := " 11111111-1111-4111-8111-111111111111 "

	input, err := normalizeFormInput(FormInput{
		Title: "Feedback form",
		Fields: []FieldInput{
			{
				ID:    &fieldID,
				Type:  FieldTypeText,
				Label: "Name",
			},
		},
	})
	if err != nil {
		t.Fatalf("normalize form input: %v", err)
	}

	if input.Fields[0].ID == nil || *input.Fields[0].ID != "11111111-1111-4111-8111-111111111111" {
		t.Fatalf("expected field id to be normalized, got %#v", input.Fields[0].ID)
	}
}

func TestNormalizeFormInputRejectsInvalidFieldID(t *testing.T) {
	fieldID := "not-a-uuid"

	_, err := normalizeFormInput(FormInput{
		Title: "Feedback form",
		Fields: []FieldInput{
			{
				ID:    &fieldID,
				Type:  FieldTypeText,
				Label: "Name",
			},
		},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNormalizeFormInputRejectsInvalidControllerEmail(t *testing.T) {
	controllerEmail := "not-an-email"

	_, err := normalizeFormInput(FormInput{
		Title:           "Feedback form",
		ControllerEmail: &controllerEmail,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidatePublishableInputRequiresPrivacyNotice(t *testing.T) {
	controllerEmail := "privacy@example.com"
	privacyPurpose := "Answer customer feedback"

	err := validatePublishableInput(FormInput{
		ControllerEmail: &controllerEmail,
		PrivacyPurpose:  &privacyPurpose,
		Fields: []FieldInput{
			{Type: FieldTypeText, Label: "Name"},
		},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
