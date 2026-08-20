package authapi

import (
	"errors"
	"fmt"
	"time"

	"github.com/warlck/food-flow/app/sdk/errs"
)

// Login failure errors. Every credential failure returns
// errInvalidCredentials so the response does not reveal which part of the
// credentials was wrong or whether the account exists.
var (
	errInvalidCredentials = errors.New("invalid credentials")
	errAdminRequired      = errors.New("admin role required")
)

// LoginRequest represents the credentials submitted to the login endpoint.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Validate checks the data in the model is considered clean.
func (app LoginRequest) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// LoginResponse represents the token issued after a successful login.
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
