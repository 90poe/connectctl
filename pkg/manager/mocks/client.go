// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"net/http"
	"sync"

	"github.com/90poe/connectctl/pkg/client/connect"
)

type FakeClient struct {
	CreateConnectorStub        func(connect.Connector) (*http.Response, error)
	createConnectorMutex       sync.RWMutex
	createConnectorArgsForCall []struct {
		arg1 connect.Connector
	}
	createConnectorReturns struct {
		result1 *http.Response
		result2 error
	}
	createConnectorReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	DeleteConnectorStub        func(string) (*http.Response, error)
	deleteConnectorMutex       sync.RWMutex
	deleteConnectorArgsForCall []struct {
		arg1 string
	}
	deleteConnectorReturns struct {
		result1 *http.Response
		result2 error
	}
	deleteConnectorReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	GetClusterInfoStub        func() (*connect.ClusterInfo, *http.Response, error)
	getClusterInfoMutex       sync.RWMutex
	getClusterInfoArgsForCall []struct {
	}
	getClusterInfoReturns struct {
		result1 *connect.ClusterInfo
		result2 *http.Response
		result3 error
	}
	getClusterInfoReturnsOnCall map[int]struct {
		result1 *connect.ClusterInfo
		result2 *http.Response
		result3 error
	}
	GetConnectorStub        func(string) (*connect.Connector, *http.Response, error)
	getConnectorMutex       sync.RWMutex
	getConnectorArgsForCall []struct {
		arg1 string
	}
	getConnectorReturns struct {
		result1 *connect.Connector
		result2 *http.Response
		result3 error
	}
	getConnectorReturnsOnCall map[int]struct {
		result1 *connect.Connector
		result2 *http.Response
		result3 error
	}
	GetConnectorStatusStub        func(string) (*connect.ConnectorStatus, *http.Response, error)
	getConnectorStatusMutex       sync.RWMutex
	getConnectorStatusArgsForCall []struct {
		arg1 string
	}
	getConnectorStatusReturns struct {
		result1 *connect.ConnectorStatus
		result2 *http.Response
		result3 error
	}
	getConnectorStatusReturnsOnCall map[int]struct {
		result1 *connect.ConnectorStatus
		result2 *http.Response
		result3 error
	}
	ListConnectorsStub        func() ([]string, *http.Response, error)
	listConnectorsMutex       sync.RWMutex
	listConnectorsArgsForCall []struct {
	}
	listConnectorsReturns struct {
		result1 []string
		result2 *http.Response
		result3 error
	}
	listConnectorsReturnsOnCall map[int]struct {
		result1 []string
		result2 *http.Response
		result3 error
	}
	ListPluginsStub        func() ([]*connect.Plugin, *http.Response, error)
	listPluginsMutex       sync.RWMutex
	listPluginsArgsForCall []struct {
	}
	listPluginsReturns struct {
		result1 []*connect.Plugin
		result2 *http.Response
		result3 error
	}
	listPluginsReturnsOnCall map[int]struct {
		result1 []*connect.Plugin
		result2 *http.Response
		result3 error
	}
	PauseConnectorStub        func(string) (*http.Response, error)
	pauseConnectorMutex       sync.RWMutex
	pauseConnectorArgsForCall []struct {
		arg1 string
	}
	pauseConnectorReturns struct {
		result1 *http.Response
		result2 error
	}
	pauseConnectorReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	RestartConnectorStub        func(string) (*http.Response, error)
	restartConnectorMutex       sync.RWMutex
	restartConnectorArgsForCall []struct {
		arg1 string
	}
	restartConnectorReturns struct {
		result1 *http.Response
		result2 error
	}
	restartConnectorReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	RestartConnectorTaskStub        func(string, int) (*http.Response, error)
	restartConnectorTaskMutex       sync.RWMutex
	restartConnectorTaskArgsForCall []struct {
		arg1 string
		arg2 int
	}
	restartConnectorTaskReturns struct {
		result1 *http.Response
		result2 error
	}
	restartConnectorTaskReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	ResumeConnectorStub        func(string) (*http.Response, error)
	resumeConnectorMutex       sync.RWMutex
	resumeConnectorArgsForCall []struct {
		arg1 string
	}
	resumeConnectorReturns struct {
		result1 *http.Response
		result2 error
	}
	resumeConnectorReturnsOnCall map[int]struct {
		result1 *http.Response
		result2 error
	}
	UpdateConnectorConfigStub        func(string, connect.ConnectorConfig) (*connect.Connector, *http.Response, error)
	updateConnectorConfigMutex       sync.RWMutex
	updateConnectorConfigArgsForCall []struct {
		arg1 string
		arg2 connect.ConnectorConfig
	}
	updateConnectorConfigReturns struct {
		result1 *connect.Connector
		result2 *http.Response
		result3 error
	}
	updateConnectorConfigReturnsOnCall map[int]struct {
		result1 *connect.Connector
		result2 *http.Response
		result3 error
	}
	ValidatePluginsStub        func(connect.ConnectorConfig) (*connect.ConfigValidation, *http.Response, error)
	validatePluginsMutex       sync.RWMutex
	validatePluginsArgsForCall []struct {
		arg1 connect.ConnectorConfig
	}
	validatePluginsReturns struct {
		result1 *connect.ConfigValidation
		result2 *http.Response
		result3 error
	}
	validatePluginsReturnsOnCall map[int]struct {
		result1 *connect.ConfigValidation
		result2 *http.Response
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) CreateConnector(arg1 connect.Connector) (*http.Response, error) {
	fake.createConnectorMutex.Lock()
	ret, specificReturn := fake.createConnectorReturnsOnCall[len(fake.createConnectorArgsForCall)]
	fake.createConnectorArgsForCall = append(fake.createConnectorArgsForCall, struct {
		arg1 connect.Connector
	}{arg1})
	fake.recordInvocation("CreateConnector", []interface{}{arg1})
	fake.createConnectorMutex.Unlock()
	if fake.CreateConnectorStub != nil {
		return fake.CreateConnectorStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.createConnectorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) CreateConnectorCallCount() int {
	fake.createConnectorMutex.RLock()
	defer fake.createConnectorMutex.RUnlock()
	return len(fake.createConnectorArgsForCall)
}

func (fake *FakeClient) CreateConnectorCalls(stub func(connect.Connector) (*http.Response, error)) {
	fake.createConnectorMutex.Lock()
	defer fake.createConnectorMutex.Unlock()
	fake.CreateConnectorStub = stub
}

func (fake *FakeClient) CreateConnectorArgsForCall(i int) connect.Connector {
	fake.createConnectorMutex.RLock()
	defer fake.createConnectorMutex.RUnlock()
	argsForCall := fake.createConnectorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) CreateConnectorReturns(result1 *http.Response, result2 error) {
	fake.createConnectorMutex.Lock()
	defer fake.createConnectorMutex.Unlock()
	fake.CreateConnectorStub = nil
	fake.createConnectorReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) CreateConnectorReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.createConnectorMutex.Lock()
	defer fake.createConnectorMutex.Unlock()
	fake.CreateConnectorStub = nil
	if fake.createConnectorReturnsOnCall == nil {
		fake.createConnectorReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.createConnectorReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DeleteConnector(arg1 string) (*http.Response, error) {
	fake.deleteConnectorMutex.Lock()
	ret, specificReturn := fake.deleteConnectorReturnsOnCall[len(fake.deleteConnectorArgsForCall)]
	fake.deleteConnectorArgsForCall = append(fake.deleteConnectorArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("DeleteConnector", []interface{}{arg1})
	fake.deleteConnectorMutex.Unlock()
	if fake.DeleteConnectorStub != nil {
		return fake.DeleteConnectorStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deleteConnectorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) DeleteConnectorCallCount() int {
	fake.deleteConnectorMutex.RLock()
	defer fake.deleteConnectorMutex.RUnlock()
	return len(fake.deleteConnectorArgsForCall)
}

func (fake *FakeClient) DeleteConnectorCalls(stub func(string) (*http.Response, error)) {
	fake.deleteConnectorMutex.Lock()
	defer fake.deleteConnectorMutex.Unlock()
	fake.DeleteConnectorStub = stub
}

func (fake *FakeClient) DeleteConnectorArgsForCall(i int) string {
	fake.deleteConnectorMutex.RLock()
	defer fake.deleteConnectorMutex.RUnlock()
	argsForCall := fake.deleteConnectorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) DeleteConnectorReturns(result1 *http.Response, result2 error) {
	fake.deleteConnectorMutex.Lock()
	defer fake.deleteConnectorMutex.Unlock()
	fake.DeleteConnectorStub = nil
	fake.deleteConnectorReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DeleteConnectorReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.deleteConnectorMutex.Lock()
	defer fake.deleteConnectorMutex.Unlock()
	fake.DeleteConnectorStub = nil
	if fake.deleteConnectorReturnsOnCall == nil {
		fake.deleteConnectorReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.deleteConnectorReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetClusterInfo() (*connect.ClusterInfo, *http.Response, error) {
	fake.getClusterInfoMutex.Lock()
	ret, specificReturn := fake.getClusterInfoReturnsOnCall[len(fake.getClusterInfoArgsForCall)]
	fake.getClusterInfoArgsForCall = append(fake.getClusterInfoArgsForCall, struct {
	}{})
	fake.recordInvocation("GetClusterInfo", []interface{}{})
	fake.getClusterInfoMutex.Unlock()
	if fake.GetClusterInfoStub != nil {
		return fake.GetClusterInfoStub()
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.getClusterInfoReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) GetClusterInfoCallCount() int {
	fake.getClusterInfoMutex.RLock()
	defer fake.getClusterInfoMutex.RUnlock()
	return len(fake.getClusterInfoArgsForCall)
}

func (fake *FakeClient) GetClusterInfoCalls(stub func() (*connect.ClusterInfo, *http.Response, error)) {
	fake.getClusterInfoMutex.Lock()
	defer fake.getClusterInfoMutex.Unlock()
	fake.GetClusterInfoStub = stub
}

func (fake *FakeClient) GetClusterInfoReturns(result1 *connect.ClusterInfo, result2 *http.Response, result3 error) {
	fake.getClusterInfoMutex.Lock()
	defer fake.getClusterInfoMutex.Unlock()
	fake.GetClusterInfoStub = nil
	fake.getClusterInfoReturns = struct {
		result1 *connect.ClusterInfo
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) GetClusterInfoReturnsOnCall(i int, result1 *connect.ClusterInfo, result2 *http.Response, result3 error) {
	fake.getClusterInfoMutex.Lock()
	defer fake.getClusterInfoMutex.Unlock()
	fake.GetClusterInfoStub = nil
	if fake.getClusterInfoReturnsOnCall == nil {
		fake.getClusterInfoReturnsOnCall = make(map[int]struct {
			result1 *connect.ClusterInfo
			result2 *http.Response
			result3 error
		})
	}
	fake.getClusterInfoReturnsOnCall[i] = struct {
		result1 *connect.ClusterInfo
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) GetConnector(arg1 string) (*connect.Connector, *http.Response, error) {
	fake.getConnectorMutex.Lock()
	ret, specificReturn := fake.getConnectorReturnsOnCall[len(fake.getConnectorArgsForCall)]
	fake.getConnectorArgsForCall = append(fake.getConnectorArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetConnector", []interface{}{arg1})
	fake.getConnectorMutex.Unlock()
	if fake.GetConnectorStub != nil {
		return fake.GetConnectorStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.getConnectorReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) GetConnectorCallCount() int {
	fake.getConnectorMutex.RLock()
	defer fake.getConnectorMutex.RUnlock()
	return len(fake.getConnectorArgsForCall)
}

func (fake *FakeClient) GetConnectorCalls(stub func(string) (*connect.Connector, *http.Response, error)) {
	fake.getConnectorMutex.Lock()
	defer fake.getConnectorMutex.Unlock()
	fake.GetConnectorStub = stub
}

func (fake *FakeClient) GetConnectorArgsForCall(i int) string {
	fake.getConnectorMutex.RLock()
	defer fake.getConnectorMutex.RUnlock()
	argsForCall := fake.getConnectorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) GetConnectorReturns(result1 *connect.Connector, result2 *http.Response, result3 error) {
	fake.getConnectorMutex.Lock()
	defer fake.getConnectorMutex.Unlock()
	fake.GetConnectorStub = nil
	fake.getConnectorReturns = struct {
		result1 *connect.Connector
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) GetConnectorReturnsOnCall(i int, result1 *connect.Connector, result2 *http.Response, result3 error) {
	fake.getConnectorMutex.Lock()
	defer fake.getConnectorMutex.Unlock()
	fake.GetConnectorStub = nil
	if fake.getConnectorReturnsOnCall == nil {
		fake.getConnectorReturnsOnCall = make(map[int]struct {
			result1 *connect.Connector
			result2 *http.Response
			result3 error
		})
	}
	fake.getConnectorReturnsOnCall[i] = struct {
		result1 *connect.Connector
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) GetConnectorStatus(arg1 string) (*connect.ConnectorStatus, *http.Response, error) {
	fake.getConnectorStatusMutex.Lock()
	ret, specificReturn := fake.getConnectorStatusReturnsOnCall[len(fake.getConnectorStatusArgsForCall)]
	fake.getConnectorStatusArgsForCall = append(fake.getConnectorStatusArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetConnectorStatus", []interface{}{arg1})
	fake.getConnectorStatusMutex.Unlock()
	if fake.GetConnectorStatusStub != nil {
		return fake.GetConnectorStatusStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.getConnectorStatusReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) GetConnectorStatusCallCount() int {
	fake.getConnectorStatusMutex.RLock()
	defer fake.getConnectorStatusMutex.RUnlock()
	return len(fake.getConnectorStatusArgsForCall)
}

func (fake *FakeClient) GetConnectorStatusCalls(stub func(string) (*connect.ConnectorStatus, *http.Response, error)) {
	fake.getConnectorStatusMutex.Lock()
	defer fake.getConnectorStatusMutex.Unlock()
	fake.GetConnectorStatusStub = stub
}

func (fake *FakeClient) GetConnectorStatusArgsForCall(i int) string {
	fake.getConnectorStatusMutex.RLock()
	defer fake.getConnectorStatusMutex.RUnlock()
	argsForCall := fake.getConnectorStatusArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) GetConnectorStatusReturns(result1 *connect.ConnectorStatus, result2 *http.Response, result3 error) {
	fake.getConnectorStatusMutex.Lock()
	defer fake.getConnectorStatusMutex.Unlock()
	fake.GetConnectorStatusStub = nil
	fake.getConnectorStatusReturns = struct {
		result1 *connect.ConnectorStatus
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) GetConnectorStatusReturnsOnCall(i int, result1 *connect.ConnectorStatus, result2 *http.Response, result3 error) {
	fake.getConnectorStatusMutex.Lock()
	defer fake.getConnectorStatusMutex.Unlock()
	fake.GetConnectorStatusStub = nil
	if fake.getConnectorStatusReturnsOnCall == nil {
		fake.getConnectorStatusReturnsOnCall = make(map[int]struct {
			result1 *connect.ConnectorStatus
			result2 *http.Response
			result3 error
		})
	}
	fake.getConnectorStatusReturnsOnCall[i] = struct {
		result1 *connect.ConnectorStatus
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) ListConnectors() ([]string, *http.Response, error) {
	fake.listConnectorsMutex.Lock()
	ret, specificReturn := fake.listConnectorsReturnsOnCall[len(fake.listConnectorsArgsForCall)]
	fake.listConnectorsArgsForCall = append(fake.listConnectorsArgsForCall, struct {
	}{})
	fake.recordInvocation("ListConnectors", []interface{}{})
	fake.listConnectorsMutex.Unlock()
	if fake.ListConnectorsStub != nil {
		return fake.ListConnectorsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.listConnectorsReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) ListConnectorsCallCount() int {
	fake.listConnectorsMutex.RLock()
	defer fake.listConnectorsMutex.RUnlock()
	return len(fake.listConnectorsArgsForCall)
}

func (fake *FakeClient) ListConnectorsCalls(stub func() ([]string, *http.Response, error)) {
	fake.listConnectorsMutex.Lock()
	defer fake.listConnectorsMutex.Unlock()
	fake.ListConnectorsStub = stub
}

func (fake *FakeClient) ListConnectorsReturns(result1 []string, result2 *http.Response, result3 error) {
	fake.listConnectorsMutex.Lock()
	defer fake.listConnectorsMutex.Unlock()
	fake.ListConnectorsStub = nil
	fake.listConnectorsReturns = struct {
		result1 []string
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) ListConnectorsReturnsOnCall(i int, result1 []string, result2 *http.Response, result3 error) {
	fake.listConnectorsMutex.Lock()
	defer fake.listConnectorsMutex.Unlock()
	fake.ListConnectorsStub = nil
	if fake.listConnectorsReturnsOnCall == nil {
		fake.listConnectorsReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 *http.Response
			result3 error
		})
	}
	fake.listConnectorsReturnsOnCall[i] = struct {
		result1 []string
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) ListPlugins() ([]*connect.Plugin, *http.Response, error) {
	fake.listPluginsMutex.Lock()
	ret, specificReturn := fake.listPluginsReturnsOnCall[len(fake.listPluginsArgsForCall)]
	fake.listPluginsArgsForCall = append(fake.listPluginsArgsForCall, struct {
	}{})
	fake.recordInvocation("ListPlugins", []interface{}{})
	fake.listPluginsMutex.Unlock()
	if fake.ListPluginsStub != nil {
		return fake.ListPluginsStub()
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.listPluginsReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) ListPluginsCallCount() int {
	fake.listPluginsMutex.RLock()
	defer fake.listPluginsMutex.RUnlock()
	return len(fake.listPluginsArgsForCall)
}

func (fake *FakeClient) ListPluginsCalls(stub func() ([]*connect.Plugin, *http.Response, error)) {
	fake.listPluginsMutex.Lock()
	defer fake.listPluginsMutex.Unlock()
	fake.ListPluginsStub = stub
}

func (fake *FakeClient) ListPluginsReturns(result1 []*connect.Plugin, result2 *http.Response, result3 error) {
	fake.listPluginsMutex.Lock()
	defer fake.listPluginsMutex.Unlock()
	fake.ListPluginsStub = nil
	fake.listPluginsReturns = struct {
		result1 []*connect.Plugin
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) ListPluginsReturnsOnCall(i int, result1 []*connect.Plugin, result2 *http.Response, result3 error) {
	fake.listPluginsMutex.Lock()
	defer fake.listPluginsMutex.Unlock()
	fake.ListPluginsStub = nil
	if fake.listPluginsReturnsOnCall == nil {
		fake.listPluginsReturnsOnCall = make(map[int]struct {
			result1 []*connect.Plugin
			result2 *http.Response
			result3 error
		})
	}
	fake.listPluginsReturnsOnCall[i] = struct {
		result1 []*connect.Plugin
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PauseConnector(arg1 string) (*http.Response, error) {
	fake.pauseConnectorMutex.Lock()
	ret, specificReturn := fake.pauseConnectorReturnsOnCall[len(fake.pauseConnectorArgsForCall)]
	fake.pauseConnectorArgsForCall = append(fake.pauseConnectorArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("PauseConnector", []interface{}{arg1})
	fake.pauseConnectorMutex.Unlock()
	if fake.PauseConnectorStub != nil {
		return fake.PauseConnectorStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.pauseConnectorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) PauseConnectorCallCount() int {
	fake.pauseConnectorMutex.RLock()
	defer fake.pauseConnectorMutex.RUnlock()
	return len(fake.pauseConnectorArgsForCall)
}

func (fake *FakeClient) PauseConnectorCalls(stub func(string) (*http.Response, error)) {
	fake.pauseConnectorMutex.Lock()
	defer fake.pauseConnectorMutex.Unlock()
	fake.PauseConnectorStub = stub
}

func (fake *FakeClient) PauseConnectorArgsForCall(i int) string {
	fake.pauseConnectorMutex.RLock()
	defer fake.pauseConnectorMutex.RUnlock()
	argsForCall := fake.pauseConnectorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) PauseConnectorReturns(result1 *http.Response, result2 error) {
	fake.pauseConnectorMutex.Lock()
	defer fake.pauseConnectorMutex.Unlock()
	fake.PauseConnectorStub = nil
	fake.pauseConnectorReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PauseConnectorReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.pauseConnectorMutex.Lock()
	defer fake.pauseConnectorMutex.Unlock()
	fake.PauseConnectorStub = nil
	if fake.pauseConnectorReturnsOnCall == nil {
		fake.pauseConnectorReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.pauseConnectorReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) RestartConnector(arg1 string) (*http.Response, error) {
	fake.restartConnectorMutex.Lock()
	ret, specificReturn := fake.restartConnectorReturnsOnCall[len(fake.restartConnectorArgsForCall)]
	fake.restartConnectorArgsForCall = append(fake.restartConnectorArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("RestartConnector", []interface{}{arg1})
	fake.restartConnectorMutex.Unlock()
	if fake.RestartConnectorStub != nil {
		return fake.RestartConnectorStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.restartConnectorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) RestartConnectorCallCount() int {
	fake.restartConnectorMutex.RLock()
	defer fake.restartConnectorMutex.RUnlock()
	return len(fake.restartConnectorArgsForCall)
}

func (fake *FakeClient) RestartConnectorCalls(stub func(string) (*http.Response, error)) {
	fake.restartConnectorMutex.Lock()
	defer fake.restartConnectorMutex.Unlock()
	fake.RestartConnectorStub = stub
}

func (fake *FakeClient) RestartConnectorArgsForCall(i int) string {
	fake.restartConnectorMutex.RLock()
	defer fake.restartConnectorMutex.RUnlock()
	argsForCall := fake.restartConnectorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) RestartConnectorReturns(result1 *http.Response, result2 error) {
	fake.restartConnectorMutex.Lock()
	defer fake.restartConnectorMutex.Unlock()
	fake.RestartConnectorStub = nil
	fake.restartConnectorReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) RestartConnectorReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.restartConnectorMutex.Lock()
	defer fake.restartConnectorMutex.Unlock()
	fake.RestartConnectorStub = nil
	if fake.restartConnectorReturnsOnCall == nil {
		fake.restartConnectorReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.restartConnectorReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) RestartConnectorTask(arg1 string, arg2 int) (*http.Response, error) {
	fake.restartConnectorTaskMutex.Lock()
	ret, specificReturn := fake.restartConnectorTaskReturnsOnCall[len(fake.restartConnectorTaskArgsForCall)]
	fake.restartConnectorTaskArgsForCall = append(fake.restartConnectorTaskArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	fake.recordInvocation("RestartConnectorTask", []interface{}{arg1, arg2})
	fake.restartConnectorTaskMutex.Unlock()
	if fake.RestartConnectorTaskStub != nil {
		return fake.RestartConnectorTaskStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.restartConnectorTaskReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) RestartConnectorTaskCallCount() int {
	fake.restartConnectorTaskMutex.RLock()
	defer fake.restartConnectorTaskMutex.RUnlock()
	return len(fake.restartConnectorTaskArgsForCall)
}

func (fake *FakeClient) RestartConnectorTaskCalls(stub func(string, int) (*http.Response, error)) {
	fake.restartConnectorTaskMutex.Lock()
	defer fake.restartConnectorTaskMutex.Unlock()
	fake.RestartConnectorTaskStub = stub
}

func (fake *FakeClient) RestartConnectorTaskArgsForCall(i int) (string, int) {
	fake.restartConnectorTaskMutex.RLock()
	defer fake.restartConnectorTaskMutex.RUnlock()
	argsForCall := fake.restartConnectorTaskArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) RestartConnectorTaskReturns(result1 *http.Response, result2 error) {
	fake.restartConnectorTaskMutex.Lock()
	defer fake.restartConnectorTaskMutex.Unlock()
	fake.RestartConnectorTaskStub = nil
	fake.restartConnectorTaskReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) RestartConnectorTaskReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.restartConnectorTaskMutex.Lock()
	defer fake.restartConnectorTaskMutex.Unlock()
	fake.RestartConnectorTaskStub = nil
	if fake.restartConnectorTaskReturnsOnCall == nil {
		fake.restartConnectorTaskReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.restartConnectorTaskReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ResumeConnector(arg1 string) (*http.Response, error) {
	fake.resumeConnectorMutex.Lock()
	ret, specificReturn := fake.resumeConnectorReturnsOnCall[len(fake.resumeConnectorArgsForCall)]
	fake.resumeConnectorArgsForCall = append(fake.resumeConnectorArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ResumeConnector", []interface{}{arg1})
	fake.resumeConnectorMutex.Unlock()
	if fake.ResumeConnectorStub != nil {
		return fake.ResumeConnectorStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.resumeConnectorReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ResumeConnectorCallCount() int {
	fake.resumeConnectorMutex.RLock()
	defer fake.resumeConnectorMutex.RUnlock()
	return len(fake.resumeConnectorArgsForCall)
}

func (fake *FakeClient) ResumeConnectorCalls(stub func(string) (*http.Response, error)) {
	fake.resumeConnectorMutex.Lock()
	defer fake.resumeConnectorMutex.Unlock()
	fake.ResumeConnectorStub = stub
}

func (fake *FakeClient) ResumeConnectorArgsForCall(i int) string {
	fake.resumeConnectorMutex.RLock()
	defer fake.resumeConnectorMutex.RUnlock()
	argsForCall := fake.resumeConnectorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) ResumeConnectorReturns(result1 *http.Response, result2 error) {
	fake.resumeConnectorMutex.Lock()
	defer fake.resumeConnectorMutex.Unlock()
	fake.ResumeConnectorStub = nil
	fake.resumeConnectorReturns = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ResumeConnectorReturnsOnCall(i int, result1 *http.Response, result2 error) {
	fake.resumeConnectorMutex.Lock()
	defer fake.resumeConnectorMutex.Unlock()
	fake.ResumeConnectorStub = nil
	if fake.resumeConnectorReturnsOnCall == nil {
		fake.resumeConnectorReturnsOnCall = make(map[int]struct {
			result1 *http.Response
			result2 error
		})
	}
	fake.resumeConnectorReturnsOnCall[i] = struct {
		result1 *http.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) UpdateConnectorConfig(arg1 string, arg2 connect.ConnectorConfig) (*connect.Connector, *http.Response, error) {
	fake.updateConnectorConfigMutex.Lock()
	ret, specificReturn := fake.updateConnectorConfigReturnsOnCall[len(fake.updateConnectorConfigArgsForCall)]
	fake.updateConnectorConfigArgsForCall = append(fake.updateConnectorConfigArgsForCall, struct {
		arg1 string
		arg2 connect.ConnectorConfig
	}{arg1, arg2})
	fake.recordInvocation("UpdateConnectorConfig", []interface{}{arg1, arg2})
	fake.updateConnectorConfigMutex.Unlock()
	if fake.UpdateConnectorConfigStub != nil {
		return fake.UpdateConnectorConfigStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.updateConnectorConfigReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) UpdateConnectorConfigCallCount() int {
	fake.updateConnectorConfigMutex.RLock()
	defer fake.updateConnectorConfigMutex.RUnlock()
	return len(fake.updateConnectorConfigArgsForCall)
}

func (fake *FakeClient) UpdateConnectorConfigCalls(stub func(string, connect.ConnectorConfig) (*connect.Connector, *http.Response, error)) {
	fake.updateConnectorConfigMutex.Lock()
	defer fake.updateConnectorConfigMutex.Unlock()
	fake.UpdateConnectorConfigStub = stub
}

func (fake *FakeClient) UpdateConnectorConfigArgsForCall(i int) (string, connect.ConnectorConfig) {
	fake.updateConnectorConfigMutex.RLock()
	defer fake.updateConnectorConfigMutex.RUnlock()
	argsForCall := fake.updateConnectorConfigArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) UpdateConnectorConfigReturns(result1 *connect.Connector, result2 *http.Response, result3 error) {
	fake.updateConnectorConfigMutex.Lock()
	defer fake.updateConnectorConfigMutex.Unlock()
	fake.UpdateConnectorConfigStub = nil
	fake.updateConnectorConfigReturns = struct {
		result1 *connect.Connector
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) UpdateConnectorConfigReturnsOnCall(i int, result1 *connect.Connector, result2 *http.Response, result3 error) {
	fake.updateConnectorConfigMutex.Lock()
	defer fake.updateConnectorConfigMutex.Unlock()
	fake.UpdateConnectorConfigStub = nil
	if fake.updateConnectorConfigReturnsOnCall == nil {
		fake.updateConnectorConfigReturnsOnCall = make(map[int]struct {
			result1 *connect.Connector
			result2 *http.Response
			result3 error
		})
	}
	fake.updateConnectorConfigReturnsOnCall[i] = struct {
		result1 *connect.Connector
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) ValidatePlugins(arg1 connect.ConnectorConfig) (*connect.ConfigValidation, *http.Response, error) {
	fake.validatePluginsMutex.Lock()
	ret, specificReturn := fake.validatePluginsReturnsOnCall[len(fake.validatePluginsArgsForCall)]
	fake.validatePluginsArgsForCall = append(fake.validatePluginsArgsForCall, struct {
		arg1 connect.ConnectorConfig
	}{arg1})
	fake.recordInvocation("ValidatePlugins", []interface{}{arg1})
	fake.validatePluginsMutex.Unlock()
	if fake.ValidatePluginsStub != nil {
		return fake.ValidatePluginsStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.validatePluginsReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) ValidatePluginsCallCount() int {
	fake.validatePluginsMutex.RLock()
	defer fake.validatePluginsMutex.RUnlock()
	return len(fake.validatePluginsArgsForCall)
}

func (fake *FakeClient) ValidatePluginsCalls(stub func(connect.ConnectorConfig) (*connect.ConfigValidation, *http.Response, error)) {
	fake.validatePluginsMutex.Lock()
	defer fake.validatePluginsMutex.Unlock()
	fake.ValidatePluginsStub = stub
}

func (fake *FakeClient) ValidatePluginsArgsForCall(i int) connect.ConnectorConfig {
	fake.validatePluginsMutex.RLock()
	defer fake.validatePluginsMutex.RUnlock()
	argsForCall := fake.validatePluginsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) ValidatePluginsReturns(result1 *connect.ConfigValidation, result2 *http.Response, result3 error) {
	fake.validatePluginsMutex.Lock()
	defer fake.validatePluginsMutex.Unlock()
	fake.ValidatePluginsStub = nil
	fake.validatePluginsReturns = struct {
		result1 *connect.ConfigValidation
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) ValidatePluginsReturnsOnCall(i int, result1 *connect.ConfigValidation, result2 *http.Response, result3 error) {
	fake.validatePluginsMutex.Lock()
	defer fake.validatePluginsMutex.Unlock()
	fake.ValidatePluginsStub = nil
	if fake.validatePluginsReturnsOnCall == nil {
		fake.validatePluginsReturnsOnCall = make(map[int]struct {
			result1 *connect.ConfigValidation
			result2 *http.Response
			result3 error
		})
	}
	fake.validatePluginsReturnsOnCall[i] = struct {
		result1 *connect.ConfigValidation
		result2 *http.Response
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createConnectorMutex.RLock()
	defer fake.createConnectorMutex.RUnlock()
	fake.deleteConnectorMutex.RLock()
	defer fake.deleteConnectorMutex.RUnlock()
	fake.getClusterInfoMutex.RLock()
	defer fake.getClusterInfoMutex.RUnlock()
	fake.getConnectorMutex.RLock()
	defer fake.getConnectorMutex.RUnlock()
	fake.getConnectorStatusMutex.RLock()
	defer fake.getConnectorStatusMutex.RUnlock()
	fake.listConnectorsMutex.RLock()
	defer fake.listConnectorsMutex.RUnlock()
	fake.listPluginsMutex.RLock()
	defer fake.listPluginsMutex.RUnlock()
	fake.pauseConnectorMutex.RLock()
	defer fake.pauseConnectorMutex.RUnlock()
	fake.restartConnectorMutex.RLock()
	defer fake.restartConnectorMutex.RUnlock()
	fake.restartConnectorTaskMutex.RLock()
	defer fake.restartConnectorTaskMutex.RUnlock()
	fake.resumeConnectorMutex.RLock()
	defer fake.resumeConnectorMutex.RUnlock()
	fake.updateConnectorConfigMutex.RLock()
	defer fake.updateConnectorConfigMutex.RUnlock()
	fake.validatePluginsMutex.RLock()
	defer fake.validatePluginsMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
