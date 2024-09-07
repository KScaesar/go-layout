package utility

import (
	"context"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ServeObservability(port string) {
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

	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{Addr: ":" + port, Handler: http.DefaultServeMux}
	DefaultShutdown.StopService("o11y", func() error {
		return server.Shutdown(context.Background())
	})
	go func() {
		err := server.ListenAndServe()
		DefaultShutdown.Notify(err)
	}()
}

func GinHttpObservability(svcName string) func(c *gin.Context) {
	// https://github.com/prometheus/prometheus/tree/main/docs/querying
	// https://github.com/slok/go-http-metrics/blob/master/metrics/prometheus/prometheus.go#L76-L99
	// https://github.com/slok/go-http-metrics/blob/master/middleware/middleware.go#L98-L105
	// https://github.com/brancz/prometheus-example-app

	// Throughput, Error Rate
	HttpRequestsTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests",
	}, []string{"code", "method", "handler"})

	// Latency
	HttpResponseSecond := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Histogram of response time for HTTP in seconds",
		Buckets:   []float64{0.1, 0.2, 0.4, 0.6, 0.8, 1, 5, 10, 30, 60}, // 100 ms ~ 60 s
	}, []string{"code", "method", "handler"})

	HttpRequestsInflight := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "The number of inflight requests being handled at the same time",
	})

	return func(c *gin.Context) {
		method := c.Request.Method
		handler := c.Request.URL.Path

		HttpRequestsInflight.Inc()
		// HttpRequestsInflight.WithLabelValues(method, handler).Add(1)
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		// HttpRequestsInflight.WithLabelValues(method, handler).Add(-1)
		HttpRequestsInflight.Desc()

		code := strconv.Itoa(c.Writer.Status())
		HttpRequestsTotal.WithLabelValues(code, method, handler).Inc()
		HttpResponseSecond.WithLabelValues(code, method, handler).Observe(duration)
	}
}
