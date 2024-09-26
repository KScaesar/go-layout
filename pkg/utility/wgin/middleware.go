package wgin

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/slog-gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

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
		c.Next()
	}
}

func O11YMetric(svcName string) gin.HandlerFunc {
	// Throughput
	HttpRequestsTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests",
	}, []string{"method", "route"})

	HttpErrorsTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_total_errors",
		Help:      "Total number of HTTP errors",
	}, []string{"method", "route", "code"})

	HttpRequestsInflight := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "The number of inflight requests being handled at the same time",
	}, []string{"method", "route"})

	// Latency
	HttpResponseSecond := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Histogram of response time for HTTP in seconds",
		Buckets:   []float64{0.05, 0.2, 0.4, 0.6, 0.8, 1, 5, 10, 30}, // 50 ms ~ 30 s
	}, []string{"method", "route", "code"})

	return func(c *gin.Context) {
		method := c.Request.Method
		route := c.FullPath()

		// metric1
		start := time.Now()

		// metric2
		HttpRequestsTotal.WithLabelValues(method, route).Inc()

		// metric3
		HttpRequestsInflight.WithLabelValues(method, route).Add(1)

		c.Next()

		// metric3
		HttpRequestsInflight.WithLabelValues(method, route).Add(-1)

		// metric2
		code := strconv.Itoa(c.Writer.Status())
		if code[0] == '4' || code[0] == '5' {
			HttpErrorsTotal.WithLabelValues(method, route, code).Inc()
		}

		// metric1
		duration := time.Since(start).Seconds()
		span := trace.SpanFromContext(c.Request.Context())
		traceId := span.SpanContext().TraceID()
		if traceId.IsValid() {
			traceLabels := prometheus.Labels{"trace_id": traceId.String()}
			HttpResponseSecond.WithLabelValues(method, route, code).(prometheus.ExemplarObserver).ObserveWithExemplar(duration, traceLabels)
		} else {
			HttpResponseSecond.WithLabelValues(method, route, code).Observe(duration)
		}

		return
	}
}

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
			slog.String("route", c.FullPath()),
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

		if len(c.Errors) > 0 || c.Writer.Status() >= http.StatusBadRequest {
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
