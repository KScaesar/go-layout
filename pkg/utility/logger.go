package utility

import (
	"log/slog"

	"golang.org/x/net/context"
)

type logKey struct{}

func CtxWithLogger(ctx context.Context, v *slog.Logger) context.Context {
	return context.WithValue(ctx, logKey{}, v)
}

func CtxGetLogger(ctx context.Context) (logger *slog.Logger) {
	logger, ok := ctx.Value(logKey{}).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return logger
}
