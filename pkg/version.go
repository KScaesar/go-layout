package pkg

import (
	"log/slog"
	"runtime/debug"
)

var (
	commit  string
	release string
)

var defaultVersion = newVersion("CRM")

func Version() version {
	return defaultVersion
}

func newVersion(name string) version {
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

	return version{
		ServiceName: name,
		Commit:      Commit,
		Release:     release,
		GoVersion:   goVersion,
	}
}

type version struct {
	ServiceName string `json:"name"`
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
