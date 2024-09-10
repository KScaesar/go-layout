package adapters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

func NewMySqlGorm(conf *configs.MySql) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.DSN()), &gorm.Config{
		Logger:                                   nil,
		NowFunc:                                  nil,
		DryRun:                                   false,
		DisableForeignKeyConstraintWhenMigrating: false,
		IgnoreRelationshipsWhenMigrating:         false,
	})
	if err != nil {
		return nil, fmt.Errorf("connect mysql: %w", err)
	}

	stdDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get stdDB: %w", err)
	}

	err = stdDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	if conf.Debug {
		db = db.Debug()
	}

	pkg.Shutdown().AddPriorityShutdownAction(2, "mysql", stdDB.Close)
	return db, nil
}

func NewRedis(conf *configs.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:           conf.Address(),
		Username:       conf.User,
		Password:       conf.Password,
		DB:             0,
		PoolSize:       0,
		MinIdleConns:   5,
		MaxIdleConns:   10,
		MaxActiveConns: 10,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	pkg.Shutdown().AddPriorityShutdownAction(2, "redis", client.Close)
	return client, nil
}

func NewHttpClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConnsPerHost: 5,
	}
	client := &http.Client{
		Transport: transport,
	}
	return client
}

func NewMessageProducer() {

}

func NewMessageConsumer() {

}
