package mocks

import "github.com/stretchr/testify/mock"

type URLProvider struct {
	mock.Mock
}

func (_m *URLProvider) GetURL(alias string) (string, error) {
	ret := _m.Called(alias)

	var r0 string
	var r1 error

	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(alias)
	}

	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewURLProvider interface {
	mock.TestingT
	Cleanup(func())
}

func NewURLProvider(t mockConstructorTestingTNewURLProvider) *URLProvider {
	mock := &URLProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
