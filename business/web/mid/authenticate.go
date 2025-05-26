package mid

import (
	"context"
	"net/http"

	"github.com/warlck/food-flow/business/web/auth"
	"github.com/warlck/food-flow/foundation/web"
)

// Authenticate is a middleware function that integrates with an authentication client
// to validate user credentials and attach user data to the request context.
func Authenticate(a *auth.Auth) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, err := a.Authenticate(r.Context(), r.Header.Get("Authorization"))
			if err != nil {
				return auth.NewAuthError("authentication failed: %v", err)
			}
			ctx = auth.SetClaims(ctx, claims)
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
