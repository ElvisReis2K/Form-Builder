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
		INSERT INTO form_responses (form_id, answers)
		VALUES ($1, $2::jsonb)
		RETURNING id::text, form_id::text, answers, submitted_at
	`, formID, string(answersJSON)))
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func (repo *Repository) ListByForm(ctx context.Context, ownerID string, formID string) ([]Response, error) {
	rows, err := repo.db.Query(ctx, `
		SELECT response.id::text, response.form_id::text, response.answers, response.submitted_at
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
