package forms

import "time"

type FormStatus string

const (
	FormStatusDraft     FormStatus = "draft"
	FormStatusPublished FormStatus = "published"
)

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeTextarea FieldType = "textarea"
	FieldTypeEmail    FieldType = "email"
	FieldTypeNumber   FieldType = "number"
	FieldTypePhone    FieldType = "phone"
	FieldTypeSelect   FieldType = "select"
	FieldTypeCheckbox FieldType = "checkbox"
)

type Form struct {
	ID              string
	OwnerID         string
	Title           string
	Description     *string
	ControllerEmail *string
	PrivacyPurpose  *string
	RetentionPolicy *string
	Status          FormStatus
	PublicSlug      *string
	PublishedAt     *time.Time
	Fields          []Field
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Field struct {
	ID          string
	FormID      string
	Position    int
	Type        FieldType
	Label       string
	Required    bool
	Placeholder *string
	Options     []string
	Config      map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FormInput struct {
	Title           string
	Description     *string
	ControllerEmail *string
	PrivacyPurpose  *string
	RetentionPolicy *string
	Fields          []FieldInput
}

type FieldInput struct {
	Type        FieldType
	Label       string
	Required    bool
	Placeholder *string
	Options     []string
	Config      map[string]any
}
