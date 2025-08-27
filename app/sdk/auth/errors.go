package auth

import (
	"errors"
	"fmt"
)

// authError represents an error that occurs during authentication or authorization.
type authError struct {
	message string
}

// NewAuthError creates a new authError with the given format and arguments.
func NewAuthError(format string, args ...any) error {
	return &authError{message: fmt.Sprintf(format, args...)}
}

// Error implements the error interface. It return the message of the wrapped error.
func (e *authError) Error() string {
	return e.message
}

// IsAuthError checks if the error is an authError.
func IsAuthError(err error) bool {
	var e *authError
	return errors.As(err, &e)
}
