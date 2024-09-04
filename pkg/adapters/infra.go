package adapters

import (
	"context"
	"net/http"

	"github.com/KScaesar/go-layout/configs"
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
	return db, nil
}

func NewRedis(conf *configs.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Address(),
		Username: conf.User,
		Password: conf.Password,
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
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
