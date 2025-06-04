package mid

import (
	"context"
	"net/http"

	"github.com/warlck/food-flow/business/web/auth"
	"github.com/warlck/food-flow/business/web/response"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/validate"
	"github.com/warlck/food-flow/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "error", "msg", err)

				var status int
				var appErr response.ErrorDocument

				switch {
				case response.IsError(err):
					reqErr := response.GetError(err)
					if validate.IsFieldErrors(reqErr.Err) {
						fieldErrors := validate.GetFieldErrors(reqErr.Err)
						appErr = response.ErrorDocument{
							Error:  "data validation error",
							Fields: fieldErrors.Fields(),
						}
						status = reqErr.Status
						break
					}
					appErr = response.ErrorDocument{
						Error: reqErr.Error(),
					}
					status = reqErr.Status
				case auth.IsAuthError(err):
					appErr = response.ErrorDocument{
						Error: http.StatusText(http.StatusUnauthorized),
					}
					status = http.StatusUnauthorized
				default:
					appErr = response.ErrorDocument{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, appErr, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}
			return nil
		}

		return h
	}

	return m
}
