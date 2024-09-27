// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KScaesar/go-layout/pkg/utility/dataflow (interfaces: Producer)
//
// Generated by this command:
//
//	mockgen -typed -package=dataflow -destination=pubsub_mock.go github.com/KScaesar/go-layout/pkg/utility/dataflow Producer
//

// Package dataflow is a generated GoMock package.
package dataflow

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockProducer) Send(arg0 ...*Message) error {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockProducerMockRecorder) Send(arg0 ...any) *MockProducerSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockProducer)(nil).Send), arg0...)
	return &MockProducerSendCall{Call: call}
}

// MockProducerSendCall wrap *gomock.Call
type MockProducerSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProducerSendCall) Return(arg0 error) *MockProducerSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProducerSendCall) Do(f func(...*Message) error) *MockProducerSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProducerSendCall) DoAndReturn(f func(...*Message) error) *MockProducerSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendWithCtx mocks base method.
func (m *MockProducer) SendWithCtx(arg0 context.Context, arg1 ...*Message) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendWithCtx", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendWithCtx indicates an expected call of SendWithCtx.
func (mr *MockProducerMockRecorder) SendWithCtx(arg0 any, arg1 ...any) *MockProducerSendWithCtxCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendWithCtx", reflect.TypeOf((*MockProducer)(nil).SendWithCtx), varargs...)
	return &MockProducerSendWithCtxCall{Call: call}
}

// MockProducerSendWithCtxCall wrap *gomock.Call
type MockProducerSendWithCtxCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockProducerSendWithCtxCall) Return(arg0 error) *MockProducerSendWithCtxCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockProducerSendWithCtxCall) Do(f func(context.Context, ...*Message) error) *MockProducerSendWithCtxCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockProducerSendWithCtxCall) DoAndReturn(f func(context.Context, ...*Message) error) *MockProducerSendWithCtxCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}