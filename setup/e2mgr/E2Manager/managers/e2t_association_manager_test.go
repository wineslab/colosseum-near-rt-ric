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
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

const RanName = "test"

func initE2TAssociationManagerTest(t *testing.T) (*E2TAssociationManager, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.HttpClientMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)

	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, log)
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	manager := NewE2TAssociationManager(log, rnibDataService, e2tInstancesManager, rmClient)

	return manager, readerMock, writerMock, httpClientMock
}

func mockHttpClient(httpClientMock *mocks.HttpClientMock, apiSuffix string, isSuccessful bool) {
	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress, RanName)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	var respStatusCode int
	if isSuccessful {
		respStatusCode = http.StatusCreated
	} else {
		respStatusCode = http.StatusBadRequest
	}
	httpClientMock.On("Post", apiSuffix, "application/json", body).Return(&http.Response{StatusCode: respStatusCode, Body: respBody}, nil)
}

func TestAssociateRanSuccess(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClient(httpClientMock, clients.AssociateRanToE2TInstanceApiSuffix, true)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress}
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	updatedE2tInstance := *e2tInstance
	updatedE2tInstance.AssociatedRanList = append(updatedE2tInstance.AssociatedRanList, RanName)
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)

	err := manager.AssociateRan(E2TAddress, nb)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestAssociateRanRoutingManagerError(t *testing.T) {
	manager, _, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClient(httpClientMock, clients.AssociateRanToE2TInstanceApiSuffix, false)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
	writerMock.On("UpdateNodebInfo", nb).Return(nil)

	err := manager.AssociateRan(E2TAddress, nb)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RoutingManagerError{}, err)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestAssociateRanUpdateNodebError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClient(httpClientMock, clients.AssociateRanToE2TInstanceApiSuffix, true)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(e2managererrors.NewRnibDbError())

	err := manager.AssociateRan(E2TAddress, nb)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestAssociateRanGetE2tInstanceError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClient(httpClientMock, clients.AssociateRanToE2TInstanceApiSuffix, true)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	var e2tInstance *entities.E2TInstance
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, errors.New("test"))

	err := manager.AssociateRan(E2TAddress, nb)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestAssociateRanSaveE2tInstanceError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClient(httpClientMock, clients.AssociateRanToE2TInstanceApiSuffix, true)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress}
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	updatedE2tInstance := *e2tInstance
	updatedE2tInstance.AssociatedRanList = append(updatedE2tInstance.AssociatedRanList, RanName)
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(errors.New("test"))

	err := manager.AssociateRan(E2TAddress, nb)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestDissociateRanSuccess(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClient(httpClientMock, clients.DissociateRanE2TInstanceApiSuffix, true)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress}
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	updatedE2tInstance := *e2tInstance
	updatedE2tInstance.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)

	err := manager.DissociateRan(E2TAddress, RanName)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestDissociateRanGetNodebError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	var nb *entities.NodebInfo
	readerMock.On("GetNodeb", RanName).Return(nb, e2managererrors.NewRnibDbError())

	err := manager.DissociateRan(E2TAddress, RanName)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestDissociateRanUpdateNodebError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(e2managererrors.NewRnibDbError())

	err := manager.DissociateRan(E2TAddress, RanName)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestDissociateRanGetE2tInstanceError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	var e2tInstance *entities.E2TInstance
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, errors.New("test"))

	err := manager.DissociateRan(E2TAddress, RanName)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestDissociateRanSaveE2tInstanceError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress}
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	updatedE2tInstance := *e2tInstance
	updatedE2tInstance.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(errors.New("test"))

	err := manager.DissociateRan(E2TAddress, RanName)

	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestDissociateRanRoutingManagerError(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClient(httpClientMock, clients.DissociateRanE2TInstanceApiSuffix, false)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress}
	e2tInstance.AssociatedRanList = append(e2tInstance.AssociatedRanList, RanName)
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	updatedE2tInstance := *e2tInstance
	updatedE2tInstance.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &updatedE2tInstance).Return(nil)

	err := manager.DissociateRan(E2TAddress, RanName)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestRemoveE2tInstanceSuccessWithOrphans(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)

	ranNamesToBeDissociated := []string{RanName, "test1"}
	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, ranNamesToBeDissociated, nil)
	mockHttpClientDelete(httpClientMock, data, true)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tAddressesNew := []string{}
	writerMock.On("SaveE2TAddresses", e2tAddressesNew).Return(nil)

	e2tInstance1 := &entities.E2TInstance{Address: E2TAddress, AssociatedRanList:ranNamesToBeDissociated}
	err := manager.RemoveE2tInstance(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestRemoveE2tInstanceFailureRoutingManager(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test1"}, nil)
	mockHttpClientDelete(httpClientMock, data, false)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	e2tAddresses := []string{E2TAddress}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tAddressesNew := []string{}
	writerMock.On("SaveE2TAddresses", e2tAddressesNew).Return(nil)

	e2tInstance1 := &entities.E2TInstance{Address: E2TAddress, AssociatedRanList:[]string{"test1"}}
	//readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance1, e2managererrors.NewRnibDbError())
	err := manager.RemoveE2tInstance(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestRemoveE2tInstanceFailureInE2TInstanceManager(t *testing.T) {

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test1"}, nil)
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)
	mockHttpClientDelete(httpClientMock, data, true)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	var e2tAddresses []string
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, e2managererrors.NewRnibDbError())

	e2tInstance1 := &entities.E2TInstance{Address: E2TAddress, AssociatedRanList:[]string{"test1"}}
	err := manager.RemoveE2tInstance(e2tInstance1)

	assert.NotNil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestRemoveE2tInstanceFailureInE2tInstanceAddRansToInstance(t *testing.T) {
	manager, readerMock, writerMock, httpClientMock := initE2TAssociationManagerTest(t)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, nil, nil)
	mockHttpClientDelete(httpClientMock, data, true)

	writerMock.On("RemoveE2TInstance", E2TAddress).Return(nil)
	e2tAddresses := []string{E2TAddress, E2TAddress2, E2TAddress3}
	readerMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tAddressesNew := []string{E2TAddress2, E2TAddress3}
	writerMock.On("SaveE2TAddresses", e2tAddressesNew).Return(nil)

	e2tInstance1 := &entities.E2TInstance{Address: E2TAddress}
	err := manager.RemoveE2tInstance(e2tInstance1)

	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func mockHttpClientDelete(httpClientMock *mocks.HttpClientMock, data *models.RoutingManagerDeleteRequestModel, isSuccessful bool) {

	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	var respStatusCode int
	if isSuccessful {
		respStatusCode = http.StatusCreated
	} else {
		respStatusCode = http.StatusBadRequest
	}
	httpClientMock.On("Delete", "e2t", "application/json", body).Return(&http.Response{StatusCode: respStatusCode, Body: respBody}, nil)
}
