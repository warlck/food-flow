package mux

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/api/services/auth/handlers/authapi"
	"github.com/warlck/food-flow/api/services/auth/handlers/checkapi"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Auth     *auth.Auth
	Build    string
	DB       *sqlx.DB
	Log      *logger.Logger
	Shutdown chan os.Signal
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	checkapi.Routes(app, checkapi.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	authapi.Routes(app, authapi.Config{
		Auth: cfg.Auth,
	})

	return app
}
