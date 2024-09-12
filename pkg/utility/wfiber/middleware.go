package wfiber

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/gofiber/fiber/v2"

	"github.com/KScaesar/go-layout/pkg/utility/wlog"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/slog-fiber"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func O11YTrace(enableTrace bool) fiber.Handler {
	skipMethods := map[string]bool{
		http.MethodHead:    true,
		http.MethodConnect: true,
		http.MethodOptions: true,
	}

	return func(c *fiber.Ctx) error {
		if !enableTrace || skipMethods[c.Method()] {
			return c.Next()
		}

		ctx, span := otel.Tracer("").Start(c.Context(), "")
		c.SetUserContext(ctx)
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", string(c.Context().Method())),
			attribute.String("http.path", c.Path()),
		)

		return c.Next()
	}
}

func O11YMetric(svcName string, enableTrace bool) fiber.Handler {
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

	skipMethods := map[string]bool{
		http.MethodHead:    true,
		http.MethodConnect: true,
		http.MethodOptions: true,
	}

	return func(c *fiber.Ctx) error {
		if skipMethods[c.Method()] {
			return c.Next()
		}

		var traceValues []string
		if enableTrace {
			span := trace.SpanFromContext(c.UserContext())
			traceID := span.SpanContext().TraceID().String()
			spanID := span.SpanContext().SpanID().String()
			traceValues = []string{traceID, spanID}
		}

		method := c.Method()
		handler := c.Path()

		HttpRequestsInflight.WithLabelValues(method, handler).Add(1)
		start := time.Now()
		defer func() {
			duration := time.Since(start).Seconds()
			HttpRequestsInflight.WithLabelValues(method, handler).Add(-1)

			code := strconv.Itoa(c.Response().StatusCode())

			values := append(make([]string, 0, 5), code, method, handler)
			values = append(values, traceValues...)

			HttpRequestsTotal.WithLabelValues(values...).Inc()
			HttpResponseSecond.WithLabelValues(values...).Observe(duration)
		}()

		return c.Next()
	}
}

func O11YLogger(debug bool, enableTrace bool, wlogger *wlog.Logger) (fiber.Handler, fiber.Handler) {
	var config slogfiber.Config
	config.WithRequestID = true

	if debug {
		config = slogfiber.Config{
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
	slogfiber.RequestIDKey = "req_id"

	handler1 := slogfiber.NewWithConfig(wlogger.Logger, config)

	handler2 := func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		reqId := string(c.Response().Header.Peek(slogfiber.RequestIDHeaderKey))
		requestAttributes := []slog.Attr{
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
		}
		logger := wlogger.With(
			slog.Any("request", slog.GroupValue(requestAttributes...)),
			slog.String(slogfiber.RequestIDKey, reqId),
		)

		if enableTrace {
			span := trace.SpanFromContext(ctx)
			traceId := span.SpanContext().TraceID().String()
			spanId := span.SpanContext().SpanID().String()

			logger = logger.With(
				slog.String(slogfiber.TraceIDKey, traceId),
				slog.String(slogfiber.SpanIDKey, spanId),
			)
		}

		c.SetUserContext(wlogger.CtxWithLogger(ctx, logger))
		return c.Next()
	}

	return handler1, handler2
}

// GormTX
//
// 若 skip == nil, 所有條件都會使用 tx
// 若 skip != nil, 滿足條件的, 將不會啟動 tx
func GormTX(db *gorm.DB, skip func(ctx *fiber.Ctx) bool, wlogger *wlog.Logger) fiber.Handler {
	skipMethods := map[string]bool{
		http.MethodHead:    true,
		http.MethodConnect: true,
		http.MethodOptions: true,
		http.MethodTrace:   true,
	}

	return func(c *fiber.Ctx) error {
		canSkip := db == nil ||
			skipMethods[c.Method()] ||
			(skip != nil && skip(c))

		if canSkip {
			return c.Next()
		}

		stdCtx := c.UserContext()
		logger := wlogger.CtxGetLogger(stdCtx)

		tx := db.Begin()
		err := tx.Error
		if err != nil {
			logger.Error("gorm tx begin failed", "err", err)
			return err
		}

		c.SetUserContext(utility.CtxWithGormTX(stdCtx, db, tx))

		err = c.Next()
		if err != nil {
			Err := tx.Rollback().Error
			if Err != nil {
				logger.Error("gorm tx rollback failed", "err", Err)
			}
			return err
		}

		Err := tx.Commit().Error
		if Err != nil {
			logger.Error("gorm tx commit failed", "err", Err)
			return Err
		}
		return nil
	}
}
