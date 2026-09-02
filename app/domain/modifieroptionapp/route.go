package modifieroptionapp

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build             string
	Log               *logger.Logger
	AuthClient        *authclient.Client
	ModifierOptionBus *modifieroptionbus.Business
	ModifierGroupBus  *modifiergroupbus.Business
	RestaurantBus     *restaurantbus.Business
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.ModifierOptionBus, cfg.ModifierGroupBus, cfg.RestaurantBus)
	app.HandleFunc(http.MethodGet, version, "/modifieroptions", api.query, authen)
	app.HandleFunc(http.MethodPut, version, "/modifieroptions/order", api.reorder, authen, ruleAdmin)
	app.HandleFunc(http.MethodGet, version, "/modifieroptions/{modifier_option_id}", api.queryByID, authen)
	app.HandleFunc(http.MethodPost, version, "/modifieroptions", api.create, authen, ruleAdmin)
	app.HandleFunc(http.MethodPut, version, "/modifieroptions/{modifier_option_id}", api.update, authen, ruleAdmin)
	app.HandleFunc(http.MethodDelete, version, "/modifieroptions/{modifier_option_id}", api.delete, authen, ruleAdmin)
}
