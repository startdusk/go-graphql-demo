// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	data "github.com/startdusk/twitter/data"
	mock "github.com/stretchr/testify/mock"
)

// RefreshTokenRepo is an autogenerated mock type for the RefreshTokenRepo type
type RefreshTokenRepo struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, params
func (_m *RefreshTokenRepo) Create(ctx context.Context, params data.CreateRefreshTokenParams) (data.RefreshToken, error) {
	ret := _m.Called(ctx, params)

	var r0 data.RefreshToken
	if rf, ok := ret.Get(0).(func(context.Context, data.CreateRefreshTokenParams) data.RefreshToken); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Get(0).(data.RefreshToken)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, data.CreateRefreshTokenParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *RefreshTokenRepo) GetByID(ctx context.Context, id string) (data.RefreshToken, error) {
	ret := _m.Called(ctx, id)

	var r0 data.RefreshToken
	if rf, ok := ret.Get(0).(func(context.Context, string) data.RefreshToken); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(data.RefreshToken)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRefreshTokenRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewRefreshTokenRepo creates a new instance of RefreshTokenRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRefreshTokenRepo(t mockConstructorTestingTNewRefreshTokenRepo) *RefreshTokenRepo {
	mock := &RefreshTokenRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
