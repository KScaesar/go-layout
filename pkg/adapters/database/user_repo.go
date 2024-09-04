package database

import (
	"context"

	"github.com/KScaesar/go-layout/pkg/app"
	"github.com/KScaesar/go-layout/pkg/utility"
)

func NewUserRepository(userMySQL *UserMySQL, cache *UserRedis) *UserRepository {
	return &UserRepository{
		UserMySQL: userMySQL,
		cache:     cache,
	}
}

type UserRepository struct {
	*UserMySQL
	cache *UserRedis

	singleFlight utility.Singleflight
}

func (repo *UserRepository) UpdateUser(ctx context.Context, user *app.User) error {
	err := repo.UserMySQL.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return repo.cache.DeleteUser(ctx, user)
}

func (repo *UserRepository) DeleteUser(ctx context.Context, user *app.User) error {
	err := repo.UserMySQL.DeleteUser(ctx, user)
	if err != nil {
		return err
	}
	return repo.cache.DeleteUser(ctx, user)
}

func (repo *UserRepository) QueryUserById(ctx context.Context, userId string) (app.UserResponse, error) {
	proxy := utility.ReadProxy[app.UserResponse, func(key string) (app.UserResponse, error), func(string, *app.UserResponse) error]{
		ReadReplica: func(key string) (resp app.UserResponse, err error) {
			return repo.cache.QueryUser(ctx, key)
		},
		ReadPrimary: func(key string) (resp app.UserResponse, err error) {
			return repo.UserMySQL.QueryUserById(ctx, userId)
		},
		WriteReplica: func(key string, resp *app.UserResponse) error {
			return repo.cache.SetUser(ctx, key, resp)
		},
		SingleFlight: &repo.singleFlight,
	}
	return proxy.SafeReadPrimaryAndReplicaNode("userId:" + userId)
}

func (repo *UserRepository) QueryMultiUserByFilter(ctx context.Context, filter *app.QueryMultiUserRequest) (app.MultiUserResponse, error) {
	proxy := utility.ReadProxy[app.MultiUserResponse, func(key string) (app.MultiUserResponse, error), func(string, *app.MultiUserResponse) error]{
		ReadReplica: func(key string) (resp app.MultiUserResponse, err error) {
			return repo.cache.QueryMultiUser(ctx, key)
		},
		ReadPrimary: func(key string) (resp app.MultiUserResponse, err error) {
			return repo.UserMySQL.QueryMultiUserByFilter(ctx, filter)
		},
		WriteReplica: func(key string, resp *app.MultiUserResponse) error {
			return repo.cache.SetMultiUser(ctx, key, resp)
		},
		SingleFlight: &repo.singleFlight,
	}
	return proxy.SafeReadPrimaryNode("xxxKey")
}
