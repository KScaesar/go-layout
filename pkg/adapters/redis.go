package adapters

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"

	"github.com/KScaesar/go-layout/pkg"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func NewRedis(conf *pkg.Redis) (*redis.Client, error) {
	client := redis.NewClient(newRedisOptions(conf, 0))

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	id := fmt.Sprintf("redis(%p)", client)
	pkg.Shutdown().AddPriorityShutdownAction(2, id, client.Close)
	return client, nil
}

func newRedisOptions(conf *pkg.Redis, db int) (opt *redis.Options) {
	switch db {
	default:
		opt = &redis.Options{
			Addr:           conf.Address(),
			Username:       conf.User,
			Password:       conf.Password,
			DB:             db,
			PoolSize:       8,
			MinIdleConns:   4,
			MaxIdleConns:   8,
			MaxActiveConns: 8,
		}
	}
	return opt
}

//

func ConvertErrorFromRedis(err error) error {
	switch {
	case errors.Is(err, redis.Nil):
		return pkg.ErrNotExists
	default:
		return pkg.ErrDatabase
	}
}

//

func GetRedisStringByType[T any](
	// dependency
	client *redis.Client,
	unmarshal utility.Unmarshal,
	logger *slog.Logger,

	// parameter
	ctx context.Context,
	key string,
) (resp T, Err error) {

	bData, err := client.Get(ctx, key).Bytes()
	if err != nil {
		logger.Error("redis get by string", slog.Any("err", err))
		Err = ConvertErrorFromRedis(err)
		return
	}

	err = unmarshal(bData, &resp)
	if err != nil {
		logger.Error(err.Error(), slog.Any("cause", unmarshal))
		Err = fmt.Errorf("unmarshal when redis.client.Get: %w", pkg.ErrSystem)
		return
	}
	return
}
