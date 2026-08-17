package imageapp

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/imagebus"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/storage"
	"github.com/warlck/food-flow/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build      string
	Log        *logger.Logger
	AuthClient *authclient.Client
	ImageBus   *imagebus.Business
	// LocalStore is set only when the local storage backend is active; it
	// enables the development upload/download endpoints.
	LocalStore storage.LocalStore
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAdmin := mid.Authorize(cfg.AuthClient, auth.RuleAdminOnly)

	api := newApp(cfg.ImageBus, cfg.LocalStore)

	// Admin image management endpoints
	app.HandleFunc(http.MethodPost, version, "/images/upload-url", api.createUpload, authen, ruleAdmin)
	app.HandleFunc(http.MethodPost, version, "/images/{image_id}/complete", api.complete, authen, ruleAdmin)
	app.HandleFunc(http.MethodGet, version, "/images", api.query, authen)
	app.HandleFunc(http.MethodDelete, version, "/images/{image_id}", api.delete, authen, ruleAdmin)

	// Local development backend endpoints; unused when GCS is configured.
	if cfg.LocalStore != nil {
		app.HandleFunc(http.MethodPut, version, "/images/local/{path...}", api.localUpload, authen, ruleAdmin)
		app.HandleFunc(http.MethodGet, version, "/images/local/{path...}", api.localDownload)
	}
}
