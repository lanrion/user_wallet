// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/gz00168ml/Documents/user_wallet/pkg/wallet/consumer.go

// Package mock_wallet is a generated GoMock package.
package wallet

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	kfkmodule "user_wallet/pkg/internal/kfkmodule"
)

// MockConsumer is a mock of Consumer interface.
type MockConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerMockRecorder
}

// MockConsumerMockRecorder is the mock recorder for MockConsumer.
type MockConsumerMockRecorder struct {
	mock *MockConsumer
}

// NewMockConsumer creates a new mock instance.
func NewMockConsumer(ctrl *gomock.Controller) *MockConsumer {
	mock := &MockConsumer{ctrl: ctrl}
	mock.recorder = &MockConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumer) EXPECT() *MockConsumerMockRecorder {
	return m.recorder
}

// Produce mocks base method.
func (m *MockConsumer) Produce(arg0 *kfkmodule.PushData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Produce", arg0)
}

// Produce indicates an expected call of Produce.
func (mr *MockConsumerMockRecorder) Produce(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Produce", reflect.TypeOf((*MockConsumer)(nil).Produce), arg0)
}

// PushBack mocks base method.
func (m *MockConsumer) PushBack(ctx context.Context, job *Job) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushBack", ctx, job)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushBack indicates an expected call of PushBack.
func (mr *MockConsumerMockRecorder) PushBack(ctx, job interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushBack", reflect.TypeOf((*MockConsumer)(nil).PushBack), ctx, job)
}

// Start mocks base method.
func (m *MockConsumer) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockConsumerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockConsumer)(nil).Start))
}
