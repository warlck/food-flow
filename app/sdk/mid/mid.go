package mid

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/business/domain/userbus"
)

type ctxKey int

// claimsKey is the key for the claims in the context.
const (
	claimsKey ctxKey = iota + 1
	userIDKey
	userKey
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
	if ok && v != uuid.Nil {
		return v, nil
	}

	claims := GetClaims(ctx)
	if claims.Subject != "" {
		return uuid.Parse(claims.Subject)
	}

	return uuid.UUID{}, errors.New("user id not found in context")
}

func setUser(ctx context.Context, usr userbus.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

// GetUser returns the user from the context.
func GetUser(ctx context.Context) (userbus.User, error) {
	v, ok := ctx.Value(userKey).(userbus.User)
	if !ok {
		return userbus.User{}, errors.New("user not found in context")
	}

	return v, nil
}
