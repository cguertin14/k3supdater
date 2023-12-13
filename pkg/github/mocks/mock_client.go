// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cguertin14/k3supdater/pkg/github (interfaces: Client)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	legacy "github.com/cguertin14/k3supdater/pkg/github"
	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/v57/github"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CreateBranch mocks base method.
func (m *MockClient) CreateBranch(arg0 context.Context, arg1 legacy.CreateBranchRequest) (*github.Reference, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBranch", arg0, arg1)
	ret0, _ := ret[0].(*github.Reference)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateBranch indicates an expected call of CreateBranch.
func (mr *MockClientMockRecorder) CreateBranch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBranch", reflect.TypeOf((*MockClient)(nil).CreateBranch), arg0, arg1)
}

// CreatePullRequest mocks base method.
func (m *MockClient) CreatePullRequest(arg0 context.Context, arg1 legacy.CreatePRRequest) (*github.PullRequest, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePullRequest", arg0, arg1)
	ret0, _ := ret[0].(*github.PullRequest)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreatePullRequest indicates an expected call of CreatePullRequest.
func (mr *MockClientMockRecorder) CreatePullRequest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePullRequest", reflect.TypeOf((*MockClient)(nil).CreatePullRequest), arg0, arg1)
}

// GetBranch mocks base method.
func (m *MockClient) GetBranch(arg0 context.Context, arg1 legacy.GetBranchRequest) (*github.Reference, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranch", arg0, arg1)
	ret0, _ := ret[0].(*github.Reference)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBranch indicates an expected call of GetBranch.
func (mr *MockClientMockRecorder) GetBranch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranch", reflect.TypeOf((*MockClient)(nil).GetBranch), arg0, arg1)
}

// GetRepositoryContents mocks base method.
func (m *MockClient) GetRepositoryContents(arg0 context.Context, arg1 legacy.GetRepositoryContentsRequest) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepositoryContents", arg0, arg1)
	ret0, _ := ret[0].(*github.RepositoryContent)
	ret1, _ := ret[1].([]*github.RepositoryContent)
	ret2, _ := ret[2].(*github.Response)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetRepositoryContents indicates an expected call of GetRepositoryContents.
func (mr *MockClientMockRecorder) GetRepositoryContents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepositoryContents", reflect.TypeOf((*MockClient)(nil).GetRepositoryContents), arg0, arg1)
}

// GetRepositoryReleases mocks base method.
func (m *MockClient) GetRepositoryReleases(arg0 context.Context, arg1 legacy.CommonRequest) ([]*github.RepositoryRelease, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepositoryReleases", arg0, arg1)
	ret0, _ := ret[0].([]*github.RepositoryRelease)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetRepositoryReleases indicates an expected call of GetRepositoryReleases.
func (mr *MockClientMockRecorder) GetRepositoryReleases(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepositoryReleases", reflect.TypeOf((*MockClient)(nil).GetRepositoryReleases), arg0, arg1)
}

// UpdateFile mocks base method.
func (m *MockClient) UpdateFile(arg0 context.Context, arg1 legacy.UpdateFileRequest) (*github.RepositoryContentResponse, *github.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFile", arg0, arg1)
	ret0, _ := ret[0].(*github.RepositoryContentResponse)
	ret1, _ := ret[1].(*github.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateFile indicates an expected call of UpdateFile.
func (mr *MockClientMockRecorder) UpdateFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFile", reflect.TypeOf((*MockClient)(nil).UpdateFile), arg0, arg1)
}
