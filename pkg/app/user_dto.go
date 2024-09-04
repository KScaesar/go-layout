package app

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
