// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/modelgeneration (interfaces: APIFacade,State,Model,Generation,Application)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	modelgeneration "github.com/juju/juju/apiserver/facades/client/modelgeneration"
	params "github.com/juju/juju/apiserver/params"
	model "github.com/juju/juju/core/model"
	charm_v6 "gopkg.in/juju/charm.v6"
	names_v2 "gopkg.in/juju/names.v2"
	reflect "reflect"
)

// MockAPIFacade is a mock of APIFacade interface
type MockAPIFacade struct {
	ctrl     *gomock.Controller
	recorder *MockAPIFacadeMockRecorder
}

// MockAPIFacadeMockRecorder is the mock recorder for MockAPIFacade
type MockAPIFacadeMockRecorder struct {
	mock *MockAPIFacade
}

// NewMockAPIFacade creates a new mock instance
func NewMockAPIFacade(ctrl *gomock.Controller) *MockAPIFacade {
	mock := &MockAPIFacade{ctrl: ctrl}
	mock.recorder = &MockAPIFacadeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAPIFacade) EXPECT() *MockAPIFacadeMockRecorder {
	return m.recorder
}

// AddGeneration mocks base method
func (m *MockAPIFacade) AddGeneration(arg0 params.Entity) (params.ErrorResult, error) {
	ret := m.ctrl.Call(m, "AddGeneration", arg0)
	ret0, _ := ret[0].(params.ErrorResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddGeneration indicates an expected call of AddGeneration
func (mr *MockAPIFacadeMockRecorder) AddGeneration(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddGeneration", reflect.TypeOf((*MockAPIFacade)(nil).AddGeneration), arg0)
}

// AdvanceGeneration mocks base method
func (m *MockAPIFacade) AdvanceGeneration(arg0 params.AdvanceGenerationArg) (params.AdvanceGenerationResult, error) {
	ret := m.ctrl.Call(m, "AdvanceGeneration", arg0)
	ret0, _ := ret[0].(params.AdvanceGenerationResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AdvanceGeneration indicates an expected call of AdvanceGeneration
func (mr *MockAPIFacadeMockRecorder) AdvanceGeneration(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdvanceGeneration", reflect.TypeOf((*MockAPIFacade)(nil).AdvanceGeneration), arg0)
}

// CancelGeneration mocks base method
func (m *MockAPIFacade) CancelGeneration(arg0 params.Entity) (params.ErrorResult, error) {
	ret := m.ctrl.Call(m, "CancelGeneration", arg0)
	ret0, _ := ret[0].(params.ErrorResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CancelGeneration indicates an expected call of CancelGeneration
func (mr *MockAPIFacadeMockRecorder) CancelGeneration(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelGeneration", reflect.TypeOf((*MockAPIFacade)(nil).CancelGeneration), arg0)
}

// GenerationInfo mocks base method
func (m *MockAPIFacade) GenerationInfo(arg0 params.Entity) (params.GenerationResult, error) {
	ret := m.ctrl.Call(m, "GenerationInfo", arg0)
	ret0, _ := ret[0].(params.GenerationResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerationInfo indicates an expected call of GenerationInfo
func (mr *MockAPIFacadeMockRecorder) GenerationInfo(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerationInfo", reflect.TypeOf((*MockAPIFacade)(nil).GenerationInfo), arg0)
}

// HasNextGeneration mocks base method
func (m *MockAPIFacade) HasNextGeneration(arg0 params.Entity) (params.BoolResult, error) {
	ret := m.ctrl.Call(m, "HasNextGeneration", arg0)
	ret0, _ := ret[0].(params.BoolResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasNextGeneration indicates an expected call of HasNextGeneration
func (mr *MockAPIFacadeMockRecorder) HasNextGeneration(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasNextGeneration", reflect.TypeOf((*MockAPIFacade)(nil).HasNextGeneration), arg0)
}

// SwitchGeneration mocks base method
func (m *MockAPIFacade) SwitchGeneration(arg0 params.GenerationVersionArg) (params.ErrorResult, error) {
	ret := m.ctrl.Call(m, "SwitchGeneration", arg0)
	ret0, _ := ret[0].(params.ErrorResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SwitchGeneration indicates an expected call of SwitchGeneration
func (mr *MockAPIFacadeMockRecorder) SwitchGeneration(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SwitchGeneration", reflect.TypeOf((*MockAPIFacade)(nil).SwitchGeneration), arg0)
}

// MockState is a mock of State interface
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// Application mocks base method
func (m *MockState) Application(arg0 string) (modelgeneration.Application, error) {
	ret := m.ctrl.Call(m, "Application", arg0)
	ret0, _ := ret[0].(modelgeneration.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Application indicates an expected call of Application
func (mr *MockStateMockRecorder) Application(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockState)(nil).Application), arg0)
}

// ControllerTag mocks base method
func (m *MockState) ControllerTag() names_v2.ControllerTag {
	ret := m.ctrl.Call(m, "ControllerTag")
	ret0, _ := ret[0].(names_v2.ControllerTag)
	return ret0
}

// ControllerTag indicates an expected call of ControllerTag
func (mr *MockStateMockRecorder) ControllerTag() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerTag", reflect.TypeOf((*MockState)(nil).ControllerTag))
}

// Model mocks base method
func (m *MockState) Model() (modelgeneration.Model, error) {
	ret := m.ctrl.Call(m, "Model")
	ret0, _ := ret[0].(modelgeneration.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Model indicates an expected call of Model
func (mr *MockStateMockRecorder) Model() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Model", reflect.TypeOf((*MockState)(nil).Model))
}

// MockModel is a mock of Model interface
type MockModel struct {
	ctrl     *gomock.Controller
	recorder *MockModelMockRecorder
}

// MockModelMockRecorder is the mock recorder for MockModel
type MockModelMockRecorder struct {
	mock *MockModel
}

// NewMockModel creates a new mock instance
func NewMockModel(ctrl *gomock.Controller) *MockModel {
	mock := &MockModel{ctrl: ctrl}
	mock.recorder = &MockModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockModel) EXPECT() *MockModelMockRecorder {
	return m.recorder
}

// AddGeneration mocks base method
func (m *MockModel) AddGeneration() error {
	ret := m.ctrl.Call(m, "AddGeneration")
	ret0, _ := ret[0].(error)
	return ret0
}

// AddGeneration indicates an expected call of AddGeneration
func (mr *MockModelMockRecorder) AddGeneration() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddGeneration", reflect.TypeOf((*MockModel)(nil).AddGeneration))
}

// HasNextGeneration mocks base method
func (m *MockModel) HasNextGeneration() (bool, error) {
	ret := m.ctrl.Call(m, "HasNextGeneration")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HasNextGeneration indicates an expected call of HasNextGeneration
func (mr *MockModelMockRecorder) HasNextGeneration() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasNextGeneration", reflect.TypeOf((*MockModel)(nil).HasNextGeneration))
}

// NextGeneration mocks base method
func (m *MockModel) NextGeneration() (modelgeneration.Generation, error) {
	ret := m.ctrl.Call(m, "NextGeneration")
	ret0, _ := ret[0].(modelgeneration.Generation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NextGeneration indicates an expected call of NextGeneration
func (mr *MockModelMockRecorder) NextGeneration() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextGeneration", reflect.TypeOf((*MockModel)(nil).NextGeneration))
}

// MockGeneration is a mock of Generation interface
type MockGeneration struct {
	ctrl     *gomock.Controller
	recorder *MockGenerationMockRecorder
}

// MockGenerationMockRecorder is the mock recorder for MockGeneration
type MockGenerationMockRecorder struct {
	mock *MockGeneration
}

// NewMockGeneration creates a new mock instance
func NewMockGeneration(ctrl *gomock.Controller) *MockGeneration {
	mock := &MockGeneration{ctrl: ctrl}
	mock.recorder = &MockGenerationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGeneration) EXPECT() *MockGenerationMockRecorder {
	return m.recorder
}

// AssignAllUnits mocks base method
func (m *MockGeneration) AssignAllUnits(arg0 string) error {
	ret := m.ctrl.Call(m, "AssignAllUnits", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignAllUnits indicates an expected call of AssignAllUnits
func (mr *MockGenerationMockRecorder) AssignAllUnits(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignAllUnits", reflect.TypeOf((*MockGeneration)(nil).AssignAllUnits), arg0)
}

// AssignUnit mocks base method
func (m *MockGeneration) AssignUnit(arg0 string) error {
	ret := m.ctrl.Call(m, "AssignUnit", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignUnit indicates an expected call of AssignUnit
func (mr *MockGenerationMockRecorder) AssignUnit(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignUnit", reflect.TypeOf((*MockGeneration)(nil).AssignUnit), arg0)
}

// AssignedUnits mocks base method
func (m *MockGeneration) AssignedUnits() map[string][]string {
	ret := m.ctrl.Call(m, "AssignedUnits")
	ret0, _ := ret[0].(map[string][]string)
	return ret0
}

// AssignedUnits indicates an expected call of AssignedUnits
func (mr *MockGenerationMockRecorder) AssignedUnits() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignedUnits", reflect.TypeOf((*MockGeneration)(nil).AssignedUnits))
}

// AutoComplete mocks base method
func (m *MockGeneration) AutoComplete() (bool, error) {
	ret := m.ctrl.Call(m, "AutoComplete")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AutoComplete indicates an expected call of AutoComplete
func (mr *MockGenerationMockRecorder) AutoComplete() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AutoComplete", reflect.TypeOf((*MockGeneration)(nil).AutoComplete))
}

// MakeCurrent mocks base method
func (m *MockGeneration) MakeCurrent() error {
	ret := m.ctrl.Call(m, "MakeCurrent")
	ret0, _ := ret[0].(error)
	return ret0
}

// MakeCurrent indicates an expected call of MakeCurrent
func (mr *MockGenerationMockRecorder) MakeCurrent() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeCurrent", reflect.TypeOf((*MockGeneration)(nil).MakeCurrent))
}

// Refresh mocks base method
func (m *MockGeneration) Refresh() error {
	ret := m.ctrl.Call(m, "Refresh")
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh
func (mr *MockGenerationMockRecorder) Refresh() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockGeneration)(nil).Refresh))
}

// MockApplication is a mock of Application interface
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// CharmConfig mocks base method
func (m *MockApplication) CharmConfig(arg0 model.GenerationVersion) (charm_v6.Settings, error) {
	ret := m.ctrl.Call(m, "CharmConfig", arg0)
	ret0, _ := ret[0].(charm_v6.Settings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CharmConfig indicates an expected call of CharmConfig
func (mr *MockApplicationMockRecorder) CharmConfig(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmConfig", reflect.TypeOf((*MockApplication)(nil).CharmConfig), arg0)
}
