// Code generated by MockGen. DO NOT EDIT.
// Source: config.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockflagSet is a mock of flagSet interface
type MockflagSet struct {
	ctrl     *gomock.Controller
	recorder *MockflagSetMockRecorder
}

// MockflagSetMockRecorder is the mock recorder for MockflagSet
type MockflagSetMockRecorder struct {
	mock *MockflagSet
}

// NewMockflagSet creates a new mock instance
func NewMockflagSet(ctrl *gomock.Controller) *MockflagSet {
	mock := &MockflagSet{ctrl: ctrl}
	mock.recorder = &MockflagSetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockflagSet) EXPECT() *MockflagSetMockRecorder {
	return m.recorder
}

// SetOutput mocks base method
func (m *MockflagSet) SetOutput(arg0 io.Writer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOutput", arg0)
}

// SetOutput indicates an expected call of SetOutput
func (mr *MockflagSetMockRecorder) SetOutput(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOutput", reflect.TypeOf((*MockflagSet)(nil).SetOutput), arg0)
}

// String mocks base method
func (m *MockflagSet) String(arg0, arg1, arg2 string) *string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String", arg0, arg1, arg2)
	ret0, _ := ret[0].(*string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockflagSetMockRecorder) String(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockflagSet)(nil).String), arg0, arg1, arg2)
}

// Bool mocks base method
func (m *MockflagSet) Bool(arg0 string, arg1 bool, arg2 string) *bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bool", arg0, arg1, arg2)
	ret0, _ := ret[0].(*bool)
	return ret0
}

// Bool indicates an expected call of Bool
func (mr *MockflagSetMockRecorder) Bool(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bool", reflect.TypeOf((*MockflagSet)(nil).Bool), arg0, arg1, arg2)
}

// Uint mocks base method
func (m *MockflagSet) Uint(arg0 string, arg1 uint, arg2 string) *uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Uint", arg0, arg1, arg2)
	ret0, _ := ret[0].(*uint)
	return ret0
}

// Uint indicates an expected call of Uint
func (mr *MockflagSetMockRecorder) Uint(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Uint", reflect.TypeOf((*MockflagSet)(nil).Uint), arg0, arg1, arg2)
}

// Parse mocks base method
func (m *MockflagSet) Parse(arg0 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Parse indicates an expected call of Parse
func (mr *MockflagSetMockRecorder) Parse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockflagSet)(nil).Parse), arg0)
}

// PrintDefaults mocks base method
func (m *MockflagSet) PrintDefaults() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PrintDefaults")
}

// PrintDefaults indicates an expected call of PrintDefaults
func (mr *MockflagSetMockRecorder) PrintDefaults() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrintDefaults", reflect.TypeOf((*MockflagSet)(nil).PrintDefaults))
}
