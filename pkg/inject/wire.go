//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/adapters"
	"github.com/KScaesar/go-layout/pkg/adapters/database"
	"github.com/KScaesar/go-layout/pkg/app"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// https://github.com/google/wire/tree/main/docs
// https://github.com/google/wire/tree/main/_tutorial

//go:generate wire gen

func NewInfra(conf *configs.Config) (*Infra, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*configs.Config),
			"MySql",
			"Redis",
		),

		adapters.NewMySql,
		adapters.NewRedis,

		wire.Struct(new(Infra), "*"),
	))
}

type Infra struct {
	MySql *gorm.DB
	Redis *redis.Client
}

func NewService(infra *Infra) *Service {
	panic(wire.Build(
		wire.FieldsOf(new(*Infra),
			"MySql",
			"Redis",
		),

		database.NewUserMySQL,
		database.NewUserRedis,
		database.NewUserRepository,
		wire.Bind(new(app.UserRepository), new(*database.UserRepository)),

		app.NewUserUseCase,
		wire.Bind(new(app.UserService), new(*app.UserUseCase)),

		wire.Struct(new(Service), "*"),
	))
}

type Service struct {
	app.UserService
}