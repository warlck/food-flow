package mid

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/auth"
)

type ctxKey int

// claimsKey is the key for the claims in the context.
const (
	claimsKey ctxKey = iota + 1
	userIDKey
)

// setClaims adds the claims to the context.
func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimsKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

// setUserID adds the user id to the context.
func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the user id from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}

	return v, nil
}
