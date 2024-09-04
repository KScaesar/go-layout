package database

import (
	"context"

	"github.com/KScaesar/go-layout/pkg/app"
	"gorm.io/gorm"
)

func NewUserMySQL(db *gorm.DB) *UserMySQL {
	return &UserMySQL{db: db}
}

type UserMySQL struct {
	db *gorm.DB
}

func (repo *UserMySQL) LockUserById(ctx context.Context, userId string) (app.User, error) {
	// TODO implement me
	panic("implement me")
}

func (repo *UserMySQL) CreteUser(ctx context.Context, user *app.User) error {
	// TODO implement me
	panic("implement me")
}

func (repo *UserMySQL) UpdateUser(ctx context.Context, user *app.User) error {
	// TODO implement me
	panic("implement me")
}

func (repo *UserMySQL) DeleteUser(ctx context.Context, user *app.User) error {
	// TODO implement me
	panic("implement me")
}

func (repo *UserMySQL) QueryUserById(ctx context.Context, userId string) (app.UserResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (repo *UserMySQL) LoginUser(ctx context.Context, req *app.LoginUserRequest) (app.UserResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (repo *UserMySQL) QueryMultiUserByFilter(ctx context.Context, filter *app.QueryMultiUserRequest) (app.MultiUserResponse, error) {
	// TODO implement me
	panic("implement me")
}
