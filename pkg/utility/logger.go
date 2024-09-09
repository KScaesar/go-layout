package utility

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/context"
)

//

type logKey struct{}

func CtxWithLogger(ctx context.Context, v *slog.Logger) context.Context {
	return context.WithValue(ctx, logKey{}, v)
}

func CtxGetLogger(ctx context.Context) (logger *slog.Logger) {
	v, ok := ctx.Value(logKey{}).(*slog.Logger)
	if !ok {
		return DefaultLogger().Logger
	}
	return v
}

//

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

func NewWrapLogger(w io.Writer, conf *LoggerConfig) *WrapLogger {
	lvl := &slog.LevelVar{}
	lvl.Set(conf.Level())

	pool := NewPool(func() *bytes.Buffer {
		return &bytes.Buffer{}
	})
	stdReplace := func(groups []string, a slog.Attr) slog.Attr {
		if !conf.JsonFormat && a.Key == slog.SourceKey {
			src, ok := a.Value.Any().(*slog.Source)
			if !ok {
				return a
			}

			buf := pool.Get()
			defer func() {
				buf.Reset()
				pool.Put(buf)
			}()

			buf.WriteString(src.File)
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(src.Line))
			a.Value = slog.StringValue(buf.String())
		}
		if a.Value.Kind() == slog.KindTime {
			a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
		}
		if a.Value.Kind() == slog.KindDuration {
			a.Value = slog.StringValue(a.Value.Duration().String())
		}
		return a
	}

	var handler slog.Handler
	if conf.JsonFormat {
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   conf.AddSource,
			Level:       lvl,
			ReplaceAttr: stdReplace,
		})
	} else {
		noColor := true
		fd, ok := w.(interface{ Fd() uintptr })
		if ok {
			noColor = !isatty.IsTerminal(fd.Fd())
		}
		handler = tint.NewHandler(w, &tint.Options{
			AddSource:   conf.AddSource,
			Level:       lvl,
			ReplaceAttr: stdReplace,
			TimeFormat:  time.RFC3339,
			NoColor:     noColor,
		})
	}

	return &WrapLogger{
		lvl:    lvl,
		Logger: slog.New(handler),
	}
}

type WrapLogger struct {
	lvl *slog.LevelVar
	*slog.Logger
}

func (l *WrapLogger) SetLevel(lvl slog.Level) {
	l.lvl.Set(lvl)
}

func (l WrapLogger) SetStdDefaultLogger() {
	slog.SetDefault(l.Logger)
}

func (l WrapLogger) SetStdDefaultLevel() {
	slog.SetLogLoggerLevel(l.lvl.Level())
}

//

func LoggerWhenDebug() *WrapLogger {
	const debug = -4
	conf := &LoggerConfig{
		AddSource:  true,
		JsonFormat: false,
		Level_:     debug,
	}
	logger := NewWrapLogger(os.Stdout, conf)
	return logger
}

func LoggerWhenGoTest() *WrapLogger {
	const warn = 4
	conf := &LoggerConfig{
		AddSource:  true,
		JsonFormat: false,
		Level_:     warn,
	}
	logger := NewWrapLogger(os.Stdout, conf)
	return logger
}

var defaultLogger = LoggerWhenGoTest()

func DefaultLogger() *WrapLogger {
	return defaultLogger
}

func SetDefaultLogger(logger *WrapLogger) {
	defaultLogger = logger
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
		sloggin.NewWithConfig(DefaultLogger().Logger, config),

		func(c *gin.Context) {
			ctx := c.Request.Context()

			reqId := c.Writer.Header().Get(sloggin.RequestIDHeaderKey)
			requestAttributes := []slog.Attr{
				slog.String("method", c.Request.Method),
				slog.String("path", c.Request.URL.Path),
			}
			logger := DefaultLogger().With(
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
