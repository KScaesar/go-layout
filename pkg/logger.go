package pkg

import (
	"sync/atomic"

	"github.com/KScaesar/go-layout/pkg/utility"
)

var defaultLogger atomic.Pointer[utility.WrapLogger]

func DefaultLogger() *utility.WrapLogger {
	return defaultLogger.Load()
}

func SetDefaultLogger(logger *utility.WrapLogger) {
	defaultLogger.Store(logger)
}
