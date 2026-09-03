package categoryapp

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build         string
	Log           *logger.Logger
	AuthClient    *authclient.Client
	CategoryBus   *categorybus.Business
	RestaurantBus *restaurantbus.Business
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.CategoryBus, cfg.RestaurantBus)
	app.HandleFunc(http.MethodGet, version, "/categories", api.query, authen)
	app.HandleFunc(http.MethodPut, version, "/categories/order", api.reorder, authen, ruleAdmin)
	app.HandleFunc(http.MethodGet, version, "/categories/{category_id}", api.queryByID, authen)
	app.HandleFunc(http.MethodPost, version, "/categories", api.create, authen, ruleAdmin)
	app.HandleFunc(http.MethodPut, version, "/categories/{category_id}", api.update, authen, ruleAdmin)
	app.HandleFunc(http.MethodDelete, version, "/categories/{category_id}", api.delete, authen, ruleAdmin)
}
