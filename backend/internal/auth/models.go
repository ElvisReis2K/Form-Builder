package auth

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

type AuthResult struct {
	User      User
	Token     string
	ExpiresAt time.Time
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type GoogleIdentityInput struct {
	Subject string
	Email   string
	Name    string
}
