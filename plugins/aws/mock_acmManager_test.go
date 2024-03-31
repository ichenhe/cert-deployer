// Code generated by mockery v2.42.1. DO NOT EDIT.

package aws

import (
	context "context"

	domain "github.com/ichenhe/cert-deployer/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockacmManager is an autogenerated mock type for the acmManager type
type MockacmManager struct {
	mock.Mock
}

type MockacmManager_Expecter struct {
	mock *mock.Mock
}

func (_m *MockacmManager) EXPECT() *MockacmManager_Expecter {
	return &MockacmManager_Expecter{mock: &_m.Mock}
}

// DeleteManagedCertFromAcmIfUnused provides a mock function with given fields: ctx, certArn
func (_m *MockacmManager) DeleteManagedCertFromAcmIfUnused(ctx context.Context, certArn *string) (bool, error) {
	ret := _m.Called(ctx, certArn)

	if len(ret) == 0 {
		panic("no return value specified for DeleteManagedCertFromAcmIfUnused")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *string) (bool, error)); ok {
		return rf(ctx, certArn)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *string) bool); ok {
		r0 = rf(ctx, certArn)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *string) error); ok {
		r1 = rf(ctx, certArn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockacmManager_DeleteManagedCertFromAcmIfUnused_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteManagedCertFromAcmIfUnused'
type MockacmManager_DeleteManagedCertFromAcmIfUnused_Call struct {
	*mock.Call
}

// DeleteManagedCertFromAcmIfUnused is a helper method to define mock.On call
//   - ctx context.Context
//   - certArn *string
func (_e *MockacmManager_Expecter) DeleteManagedCertFromAcmIfUnused(ctx interface{}, certArn interface{}) *MockacmManager_DeleteManagedCertFromAcmIfUnused_Call {
	return &MockacmManager_DeleteManagedCertFromAcmIfUnused_Call{Call: _e.mock.On("DeleteManagedCertFromAcmIfUnused", ctx, certArn)}
}

func (_c *MockacmManager_DeleteManagedCertFromAcmIfUnused_Call) Run(run func(ctx context.Context, certArn *string)) *MockacmManager_DeleteManagedCertFromAcmIfUnused_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*string))
	})
	return _c
}

func (_c *MockacmManager_DeleteManagedCertFromAcmIfUnused_Call) Return(deleted bool, err error) *MockacmManager_DeleteManagedCertFromAcmIfUnused_Call {
	_c.Call.Return(deleted, err)
	return _c
}

func (_c *MockacmManager_DeleteManagedCertFromAcmIfUnused_Call) RunAndReturn(run func(context.Context, *string) (bool, error)) *MockacmManager_DeleteManagedCertFromAcmIfUnused_Call {
	_c.Call.Return(run)
	return _c
}

// FindCertInACM provides a mock function with given fields: ctx, certBundle
func (_m *MockacmManager) FindCertInACM(ctx context.Context, certBundle domain.CertificateBundle) (string, bool, error) {
	ret := _m.Called(ctx, certBundle)

	if len(ret) == 0 {
		panic("no return value specified for FindCertInACM")
	}

	var r0 string
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.CertificateBundle) (string, bool, error)); ok {
		return rf(ctx, certBundle)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.CertificateBundle) string); ok {
		r0 = rf(ctx, certBundle)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.CertificateBundle) bool); ok {
		r1 = rf(ctx, certBundle)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, domain.CertificateBundle) error); ok {
		r2 = rf(ctx, certBundle)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// MockacmManager_FindCertInACM_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindCertInACM'
type MockacmManager_FindCertInACM_Call struct {
	*mock.Call
}

// FindCertInACM is a helper method to define mock.On call
//   - ctx context.Context
//   - certBundle domain.CertificateBundle
func (_e *MockacmManager_Expecter) FindCertInACM(ctx interface{}, certBundle interface{}) *MockacmManager_FindCertInACM_Call {
	return &MockacmManager_FindCertInACM_Call{Call: _e.mock.On("FindCertInACM", ctx, certBundle)}
}

func (_c *MockacmManager_FindCertInACM_Call) Run(run func(ctx context.Context, certBundle domain.CertificateBundle)) *MockacmManager_FindCertInACM_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.CertificateBundle))
	})
	return _c
}

func (_c *MockacmManager_FindCertInACM_Call) Return(arn string, fromCache bool, err error) *MockacmManager_FindCertInACM_Call {
	_c.Call.Return(arn, fromCache, err)
	return _c
}

func (_c *MockacmManager_FindCertInACM_Call) RunAndReturn(run func(context.Context, domain.CertificateBundle) (string, bool, error)) *MockacmManager_FindCertInACM_Call {
	_c.Call.Return(run)
	return _c
}

// ImportCertificate provides a mock function with given fields: ctx, certBundle, key
func (_m *MockacmManager) ImportCertificate(ctx context.Context, certBundle domain.CertificateBundle, key []byte) (string, error) {
	ret := _m.Called(ctx, certBundle, key)

	if len(ret) == 0 {
		panic("no return value specified for ImportCertificate")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.CertificateBundle, []byte) (string, error)); ok {
		return rf(ctx, certBundle, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.CertificateBundle, []byte) string); ok {
		r0 = rf(ctx, certBundle, key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.CertificateBundle, []byte) error); ok {
		r1 = rf(ctx, certBundle, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockacmManager_ImportCertificate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ImportCertificate'
type MockacmManager_ImportCertificate_Call struct {
	*mock.Call
}

// ImportCertificate is a helper method to define mock.On call
//   - ctx context.Context
//   - certBundle domain.CertificateBundle
//   - key []byte
func (_e *MockacmManager_Expecter) ImportCertificate(ctx interface{}, certBundle interface{}, key interface{}) *MockacmManager_ImportCertificate_Call {
	return &MockacmManager_ImportCertificate_Call{Call: _e.mock.On("ImportCertificate", ctx, certBundle, key)}
}

func (_c *MockacmManager_ImportCertificate_Call) Run(run func(ctx context.Context, certBundle domain.CertificateBundle, key []byte)) *MockacmManager_ImportCertificate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.CertificateBundle), args[2].([]byte))
	})
	return _c
}

func (_c *MockacmManager_ImportCertificate_Call) Return(arn string, err error) *MockacmManager_ImportCertificate_Call {
	_c.Call.Return(arn, err)
	return _c
}

func (_c *MockacmManager_ImportCertificate_Call) RunAndReturn(run func(context.Context, domain.CertificateBundle, []byte) (string, error)) *MockacmManager_ImportCertificate_Call {
	_c.Call.Return(run)
	return _c
}

// RemoveCertFromCache provides a mock function with given fields: arn
func (_m *MockacmManager) RemoveCertFromCache(arn string) {
	_m.Called(arn)
}

// MockacmManager_RemoveCertFromCache_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RemoveCertFromCache'
type MockacmManager_RemoveCertFromCache_Call struct {
	*mock.Call
}

// RemoveCertFromCache is a helper method to define mock.On call
//   - arn string
func (_e *MockacmManager_Expecter) RemoveCertFromCache(arn interface{}) *MockacmManager_RemoveCertFromCache_Call {
	return &MockacmManager_RemoveCertFromCache_Call{Call: _e.mock.On("RemoveCertFromCache", arn)}
}

func (_c *MockacmManager_RemoveCertFromCache_Call) Run(run func(arn string)) *MockacmManager_RemoveCertFromCache_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockacmManager_RemoveCertFromCache_Call) Return() *MockacmManager_RemoveCertFromCache_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockacmManager_RemoveCertFromCache_Call) RunAndReturn(run func(string)) *MockacmManager_RemoveCertFromCache_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockacmManager creates a new instance of MockacmManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockacmManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockacmManager {
	mock := &MockacmManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}