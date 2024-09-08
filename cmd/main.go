package main

import (
	"log/slog"
	"os"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/inject"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func main() {
	svc := "CRM"

	utility.LoggerWhenDebug()
	conf := configs.MustLoadConfig("./configs/example.conf")
	utility.InitLogger(svc, &conf.Logger)

	slog.Default().Debug("print config", slog.Any("conf", conf))

	err := utility.ServeObservability(svc, &conf.O11Y)
	if err != nil {
		slog.Default().Error("serve o11y fail", slog.Any("err", err))
		os.Exit(1)
	}

	infra, err := inject.NewInfra(conf)
	if err != nil {
		slog.Default().Error("create infra fail", slog.Any("err", err))
		os.Exit(1)
	}

	service := inject.NewService(svc, conf, infra)
	mux := inject.NewHttpMux(conf, infra.MySql, service)
	inject.ServeApiServer(conf.Http.Port, mux)

	slog.Default().Info("service start")
	utility.DefaultShutdown.Serve()
}
