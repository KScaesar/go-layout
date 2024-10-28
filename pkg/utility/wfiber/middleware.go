package wfiber

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/slog-fiber"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
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

		ctx, span := otel.Tracer("").Start(c.Context(), "", trace.WithSpanKind(trace.SpanKindServer))
		c.SetUserContext(ctx)
		defer span.End()
		return c.Next()
	}
}

func O11YMetric(svcName string) fiber.Handler {
	// metric2-a
	HttpRequestsTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests",
	}, []string{"method", "route"})

	// metric2-b
	HttpErrorsTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_total_errors",
		Help:      "Total number of HTTP errors",
	}, []string{"code", "method", "route"})

	// metric3
	HttpRequestsInflight := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "The number of inflight Http requests being handled at the same time",
	}, []string{"method", "route"})

	// metric1
	HttpResponseSecond := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: svcName,
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Histogram of response time for HTTP in seconds",
		Buckets:   []float64{0.05, 0.2, 0.4, 0.6, 0.8, 1, 5, 10, 30}, // 50 ms ~ 30 s
	}, []string{"code", "method", "route"})

	return func(c *fiber.Ctx) error {
		method := c.Method()
		route := c.Route().Path

		// metric1
		start := time.Now()

		// metric2-a
		HttpRequestsTotal.WithLabelValues(method, route).Inc()

		// metric3
		HttpRequestsInflight.WithLabelValues(method, route).Add(1)

		err := c.Next()

		// metric3
		HttpRequestsInflight.WithLabelValues(method, route).Add(-1)

		// metric2-b
		code := strconv.Itoa(c.Response().StatusCode())
		if code[0] == '4' || code[0] == '5' {
			HttpErrorsTotal.WithLabelValues(code, method, route).Inc()
		}

		// metric1
		duration := time.Since(start).Seconds()
		span := trace.SpanFromContext(c.UserContext())
		traceId := span.SpanContext().TraceID()
		if traceId.IsValid() {
			traceLabels := prometheus.Labels{"trace_id": traceId.String()}
			HttpResponseSecond.WithLabelValues(code, method, route).(prometheus.ExemplarObserver).ObserveWithExemplar(duration, traceLabels)
		} else {
			HttpResponseSecond.WithLabelValues(code, method, route).Observe(duration)
		}

		return err
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
			slog.String("route", c.Route().Path),
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
		if err != nil || c.Response().StatusCode() >= http.StatusBadRequest {
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
