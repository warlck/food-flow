package auth

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey int

// claimsKey is the key for the claims in the context.
const claimsKey ctxKey = 1

// userIDKey is the key for the user ID in the context.
const userIDKey ctxKey = 2

// SetClaims adds the claims to the context.
func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) Claims {
	claims, ok := ctx.Value(claimsKey).(Claims)
	if !ok {
		return Claims{}
	}
	return claims
}

// SetUserID adds the user ID to the context.
func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the user ID from the context.
func GetUserID(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}
