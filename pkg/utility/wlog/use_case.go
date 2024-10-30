package wlog

import (
	"io"
	"os"
)

func LoggerFactory(filename string, conf *Config) (logger *Logger, w io.WriteCloser, err error) {
	if filename != "" {
		w, err = NewRotateWriter(filename, -1)
		if err != nil {
			return nil, nil, err
		}
		handler := NewHandler(w, conf)
		logger = NewLogger(conf.LevelVar, handler)
	} else {
		w = os.Stderr
		handler := NewHandler(w, conf)
		logger = NewLogger(conf.LevelVar, handler)
	}
	return logger, w, nil
}

func NewStderrLogger(conf *Config) *Logger {
	handler := NewHandler(os.Stderr, conf)
	logger := NewLogger(conf.LevelVar, handler)
	return logger
}

func NewStderrLoggerWhenNormal(source bool) *Logger {
	conf := &Config{}
	info := 0
	conf.SetAddSource(source).
		SetLevelVar(info).
		SetNoColor(false)
	return NewStderrLogger(conf)
}

func NewStderrLoggerWhenDebug() *Logger {
	conf := &Config{}
	debug := -4
	conf.SetAddSource(true).
		SetLevelVar(debug).
		SetNoColor(false)
	return NewStderrLogger(conf)
}

func NewDiscardLogger() *Logger {
	conf := &Config{}
	handler := NewHandler(io.Discard, conf)
	logger := NewLogger(conf.LevelVar, handler)
	return logger
}
