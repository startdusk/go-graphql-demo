// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	data "github.com/startdusk/twitter/data"
	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *UserService) GetByID(ctx context.Context, id string) (data.User, error) {
	ret := _m.Called(ctx, id)

	var r0 data.User
	if rf, ok := ret.Get(0).(func(context.Context, string) data.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(data.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserService interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserService(t mockConstructorTestingTNewUserService) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
