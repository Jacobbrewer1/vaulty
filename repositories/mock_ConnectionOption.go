// Code generated by mockery v2.46.2. DO NOT EDIT.

package repositories

import mock "github.com/stretchr/testify/mock"

// MockConnectionOption is an autogenerated mock type for the ConnectionOption type
type MockConnectionOption struct {
	mock.Mock
}

// Execute provides a mock function with given fields: c
func (_m *MockConnectionOption) Execute(c *databaseConnector) {
	_m.Called(c)
}

// NewMockConnectionOption creates a new instance of MockConnectionOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConnectionOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConnectionOption {
	mock := &MockConnectionOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
