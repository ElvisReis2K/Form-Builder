package forms

import "testing"

func TestNormalizeFormInput(t *testing.T) {
	description := "  Customer feedback  "
	placeholder := "  Your name  "

	input, err := normalizeFormInput(FormInput{
		Title:       "  Feedback form  ",
		Description: &description,
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
