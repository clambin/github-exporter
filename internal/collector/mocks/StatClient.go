// Code generated by mockery v2.50.1. DO NOT EDIT.

package mocks

import (
	context "context"

	github "github.com/clambin/github-exporter/internal/stats/github"
	mock "github.com/stretchr/testify/mock"
)

// StatClient is an autogenerated mock type for the StatClient type
type StatClient struct {
	mock.Mock
}

type StatClient_Expecter struct {
	mock *mock.Mock
}

func (_m *StatClient) EXPECT() *StatClient_Expecter {
	return &StatClient_Expecter{mock: &_m.Mock}
}

// GetRepoStats provides a mock function with given fields: _a0, _a1, _a2
func (_m *StatClient) GetRepoStats(_a0 context.Context, _a1 []string, _a2 []string) ([]github.RepoStats, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for GetRepoStats")
	}

	var r0 []github.RepoStats
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string, []string) ([]github.RepoStats, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string, []string) []github.RepoStats); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]github.RepoStats)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string, []string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StatClient_GetRepoStats_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRepoStats'
type StatClient_GetRepoStats_Call struct {
	*mock.Call
}

// GetRepoStats is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []string
//   - _a2 []string
func (_e *StatClient_Expecter) GetRepoStats(_a0 interface{}, _a1 interface{}, _a2 interface{}) *StatClient_GetRepoStats_Call {
	return &StatClient_GetRepoStats_Call{Call: _e.mock.On("GetRepoStats", _a0, _a1, _a2)}
}

func (_c *StatClient_GetRepoStats_Call) Run(run func(_a0 context.Context, _a1 []string, _a2 []string)) *StatClient_GetRepoStats_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string), args[2].([]string))
	})
	return _c
}

func (_c *StatClient_GetRepoStats_Call) Return(_a0 []github.RepoStats, _a1 error) *StatClient_GetRepoStats_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StatClient_GetRepoStats_Call) RunAndReturn(run func(context.Context, []string, []string) ([]github.RepoStats, error)) *StatClient_GetRepoStats_Call {
	_c.Call.Return(run)
	return _c
}

// NewStatClient creates a new instance of StatClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStatClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *StatClient {
	mock := &StatClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
