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
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const ranName = "test"
const e2tAddress = "10.10.2.15:9800"

func initRanLostConnectionTest(t *testing.T) (*logger.Logger, *mocks.RmrMessengerMock, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *RanDisconnectionManager, *mocks.HttpClientMock) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, logger)
	httpClient := &mocks.HttpClientMock{}
	routingManagerClient := clients.NewRoutingManagerClient(logger, config, httpClient)
	e2tAssociationManager := NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManager, routingManagerClient)
	ranDisconnectionManager := NewRanDisconnectionManager(logger, configuration.ParseConfiguration(), rnibDataService, e2tAssociationManager)
	return logger, rmrMessengerMock, readerMock, writerMock, ranDisconnectionManager, httpClient
}

func TestRanDisconnectionGetNodebFailure(t *testing.T) {
	_, _, readerMock, writerMock, ranDisconnectionManager, _ := initRanLostConnectionTest(t)

	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, common.NewInternalError(errors.New("Error")))
	err := ranDisconnectionManager.DisconnectRan(ranName)
	assert.NotNil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
}

func TestShutdownRan(t *testing.T) {
	_, _, readerMock, writerMock, ranDisconnectionManager, _ := initRanLostConnectionTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_SHUT_DOWN}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	err := ranDisconnectionManager.DisconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNotCalled(t, "UpdateNodebInfo")
}

func TestShuttingdownRan(t *testing.T) {
	_, _, readerMock, writerMock, ranDisconnectionManager, _ := initRanLostConnectionTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(rnibErr)
	err := ranDisconnectionManager.DisconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestShuttingDownRanUpdateNodebInfoFailure(t *testing.T) {
	_, _, readerMock, writerMock, ranDisconnectionManager, _ := initRanLostConnectionTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_SHUT_DOWN
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(common.NewInternalError(errors.New("Error")))
	err := ranDisconnectionManager.DisconnectRan(ranName)
	assert.NotNil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestConnectingRanUpdateNodebInfoFailure(t *testing.T) {
	_, _, readerMock, writerMock, ranDisconnectionManager, _ := initRanLostConnectionTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTING}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo := *origNodebInfo
	updatedNodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo).Return(common.NewInternalError(errors.New("Error")))
	err := ranDisconnectionManager.DisconnectRan(ranName)
	assert.NotNil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
}

func TestConnectingRanDisconnectSucceeds(t *testing.T) {
	_, _, readerMock, writerMock, ranDisconnectionManager, httpClient := initRanLostConnectionTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTING, AssociatedE2TInstanceAddress: E2TAddress}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo1 := *origNodebInfo
	updatedNodebInfo1.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo1).Return(rnibErr)
	updatedNodebInfo2 := *origNodebInfo
	updatedNodebInfo2.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	updatedNodebInfo2.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo2).Return(rnibErr)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress, AssociatedRanList: []string{ranName}}
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	e2tInstanceToSave := * e2tInstance
	e2tInstanceToSave.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &e2tInstanceToSave).Return(nil)
	mockHttpClient(httpClient, clients.DissociateRanE2TInstanceApiSuffix, true)
	err := ranDisconnectionManager.DisconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
}

func TestConnectingRanDissociateFailsRmError(t *testing.T) {
	_, _, readerMock, writerMock, ranDisconnectionManager, httpClient := initRanLostConnectionTest(t)

	origNodebInfo := &entities.NodebInfo{RanName: ranName, GlobalNbId: &entities.GlobalNbId{PlmnId: "xxx", NbId: "yyy"}, ConnectionStatus: entities.ConnectionStatus_CONNECTING, AssociatedE2TInstanceAddress: E2TAddress}
	var rnibErr error
	readerMock.On("GetNodeb", ranName).Return(origNodebInfo, rnibErr)
	updatedNodebInfo1 := *origNodebInfo
	updatedNodebInfo1.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo1).Return(rnibErr)
	updatedNodebInfo2 := *origNodebInfo
	updatedNodebInfo2.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	updatedNodebInfo2.AssociatedE2TInstanceAddress = ""
	writerMock.On("UpdateNodebInfo", &updatedNodebInfo2).Return(rnibErr)
	e2tInstance := &entities.E2TInstance{Address: E2TAddress, AssociatedRanList: []string{ranName}}
	readerMock.On("GetE2TInstance", E2TAddress).Return(e2tInstance, nil)
	e2tInstanceToSave := * e2tInstance
	e2tInstanceToSave.AssociatedRanList = []string{}
	writerMock.On("SaveE2TInstance", &e2tInstanceToSave).Return(nil)
	mockHttpClient(httpClient, clients.DissociateRanE2TInstanceApiSuffix, false)
	err := ranDisconnectionManager.DisconnectRan(ranName)
	assert.Nil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
}

func TestConnectingRanDissociateFailsDbError(t *testing.T) {
	_, _, readerMock, writerMock, ranDisconnectionManager, _ := initRanLostConnectionTest(t)

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
	readerMock.On("GetE2TInstance", e2tAddress).Return(e2tInstance, common.NewInternalError(errors.New("Error")))
	err := ranDisconnectionManager.DisconnectRan(ranName)
	assert.NotNil(t, err)
	readerMock.AssertCalled(t, "GetNodeb", ranName)
	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	writerMock.AssertNotCalled(t, "SaveE2TInstance", )
}

func initRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}
