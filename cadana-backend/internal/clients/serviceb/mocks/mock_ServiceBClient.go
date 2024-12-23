// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	serviceb "cadana-backend/internal/clients/serviceb"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockServiceBClient is an autogenerated mock type for the ServiceBClient type
type MockServiceBClient struct {
	mock.Mock
}

// ExchangeRate provides a mock function with given fields: ctx, requestBody
func (_m *MockServiceBClient) ExchangeRate(ctx context.Context, requestBody serviceb.ExchangeRateRequestBody) (*serviceb.ExchangeRateResponse, error) {
	ret := _m.Called(ctx, requestBody)

	if len(ret) == 0 {
		panic("no return value specified for ExchangeRate")
	}

	var r0 *serviceb.ExchangeRateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, serviceb.ExchangeRateRequestBody) (*serviceb.ExchangeRateResponse, error)); ok {
		return rf(ctx, requestBody)
	}
	if rf, ok := ret.Get(0).(func(context.Context, serviceb.ExchangeRateRequestBody) *serviceb.ExchangeRateResponse); ok {
		r0 = rf(ctx, requestBody)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serviceb.ExchangeRateResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, serviceb.ExchangeRateRequestBody) error); ok {
		r1 = rf(ctx, requestBody)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockServiceBClient creates a new instance of MockServiceBClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockServiceBClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockServiceBClient {
	mock := &MockServiceBClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
