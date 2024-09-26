package pkg

import (
	"log/slog"
	"runtime/debug"
)

var (
	git_commit   string
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

	if len(git_commit) >= 8 {
		git_commit = git_commit[:8]
	}

	return version{
		ServiceName: "Service",
		GitCommit:   git_commit,
		Release:     release,
		ReleaseDate: release_date,
		GoVersion:   goVersion,
	}
}

type version struct {
	ServiceName string `json:"svc"`
	GitCommit   string `json:"git_commit"`
	Release     string `json:"release"`
	ReleaseDate string `json:"release_date"`
	GoVersion   string `json:"go_version"`
}

func (svc version) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("svc", svc.ServiceName),
		slog.String("git_commit", svc.GitCommit),
		slog.String("release", svc.Release),
		slog.String("release_date", svc.ReleaseDate),
		slog.String("go_version", svc.GoVersion),
	)
}
