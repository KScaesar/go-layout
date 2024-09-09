package adapters

import (
	"context"
	"net/http"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg/utility"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

func NewMySqlGorm(conf *configs.MySql) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if conf.Debug {
		db = db.Debug()
	}

	utility.DefaultShutdown().AddPriorityShutdownAction(2, "mysql", func() error {
		stdDB, err := db.DB()
		if err != nil {
			return err
		}
		return stdDB.Close()
	})
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
		return nil, err
	}

	utility.DefaultShutdown().AddPriorityShutdownAction(2, "redis", client.Close)
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
