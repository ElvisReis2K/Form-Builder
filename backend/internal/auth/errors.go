package auth

import "errors"

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrIdentityTaken      = errors.New("identity already linked")
	ErrNotFound           = errors.New("not found")
	ErrUnauthenticated    = errors.New("unauthenticated")
)

type ValidationError struct {
	Message string
}

func (err ValidationError) Error() string {
	return err.Message
}
