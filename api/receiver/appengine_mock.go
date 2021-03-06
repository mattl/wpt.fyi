// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/web-platform-tests/wpt.fyi/api/receiver (interfaces: AppEngineAPI)

// Package receiver is a generated GoMock package.
package receiver

import (
	context "context"
	io "io"
	http "net/http"
	url "net/url"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/github"
	shared "github.com/web-platform-tests/wpt.fyi/shared"
	taskqueue "google.golang.org/appengine/taskqueue"
)

// MockAppEngineAPI is a mock of AppEngineAPI interface
type MockAppEngineAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAppEngineAPIMockRecorder
}

// MockAppEngineAPIMockRecorder is the mock recorder for MockAppEngineAPI
type MockAppEngineAPIMockRecorder struct {
	mock *MockAppEngineAPI
}

// NewMockAppEngineAPI creates a new mock instance
func NewMockAppEngineAPI(ctrl *gomock.Controller) *MockAppEngineAPI {
	mock := &MockAppEngineAPI{ctrl: ctrl}
	mock.recorder = &MockAppEngineAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppEngineAPI) EXPECT() *MockAppEngineAPIMockRecorder {
	return m.recorder
}

// Context mocks base method
func (m *MockAppEngineAPI) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockAppEngineAPIMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockAppEngineAPI)(nil).Context))
}

// GetGitHubClient mocks base method
func (m *MockAppEngineAPI) GetGitHubClient() (*github.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitHubClient")
	ret0, _ := ret[0].(*github.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitHubClient indicates an expected call of GetGitHubClient
func (mr *MockAppEngineAPIMockRecorder) GetGitHubClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitHubClient", reflect.TypeOf((*MockAppEngineAPI)(nil).GetGitHubClient))
}

// GetHTTPClient mocks base method
func (m *MockAppEngineAPI) GetHTTPClient() *http.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHTTPClient")
	ret0, _ := ret[0].(*http.Client)
	return ret0
}

// GetHTTPClient indicates an expected call of GetHTTPClient
func (mr *MockAppEngineAPIMockRecorder) GetHTTPClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHTTPClient", reflect.TypeOf((*MockAppEngineAPI)(nil).GetHTTPClient))
}

// GetHostname mocks base method
func (m *MockAppEngineAPI) GetHostname() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHostname")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetHostname indicates an expected call of GetHostname
func (mr *MockAppEngineAPIMockRecorder) GetHostname() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHostname", reflect.TypeOf((*MockAppEngineAPI)(nil).GetHostname))
}

// GetResultsURL mocks base method
func (m *MockAppEngineAPI) GetResultsURL(arg0 shared.TestRunFilter) *url.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResultsURL", arg0)
	ret0, _ := ret[0].(*url.URL)
	return ret0
}

// GetResultsURL indicates an expected call of GetResultsURL
func (mr *MockAppEngineAPIMockRecorder) GetResultsURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResultsURL", reflect.TypeOf((*MockAppEngineAPI)(nil).GetResultsURL), arg0)
}

// GetResultsUploadURL mocks base method
func (m *MockAppEngineAPI) GetResultsUploadURL() *url.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResultsUploadURL")
	ret0, _ := ret[0].(*url.URL)
	return ret0
}

// GetResultsUploadURL indicates an expected call of GetResultsUploadURL
func (mr *MockAppEngineAPIMockRecorder) GetResultsUploadURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResultsUploadURL", reflect.TypeOf((*MockAppEngineAPI)(nil).GetResultsUploadURL))
}

// GetRunsURL mocks base method
func (m *MockAppEngineAPI) GetRunsURL(arg0 shared.TestRunFilter) *url.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRunsURL", arg0)
	ret0, _ := ret[0].(*url.URL)
	return ret0
}

// GetRunsURL indicates an expected call of GetRunsURL
func (mr *MockAppEngineAPIMockRecorder) GetRunsURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRunsURL", reflect.TypeOf((*MockAppEngineAPI)(nil).GetRunsURL), arg0)
}

// GetServiceHostname mocks base method
func (m *MockAppEngineAPI) GetServiceHostname(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceHostname", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetServiceHostname indicates an expected call of GetServiceHostname
func (mr *MockAppEngineAPIMockRecorder) GetServiceHostname(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceHostname", reflect.TypeOf((*MockAppEngineAPI)(nil).GetServiceHostname), arg0)
}

// GetSlowHTTPClient mocks base method
func (m *MockAppEngineAPI) GetSlowHTTPClient(arg0 time.Duration) (*http.Client, context.CancelFunc) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlowHTTPClient", arg0)
	ret0, _ := ret[0].(*http.Client)
	ret1, _ := ret[1].(context.CancelFunc)
	return ret0, ret1
}

// GetSlowHTTPClient indicates an expected call of GetSlowHTTPClient
func (mr *MockAppEngineAPIMockRecorder) GetSlowHTTPClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlowHTTPClient", reflect.TypeOf((*MockAppEngineAPI)(nil).GetSlowHTTPClient), arg0)
}

// GetUploader mocks base method
func (m *MockAppEngineAPI) GetUploader(arg0 string) (shared.Uploader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUploader", arg0)
	ret0, _ := ret[0].(shared.Uploader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUploader indicates an expected call of GetUploader
func (mr *MockAppEngineAPIMockRecorder) GetUploader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUploader", reflect.TypeOf((*MockAppEngineAPI)(nil).GetUploader), arg0)
}

// GetVersion mocks base method
func (m *MockAppEngineAPI) GetVersion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetVersion indicates an expected call of GetVersion
func (mr *MockAppEngineAPIMockRecorder) GetVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockAppEngineAPI)(nil).GetVersion))
}

// GetVersionedHostname mocks base method
func (m *MockAppEngineAPI) GetVersionedHostname() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersionedHostname")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetVersionedHostname indicates an expected call of GetVersionedHostname
func (mr *MockAppEngineAPIMockRecorder) GetVersionedHostname() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersionedHostname", reflect.TypeOf((*MockAppEngineAPI)(nil).GetVersionedHostname))
}

// IsAdmin mocks base method
func (m *MockAppEngineAPI) IsAdmin() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAdmin")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAdmin indicates an expected call of IsAdmin
func (mr *MockAppEngineAPIMockRecorder) IsAdmin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAdmin", reflect.TypeOf((*MockAppEngineAPI)(nil).IsAdmin))
}

// IsFeatureEnabled mocks base method
func (m *MockAppEngineAPI) IsFeatureEnabled(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFeatureEnabled", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsFeatureEnabled indicates an expected call of IsFeatureEnabled
func (mr *MockAppEngineAPIMockRecorder) IsFeatureEnabled(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFeatureEnabled", reflect.TypeOf((*MockAppEngineAPI)(nil).IsFeatureEnabled), arg0)
}

// IsLoggedIn mocks base method
func (m *MockAppEngineAPI) IsLoggedIn() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLoggedIn")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsLoggedIn indicates an expected call of IsLoggedIn
func (mr *MockAppEngineAPIMockRecorder) IsLoggedIn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLoggedIn", reflect.TypeOf((*MockAppEngineAPI)(nil).IsLoggedIn))
}

// LoginURL mocks base method
func (m *MockAppEngineAPI) LoginURL(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginURL", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginURL indicates an expected call of LoginURL
func (mr *MockAppEngineAPIMockRecorder) LoginURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginURL", reflect.TypeOf((*MockAppEngineAPI)(nil).LoginURL), arg0)
}

// addTestRun mocks base method
func (m *MockAppEngineAPI) addTestRun(arg0 *shared.TestRun) (*DatastoreKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "addTestRun", arg0)
	ret0, _ := ret[0].(*DatastoreKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// addTestRun indicates an expected call of addTestRun
func (mr *MockAppEngineAPIMockRecorder) addTestRun(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "addTestRun", reflect.TypeOf((*MockAppEngineAPI)(nil).addTestRun), arg0)
}

// authenticateUploader mocks base method
func (m *MockAppEngineAPI) authenticateUploader(arg0, arg1 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "authenticateUploader", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// authenticateUploader indicates an expected call of authenticateUploader
func (mr *MockAppEngineAPIMockRecorder) authenticateUploader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "authenticateUploader", reflect.TypeOf((*MockAppEngineAPI)(nil).authenticateUploader), arg0, arg1)
}

// fetchWithTimeout mocks base method
func (m *MockAppEngineAPI) fetchWithTimeout(arg0 string, arg1 time.Duration) (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "fetchWithTimeout", arg0, arg1)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// fetchWithTimeout indicates an expected call of fetchWithTimeout
func (mr *MockAppEngineAPIMockRecorder) fetchWithTimeout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "fetchWithTimeout", reflect.TypeOf((*MockAppEngineAPI)(nil).fetchWithTimeout), arg0, arg1)
}

// scheduleResultsTask mocks base method
func (m *MockAppEngineAPI) scheduleResultsTask(arg0 string, arg1 []string, arg2 string, arg3 map[string]string) (*taskqueue.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "scheduleResultsTask", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*taskqueue.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// scheduleResultsTask indicates an expected call of scheduleResultsTask
func (mr *MockAppEngineAPIMockRecorder) scheduleResultsTask(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "scheduleResultsTask", reflect.TypeOf((*MockAppEngineAPI)(nil).scheduleResultsTask), arg0, arg1, arg2, arg3)
}

// uploadToGCS mocks base method
func (m *MockAppEngineAPI) uploadToGCS(arg0 string, arg1 io.Reader, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "uploadToGCS", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// uploadToGCS indicates an expected call of uploadToGCS
func (mr *MockAppEngineAPIMockRecorder) uploadToGCS(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "uploadToGCS", reflect.TypeOf((*MockAppEngineAPI)(nil).uploadToGCS), arg0, arg1, arg2)
}
