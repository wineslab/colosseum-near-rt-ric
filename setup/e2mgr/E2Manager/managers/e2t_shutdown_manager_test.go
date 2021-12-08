//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package managers

import (
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/services"
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	//"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/client-go/kubernetes/fake"
	"net/http"
	"testing"
	"time"
)

const E2TAddress3 = "10.10.2.17:9800"

func initE2TShutdownManagerTest(t *testing.T) (*E2TShutdownManager, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.HttpClientMock, *KubernetesManager) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3, E2TInstanceDeletionTimeoutMs: 15000}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)

	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, log)
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	associationManager := NewE2TAssociationManager(log, rnibDataService, e2tInstancesManager, rmClient)
	//kubernetesManager := initKubernetesManagerTest(t)

	/*shutdownManager := NewE2TShutdownManager(log, config, rnibDataService, e2tInstancesManager, associationManager, kubernetesManager)

	return shutdownManager, readerMock, writerMock, httpClientMock, kubernetesManager*/
	shutdownManager := NewE2TShutdownManager(log, config, rnibDataService, e2tInstancesManager, associationManager, nil)

	return shutdownManager, readerMock, writerMock, httpClientMock, nil
}

func TestShutdownSuccess1OutOf3Instances(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_ := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	e2tInstance3 := entities.NewE2TInstance(E2TAddress3, PodName)
	e2tInstance3.State = entities.Active
	e2tInstance3.AssociatedRanList = []string{"test4"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)
	nodeb5 := &entities.NodebInfo{RanName:"test5", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test5").Return(nodeb5, nil)

	e2tAddresses := []string{E2TAddress, E2TAddress2,E2TAddress3}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, e2tInstance1.AssociatedRanList, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	writerMock.On("SaveE2TAddresses", []string{E2TAddress2,E2TAddress3}).Return(nil)

	nodeb1connected := *nodeb1
	nodeb1connected.AssociatedE2TInstanceAddress = ""
	nodeb1connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb1connected).Return(nil)
	nodeb2connected := *nodeb2
	nodeb2connected.AssociatedE2TInstanceAddress = ""
	nodeb2connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb2connected).Return(nil)
	nodeb5connected := *nodeb5
	nodeb5connected.AssociatedE2TInstanceAddress = ""
	nodeb5connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb5connected).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestShutdownSuccess1InstanceWithoutRans(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, nil, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	writerMock.On("SaveE2TAddresses", []string{}).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestShutdownSuccess1Instance2Rans(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test1", "test2"}, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	writerMock.On("SaveE2TAddresses", []string{}).Return(nil)

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = ""
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb1new).Return(nil)
	nodeb2new := *nodeb2
	nodeb2new.AssociatedE2TInstanceAddress = ""
	nodeb2new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb2new).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)

}

func TestShutdownE2tInstanceAlreadyBeingDeleted(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.ToBeDeleted
	e2tInstance1.AssociatedRanList = []string{"test1"}
	e2tInstance1.DeletionTimestamp = time.Now().UnixNano()

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)

}

func TestShutdownFailureMarkInstanceAsToBeDeleted(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(e2managererrors.NewRnibDbError())

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.NotNil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)

}

func TestShutdownFailureRoutingManagerError(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	e2tInstance3 := entities.NewE2TInstance(E2TAddress3, PodName)
	e2tInstance3.State = entities.Active
	e2tInstance3.AssociatedRanList = []string{"test4"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)
	nodeb5 := &entities.NodebInfo{RanName:"test5", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test5").Return(nodeb5, nil)

	e2tAddresses := []string{E2TAddress, E2TAddress2,E2TAddress3}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, e2tInstance1.AssociatedRanList, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	writerMock.On("SaveE2TAddresses", []string{E2TAddress2,E2TAddress3}).Return(nil)

	nodeb1connected := *nodeb1
	nodeb1connected.AssociatedE2TInstanceAddress = ""
	nodeb1connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb1connected).Return(nil)
	nodeb2connected := *nodeb2
	nodeb2connected.AssociatedE2TInstanceAddress = ""
	nodeb2connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb2connected).Return(nil)
	nodeb5connected := *nodeb5
	nodeb5connected.AssociatedE2TInstanceAddress = ""
	nodeb5connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb5connected).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)

}

func TestShutdownFailureInClearNodebsAssociation(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = ""
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb1new).Return(common.NewInternalError(fmt.Errorf("for tests")))

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.NotNil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestShutdownResourceNotFoundErrorInGetNodeb(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	var nodeb2 *entities.NodebInfo
	readerMock.On("GetNodeb", "test2").Return(nodeb2, common.NewResourceNotFoundError("for testing"))

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = ""
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb1new).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.NotNil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestShutdownResourceGeneralErrorInGetNodeb(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	var nodeb1 *entities.NodebInfo
	readerMock.On("GetNodeb", "test1").Return(nodeb1, common.NewInternalError(fmt.Errorf("for testing")))
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test1", "test2"}, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	writerMock.On("SaveE2TAddresses", []string{}).Return(nil)

	nodeb2new := *nodeb2
	nodeb2new.AssociatedE2TInstanceAddress = ""
	nodeb2new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb2new).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)

}

func TestShutdownFailureInRemoveE2TInstance(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,_  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test1", "test2", "test5"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.State = entities.Active
	e2tInstance2.AssociatedRanList = []string{"test3"}
	e2tInstance3 := entities.NewE2TInstance(E2TAddress3, PodName)
	e2tInstance3.State = entities.Active
	e2tInstance3.AssociatedRanList = []string{"test4"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test1").Return(nodeb1, nil)
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_SHUTTING_DOWN, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)
	nodeb5 := &entities.NodebInfo{RanName:"test5", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test5").Return(nodeb5, nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, e2tInstance1.AssociatedRanList, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(common.NewInternalError(fmt.Errorf("for tests")))

	nodeb1connected := *nodeb1
	nodeb1connected.AssociatedE2TInstanceAddress = ""
	nodeb1connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb1connected).Return(nil)
	nodeb2connected := *nodeb2
	nodeb2connected.AssociatedE2TInstanceAddress = ""
	nodeb2connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb2connected).Return(nil)
	nodeb5connected := *nodeb5
	nodeb5connected.AssociatedE2TInstanceAddress = ""
	nodeb5connected.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb5connected).Return(nil)

	err := shutdownManager.Shutdown(e2tInstance1)

	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}
/*
func TestShutdownSuccess2Instance2Rans(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,kubernetesManager  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test2"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress2, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test2"}, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	writerMock.On("SaveE2TAddresses", []string{}).Return(nil)

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = ""
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	nodeb2new := *nodeb2
	nodeb2new.AssociatedE2TInstanceAddress = ""
	nodeb2new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb2new).Return(nil)

	test := TestStruct{
		description: "namespace, 2 pods in Oran",
		namespace:   "oran",
		objs:        []runtime.Object{pod("oran", PodName), pod("oran", "e2t_2"), pod("some-namespace", "POD_Test_1")},
	}

	t.Run(test.description, func(t *testing.T) {
		kubernetesManager.ClientSet = fake.NewSimpleClientset(test.objs...)

		err := shutdownManager.Shutdown(e2tInstance1)

		assert.Nil(t, err)
		readerMock.AssertExpectations(t)
		writerMock.AssertExpectations(t)
		httpClientMock.AssertExpectations(t)
	})
}

func TestShutdownSuccess2Instance2RansNoPod(t *testing.T) {
	shutdownManager, readerMock, writerMock, httpClientMock,kubernetesManager  := initE2TShutdownManagerTest(t)

	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.State = entities.Active
	e2tInstance1.AssociatedRanList = []string{"test2"}
	writerMock.On("SaveE2TInstance", mock.MatchedBy(func(e2tInstance *entities.E2TInstance) bool { return e2tInstance.Address == E2TAddress && e2tInstance.State == entities.ToBeDeleted })).Return(nil)

	nodeb1 := &entities.NodebInfo{RanName:"test1", AssociatedE2TInstanceAddress:E2TAddress2, ConnectionStatus:entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	nodeb2 := &entities.NodebInfo{RanName:"test2", AssociatedE2TInstanceAddress:E2TAddress, ConnectionStatus:entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", "test2").Return(nodeb2, nil)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test2"}, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	readerMock.On("GetE2TAddresses").Return([]string{E2TAddress}, nil)
	writerMock.On("SaveE2TAddresses", []string{}).Return(nil)

	nodeb1new := *nodeb1
	nodeb1new.AssociatedE2TInstanceAddress = ""
	nodeb1new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	nodeb2new := *nodeb2
	nodeb2new.AssociatedE2TInstanceAddress = ""
	nodeb2new.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &nodeb2new).Return(nil)

	test := TestStruct{
		description: "namespace, 2 pods in Oran",
		namespace:   "oran",
		objs:        []runtime.Object{pod("oran", "e2t_2"), pod("some-namespace", "POD_Test_1")},
	}

	t.Run(test.description, func(t *testing.T) {
		kubernetesManager.ClientSet = fake.NewSimpleClientset(test.objs...)

		err := shutdownManager.Shutdown(e2tInstance1)

		assert.Nil(t, err)
		readerMock.AssertExpectations(t)
		writerMock.AssertExpectations(t)
		httpClientMock.AssertExpectations(t)
	})
}*/