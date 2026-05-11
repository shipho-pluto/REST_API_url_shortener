package mocks

import "github.com/stretchr/testify/mock"

type URLCacher struct {
	mock.Mock
}

func (_m *URLCacher) CacheURL(urlToSave string, alias string) error {
	ret := _m.Called(urlToSave, alias)

	var r0 error

	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(urlToSave, alias)
	} else {
		r0 = ret.Error(1)
	}

	return r0
}

func (_m *URLCacher) GetURL(alias string) (string, error) {
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

type mockConstructorTestingTNewURLCacher interface {
	mock.TestingT
	Cleanup(func())
}

func NewURLCacher(t mockConstructorTestingTNewURLCacher) *URLCacher {
	mock := &URLCacher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
