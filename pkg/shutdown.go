package pkg

import (
	"context"
	"log/slog"
	"sync/atomic"

	"github.com/KScaesar/go-layout/pkg/utility"
)

var defaultShutdown atomic.Pointer[utility.Shutdown]

func Shutdown() *utility.Shutdown {
	return defaultShutdown.Load()
}

func SetShutdown(shutdown *utility.Shutdown) {
	defaultShutdown.Store(shutdown)
}

func NewShutdownWhenDefault(l *slog.Logger, waitSeconds int) *utility.Shutdown {
	shutdown := utility.NewShutdown(context.Background(), waitSeconds)
	shutdown.Logger = l
	return shutdown
}
