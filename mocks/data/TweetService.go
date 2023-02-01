// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	data "github.com/startdusk/twitter/data"
	mock "github.com/stretchr/testify/mock"
)

// TweetService is an autogenerated mock type for the TweetService type
type TweetService struct {
	mock.Mock
}

// All provides a mock function with given fields: ctx
func (_m *TweetService) All(ctx context.Context) ([]data.Tweet, error) {
	ret := _m.Called(ctx)

	var r0 []data.Tweet
	if rf, ok := ret.Get(0).(func(context.Context) []data.Tweet); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]data.Tweet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, input
func (_m *TweetService) Create(ctx context.Context, input data.CreateTweetInput) (data.Tweet, error) {
	ret := _m.Called(ctx, input)

	var r0 data.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, data.CreateTweetInput) data.Tweet); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(data.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, data.CreateTweetInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *TweetService) GetByID(ctx context.Context, id string) (data.Tweet, error) {
	ret := _m.Called(ctx, id)

	var r0 data.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, string) data.Tweet); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(data.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTweetService interface {
	mock.TestingT
	Cleanup(func())
}

// NewTweetService creates a new instance of TweetService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTweetService(t mockConstructorTestingTNewTweetService) *TweetService {
	mock := &TweetService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}