package utility

import (
	"context"
	"log/slog"
	"os"
)

func init() {
	Init(true, false)
}

func Init(addSource bool, jsonFormat bool) {
	opts := &slog.HandlerOptions{AddSource: addSource}

	if jsonFormat {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
	}

	// because shutdown use slog.Default()
	DefaultShutdown = NewShutdown(context.Background(), 0)
}
