package wlog

import (
	"io"
	"os"
)

// NewSmartLogger 依照輸入值判斷使用哪一種類型的 Logger
func NewSmartLogger(filename string, conf *Config) (logger *Logger, w io.WriteCloser, err error) {
	if filename != "" {
		w, err = NewRotateWriter(filename, -1)
		if err != nil {
			return nil, nil, err
		}
		logger = NewFileLogger(w, conf)
	} else {
		w = os.Stderr
		logger = NewConsoleLogger(w, conf)
	}
	return logger, w, nil
}

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
