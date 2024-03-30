// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocker

import (
	context "context"

	domain "github.com/ichenhe/cert-deployer/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockDeploymentExecutor is an autogenerated mock type for the DeploymentExecutor type
type MockDeploymentExecutor struct {
	mock.Mock
}

type MockDeploymentExecutor_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDeploymentExecutor) EXPECT() *MockDeploymentExecutor_Expecter {
	return &MockDeploymentExecutor_Expecter{mock: &_m.Mock}
}

// ExecuteDeployment provides a mock function with given fields: ctx, deployment
func (_m *MockDeploymentExecutor) ExecuteDeployment(ctx context.Context, deployment domain.Deployment) error {
	ret := _m.Called(ctx, deployment)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteDeployment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Deployment) error); ok {
		r0 = rf(ctx, deployment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDeploymentExecutor_ExecuteDeployment_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteDeployment'
type MockDeploymentExecutor_ExecuteDeployment_Call struct {
	*mock.Call
}

// ExecuteDeployment is a helper method to define mock.On call
//   - ctx context.Context
//   - deployment domain.Deployment
func (_e *MockDeploymentExecutor_Expecter) ExecuteDeployment(ctx interface{}, deployment interface{}) *MockDeploymentExecutor_ExecuteDeployment_Call {
	return &MockDeploymentExecutor_ExecuteDeployment_Call{Call: _e.mock.On("ExecuteDeployment", ctx, deployment)}
}

func (_c *MockDeploymentExecutor_ExecuteDeployment_Call) Run(run func(ctx context.Context, deployment domain.Deployment)) *MockDeploymentExecutor_ExecuteDeployment_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Deployment))
	})
	return _c
}

func (_c *MockDeploymentExecutor_ExecuteDeployment_Call) Return(_a0 error) *MockDeploymentExecutor_ExecuteDeployment_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDeploymentExecutor_ExecuteDeployment_Call) RunAndReturn(run func(context.Context, domain.Deployment) error) *MockDeploymentExecutor_ExecuteDeployment_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDeploymentExecutor creates a new instance of MockDeploymentExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDeploymentExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDeploymentExecutor {
	mock := &MockDeploymentExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
