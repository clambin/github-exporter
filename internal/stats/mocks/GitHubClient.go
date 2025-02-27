// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	context "context"

	github "github.com/clambin/github-exporter/internal/stats/github"
	mock "github.com/stretchr/testify/mock"
)

// GitHubClient is an autogenerated mock type for the GitHubClient type
type GitHubClient struct {
	mock.Mock
}

type GitHubClient_Expecter struct {
	mock *mock.Mock
}

func (_m *GitHubClient) EXPECT() *GitHubClient_Expecter {
	return &GitHubClient_Expecter{mock: &_m.Mock}
}

// GetPullRequestCount provides a mock function with given fields: _a0, _a1, _a2
func (_m *GitHubClient) GetPullRequestCount(_a0 context.Context, _a1 string, _a2 string) (int, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for GetPullRequestCount")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (int, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) int); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GitHubClient_GetPullRequestCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPullRequestCount'
type GitHubClient_GetPullRequestCount_Call struct {
	*mock.Call
}

// GetPullRequestCount is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
//   - _a2 string
func (_e *GitHubClient_Expecter) GetPullRequestCount(_a0 interface{}, _a1 interface{}, _a2 interface{}) *GitHubClient_GetPullRequestCount_Call {
	return &GitHubClient_GetPullRequestCount_Call{Call: _e.mock.On("GetPullRequestCount", _a0, _a1, _a2)}
}

func (_c *GitHubClient_GetPullRequestCount_Call) Run(run func(_a0 context.Context, _a1 string, _a2 string)) *GitHubClient_GetPullRequestCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *GitHubClient_GetPullRequestCount_Call) Return(_a0 int, _a1 error) *GitHubClient_GetPullRequestCount_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GitHubClient_GetPullRequestCount_Call) RunAndReturn(run func(context.Context, string, string) (int, error)) *GitHubClient_GetPullRequestCount_Call {
	_c.Call.Return(run)
	return _c
}

// GetRepoStats provides a mock function with given fields: _a0, _a1, _a2
func (_m *GitHubClient) GetRepoStats(_a0 context.Context, _a1 string, _a2 string) (github.RepoStats, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for GetRepoStats")
	}

	var r0 github.RepoStats
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (github.RepoStats, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) github.RepoStats); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(github.RepoStats)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GitHubClient_GetRepoStats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRepoStats'
type GitHubClient_GetRepoStats_Call struct {
	*mock.Call
}

// GetRepoStats is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
//   - _a2 string
func (_e *GitHubClient_Expecter) GetRepoStats(_a0 interface{}, _a1 interface{}, _a2 interface{}) *GitHubClient_GetRepoStats_Call {
	return &GitHubClient_GetRepoStats_Call{Call: _e.mock.On("GetRepoStats", _a0, _a1, _a2)}
}

func (_c *GitHubClient_GetRepoStats_Call) Run(run func(_a0 context.Context, _a1 string, _a2 string)) *GitHubClient_GetRepoStats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *GitHubClient_GetRepoStats_Call) Return(_a0 github.RepoStats, _a1 error) *GitHubClient_GetRepoStats_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GitHubClient_GetRepoStats_Call) RunAndReturn(run func(context.Context, string, string) (github.RepoStats, error)) *GitHubClient_GetRepoStats_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserRepoNames provides a mock function with given fields: _a0, _a1
func (_m *GitHubClient) GetUserRepoNames(_a0 context.Context, _a1 string) ([]string, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetUserRepoNames")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]string, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GitHubClient_GetUserRepoNames_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserRepoNames'
type GitHubClient_GetUserRepoNames_Call struct {
	*mock.Call
}

// GetUserRepoNames is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *GitHubClient_Expecter) GetUserRepoNames(_a0 interface{}, _a1 interface{}) *GitHubClient_GetUserRepoNames_Call {
	return &GitHubClient_GetUserRepoNames_Call{Call: _e.mock.On("GetUserRepoNames", _a0, _a1)}
}

func (_c *GitHubClient_GetUserRepoNames_Call) Run(run func(_a0 context.Context, _a1 string)) *GitHubClient_GetUserRepoNames_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *GitHubClient_GetUserRepoNames_Call) Return(_a0 []string, _a1 error) *GitHubClient_GetUserRepoNames_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GitHubClient_GetUserRepoNames_Call) RunAndReturn(run func(context.Context, string) ([]string, error)) *GitHubClient_GetUserRepoNames_Call {
	_c.Call.Return(run)
	return _c
}

// NewGitHubClient creates a new instance of GitHubClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGitHubClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *GitHubClient {
	mock := &GitHubClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
