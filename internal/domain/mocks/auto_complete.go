// Code generated by MockGen. DO NOT EDIT.
// Source: autocomplete_repo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	domain "gitlab.com/igor.tumanov1/theboatscom/internal/domain"
	reflect "reflect"
)

// MockIAutoCompleteRepository is a mock of IAutoCompleteRepository interface
type MockIAutoCompleteRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIAutoCompleteRepositoryMockRecorder
}

// MockIAutoCompleteRepositoryMockRecorder is the mock recorder for MockIAutoCompleteRepository
type MockIAutoCompleteRepositoryMockRecorder struct {
	mock *MockIAutoCompleteRepository
}

// NewMockIAutoCompleteRepository creates a new mock instance
func NewMockIAutoCompleteRepository(ctrl *gomock.Controller) *MockIAutoCompleteRepository {
	mock := &MockIAutoCompleteRepository{ctrl: ctrl}
	mock.recorder = &MockIAutoCompleteRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIAutoCompleteRepository) EXPECT() *MockIAutoCompleteRepositoryMockRecorder {
	return m.recorder
}

// Find mocks base method
func (m *MockIAutoCompleteRepository) Find(ctx context.Context, filter *domain.AutoCompleteFilter) ([]*domain.AutoCompleteEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, filter)
	ret0, _ := ret[0].([]*domain.AutoCompleteEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockIAutoCompleteRepositoryMockRecorder) Find(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockIAutoCompleteRepository)(nil).Find), ctx, filter)
}
