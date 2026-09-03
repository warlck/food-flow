// Package all binds all the routes into the specified app.
package all

import (
	addonapi "github.com/warlck/food-flow/app/domain/addonapp"
	categoryapi "github.com/warlck/food-flow/app/domain/categoryapp"
	checkapi "github.com/warlck/food-flow/app/domain/checkapp"
	imageapi "github.com/warlck/food-flow/app/domain/imageapp"
	menuitemapi "github.com/warlck/food-flow/app/domain/menuitemapp"
	modifiergroupapi "github.com/warlck/food-flow/app/domain/modifiergroupapp"
	modifieroptionapi "github.com/warlck/food-flow/app/domain/modifieroptionapp"
	orderapi "github.com/warlck/food-flow/app/domain/orderapp"
	organizationapi "github.com/warlck/food-flow/app/domain/organizationapp"
	promoapi "github.com/warlck/food-flow/app/domain/promoapp"
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
		AuthClient:        cfg.AuthClient,
		RestaurantBus:     cfg.RestaurantBus,
		CategoryBus:       cfg.CategoryBus,
		MenuItemBus:       cfg.MenuItemBus,
		ModifierGroupBus:  cfg.ModifierGroupBus,
		ModifierOptionBus: cfg.ModifierOptionBus,
		AddonBus:          cfg.AddonBus,
		Log:               cfg.Log,
	})

	organizationapi.Routes(app, organizationapi.Config{
		AuthClient: cfg.AuthClient,
		OrgBus:     cfg.OrgBus,
		Log:        cfg.Log,
	})

	categoryapi.Routes(app, categoryapi.Config{
		AuthClient:    cfg.AuthClient,
		CategoryBus:   cfg.CategoryBus,
		RestaurantBus: cfg.RestaurantBus,
		Log:           cfg.Log,
	})

	menuitemapi.Routes(app, menuitemapi.Config{
		AuthClient:    cfg.AuthClient,
		MenuItemBus:   cfg.MenuItemBus,
		RestaurantBus: cfg.RestaurantBus,
		CategoryBus:   cfg.CategoryBus,
		Log:           cfg.Log,
	})

	modifiergroupapi.Routes(app, modifiergroupapi.Config{
		AuthClient:       cfg.AuthClient,
		ModifierGroupBus: cfg.ModifierGroupBus,
		RestaurantBus:    cfg.RestaurantBus,
		MenuItemBus:      cfg.MenuItemBus,
		Log:              cfg.Log,
	})

	modifieroptionapi.Routes(app, modifieroptionapi.Config{
		AuthClient:        cfg.AuthClient,
		ModifierOptionBus: cfg.ModifierOptionBus,
		ModifierGroupBus:  cfg.ModifierGroupBus,
		RestaurantBus:     cfg.RestaurantBus,
		Log:               cfg.Log,
	})

	addonapi.Routes(app, addonapi.Config{
		AuthClient:    cfg.AuthClient,
		AddonBus:      cfg.AddonBus,
		MenuItemBus:   cfg.MenuItemBus,
		RestaurantBus: cfg.RestaurantBus,
		Log:           cfg.Log,
	})

	promoapi.Routes(app, promoapi.Config{
		Build:         cfg.Build,
		Log:           cfg.Log,
		AuthClient:    cfg.AuthClient,
		PromoBus:      cfg.PromoBus,
		RestaurantBus: cfg.RestaurantBus,
	})

	imageapi.Routes(app, imageapi.Config{
		Build:         cfg.Build,
		Log:           cfg.Log,
		AuthClient:    cfg.AuthClient,
		ImageBus:      cfg.ImageBus,
		RestaurantBus: cfg.RestaurantBus,
		LocalStore:    cfg.ImageLocalStore,
	})

	orderapi.Routes(app, orderapi.Config{
		Build:               cfg.Build,
		Log:                 cfg.Log,
		AuthClient:          cfg.AuthClient,
		OrderBus:            cfg.OrderBus,
		RestaurantBus:       cfg.RestaurantBus,
		MenuItemBus:         cfg.MenuItemBus,
		CategoryBus:         cfg.CategoryBus,
		StripeSecretKey:     cfg.StripeSecretKey,
		StripeWebhookSecret: cfg.StripeWebhookSecret,
	})
}
