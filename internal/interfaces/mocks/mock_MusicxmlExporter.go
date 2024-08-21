// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	musicmodel "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
	model "github.com/tomvodi/limepipes/internal/musicxml/model"
)

// MusicxmlExporter is an autogenerated mock type for the MusicxmlExporter type
type MusicxmlExporter struct {
	mock.Mock
}

type MusicxmlExporter_Expecter struct {
	mock *mock.Mock
}

func (_m *MusicxmlExporter) EXPECT() *MusicxmlExporter_Expecter {
	return &MusicxmlExporter_Expecter{mock: &_m.Mock}
}

// Export provides a mock function with given fields: musicModel
func (_m *MusicxmlExporter) Export(musicModel musicmodel.MusicModel) (*model.Score, error) {
	ret := _m.Called(musicModel)

	if len(ret) == 0 {
		panic("no return value specified for Export")
	}

	var r0 *model.Score
	var r1 error
	if rf, ok := ret.Get(0).(func(musicmodel.MusicModel) (*model.Score, error)); ok {
		return rf(musicModel)
	}
	if rf, ok := ret.Get(0).(func(musicmodel.MusicModel) *model.Score); ok {
		r0 = rf(musicModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Score)
		}
	}

	if rf, ok := ret.Get(1).(func(musicmodel.MusicModel) error); ok {
		r1 = rf(musicModel)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MusicxmlExporter_Export_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Export'
type MusicxmlExporter_Export_Call struct {
	*mock.Call
}

// Export is a helper method to define mock.On call
//   - musicModel musicmodel.MusicModel
func (_e *MusicxmlExporter_Expecter) Export(musicModel interface{}) *MusicxmlExporter_Export_Call {
	return &MusicxmlExporter_Export_Call{Call: _e.mock.On("Export", musicModel)}
}

func (_c *MusicxmlExporter_Export_Call) Run(run func(musicModel musicmodel.MusicModel)) *MusicxmlExporter_Export_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(musicmodel.MusicModel))
	})
	return _c
}

func (_c *MusicxmlExporter_Export_Call) Return(_a0 *model.Score, _a1 error) *MusicxmlExporter_Export_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MusicxmlExporter_Export_Call) RunAndReturn(run func(musicmodel.MusicModel) (*model.Score, error)) *MusicxmlExporter_Export_Call {
	_c.Call.Return(run)
	return _c
}

// NewMusicxmlExporter creates a new instance of MusicxmlExporter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMusicxmlExporter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MusicxmlExporter {
	mock := &MusicxmlExporter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
