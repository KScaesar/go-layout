package wlog

import (
	"io"
	"os"
)

func NewFileLogger(w io.Writer, conf *Config) *Logger {
	handler := NewHandler(w, true, conf)
	return NewLogger(conf.LevelVar, handler)
}

func NewConsoleLogger(w io.Writer, conf *Config) *Logger {
	handler := NewHandler(w, false, conf)
	return NewLogger(conf.LevelVar, handler)
}

func NewStderrLogger(conf *Config) *Logger {
	return NewConsoleLogger(os.Stderr, conf)
}

func NewStderrLoggerWhenNormal(source bool) *Logger {
	conf := &Config{}
	info := 0
	conf.SetAddSource(source).SetLevelVar(info)
	return NewStderrLogger(conf)
}

func NewStderrLoggerWhenDebug() *Logger {
	conf := &Config{}
	debug := -4
	conf.SetAddSource(true).SetLevelVar(debug)
	return NewStderrLogger(conf)
}

func NewStderrLoggerWhenIntegrationTest() *Logger {
	conf := &Config{}
	warn := 4
	conf.SetAddSource(true).SetLevelVar(warn)
	return NewStderrLogger(conf)
}
