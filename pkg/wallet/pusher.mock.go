// Code generated by MockGen. DO NOT EDIT.
// Source: /Users/gz00168ml/Documents/user_wallet/pkg/wallet/pusher.go

// Package mock_wallet is a generated GoMock package.
package wallet

import (
	reflect "reflect"
	kfkmodule "user_wallet/pkg/internal/kfkmodule"

	gomock "github.com/golang/mock/gomock"
)

// MockPusher is a mock of Pusher interface.
type MockPusher struct {
	ctrl     *gomock.Controller
	recorder *MockPusherMockRecorder
}

// MockPusherMockRecorder is the mock recorder for MockPusher.
type MockPusherMockRecorder struct {
	mock *MockPusher
}

// NewMockPusher creates a new mock instance.
func NewMockPusher(ctrl *gomock.Controller) *MockPusher {
	mock := &MockPusher{ctrl: ctrl}
	mock.recorder = &MockPusherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPusher) EXPECT() *MockPusherMockRecorder {
	return m.recorder
}

// PushWallet mocks base method.
func (m *MockPusher) PushWallet(d *kfkmodule.PushData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PushWallet", d)
}

// PushWallet indicates an expected call of PushWallet.
func (mr *MockPusherMockRecorder) PushWallet(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushWallet", reflect.TypeOf((*MockPusher)(nil).PushWallet), d)
}
