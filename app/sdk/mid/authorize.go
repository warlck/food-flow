package mid

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/business/web/response"
	"github.com/warlck/food-flow/foundation/web"
)

// Authorize is a middleware function that integrates with an authentication client
// to validate user credentials and attach user data to the request context.
func Authorize(client *authclient.Client, rule string) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims := GetClaims(ctx)

			if claims.Subject == "" {
				return auth.NewAuthError("unauthorized: no claims provided")
			}

			var userID uuid.UUID
			id := web.Param(r, "userID")
			if id != "" {
				userID, err := uuid.Parse(id)
				if err != nil {
					return response.NewError(errors.New("invalid userID"), http.StatusBadRequest)
				}
				ctx = setUserID(ctx, userID)
			}

			ath := authclient.Authorize{
				UserID: userID,
				Rule:   rule,
				Claims: claims,
			}
			if err := client.Authorize(ctx, ath); err != nil {
				return auth.NewAuthError("unauthorized: %v", err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
