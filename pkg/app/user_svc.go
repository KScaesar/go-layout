package app

import (
	"context"

	"github.com/KScaesar/go-layout/pkg/utility/dataflow"
)

//go:generate mockgen -typed -package=app -destination=user_svc_mock.go -source=user_svc.go

type UserRepository interface {
	LockUserById(ctx context.Context, userId string) (User, error)
	CreteUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, user *User) error

	QueryUserById(ctx context.Context, userId string) (UserResponse, error)
	LoginUser(ctx context.Context, req *LoginUserRequest) (UserResponse, error)
	QueryMultiUserByFilter(ctx context.Context, filter *QueryMultiUserRequest) (MultiUserResponse, error)
}

type UserService interface {
	RegisterUser(ctx context.Context, req *RegisterUserRequest) error
	UpdateUserInfo(ctx context.Context, userId string, req *UpdateUserInfoRequest) error
	UpdateUserPassword(ctx context.Context, req *UpdateUserPasswordRequest) error
	ResetUserPassword(ctx context.Context, req *ResetUserPasswordRequest) error
	DeleteUser(ctx context.Context, req *DeleteUserRequest) error

	LoginUser(ctx context.Context, req *LoginUserRequest) (UserResponse, error)
	QueryMultiUser(ctx context.Context, filter *QueryMultiUserRequest) (MultiUserResponse, error)
}

func NewUserUseCase(userRepo UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

type UserUseCase struct {
	userRepo UserRepository
	bus      dataflow.Producer
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, req *RegisterUserRequest) error {
	user, err := RegisterUser(req)
	if err != nil {
		return err
	}

	err = uc.userRepo.CreteUser(ctx, user)
	if err != nil {
		return err
	}

	event := NewRegisteredUserEvent(user)
	return uc.bus.SendWithCtx(ctx, event)
}

func (uc *UserUseCase) UpdateUserInfo(ctx context.Context, userId string, req *UpdateUserInfoRequest) error {
	user, err := uc.userRepo.LockUserById(ctx, userId)
	if err != nil {
		return err
	}

	err = user.UpdateInfo(req)
	if err != nil {
		return err
	}

	return uc.userRepo.UpdateUser(ctx, &user)
}

func (uc *UserUseCase) UpdateUserPassword(ctx context.Context, req *UpdateUserPasswordRequest) error {
	// TODO implement me
	panic("implement me")
}

func (uc *UserUseCase) ResetUserPassword(ctx context.Context, req *ResetUserPasswordRequest) error {
	// TODO implement me
	panic("implement me")
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, req *DeleteUserRequest) error {
	// TODO implement me
	panic("implement me")
}

func (uc *UserUseCase) LoginUser(ctx context.Context, req *LoginUserRequest) (UserResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (uc *UserUseCase) QueryMultiUser(ctx context.Context, filter *QueryMultiUserRequest) (MultiUserResponse, error) {
	// TODO implement me
	panic("implement me")
}
