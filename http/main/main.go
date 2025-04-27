package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/tinello/golang-openapi/core"
	http_delivery "github.com/tinello/golang-openapi/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {

	shutdown := initTracer()
	defer shutdown()

	provider := core.GetProviderInstance()

	server := &http.Server{
		Handler: otelhttp.NewHandler(http_delivery.NewServer(&provider, generateId), "otelhandler"),
	}

	listenAddress := ":" + core.MustGetEnv("PORT")
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf(
		"%s:%s server started on port %s...\n",
		"golang-openapi",
		http_delivery.GetApplicationVersion(),
		listenAddress)
	log.Fatalln(server.Serve(listener))
}

func generateId() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}

func initTracer() func() {
	ctx := context.Background()
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(core.MustGetEnv("OTEL_EXPORTER_OTLP_ENDPOINT")),
		otlptracegrpc.WithInsecure(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(core.MustGetEnv("OTEL_SERVICE_NAME")),
		)),
	)

	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}
