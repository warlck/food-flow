// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"os"

	"github.com/warlck/food-flow/app/services/sales-api/route/sys/checkapi"
	"github.com/warlck/food-flow/business/web/auth"
	"github.com/warlck/food-flow/business/web/mid"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build    string
	Log      *logger.Logger
	Auth     *auth.Auth
	Shutdown chan os.Signal
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	checkapi.Routes(app, cfg.Auth, cfg.Build, cfg.Log)
	return app
}
