package checkapi

import (
	"net/http"

	"github.com/warlck/food-flow/business/web/auth"
	"github.com/warlck/food-flow/business/web/mid"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, a *auth.Auth, build string, log *logger.Logger) {
	const version = "v1"
	api := newAPI(build, log)

	authenticate := mid.Authenticate(a)
	authAdminOnly := mid.Authorize(a, auth.RuleAdminOnly)
	app.HandleFuncNoMiddleware(http.MethodGet, version, "/liveness", api.liveness)
	app.HandleFuncNoMiddleware(http.MethodGet, version, "/readiness", api.readiness)
	app.HandleFunc(http.MethodGet, version, "/testerror", api.testError, authenticate, authAdminOnly)
}
