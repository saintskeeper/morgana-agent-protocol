package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Provider holds the telemetry provider and tracer
type Provider struct {
	TracerProvider *sdktrace.TracerProvider
	Tracer         trace.Tracer
}

// Config holds telemetry configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	ExporterType   string // "otlp", "stdout", "none"
	OTLPEndpoint   string
	Debug          bool
}

// NewProvider creates a new telemetry provider
func NewProvider(ctx context.Context, cfg Config) (*Provider, error) {
	// Create resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	// Create exporter based on type
	var exporter sdktrace.SpanExporter
	switch cfg.ExporterType {
	case "otlp":
		exporter, err = createOTLPExporter(ctx, cfg.OTLPEndpoint)
		if err != nil {
			return nil, fmt.Errorf("creating OTLP exporter: %w", err)
		}
	case "stdout":
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithWriter(os.Stdout),
		)
		if err != nil {
			return nil, fmt.Errorf("creating stdout exporter: %w", err)
		}
	case "none":
		// No exporter - useful for testing
		exporter = nil
	default:
		return nil, fmt.Errorf("unknown exporter type: %s", cfg.ExporterType)
	}

	// Create tracer provider options
	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(res),
	}

	if exporter != nil {
		if cfg.Debug {
			// In debug mode, export spans immediately
			opts = append(opts, sdktrace.WithBatcher(exporter,
				sdktrace.WithBatchTimeout(time.Second),
			))
		} else {
			// In production, batch spans
			opts = append(opts, sdktrace.WithBatcher(exporter))
		}
	}

	// Add sampler
	if cfg.Debug {
		opts = append(opts, sdktrace.WithSampler(sdktrace.AlwaysSample()))
	} else {
		opts = append(opts, sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)))
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(opts...)

	// Set global provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Provider{
		TracerProvider: tp,
		Tracer:         tp.Tracer("morgana-protocol"),
	}, nil
}

// Shutdown gracefully shuts down the provider
func (p *Provider) Shutdown(ctx context.Context) error {
	return p.TracerProvider.Shutdown(ctx)
}

// createOTLPExporter creates an OTLP exporter
func createOTLPExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connecting to OTLP endpoint: %w", err)
	}

	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(conn),
	))
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	return exporter, nil
}

// AgentAttributes returns common attributes for agent spans
func AgentAttributes(agentType, taskID string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("agent.type", agentType),
		attribute.String("agent.task_id", taskID),
		attribute.String("agent.framework", "morgana-protocol"),
	}
}

// PromptAttributes returns attributes for prompt-related spans
func PromptAttributes(promptLen int, truncated bool) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Int("prompt.length", promptLen),
		attribute.Bool("prompt.truncated", truncated),
	}
}

// ResultAttributes returns attributes for result spans
func ResultAttributes(success bool, outputLen int, executionTimeMs int64) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Bool("result.success", success),
		attribute.Int("result.output_length", outputLen),
		attribute.Int64("result.execution_time_ms", executionTimeMs),
	}
}

// ErrorAttributes returns attributes for error spans
func ErrorAttributes(err error) []attribute.KeyValue {
	if err == nil {
		return nil
	}
	return []attribute.KeyValue{
		attribute.String("error.type", fmt.Sprintf("%T", err)),
		attribute.String("error.message", err.Error()),
	}
}
