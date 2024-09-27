// Code generated by MockGen. DO NOT EDIT.
// Source: user_svc.go
//
// Generated by this command:
//
//	mockgen -typed -package=app -destination=user_svc_mock.go -source=user_svc.go
//

// Package app is a generated GoMock package.
package app

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CreteUser mocks base method.
func (m *MockUserRepository) CreteUser(ctx context.Context, user *User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreteUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreteUser indicates an expected call of CreteUser.
func (mr *MockUserRepositoryMockRecorder) CreteUser(ctx, user any) *MockUserRepositoryCreteUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreteUser", reflect.TypeOf((*MockUserRepository)(nil).CreteUser), ctx, user)
	return &MockUserRepositoryCreteUserCall{Call: call}
}

// MockUserRepositoryCreteUserCall wrap *gomock.Call
type MockUserRepositoryCreteUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserRepositoryCreteUserCall) Return(arg0 error) *MockUserRepositoryCreteUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserRepositoryCreteUserCall) Do(f func(context.Context, *User) error) *MockUserRepositoryCreteUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserRepositoryCreteUserCall) DoAndReturn(f func(context.Context, *User) error) *MockUserRepositoryCreteUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// DeleteUser mocks base method.
func (m *MockUserRepository) DeleteUser(ctx context.Context, user *User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserRepositoryMockRecorder) DeleteUser(ctx, user any) *MockUserRepositoryDeleteUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserRepository)(nil).DeleteUser), ctx, user)
	return &MockUserRepositoryDeleteUserCall{Call: call}
}

// MockUserRepositoryDeleteUserCall wrap *gomock.Call
type MockUserRepositoryDeleteUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserRepositoryDeleteUserCall) Return(arg0 error) *MockUserRepositoryDeleteUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserRepositoryDeleteUserCall) Do(f func(context.Context, *User) error) *MockUserRepositoryDeleteUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserRepositoryDeleteUserCall) DoAndReturn(f func(context.Context, *User) error) *MockUserRepositoryDeleteUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LockUserById mocks base method.
func (m *MockUserRepository) LockUserById(ctx context.Context, userId string) (User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockUserById", ctx, userId)
	ret0, _ := ret[0].(User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LockUserById indicates an expected call of LockUserById.
func (mr *MockUserRepositoryMockRecorder) LockUserById(ctx, userId any) *MockUserRepositoryLockUserByIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockUserById", reflect.TypeOf((*MockUserRepository)(nil).LockUserById), ctx, userId)
	return &MockUserRepositoryLockUserByIdCall{Call: call}
}

// MockUserRepositoryLockUserByIdCall wrap *gomock.Call
type MockUserRepositoryLockUserByIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserRepositoryLockUserByIdCall) Return(arg0 User, arg1 error) *MockUserRepositoryLockUserByIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserRepositoryLockUserByIdCall) Do(f func(context.Context, string) (User, error)) *MockUserRepositoryLockUserByIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserRepositoryLockUserByIdCall) DoAndReturn(f func(context.Context, string) (User, error)) *MockUserRepositoryLockUserByIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LoginUser mocks base method.
func (m *MockUserRepository) LoginUser(ctx context.Context, req *LoginUserRequest) (UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", ctx, req)
	ret0, _ := ret[0].(UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockUserRepositoryMockRecorder) LoginUser(ctx, req any) *MockUserRepositoryLoginUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockUserRepository)(nil).LoginUser), ctx, req)
	return &MockUserRepositoryLoginUserCall{Call: call}
}

// MockUserRepositoryLoginUserCall wrap *gomock.Call
type MockUserRepositoryLoginUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserRepositoryLoginUserCall) Return(arg0 UserResponse, arg1 error) *MockUserRepositoryLoginUserCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserRepositoryLoginUserCall) Do(f func(context.Context, *LoginUserRequest) (UserResponse, error)) *MockUserRepositoryLoginUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserRepositoryLoginUserCall) DoAndReturn(f func(context.Context, *LoginUserRequest) (UserResponse, error)) *MockUserRepositoryLoginUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// QueryMultiUserByFilter mocks base method.
func (m *MockUserRepository) QueryMultiUserByFilter(ctx context.Context, filter *QueryMultiUserRequest) (MultiUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryMultiUserByFilter", ctx, filter)
	ret0, _ := ret[0].(MultiUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryMultiUserByFilter indicates an expected call of QueryMultiUserByFilter.
func (mr *MockUserRepositoryMockRecorder) QueryMultiUserByFilter(ctx, filter any) *MockUserRepositoryQueryMultiUserByFilterCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryMultiUserByFilter", reflect.TypeOf((*MockUserRepository)(nil).QueryMultiUserByFilter), ctx, filter)
	return &MockUserRepositoryQueryMultiUserByFilterCall{Call: call}
}

// MockUserRepositoryQueryMultiUserByFilterCall wrap *gomock.Call
type MockUserRepositoryQueryMultiUserByFilterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserRepositoryQueryMultiUserByFilterCall) Return(arg0 MultiUserResponse, arg1 error) *MockUserRepositoryQueryMultiUserByFilterCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserRepositoryQueryMultiUserByFilterCall) Do(f func(context.Context, *QueryMultiUserRequest) (MultiUserResponse, error)) *MockUserRepositoryQueryMultiUserByFilterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserRepositoryQueryMultiUserByFilterCall) DoAndReturn(f func(context.Context, *QueryMultiUserRequest) (MultiUserResponse, error)) *MockUserRepositoryQueryMultiUserByFilterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// QueryUserById mocks base method.
func (m *MockUserRepository) QueryUserById(ctx context.Context, userId string) (UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryUserById", ctx, userId)
	ret0, _ := ret[0].(UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryUserById indicates an expected call of QueryUserById.
func (mr *MockUserRepositoryMockRecorder) QueryUserById(ctx, userId any) *MockUserRepositoryQueryUserByIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryUserById", reflect.TypeOf((*MockUserRepository)(nil).QueryUserById), ctx, userId)
	return &MockUserRepositoryQueryUserByIdCall{Call: call}
}

// MockUserRepositoryQueryUserByIdCall wrap *gomock.Call
type MockUserRepositoryQueryUserByIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserRepositoryQueryUserByIdCall) Return(arg0 UserResponse, arg1 error) *MockUserRepositoryQueryUserByIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserRepositoryQueryUserByIdCall) Do(f func(context.Context, string) (UserResponse, error)) *MockUserRepositoryQueryUserByIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserRepositoryQueryUserByIdCall) DoAndReturn(f func(context.Context, string) (UserResponse, error)) *MockUserRepositoryQueryUserByIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateUser mocks base method.
func (m *MockUserRepository) UpdateUser(ctx context.Context, user *User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserRepositoryMockRecorder) UpdateUser(ctx, user any) *MockUserRepositoryUpdateUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserRepository)(nil).UpdateUser), ctx, user)
	return &MockUserRepositoryUpdateUserCall{Call: call}
}

// MockUserRepositoryUpdateUserCall wrap *gomock.Call
type MockUserRepositoryUpdateUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserRepositoryUpdateUserCall) Return(arg0 error) *MockUserRepositoryUpdateUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserRepositoryUpdateUserCall) Do(f func(context.Context, *User) error) *MockUserRepositoryUpdateUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserRepositoryUpdateUserCall) DoAndReturn(f func(context.Context, *User) error) *MockUserRepositoryUpdateUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// DeleteUser mocks base method.
func (m *MockUserService) DeleteUser(ctx context.Context, req *DeleteUserRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserServiceMockRecorder) DeleteUser(ctx, req any) *MockUserServiceDeleteUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserService)(nil).DeleteUser), ctx, req)
	return &MockUserServiceDeleteUserCall{Call: call}
}

// MockUserServiceDeleteUserCall wrap *gomock.Call
type MockUserServiceDeleteUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserServiceDeleteUserCall) Return(arg0 error) *MockUserServiceDeleteUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserServiceDeleteUserCall) Do(f func(context.Context, *DeleteUserRequest) error) *MockUserServiceDeleteUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserServiceDeleteUserCall) DoAndReturn(f func(context.Context, *DeleteUserRequest) error) *MockUserServiceDeleteUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LoginUser mocks base method.
func (m *MockUserService) LoginUser(ctx context.Context, req *LoginUserRequest) (UserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", ctx, req)
	ret0, _ := ret[0].(UserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockUserServiceMockRecorder) LoginUser(ctx, req any) *MockUserServiceLoginUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockUserService)(nil).LoginUser), ctx, req)
	return &MockUserServiceLoginUserCall{Call: call}
}

// MockUserServiceLoginUserCall wrap *gomock.Call
type MockUserServiceLoginUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserServiceLoginUserCall) Return(arg0 UserResponse, arg1 error) *MockUserServiceLoginUserCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserServiceLoginUserCall) Do(f func(context.Context, *LoginUserRequest) (UserResponse, error)) *MockUserServiceLoginUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserServiceLoginUserCall) DoAndReturn(f func(context.Context, *LoginUserRequest) (UserResponse, error)) *MockUserServiceLoginUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// QueryMultiUser mocks base method.
func (m *MockUserService) QueryMultiUser(ctx context.Context, filter *QueryMultiUserRequest) (MultiUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryMultiUser", ctx, filter)
	ret0, _ := ret[0].(MultiUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryMultiUser indicates an expected call of QueryMultiUser.
func (mr *MockUserServiceMockRecorder) QueryMultiUser(ctx, filter any) *MockUserServiceQueryMultiUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryMultiUser", reflect.TypeOf((*MockUserService)(nil).QueryMultiUser), ctx, filter)
	return &MockUserServiceQueryMultiUserCall{Call: call}
}

// MockUserServiceQueryMultiUserCall wrap *gomock.Call
type MockUserServiceQueryMultiUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserServiceQueryMultiUserCall) Return(arg0 MultiUserResponse, arg1 error) *MockUserServiceQueryMultiUserCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserServiceQueryMultiUserCall) Do(f func(context.Context, *QueryMultiUserRequest) (MultiUserResponse, error)) *MockUserServiceQueryMultiUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserServiceQueryMultiUserCall) DoAndReturn(f func(context.Context, *QueryMultiUserRequest) (MultiUserResponse, error)) *MockUserServiceQueryMultiUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RegisterUser mocks base method.
func (m *MockUserService) RegisterUser(ctx context.Context, req *RegisterUserRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockUserServiceMockRecorder) RegisterUser(ctx, req any) *MockUserServiceRegisterUserCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockUserService)(nil).RegisterUser), ctx, req)
	return &MockUserServiceRegisterUserCall{Call: call}
}

// MockUserServiceRegisterUserCall wrap *gomock.Call
type MockUserServiceRegisterUserCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserServiceRegisterUserCall) Return(arg0 error) *MockUserServiceRegisterUserCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserServiceRegisterUserCall) Do(f func(context.Context, *RegisterUserRequest) error) *MockUserServiceRegisterUserCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserServiceRegisterUserCall) DoAndReturn(f func(context.Context, *RegisterUserRequest) error) *MockUserServiceRegisterUserCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ResetUserPassword mocks base method.
func (m *MockUserService) ResetUserPassword(ctx context.Context, req *ResetUserPasswordRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetUserPassword", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetUserPassword indicates an expected call of ResetUserPassword.
func (mr *MockUserServiceMockRecorder) ResetUserPassword(ctx, req any) *MockUserServiceResetUserPasswordCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetUserPassword", reflect.TypeOf((*MockUserService)(nil).ResetUserPassword), ctx, req)
	return &MockUserServiceResetUserPasswordCall{Call: call}
}

// MockUserServiceResetUserPasswordCall wrap *gomock.Call
type MockUserServiceResetUserPasswordCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserServiceResetUserPasswordCall) Return(arg0 error) *MockUserServiceResetUserPasswordCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserServiceResetUserPasswordCall) Do(f func(context.Context, *ResetUserPasswordRequest) error) *MockUserServiceResetUserPasswordCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserServiceResetUserPasswordCall) DoAndReturn(f func(context.Context, *ResetUserPasswordRequest) error) *MockUserServiceResetUserPasswordCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateUserInfo mocks base method.
func (m *MockUserService) UpdateUserInfo(ctx context.Context, userId string, req *UpdateUserInfoRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserInfo", ctx, userId, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserInfo indicates an expected call of UpdateUserInfo.
func (mr *MockUserServiceMockRecorder) UpdateUserInfo(ctx, userId, req any) *MockUserServiceUpdateUserInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserInfo", reflect.TypeOf((*MockUserService)(nil).UpdateUserInfo), ctx, userId, req)
	return &MockUserServiceUpdateUserInfoCall{Call: call}
}

// MockUserServiceUpdateUserInfoCall wrap *gomock.Call
type MockUserServiceUpdateUserInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserServiceUpdateUserInfoCall) Return(arg0 error) *MockUserServiceUpdateUserInfoCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserServiceUpdateUserInfoCall) Do(f func(context.Context, string, *UpdateUserInfoRequest) error) *MockUserServiceUpdateUserInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserServiceUpdateUserInfoCall) DoAndReturn(f func(context.Context, string, *UpdateUserInfoRequest) error) *MockUserServiceUpdateUserInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UpdateUserPassword mocks base method.
func (m *MockUserService) UpdateUserPassword(ctx context.Context, req *UpdateUserPasswordRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockUserServiceMockRecorder) UpdateUserPassword(ctx, req any) *MockUserServiceUpdateUserPasswordCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockUserService)(nil).UpdateUserPassword), ctx, req)
	return &MockUserServiceUpdateUserPasswordCall{Call: call}
}

// MockUserServiceUpdateUserPasswordCall wrap *gomock.Call
type MockUserServiceUpdateUserPasswordCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockUserServiceUpdateUserPasswordCall) Return(arg0 error) *MockUserServiceUpdateUserPasswordCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockUserServiceUpdateUserPasswordCall) Do(f func(context.Context, *UpdateUserPasswordRequest) error) *MockUserServiceUpdateUserPasswordCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockUserServiceUpdateUserPasswordCall) DoAndReturn(f func(context.Context, *UpdateUserPasswordRequest) error) *MockUserServiceUpdateUserPasswordCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
