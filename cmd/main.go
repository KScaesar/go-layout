package main

import (
	"log/slog"
	"os"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/inject"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func init() {
	logger := wlog.NewLoggerWhenNormalRun(true)
	logger.Logger = logger.With(slog.Any("version", pkg.Version()))
	pkg.SetDefaultLogger(logger)
}

func main() {
	var err error

	conf := configs.MustLoadConfig("./configs/example.conf", pkg.DefaultLogger().Logger)

	// Init is required before get default global variables
	pkg.Init(conf)
	logger := pkg.DefaultLogger()
	shutdown := pkg.DefaultShutdown()

	go shutdown.Serve()
	defer func() {
		shutdown.Notify(err)
	}()

	if err = utility.ServeObservability(
		pkg.Version().ServiceName,
		&conf.O11Y,
		logger.Logger,
		shutdown,
	); err != nil {
		logger.Error("serve o11y failed", slog.Any("err", err))
		os.Exit(1)
	}

	infra, Err := inject.NewInfra(conf)
	if Err != nil {
		err = Err
		logger.Error("create infra failed", slog.Any("err", Err))
		os.Exit(1)
	}

	svc := inject.NewService(conf, infra)
	mux := inject.NewHttpMux(conf, infra.MySql, svc)
	inject.ServeApiServer(conf.Http.Port, mux)

	<-shutdown.WaitChannel()
}
