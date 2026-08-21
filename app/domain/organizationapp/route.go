package organizationapp

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	AuthClient *authclient.Client
	OrgBus     *organizationbus.Business
	Log        *logger.Logger
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)
	// We might only need authentication to list organizations they belong to.

	appCtx := newApp(cfg.OrgBus)

	app.HandleFunc(http.MethodGet, version, "/organizations/me", appCtx.queryMyOrgs, authen, ruleAdmin)
}
