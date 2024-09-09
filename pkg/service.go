package pkg

import (
	"log/slog"
	"runtime/debug"
	"sync/atomic"
)

var (
	commit string
)

var defaultService atomic.Pointer[service]

func Service() service {
	return *(defaultService.Load())
}

func newService(name string) *service {
	Commit := commit
	if Commit == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					Commit = setting.Value
				}
			}
		}
	}

	return &service{
		Name:   name,
		Commit: Commit,
	}
}

type service struct {
	Name   string `json:"name"`
	Commit string `json:"commit"`
}

func (svc service) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", svc.Name),
		slog.String("commit", svc.Commit),
	)
}
