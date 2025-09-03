package mid

import (
	"context"
	"errors"
	"net/http"

	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/foundation/logger"
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
				var appErr *errs.Error
				if !errors.As(err, &appErr) {
					appErr = errs.Newf(errs.Internal, "Internal Server Error")
				}

				if err := web.Respond(ctx, w, appErr, errs.ErrCodeHttpStatus(appErr.Code)); err != nil {
					return err
				}

			}
			return nil
		}

		return h
	}

	return m
}
