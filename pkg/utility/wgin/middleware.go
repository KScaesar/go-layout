package wgin

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/slog-gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func O11YLogger(debug bool, enableTrace bool, Logger *wlog.Logger) (gin.HandlerFunc, gin.HandlerFunc) {
	var config sloggin.Config
	config.WithRequestID = true

	if debug {
		config = sloggin.Config{
			WithUserAgent:      false,
			WithRequestID:      true,
			WithRequestBody:    true,
			WithRequestHeader:  true,
			WithResponseBody:   false,
			WithResponseHeader: false,
		}
	}

	if enableTrace {
		config.WithTraceID = true
		config.WithSpanID = true
	}
	sloggin.RequestIDKey = "req_id"

	h1 := sloggin.NewWithConfig(Logger.Logger, config)

	h2 := func(c *gin.Context) {
		ctx := c.Request.Context()

		reqId := c.Writer.Header().Get(sloggin.RequestIDHeaderKey)
		requestAttributes := []slog.Attr{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		}
		logger := Logger.With(
			slog.Any("request", slog.GroupValue(requestAttributes...)),
			slog.String(sloggin.RequestIDKey, reqId),
		)

		if enableTrace {
			span := trace.SpanFromContext(ctx)
			traceId := span.SpanContext().TraceID().String()
			spanId := span.SpanContext().SpanID().String()

			logger = logger.With(
				slog.String(sloggin.TraceIDKey, traceId),
				slog.String(sloggin.SpanIDKey, spanId),
			)
		}

		c.Request = c.Request.WithContext(Logger.CtxWithLogger(ctx, logger))
		c.Next()
	}
	return h1, h2
}

func O11YTrace(enableTrace bool) func(c *gin.Context) {
	skipMethods := map[string]bool{
		http.MethodHead:    true,
		http.MethodConnect: true,
		http.MethodOptions: true,
	}
	return func(c *gin.Context) {
		if !enableTrace || skipMethods[c.Request.Method] {
			c.Next()
			return
		}

		ctx, span := otel.Tracer("").Start(c.Request.Context(), "", trace.WithSpanKind(trace.SpanKindServer))
		c.Request = c.Request.WithContext(ctx)
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.path", c.FullPath()),
		)

		c.Next()
	}
}

func O11YMetric(svcName string) func(c *gin.Context) {
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
		Buckets:   []float64{0.05, 0.2, 0.4, 0.6, 0.8, 1, 5, 10, 30}, // 50 ms ~ 30 s
	}, []string{"code", "method", "handler"})

	HttpRequestsInflight := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "The number of inflight requests being handled at the same time",
	}, []string{"method", "handler"})

	skipMethods := map[string]bool{
		http.MethodHead:    true,
		http.MethodConnect: true,
		http.MethodOptions: true,
	}
	return func(c *gin.Context) {
		if skipMethods[c.Request.Method] {
			c.Next()
			return
		}

		span := trace.SpanFromContext(c.Request.Context())
		traceId := span.SpanContext().TraceID()
		spanId := span.SpanContext().SpanID()

		method := c.Request.Method
		handler := c.FullPath()

		HttpRequestsInflight.WithLabelValues(method, handler).Add(1)
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			HttpRequestsInflight.WithLabelValues(method, handler).Add(-1)

			code := strconv.Itoa(c.Writer.Status())
			labels := []string{code, method, handler}
			HttpRequestsTotal.WithLabelValues(labels...).Inc()

			if traceId.IsValid() && spanId.IsValid() {
				traceLabels := prometheus.Labels{
					"trace_id": traceId.String(),
					"span_id":  spanId.String(),
				}
				HttpResponseSecond.WithLabelValues(labels...).(prometheus.ExemplarObserver).ObserveWithExemplar(duration, traceLabels)
			} else {
				HttpResponseSecond.WithLabelValues(labels...).Observe(duration)
			}
		}()

		c.Next()
	}
}

// GormTX
//
// 若 skip == nil, 所有條件都會使用 tx
// 若 skip != nil, 滿足條件的, 將不會啟動 tx
func GormTX(db *gorm.DB, skip func(ctx *gin.Context) bool, wlogger *wlog.Logger) gin.HandlerFunc {
	skipMethods := map[string]bool{
		http.MethodHead:    true,
		http.MethodConnect: true,
		http.MethodOptions: true,
		http.MethodTrace:   true,
	}

	return func(c *gin.Context) {
		canSkip := db == nil ||
			skipMethods[c.Request.Method] ||
			(skip != nil && skip(c))

		if canSkip {
			c.Next()
			return
		}

		stdCtx := c.Request.Context()
		logger := wlogger.CtxGetLogger(stdCtx)

		tx := db.Begin()
		err := tx.Error
		if err != nil {
			logger.Error("gorm tx begin failed", "err", err)
			c.Error(err)
			return
		}

		c.Request = c.Request.WithContext(utility.CtxWithGormTX(stdCtx, db, tx))

		c.Next()

		if len(c.Errors) > 0 {
			Err := tx.Rollback().Error
			if Err != nil {
				logger.Error("gorm tx rollback failed", "err", Err)
			}
			return
		}

		Err := tx.Commit().Error
		if Err != nil {
			logger.Error("gorm tx commit failed", "err", Err)
			c.Error(Err)
			return
		}
	}
}
