// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/4chain-ag/go-wallet-toolbox/pkg/wdk (interfaces: WalletStorageWriter)
//
// Generated by this command:
//
//	mockgen -destination=../internal/mocks/mock_wallet_storage_writer.go -package=mocks github.com/4chain-ag/go-wallet-toolbox/pkg/wdk WalletStorageWriter
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	wdk "github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	gomock "go.uber.org/mock/gomock"
)

// MockWalletStorageWriter is a mock of WalletStorageWriter interface.
type MockWalletStorageWriter struct {
	ctrl     *gomock.Controller
	recorder *MockWalletStorageWriterMockRecorder
	isgomock struct{}
}

// MockWalletStorageWriterMockRecorder is the mock recorder for MockWalletStorageWriter.
type MockWalletStorageWriterMockRecorder struct {
	mock *MockWalletStorageWriter
}

// NewMockWalletStorageWriter creates a new mock instance.
func NewMockWalletStorageWriter(ctrl *gomock.Controller) *MockWalletStorageWriter {
	mock := &MockWalletStorageWriter{ctrl: ctrl}
	mock.recorder = &MockWalletStorageWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletStorageWriter) EXPECT() *MockWalletStorageWriterMockRecorder {
	return m.recorder
}

// CreateAction mocks base method.
func (m *MockWalletStorageWriter) CreateAction(auth wdk.AuthID, args wdk.ValidCreateActionArgs) (*wdk.StorageCreateActionResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAction", auth, args)
	ret0, _ := ret[0].(*wdk.StorageCreateActionResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAction indicates an expected call of CreateAction.
func (mr *MockWalletStorageWriterMockRecorder) CreateAction(auth, args any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAction", reflect.TypeOf((*MockWalletStorageWriter)(nil).CreateAction), auth, args)
}

// FindOrInsertUser mocks base method.
func (m *MockWalletStorageWriter) FindOrInsertUser(identityKey string) (*wdk.FindOrInsertUserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrInsertUser", identityKey)
	ret0, _ := ret[0].(*wdk.FindOrInsertUserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrInsertUser indicates an expected call of FindOrInsertUser.
func (mr *MockWalletStorageWriterMockRecorder) FindOrInsertUser(identityKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrInsertUser", reflect.TypeOf((*MockWalletStorageWriter)(nil).FindOrInsertUser), identityKey)
}

// InsertCertificateAuth mocks base method.
func (m *MockWalletStorageWriter) InsertCertificateAuth(auth wdk.AuthID, certificate *wdk.TableCertificateX) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertCertificateAuth", auth, certificate)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertCertificateAuth indicates an expected call of InsertCertificateAuth.
func (mr *MockWalletStorageWriterMockRecorder) InsertCertificateAuth(auth, certificate any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertCertificateAuth", reflect.TypeOf((*MockWalletStorageWriter)(nil).InsertCertificateAuth), auth, certificate)
}

// ListCertificates mocks base method.
func (m *MockWalletStorageWriter) ListCertificates(auth wdk.AuthID, args wdk.ListCertificatesArgs) (*wdk.ListCertificatesResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCertificates", auth, args)
	ret0, _ := ret[0].(*wdk.ListCertificatesResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCertificates indicates an expected call of ListCertificates.
func (mr *MockWalletStorageWriterMockRecorder) ListCertificates(auth, args any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCertificates", reflect.TypeOf((*MockWalletStorageWriter)(nil).ListCertificates), auth, args)
}

// MakeAvailable mocks base method.
func (m *MockWalletStorageWriter) MakeAvailable() (*wdk.TableSettings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeAvailable")
	ret0, _ := ret[0].(*wdk.TableSettings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeAvailable indicates an expected call of MakeAvailable.
func (mr *MockWalletStorageWriterMockRecorder) MakeAvailable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeAvailable", reflect.TypeOf((*MockWalletStorageWriter)(nil).MakeAvailable))
}

// Migrate mocks base method.
func (m *MockWalletStorageWriter) Migrate(storageName, storageIdentityKey string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Migrate", storageName, storageIdentityKey)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Migrate indicates an expected call of Migrate.
func (mr *MockWalletStorageWriterMockRecorder) Migrate(storageName, storageIdentityKey any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockWalletStorageWriter)(nil).Migrate), storageName, storageIdentityKey)
}

// RelinquishCertificate mocks base method.
func (m *MockWalletStorageWriter) RelinquishCertificate(auth wdk.AuthID, args wdk.RelinquishCertificateArgs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RelinquishCertificate", auth, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// RelinquishCertificate indicates an expected call of RelinquishCertificate.
func (mr *MockWalletStorageWriterMockRecorder) RelinquishCertificate(auth, args any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RelinquishCertificate", reflect.TypeOf((*MockWalletStorageWriter)(nil).RelinquishCertificate), auth, args)
}
