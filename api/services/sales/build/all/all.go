// Package all binds all the routes into the specified app.
package all

import (
	"github.com/warlck/food-flow/app/domain/checkapi"
	"github.com/warlck/food-flow/app/domain/userapi"
	"github.com/warlck/food-flow/app/sdk/mux"
	"github.com/warlck/food-flow/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {
	checkapi.Routes(app, checkapi.Config{
		Build:      cfg.Build,
		Log:        cfg.Log,
		DB:         cfg.DB,
		AuthClient: cfg.AuthClient,
	})

	userapi.Routes(app, userapi.Config{
		AuthClient: cfg.AuthClient,
		Build:      cfg.Build,
		UserBus:    cfg.UserBus,
		Log:        cfg.Log,
	})
}
