// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/yonesko/s3-test-task/internal/filegateway/client/filestorage (interfaces: Client)

// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/yonesko/s3-test-task/internal/model"
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

// GetFile mocks base method.
func (m *MockClient) GetFile(arg0 context.Context, arg1, arg2 string) (io.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFile indicates an expected call of GetFile.
func (mr *MockClientMockRecorder) GetFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockClient)(nil).GetFile), arg0, arg1, arg2)
}

// SaveFile mocks base method.
func (m *MockClient) SaveFile(arg0 context.Context, arg1 string, arg2 model.File) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveFile", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveFile indicates an expected call of SaveFile.
func (mr *MockClientMockRecorder) SaveFile(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveFile", reflect.TypeOf((*MockClient)(nil).SaveFile), arg0, arg1, arg2)
}
