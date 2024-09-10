package main

import (
	"log/slog"
	"os"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/inject"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func init() {
	logger := utility.LoggerWhenDebug()
	logger.Logger = logger.With(slog.Any("version", pkg.DefaultVersion()))
	pkg.SetDefaultLogger(logger)
}

// Init initializes the necessary default global variables
func Init(conf *configs.Config) {
	writer := os.Stdout
	logger := utility.NewWrapLogger(writer, &conf.Logger)
	logger.Logger = logger.With(slog.String("svc", pkg.DefaultVersion().ServiceName))
	logger.SetStdDefaultLevel()
	logger.SetStdDefaultLogger()
	pkg.SetDefaultLogger(logger)
	pkg.SetDefaultShutdown(pkg.NewShutdown(logger.Logger, 0))
	return
}

func main() {
	conf := configs.MustLoadConfig("./configs/example.conf")
	Init(conf)
	logger := pkg.DefaultLogger()

	if err := utility.ServeObservability(
		pkg.DefaultVersion().ServiceName,
		&conf.O11Y,
		pkg.DefaultLogger().Logger,
		pkg.DefaultShutdown(),
	); err != nil {
		logger.Error("serve o11y fail", slog.Any("err", err))
		os.Exit(1)
	}

	infra, err := inject.NewInfra(conf)
	if err != nil {
		logger.Error("create infra fail", slog.Any("err", err))
		os.Exit(1)
	}

	service := inject.NewService(conf, infra)
	mux := inject.NewHttpMux(conf, infra.MySql, service)
	inject.ServeApiServer(conf.Http.Port, mux)

	pkg.DefaultShutdown().Serve()
}
