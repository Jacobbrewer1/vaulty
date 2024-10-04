// Code generated by mockery v2.46.1. DO NOT EDIT.

package repositories

import mock "github.com/stretchr/testify/mock"

// MockDatabaseConnector is an autogenerated mock type for the DatabaseConnector type
type MockDatabaseConnector struct {
	mock.Mock
}

// ConnectDB provides a mock function with given fields:
func (_m *MockDatabaseConnector) ConnectDB() (*Database, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ConnectDB")
	}

	var r0 *Database
	var r1 error
	if rf, ok := ret.Get(0).(func() (*Database, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *Database); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Database)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockDatabaseConnector creates a new instance of MockDatabaseConnector. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabaseConnector(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabaseConnector {
	mock := &MockDatabaseConnector{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
