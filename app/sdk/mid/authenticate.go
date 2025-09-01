package mid

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/foundation/web"
)

// Authenticate is a middleware function that integrates with an authentication client
// to validate user credentials and attach user data to the request context.
func Authenticate(ac *authclient.Client) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			resp, err := ac.Authenticate(r.Context(), r.Header.Get("Authorization"))
			if err != nil {
				return auth.NewAuthError("authentication failed: %v", err)
			}
			ctx = setClaims(ctx, resp.Claims)
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// Bearer processes JWT authentication logic.
func Bearer(ath *auth.Auth) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			claims, err := ath.Authenticate(ctx, r.Header.Get("authorization"))
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			if claims.Subject == "" {
				return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
			}

			ctx = setUserID(ctx, subjectID)
			ctx = setClaims(ctx, claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

// Basic processes basic authentication logic.
func Basic(ath *auth.Auth) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "5a4f42ef-5439-423d-881f-4a3628efeaf1", // TODO: Remove hardcoded subject and replace with actual user ID
					Issuer:    ath.Issuer(),
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				},
				Roles: []string{"ADMIN"},
			}

			subjectID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return errs.Newf(errs.Unauthenticated, "parsing subject: %s", err)
			}

			ctx = setUserID(ctx, subjectID)
			ctx = setClaims(ctx, claims)

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}

func parseBasicAuth(auth string) (string, string, bool) {
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}

	c, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", false
	}

	return username, password, true
}
