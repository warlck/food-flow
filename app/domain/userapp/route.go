package userapi

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build      string
	Log        *logger.Logger
	AuthClient *authclient.Client
	UserBus    *userbus.Business
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)
	ruleAuthorizeUser := mid.AuthorizeUser(cfg.AuthClient, cfg.UserBus, auth.RuleAdminOrSubject)
	ruleAuthorizeAdmin := mid.AuthorizeUser(cfg.AuthClient, cfg.UserBus, auth.RuleAdminOnly)

	api := newApp(cfg.UserBus)
	app.HandleFunc(http.MethodGet, version, "/users", api.query, authen, ruleAdmin)
	app.HandleFunc(http.MethodGet, version, "/users/{user_id}", api.queryByID, authen, ruleAuthorizeUser)
	app.HandleFunc(http.MethodPost, version, "/users", api.create, authen, ruleAdmin)
	app.HandleFunc(http.MethodPut, version, "/users/role/{user_id}", api.updateRole, authen, ruleAuthorizeAdmin)
	app.HandleFunc(http.MethodPut, version, "/users/{user_id}", api.update, authen, ruleAuthorizeUser)
	app.HandleFunc(http.MethodDelete, version, "/users/{user_id}", api.delete, authen, ruleAuthorizeUser)
}
