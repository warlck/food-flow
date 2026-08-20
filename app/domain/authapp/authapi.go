package authapi

import (
	"context"
	"net/http"
	"net/mail"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/types/role"
	"github.com/warlck/food-flow/foundation/web"
)

type api struct {
	auth    *auth.Auth
	userBus *userbus.Business
}

func newAPI(a *auth.Auth, userBus *userbus.Business) *api {
	return &api{
		auth:    a,
		userBus: userBus,
	}
}

// login authenticates a user by email and password and issues a JWT signed
// with the server's active key. Only users with the ADMIN role can log in.
func (a *api) login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var lr LoginRequest
	if err := web.Decode(r, &lr); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	email, err := mail.ParseAddress(lr.Email)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Authenticate(ctx, *email, lr.Password)
	if err != nil {
		return errs.New(errs.Unauthenticated, errInvalidCredentials)
	}

	if !slices.Contains(usr.Roles, role.Admin) {
		return errs.New(errs.PermissionDenied, errAdminRequired)
	}

	now := time.Now().UTC()
	expiresAt := now.Add(a.auth.TokenTTL())

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Roles: role.ParseToString(usr.Roles),
	}

	tkn, err := a.auth.GenerateToken(a.auth.ActiveKID(), claims)
	if err != nil {
		return errs.Newf(errs.Internal, "generating token: %s", err)
	}

	resp := LoginResponse{
		Token:     tkn,
		ExpiresAt: expiresAt,
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func (a *api) authenticate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// The middleware is actually handling the authentication. So if the code
	// gets to this handler, authentication passed.

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	resp := authclient.AuthenticateResp{
		UserID: userID,
		Claims: mid.GetClaims(ctx),
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func (a *api) authorize(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var auth authclient.Authorize
	if err := web.Decode(r, &auth); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := a.auth.Authorize(ctx, auth.Claims, auth.UserID, auth.Rule); err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]", auth.Claims.Roles, auth.Rule)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
