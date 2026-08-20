package authapi

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Auth    *auth.Auth
	UserBus *userbus.Business
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	bearer := mid.Bearer(cfg.Auth)

	api := newAPI(cfg.Auth, cfg.UserBus)

	app.HandleFunc(http.MethodPost, version, "/auth/login", api.login)
	app.HandleFunc(http.MethodGet, version, "/auth/authenticate", api.authenticate, bearer)
	app.HandleFunc(http.MethodPost, version, "/auth/authorize", api.authorize)
}
