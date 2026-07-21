package mux

import (
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

type BusConfig struct {
	UserBus             *userbus.Business
	RestaurantBus       *restaurantbus.Business
	CategoryBus         *categorybus.Business
	MenuItemBus         *menuitembus.Business
	OrderBus            *orderbus.Business
	AddonBus            *addonbus.Business
	StripeSecretKey     string
	StripeWebhookSecret string
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build      string
	Log        *logger.Logger
	AuthClient *authclient.Client
	DB         *sqlx.DB
	Auth       *auth.Auth
	BusConfig
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type RouteAdder interface {
	Add(app *web.App, cfg Config)
}

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(cfg Config, routeAdder RouteAdder) *web.App {
	// Start the span before any request logging or error handling so every
	// request log contains the OpenTelemetry trace ID used by Tempo.
	app := web.NewApp(cfg.Log.Info, mid.OTEL(), mid.Logger(cfg.Log), mid.Metrics(), mid.Errors(cfg.Log), mid.Panics())

	routeAdder.Add(app, cfg)
	return app
}
