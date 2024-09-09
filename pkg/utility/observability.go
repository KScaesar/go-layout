package utility

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type O11YConfig struct {
	MetricPort_ string  `yaml:"MetricPort"`
	EnableTrace bool    `yaml:"EnableTrace"`
	TraceHost   string  `yaml:"TraceHost"`
	TracePort   string  `yaml:"TracePort"`
	SampleRate  float64 `yaml:"SampleRate"` // 0 ~ 1
}

func (o *O11YConfig) MetricPort() string {
	if o.MetricPort_ == "" {
		const DefaultMetricPort = "2112"
		o.MetricPort_ = DefaultMetricPort
	}
	return o.MetricPort_
}

func (o O11YConfig) TraceAddress() string {
	return fmt.Sprintf("%v:%v", o.TraceHost, o.TracePort)
}

//

func ServeObservability(svcName string, conf *O11YConfig) error {
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

		DefaultShutdown().AddPriorityShutdownAction(2, "trace", func() error {
			return provider.Shutdown(ctx)
		})
	}

	// base pprof
	// https://cs.opensource.google/go/go/+/refs/tags/go1.23.0:src/net/http/pprof/pprof.go;l=100-104
	//
	// custom pprof
	// https://pkg.go.dev/runtime/pprof#Profile
	http.Handle("/debug/pprof/profile/goroutine", pprof.Handler("goroutine"))
	http.Handle("/debug/pprof/profile/heap", pprof.Handler("heap"))
	http.Handle("/debug/pprof/profile/allocs", pprof.Handler("allocs"))
	http.Handle("/debug/pprof/profile/threadcreate", pprof.Handler("threadcreate"))
	http.Handle("/debug/pprof/profile/block", pprof.Handler("block"))
	http.Handle("/debug/pprof/profile/mutex", pprof.Handler("mutex"))

	// metric

	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{Addr: ":" + conf.MetricPort(), Handler: http.DefaultServeMux}
	go func() {
		DefaultLogger().Info("metric start", slog.String("url", "http://localhost:"+conf.MetricPort()+"/metrics"))
		err := server.ListenAndServe()
		DefaultShutdown().Notify(err)
	}()
	DefaultShutdown().AddPriorityShutdownAction(2, "metric", func() error {
		return server.Shutdown(ctx)
	})

	return nil
}

func GinO11YTrace(enableTrace bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		if !enableTrace {
			c.Next()
			return
		}

		ctx, span := otel.Tracer("").Start(c.Request.Context(), "")
		c.Request = c.Request.WithContext(ctx)
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.path", c.Request.URL.Path),
		)

		c.Next()
	}
}

func GinO11YMetric(svcName string, enableTrace bool) func(c *gin.Context) {
	// https://github.com/prometheus/prometheus/tree/main/docs/querying
	// https://github.com/slok/go-http-metrics/blob/master/metrics/prometheus/prometheus.go#L76-L99
	// https://github.com/slok/go-http-metrics/blob/master/middleware/middleware.go#L98-L105
	// https://github.com/brancz/prometheus-example-app

	var traceKeys []string
	if enableTrace {
		traceKeys = []string{"trace_id", "span_id"}
	}

	// Throughput, Error Rate
	HttpRequestsTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests",
	}, append([]string{"code", "method", "handler"}, traceKeys...))

	// Latency
	HttpResponseSecond := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Histogram of response time for HTTP in seconds",
		Buckets:   []float64{0.05, 0.2, 0.4, 0.6, 0.8, 1, 5, 10, 30}, // 50 ms ~ 30 s
	}, append([]string{"code", "method", "handler"}, traceKeys...))

	HttpRequestsInflight := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "The number of inflight requests being handled at the same time",
	}, []string{"method", "handler"})

	return func(c *gin.Context) {

		var traceValues []string
		if enableTrace {
			span := trace.SpanFromContext(c.Request.Context())
			traceID := span.SpanContext().TraceID().String()
			spanID := span.SpanContext().SpanID().String()
			traceValues = []string{traceID, spanID}
		}

		method := c.Request.Method
		handler := c.Request.URL.Path

		HttpRequestsInflight.WithLabelValues(method, handler).Add(1)
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			HttpRequestsInflight.WithLabelValues(method, handler).Add(-1)

			code := strconv.Itoa(c.Writer.Status())

			values := append(make([]string, 0, 5), code, method, handler)
			values = append(values, traceValues...)

			HttpRequestsTotal.WithLabelValues(values...).Inc()
			HttpResponseSecond.WithLabelValues(values...).Observe(duration)
		}()

		c.Next()
	}
}
