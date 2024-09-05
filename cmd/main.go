package main

import (
	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/inject"
)

func main() {
	conf := configs.MustLoadConfig("./configs/example.conf")
	infra, err := inject.NewInfra(conf)
	if err != nil {
		panic(err)
	}

	svc := inject.NewService(infra)
	mux := inject.NewHttpMux(conf.Http.Debug, svc)
	server := inject.NewHttpServer(conf.Http.Port, mux)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
