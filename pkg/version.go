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

var defaultVersion atomic.Pointer[Version]

func DefaultVersion() Version {
	return *(defaultVersion.Load())
}

func newVersion(name string) *Version {
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

	return &Version{
		ServiceName: name,
		Commit:      Commit,
		Release:     release,
		GoVersion:   goVersion,
	}
}

type Version struct {
	ServiceName string `json:"name"`
	Commit      string `json:"commit"`
	Release     string `json:"release"`
	GoVersion   string `json:"go_version"`
}

func (svc Version) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("svc", svc.ServiceName),
		slog.String("commit", svc.Commit),
		slog.String("release", svc.Release),
		slog.String("go_version", svc.GoVersion),
	)
}
