package promoapp

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build         string
	Log           *logger.Logger
	AuthClient    *authclient.Client
	PromoBus      *promobus.Business
	RestaurantBus *restaurantbus.Business
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.PromoBus, cfg.RestaurantBus)

	// Public endpoint for customer checkout validation
	app.HandleFunc(http.MethodPost, version, "/promotions/validate", api.validate)

	// Admin CRUD endpoints
	app.HandleFunc(http.MethodGet, version, "/promotions", api.query, authen)
	app.HandleFunc(http.MethodGet, version, "/promotions/{promotion_id}", api.queryByID, authen)
	app.HandleFunc(http.MethodPost, version, "/promotions", api.create, authen, ruleAdmin)
	app.HandleFunc(http.MethodPut, version, "/promotions/{promotion_id}", api.update, authen, ruleAdmin)
	app.HandleFunc(http.MethodDelete, version, "/promotions/{promotion_id}", api.delete, authen, ruleAdmin)
}
