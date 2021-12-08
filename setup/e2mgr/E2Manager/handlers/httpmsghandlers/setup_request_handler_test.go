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

package httpmsghandlers

import (
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/e2pdus"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"testing"
)

const E2TAddress = "10.0.2.15:8989"
const RanName = "test"
const BaseRMUrl = "http://10.10.2.15:12020/routingmanager"

func initSetupRequestTest(t *testing.T, protocol entities.E2ApplicationProtocol) (*mocks.RnibReaderMock, *mocks.RnibWriterMock, *SetupRequestHandler, *mocks.E2TInstancesManagerMock, *mocks.RanSetupManagerMock, *mocks.HttpClientMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	config.RoutingManager.BaseUrl = BaseRMUrl

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)

	ranSetupManagerMock := &mocks.RanSetupManagerMock{}
	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	httpClientMock := &mocks.HttpClientMock{}
	mockHttpClientAssociateRan(httpClientMock)
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	e2tAssociationManager := managers.NewE2TAssociationManager(log, rnibDataService, e2tInstancesManagerMock, rmClient)
	handler := NewSetupRequestHandler(log, rnibDataService, ranSetupManagerMock, protocol, e2tInstancesManagerMock, e2tAssociationManager)

	return readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, httpClientMock
}

func initSetupRequestTestBasicMocks(t *testing.T, protocol entities.E2ApplicationProtocol) (*mocks.RnibReaderMock, *mocks.RnibWriterMock, *SetupRequestHandler, *mocks.RmrMessengerMock, *mocks.HttpClientMock, *mocks.E2TInstancesManagerMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	config.RoutingManager.BaseUrl = BaseRMUrl
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := getRmrSender(rmrMessengerMock, log)
	ranSetupManager := managers.NewRanSetupManager(log, rmrSender, rnibDataService)
	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	e2tAssociationManager := managers.NewE2TAssociationManager(log, rnibDataService, e2tInstancesManagerMock, rmClient)
	handler := NewSetupRequestHandler(log, rnibDataService, ranSetupManager, protocol, e2tInstancesManagerMock, e2tAssociationManager)

	return readerMock, writerMock, handler, rmrMessengerMock, httpClientMock, e2tInstancesManagerMock
}

func mockHttpClientAssociateRan(httpClientMock *mocks.HttpClientMock) {
	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress, RanName)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := BaseRMUrl + clients.AssociateRanToE2TInstanceApiSuffix
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)
}

func TestX2SetupHandleNoPortError(t *testing.T) {
	readerMock, _, handler, _, _, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{"127.0.0.1", 0, RanName,}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RequestValidationError{}, err)
	readerMock.AssertNotCalled(t, "GetNodeb")
}

func TestX2SetupHandleNoRanNameError(t *testing.T) {
	readerMock, _, handler, _, _, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{RanPort: 8080, RanIp: "127.0.0.1"}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RequestValidationError{}, err)
	readerMock.AssertNotCalled(t, "GetNodeb")
}

func TestX2SetupHandleNoIpError(t *testing.T) {
	readerMock, _, handler, _, _, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{RanPort: 8080, RanName: RanName}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RequestValidationError{}, err)
	readerMock.AssertNotCalled(t, "GetNodeb")
}

func TestX2SetupHandleInvalidIpError(t *testing.T) {
	readerMock, _, handler, _, _, _:= initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	sr := models.SetupRequest{RanPort: 8080, RanName: RanName, RanIp: "invalid ip"}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RequestValidationError{}, err)
	readerMock.AssertNotCalled(t, "GetNodeb")
}

func TestSetupGetNodebFailure(t *testing.T) {
	readerMock, _, handler, _, _, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)

	rnibErr := &common.ValidationError{}
	nb := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	readerMock.On("GetNodeb", RanName).Return(nb, rnibErr)

	sr := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(sr)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
}

func TestSetupNewRanSelectE2TInstancesDbError(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return("", e2managererrors.NewRnibDbError())
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance")
	writerMock.AssertNotCalled(t, "SaveNodeb")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupNewRanSelectE2TInstancesNoInstances(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nodebInfo, _ := createInitialNodeInfo(&setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	writerMock.On("SaveNodeb", mock.Anything, mock.Anything).Return(nil)
	updatedNb := *nodebInfo
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", &updatedNb, entities.ConnectionStatus_CONNECTING).Return(nil)
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.Nil(t, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
	ranSetupManagerMock.AssertExpectations(t)
}

func TestSetupNewRanAssociateRanFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, httpClientMock := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(e2managererrors.NewRnibDbError())
	setupRequest := &models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nb, nbIdentity := createInitialNodeInfo(setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("SaveNodeb", nbIdentity, mock.Anything).Return(nil)
	writerMock.On("UpdateNodebInfo", nb).Return(nil)
	nb.AssociatedE2TInstanceAddress = E2TAddress
	mockHttpClientAssociateRan(httpClientMock)
	updatedNb := *nb
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress

	_, err := handler.Handle(*setupRequest)
	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	e2tInstancesManagerMock.AssertExpectations(t)
	ranSetupManagerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}

func TestSetupNewRanSaveNodebFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nodebInfo, nbIdentity := createInitialNodeInfo(&setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(common.NewInternalError(fmt.Errorf("")))
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupNewRanSetupDbError(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(e2managererrors.NewRnibDbError())
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nodebInfo, nbIdentity := createInitialNodeInfo(&setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(nil)
	updatedNb := *nodebInfo
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	_, err := handler.Handle(setupRequest)
	assert.NotNil(t, err)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	ranSetupManagerMock.AssertExpectations(t)
}

func TestSetupNewRanSetupRmrError(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nodebInfo, nbIdentity := createInitialNodeInfo(&setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(nil)
	updatedNb := *nodebInfo
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", &updatedNb, entities.ConnectionStatus_CONNECTING).Return(e2managererrors.NewRmrError())
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.RmrError{}, err)
}

func TestSetupNewRanSetupSuccess(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError(""))
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	nodebInfo, nbIdentity := createInitialNodeInfo(&setupRequest, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	writerMock.On("SaveNodeb", nbIdentity, nodebInfo).Return(nil)
	updatedNb := *nodebInfo
	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
	updatedNb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", &updatedNb, entities.ConnectionStatus_CONNECTING).Return(nil)
	_, err := handler.Handle(setupRequest)
	assert.Nil(t, err)
}

func TestX2SetupExistingRanShuttingDown(t *testing.T) {
	readerMock, _, handler, e2tInstancesManagerMock, ranSetupManagerMock , _:= initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}, nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.WrongStateError{}, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "SelectE2TInstance")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestEndcSetupExistingRanShuttingDown(t *testing.T) {
	readerMock, _, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST)
	readerMock.On("GetNodeb", RanName).Return(&entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}, nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.WrongStateError{}, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "SelectE2TInstance")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceSelectDbError(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return("", e2managererrors.NewRnibDbError())
	updatedNb := *nb
	updatedNb.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceSelectNoInstanceError(t *testing.T) {
	readerMock, writerMock, handler, rmrMessengerMock, httpClientMock,e2tInstancesManagerMock:= initSetupRequestTestBasicMocks(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	readerMock.On("GetE2TAddresses").Return([]string{}, nil)
	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	e2tInstancesManagerMock.On("AddRansToInstance", "10.0.2.15:8989", []string{"test"}).Return(nil)
	mockHttpClientAssociateRan(httpClientMock)
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.InternalError{}, err)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
	writerMock.AssertExpectations(t)
}

func TestSetupExistingRanWithoutAssocE2TInstanceSelectNoInstanceErrorUpdateFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return("", e2managererrors.NewE2TInstanceAbsenceError())
	updatedNb := *nb
	updatedNb.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(common.NewInternalError(fmt.Errorf("")))
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.E2TInstanceAbsenceError{}, err)
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithoutAssocE2TInstanceSelectErrorAlreadyDisconnected(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: "", ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, e2managererrors.NewE2TInstanceAbsenceError())
	setupRequest := models.SetupRequest{"127.0.0.1", 8080, RanName,}
	_, err := handler.Handle(setupRequest)
	assert.IsType(t, &e2managererrors.E2TInstanceAbsenceError{}, err)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

//func TestSetupExistingRanWithoutAssocE2TInstanceAssociateRanFailure(t *testing.T) {
//	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
//	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
//	readerMock.On("GetNodeb", RanName).Return(nb, nil)
//	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
//	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(e2managererrors.NewRnibDbError())
//	writerMock.On("UpdateNodebInfo", mock.Anything).Return(nil)
//	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
//	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
//	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
//	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
//}

//func TestSetupExistingRanWithoutAssocE2TInstanceAssociateRanSucceedsUpdateNodebFails(t *testing.T) {
//	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
//	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
//	readerMock.On("GetNodeb", RanName).Return(nb, nil)
//	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
//	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(nil)
//	updatedNb := *nb
//	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
//	writerMock.On("UpdateNodebInfo", &updatedNb).Return(common.NewInternalError(fmt.Errorf("")))
//	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
//	assert.IsType(t, /* &e2managererrors.RnibDbError{} */&common.InternalError{}, err)
//	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
//}

//func TestSetupExistingRanWithoutAssocE2TInstanceExecuteSetupFailure(t *testing.T) {
//	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
//	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
//	readerMock.On("GetNodeb", RanName).Return(nb, nil)
//	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
//	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(nil)
//	updatedNb := *nb
//	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
//	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
//	ranSetupManagerMock.On("ExecuteSetup", &updatedNb, entities.ConnectionStatus_CONNECTING).Return(e2managererrors.NewRnibDbError())
//	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
//	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
//}
//
//func TestSetupExistingRanWithoutAssocE2TInstanceSuccess(t *testing.T) {
//	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
//	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: ""}
//	readerMock.On("GetNodeb", RanName).Return(nb, nil)
//	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
//	e2tInstancesManagerMock.On("AddRansToInstance", E2TAddress, []string{RanName}).Return(nil)
//	updatedNb := *nb
//	updatedNb.AssociatedE2TInstanceAddress = E2TAddress
//	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
//	ranSetupManagerMock.On("ExecuteSetup", &updatedNb, entities.ConnectionStatus_CONNECTING).Return(nil)
//	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
//	assert.Nil(t, err)
//}

func TestSetupExistingRanWithAssocE2TInstanceUpdateNodebFailure(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: E2TAddress}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(common.NewInternalError(fmt.Errorf("")))
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.RnibDbError{}, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "SelectE2TInstance")
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance")
	ranSetupManagerMock.AssertNotCalled(t, "ExecuteSetup")
}

func TestSetupExistingRanWithAssocE2TInstanceExecuteSetupRmrError(t *testing.T) {
	readerMock, writerMock, handler, rmrMessengerMock, _, _ := initSetupRequestTestBasicMocks(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: E2TAddress, ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	updatedNb3 := updatedNb
	updatedNb3.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	payload := e2pdus.PackedX2setupRequest
	xaction := []byte(RanName)
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), RanName, &payload, &xaction, nil)
	rmrMessengerMock.On("SendMsg",mock.Anything, true).Return(msg, e2managererrors.NewRmrError())
	writerMock.On("UpdateNodebInfo", &updatedNb3).Return(nil)
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.IsType(t, &e2managererrors.RmrError{}, err)
	writerMock.AssertExpectations(t)
}

func TestSetupExistingRanWithAssocE2TInstanceConnectedSuccess(t *testing.T) {
	readerMock, writerMock, handler, e2tInstancesManagerMock, ranSetupManagerMock, _ := initSetupRequestTest(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: E2TAddress, ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	updatedNb := *nb
	writerMock.On("UpdateNodebInfo", &updatedNb).Return(nil)
	ranSetupManagerMock.On("ExecuteSetup", &updatedNb, entities.ConnectionStatus_CONNECTED).Return(nil)
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.Nil(t, err)
	e2tInstancesManagerMock.AssertNotCalled(t, "SelectE2TInstance")
	e2tInstancesManagerMock.AssertNotCalled(t, "AddRansToInstance")
}

func TestSetupExistingRanWithoutAssocE2TInstanceExecuteRoutingManagerError(t *testing.T) {
	readerMock, writerMock, handler, rmrMessengerMock, httpClientMock, e2tInstancesManagerMock := initSetupRequestTestBasicMocks(t, entities.E2ApplicationProtocol_X2_SETUP_REQUEST)
	nb := &entities.NodebInfo{RanName: RanName, AssociatedE2TInstanceAddress: "", ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol:entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	readerMock.On("GetNodeb", RanName).Return(nb, nil)
	writerMock.On("UpdateNodebInfo", nb).Return(nil)
	e2tInstancesManagerMock.On("SelectE2TInstance").Return(E2TAddress, nil)
	mockHttpClientAssociateRan(httpClientMock)
	e2tInstancesManagerMock.On("AddRansToInstance", mock.Anything, mock.Anything).Return(nil)
	msg := &rmrCgo.MBuf{}
	var errNIl error
	rmrMessengerMock.On("SendMsg",mock.Anything, true).Return(msg, errNIl)
	_, err := handler.Handle(models.SetupRequest{"127.0.0.1", 8080, RanName,})
	assert.Nil(t, err)
	writerMock.AssertExpectations(t)
	readerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
}
