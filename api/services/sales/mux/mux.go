// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/api/services/sales/handlers/checkapi"
	"github.com/warlck/food-flow/api/services/sales/handlers/userapi"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build    string
	Log      *logger.Logger
	Auth     *authclient.Client
	Shutdown chan os.Signal
	DB       *sqlx.DB
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	checkapi.Routes(app, checkapi.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
		Auth:  cfg.Auth,
	})

	userapi.Routes(app, userapi.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
		Auth:  cfg.Auth,
	})
	return app
}
