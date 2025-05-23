package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/warlck/food-flow/business/web/metrics"
	"github.com/warlck/food-flow/foundation/web"
)

// Panics recovers from panics and logs the panic.
func Panics() web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC: [%v] TRACE: [%s]", rec, string(trace))

					metrics.AddPanics(ctx)
				}
			}()

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
