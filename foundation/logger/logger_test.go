package logger_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/warlck/food-flow/foundation/logger"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestLoggerExtractsTraceAndSpanIDs(t *testing.T) {
	tp := sdktrace.NewTracerProvider()
	otel.SetTracerProvider(tp)

	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", nil)

	tracer := otel.Tracer("test-tracer")
	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	log.Info(ctx, "test message")

	output := buf.String()

	if !strings.Contains(output, "trace_id") {
		t.Errorf("Expected trace_id in output, got: %s", output)
	}

	if !strings.Contains(output, "span_id") {
		t.Errorf("Expected span_id in output, got: %s", output)
	}

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	if !strings.Contains(output, traceID) {
		t.Errorf("Expected trace ID %s in output, got: %s", traceID, output)
	}

	if !strings.Contains(output, spanID) {
		t.Errorf("Expected span ID %s in output, got: %s", spanID, output)
	}
}
