// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// NonceStore is an autogenerated mock type for the NonceStore type
type NonceStore struct {
	mock.Mock
}

// Accept provides a mock function with given fields: endpoint, nonce
func (_m *NonceStore) Accept(endpoint string, nonce string) error {
	ret := _m.Called(endpoint, nonce)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(endpoint, nonce)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewNonceStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewNonceStore creates a new instance of NonceStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewNonceStore(t mockConstructorTestingTNewNonceStore) *NonceStore {
	mock := &NonceStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
