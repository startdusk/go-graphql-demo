// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	graph "github.com/startdusk/twitter/cmd/twitter/graph"
	mock "github.com/stretchr/testify/mock"
)

// TweetResolver is an autogenerated mock type for the TweetResolver type
type TweetResolver struct {
	mock.Mock
}

// User provides a mock function with given fields: ctx, obj
func (_m *TweetResolver) User(ctx context.Context, obj *graph.Tweet) (*graph.User, error) {
	ret := _m.Called(ctx, obj)

	var r0 *graph.User
	if rf, ok := ret.Get(0).(func(context.Context, *graph.Tweet) *graph.User); ok {
		r0 = rf(ctx, obj)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graph.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *graph.Tweet) error); ok {
		r1 = rf(ctx, obj)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTweetResolver interface {
	mock.TestingT
	Cleanup(func())
}

// NewTweetResolver creates a new instance of TweetResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTweetResolver(t mockConstructorTestingTNewTweetResolver) *TweetResolver {
	mock := &TweetResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
