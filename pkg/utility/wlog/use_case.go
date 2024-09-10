package wlog

import (
	"os"
)

func LoggerWhenDebug() *Logger {
	const debug = -4
	conf := &LoggerConfig{
		AddSource:  true,
		JsonFormat: false,
		Level_:     debug,
	}
	logger := NewLogger(os.Stdout, conf)
	return logger
}

func LoggerWhenGoTest(source bool) *Logger {
	const info = 0
	conf := &LoggerConfig{
		AddSource:  source,
		JsonFormat: false,
		Level_:     info,
	}
	logger := NewLogger(os.Stdout, conf)
	return logger
}

func LoggerWhenContinuousIntegration() *Logger {
	const warn = 4
	conf := &LoggerConfig{
		AddSource:  true,
		JsonFormat: false,
		Level_:     warn,
	}
	logger := NewLogger(os.Stdout, conf)
	return logger
}
