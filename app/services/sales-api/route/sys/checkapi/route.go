package checkapi

import "github.com/warlck/food-flow/foundation/web"

// Routes adds specific routes for this group.
func Routes(app *web.App) {
	app.HandleFuncNoMiddleware("GET /liveness", liveness)
	app.HandleFuncNoMiddleware("GET /readiness", readiness)
	app.HandleFuncNoMiddleware("GET /testerror", testError)
}
