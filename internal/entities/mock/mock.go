// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/weaveworks/magalix-policy-agent/pkg/domain (interfaces: EntitiesSource)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/weaveworks/magalix-policy-agent/pkg/domain"
)

// MockEntitiesSource is a mock of EntitiesSource interface.
type MockEntitiesSource struct {
	ctrl     *gomock.Controller
	recorder *MockEntitiesSourceMockRecorder
}

// MockEntitiesSourceMockRecorder is the mock recorder for MockEntitiesSource.
type MockEntitiesSourceMockRecorder struct {
	mock *MockEntitiesSource
}

// NewMockEntitiesSource creates a new mock instance.
func NewMockEntitiesSource(ctrl *gomock.Controller) *MockEntitiesSource {
	mock := &MockEntitiesSource{ctrl: ctrl}
	mock.recorder = &MockEntitiesSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntitiesSource) EXPECT() *MockEntitiesSourceMockRecorder {
	return m.recorder
}

// Kind mocks base method.
func (m *MockEntitiesSource) Kind() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Kind")
	ret0, _ := ret[0].(string)
	return ret0
}

// Kind indicates an expected call of Kind.
func (mr *MockEntitiesSourceMockRecorder) Kind() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Kind", reflect.TypeOf((*MockEntitiesSource)(nil).Kind))
}

// List mocks base method.
func (m *MockEntitiesSource) List(arg0 context.Context, arg1 *domain.ListOptions) (*domain.EntitiesList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*domain.EntitiesList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockEntitiesSourceMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockEntitiesSource)(nil).List), arg0, arg1)
}
