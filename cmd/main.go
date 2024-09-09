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
	logger.Logger = logger.With(slog.Any("svc", pkg.Service))
	utility.SetDefaultLogger(logger)
}

// Init initializes the necessary default global variables
func Init(conf *configs.Config) {
	writer := os.Stdout
	logger := utility.NewWrapLogger(writer, &conf.Logger)
	logger.Logger = logger.With(slog.String("svc", pkg.Service.Name))
	logger.SetStdDefaultLevel()
	logger.SetStdDefaultLogger()
	utility.SetDefaultLogger(logger)
	return
}

func main() {
	conf := configs.MustLoadConfig("./configs/example.conf")
	Init(conf)
	logger := utility.DefaultLogger()
	logger.Debug("print config", slog.Any("conf", conf))

	err := utility.ServeObservability(pkg.Service.Name, &conf.O11Y)
	if err != nil {
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

	utility.DefaultShutdown().Serve()
}
