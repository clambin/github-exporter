// Code generated by mockery v2.29.0. DO NOT EDIT.

package mocks

import (
	context "context"
	github2 "github.com/clambin/github-exporter/pkg/github"

	mock "github.com/stretchr/testify/mock"
)

// GitHubClient is an autogenerated mock type for the GitHubClient type
type GitHubClient struct {
	mock.Mock
}

// GetPullRequests provides a mock function with given fields: _a0, _a1
func (_m *GitHubClient) GetPullRequests(_a0 context.Context, _a1 string) ([]github2.PullRequest, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []github2.PullRequest
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]github2.PullRequest, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []github2.PullRequest); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]github2.PullRequest)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRepo provides a mock function with given fields: _a0, _a1
func (_m *GitHubClient) GetRepo(_a0 context.Context, _a1 string) (github2.Repo, error) {
	ret := _m.Called(_a0, _a1)

	var r0 github2.Repo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (github2.Repo, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) github2.Repo); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(github2.Repo)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserRepos provides a mock function with given fields: _a0, _a1
func (_m *GitHubClient) GetUserRepos(_a0 context.Context, _a1 string) ([]github2.Repo, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []github2.Repo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]github2.Repo, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []github2.Repo); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]github2.Repo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewGitHubClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewGitHubClient creates a new instance of GitHubClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGitHubClient(t mockConstructorTestingTNewGitHubClient) *GitHubClient {
	mock := &GitHubClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}