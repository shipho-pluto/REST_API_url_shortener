package mocks

import "github.com/stretchr/testify/mock"

type URLSaver struct {
	mock.Mock
}

func (_m *URLSaver) SaveURL(urlToSave string, alias string) error {
	ret := _m.Called(urlToSave, alias)

	var r0 error

	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(urlToSave, alias)
	} else {
		r0 = ret.Error(1)
	}

	return r0
}

type mockConstructorTestingTNewURLSaver interface {
	mock.TestingT
	Cleanup(func())
}

func NewURLSaver(t mockConstructorTestingTNewURLSaver) *URLSaver {
	mock := &URLSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
