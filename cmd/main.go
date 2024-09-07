package main

import (
	"log/slog"
	"os"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/inject"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func main() {
	conf := configs.MustLoadConfig("./configs/example.conf")

	opts := &slog.HandlerOptions{AddSource: true}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))

	infra, err := inject.NewInfra(conf)
	if err != nil {
		panic(err)
	}

	svc := inject.NewService("CRM", &conf.Biz, infra)
	mux := inject.NewHttpMux(conf, svc)
	inject.ServeHttpServer(conf.Http.Port, mux)

	if conf.O11Y.Enable {
		utility.ServeObservability(conf.O11Y.MetricPort_())
	}

	utility.DefaultShutdown.Serve()
}
