package userapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	user "github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/web/auth"
	"github.com/warlck/food-flow/business/web/response"
	"github.com/warlck/food-flow/foundation/web"
)

// Handlers manages the set of user endpoints.
type api struct {
	user *user.Business
	auth *auth.Auth
}

// New constructs a handlers for route access.
func newAPI(user *user.Business, auth *auth.Auth) *api {
	return &api{
		user: user,
		auth: auth,
	}
}

// Create adds a new user to the system.
func (h *api) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewUser
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nb, err := toBusNewUser(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nb)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}
