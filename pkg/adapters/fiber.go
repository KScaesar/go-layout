package adapters

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/trace"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/dataflow"
)

var FiberMetadata = newFiberMetadataKey()

func newFiberMetadataKey() *fiberMetadataKey {
	return &fiberMetadataKey{
		errCode: "err_code",
	}
}

type fiberMetadataKey struct {
	errCode string
}

func (key *fiberMetadataKey) GetErrorCode(c *fiber.Ctx) int {
	errCode, ok := c.Context().UserValue(key.errCode).(int)
	if !ok {
		const successCode = 0
		return successCode
	}
	return errCode
}

func (key *fiberMetadataKey) SetErrorCode(c *fiber.Ctx, errCode int) {
	c.Context().SetUserValue(key.errCode, errCode)
}

//

func HandleErrorByFiber(c *fiber.Ctx, err error) error {
	myErr, ok := utility.UnwrapCustomError(err)
	if !ok {
		Err, isFixed := fixUnknownError(err)
		if isFixed {
			myErr, ok = utility.UnwrapCustomError(Err)
			err = Err
		}
	}

	if !ok {
		logger := pkg.Logger().CtxGetLogger(c.UserContext())
		logger.Warn("capture unknown error", slog.Any("err", err))
	}

	FiberMetadata.SetErrorCode(c, myErr.ErrorCode())

	DefaultErrorResponse := fiber.Map{
		"code": myErr.ErrorCode(),
		"msg":  err.Error(),
	}
	return c.Status(myErr.HttpStatus()).JSON(DefaultErrorResponse)
}

func fixUnknownError(err error) (Err error, isFixed bool) {
	var fiberErr *fiber.Error

	switch {
	case errors.As(err, &fiberErr):
		Err, isFixed = fixFiberError(fiberErr)
		if isFixed {
			return Err, true
		}
	}

	return err, false
}

func fixFiberError(err *fiber.Error) (error, bool) {
	switch err.Code {
	case fiber.StatusNotFound:
		return fmt.Errorf("%w: %w", pkg.ErrNotExists, err), true
	}
	return err, false
}

//

func ParseQueryByFiber(c *fiber.Ctx, req any, logger *slog.Logger) (bool, error) {
	err := c.QueryParser(req)
	if err != nil {
		logger.Error(err.Error(), slog.Any("cause", c.QueryParser))
		return false, HandleErrorByFiber(c, pkg.ErrInvalidParam)
	}
	return true, nil
}

//

func dataflowO11YMetric() dataflow.Middleware {
	svcName := pkg.Version().ServiceName

	// metric2-a
	RPCRequestsTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: svcName,
		Subsystem: "dataflow",
		Name:      "requests_total",
		Help:      "Total number of RPC requests",
	}, []string{"subject"})

	// metric2-b
	RPCErrorsTotal := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: svcName,
		Subsystem: "dataflow",
		Name:      "requests_total_errors",
		Help:      "Total number of RPC errors",
	}, []string{"err_code", "subject"})

	// metric3
	RPCRequestsInflight := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: svcName,
		Subsystem: "dataflow",
		Name:      "requests_in_flight",
		Help:      "The number of inflight RPC requests being handled at the same time",
	}, []string{"subject"})

	// metric1
	RPCResponseSecond := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: svcName,
		Subsystem: "dataflow",
		Name:      "request_duration_seconds",
		Help:      "Histogram of response time for RPC in seconds",
		Buckets:   []float64{0.05, 0.2, 0.4, 0.6, 0.8, 1, 5, 10, 30}, // 50 ms ~ 30 s
	}, []string{"err_code", "subject"})

	return func(next dataflow.HandleFunc) dataflow.HandleFunc {
		return func(msg *dataflow.Message, dep any) error {
			c := msg.RawInfra.(*fiber.Ctx)
			subject := msg.Subject

			// metric1
			start := time.Now()

			// metric2-a
			RPCRequestsTotal.WithLabelValues(subject).Inc()

			// metric3
			RPCRequestsInflight.WithLabelValues(subject).Add(1)

			err := next(msg, dep)

			// metric3
			RPCRequestsInflight.WithLabelValues(subject).Add(-1)

			// metric2-b
			errCode := strconv.Itoa(FiberMetadata.GetErrorCode(c))
			if errCode != "0" {
				RPCErrorsTotal.WithLabelValues(errCode, subject).Inc()
			}

			// metric1
			duration := time.Since(start).Seconds()
			span := trace.SpanFromContext(msg.Ctx)
			traceId := span.SpanContext().TraceID()
			if traceId.IsValid() {
				traceLabels := prometheus.Labels{"trace_id": traceId.String()}
				RPCResponseSecond.WithLabelValues(errCode, subject).(prometheus.ExemplarObserver).ObserveWithExemplar(duration, traceLabels)
			} else {
				RPCResponseSecond.WithLabelValues(errCode, subject).Observe(duration)
			}

			return err
		}
	}
}

func dataflowO11YLogger() dataflow.Middleware {
	return func(next dataflow.HandleFunc) dataflow.HandleFunc {
		return func(msg *dataflow.Message, dep any) error {
			c := msg.RawInfra.(*fiber.Ctx)
			ctx := c.UserContext()

			logger := pkg.Logger().CtxGetLogger(ctx).With(
				slog.Any("dataflow", slog.GroupValue(
					slog.String("subject", msg.Subject),
				)),
			)

			c.SetUserContext(pkg.Logger().CtxWithLogger(ctx, logger))
			msg.RawInfra = c

			err := next(msg, dep)

			errCode := strconv.Itoa(FiberMetadata.GetErrorCode(c))
			if errCode != "0" {
				logger.Error("dataflow failed", slog.String("err_code", errCode))
			} else {
				logger.Info("dataflow succeeded")
			}

			return err
		}
	}
}
