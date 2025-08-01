package adapters

import (
	"encoding/json"
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
	"github.com/KScaesar/go-layout/pkg/utility/wfiber"
)

var FiberO11YMetric = wfiber.NewO11YMetric(pkg.Version().ServiceName)

//

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

func NewErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//

func HandleErrorByFiber(c *fiber.Ctx, err error) error {
	logger := pkg.Logger().CtxGetLogger(c.UserContext())
	logger.Error(err.Error())

	errCode, httpStatus, Err, ok := unwrapError(err)
	if !ok {
		logger.Warn("capture unknown error", slog.Any("err", Err))
	}

	FiberMetadata.SetErrorCode(c, errCode)
	errorResponse := NewErrorResponse(errCode, Err.Error())
	body := fiber.Map{"error": errorResponse}
	return c.Status(httpStatus).JSON(body)
}

func unwrapError(err error) (errCode int, httpStatus int, Err error, ok bool) {
	myErr, ok := utility.UnwrapCustomError(err)
	if !ok {
		errByFixing, isFixed := fixUnknownError(err)
		if isFixed {
			myErr, ok = utility.UnwrapCustomError(errByFixing)
			err = errByFixing
		}
	}
	return myErr.ErrorCode(), myErr.HttpStatus(), err, ok
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
	case fiber.StatusMethodNotAllowed:
		return pkg.ErrInvalidHttpMethod, true
	}
	return err, false
}

//

func ParseQueryString(c *fiber.Ctx, req any, logger *slog.Logger) error {
	err := c.QueryParser(req)
	if err != nil {
		logger.Error(err.Error(), slog.Any("cause", ParseQueryString))
		return fmt.Errorf("Parse QueryString: %w", pkg.ErrInvalidParam)
	}
	return nil
}

func ParseJsonBody(c *fiber.Ctx, req any, logger *slog.Logger) error {
	err := json.Unmarshal(c.BodyRaw(), &req)
	if err != nil {
		logger.Error(err.Error(), slog.Any("cause", ParseJsonBody))
		return fmt.Errorf("%v: %w", err.Error(), pkg.ErrInvalidParam)
	}
	return nil
}

//

var dataflowO11YMetric = newDataflowO11YMetric(pkg.Version().ServiceName)

func newDataflowO11YMetric(svcName string) *_dataflowO11YMetric {
	return &_dataflowO11YMetric{
		ResponseSecond: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: svcName,
			Subsystem: "dataflow",
			Name:      "request_duration_seconds",
			Help:      "Histogram of response time for RPC in seconds",
			Buckets:   []float64{0.05, 0.2, 0.4, 0.6, 0.8, 1, 5, 10, 30}, // 50 ms ~ 30 s
		}, []string{"err_code", "subject"}),

		RequestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: svcName,
			Subsystem: "dataflow",
			Name:      "requests_total",
			Help:      "Total number of RPC requests",
		}, []string{"subject"}),
		ErrorsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: svcName,
			Subsystem: "dataflow",
			Name:      "requests_total_errors",
			Help:      "Total number of RPC errors",
		}, []string{"err_code", "subject"}),

		RequestsInflight: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: svcName,
			Subsystem: "dataflow",
			Name:      "requests_in_flight",
			Help:      "The number of inflight RPC requests being handled at the same time",
		}, []string{"subject"}),
	}
}

type _dataflowO11YMetric struct {
	// metric1
	ResponseSecond *prometheus.HistogramVec

	// metric2-a
	RequestsTotal *prometheus.CounterVec
	// metric2-b
	ErrorsTotal *prometheus.CounterVec

	// metric3
	RequestsInflight *prometheus.GaugeVec
}

func (m *_dataflowO11YMetric) Middleware(next dataflow.HandleFunc) dataflow.HandleFunc {
	return func(ingress *dataflow.Message, dep any) error {
		c := ingress.RawInfra.(*fiber.Ctx)
		subject := ingress.Subject

		// metric1
		start := time.Now()

		// metric2-a
		m.RequestsTotal.WithLabelValues(subject).Inc()

		// metric3
		m.RequestsInflight.WithLabelValues(subject).Add(1)

		err := next(ingress, dep)

		// metric3
		m.RequestsInflight.WithLabelValues(subject).Add(-1)

		// metric2-b
		errCode := strconv.Itoa(FiberMetadata.GetErrorCode(c))
		if errCode != "0" {
			m.ErrorsTotal.WithLabelValues(errCode, subject).Inc()
		}

		// metric1
		duration := time.Since(start).Seconds()
		span := trace.SpanFromContext(ingress.Ctx)
		traceId := span.SpanContext().TraceID()
		if traceId.IsValid() {
			traceLabels := prometheus.Labels{"trace_id": traceId.String()}
			m.ResponseSecond.WithLabelValues(errCode, subject).(prometheus.ExemplarObserver).ObserveWithExemplar(duration, traceLabels)
		} else {
			m.ResponseSecond.WithLabelValues(errCode, subject).Observe(duration)
		}

		return err
	}
}

func dataflowO11YLogger() dataflow.Middleware {
	return func(next dataflow.HandleFunc) dataflow.HandleFunc {
		return func(ingress *dataflow.Message, dep any) error {
			c := ingress.RawInfra.(*fiber.Ctx)
			ctx := c.UserContext()

			logger := pkg.Logger().CtxGetLogger(ctx).With(
				slog.Any("dataflow", slog.GroupValue(
					slog.String("subject", ingress.Subject),
				)),
			)

			c.SetUserContext(pkg.Logger().CtxWithLogger(ctx, logger))
			ingress.RawInfra = c

			err := next(ingress, dep)

			errCode := strconv.Itoa(FiberMetadata.GetErrorCode(c))
			logger.Info("dataflow finish", slog.String("err_code", errCode))

			return err
		}
	}
}
