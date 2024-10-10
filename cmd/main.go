package main

import (
	"log/slog"
	"os"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/inject"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/KScaesar/go-layout/pkg/utility/wlog"
)

func init() {
	logger := wlog.NewStderrLoggerWhenNormal(true)
	logger.Logger = logger.With(slog.Any("version", pkg.Version()))
	pkg.Logger().PointToNew(logger)
}

func main() {
	var err error
	defer func() {
		if err != nil {
			os.Exit(1)
		}
	}()

	conf := pkg.MustLoadConfig()

	// Init is required before get default global variables
	logWriter := pkg.Init(conf)
	defer logWriter.Close()

	shutdown := pkg.Shutdown()
	defer func() {
		if err != nil {
			shutdown.Notify(err)
			<-shutdown.WaitChannel()
		}
	}()

	infra, err := inject.NewInfra(conf)
	if err != nil {
		return
	}
	svc := inject.NewService(conf, infra)
	mux := inject.NewFiberRouter(conf, infra.MySql, svc)
	// mux := inject.NewGinRouter(conf, infra.MySql, svc)

	// server start
	utility.ServeO11YMetric(conf.O11Y.Port, shutdown, pkg.Logger().Logger)
	inject.ServeFiber(conf.Http.Port, conf.Http.Debug, mux)
	// inject.ServeGin(conf.Http.Port, mux)

	<-shutdown.WaitChannel()
}
