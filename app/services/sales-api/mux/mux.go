// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"os"

	"github.com/warlck/food-flow/app/services/sales-api/route/sys/checkapi"
	"github.com/warlck/food-flow/business/web/auth"
	"github.com/warlck/food-flow/business/web/mid"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
)

// WebAPI constructs a http.Handler with all application routes bound.
func WebAPI(log *logger.Logger, auth *auth.Auth, shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics())

	checkapi.Routes(mux, auth)
	return mux
}
