package pkg

import (
	"log/slog"
	"runtime/debug"
	"sync/atomic"
)

var (
	commit  string
	release string
)

var defaultService atomic.Pointer[service]

func Service() service {
	return *(defaultService.Load())
}

func newService(name string) *service {
	Commit := commit
	goVersion := ""
	if Commit == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			goVersion = info.GoVersion

			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					Commit = setting.Value[:8]
				}
			}
		}
	}

	return &service{
		Name:      name,
		Commit:    Commit,
		Release:   release,
		GoVersion: goVersion,
	}
}

type service struct {
	Name      string `json:"name"`
	Commit    string `json:"commit"`
	Release   string `json:"release"`
	GoVersion string `json:"go_version"`
}

func (svc service) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", svc.Name),
		slog.String("commit", svc.Commit),
		slog.String("release", svc.Release),
		slog.String("go_version", svc.GoVersion),
	)
}
