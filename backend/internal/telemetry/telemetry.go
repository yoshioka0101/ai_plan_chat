package telemetry

import (
	"context"
	"log/slog"

	"github.com/yoshioka0101/ai_plan_chat/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

// Init configures the OpenTelemetry SDK. It returns a shutdown function.
func Init(ctx context.Context, cfg *config.Config) (func(context.Context) error, error) {
	if !cfg.Telemetry.Enabled {
		return func(context.Context) error { return nil }, nil
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL(cfg.Telemetry.OTLPEndpoint),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
	if err != nil {
		return nil, err
	}

	attributes := []attribute.KeyValue{
		semconv.ServiceName(cfg.Telemetry.ServiceName),
		semconv.ServiceVersion(cfg.Telemetry.ServiceVersion),
	}
	if cfg.Telemetry.Environment != "" {
		attributes = append(attributes, semconv.DeploymentEnvironment(cfg.Telemetry.Environment))
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(attributes...),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, err
	}

	sampleRatio := cfg.Telemetry.SampleRatio
	if sampleRatio < 0 || sampleRatio > 1 {
		sampleRatio = 1
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampleRatio))),
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	slog.Info("OpenTelemetry initialized",
		"service", cfg.Telemetry.ServiceName,
		"endpoint", cfg.Telemetry.OTLPEndpoint,
		"sample_ratio", sampleRatio,
	)

	return tp.Shutdown, nil
}
