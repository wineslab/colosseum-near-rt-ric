		//// Copyright 2019 AT&T Intellectual Property
//// Copyright 2019 Nokia
////
//// Licensed under the Apache License, Version 2.0 (the "License");
//// you may not use this file except in compliance with the License.
//// You may obtain a copy of the License at
////
////      http://www.apache.org/licenses/LICENSE-2.0
////
//// Unless required by applicable law or agreed to in writing, software
//// distributed under the License is distributed on an "AS IS" BASIS,
//// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//// See the License for the specific language governing permissions and
//// limitations under the License.

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package rmrmsghandlers

import (
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/services"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"testing"
)

const ranName = "test"
const e2tAddress = "10.10.2.15:9800"

func TestLostConnectionHandlerSuccess(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)

	notificationRequest := models.NotificationRequest{RanName: ranName}
	ranDisconnectionManagerMock := &mocks.RanDisconnectionManagerMock{}
	ranDisconnectionManagerMock.On("DisconnectRan", ranName).Return(nil)
	handler := NewRanLostConnectionHandler(logger, ranDisconnectionManagerMock)
	handler.Handle(&notificationRequest)
	ranDisconnectionManagerMock.AssertCalled(t, "DisconnectRan", ranName)
}

func TestLostConnectionHandlerFailure(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)

	notificationRequest := models.NotificationRequest{RanName: ranName}
	ranDisconnectionManagerMock := &mocks.RanDisconnectionManagerMock{}
	ranDisconnectionManagerMock.On("DisconnectRan", ranName).Return(errors.New("error"))
	handler := NewRanLostConnectionHandler(logger, ranDisconnectionManagerMock)
	handler.Handle(&notificationRequest)
	ranDisconnectionManagerMock.AssertCalled(t, "DisconnectRan", ranName)
}

func setupLostConnectionHandlerTestWithRealDisconnectionManager(t *testing.T, isSuccessfulHttpPost bool) (RanLostConnectionHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.HttpClientMock) {
	logger, _ := logger.InitLogger(logger.InfoLevel)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, logger)
	httpClientMock := &mocks.HttpClientMock{}
	routingManagerClient := clients.NewRoutingManagerClient(logger, config, httpClientMock)
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManager, routingManagerClient)
	ranDisconnectionManager := managers.NewRanDisconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, e2tAssociationManager)
	handler := NewRanLostConnectionHandler(logger, ranDisconnectionManager)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTING, AssociatedE2TInstanceAddress: e2tAddress}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo1 := *origNodebInfo
	updatedNodebInfo1.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo1).Return(rnibErr)
	updatedNodebInfo2 := *origNodebInfo
	updatedNodebInfo2.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	updatedNodebInfo2.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo2).Return(rnibErr)
	e2tInstance := &entities.E2TInstance{Address: e2tAddress, AssociatedRanList: []string{ranName}}
	readerMock.On("GetE2TInstance", e2tAddress).Return(e2tInstance, nil)
	e2tInstanceToSave := *e2tInstance
	e2tInstanceToSave.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &e2tInstanceToSave).Return(nil)
	mockHttpClient(httpClientMock, isSuccessfulHttpPost)

	return handler, readerMock, writerMock, httpClientMock
}

func mockHttpClient(httpClientMock *mocks.HttpClientMock, isSuccessful bool) {
	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(e2tAddress, RanName)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	var respStatusCode int
	if isSuccessful {
		respStatusCode = http.StatusCreated
	} else {
		respStatusCode = http.StatusBadRequest
	}
	httpClientMock.On("Post", clients.DissociateRanE2TInstanceApiSuffix, "application/json", body).Return(&http.Response{StatusCode: respStatusCode, Body: respBody}, nil)
}

func TestLostConnectionHandlerFailureWithRealDisconnectionManager(t *testing.T) {
	handler, readerMock, writerMock, httpClientMock := setupLostConnectionHandlerTestWithRealDisconnectionManager(t, false)

	notificationRequest := models.NotificationRequest{RanName: ranName}
	handler.Handle(&notificationRequest)

	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
}

func TestLostConnectionHandlerSuccessWithRealDisconnectionManager(t *testing.T) {
	handler, readerMock, writerMock, httpClientMock := setupLostConnectionHandlerTestWithRealDisconnectionManager(t, true)

	notificationRequest := models.NotificationRequest{RanName: ranName}
	handler.Handle(&notificationRequest)

	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)
	httpClientMock.AssertExpectations(t)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
}
