package database

import (
	"context"

	"github.com/KScaesar/go-layout/pkg/app"
	"github.com/redis/go-redis/v9"
)

func NewUserRedis(client *redis.Client) *UserRedis {
	return &UserRedis{client: client}
}

type UserRedis struct {
	client *redis.Client
}

func (repo *UserRedis) SetUser(ctx context.Context, key string, resp *app.UserResponse) error {
	// TODO implement me
	panic("implement me")
}

func (repo *UserRedis) SetMultiUser(ctx context.Context, key string, resp *app.MultiUserResponse) error {
	// TODO implement me
	panic("implement me")
}

func (repo *UserRedis) DeleteUser(ctx context.Context, resp *app.User) error {
	// TODO implement me
	panic("implement me")
}

func (repo *UserRedis) QueryUser(ctx context.Context, key string) (app.UserResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (repo *UserRedis) QueryMultiUser(ctx context.Context, key string) (app.MultiUserResponse, error) {
	// TODO implement me
	panic("implement me")
}
