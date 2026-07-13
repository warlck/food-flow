package mid_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/warlck/food-flow/app/sdk/mid"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func TestOTELMiddleware(t *testing.T) {
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	middleware := mid.OTEL()

	var spanCtx trace.SpanContext

	h := middleware(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		span := trace.SpanFromContext(ctx)
		spanCtx = span.SpanContext()
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	err := h(context.Background(), w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !spanCtx.IsValid() {
		t.Error("expected valid span context in handler")
	}
}
