package organizationapp

import (
	"context"
	"net/http"

	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/foundation/web"
)

type app struct {
	orgBus *organizationbus.Business
}

func newApp(orgBus *organizationbus.Business) *app {
	return &app{
		orgBus: orgBus,
	}
}

// queryMyOrgs returns the organizations the authenticated user belongs to.
func (a *app) queryMyOrgs(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.Newf(errs.Unauthenticated, "user not authenticated: %s", err)
	}

	orgs, err := a.orgBus.QueryOrgsForUser(ctx, userID)
	if err != nil {
		return errs.Newf(errs.Internal, "querymyorgs: %s", err)
	}

	return web.Respond(ctx, w, toAppOrganizations(orgs), http.StatusOK)
}
