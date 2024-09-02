// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	apimodel "github.com/tomvodi/limepipes/internal/apigen/apimodel"
	common "github.com/tomvodi/limepipes/internal/common"

	file_type "github.com/tomvodi/limepipes-plugin-api/plugin/v1/file_type"

	messages "github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"

	mock "github.com/stretchr/testify/mock"

	model "github.com/tomvodi/limepipes/internal/database/model"

	uuid "github.com/google/uuid"
)

// DataService is an autogenerated mock type for the DataService type
type DataService struct {
	mock.Mock
}

type DataService_Expecter struct {
	mock *mock.Mock
}

func (_m *DataService) EXPECT() *DataService_Expecter {
	return &DataService_Expecter{mock: &_m.Mock}
}

// AddFileToTune provides a mock function with given fields: tuneID, tFile
func (_m *DataService) AddFileToTune(tuneID uuid.UUID, tFile *model.TuneFile) error {
	ret := _m.Called(tuneID, tFile)

	if len(ret) == 0 {
		panic("no return value specified for AddFileToTune")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, *model.TuneFile) error); ok {
		r0 = rf(tuneID, tFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DataService_AddFileToTune_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddFileToTune'
type DataService_AddFileToTune_Call struct {
	*mock.Call
}

// AddFileToTune is a helper method to define mock.On call
//   - tuneID uuid.UUID
//   - tFile *model.TuneFile
func (_e *DataService_Expecter) AddFileToTune(tuneID interface{}, tFile interface{}) *DataService_AddFileToTune_Call {
	return &DataService_AddFileToTune_Call{Call: _e.mock.On("AddFileToTune", tuneID, tFile)}
}

func (_c *DataService_AddFileToTune_Call) Run(run func(tuneID uuid.UUID, tFile *model.TuneFile)) *DataService_AddFileToTune_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(*model.TuneFile))
	})
	return _c
}

func (_c *DataService_AddFileToTune_Call) Return(_a0 error) *DataService_AddFileToTune_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataService_AddFileToTune_Call) RunAndReturn(run func(uuid.UUID, *model.TuneFile) error) *DataService_AddFileToTune_Call {
	_c.Call.Return(run)
	return _c
}

// AssignTunesToMusicSet provides a mock function with given fields: setID, tuneIDs
func (_m *DataService) AssignTunesToMusicSet(setID uuid.UUID, tuneIDs []uuid.UUID) (*apimodel.MusicSet, error) {
	ret := _m.Called(setID, tuneIDs)

	if len(ret) == 0 {
		panic("no return value specified for AssignTunesToMusicSet")
	}

	var r0 *apimodel.MusicSet
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, []uuid.UUID) (*apimodel.MusicSet, error)); ok {
		return rf(setID, tuneIDs)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, []uuid.UUID) *apimodel.MusicSet); ok {
		r0 = rf(setID, tuneIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apimodel.MusicSet)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, []uuid.UUID) error); ok {
		r1 = rf(setID, tuneIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_AssignTunesToMusicSet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AssignTunesToMusicSet'
type DataService_AssignTunesToMusicSet_Call struct {
	*mock.Call
}

// AssignTunesToMusicSet is a helper method to define mock.On call
//   - setID uuid.UUID
//   - tuneIDs []uuid.UUID
func (_e *DataService_Expecter) AssignTunesToMusicSet(setID interface{}, tuneIDs interface{}) *DataService_AssignTunesToMusicSet_Call {
	return &DataService_AssignTunesToMusicSet_Call{Call: _e.mock.On("AssignTunesToMusicSet", setID, tuneIDs)}
}

func (_c *DataService_AssignTunesToMusicSet_Call) Run(run func(setID uuid.UUID, tuneIDs []uuid.UUID)) *DataService_AssignTunesToMusicSet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].([]uuid.UUID))
	})
	return _c
}

func (_c *DataService_AssignTunesToMusicSet_Call) Return(_a0 *apimodel.MusicSet, _a1 error) *DataService_AssignTunesToMusicSet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_AssignTunesToMusicSet_Call) RunAndReturn(run func(uuid.UUID, []uuid.UUID) (*apimodel.MusicSet, error)) *DataService_AssignTunesToMusicSet_Call {
	_c.Call.Return(run)
	return _c
}

// CreateMusicSet provides a mock function with given fields: tune, importFile
func (_m *DataService) CreateMusicSet(tune apimodel.CreateSet, importFile *model.ImportFile) (*apimodel.MusicSet, error) {
	ret := _m.Called(tune, importFile)

	if len(ret) == 0 {
		panic("no return value specified for CreateMusicSet")
	}

	var r0 *apimodel.MusicSet
	var r1 error
	if rf, ok := ret.Get(0).(func(apimodel.CreateSet, *model.ImportFile) (*apimodel.MusicSet, error)); ok {
		return rf(tune, importFile)
	}
	if rf, ok := ret.Get(0).(func(apimodel.CreateSet, *model.ImportFile) *apimodel.MusicSet); ok {
		r0 = rf(tune, importFile)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apimodel.MusicSet)
		}
	}

	if rf, ok := ret.Get(1).(func(apimodel.CreateSet, *model.ImportFile) error); ok {
		r1 = rf(tune, importFile)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_CreateMusicSet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMusicSet'
type DataService_CreateMusicSet_Call struct {
	*mock.Call
}

// CreateMusicSet is a helper method to define mock.On call
//   - tune apimodel.CreateSet
//   - importFile *model.ImportFile
func (_e *DataService_Expecter) CreateMusicSet(tune interface{}, importFile interface{}) *DataService_CreateMusicSet_Call {
	return &DataService_CreateMusicSet_Call{Call: _e.mock.On("CreateMusicSet", tune, importFile)}
}

func (_c *DataService_CreateMusicSet_Call) Run(run func(tune apimodel.CreateSet, importFile *model.ImportFile)) *DataService_CreateMusicSet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(apimodel.CreateSet), args[1].(*model.ImportFile))
	})
	return _c
}

func (_c *DataService_CreateMusicSet_Call) Return(_a0 *apimodel.MusicSet, _a1 error) *DataService_CreateMusicSet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_CreateMusicSet_Call) RunAndReturn(run func(apimodel.CreateSet, *model.ImportFile) (*apimodel.MusicSet, error)) *DataService_CreateMusicSet_Call {
	_c.Call.Return(run)
	return _c
}

// CreateTune provides a mock function with given fields: tune, importFile
func (_m *DataService) CreateTune(tune apimodel.CreateTune, importFile *model.ImportFile) (*apimodel.Tune, error) {
	ret := _m.Called(tune, importFile)

	if len(ret) == 0 {
		panic("no return value specified for CreateTune")
	}

	var r0 *apimodel.Tune
	var r1 error
	if rf, ok := ret.Get(0).(func(apimodel.CreateTune, *model.ImportFile) (*apimodel.Tune, error)); ok {
		return rf(tune, importFile)
	}
	if rf, ok := ret.Get(0).(func(apimodel.CreateTune, *model.ImportFile) *apimodel.Tune); ok {
		r0 = rf(tune, importFile)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apimodel.Tune)
		}
	}

	if rf, ok := ret.Get(1).(func(apimodel.CreateTune, *model.ImportFile) error); ok {
		r1 = rf(tune, importFile)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_CreateTune_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateTune'
type DataService_CreateTune_Call struct {
	*mock.Call
}

// CreateTune is a helper method to define mock.On call
//   - tune apimodel.CreateTune
//   - importFile *model.ImportFile
func (_e *DataService_Expecter) CreateTune(tune interface{}, importFile interface{}) *DataService_CreateTune_Call {
	return &DataService_CreateTune_Call{Call: _e.mock.On("CreateTune", tune, importFile)}
}

func (_c *DataService_CreateTune_Call) Run(run func(tune apimodel.CreateTune, importFile *model.ImportFile)) *DataService_CreateTune_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(apimodel.CreateTune), args[1].(*model.ImportFile))
	})
	return _c
}

func (_c *DataService_CreateTune_Call) Return(_a0 *apimodel.Tune, _a1 error) *DataService_CreateTune_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_CreateTune_Call) RunAndReturn(run func(apimodel.CreateTune, *model.ImportFile) (*apimodel.Tune, error)) *DataService_CreateTune_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteFileFromTune provides a mock function with given fields: tuneID, fType
func (_m *DataService) DeleteFileFromTune(tuneID uuid.UUID, fType file_type.Type) error {
	ret := _m.Called(tuneID, fType)

	if len(ret) == 0 {
		panic("no return value specified for DeleteFileFromTune")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, file_type.Type) error); ok {
		r0 = rf(tuneID, fType)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DataService_DeleteFileFromTune_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteFileFromTune'
type DataService_DeleteFileFromTune_Call struct {
	*mock.Call
}

// DeleteFileFromTune is a helper method to define mock.On call
//   - tuneID uuid.UUID
//   - fType file_type.Type
func (_e *DataService_Expecter) DeleteFileFromTune(tuneID interface{}, fType interface{}) *DataService_DeleteFileFromTune_Call {
	return &DataService_DeleteFileFromTune_Call{Call: _e.mock.On("DeleteFileFromTune", tuneID, fType)}
}

func (_c *DataService_DeleteFileFromTune_Call) Run(run func(tuneID uuid.UUID, fType file_type.Type)) *DataService_DeleteFileFromTune_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(file_type.Type))
	})
	return _c
}

func (_c *DataService_DeleteFileFromTune_Call) Return(_a0 error) *DataService_DeleteFileFromTune_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataService_DeleteFileFromTune_Call) RunAndReturn(run func(uuid.UUID, file_type.Type) error) *DataService_DeleteFileFromTune_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteMusicSet provides a mock function with given fields: id
func (_m *DataService) DeleteMusicSet(id uuid.UUID) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMusicSet")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DataService_DeleteMusicSet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteMusicSet'
type DataService_DeleteMusicSet_Call struct {
	*mock.Call
}

// DeleteMusicSet is a helper method to define mock.On call
//   - id uuid.UUID
func (_e *DataService_Expecter) DeleteMusicSet(id interface{}) *DataService_DeleteMusicSet_Call {
	return &DataService_DeleteMusicSet_Call{Call: _e.mock.On("DeleteMusicSet", id)}
}

func (_c *DataService_DeleteMusicSet_Call) Run(run func(id uuid.UUID)) *DataService_DeleteMusicSet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *DataService_DeleteMusicSet_Call) Return(_a0 error) *DataService_DeleteMusicSet_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataService_DeleteMusicSet_Call) RunAndReturn(run func(uuid.UUID) error) *DataService_DeleteMusicSet_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteTune provides a mock function with given fields: id
func (_m *DataService) DeleteTune(id uuid.UUID) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTune")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DataService_DeleteTune_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteTune'
type DataService_DeleteTune_Call struct {
	*mock.Call
}

// DeleteTune is a helper method to define mock.On call
//   - id uuid.UUID
func (_e *DataService_Expecter) DeleteTune(id interface{}) *DataService_DeleteTune_Call {
	return &DataService_DeleteTune_Call{Call: _e.mock.On("DeleteTune", id)}
}

func (_c *DataService_DeleteTune_Call) Run(run func(id uuid.UUID)) *DataService_DeleteTune_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *DataService_DeleteTune_Call) Return(_a0 error) *DataService_DeleteTune_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DataService_DeleteTune_Call) RunAndReturn(run func(uuid.UUID) error) *DataService_DeleteTune_Call {
	_c.Call.Return(run)
	return _c
}

// GetImportFileByHash provides a mock function with given fields: fHash
func (_m *DataService) GetImportFileByHash(fHash string) (*model.ImportFile, error) {
	ret := _m.Called(fHash)

	if len(ret) == 0 {
		panic("no return value specified for GetImportFileByHash")
	}

	var r0 *model.ImportFile
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*model.ImportFile, error)); ok {
		return rf(fHash)
	}
	if rf, ok := ret.Get(0).(func(string) *model.ImportFile); ok {
		r0 = rf(fHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ImportFile)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(fHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_GetImportFileByHash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetImportFileByHash'
type DataService_GetImportFileByHash_Call struct {
	*mock.Call
}

// GetImportFileByHash is a helper method to define mock.On call
//   - fHash string
func (_e *DataService_Expecter) GetImportFileByHash(fHash interface{}) *DataService_GetImportFileByHash_Call {
	return &DataService_GetImportFileByHash_Call{Call: _e.mock.On("GetImportFileByHash", fHash)}
}

func (_c *DataService_GetImportFileByHash_Call) Run(run func(fHash string)) *DataService_GetImportFileByHash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *DataService_GetImportFileByHash_Call) Return(_a0 *model.ImportFile, _a1 error) *DataService_GetImportFileByHash_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_GetImportFileByHash_Call) RunAndReturn(run func(string) (*model.ImportFile, error)) *DataService_GetImportFileByHash_Call {
	_c.Call.Return(run)
	return _c
}

// GetMusicSet provides a mock function with given fields: id
func (_m *DataService) GetMusicSet(id uuid.UUID) (*apimodel.MusicSet, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetMusicSet")
	}

	var r0 *apimodel.MusicSet
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*apimodel.MusicSet, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *apimodel.MusicSet); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apimodel.MusicSet)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_GetMusicSet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMusicSet'
type DataService_GetMusicSet_Call struct {
	*mock.Call
}

// GetMusicSet is a helper method to define mock.On call
//   - id uuid.UUID
func (_e *DataService_Expecter) GetMusicSet(id interface{}) *DataService_GetMusicSet_Call {
	return &DataService_GetMusicSet_Call{Call: _e.mock.On("GetMusicSet", id)}
}

func (_c *DataService_GetMusicSet_Call) Run(run func(id uuid.UUID)) *DataService_GetMusicSet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *DataService_GetMusicSet_Call) Return(_a0 *apimodel.MusicSet, _a1 error) *DataService_GetMusicSet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_GetMusicSet_Call) RunAndReturn(run func(uuid.UUID) (*apimodel.MusicSet, error)) *DataService_GetMusicSet_Call {
	_c.Call.Return(run)
	return _c
}

// GetTune provides a mock function with given fields: id
func (_m *DataService) GetTune(id uuid.UUID) (*apimodel.Tune, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetTune")
	}

	var r0 *apimodel.Tune
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) (*apimodel.Tune, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) *apimodel.Tune); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apimodel.Tune)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_GetTune_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTune'
type DataService_GetTune_Call struct {
	*mock.Call
}

// GetTune is a helper method to define mock.On call
//   - id uuid.UUID
func (_e *DataService_Expecter) GetTune(id interface{}) *DataService_GetTune_Call {
	return &DataService_GetTune_Call{Call: _e.mock.On("GetTune", id)}
}

func (_c *DataService_GetTune_Call) Run(run func(id uuid.UUID)) *DataService_GetTune_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *DataService_GetTune_Call) Return(_a0 *apimodel.Tune, _a1 error) *DataService_GetTune_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_GetTune_Call) RunAndReturn(run func(uuid.UUID) (*apimodel.Tune, error)) *DataService_GetTune_Call {
	_c.Call.Return(run)
	return _c
}

// GetTuneFile provides a mock function with given fields: tuneID, fType
func (_m *DataService) GetTuneFile(tuneID uuid.UUID, fType file_type.Type) (*model.TuneFile, error) {
	ret := _m.Called(tuneID, fType)

	if len(ret) == 0 {
		panic("no return value specified for GetTuneFile")
	}

	var r0 *model.TuneFile
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, file_type.Type) (*model.TuneFile, error)); ok {
		return rf(tuneID, fType)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, file_type.Type) *model.TuneFile); ok {
		r0 = rf(tuneID, fType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.TuneFile)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, file_type.Type) error); ok {
		r1 = rf(tuneID, fType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_GetTuneFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTuneFile'
type DataService_GetTuneFile_Call struct {
	*mock.Call
}

// GetTuneFile is a helper method to define mock.On call
//   - tuneID uuid.UUID
//   - fType file_type.Type
func (_e *DataService_Expecter) GetTuneFile(tuneID interface{}, fType interface{}) *DataService_GetTuneFile_Call {
	return &DataService_GetTuneFile_Call{Call: _e.mock.On("GetTuneFile", tuneID, fType)}
}

func (_c *DataService_GetTuneFile_Call) Run(run func(tuneID uuid.UUID, fType file_type.Type)) *DataService_GetTuneFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(file_type.Type))
	})
	return _c
}

func (_c *DataService_GetTuneFile_Call) Return(_a0 *model.TuneFile, _a1 error) *DataService_GetTuneFile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_GetTuneFile_Call) RunAndReturn(run func(uuid.UUID, file_type.Type) (*model.TuneFile, error)) *DataService_GetTuneFile_Call {
	_c.Call.Return(run)
	return _c
}

// GetTuneFiles provides a mock function with given fields: tuneID
func (_m *DataService) GetTuneFiles(tuneID uuid.UUID) ([]*model.TuneFile, error) {
	ret := _m.Called(tuneID)

	if len(ret) == 0 {
		panic("no return value specified for GetTuneFiles")
	}

	var r0 []*model.TuneFile
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) ([]*model.TuneFile, error)); ok {
		return rf(tuneID)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID) []*model.TuneFile); ok {
		r0 = rf(tuneID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.TuneFile)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(tuneID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_GetTuneFiles_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTuneFiles'
type DataService_GetTuneFiles_Call struct {
	*mock.Call
}

// GetTuneFiles is a helper method to define mock.On call
//   - tuneID uuid.UUID
func (_e *DataService_Expecter) GetTuneFiles(tuneID interface{}) *DataService_GetTuneFiles_Call {
	return &DataService_GetTuneFiles_Call{Call: _e.mock.On("GetTuneFiles", tuneID)}
}

func (_c *DataService_GetTuneFiles_Call) Run(run func(tuneID uuid.UUID)) *DataService_GetTuneFiles_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID))
	})
	return _c
}

func (_c *DataService_GetTuneFiles_Call) Return(_a0 []*model.TuneFile, _a1 error) *DataService_GetTuneFiles_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_GetTuneFiles_Call) RunAndReturn(run func(uuid.UUID) ([]*model.TuneFile, error)) *DataService_GetTuneFiles_Call {
	_c.Call.Return(run)
	return _c
}

// ImportTunes provides a mock function with given fields: tunes, fileInfo
func (_m *DataService) ImportTunes(tunes []*messages.ImportedTune, fileInfo *common.ImportFileInfo) ([]*apimodel.ImportTune, *apimodel.BasicMusicSet, error) {
	ret := _m.Called(tunes, fileInfo)

	if len(ret) == 0 {
		panic("no return value specified for ImportTunes")
	}

	var r0 []*apimodel.ImportTune
	var r1 *apimodel.BasicMusicSet
	var r2 error
	if rf, ok := ret.Get(0).(func([]*messages.ImportedTune, *common.ImportFileInfo) ([]*apimodel.ImportTune, *apimodel.BasicMusicSet, error)); ok {
		return rf(tunes, fileInfo)
	}
	if rf, ok := ret.Get(0).(func([]*messages.ImportedTune, *common.ImportFileInfo) []*apimodel.ImportTune); ok {
		r0 = rf(tunes, fileInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*apimodel.ImportTune)
		}
	}

	if rf, ok := ret.Get(1).(func([]*messages.ImportedTune, *common.ImportFileInfo) *apimodel.BasicMusicSet); ok {
		r1 = rf(tunes, fileInfo)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*apimodel.BasicMusicSet)
		}
	}

	if rf, ok := ret.Get(2).(func([]*messages.ImportedTune, *common.ImportFileInfo) error); ok {
		r2 = rf(tunes, fileInfo)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DataService_ImportTunes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ImportTunes'
type DataService_ImportTunes_Call struct {
	*mock.Call
}

// ImportTunes is a helper method to define mock.On call
//   - tunes []*messages.ImportedTune
//   - fileInfo *common.ImportFileInfo
func (_e *DataService_Expecter) ImportTunes(tunes interface{}, fileInfo interface{}) *DataService_ImportTunes_Call {
	return &DataService_ImportTunes_Call{Call: _e.mock.On("ImportTunes", tunes, fileInfo)}
}

func (_c *DataService_ImportTunes_Call) Run(run func(tunes []*messages.ImportedTune, fileInfo *common.ImportFileInfo)) *DataService_ImportTunes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]*messages.ImportedTune), args[1].(*common.ImportFileInfo))
	})
	return _c
}

func (_c *DataService_ImportTunes_Call) Return(_a0 []*apimodel.ImportTune, _a1 *apimodel.BasicMusicSet, _a2 error) *DataService_ImportTunes_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *DataService_ImportTunes_Call) RunAndReturn(run func([]*messages.ImportedTune, *common.ImportFileInfo) ([]*apimodel.ImportTune, *apimodel.BasicMusicSet, error)) *DataService_ImportTunes_Call {
	_c.Call.Return(run)
	return _c
}

// MusicSets provides a mock function with given fields:
func (_m *DataService) MusicSets() ([]*apimodel.MusicSet, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MusicSets")
	}

	var r0 []*apimodel.MusicSet
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*apimodel.MusicSet, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*apimodel.MusicSet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*apimodel.MusicSet)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_MusicSets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MusicSets'
type DataService_MusicSets_Call struct {
	*mock.Call
}

// MusicSets is a helper method to define mock.On call
func (_e *DataService_Expecter) MusicSets() *DataService_MusicSets_Call {
	return &DataService_MusicSets_Call{Call: _e.mock.On("MusicSets")}
}

func (_c *DataService_MusicSets_Call) Run(run func()) *DataService_MusicSets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DataService_MusicSets_Call) Return(_a0 []*apimodel.MusicSet, _a1 error) *DataService_MusicSets_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_MusicSets_Call) RunAndReturn(run func() ([]*apimodel.MusicSet, error)) *DataService_MusicSets_Call {
	_c.Call.Return(run)
	return _c
}

// Tunes provides a mock function with given fields:
func (_m *DataService) Tunes() ([]*apimodel.Tune, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Tunes")
	}

	var r0 []*apimodel.Tune
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*apimodel.Tune, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*apimodel.Tune); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*apimodel.Tune)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_Tunes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Tunes'
type DataService_Tunes_Call struct {
	*mock.Call
}

// Tunes is a helper method to define mock.On call
func (_e *DataService_Expecter) Tunes() *DataService_Tunes_Call {
	return &DataService_Tunes_Call{Call: _e.mock.On("Tunes")}
}

func (_c *DataService_Tunes_Call) Run(run func()) *DataService_Tunes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DataService_Tunes_Call) Return(_a0 []*apimodel.Tune, _a1 error) *DataService_Tunes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_Tunes_Call) RunAndReturn(run func() ([]*apimodel.Tune, error)) *DataService_Tunes_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMusicSet provides a mock function with given fields: id, tune
func (_m *DataService) UpdateMusicSet(id uuid.UUID, tune apimodel.UpdateSet) (*apimodel.MusicSet, error) {
	ret := _m.Called(id, tune)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMusicSet")
	}

	var r0 *apimodel.MusicSet
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, apimodel.UpdateSet) (*apimodel.MusicSet, error)); ok {
		return rf(id, tune)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, apimodel.UpdateSet) *apimodel.MusicSet); ok {
		r0 = rf(id, tune)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apimodel.MusicSet)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, apimodel.UpdateSet) error); ok {
		r1 = rf(id, tune)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_UpdateMusicSet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMusicSet'
type DataService_UpdateMusicSet_Call struct {
	*mock.Call
}

// UpdateMusicSet is a helper method to define mock.On call
//   - id uuid.UUID
//   - tune apimodel.UpdateSet
func (_e *DataService_Expecter) UpdateMusicSet(id interface{}, tune interface{}) *DataService_UpdateMusicSet_Call {
	return &DataService_UpdateMusicSet_Call{Call: _e.mock.On("UpdateMusicSet", id, tune)}
}

func (_c *DataService_UpdateMusicSet_Call) Run(run func(id uuid.UUID, tune apimodel.UpdateSet)) *DataService_UpdateMusicSet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(apimodel.UpdateSet))
	})
	return _c
}

func (_c *DataService_UpdateMusicSet_Call) Return(_a0 *apimodel.MusicSet, _a1 error) *DataService_UpdateMusicSet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_UpdateMusicSet_Call) RunAndReturn(run func(uuid.UUID, apimodel.UpdateSet) (*apimodel.MusicSet, error)) *DataService_UpdateMusicSet_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateTune provides a mock function with given fields: id, tune
func (_m *DataService) UpdateTune(id uuid.UUID, tune apimodel.UpdateTune) (*apimodel.Tune, error) {
	ret := _m.Called(id, tune)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTune")
	}

	var r0 *apimodel.Tune
	var r1 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, apimodel.UpdateTune) (*apimodel.Tune, error)); ok {
		return rf(id, tune)
	}
	if rf, ok := ret.Get(0).(func(uuid.UUID, apimodel.UpdateTune) *apimodel.Tune); ok {
		r0 = rf(id, tune)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apimodel.Tune)
		}
	}

	if rf, ok := ret.Get(1).(func(uuid.UUID, apimodel.UpdateTune) error); ok {
		r1 = rf(id, tune)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DataService_UpdateTune_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateTune'
type DataService_UpdateTune_Call struct {
	*mock.Call
}

// UpdateTune is a helper method to define mock.On call
//   - id uuid.UUID
//   - tune apimodel.UpdateTune
func (_e *DataService_Expecter) UpdateTune(id interface{}, tune interface{}) *DataService_UpdateTune_Call {
	return &DataService_UpdateTune_Call{Call: _e.mock.On("UpdateTune", id, tune)}
}

func (_c *DataService_UpdateTune_Call) Run(run func(id uuid.UUID, tune apimodel.UpdateTune)) *DataService_UpdateTune_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uuid.UUID), args[1].(apimodel.UpdateTune))
	})
	return _c
}

func (_c *DataService_UpdateTune_Call) Return(_a0 *apimodel.Tune, _a1 error) *DataService_UpdateTune_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DataService_UpdateTune_Call) RunAndReturn(run func(uuid.UUID, apimodel.UpdateTune) (*apimodel.Tune, error)) *DataService_UpdateTune_Call {
	_c.Call.Return(run)
	return _c
}

// NewDataService creates a new instance of DataService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataService(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataService {
	mock := &DataService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
