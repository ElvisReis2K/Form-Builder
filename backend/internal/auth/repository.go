package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type Repository struct {
	db *pgxpool.Pool
}

type userWithPassword struct {
	User
	PasswordHash string
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) CreateUser(ctx context.Context, email string, name string, passwordHash string) (User, error) {
	var user User
	err := repo.db.QueryRow(ctx, `
		INSERT INTO users (email, name, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id::text, email, name, created_at
	`, email, name, passwordHash).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	if isUniqueViolation(err) {
		return User{}, ErrEmailTaken
	}
	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (repo *Repository) FindUserByEmail(ctx context.Context, email string) (userWithPassword, error) {
	var user userWithPassword
	err := repo.db.QueryRow(ctx, `
		SELECT id::text, email, name, COALESCE(password_hash, ''), created_at
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return userWithPassword{}, ErrNotFound
	}
	if err != nil {
		return userWithPassword{}, fmt.Errorf("find user by email: %w", err)
	}

	return user, nil
}

func (repo *Repository) CreateSession(ctx context.Context, userID string, tokenHash string, expiresAt time.Time) error {
	_, err := repo.db.Exec(ctx, `
		INSERT INTO sessions (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	return nil
}

func (repo *Repository) FindUserBySessionHash(ctx context.Context, tokenHash string, now time.Time) (User, error) {
	var user User
	err := repo.db.QueryRow(ctx, `
		SELECT u.id::text, u.email, u.name, u.created_at
		FROM sessions s
		INNER JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1
			AND s.revoked_at IS NULL
			AND s.expires_at > $2
	`, tokenHash, now).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("find user by session: %w", err)
	}

	return user, nil
}

func (repo *Repository) RevokeSession(ctx context.Context, tokenHash string) error {
	_, err := repo.db.Exec(ctx, `
		UPDATE sessions
		SET revoked_at = now()
		WHERE token_hash = $1
			AND revoked_at IS NULL
	`, tokenHash)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}

	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}
