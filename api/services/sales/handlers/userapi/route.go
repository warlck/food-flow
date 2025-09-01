package userapi

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/domain/userbus/stores/userdb"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	usrBus := userbus.NewBusiness(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))
	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)

	api := newAPI(usrBus)
	app.HandleFunc(http.MethodPost, version, "/users", api.create, authen, ruleAdmin)
	app.HandleFunc(http.MethodGet, version, "/users", api.query, authen, ruleAdmin)
	app.HandleFunc(http.MethodGet, version, "/users/{user_id}", api.queryByID, authen)
}
