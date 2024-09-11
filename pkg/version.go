package pkg

import (
	"log/slog"
	"runtime/debug"
)

var (
	release string
)

var defaultVersion = newVersion()

func Version() version {
	return defaultVersion
}

func newVersion() version {
	commit := ""
	goVersion := ""
	if info, ok := debug.ReadBuildInfo(); ok {
		goVersion = info.GoVersion

		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit = setting.Value
			}
		}
	}

	if len(commit) >= 8 {
		commit = commit[:8]
	}

	return version{
		ServiceName: "CRM",
		Commit:      commit,
		Release:     release,
		GoVersion:   goVersion,
	}
}

type version struct {
	ServiceName string `json:"svc"`
	Commit      string `json:"commit"`
	Release     string `json:"release"`
	GoVersion   string `json:"go_version"`
}

func (svc version) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("svc", svc.ServiceName),
		slog.String("commit", svc.Commit),
		slog.String("release", svc.Release),
		slog.String("go_version", svc.GoVersion),
	)
}
