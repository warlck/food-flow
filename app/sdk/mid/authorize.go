package mid

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/foundation/web"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// Authorize is a middleware function that integrates with an authentication client
// to validate user credentials and attach user data to the request context.
func Authorize(client *authclient.Client, rule string) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims := GetClaims(ctx)

			if claims.Subject == "" {
				return errs.New(errs.Unauthenticated, auth.NewAuthError("unauthorized: no claims provided"))
			}

			var userID uuid.UUID
			id := web.Param(r, "userID")
			if id != "" {
				userID, err := uuid.Parse(id)
				if err != nil {
					return errs.New(errs.Unauthenticated, err)
				}
				ctx = setUserID(ctx, userID)
			}

			ath := authclient.Authorize{
				UserID: userID,
				Rule:   rule,
				Claims: claims,
			}
			if err := client.Authorize(ctx, ath); err != nil {
				// Preserve the auth service's error code (403 for an
				// authenticated-but-not-allowed caller); only transport
				// failures collapse to 401 (fail-closed).
				var appErr *errs.Error
				if errors.As(err, &appErr) {
					return appErr
				}
				return errs.New(errs.Unauthenticated, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// AuthorizeUser executes the specified role and extracts the specified
// user from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(client *authclient.Client, userBus *userbus.Business, rule string) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			id := web.Param(r, "user_id")

			var userID uuid.UUID

			if id != "" {
				var err error
				userID, err = uuid.Parse(id)
				if err != nil {
					return errs.New(errs.Unauthenticated, ErrInvalidID)
				}

				usr, err := userBus.QueryByID(ctx, userID)
				if err != nil {
					switch {
					case errors.Is(err, userbus.ErrNotFound):
						return errs.New(errs.Unauthenticated, err)
					default:
						return errs.Newf(errs.Unauthenticated, "querybyid: userID[%s]: %s", userID, err)
					}
				}

				ctx = setUser(ctx, usr)
			}

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			auth := authclient.Authorize{
				Claims: GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
			}

			if err := client.Authorize(ctx, auth); err != nil {
				var appErr *errs.Error
				if errors.As(err, &appErr) {
					return appErr
				}
				return errs.New(errs.Unauthenticated, err)
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
