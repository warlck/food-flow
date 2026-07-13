package mid

import (
	"context"
	"net/http"

	"github.com/warlck/food-flow/foundation/web"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// OTEL constructs a middleware that traces incoming requests.
func OTEL() web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))

			tracer := otel.Tracer("http-handler")
			ctx, span := tracer.Start(ctx, r.URL.Path, trace.WithAttributes(
				semconv.HTTPRouteKey.String(r.URL.Path),
				semconv.HTTPRequestMethodKey.String(r.Method),
			))
			defer span.End()

			err := handler(ctx, w, r)

			if err != nil {
				span.RecordError(err)
			}

			return err
		}

		return h
	}

	return m
}
