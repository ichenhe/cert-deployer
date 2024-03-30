// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocker

import (
	context "context"

	domain "github.com/ichenhe/cert-deployer/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockDeployer is an autogenerated mock type for the Deployer type
type MockDeployer struct {
	mock.Mock
}

type MockDeployer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDeployer) EXPECT() *MockDeployer_Expecter {
	return &MockDeployer_Expecter{mock: &_m.Mock}
}

// Deploy provides a mock function with given fields: ctx, assets, cert, key, callback
func (_m *MockDeployer) Deploy(ctx context.Context, assets []domain.Asseter, cert []byte, key []byte, callback *domain.DeployCallback) error {
	ret := _m.Called(ctx, assets, cert, key, callback)

	if len(ret) == 0 {
		panic("no return value specified for Deploy")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []domain.Asseter, []byte, []byte, *domain.DeployCallback) error); ok {
		r0 = rf(ctx, assets, cert, key, callback)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDeployer_Deploy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Deploy'
type MockDeployer_Deploy_Call struct {
	*mock.Call
}

// Deploy is a helper method to define mock.On call
//   - ctx context.Context
//   - assets []domain.Asseter
//   - cert []byte
//   - key []byte
//   - callback *domain.DeployCallback
func (_e *MockDeployer_Expecter) Deploy(ctx interface{}, assets interface{}, cert interface{}, key interface{}, callback interface{}) *MockDeployer_Deploy_Call {
	return &MockDeployer_Deploy_Call{Call: _e.mock.On("Deploy", ctx, assets, cert, key, callback)}
}

func (_c *MockDeployer_Deploy_Call) Run(run func(ctx context.Context, assets []domain.Asseter, cert []byte, key []byte, callback *domain.DeployCallback)) *MockDeployer_Deploy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]domain.Asseter), args[2].([]byte), args[3].([]byte), args[4].(*domain.DeployCallback))
	})
	return _c
}

func (_c *MockDeployer_Deploy_Call) Return(_a0 error) *MockDeployer_Deploy_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDeployer_Deploy_Call) RunAndReturn(run func(context.Context, []domain.Asseter, []byte, []byte, *domain.DeployCallback) error) *MockDeployer_Deploy_Call {
	_c.Call.Return(run)
	return _c
}

// IsAssetTypeSupported provides a mock function with given fields: assetType
func (_m *MockDeployer) IsAssetTypeSupported(assetType string) bool {
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

// MockDeployer_IsAssetTypeSupported_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsAssetTypeSupported'
type MockDeployer_IsAssetTypeSupported_Call struct {
	*mock.Call
}

// IsAssetTypeSupported is a helper method to define mock.On call
//   - assetType string
func (_e *MockDeployer_Expecter) IsAssetTypeSupported(assetType interface{}) *MockDeployer_IsAssetTypeSupported_Call {
	return &MockDeployer_IsAssetTypeSupported_Call{Call: _e.mock.On("IsAssetTypeSupported", assetType)}
}

func (_c *MockDeployer_IsAssetTypeSupported_Call) Run(run func(assetType string)) *MockDeployer_IsAssetTypeSupported_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockDeployer_IsAssetTypeSupported_Call) Return(_a0 bool) *MockDeployer_IsAssetTypeSupported_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDeployer_IsAssetTypeSupported_Call) RunAndReturn(run func(string) bool) *MockDeployer_IsAssetTypeSupported_Call {
	_c.Call.Return(run)
	return _c
}

// ListApplicableAssets provides a mock function with given fields: ctx, assetType, cert
func (_m *MockDeployer) ListApplicableAssets(ctx context.Context, assetType string, cert []byte) ([]domain.Asseter, error) {
	ret := _m.Called(ctx, assetType, cert)

	if len(ret) == 0 {
		panic("no return value specified for ListApplicableAssets")
	}

	var r0 []domain.Asseter
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte) ([]domain.Asseter, error)); ok {
		return rf(ctx, assetType, cert)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte) []domain.Asseter); ok {
		r0 = rf(ctx, assetType, cert)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Asseter)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, []byte) error); ok {
		r1 = rf(ctx, assetType, cert)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDeployer_ListApplicableAssets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListApplicableAssets'
type MockDeployer_ListApplicableAssets_Call struct {
	*mock.Call
}

// ListApplicableAssets is a helper method to define mock.On call
//   - ctx context.Context
//   - assetType string
//   - cert []byte
func (_e *MockDeployer_Expecter) ListApplicableAssets(ctx interface{}, assetType interface{}, cert interface{}) *MockDeployer_ListApplicableAssets_Call {
	return &MockDeployer_ListApplicableAssets_Call{Call: _e.mock.On("ListApplicableAssets", ctx, assetType, cert)}
}

func (_c *MockDeployer_ListApplicableAssets_Call) Run(run func(ctx context.Context, assetType string, cert []byte)) *MockDeployer_ListApplicableAssets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]byte))
	})
	return _c
}

func (_c *MockDeployer_ListApplicableAssets_Call) Return(_a0 []domain.Asseter, _a1 error) *MockDeployer_ListApplicableAssets_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDeployer_ListApplicableAssets_Call) RunAndReturn(run func(context.Context, string, []byte) ([]domain.Asseter, error)) *MockDeployer_ListApplicableAssets_Call {
	_c.Call.Return(run)
	return _c
}

// ListAssets provides a mock function with given fields: ctx, assetType
func (_m *MockDeployer) ListAssets(ctx context.Context, assetType string) ([]domain.Asseter, error) {
	ret := _m.Called(ctx, assetType)

	if len(ret) == 0 {
		panic("no return value specified for ListAssets")
	}

	var r0 []domain.Asseter
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]domain.Asseter, error)); ok {
		return rf(ctx, assetType)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []domain.Asseter); ok {
		r0 = rf(ctx, assetType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Asseter)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, assetType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDeployer_ListAssets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAssets'
type MockDeployer_ListAssets_Call struct {
	*mock.Call
}

// ListAssets is a helper method to define mock.On call
//   - ctx context.Context
//   - assetType string
func (_e *MockDeployer_Expecter) ListAssets(ctx interface{}, assetType interface{}) *MockDeployer_ListAssets_Call {
	return &MockDeployer_ListAssets_Call{Call: _e.mock.On("ListAssets", ctx, assetType)}
}

func (_c *MockDeployer_ListAssets_Call) Run(run func(ctx context.Context, assetType string)) *MockDeployer_ListAssets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDeployer_ListAssets_Call) Return(_a0 []domain.Asseter, _a1 error) *MockDeployer_ListAssets_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDeployer_ListAssets_Call) RunAndReturn(run func(context.Context, string) ([]domain.Asseter, error)) *MockDeployer_ListAssets_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDeployer creates a new instance of MockDeployer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDeployer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDeployer {
	mock := &MockDeployer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
