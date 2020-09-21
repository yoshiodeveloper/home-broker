// Code generated by MockGen. DO NOT EDIT.
// Source: ./users/db.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	users "home-broker/users"
	reflect "reflect"
)

// MockUserDBInterface is a mock of UserDBInterface interface
type MockUserDBInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserDBInterfaceMockRecorder
}

// MockUserDBInterfaceMockRecorder is the mock recorder for MockUserDBInterface
type MockUserDBInterfaceMockRecorder struct {
	mock *MockUserDBInterface
}

// NewMockUserDBInterface creates a new mock instance
func NewMockUserDBInterface(ctrl *gomock.Controller) *MockUserDBInterface {
	mock := &MockUserDBInterface{ctrl: ctrl}
	mock.recorder = &MockUserDBInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserDBInterface) EXPECT() *MockUserDBInterfaceMockRecorder {
	return m.recorder
}

// GetByID mocks base method
func (m *MockUserDBInterface) GetByID(id users.UserID) (*users.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(*users.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockUserDBInterfaceMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUserDBInterface)(nil).GetByID), id)
}

// Insert mocks base method
func (m *MockUserDBInterface) Insert(entity users.User) (*users.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", entity)
	ret0, _ := ret[0].(*users.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert
func (mr *MockUserDBInterfaceMockRecorder) Insert(entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserDBInterface)(nil).Insert), entity)
}