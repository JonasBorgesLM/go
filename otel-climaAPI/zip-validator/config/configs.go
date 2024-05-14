package configs

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LoadConfig carrega o arquivo .env a partir do caminho especificado
func LoadConfig(path string) error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("erro ao ler o arquivo de configuração: %w", err)
	}

	viper.AutomaticEnv()
	return nil
}

// Setup configura o OpenTelemetry para o rastreamento
func Setup() (func(ctx context.Context) error, error) {
	ctx := context.Background()

	// Cria um novo recurso com atributos
	res, err := createResource(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar recurso: %w", err)
	}

	// Cria um novo exportador de rastreamento
	traceExporter, err := createTraceExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar exportador de rastreamento: %w", err)
	}

	// Cria um novo provedor de rastreamento
	tp := createTracerProvider(res, traceExporter)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown, nil
}

// createResource cria um novo recurso com atributos de serviço
func createResource(ctx context.Context) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(viper.GetString("SERVICE")),
		),
	)
}

// createTraceExporter cria um novo exportador de rastreamento gRPC
func createTraceExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		viper.GetString("OTEL_COLLECTOR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar conexão gRPC com o coletor: %w", err)
	}

	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

// createTracerProvider cria um novo provedor de rastreamento
func createTracerProvider(res *resource.Resource, exporter sdktrace.SpanExporter) *sdktrace.TracerProvider {
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
}
