package wlog

import (
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
	"golang.org/x/net/context"
)

type Config struct {
	// Debug = -4
	// Info  = 0
	// Warn  = 4
	// Error = 8
	Level *int `yaml:"Level"`

	AddSource  *bool `yaml:"AddSource"`
	JsonFormat *bool `yaml:"JsonFormat"`
	NoColor    *bool `yaml:"NoColor"`

	Formats  []FormatFunc   `yaml:"-" json:"-"`
	LevelVar *slog.LevelVar `yaml:"-" json:"-"`
}

func (conf *Config) defaultValue() {
	if conf.AddSource == nil {
		conf.SetAddSource(false)
	}

	if conf.JsonFormat == nil {
		conf.SetJsonFormat(false)
	}

	if conf.NoColor == nil {
		conf.SetNoColor(true)
	}

	if conf.Formats == nil {
		conf.SetFormats(DefaultFormats...)
	}

	if conf.LevelVar == nil {
		if conf.Level == nil {
			info := 0
			conf.Level = &info
		}
		conf.SetLevelVar(*conf.Level)
	}
}

func (conf *Config) SetAddSource(add bool) *Config {
	conf.AddSource = &add
	return conf
}

func (conf *Config) SetJsonFormat(json bool) *Config {
	conf.JsonFormat = &json
	return conf
}

func (conf *Config) SetNoColor(noColor bool) *Config {
	conf.NoColor = &noColor
	return conf
}

func (conf *Config) SetFormats(formats ...FormatFunc) *Config {
	conf.Formats = formats
	return conf
}

func (conf *Config) SetLevelVar(level int) *Config {
	conf.Level = &level
	lvl := &slog.LevelVar{}
	lvl.Set(slog.Level(*conf.Level))
	conf.LevelVar = lvl
	return conf
}

//

func NewHandler(w io.Writer, conf *Config) slog.Handler {
	conf.defaultValue()

	replace := func(groups []string, a slog.Attr) slog.Attr {
		for _, format := range conf.Formats {
			attr, ok := format(groups, a)
			if ok {
				return attr
			}
		}
		return a
	}

	var handler slog.Handler
	if *conf.JsonFormat {
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource:   *conf.AddSource,
			Level:       conf.LevelVar,
			ReplaceAttr: replace,
		})
	} else {
		handler = tint.NewHandler(w, &tint.Options{
			AddSource:   *conf.AddSource,
			Level:       conf.LevelVar,
			ReplaceAttr: replace,
			TimeFormat:  time.RFC3339,
			NoColor:     *conf.NoColor,
		})
	}

	return handler
}

func NewLogger(lvl *slog.LevelVar, handlers ...slog.Handler) *Logger {
	return &Logger{
		lvl:    lvl,
		Logger: slog.New(slogmulti.Fanout(handlers...)),
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
