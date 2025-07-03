package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/CrowderSoup/emoji-match/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
)

func NewTraceExporter(config *config.Config) (*otlptrace.Exporter, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(config.TracingEndpoint),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	return exporter, nil
}

func NewTracerProvider(exporter *otlptrace.Exporter, config *config.Config) (*trace.TracerProvider, error) {
	tp := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(config.ServiceName),
			),
		),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

func InvokeTracerProvider(lc fx.Lifecycle, tp *trace.TracerProvider) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return tp.Shutdown(ctx)
			},
		},
	)
}
