package responses

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Create(ctx context.Context, formID string, answers map[string]any) (Response, error) {
	answersJSON, err := json.Marshal(answers)
	if err != nil {
		return Response{}, fmt.Errorf("marshal response answers: %w", err)
	}

	response, err := scanResponse(repo.db.QueryRow(ctx, `
		INSERT INTO form_responses (form_id, answers, privacy_acknowledged_at)
		VALUES ($1, $2::jsonb, now())
		RETURNING id::text, form_id::text, answers, privacy_acknowledged_at, submitted_at
	`, formID, string(answersJSON)))
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func (repo *Repository) ListByForm(ctx context.Context, ownerID string, formID string) ([]Response, error) {
	rows, err := repo.db.Query(ctx, `
		SELECT response.id::text, response.form_id::text, response.answers, response.privacy_acknowledged_at, response.submitted_at
		FROM form_responses response
		INNER JOIN forms form ON form.id = response.form_id
		WHERE response.form_id = $1 AND form.owner_id = $2
		ORDER BY response.submitted_at DESC
	`, formID, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list form responses: %w", err)
	}
	defer rows.Close()

	return scanResponseRows(rows)
}

func (repo *Repository) DeleteByID(ctx context.Context, ownerID string, formID string, responseID string) error {
	command, err := repo.db.Exec(ctx, `
		DELETE FROM form_responses response
		USING forms form
		WHERE response.id = $1
			AND response.form_id = $2
			AND form.id = response.form_id
			AND form.owner_id = $3
	`, responseID, formID, ownerID)
	if err != nil {
		return fmt.Errorf("delete form response: %w", err)
	}

	if command.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func scanResponseRows(rows pgx.Rows) ([]Response, error) {
	responses := []Response{}
	for rows.Next() {
		response, err := scanResponse(rows)
		if err != nil {
			return nil, err
		}

		responses = append(responses, response)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate form responses: %w", err)
	}

	return responses, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanResponse(row scanner) (Response, error) {
	var response Response
	var answersJSON []byte

	err := row.Scan(
		&response.ID,
		&response.FormID,
		&answersJSON,
		&response.PrivacyAcknowledgedAt,
		&response.SubmittedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Response{}, ErrNotFound
	}
	if err != nil {
		return Response{}, fmt.Errorf("scan form response: %w", err)
	}

	if err := json.Unmarshal(answersJSON, &response.Answers); err != nil {
		return Response{}, fmt.Errorf("decode response answers: %w", err)
	}
	if response.Answers == nil {
		response.Answers = map[string]any{}
	}

	return response, nil
}
