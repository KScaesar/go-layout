package utility

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

type LoggerConfig struct {
	AddSource  bool `yaml:"AddSource"`
	JsonFormat bool `yaml:"JsonFormat"`

	// Debug = -4
	// Info  = 0
	// Warn  = 4
	// Error = 8
	Level_ int `yaml:"Level"`
}

func (l LoggerConfig) Level() slog.Level {
	return slog.Level(l.Level_)
}

//

var loggerLevel = &slog.LevelVar{}

func SetLoggerLevel(lvl slog.Level) {
	loggerLevel.Set(lvl)
	slog.SetLogLoggerLevel(lvl)
}

func GetLoggerLevel() slog.Leveler {
	return loggerLevel
}

//

func NewLoggerHandlerOptions(source bool) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource: source,
		Level:     loggerLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}
			a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			return a
		},
	}
}

//

func InitLogger(svc string, conf *LoggerConfig) {
	SetLoggerLevel(conf.Level())
	opts := NewLoggerHandlerOptions(conf.AddSource)

	var logger *slog.Logger
	if conf.JsonFormat {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	if svc != "" {
		logger = logger.With(slog.String("svc", svc))
	}

	// slog.SetDefault(otelslog.NewLogger("app_or_package_name"))
	slog.SetDefault(logger)
}

func LoggerWhenDebug() {
	debug := -4
	conf := &LoggerConfig{
		AddSource:  true,
		JsonFormat: false,
		Level_:     debug,
	}
	InitLogger("", conf)
}

func LoggerWhenGoTest() {
	warn := 4
	conf := &LoggerConfig{
		AddSource:  true,
		JsonFormat: false,
		Level_:     warn,
	}
	InitLogger("", conf)
}

//

type logKey struct{}

func CtxWithLogger(ctx context.Context, v *slog.Logger) context.Context {
	return context.WithValue(ctx, logKey{}, v)
}

func CtxGetLogger(ctx context.Context) (logger *slog.Logger) {
	logger, ok := ctx.Value(logKey{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}

//

func GinO11YLogger(debug bool, enableTrace bool) []gin.HandlerFunc {
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
		sloggin.NewWithConfig(slog.Default(), config),

		func(c *gin.Context) {
			ctx := c.Request.Context()

			reqId := c.Writer.Header().Get(sloggin.RequestIDHeaderKey)
			requestAttributes := []slog.Attr{
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
			}
			logger := slog.Default().With(
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

			c.Request = c.Request.WithContext(CtxWithLogger(ctx, logger))
			c.Next()
		},
	}
}
