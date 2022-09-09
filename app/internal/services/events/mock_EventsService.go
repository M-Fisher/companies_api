// Code generated by mockery v2.10.6. DO NOT EDIT.

package events

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockEventsService is an autogenerated mock type for the EventsService type
type MockEventsService struct {
	mock.Mock
}

// GetStatus provides a mock function with given fields:
func (_m *MockEventsService) GetStatus() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendEvent provides a mock function with given fields: ctx, event, data
func (_m *MockEventsService) SendEvent(ctx context.Context, event EventName, data []byte) error {
	ret := _m.Called(ctx, event, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, EventName, []byte) error); ok {
		r0 = rf(ctx, event, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *MockEventsService) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
