package checkapi

import (
	"github.com/warlck/food-flow/business/web/auth"
	"github.com/warlck/food-flow/business/web/mid"
	"github.com/warlck/food-flow/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, a *auth.Auth) {
	authenticate := mid.Authenticate(a)
	authAdminOnly := mid.Authorize(a, auth.RuleAdminOnly)
	app.HandleFuncNoMiddleware("GET /liveness", liveness)
	app.HandleFuncNoMiddleware("GET /readiness", readiness)
	app.HandleFunc("GET /testerror", testError, authenticate, authAdminOnly)
}
