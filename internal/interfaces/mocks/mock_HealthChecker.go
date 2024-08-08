// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// HealthChecker is an autogenerated mock type for the HealthChecker type
type HealthChecker struct {
	mock.Mock
}

type HealthChecker_Expecter struct {
	mock *mock.Mock
}

func (_m *HealthChecker) EXPECT() *HealthChecker_Expecter {
	return &HealthChecker_Expecter{mock: &_m.Mock}
}

// GetCheckHandler provides a mock function with given fields:
func (_m *HealthChecker) GetCheckHandler() (http.Handler, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetCheckHandler")
	}

	var r0 http.Handler
	var r1 error
	if rf, ok := ret.Get(0).(func() (http.Handler, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() http.Handler); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Handler)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HealthChecker_GetCheckHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCheckHandler'
type HealthChecker_GetCheckHandler_Call struct {
	*mock.Call
}

// GetCheckHandler is a helper method to define mock.On call
func (_e *HealthChecker_Expecter) GetCheckHandler() *HealthChecker_GetCheckHandler_Call {
	return &HealthChecker_GetCheckHandler_Call{Call: _e.mock.On("GetCheckHandler")}
}

func (_c *HealthChecker_GetCheckHandler_Call) Run(run func()) *HealthChecker_GetCheckHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *HealthChecker_GetCheckHandler_Call) Return(_a0 http.Handler, _a1 error) *HealthChecker_GetCheckHandler_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *HealthChecker_GetCheckHandler_Call) RunAndReturn(run func() (http.Handler, error)) *HealthChecker_GetCheckHandler_Call {
	_c.Call.Return(run)
	return _c
}

// NewHealthChecker creates a new instance of HealthChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHealthChecker(t interface {
	mock.TestingT
	Cleanup(func())
}) *HealthChecker {
	mock := &HealthChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
