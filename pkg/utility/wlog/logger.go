package wlog

import (
	"io"
	"log/slog"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"golang.org/x/net/context"
)

type Config struct {
	AddSource  bool `yaml:"AddSource"`
	JsonFormat bool `yaml:"JsonFormat"`

	// Debug = -4
	// Info  = 0
	// Warn  = 4
	// Error = 8
	Level_ int `yaml:"Level"`
}

func (conf Config) Level() slog.Level {
	return slog.Level(conf.Level_)
}

func NewLogger(w io.Writer, conf *Config) *Logger {
	formats := []FormatFunc{
		FormatSource(conf.JsonFormat),
		FormatKindTime(),
		FormatKindDuration(),
		FormatTypeFunc(),
	}

	stdReplace := func(groups []string, a slog.Attr) slog.Attr {
		for _, format := range formats {
			attr, ok := format(groups, a)
			if ok {
				return attr
			}
		}
		return a
	}

	lvl := &slog.LevelVar{}
	lvl.Set(conf.Level())

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

	return &Logger{
		lvl:    lvl,
		Logger: slog.New(handler),
	}
}

type Logger struct {
	lvl *slog.LevelVar
	*slog.Logger
}

func (l *Logger) CtxWithLogger(ctx context.Context, v *slog.Logger) context.Context {
	return context.WithValue(ctx, l.Logger, v)
}

func (l *Logger) CtxGetLogger(ctx context.Context) (logger *slog.Logger) {
	v, ok := ctx.Value(l.Logger).(*slog.Logger)
	if !ok {
		return l.Logger
	}
	return v
}

func (l Logger) Level() slog.Level {
	return l.lvl.Level()
}

func (l *Logger) SetLevel(lvl slog.Level) {
	l.lvl.Set(lvl)
}

func (l Logger) SetStdDefaultLogger() {
	slog.SetDefault(l.Logger)
}

func (l Logger) SetStdDefaultLevel() {
	slog.SetLogLoggerLevel(l.lvl.Level())
}
