package forms

import "errors"

var (
	ErrNotFound  = errors.New("form not found")
	ErrSlugTaken = errors.New("form public slug already exists")
)

type ValidationError struct {
	Message string
}

func (err ValidationError) Error() string {
	return err.Message
}
