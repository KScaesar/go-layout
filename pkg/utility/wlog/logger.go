package wlog

import (
	"io"
	"log/slog"
	"sync"
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

func NewLogger(w io.Writer, conf *Config, formats ...FormatFunc) *Logger {
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
	mu  sync.Mutex
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

func (l *Logger) Level() slog.Level {
	return l.lvl.Level()
}

func (l *Logger) SetLevel(lvl slog.Level) {
	l.lvl.Set(lvl)
}

func (l *Logger) SetStdDefaultLevel() {
	l.mu.Lock()
	defer l.mu.Unlock()
	slog.SetLogLoggerLevel(l.lvl.Level())
}

// SetStdDefaultLogger
// 將標準庫的預設值, 以我方的物件為基準, 控制 logger 行為
func (l *Logger) SetStdDefaultLogger() {
	l.mu.Lock()
	defer l.mu.Unlock()
	slog.SetDefault(l.Logger)
}

// PointToNew
// 通過改變指標的指向, 讓所有引用此指標的其他元件獲得最新的狀態
//
// 以下是命名方式的分類
//
// 改變指標本身:
// Replace, Set
//
// 改變指標的指向:
// Redirect, PointTo
func (l *Logger) PointToNew(new *Logger) {
	l.mu.Lock()
	defer l.mu.Unlock()

	*l.lvl = *(new.lvl)
	*l.Logger = *(new.Logger)
}
