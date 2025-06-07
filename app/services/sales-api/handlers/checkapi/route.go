package checkapi

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"
	api := newAPI(cfg.Build, cfg.Log, cfg.DB)

	app.HandleFuncNoMiddleware(http.MethodGet, version, "/liveness", api.liveness)
	app.HandleFuncNoMiddleware(http.MethodGet, version, "/readiness", api.readiness)
}
