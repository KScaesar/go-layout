package pkg

import (
	"context"
	"log/slog"
	"sync/atomic"

	"github.com/KScaesar/go-layout/pkg/utility"
)

var defaultShutdown atomic.Pointer[utility.Shutdown]

func DefaultShutdown() *utility.Shutdown {
	return defaultShutdown.Load()
}

func SetDefaultShutdown(l *slog.Logger) {
	defaultShutdown.Store(NewShutdown(l))
}

func NewShutdown(l *slog.Logger) *utility.Shutdown {
	shutdown := utility.NewShutdown(context.Background(), 0)
	shutdown.Logger = l
	return shutdown
}
