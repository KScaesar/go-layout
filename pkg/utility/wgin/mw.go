package wgin

import (
	"log/slog"
	"net/http"
	"slices"
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

func O11YLogger(debug bool, enableTrace bool, Logger *wlog.Logger) []gin.HandlerFunc {
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

	return []gin.HandlerFunc{
		sloggin.NewWithConfig(Logger.Logger, config),

		func(c *gin.Context) {
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
		},
	}
}

func O11YTrace(enableTrace bool) func(c *gin.Context) {
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

func O11YMetric(svcName string, enableTrace bool) func(c *gin.Context) {
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

// GormTransaction
//
// 若 skipPaths 長度為 0，表示所有 path 都會使用 tx
// 若 skipPaths 長度不為 0，表示 skipPaths 中的路徑將不會啟動 tx
func GormTransaction(db *gorm.DB, skipPaths []string) gin.HandlerFunc {
	skipMethods := map[string]bool{
		http.MethodHead:    true,
		http.MethodConnect: true,
		http.MethodOptions: true,
		http.MethodTrace:   true,

		http.MethodGet:    false,
		http.MethodPost:   false,
		http.MethodPut:    false,
		http.MethodPatch:  false,
		http.MethodDelete: false,
	}

	return func(c *gin.Context) {
		canSkip := db == nil ||
			skipMethods[c.Request.Method] ||
			(len(skipPaths) != 0 && slices.Contains(skipPaths, c.Request.URL.Path))

		if canSkip {
			c.Next()
			return
		}

		tx := db.Begin()
		err := tx.Error
		if err != nil {
			c.Error(err)
			return
		}

		c.Request = c.Request.WithContext(utility.CtxWithGormTX(c.Request.Context(), db, tx))
		c.Next()

		if len(c.Errors) > 0 {
			err := tx.Rollback().Error
			if err != nil {

			}
			return
		}

		err = tx.Commit().Error
		if err != nil {
			c.Error(err)
			return
		}
	}
}
