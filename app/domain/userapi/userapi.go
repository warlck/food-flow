package userapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// Handlers manages the set of user endpoints.
type api struct {
	userBus *userbus.Business
}

// New constructs a handlers for route access.
func newAPI(user *userbus.Business) *api {
	return &api{
		userBus: user,
	}
}

// Create adds a new user to the system.
func (a *api) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewUser
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nb, err := toBusNewUser(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, nb)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}

// Query retrieves a list of users based on query parameters.
func (a *api) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}
	pg, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	usrs, err := a.userBus.Query(ctx, filter, orderBy, pg.Number(), pg.RowsPerPage())
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.userBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	users := query.NewResult(toAppUsers(usrs), total, pg)
	return web.Respond(ctx, w, users, http.StatusOK)
}

// QueryByID retrieves a user by its ID.
func (a *api) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userIDStr := web.Param(r, "user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errs.NewFieldErrors("user_id", err)
	}

	usr, err := a.userBus.QueryByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("querybyid: userID[%s]: %w", userID, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
}
