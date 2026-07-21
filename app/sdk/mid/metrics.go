package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/warlck/food-flow/app/sdk/metrics"
	"github.com/warlck/food-flow/foundation/web"
)

func Metrics() web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)
			started := time.Now()

			err := handler(ctx, w, r)

			n := metrics.AddRequests(ctx)
			if n%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}

			statusCode := web.GetValues(ctx).StatusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
				if err != nil {
					statusCode = http.StatusInternalServerError
				}
			}
			metrics.RecordHTTPRequest(r, statusCode, time.Since(started))

			if statusCode >= http.StatusBadRequest {
				metrics.AddErrors(ctx)
			}

			return err
		}

		return h
	}

	return m
}
