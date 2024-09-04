// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package inject

import (
	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/adapters"
	"github.com/KScaesar/go-layout/pkg/adapters/database"
	"github.com/KScaesar/go-layout/pkg/app"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func NewInfra(conf *configs.Config) (*Infra, error) {
	mySql := &conf.MySql
	db, err := adapters.NewMySql(mySql)
	if err != nil {
		return nil, err
	}
	redis := &conf.Redis
	client, err := adapters.NewRedis(redis)
	if err != nil {
		return nil, err
	}
	infra := &Infra{
		MySql: db,
		Redis: client,
	}
	return infra, nil
}

func NewService(infra *Infra) *Service {
	db := infra.MySql
	userMySQL := database.NewUserMySQL(db)
	client := infra.Redis
	userRedis := database.NewUserRedis(client)
	userRepository := database.NewUserRepository(userMySQL, userRedis)
	userUseCase := app.NewUserUseCase(userRepository)
	service := &Service{
		UserService: userUseCase,
	}
	return service
}

// wire.go:

type Infra struct {
	MySql *gorm.DB
	Redis *redis.Client
}

type Service struct {
	app.UserService
}