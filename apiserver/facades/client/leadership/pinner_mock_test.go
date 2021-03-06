// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/core/leadership (interfaces: Pinner)

// Package leadership_test is a generated GoMock package.
package leadership_test

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPinner is a mock of Pinner interface
type MockPinner struct {
	ctrl     *gomock.Controller
	recorder *MockPinnerMockRecorder
}

// MockPinnerMockRecorder is the mock recorder for MockPinner
type MockPinnerMockRecorder struct {
	mock *MockPinner
}

// NewMockPinner creates a new mock instance
func NewMockPinner(ctrl *gomock.Controller) *MockPinner {
	mock := &MockPinner{ctrl: ctrl}
	mock.recorder = &MockPinnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPinner) EXPECT() *MockPinnerMockRecorder {
	return m.recorder
}

// PinLeadership mocks base method
func (m *MockPinner) PinLeadership(arg0 string) error {
	ret := m.ctrl.Call(m, "PinLeadership", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PinLeadership indicates an expected call of PinLeadership
func (mr *MockPinnerMockRecorder) PinLeadership(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PinLeadership", reflect.TypeOf((*MockPinner)(nil).PinLeadership), arg0)
}

// UnpinLeadership mocks base method
func (m *MockPinner) UnpinLeadership(arg0 string) error {
	ret := m.ctrl.Call(m, "UnpinLeadership", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnpinLeadership indicates an expected call of UnpinLeadership
func (mr *MockPinnerMockRecorder) UnpinLeadership(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnpinLeadership", reflect.TypeOf((*MockPinner)(nil).UnpinLeadership), arg0)
}
