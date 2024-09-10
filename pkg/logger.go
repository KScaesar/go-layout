package pkg

import (
	"sync/atomic"

	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

var defaultLogger atomic.Pointer[wlog.Logger]

func Logger() *wlog.Logger {
	return defaultLogger.Load()
}

func SetLogger(logger *wlog.Logger) {
	defaultLogger.Store(logger)
}
