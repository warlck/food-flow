package authapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/foundation/web"
)

type api struct {
	auth *auth.Auth
}

func newAPI(a *auth.Auth) *api {
	return &api{
		auth: a,
	}
}

type token struct {
	Token string `json:"token"`
}

func (api *api) token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	kid := web.Param(r, "kid")
	if kid == "" {
		return errs.NewFieldErrors("kid", errors.New("missing kid"))
	}

	// The BearerBasic middleware function generates the claims.
	claims := mid.GetClaims(ctx)

	tkn, err := api.auth.GenerateToken(kid, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return web.Respond(ctx, w, token{Token: tkn}, http.StatusOK)
}

func (a *api) authenticate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// The middleware is actually handling the authentication. So if the code
	// gets to this handler, authentication passed.

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.New(errs.Unauthenticated, err)
	}

	resp := struct {
		UserID uuid.UUID   `json:"userID"`
		Claims auth.Claims `json:"claims"`
	}{
		UserID: userID,
		Claims: mid.GetClaims(ctx),
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func (a *api) authorize(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var auth struct {
		Claims auth.Claims `json:"claims"`
		UserID uuid.UUID   `json:"userID"`
		Rule   string      `json:"rule"`
	}
	if err := web.Decode(r, &auth); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if err := a.auth.Authorize(ctx, auth.Claims, auth.UserID, auth.Rule); err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]", auth.Claims.Roles, auth.Rule)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
