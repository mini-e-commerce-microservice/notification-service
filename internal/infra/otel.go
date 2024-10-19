package infra

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/mini-e-commerce-microservice/notification-service/generated/proto/secret_proto"
	"github.com/mini-e-commerce-microservice/notification-service/internal/util/primitive"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"time"
)

func NewOtelCollector(cred *secret_proto.Otel, tracerName string) primitive.CloseFn {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cred.Username, cred.Password)))
	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cred.Endpoint),
		otlptracegrpc.WithHeaders(map[string]string{
			"Authorization": authHeader,
		}),
	)

	traceExporter, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		panic(err)
	}

	provider := &otelProvider{
		name: tracerName,
	}

	traceProvide, closeFnTracer, err := provider.start(traceExporter)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed initializing the tracer provider")
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(traceProvide)

	return closeFnTracer
}

type otelProvider struct {
	name     string
	exporter trace.SpanExporter
}

func (o *otelProvider) start(exp trace.SpanExporter) (*trace.TracerProvider, primitive.CloseFn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(o.name),
		),
	)
	if err != nil {
		err = fmt.Errorf("failed to created resource: %w", err)
		return nil, nil, err
	}

	o.exporter = exp
	bsp := trace.NewBatchSpanProcessor(o.exporter)

	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)

	return provider, func(ctx context.Context) (err error) {
		log.Info().Msg("starting shutdown export and provider")
		ctxClosure, cancelClosure := context.WithTimeout(ctx, 5*time.Second)
		defer cancelClosure()

		if err = o.exporter.Shutdown(ctxClosure); err != nil {
			return err
		}

		if err = provider.Shutdown(ctxClosure); err != nil {
			return err
		}

		log.Info().Msg("shutdown export and provider successfully")

		return
	}, nil
}
