// Code generated by mockery v2.42.1. DO NOT EDIT.

package main

import (
	domain "github.com/ichenhe/cert-deployer/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockcommandExecutor is an autogenerated mock type for the commandExecutor type
type MockcommandExecutor struct {
	mock.Mock
}

type MockcommandExecutor_Expecter struct {
	mock *mock.Mock
}

func (_m *MockcommandExecutor) EXPECT() *MockcommandExecutor_Expecter {
	return &MockcommandExecutor_Expecter{mock: &_m.Mock}
}

// customDeploy provides a mock function with given fields: providers, rawTypes, cert, key
func (_m *MockcommandExecutor) customDeploy(providers map[string]domain.CloudProvider, rawTypes []string, cert []byte, key []byte) {
	_m.Called(providers, rawTypes, cert, key)
}

// MockcommandExecutor_customDeploy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'customDeploy'
type MockcommandExecutor_customDeploy_Call struct {
	*mock.Call
}

// customDeploy is a helper method to define mock.On call
//   - providers map[string]domain.CloudProvider
//   - rawTypes []string
//   - cert []byte
//   - key []byte
func (_e *MockcommandExecutor_Expecter) customDeploy(providers interface{}, rawTypes interface{}, cert interface{}, key interface{}) *MockcommandExecutor_customDeploy_Call {
	return &MockcommandExecutor_customDeploy_Call{Call: _e.mock.On("customDeploy", providers, rawTypes, cert, key)}
}

func (_c *MockcommandExecutor_customDeploy_Call) Run(run func(providers map[string]domain.CloudProvider, rawTypes []string, cert []byte, key []byte)) *MockcommandExecutor_customDeploy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]domain.CloudProvider), args[1].([]string), args[2].([]byte), args[3].([]byte))
	})
	return _c
}

func (_c *MockcommandExecutor_customDeploy_Call) Return() *MockcommandExecutor_customDeploy_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockcommandExecutor_customDeploy_Call) RunAndReturn(run func(map[string]domain.CloudProvider, []string, []byte, []byte)) *MockcommandExecutor_customDeploy_Call {
	_c.Call.Return(run)
	return _c
}

// executeDeployments provides a mock function with given fields: appConfig, deploymentIds
func (_m *MockcommandExecutor) executeDeployments(appConfig *domain.AppConfig, deploymentIds []string) {
	_m.Called(appConfig, deploymentIds)
}

// MockcommandExecutor_executeDeployments_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'executeDeployments'
type MockcommandExecutor_executeDeployments_Call struct {
	*mock.Call
}

// executeDeployments is a helper method to define mock.On call
//   - appConfig *domain.AppConfig
//   - deploymentIds []string
func (_e *MockcommandExecutor_Expecter) executeDeployments(appConfig interface{}, deploymentIds interface{}) *MockcommandExecutor_executeDeployments_Call {
	return &MockcommandExecutor_executeDeployments_Call{Call: _e.mock.On("executeDeployments", appConfig, deploymentIds)}
}

func (_c *MockcommandExecutor_executeDeployments_Call) Run(run func(appConfig *domain.AppConfig, deploymentIds []string)) *MockcommandExecutor_executeDeployments_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.AppConfig), args[1].([]string))
	})
	return _c
}

func (_c *MockcommandExecutor_executeDeployments_Call) Return() *MockcommandExecutor_executeDeployments_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockcommandExecutor_executeDeployments_Call) RunAndReturn(run func(*domain.AppConfig, []string)) *MockcommandExecutor_executeDeployments_Call {
	_c.Call.Return(run)
	return _c
}

// registerTriggers provides a mock function with given fields: providers, deployments, triggerDefs
func (_m *MockcommandExecutor) registerTriggers(providers map[string]domain.CloudProvider, deployments map[string]domain.Deployment, triggerDefs map[string]domain.TriggerDefiner) []domain.Trigger {
	ret := _m.Called(providers, deployments, triggerDefs)

	if len(ret) == 0 {
		panic("no return value specified for registerTriggers")
	}

	var r0 []domain.Trigger
	if rf, ok := ret.Get(0).(func(map[string]domain.CloudProvider, map[string]domain.Deployment, map[string]domain.TriggerDefiner) []domain.Trigger); ok {
		r0 = rf(providers, deployments, triggerDefs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Trigger)
		}
	}

	return r0
}

// MockcommandExecutor_registerTriggers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'registerTriggers'
type MockcommandExecutor_registerTriggers_Call struct {
	*mock.Call
}

// registerTriggers is a helper method to define mock.On call
//   - providers map[string]domain.CloudProvider
//   - deployments map[string]domain.Deployment
//   - triggerDefs map[string]domain.TriggerDefiner
func (_e *MockcommandExecutor_Expecter) registerTriggers(providers interface{}, deployments interface{}, triggerDefs interface{}) *MockcommandExecutor_registerTriggers_Call {
	return &MockcommandExecutor_registerTriggers_Call{Call: _e.mock.On("registerTriggers", providers, deployments, triggerDefs)}
}

func (_c *MockcommandExecutor_registerTriggers_Call) Run(run func(providers map[string]domain.CloudProvider, deployments map[string]domain.Deployment, triggerDefs map[string]domain.TriggerDefiner)) *MockcommandExecutor_registerTriggers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]domain.CloudProvider), args[1].(map[string]domain.Deployment), args[2].(map[string]domain.TriggerDefiner))
	})
	return _c
}

func (_c *MockcommandExecutor_registerTriggers_Call) Return(registeredTriggers []domain.Trigger) *MockcommandExecutor_registerTriggers_Call {
	_c.Call.Return(registeredTriggers)
	return _c
}

func (_c *MockcommandExecutor_registerTriggers_Call) RunAndReturn(run func(map[string]domain.CloudProvider, map[string]domain.Deployment, map[string]domain.TriggerDefiner) []domain.Trigger) *MockcommandExecutor_registerTriggers_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockcommandExecutor creates a new instance of MockcommandExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockcommandExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockcommandExecutor {
	mock := &MockcommandExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
