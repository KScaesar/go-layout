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
	pkg.Logger().PointToNew(logger)
}

func main() {
	conf := configs.MustLoadConfig(pkg.Logger().Logger)

	// Init is required before get default global variables
	pkg.Init(conf)
	logger := pkg.Logger()
	shutdown := pkg.Shutdown()
	go shutdown.Serve()

	logger.Debug("show config", slog.Any("conf", conf))
	// pkg.ErrorRegistry.ShowErrors()

	var err error
	defer func() {
		if err != nil {
			shutdown.Notify(err)
			<-shutdown.WaitChannel()
			os.Exit(1)
		}
	}()

	if err = utility.ServeObservability(
		pkg.Version().ServiceName,
		&conf.O11Y,
		logger.Logger,
		shutdown,
	); err != nil {
		logger.Error("serve o11y failed", slog.Any("err", err))
		return
	}

	infra, Err := inject.NewInfra(conf)
	if Err != nil {
		err = Err
		logger.Error("create infra failed", slog.Any("err", Err))
		return
	}

	svc := inject.NewService(conf, infra)
	mux := inject.NewGinRouter(conf, infra.MySql, svc)
	inject.ServeGin(conf.Http.Port, mux)

	<-shutdown.WaitChannel()
}
