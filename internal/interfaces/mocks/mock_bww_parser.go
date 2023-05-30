// Code generated by MockGen. DO NOT EDIT.
// Source: bww_parser.go

// Package mock_interfaces is a generated GoMock package.
package mock_interfaces

import (
	music_model "banduslib/internal/common/music_model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBwwParser is a mock of BwwParser interface.
type MockBwwParser struct {
	ctrl     *gomock.Controller
	recorder *MockBwwParserMockRecorder
}

// MockBwwParserMockRecorder is the mock recorder for MockBwwParser.
type MockBwwParserMockRecorder struct {
	mock *MockBwwParser
}

// NewMockBwwParser creates a new mock instance.
func NewMockBwwParser(ctrl *gomock.Controller) *MockBwwParser {
	mock := &MockBwwParser{ctrl: ctrl}
	mock.recorder = &MockBwwParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBwwParser) EXPECT() *MockBwwParserMockRecorder {
	return m.recorder
}

// ParseBwwData mocks base method.
func (m *MockBwwParser) ParseBwwData(data []byte) (music_model.MusicModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseBwwData", data)
	ret0, _ := ret[0].(music_model.MusicModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseBwwData indicates an expected call of ParseBwwData.
func (mr *MockBwwParserMockRecorder) ParseBwwData(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseBwwData", reflect.TypeOf((*MockBwwParser)(nil).ParseBwwData), data)
}
