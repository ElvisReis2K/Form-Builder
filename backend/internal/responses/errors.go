package responses

import "errors"

var ErrNotFound = errors.New("response resource not found")

type ValidationError struct {
	Message string
}

func (err ValidationError) Error() string {
	return err.Message
}
