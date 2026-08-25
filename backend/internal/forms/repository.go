package forms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) ListByOwner(ctx context.Context, ownerID string) ([]Form, error) {
	rows, err := repo.db.Query(ctx, `
		SELECT id::text, owner_id::text, title, description, controller_email, privacy_purpose, retention_policy, status, public_slug, published_at, created_at, updated_at
		FROM forms
		WHERE owner_id = $1
		ORDER BY updated_at DESC
	`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list forms: %w", err)
	}
	defer rows.Close()

	forms, err := scanFormRows(rows)
	if err != nil {
		return nil, err
	}

	for index := range forms {
		fields, err := repo.ListFields(ctx, forms[index].ID)
		if err != nil {
			return nil, err
		}

		forms[index].Fields = fields
	}

	return forms, nil
}

func (repo *Repository) Create(ctx context.Context, ownerID string, input FormInput) (Form, error) {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return Form{}, fmt.Errorf("begin create form: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	form, err := insertForm(ctx, tx, ownerID, input)
	if err != nil {
		return Form{}, err
	}

	fields, err := replaceFields(ctx, tx, form.ID, input.Fields)
	if err != nil {
		return Form{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Form{}, fmt.Errorf("commit create form: %w", err)
	}

	form.Fields = fields
	return form, nil
}

func (repo *Repository) GetByOwner(ctx context.Context, ownerID string, formID string) (Form, error) {
	form, err := repo.findForm(ctx, `
		SELECT id::text, owner_id::text, title, description, controller_email, privacy_purpose, retention_policy, status, public_slug, published_at, created_at, updated_at
		FROM forms
		WHERE id = $1 AND owner_id = $2
	`, formID, ownerID)
	if err != nil {
		return Form{}, err
	}

	fields, err := repo.ListFields(ctx, form.ID)
	if err != nil {
		return Form{}, err
	}

	form.Fields = fields
	return form, nil
}

func (repo *Repository) GetPublishedBySlug(ctx context.Context, slug string) (Form, error) {
	form, err := repo.findForm(ctx, `
		SELECT id::text, owner_id::text, title, description, controller_email, privacy_purpose, retention_policy, status, public_slug, published_at, created_at, updated_at
		FROM forms
		WHERE public_slug = $1 AND status = 'published'
	`, slug)
	if err != nil {
		return Form{}, err
	}

	fields, err := repo.ListFields(ctx, form.ID)
	if err != nil {
		return Form{}, err
	}

	form.Fields = fields
	return form, nil
}

func (repo *Repository) Update(ctx context.Context, ownerID string, formID string, input FormInput) (Form, error) {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return Form{}, fmt.Errorf("begin update form: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	form, err := scanForm(tx.QueryRow(ctx, `
		UPDATE forms
		SET title = $3,
			description = $4,
			controller_email = $5,
			privacy_purpose = $6,
			retention_policy = $7,
			updated_at = now()
		WHERE id = $1 AND owner_id = $2
		RETURNING id::text, owner_id::text, title, description, controller_email, privacy_purpose, retention_policy, status, public_slug, published_at, created_at, updated_at
	`, formID, ownerID, input.Title, input.Description, input.ControllerEmail, input.PrivacyPurpose, input.RetentionPolicy))
	if err != nil {
		return Form{}, err
	}

	fields, err := replaceFields(ctx, tx, form.ID, input.Fields)
	if err != nil {
		return Form{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Form{}, fmt.Errorf("commit update form: %w", err)
	}

	form.Fields = fields
	return form, nil
}

func (repo *Repository) Delete(ctx context.Context, ownerID string, formID string) error {
	command, err := repo.db.Exec(ctx, "DELETE FROM forms WHERE id = $1 AND owner_id = $2", formID, ownerID)
	if err != nil {
		return fmt.Errorf("delete form: %w", err)
	}

	if command.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (repo *Repository) Publish(ctx context.Context, ownerID string, formID string, slug string) (Form, error) {
	form, err := scanForm(repo.db.QueryRow(ctx, `
		UPDATE forms
		SET status = 'published',
			public_slug = COALESCE(public_slug, $3),
			published_at = COALESCE(published_at, now()),
			updated_at = now()
		WHERE id = $1 AND owner_id = $2
		RETURNING id::text, owner_id::text, title, description, controller_email, privacy_purpose, retention_policy, status, public_slug, published_at, created_at, updated_at
	`, formID, ownerID, slug))
	if isUniqueViolation(err) {
		return Form{}, ErrSlugTaken
	}
	if err != nil {
		return Form{}, err
	}

	fields, err := repo.ListFields(ctx, form.ID)
	if err != nil {
		return Form{}, err
	}

	form.Fields = fields
	return form, nil
}

func (repo *Repository) Unpublish(ctx context.Context, ownerID string, formID string) (Form, error) {
	form, err := scanForm(repo.db.QueryRow(ctx, `
		UPDATE forms
		SET status = 'draft',
			published_at = NULL,
			updated_at = now()
		WHERE id = $1 AND owner_id = $2
		RETURNING id::text, owner_id::text, title, description, controller_email, privacy_purpose, retention_policy, status, public_slug, published_at, created_at, updated_at
	`, formID, ownerID))
	if err != nil {
		return Form{}, err
	}

	fields, err := repo.ListFields(ctx, form.ID)
	if err != nil {
		return Form{}, err
	}

	form.Fields = fields
	return form, nil
}

func (repo *Repository) ListFields(ctx context.Context, formID string) ([]Field, error) {
	rows, err := repo.db.Query(ctx, `
		SELECT id::text, form_id::text, position, type, label, required, placeholder, options, config, created_at, updated_at
		FROM form_fields
		WHERE form_id = $1
		ORDER BY position ASC
	`, formID)
	if err != nil {
		return nil, fmt.Errorf("list form fields: %w", err)
	}
	defer rows.Close()

	return scanFieldRows(rows)
}

func (repo *Repository) findForm(ctx context.Context, query string, args ...any) (Form, error) {
	return scanForm(repo.db.QueryRow(ctx, query, args...))
}

func insertForm(ctx context.Context, tx pgx.Tx, ownerID string, input FormInput) (Form, error) {
	return scanForm(tx.QueryRow(ctx, `
		INSERT INTO forms (owner_id, title, description, controller_email, privacy_purpose, retention_policy)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, owner_id::text, title, description, controller_email, privacy_purpose, retention_policy, status, public_slug, published_at, created_at, updated_at
	`, ownerID, input.Title, input.Description, input.ControllerEmail, input.PrivacyPurpose, input.RetentionPolicy))
}

func replaceFields(ctx context.Context, tx pgx.Tx, formID string, fields []FieldInput) ([]Field, error) {
	if _, err := tx.Exec(ctx, "DELETE FROM form_fields WHERE form_id = $1", formID); err != nil {
		return nil, fmt.Errorf("delete existing fields: %w", err)
	}

	result := make([]Field, 0, len(fields))
	for index, field := range fields {
		inserted, err := insertField(ctx, tx, formID, index, field)
		if err != nil {
			return nil, err
		}

		result = append(result, inserted)
	}

	return result, nil
}

func insertField(ctx context.Context, tx pgx.Tx, formID string, position int, field FieldInput) (Field, error) {
	optionsJSON, err := json.Marshal(field.Options)
	if err != nil {
		return Field{}, fmt.Errorf("marshal field options: %w", err)
	}

	configJSON, err := json.Marshal(field.Config)
	if err != nil {
		return Field{}, fmt.Errorf("marshal field config: %w", err)
	}

	return scanField(tx.QueryRow(ctx, `
		INSERT INTO form_fields (form_id, position, type, label, required, placeholder, options, config)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8::jsonb)
		RETURNING id::text, form_id::text, position, type, label, required, placeholder, options, config, created_at, updated_at
	`, formID, position, field.Type, field.Label, field.Required, field.Placeholder, string(optionsJSON), string(configJSON)))
}

func scanFormRows(rows pgx.Rows) ([]Form, error) {
	forms := []Form{}
	for rows.Next() {
		form, err := scanForm(rows)
		if err != nil {
			return nil, err
		}

		forms = append(forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate forms: %w", err)
	}

	return forms, nil
}

func scanFieldRows(rows pgx.Rows) ([]Field, error) {
	fields := []Field{}
	for rows.Next() {
		field, err := scanField(rows)
		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate form fields: %w", err)
	}

	return fields, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanForm(row scanner) (Form, error) {
	var form Form
	var description pgtype.Text
	var controllerEmail pgtype.Text
	var privacyPurpose pgtype.Text
	var retentionPolicy pgtype.Text
	var publicSlug pgtype.Text
	var publishedAt pgtype.Timestamptz

	err := row.Scan(
		&form.ID,
		&form.OwnerID,
		&form.Title,
		&description,
		&controllerEmail,
		&privacyPurpose,
		&retentionPolicy,
		&form.Status,
		&publicSlug,
		&publishedAt,
		&form.CreatedAt,
		&form.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Form{}, ErrNotFound
	}
	if err != nil {
		return Form{}, fmt.Errorf("scan form: %w", err)
	}

	form.Description = textPtr(description)
	form.ControllerEmail = textPtr(controllerEmail)
	form.PrivacyPurpose = textPtr(privacyPurpose)
	form.RetentionPolicy = textPtr(retentionPolicy)
	form.PublicSlug = textPtr(publicSlug)
	form.PublishedAt = timePtr(publishedAt)
	form.Fields = []Field{}

	return form, nil
}

func scanField(row scanner) (Field, error) {
	var field Field
	var placeholder pgtype.Text
	var optionsJSON []byte
	var configJSON []byte

	err := row.Scan(
		&field.ID,
		&field.FormID,
		&field.Position,
		&field.Type,
		&field.Label,
		&field.Required,
		&placeholder,
		&optionsJSON,
		&configJSON,
		&field.CreatedAt,
		&field.UpdatedAt,
	)
	if err != nil {
		return Field{}, fmt.Errorf("scan form field: %w", err)
	}

	field.Placeholder = textPtr(placeholder)
	if err := json.Unmarshal(optionsJSON, &field.Options); err != nil {
		return Field{}, fmt.Errorf("decode field options: %w", err)
	}
	if err := json.Unmarshal(configJSON, &field.Config); err != nil {
		return Field{}, fmt.Errorf("decode field config: %w", err)
	}
	if field.Config == nil {
		field.Config = map[string]any{}
	}

	return field, nil
}

func textPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}

func timePtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}
