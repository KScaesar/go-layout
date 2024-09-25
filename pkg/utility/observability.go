package utility

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"github.com/felixge/fgprof"
	"github.com/prometheus/client_golang/prometheus"
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

func InitO11YTracer(conf *O11YConfig, shutdown *Shutdown, svcName string) error {
	if !conf.EnableTrace {
		return nil
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(conf.TraceAddress()),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return fmt.Errorf("create OTLP trace exporter: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(svcName),
		)),
	)
	otel.SetTracerProvider(provider)

	shutdown.AddPriorityShutdownAction(2, "trace", func() error {
		return provider.Shutdown(context.Background())
	})
	return nil
}

func ServeO11YMetric(port string, shutdown *Shutdown, logger *slog.Logger) {
	// pprof
	// https://cs.opensource.google/go/go/+/refs/tags/go1.23.0:src/net/http/pprof/pprof.go;l=100-104
	// https://pkg.go.dev/runtime/pprof#Profile

	// fgprof
	// https://github.com/felixge/fgprof?tab=readme-ov-file#how-it-works
	http.Handle("/debug/fgprof", fgprof.Handler())

	// metric
	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	server := &http.Server{Addr: "0.0.0.0:" + port, Handler: http.DefaultServeMux}
	shutdown.AddPriorityShutdownAction(2, "metric_&_pprof", func() error {
		return server.Shutdown(context.Background())
	})

	go func() {
		logger.Info("pprof start", slog.String("url", "http://0.0.0.0:"+port+"/debug/pprof"))
		logger.Info("fgprof start", slog.String("url", "http://0.0.0.0:"+port+"/debug/fgprof?seconds=1"))
		logger.Info("metric start", slog.String("url", "http://0.0.0.0:"+port+"/metrics"))
		err := server.ListenAndServe()
		shutdown.Notify(err)
	}()
}
