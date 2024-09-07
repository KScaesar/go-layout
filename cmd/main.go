package main

import (
	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/inject"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func main() {
	conf := configs.MustLoadConfig("./configs/example.conf")
	utility.Init(conf.Logger.AddSource, conf.Logger.JsonFormat)

	infra, err := inject.NewInfra(conf)
	if err != nil {
		panic(err)
	}

	svc := inject.NewService("CRM", &conf.Biz, infra)
	mux := inject.NewHttpMux(conf, svc)
	inject.ServeApiServer(conf.Http.Port, mux)

	if conf.O11Y.Enable {
		utility.ServeObservability(conf.O11Y.MetricPort_())
	}

	utility.DefaultShutdown.Serve()
}
