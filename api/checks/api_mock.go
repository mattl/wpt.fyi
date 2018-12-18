// Code generated by MockGen. DO NOT EDIT.
// Source: api/checks/api.go

// Package checks is a generated GoMock package.
package checks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/github"
	shared "github.com/web-platform-tests/wpt.fyi/shared"
	reflect "reflect"
)

// MockAPI is a mock of API interface
type MockAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAPIMockRecorder
}

// MockAPIMockRecorder is the mock recorder for MockAPI
type MockAPIMockRecorder struct {
	mock *MockAPI
}

// NewMockAPI creates a new mock instance
func NewMockAPI(ctrl *gomock.Controller) *MockAPI {
	mock := &MockAPI{ctrl: ctrl}
	mock.recorder = &MockAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAPI) EXPECT() *MockAPIMockRecorder {
	return m.recorder
}

// Context mocks base method
func (m *MockAPI) Context() context.Context {
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockAPIMockRecorder) Context() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockAPI)(nil).Context))
}

// ScheduleResultsProcessing mocks base method
func (m *MockAPI) ScheduleResultsProcessing(sha string, browser shared.ProductSpec) error {
	ret := m.ctrl.Call(m, "ScheduleResultsProcessing", sha, browser)
	ret0, _ := ret[0].(error)
	return ret0
}

// ScheduleResultsProcessing indicates an expected call of ScheduleResultsProcessing
func (mr *MockAPIMockRecorder) ScheduleResultsProcessing(sha, browser interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduleResultsProcessing", reflect.TypeOf((*MockAPI)(nil).ScheduleResultsProcessing), sha, browser)
}

// PendingCheckRun mocks base method
func (m *MockAPI) PendingCheckRun(checkSuite shared.CheckSuite, browser shared.ProductSpec) (bool, error) {
	ret := m.ctrl.Call(m, "PendingCheckRun", checkSuite, browser)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingCheckRun indicates an expected call of PendingCheckRun
func (mr *MockAPIMockRecorder) PendingCheckRun(checkSuite, browser interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingCheckRun", reflect.TypeOf((*MockAPI)(nil).PendingCheckRun), checkSuite, browser)
}

// GetSuitesForSHA mocks base method
func (m *MockAPI) GetSuitesForSHA(sha string) ([]shared.CheckSuite, error) {
	ret := m.ctrl.Call(m, "GetSuitesForSHA", sha)
	ret0, _ := ret[0].([]shared.CheckSuite)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSuitesForSHA indicates an expected call of GetSuitesForSHA
func (mr *MockAPIMockRecorder) GetSuitesForSHA(sha interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSuitesForSHA", reflect.TypeOf((*MockAPI)(nil).GetSuitesForSHA), sha)
}

// IgnoreFailure mocks base method
func (m *MockAPI) IgnoreFailure(sender, owner, repo string, run *github.CheckRun, installation *github.Installation) error {
	ret := m.ctrl.Call(m, "IgnoreFailure", sender, owner, repo, run, installation)
	ret0, _ := ret[0].(error)
	return ret0
}

// IgnoreFailure indicates an expected call of IgnoreFailure
func (mr *MockAPIMockRecorder) IgnoreFailure(sender, owner, repo, run, installation interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IgnoreFailure", reflect.TypeOf((*MockAPI)(nil).IgnoreFailure), sender, owner, repo, run, installation)
}

// CancelRun mocks base method
func (m *MockAPI) CancelRun(sender, owner, repo string, run *github.CheckRun, installation *github.Installation) error {
	ret := m.ctrl.Call(m, "CancelRun", sender, owner, repo, run, installation)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelRun indicates an expected call of CancelRun
func (mr *MockAPIMockRecorder) CancelRun(sender, owner, repo, run, installation interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelRun", reflect.TypeOf((*MockAPI)(nil).CancelRun), sender, owner, repo, run, installation)
}

// CreateWPTCheckSuite mocks base method
func (m *MockAPI) CreateWPTCheckSuite(appID, installationID int64, sha string, prNumbers ...int) (bool, error) {
	varargs := []interface{}{appID, installationID, sha}
	for _, a := range prNumbers {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateWPTCheckSuite", varargs...)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWPTCheckSuite indicates an expected call of CreateWPTCheckSuite
func (mr *MockAPIMockRecorder) CreateWPTCheckSuite(appID, installationID, sha interface{}, prNumbers ...interface{}) *gomock.Call {
	varargs := append([]interface{}{appID, installationID, sha}, prNumbers...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWPTCheckSuite", reflect.TypeOf((*MockAPI)(nil).CreateWPTCheckSuite), varargs...)
}
