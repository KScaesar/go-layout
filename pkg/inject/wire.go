//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/adapters"
	"github.com/KScaesar/go-layout/pkg/adapters/datastore"
	"github.com/KScaesar/go-layout/pkg/app"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func NewInfra(conf *pkg.Config) (*Infra, error) {
	panic(wire.Build(
		wire.FieldsOf(new(*pkg.Config),
			"MySql",
			"Redis",
		),

		adapters.NewMySqlGorm,
		adapters.NewRedis,

		wire.Struct(new(Infra), "*"),
	))
}

type Infra struct {
	MySql *gorm.DB
	Redis *redis.Client
}

//

func NewService(conf *pkg.Config, infra *Infra) *Service {
	panic(wire.Build(
		wire.FieldsOf(new(*Infra),
			"MySql",
			"Redis",
		),
		utility.NewGormEasyTransaction,
		utility.NewGormTransaction,

		datastore.NewUserMySQL,
		datastore.NewUserRedis,
		datastore.NewUserRepository,
		wire.Bind(new(app.UserRepository), new(*datastore.UserRepository)),

		app.NewUserUseCase,
		wire.Bind(new(app.UserService), new(*app.UserUseCase)),

		wire.Struct(new(Service), "*"),
	))
}

type Service struct {
	utility.Transaction
	utility.EasyTransaction
	app.UserService
}

//go:generate wire gen

// https://github.com/google/wire/tree/main/docs
// https://github.com/google/wire/tree/main/_tutorial
