// Package all binds all the routes into the specified app.
package all

import (
	categoryapi "github.com/warlck/food-flow/app/domain/categoryapp"
	checkapi "github.com/warlck/food-flow/app/domain/checkapp"
	menuitemapi "github.com/warlck/food-flow/app/domain/menuitemapp"
	restaurantapi "github.com/warlck/food-flow/app/domain/restaurantapp"
	userapi "github.com/warlck/food-flow/app/domain/userapp"
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

	restaurantapi.Routes(app, restaurantapi.Config{
		AuthClient:    cfg.AuthClient,
		RestaurantBus: cfg.RestaurantBus,
		Log:           cfg.Log,
	})

	categoryapi.Routes(app, categoryapi.Config{
		AuthClient:  cfg.AuthClient,
		CategoryBus: cfg.CategoryBus,
		Log:         cfg.Log,
	})

	menuitemapi.Routes(app, menuitemapi.Config{
		AuthClient:  cfg.AuthClient,
		MenuItemBus: cfg.MenuItemBus,
		Log:         cfg.Log,
	})
}
