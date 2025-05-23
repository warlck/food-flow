package response

import "errors"

// ErrorDocument is the form used for API responses when an error occurs.
type ErrorDocument struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// Error is used to pass an error during the request through the call chain.
// Adding this as a separate type that implements the error interface allows
// us to handle errors in a way that is specific to the business domain.
type Error struct {
	Err    error
	Status int
}

// NewError wraps a provided error with an HTTP status code. This function
// should be used when handlers encounter expected errors.
func NewError(err error, status int) *Error {
	return &Error{
		Err:    err,
		Status: status,
	}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (r *Error) Error() string {
	return r.Err.Error()
}

// Unwrap implements the unwrap interface. It returns the wrapped error.
func (r *Error) Unwrap() error {
	return r.Err
}

// StatusCode returns our HTTP status code.
func (r *Error) StatusCode() int {
	return r.Status
}

// IsError returns true if the error is an Error type.
func IsError(err error) bool {
	var e *Error
	return errors.As(err, &e)
}

// GetError returns the error if it is an Error type. Otherwise it returns nil.
func GetError(err error) *Error {
	var e *Error
	if !errors.As(err, &e) {
		return nil
	}
	return e
}
