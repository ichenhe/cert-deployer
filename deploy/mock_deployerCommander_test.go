// Code generated by mockery v2.42.1. DO NOT EDIT.

package deploy

import (
	context "context"

	domain "github.com/ichenhe/cert-deployer/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockdeployerCommander is an autogenerated mock type for the deployerCommander type
type MockdeployerCommander struct {
	mock.Mock
}

type MockdeployerCommander_Expecter struct {
	mock *mock.Mock
}

func (_m *MockdeployerCommander) EXPECT() *MockdeployerCommander_Expecter {
	return &MockdeployerCommander_Expecter{mock: &_m.Mock}
}

// DeployToAsset provides a mock function with given fields: ctx, assetType, assetId, cert, key
func (_m *MockdeployerCommander) DeployToAsset(ctx context.Context, assetType string, assetId string, cert []byte, key []byte) error {
	ret := _m.Called(ctx, assetType, assetId, cert, key)

	if len(ret) == 0 {
		panic("no return value specified for DeployToAsset")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []byte, []byte) error); ok {
		r0 = rf(ctx, assetType, assetId, cert, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockdeployerCommander_DeployToAsset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeployToAsset'
type MockdeployerCommander_DeployToAsset_Call struct {
	*mock.Call
}

// DeployToAsset is a helper method to define mock.On call
//   - ctx context.Context
//   - assetType string
//   - assetId string
//   - cert []byte
//   - key []byte
func (_e *MockdeployerCommander_Expecter) DeployToAsset(ctx interface{}, assetType interface{}, assetId interface{}, cert interface{}, key interface{}) *MockdeployerCommander_DeployToAsset_Call {
	return &MockdeployerCommander_DeployToAsset_Call{Call: _e.mock.On("DeployToAsset", ctx, assetType, assetId, cert, key)}
}

func (_c *MockdeployerCommander_DeployToAsset_Call) Run(run func(ctx context.Context, assetType string, assetId string, cert []byte, key []byte)) *MockdeployerCommander_DeployToAsset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].([]byte), args[4].([]byte))
	})
	return _c
}

func (_c *MockdeployerCommander_DeployToAsset_Call) Return(_a0 error) *MockdeployerCommander_DeployToAsset_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockdeployerCommander_DeployToAsset_Call) RunAndReturn(run func(context.Context, string, string, []byte, []byte) error) *MockdeployerCommander_DeployToAsset_Call {
	_c.Call.Return(run)
	return _c
}

// DeployToAssetType provides a mock function with given fields: ctx, assetType, cert, key, onAssetsAcquired, onDeployResult
func (_m *MockdeployerCommander) DeployToAssetType(ctx context.Context, assetType string, cert []byte, key []byte, onAssetsAcquired func([]domain.Asseter), onDeployResult func(domain.Asseter, error)) error {
	ret := _m.Called(ctx, assetType, cert, key, onAssetsAcquired, onDeployResult)

	if len(ret) == 0 {
		panic("no return value specified for DeployToAssetType")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte, []byte, func([]domain.Asseter), func(domain.Asseter, error)) error); ok {
		r0 = rf(ctx, assetType, cert, key, onAssetsAcquired, onDeployResult)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockdeployerCommander_DeployToAssetType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeployToAssetType'
type MockdeployerCommander_DeployToAssetType_Call struct {
	*mock.Call
}

// DeployToAssetType is a helper method to define mock.On call
//   - ctx context.Context
//   - assetType string
//   - cert []byte
//   - key []byte
//   - onAssetsAcquired func([]domain.Asseter)
//   - onDeployResult func(domain.Asseter , error)
func (_e *MockdeployerCommander_Expecter) DeployToAssetType(ctx interface{}, assetType interface{}, cert interface{}, key interface{}, onAssetsAcquired interface{}, onDeployResult interface{}) *MockdeployerCommander_DeployToAssetType_Call {
	return &MockdeployerCommander_DeployToAssetType_Call{Call: _e.mock.On("DeployToAssetType", ctx, assetType, cert, key, onAssetsAcquired, onDeployResult)}
}

func (_c *MockdeployerCommander_DeployToAssetType_Call) Run(run func(ctx context.Context, assetType string, cert []byte, key []byte, onAssetsAcquired func([]domain.Asseter), onDeployResult func(domain.Asseter, error))) *MockdeployerCommander_DeployToAssetType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]byte), args[3].([]byte), args[4].(func([]domain.Asseter)), args[5].(func(domain.Asseter, error)))
	})
	return _c
}

func (_c *MockdeployerCommander_DeployToAssetType_Call) Return(_a0 error) *MockdeployerCommander_DeployToAssetType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockdeployerCommander_DeployToAssetType_Call) RunAndReturn(run func(context.Context, string, []byte, []byte, func([]domain.Asseter), func(domain.Asseter, error)) error) *MockdeployerCommander_DeployToAssetType_Call {
	_c.Call.Return(run)
	return _c
}

// IsAssetTypeSupported provides a mock function with given fields: assetType
func (_m *MockdeployerCommander) IsAssetTypeSupported(assetType string) bool {
	ret := _m.Called(assetType)

	if len(ret) == 0 {
		panic("no return value specified for IsAssetTypeSupported")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(assetType)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockdeployerCommander_IsAssetTypeSupported_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsAssetTypeSupported'
type MockdeployerCommander_IsAssetTypeSupported_Call struct {
	*mock.Call
}

// IsAssetTypeSupported is a helper method to define mock.On call
//   - assetType string
func (_e *MockdeployerCommander_Expecter) IsAssetTypeSupported(assetType interface{}) *MockdeployerCommander_IsAssetTypeSupported_Call {
	return &MockdeployerCommander_IsAssetTypeSupported_Call{Call: _e.mock.On("IsAssetTypeSupported", assetType)}
}

func (_c *MockdeployerCommander_IsAssetTypeSupported_Call) Run(run func(assetType string)) *MockdeployerCommander_IsAssetTypeSupported_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockdeployerCommander_IsAssetTypeSupported_Call) Return(_a0 bool) *MockdeployerCommander_IsAssetTypeSupported_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockdeployerCommander_IsAssetTypeSupported_Call) RunAndReturn(run func(string) bool) *MockdeployerCommander_IsAssetTypeSupported_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockdeployerCommander creates a new instance of MockdeployerCommander. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockdeployerCommander(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockdeployerCommander {
	mock := &MockdeployerCommander{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
