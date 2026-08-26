package forms

import "testing"

func TestPreserveExistingFieldIDsUsesRequestID(t *testing.T) {
	fieldID := "field-id"

	fields := preserveExistingFieldIDs(
		[]FieldInput{
			{ID: &fieldID, Type: FieldTypePhone, Label: "Novo telefone"},
		},
		[]Field{
			{ID: "field-id", Position: 0, Type: FieldTypeText, Label: "Telefone"},
		},
	)

	if fields[0].ID == nil || *fields[0].ID != "field-id" {
		t.Fatalf("expected existing field id to be preserved, got %#v", fields[0].ID)
	}
}

func TestPreserveExistingFieldIDsFallsBackToMatchingPosition(t *testing.T) {
	fields := preserveExistingFieldIDs(
		[]FieldInput{
			{Type: FieldTypeEmail, Label: "E-mail"},
		},
		[]Field{
			{ID: "email-field-id", Position: 0, Type: FieldTypeEmail, Label: "E-mail"},
		},
	)

	if fields[0].ID == nil || *fields[0].ID != "email-field-id" {
		t.Fatalf("expected matching positioned field id to be preserved, got %#v", fields[0].ID)
	}
}

func TestPreserveExistingFieldIDsDoesNotReuseDifferentField(t *testing.T) {
	fields := preserveExistingFieldIDs(
		[]FieldInput{
			{Type: FieldTypePhone, Label: "Telefone"},
		},
		[]Field{
			{ID: "email-field-id", Position: 0, Type: FieldTypeEmail, Label: "E-mail"},
		},
	)

	if fields[0].ID != nil {
		t.Fatalf("expected different field to receive a new id, got %#v", fields[0].ID)
	}
}
