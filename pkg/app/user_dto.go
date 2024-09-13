package app

import (
	"github.com/KScaesar/go-layout/pkg/utility/dataflow"
)

// write

type RegisterUserRequest struct {
}

type UpdateUserInfoRequest struct {
}

type UpdateUserPasswordRequest struct {
}

type ResetUserPasswordRequest struct {
}

type DeleteUserRequest struct {
}

// event

func NewRegisteredUserEvent(user *User) *dataflow.Message {
	return dataflow.NewBodyEgress("user.registered", &RegisteredUserEvent{})
}

type RegisteredUserEvent struct {
}

// read

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type QueryMultiUserRequest struct {
}

type MultiUserResponse []UserResponse

func ConvertUserResponse(user *User) UserResponse {
	return UserResponse{}
}

type UserResponse struct {
}
