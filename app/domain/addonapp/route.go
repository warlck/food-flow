package addonapp

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build         string
	Log           *logger.Logger
	AuthClient    *authclient.Client
	AddonBus      *addonbus.Business
	MenuItemBus   *menuitembus.Business
	RestaurantBus *restaurantbus.Business
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.AddonBus, cfg.MenuItemBus, cfg.RestaurantBus)
	app.HandleFunc(http.MethodGet, version, "/addons", api.query, authen)
	app.HandleFunc(http.MethodGet, version, "/addons/{addon_id}", api.queryByID, authen)
	app.HandleFunc(http.MethodPost, version, "/addons", api.create, authen, ruleAdmin)
	app.HandleFunc(http.MethodPut, version, "/addons/{addon_id}", api.update, authen, ruleAdmin)
	app.HandleFunc(http.MethodDelete, version, "/addons/{addon_id}", api.delete, authen, ruleAdmin)
	app.HandleFunc(http.MethodPost, version, "/addons/reorder", api.reorder, authen, ruleAdmin)
}
