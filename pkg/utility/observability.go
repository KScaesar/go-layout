package utility

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type O11YConfig struct {
	Port        string  `yaml:"Port"`
	EnableTrace bool    `yaml:"EnableTrace"`
	TraceHost   string  `yaml:"TraceHost"`
	TracePort   string  `yaml:"TracePort"`
	SampleRate  float64 `yaml:"SampleRate"` // 0 ~ 1
}

func (o O11YConfig) TraceAddress() string {
	return fmt.Sprintf("%v:%v", o.TraceHost, o.TracePort)
}

//

func ServeObservability(
	svcName string,
	conf *O11YConfig,
	logger *slog.Logger,
	shutdown *Shutdown,
) error {
	ctx := context.Background()

	if conf.EnableTrace {
		exporter, err := otlptrace.New(
			ctx,
			otlptracehttp.NewClient(
				otlptracehttp.WithEndpoint(conf.TraceAddress()),
				otlptracehttp.WithInsecure(),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}

		provider := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(svcName),
			)),
			// sdktrace.WithSampler(sdktrace.TraceIDRatioBased(conf.SampleRate)),
		)
		otel.SetTracerProvider(provider)

		shutdown.AddPriorityShutdownAction(2, "trace", func() error {
			return provider.Shutdown(ctx)
		})
	}

	// base pprof
	// https://cs.opensource.google/go/go/+/refs/tags/go1.23.0:src/net/http/pprof/pprof.go;l=100-104
	//
	// custom pprof
	// https://pkg.go.dev/runtime/pprof#Profile

	// metric
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{Addr: "0.0.0.0:" + conf.Port, Handler: http.DefaultServeMux}
	go func() {
		logger.Info("pprof start", slog.String("url", "http://0.0.0.0:"+conf.Port+"/debug/pprof"))
		logger.Info("metric start", slog.String("url", "http://0.0.0.0:"+conf.Port+"/metrics"))
		err := server.ListenAndServe()
		shutdown.Notify(err)
	}()
	shutdown.AddPriorityShutdownAction(2, "metric_&_pprof", func() error {
		return server.Shutdown(ctx)
	})

	return nil
}
