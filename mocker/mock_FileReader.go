// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocker

import mock "github.com/stretchr/testify/mock"

// MockFileReader is an autogenerated mock type for the FileReader type
type MockFileReader struct {
	mock.Mock
}

type MockFileReader_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFileReader) EXPECT() *MockFileReader_Expecter {
	return &MockFileReader_Expecter{mock: &_m.Mock}
}

// ReadFile provides a mock function with given fields: name
func (_m *MockFileReader) ReadFile(name string) ([]byte, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for ReadFile")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]byte, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFileReader_ReadFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadFile'
type MockFileReader_ReadFile_Call struct {
	*mock.Call
}

// ReadFile is a helper method to define mock.On call
//   - name string
func (_e *MockFileReader_Expecter) ReadFile(name interface{}) *MockFileReader_ReadFile_Call {
	return &MockFileReader_ReadFile_Call{Call: _e.mock.On("ReadFile", name)}
}

func (_c *MockFileReader_ReadFile_Call) Run(run func(name string)) *MockFileReader_ReadFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockFileReader_ReadFile_Call) Return(_a0 []byte, _a1 error) *MockFileReader_ReadFile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFileReader_ReadFile_Call) RunAndReturn(run func(string) ([]byte, error)) *MockFileReader_ReadFile_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFileReader creates a new instance of MockFileReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFileReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFileReader {
	mock := &MockFileReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
