package menuitemapp

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/categorybus"
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
	MenuItemBus   *menuitembus.Business
	RestaurantBus *restaurantbus.Business
	CategoryBus   *categorybus.Business
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.MenuItemBus, cfg.RestaurantBus, cfg.CategoryBus)
	app.HandleFunc(http.MethodGet, version, "/menuitems", api.query, authen)
	app.HandleFunc(http.MethodGet, version, "/menuitems/{menuitem_id}", api.queryByID, authen)
	app.HandleFunc(http.MethodPost, version, "/menuitems", api.create, authen, ruleAdmin)
	app.HandleFunc(http.MethodPut, version, "/menuitems/order", api.reorder, authen, ruleAdmin)
	app.HandleFunc(http.MethodPut, version, "/menuitems/{menuitem_id}", api.update, authen, ruleAdmin)
	app.HandleFunc(http.MethodDelete, version, "/menuitems/{menuitem_id}", api.delete, authen, ruleAdmin)
}
