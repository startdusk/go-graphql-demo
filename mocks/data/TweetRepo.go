// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	data "github.com/startdusk/twitter/data"
	mock "github.com/stretchr/testify/mock"
)

// TweetRepo is an autogenerated mock type for the TweetRepo type
type TweetRepo struct {
	mock.Mock
}

// All provides a mock function with given fields: ctx
func (_m *TweetRepo) All(ctx context.Context) ([]data.Tweet, error) {
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

// Create provides a mock function with given fields: ctx, tweet
func (_m *TweetRepo) Create(ctx context.Context, tweet data.Tweet) (data.Tweet, error) {
	ret := _m.Called(ctx, tweet)

	var r0 data.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, data.Tweet) data.Tweet); ok {
		r0 = rf(ctx, tweet)
	} else {
		r0 = ret.Get(0).(data.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, data.Tweet) error); ok {
		r1 = rf(ctx, tweet)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, tweetID
func (_m *TweetRepo) GetByID(ctx context.Context, tweetID string) (data.Tweet, error) {
	ret := _m.Called(ctx, tweetID)

	var r0 data.Tweet
	if rf, ok := ret.Get(0).(func(context.Context, string) data.Tweet); ok {
		r0 = rf(ctx, tweetID)
	} else {
		r0 = ret.Get(0).(data.Tweet)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tweetID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTweetRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewTweetRepo creates a new instance of TweetRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTweetRepo(t mockConstructorTestingTNewTweetRepo) *TweetRepo {
	mock := &TweetRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
