// Code generated by MockGen. DO NOT EDIT.
// Source: calendar_repo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	domain "gitlab.com/igor.tumanov1/theboatscom/internal/domain"
	reflect "reflect"
)

// MockICalendarRepository is a mock of ICalendarRepository interface
type MockICalendarRepository struct {
	ctrl     *gomock.Controller
	recorder *MockICalendarRepositoryMockRecorder
}

// MockICalendarRepositoryMockRecorder is the mock recorder for MockICalendarRepository
type MockICalendarRepositoryMockRecorder struct {
	mock *MockICalendarRepository
}

// NewMockICalendarRepository creates a new mock instance
func NewMockICalendarRepository(ctrl *gomock.Controller) *MockICalendarRepository {
	mock := &MockICalendarRepository{ctrl: ctrl}
	mock.recorder = &MockICalendarRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockICalendarRepository) EXPECT() *MockICalendarRepositoryMockRecorder {
	return m.recorder
}

// GetAvailabilityByIDs mocks base method
func (m *MockICalendarRepository) GetAvailabilityByIDs(ctx context.Context, ids []int64) (map[int64]bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvailabilityByIDs", ctx, ids)
	ret0, _ := ret[0].(map[int64]bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvailabilityByIDs indicates an expected call of GetAvailabilityByIDs
func (mr *MockICalendarRepositoryMockRecorder) GetAvailabilityByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvailabilityByIDs", reflect.TypeOf((*MockICalendarRepository)(nil).GetAvailabilityByIDs), ctx, ids)
}

// GetUpcomingAvailabilityDatesByIDs mocks base method
func (m *MockICalendarRepository) GetUpcomingAvailabilityDatesByIDs(ctx context.Context, ids []int64) (map[int64]*domain.CalendarEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpcomingAvailabilityDatesByIDs", ctx, ids)
	ret0, _ := ret[0].(map[int64]*domain.CalendarEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUpcomingAvailabilityDatesByIDs indicates an expected call of GetUpcomingAvailabilityDatesByIDs
func (mr *MockICalendarRepositoryMockRecorder) GetUpcomingAvailabilityDatesByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpcomingAvailabilityDatesByIDs", reflect.TypeOf((*MockICalendarRepository)(nil).GetUpcomingAvailabilityDatesByIDs), ctx, ids)
}
