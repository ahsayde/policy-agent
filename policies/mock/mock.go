// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/MagalixCorp/magalix-policy-agent/pkg/domain (interfaces: PoliciesSource)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	domain "github.com/MagalixCorp/magalix-policy-agent/pkg/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockPoliciesSource is a mock of PoliciesSource interface.
type MockPoliciesSource struct {
	ctrl     *gomock.Controller
	recorder *MockPoliciesSourceMockRecorder
}

// MockPoliciesSourceMockRecorder is the mock recorder for MockPoliciesSource.
type MockPoliciesSourceMockRecorder struct {
	mock *MockPoliciesSource
}

// NewMockPoliciesSource creates a new mock instance.
func NewMockPoliciesSource(ctrl *gomock.Controller) *MockPoliciesSource {
	mock := &MockPoliciesSource{ctrl: ctrl}
	mock.recorder = &MockPoliciesSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPoliciesSource) EXPECT() *MockPoliciesSourceMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockPoliciesSource) GetAll(arg0 context.Context) ([]domain.Policy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]domain.Policy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPoliciesSourceMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPoliciesSource)(nil).GetAll), arg0)
}
