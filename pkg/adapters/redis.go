package adapters

import (
	"context"
	"fmt"

	"github.com/KScaesar/go-layout/configs"
	"github.com/KScaesar/go-layout/pkg"
	"github.com/redis/go-redis/v9"
)

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
