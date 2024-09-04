package app

import (
	"github.com/KScaesar/go-layout/pkg/utility"
)

func RegisterUser(req *RegisterUserRequest) (*User, error) {
	user := User{}
	return &user, nil
}

type User struct {
	ByUpdated utility.MapData `gorm:"-"`
}

func (user *User) ResetPassword(req *ResetUserPasswordRequest) error {
	return nil
}

func (user *User) UpdateInfo(req *UpdateUserInfoRequest) error {
	return nil
}
