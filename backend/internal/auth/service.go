package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 8
	maxPasswordBytes  = 72
)

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type Service struct {
	repo          *Repository
	sessionSecret string
	sessionTTL    time.Duration
}

func NewService(repo *Repository, sessionSecret string, sessionTTL time.Duration) *Service {
	return &Service{
		repo:          repo,
		sessionSecret: sessionSecret,
		sessionTTL:    sessionTTL,
	}
}

func (service *Service) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	name := strings.TrimSpace(input.Name)
	email := normalizeEmail(input.Email)

	if name == "" {
		return AuthResult{}, ValidationError{Message: "nome é obrigatório"}
	}

	if err := validateEmail(email); err != nil {
		return AuthResult{}, err
	}

	if err := validatePassword(input.Password); err != nil {
		return AuthResult{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := service.repo.CreateUser(ctx, email, name, string(passwordHash))
	if errors.Is(err, ErrEmailTaken) {
		return AuthResult{}, ErrEmailTaken
	}
	if err != nil {
		return AuthResult{}, err
	}

	return service.createSession(ctx, user)
}

func (service *Service) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	email := normalizeEmail(input.Email)
	if err := validateEmail(email); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	user, err := service.repo.FindUserByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		return AuthResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return AuthResult{}, err
	}

	if user.PasswordHash == "" {
		return AuthResult{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	return service.createSession(ctx, user.User)
}

func (service *Service) EnsureDefaultAdmin(ctx context.Context, input DefaultAdminInput) error {
	name := strings.TrimSpace(input.Name)
	email := normalizeEmail(input.Email)

	if name == "" {
		return ValidationError{Message: "nome do usuário padrão é obrigatório"}
	}

	if err := validateEmail(email); err != nil {
		return err
	}

	if err := validatePassword(input.Password); err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash default admin password: %w", err)
	}

	if err := service.repo.EnsureUserPassword(ctx, email, name, string(passwordHash)); err != nil {
		return err
	}

	return nil
}

func (service *Service) LoginWithGoogle(ctx context.Context, input GoogleIdentityInput) (AuthResult, error) {
	subject := strings.TrimSpace(input.Subject)
	email := normalizeEmail(input.Email)
	name := strings.TrimSpace(input.Name)

	if subject == "" {
		return AuthResult{}, ValidationError{Message: "identificador do Google é obrigatório"}
	}

	if err := validateEmail(email); err != nil {
		return AuthResult{}, err
	}

	if name == "" {
		name = fallbackNameFromEmail(email)
	}

	user, err := service.repo.FindUserByIdentity(ctx, "google", subject)
	if err == nil {
		return service.createSession(ctx, user)
	}
	if !errors.Is(err, ErrNotFound) {
		return AuthResult{}, err
	}

	userWithPassword, err := service.repo.FindUserByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		user, err = service.repo.CreateUser(ctx, email, name, "")
		if errors.Is(err, ErrEmailTaken) {
			userWithPassword, err = service.repo.FindUserByEmail(ctx, email)
			if err != nil {
				return AuthResult{}, err
			}

			user = userWithPassword.User
		} else if err != nil {
			return AuthResult{}, err
		}
	} else if err != nil {
		return AuthResult{}, err
	} else {
		user = userWithPassword.User
	}

	if err := service.repo.LinkIdentity(ctx, user.ID, "google", subject, email); err != nil {
		if errors.Is(err, ErrIdentityTaken) {
			user, err := service.repo.FindUserByIdentity(ctx, "google", subject)
			if err != nil {
				return AuthResult{}, err
			}

			return service.createSession(ctx, user)
		}

		return AuthResult{}, err
	}

	return service.createSession(ctx, user)
}

func (service *Service) Authenticate(ctx context.Context, token string) (User, error) {
	if strings.TrimSpace(token) == "" {
		return User{}, ErrUnauthenticated
	}

	user, err := service.repo.FindUserBySessionHash(ctx, service.hashSessionToken(token), time.Now().UTC())
	if errors.Is(err, ErrNotFound) {
		return User{}, ErrUnauthenticated
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (service *Service) Logout(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return ErrUnauthenticated
	}

	return service.repo.RevokeSession(ctx, service.hashSessionToken(token))
}

func (service *Service) createSession(ctx context.Context, user User) (AuthResult, error) {
	token, err := newSessionToken()
	if err != nil {
		return AuthResult{}, err
	}

	expiresAt := time.Now().UTC().Add(service.sessionTTL)
	if err := service.repo.CreateSession(ctx, user.ID, service.hashSessionToken(token), expiresAt); err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		User:      user,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (service *Service) hashSessionToken(token string) string {
	mac := hmac.New(sha256.New, []byte(service.sessionSecret))
	_, _ = mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}

func newSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateEmail(email string) error {
	if !emailPattern.MatchString(email) {
		return ValidationError{Message: "informe um e-mail válido"}
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ValidationError{Message: "a senha deve ter pelo menos 8 caracteres"}
	}

	if len([]byte(password)) > maxPasswordBytes {
		return ValidationError{Message: "a senha deve ter no máximo 72 bytes"}
	}

	return nil
}

func fallbackNameFromEmail(email string) string {
	name, _, ok := strings.Cut(email, "@")
	if !ok || strings.TrimSpace(name) == "" {
		return "Usuário Google"
	}

	return name
}
