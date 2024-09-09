package utility

import (
	"context"
)

func init() {
	defaultShutdown.Store(NewShutdown(context.Background(), 0))
	defaultLogger.Store(LoggerWhenGoTest())
}
