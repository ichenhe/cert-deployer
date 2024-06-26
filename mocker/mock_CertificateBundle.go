// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocker

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// MockCertificateBundle is an autogenerated mock type for the CertificateBundle type
type MockCertificateBundle struct {
	mock.Mock
}

type MockCertificateBundle_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCertificateBundle) EXPECT() *MockCertificateBundle_Expecter {
	return &MockCertificateBundle_Expecter{mock: &_m.Mock}
}

// ContainsAllDomains provides a mock function with given fields: domains
func (_m *MockCertificateBundle) ContainsAllDomains(domains []string) bool {
	ret := _m.Called(domains)

	if len(ret) == 0 {
		panic("no return value specified for ContainsAllDomains")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func([]string) bool); ok {
		r0 = rf(domains)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockCertificateBundle_ContainsAllDomains_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ContainsAllDomains'
type MockCertificateBundle_ContainsAllDomains_Call struct {
	*mock.Call
}

// ContainsAllDomains is a helper method to define mock.On call
//   - domains []string
func (_e *MockCertificateBundle_Expecter) ContainsAllDomains(domains interface{}) *MockCertificateBundle_ContainsAllDomains_Call {
	return &MockCertificateBundle_ContainsAllDomains_Call{Call: _e.mock.On("ContainsAllDomains", domains)}
}

func (_c *MockCertificateBundle_ContainsAllDomains_Call) Run(run func(domains []string)) *MockCertificateBundle_ContainsAllDomains_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string))
	})
	return _c
}

func (_c *MockCertificateBundle_ContainsAllDomains_Call) Return(_a0 bool) *MockCertificateBundle_ContainsAllDomains_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_ContainsAllDomains_Call) RunAndReturn(run func([]string) bool) *MockCertificateBundle_ContainsAllDomains_Call {
	_c.Call.Return(run)
	return _c
}

// GetDomains provides a mock function with given fields:
func (_m *MockCertificateBundle) GetDomains() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDomains")
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

// MockCertificateBundle_GetDomains_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDomains'
type MockCertificateBundle_GetDomains_Call struct {
	*mock.Call
}

// GetDomains is a helper method to define mock.On call
func (_e *MockCertificateBundle_Expecter) GetDomains() *MockCertificateBundle_GetDomains_Call {
	return &MockCertificateBundle_GetDomains_Call{Call: _e.mock.On("GetDomains")}
}

func (_c *MockCertificateBundle_GetDomains_Call) Run(run func()) *MockCertificateBundle_GetDomains_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCertificateBundle_GetDomains_Call) Return(_a0 []string) *MockCertificateBundle_GetDomains_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_GetDomains_Call) RunAndReturn(run func() []string) *MockCertificateBundle_GetDomains_Call {
	_c.Call.Return(run)
	return _c
}

// GetRaw provides a mock function with given fields:
func (_m *MockCertificateBundle) GetRaw() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRaw")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// MockCertificateBundle_GetRaw_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRaw'
type MockCertificateBundle_GetRaw_Call struct {
	*mock.Call
}

// GetRaw is a helper method to define mock.On call
func (_e *MockCertificateBundle_Expecter) GetRaw() *MockCertificateBundle_GetRaw_Call {
	return &MockCertificateBundle_GetRaw_Call{Call: _e.mock.On("GetRaw")}
}

func (_c *MockCertificateBundle_GetRaw_Call) Run(run func()) *MockCertificateBundle_GetRaw_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCertificateBundle_GetRaw_Call) Return(_a0 []byte) *MockCertificateBundle_GetRaw_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_GetRaw_Call) RunAndReturn(run func() []byte) *MockCertificateBundle_GetRaw_Call {
	_c.Call.Return(run)
	return _c
}

// GetRawCert provides a mock function with given fields:
func (_m *MockCertificateBundle) GetRawCert() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRawCert")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// MockCertificateBundle_GetRawCert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRawCert'
type MockCertificateBundle_GetRawCert_Call struct {
	*mock.Call
}

// GetRawCert is a helper method to define mock.On call
func (_e *MockCertificateBundle_Expecter) GetRawCert() *MockCertificateBundle_GetRawCert_Call {
	return &MockCertificateBundle_GetRawCert_Call{Call: _e.mock.On("GetRawCert")}
}

func (_c *MockCertificateBundle_GetRawCert_Call) Run(run func()) *MockCertificateBundle_GetRawCert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCertificateBundle_GetRawCert_Call) Return(_a0 []byte) *MockCertificateBundle_GetRawCert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_GetRawCert_Call) RunAndReturn(run func() []byte) *MockCertificateBundle_GetRawCert_Call {
	_c.Call.Return(run)
	return _c
}

// GetRawChain provides a mock function with given fields:
func (_m *MockCertificateBundle) GetRawChain() []byte {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRawChain")
	}

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// MockCertificateBundle_GetRawChain_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRawChain'
type MockCertificateBundle_GetRawChain_Call struct {
	*mock.Call
}

// GetRawChain is a helper method to define mock.On call
func (_e *MockCertificateBundle_Expecter) GetRawChain() *MockCertificateBundle_GetRawChain_Call {
	return &MockCertificateBundle_GetRawChain_Call{Call: _e.mock.On("GetRawChain")}
}

func (_c *MockCertificateBundle_GetRawChain_Call) Run(run func()) *MockCertificateBundle_GetRawChain_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCertificateBundle_GetRawChain_Call) Return(_a0 []byte) *MockCertificateBundle_GetRawChain_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_GetRawChain_Call) RunAndReturn(run func() []byte) *MockCertificateBundle_GetRawChain_Call {
	_c.Call.Return(run)
	return _c
}

// GetSerialNumberHexString provides a mock function with given fields:
func (_m *MockCertificateBundle) GetSerialNumberHexString() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetSerialNumberHexString")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockCertificateBundle_GetSerialNumberHexString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSerialNumberHexString'
type MockCertificateBundle_GetSerialNumberHexString_Call struct {
	*mock.Call
}

// GetSerialNumberHexString is a helper method to define mock.On call
func (_e *MockCertificateBundle_Expecter) GetSerialNumberHexString() *MockCertificateBundle_GetSerialNumberHexString_Call {
	return &MockCertificateBundle_GetSerialNumberHexString_Call{Call: _e.mock.On("GetSerialNumberHexString")}
}

func (_c *MockCertificateBundle_GetSerialNumberHexString_Call) Run(run func()) *MockCertificateBundle_GetSerialNumberHexString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCertificateBundle_GetSerialNumberHexString_Call) Return(_a0 string) *MockCertificateBundle_GetSerialNumberHexString_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_GetSerialNumberHexString_Call) RunAndReturn(run func() string) *MockCertificateBundle_GetSerialNumberHexString_Call {
	_c.Call.Return(run)
	return _c
}

// NotAfter provides a mock function with given fields:
func (_m *MockCertificateBundle) NotAfter() *time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NotAfter")
	}

	var r0 *time.Time
	if rf, ok := ret.Get(0).(func() *time.Time); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*time.Time)
		}
	}

	return r0
}

// MockCertificateBundle_NotAfter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NotAfter'
type MockCertificateBundle_NotAfter_Call struct {
	*mock.Call
}

// NotAfter is a helper method to define mock.On call
func (_e *MockCertificateBundle_Expecter) NotAfter() *MockCertificateBundle_NotAfter_Call {
	return &MockCertificateBundle_NotAfter_Call{Call: _e.mock.On("NotAfter")}
}

func (_c *MockCertificateBundle_NotAfter_Call) Run(run func()) *MockCertificateBundle_NotAfter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCertificateBundle_NotAfter_Call) Return(_a0 *time.Time) *MockCertificateBundle_NotAfter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_NotAfter_Call) RunAndReturn(run func() *time.Time) *MockCertificateBundle_NotAfter_Call {
	_c.Call.Return(run)
	return _c
}

// NotBefore provides a mock function with given fields:
func (_m *MockCertificateBundle) NotBefore() *time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NotBefore")
	}

	var r0 *time.Time
	if rf, ok := ret.Get(0).(func() *time.Time); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*time.Time)
		}
	}

	return r0
}

// MockCertificateBundle_NotBefore_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NotBefore'
type MockCertificateBundle_NotBefore_Call struct {
	*mock.Call
}

// NotBefore is a helper method to define mock.On call
func (_e *MockCertificateBundle_Expecter) NotBefore() *MockCertificateBundle_NotBefore_Call {
	return &MockCertificateBundle_NotBefore_Call{Call: _e.mock.On("NotBefore")}
}

func (_c *MockCertificateBundle_NotBefore_Call) Run(run func()) *MockCertificateBundle_NotBefore_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCertificateBundle_NotBefore_Call) Return(_a0 *time.Time) *MockCertificateBundle_NotBefore_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_NotBefore_Call) RunAndReturn(run func() *time.Time) *MockCertificateBundle_NotBefore_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyHostname provides a mock function with given fields: hostname
func (_m *MockCertificateBundle) VerifyHostname(hostname string) bool {
	ret := _m.Called(hostname)

	if len(ret) == 0 {
		panic("no return value specified for VerifyHostname")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(hostname)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockCertificateBundle_VerifyHostname_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyHostname'
type MockCertificateBundle_VerifyHostname_Call struct {
	*mock.Call
}

// VerifyHostname is a helper method to define mock.On call
//   - hostname string
func (_e *MockCertificateBundle_Expecter) VerifyHostname(hostname interface{}) *MockCertificateBundle_VerifyHostname_Call {
	return &MockCertificateBundle_VerifyHostname_Call{Call: _e.mock.On("VerifyHostname", hostname)}
}

func (_c *MockCertificateBundle_VerifyHostname_Call) Run(run func(hostname string)) *MockCertificateBundle_VerifyHostname_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCertificateBundle_VerifyHostname_Call) Return(_a0 bool) *MockCertificateBundle_VerifyHostname_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_VerifyHostname_Call) RunAndReturn(run func(string) bool) *MockCertificateBundle_VerifyHostname_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyHostnames provides a mock function with given fields: hostnames
func (_m *MockCertificateBundle) VerifyHostnames(hostnames []string) bool {
	ret := _m.Called(hostnames)

	if len(ret) == 0 {
		panic("no return value specified for VerifyHostnames")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func([]string) bool); ok {
		r0 = rf(hostnames)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockCertificateBundle_VerifyHostnames_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyHostnames'
type MockCertificateBundle_VerifyHostnames_Call struct {
	*mock.Call
}

// VerifyHostnames is a helper method to define mock.On call
//   - hostnames []string
func (_e *MockCertificateBundle_Expecter) VerifyHostnames(hostnames interface{}) *MockCertificateBundle_VerifyHostnames_Call {
	return &MockCertificateBundle_VerifyHostnames_Call{Call: _e.mock.On("VerifyHostnames", hostnames)}
}

func (_c *MockCertificateBundle_VerifyHostnames_Call) Run(run func(hostnames []string)) *MockCertificateBundle_VerifyHostnames_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string))
	})
	return _c
}

func (_c *MockCertificateBundle_VerifyHostnames_Call) Return(_a0 bool) *MockCertificateBundle_VerifyHostnames_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_VerifyHostnames_Call) RunAndReturn(run func([]string) bool) *MockCertificateBundle_VerifyHostnames_Call {
	_c.Call.Return(run)
	return _c
}

// VerifySerialNumber provides a mock function with given fields: serialNumber
func (_m *MockCertificateBundle) VerifySerialNumber(serialNumber string) bool {
	ret := _m.Called(serialNumber)

	if len(ret) == 0 {
		panic("no return value specified for VerifySerialNumber")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(serialNumber)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockCertificateBundle_VerifySerialNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifySerialNumber'
type MockCertificateBundle_VerifySerialNumber_Call struct {
	*mock.Call
}

// VerifySerialNumber is a helper method to define mock.On call
//   - serialNumber string
func (_e *MockCertificateBundle_Expecter) VerifySerialNumber(serialNumber interface{}) *MockCertificateBundle_VerifySerialNumber_Call {
	return &MockCertificateBundle_VerifySerialNumber_Call{Call: _e.mock.On("VerifySerialNumber", serialNumber)}
}

func (_c *MockCertificateBundle_VerifySerialNumber_Call) Run(run func(serialNumber string)) *MockCertificateBundle_VerifySerialNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCertificateBundle_VerifySerialNumber_Call) Return(_a0 bool) *MockCertificateBundle_VerifySerialNumber_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCertificateBundle_VerifySerialNumber_Call) RunAndReturn(run func(string) bool) *MockCertificateBundle_VerifySerialNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCertificateBundle creates a new instance of MockCertificateBundle. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCertificateBundle(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCertificateBundle {
	mock := &MockCertificateBundle{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
