// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocker

import mock "github.com/stretchr/testify/mock"

// MockTriggerDefiner is an autogenerated mock type for the TriggerDefiner type
type MockTriggerDefiner struct {
	mock.Mock
}

type MockTriggerDefiner_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTriggerDefiner) EXPECT() *MockTriggerDefiner_Expecter {
	return &MockTriggerDefiner_Expecter{mock: &_m.Mock}
}

// GetDeploymentIds provides a mock function with given fields:
func (_m *MockTriggerDefiner) GetDeploymentIds() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDeploymentIds")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// MockTriggerDefiner_GetDeploymentIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDeploymentIds'
type MockTriggerDefiner_GetDeploymentIds_Call struct {
	*mock.Call
}

// GetDeploymentIds is a helper method to define mock.On call
func (_e *MockTriggerDefiner_Expecter) GetDeploymentIds() *MockTriggerDefiner_GetDeploymentIds_Call {
	return &MockTriggerDefiner_GetDeploymentIds_Call{Call: _e.mock.On("GetDeploymentIds")}
}

func (_c *MockTriggerDefiner_GetDeploymentIds_Call) Run(run func()) *MockTriggerDefiner_GetDeploymentIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTriggerDefiner_GetDeploymentIds_Call) Return(_a0 []string) *MockTriggerDefiner_GetDeploymentIds_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTriggerDefiner_GetDeploymentIds_Call) RunAndReturn(run func() []string) *MockTriggerDefiner_GetDeploymentIds_Call {
	_c.Call.Return(run)
	return _c
}

// GetName provides a mock function with given fields:
func (_m *MockTriggerDefiner) GetName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockTriggerDefiner_GetName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetName'
type MockTriggerDefiner_GetName_Call struct {
	*mock.Call
}

// GetName is a helper method to define mock.On call
func (_e *MockTriggerDefiner_Expecter) GetName() *MockTriggerDefiner_GetName_Call {
	return &MockTriggerDefiner_GetName_Call{Call: _e.mock.On("GetName")}
}

func (_c *MockTriggerDefiner_GetName_Call) Run(run func()) *MockTriggerDefiner_GetName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTriggerDefiner_GetName_Call) Return(_a0 string) *MockTriggerDefiner_GetName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTriggerDefiner_GetName_Call) RunAndReturn(run func() string) *MockTriggerDefiner_GetName_Call {
	_c.Call.Return(run)
	return _c
}

// GetType provides a mock function with given fields:
func (_m *MockTriggerDefiner) GetType() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetType")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockTriggerDefiner_GetType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetType'
type MockTriggerDefiner_GetType_Call struct {
	*mock.Call
}

// GetType is a helper method to define mock.On call
func (_e *MockTriggerDefiner_Expecter) GetType() *MockTriggerDefiner_GetType_Call {
	return &MockTriggerDefiner_GetType_Call{Call: _e.mock.On("GetType")}
}

func (_c *MockTriggerDefiner_GetType_Call) Run(run func()) *MockTriggerDefiner_GetType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockTriggerDefiner_GetType_Call) Return(_a0 string) *MockTriggerDefiner_GetType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockTriggerDefiner_GetType_Call) RunAndReturn(run func() string) *MockTriggerDefiner_GetType_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockTriggerDefiner creates a new instance of MockTriggerDefiner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTriggerDefiner(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTriggerDefiner {
	mock := &MockTriggerDefiner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
