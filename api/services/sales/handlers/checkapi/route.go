package checkapi

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"
	api := newAPI(cfg.Build, cfg.Log, cfg.DB)
	authen := mid.Authenticate(cfg.Auth)
	athAdminOnly := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)

	app.HandleFuncNoMiddleware(http.MethodGet, version, "/liveness", api.liveness)
	app.HandleFuncNoMiddleware(http.MethodGet, version, "/readiness", api.readiness)
	app.HandleFunc(http.MethodGet, version, "/testauth", api.liveness, authen, athAdminOnly) // TESTING ONLY, REMOVE LATER
}
