package userapi

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/domain/userbus/stores/userdb"
	"github.com/warlck/food-flow/business/web/auth"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *auth.Auth
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	usrBus := userbus.NewBusiness(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))

	hdl := New(usrBus, cfg.Auth)
	app.HandleFunc(http.MethodPost, version, "/users", hdl.Create)
}
