// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/environs (interfaces: FinalizeCloudContext)

// Package testing is a generated GoMock package.
package testing

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockFinalizeCloudContext is a mock of FinalizeCloudContext interface
type MockFinalizeCloudContext struct {
	ctrl     *gomock.Controller
	recorder *MockFinalizeCloudContextMockRecorder
}

// MockFinalizeCloudContextMockRecorder is the mock recorder for MockFinalizeCloudContext
type MockFinalizeCloudContextMockRecorder struct {
	mock *MockFinalizeCloudContext
}

// NewMockFinalizeCloudContext creates a new mock instance
func NewMockFinalizeCloudContext(ctrl *gomock.Controller) *MockFinalizeCloudContext {
	mock := &MockFinalizeCloudContext{ctrl: ctrl}
	mock.recorder = &MockFinalizeCloudContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFinalizeCloudContext) EXPECT() *MockFinalizeCloudContextMockRecorder {
	return m.recorder
}

// Verbosef mocks base method
func (m *MockFinalizeCloudContext) Verbosef(arg0 string, arg1 ...interface{}) {
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Verbosef", varargs...)
}

// Verbosef indicates an expected call of Verbosef
func (mr *MockFinalizeCloudContextMockRecorder) Verbosef(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verbosef", reflect.TypeOf((*MockFinalizeCloudContext)(nil).Verbosef), varargs...)
}
