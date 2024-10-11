package pkg

import (
	"log/slog"
	"runtime/debug"
)

var (
	release      string
	release_date string
)

var defaultVersion = newVersion()

func Version() version {
	return defaultVersion
}

func newVersion() version {
	goVersion := ""
	if info, ok := debug.ReadBuildInfo(); ok {
		goVersion = info.GoVersion
	}

	return version{
		ServiceName: "Service",
		Release:     release,
		ReleaseDate: release_date,
		GoVersion:   goVersion,
	}
}

type version struct {
	ServiceName string `json:"svc"`
	Release     string `json:"release"`
	ReleaseDate string `json:"release_date"`
	GoVersion   string `json:"go_version"`
}

func (svc version) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("svc", svc.ServiceName),
		slog.String("release", svc.Release),
		slog.String("release_date", svc.ReleaseDate),
		slog.String("go_version", svc.GoVersion),
	)
}
