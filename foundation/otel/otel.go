package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Init initializes OpenTelemetry tracer and meter providers.
func Init(ctx context.Context, serviceName string, traceURL string) (func(context.Context) error, error) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	// Set up TracerProvider
	var traceExporter trace.SpanExporter
	if traceURL != "" {
		traceExporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(traceURL))
		if err != nil {
			return nil, fmt.Errorf("creating trace exporter: %w", err)
		}
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
	)
	if traceExporter != nil {
		tracerProvider.RegisterSpanProcessor(trace.NewBatchSpanProcessor(traceExporter))
	}
	otel.SetTracerProvider(tracerProvider)

	// Set up MeterProvider with Prometheus exporter
	promExporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("creating prometheus exporter: %w", err)
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(promExporter),
	)
	otel.SetMeterProvider(meterProvider)

	// Set global propagator to tracecontext
	otel.SetTextMapPropagator(propagation.TraceContext{})

	shutdown := func(ctx context.Context) error {
		var err error
		if errTrace := tracerProvider.Shutdown(ctx); errTrace != nil {
			err = fmt.Errorf("trace provider shutdown: %w", errTrace)
		}
		if errMetric := meterProvider.Shutdown(ctx); errMetric != nil {
			if err != nil {
				err = fmt.Errorf("%w; metric provider shutdown: %v", err, errMetric)
			} else {
				err = fmt.Errorf("metric provider shutdown: %w", errMetric)
			}
		}
		return err
	}

	return shutdown, nil
}
