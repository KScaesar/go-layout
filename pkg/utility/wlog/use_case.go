package wlog

import (
	"os"
)

func NewLoggerWhenDebug() *Logger {
	const debug = -4
	conf := &Config{
		AddSource:  true,
		JsonFormat: false,
		Level_:     debug,
	}
	logger := NewLogger(os.Stdout, conf, DefaultFormat()...)
	return logger
}

func NewLoggerWhenNormalRun(source bool) *Logger {
	const info = 0
	conf := &Config{
		AddSource:  source,
		JsonFormat: false,
		Level_:     info,
	}
	logger := NewLogger(os.Stdout, conf, DefaultFormat()...)
	return logger
}

func NewLoggerWhenContinuousIntegration() *Logger {
	const warn = 4
	conf := &Config{
		AddSource:  true,
		JsonFormat: false,
		Level_:     warn,
	}
	logger := NewLogger(os.Stdout, conf, DefaultFormat()...)
	return logger
}
