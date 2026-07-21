package mux

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/web"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type routeAdderFunc func(*web.App, Config)

func (fn routeAdderFunc) Add(app *web.App, cfg Config) {
	fn(app, cfg)
}

func TestWebAPIRequestLogsUseOpenTelemetryTraceID(t *testing.T) {
	otel.SetTracerProvider(sdktrace.NewTracerProvider())

	var output bytes.Buffer
	log := logger.NewWithHandler(slog.NewJSONHandler(&output, nil))
	var traceID string
	app := WebAPI(Config{Log: log}, routeAdderFunc(func(app *web.App, _ Config) {
		app.HandleFunc(http.MethodGet, "", "/trace-test", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			traceID = trace.SpanFromContext(ctx).SpanContext().TraceID().String()
			return nil
		})
	}))

	app.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/trace-test", nil))

	if len(traceID) != 32 || strings.Contains(traceID, "-") {
		t.Fatalf("handler did not receive a valid OpenTelemetry trace ID: %q", traceID)
	}

	for _, line := range strings.Split(strings.TrimSpace(output.String()), "\n") {
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("decode request log: %v", err)
		}
		if got, _ := record["trace_id"].(string); got != traceID {
			t.Errorf("request log trace_id = %q, want %q", got, traceID)
		}
	}
}
