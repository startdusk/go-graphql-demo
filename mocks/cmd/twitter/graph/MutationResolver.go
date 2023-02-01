// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	graph "github.com/startdusk/twitter/cmd/twitter/graph"
	mock "github.com/stretchr/testify/mock"
)

// MutationResolver is an autogenerated mock type for the MutationResolver type
type MutationResolver struct {
	mock.Mock
}

// CreateReply provides a mock function with given fields: ctx, parentID, input
func (_m *MutationResolver) CreateReply(ctx context.Context, parentID string, input graph.CreatedTweetInput) (*graph.Tweet, error) {
	ret := _m.Called(ctx, parentID, input)

	var r0 *graph.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, string, graph.CreatedTweetInput) *graph.Tweet); ok {
		r0 = rf(ctx, parentID, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graph.Tweet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, graph.CreatedTweetInput) error); ok {
		r1 = rf(ctx, parentID, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateTweet provides a mock function with given fields: ctx, input
func (_m *MutationResolver) CreateTweet(ctx context.Context, input graph.CreatedTweetInput) (*graph.Tweet, error) {
	ret := _m.Called(ctx, input)

	var r0 *graph.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, graph.CreatedTweetInput) *graph.Tweet); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graph.Tweet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, graph.CreatedTweetInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteTweet provides a mock function with given fields: ctx, id
func (_m *MutationResolver) DeleteTweet(ctx context.Context, id string) (bool, error) {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, input
func (_m *MutationResolver) Login(ctx context.Context, input graph.LoginInput) (*graph.AuthResponse, error) {
	ret := _m.Called(ctx, input)

	var r0 *graph.AuthResponse
	if rf, ok := ret.Get(0).(func(context.Context, graph.LoginInput) *graph.AuthResponse); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graph.AuthResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, graph.LoginInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, input
func (_m *MutationResolver) Register(ctx context.Context, input graph.RegisterInput) (*graph.AuthResponse, error) {
	ret := _m.Called(ctx, input)

	var r0 *graph.AuthResponse
	if rf, ok := ret.Get(0).(func(context.Context, graph.RegisterInput) *graph.AuthResponse); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graph.AuthResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, graph.RegisterInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMutationResolver interface {
	mock.TestingT
	Cleanup(func())
}

// NewMutationResolver creates a new instance of MutationResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMutationResolver(t mockConstructorTestingTNewMutationResolver) *MutationResolver {
	mock := &MutationResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
