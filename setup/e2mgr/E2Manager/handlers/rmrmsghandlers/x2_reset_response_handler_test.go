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


package rmrmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/converters"
	"e2mgr/enums"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"testing"
	"time"
)

const (
	SuccessfulX2ResetResponsePackedPdu = "200700080000010011400100"
	SuccessfulX2ResetResponsePackedPduEmptyIEs = "20070003000000"
	UnsuccessfulX2ResetResponsePackedPdu = "2007000d00000100114006080000000d00"
)

func initX2ResetResponseHandlerTest(t *testing.T) (X2ResetResponseHandler, *mocks.RnibReaderMock, *mocks.RmrMessengerMock) {
	log, err := logger.InitLogger(logger.InfoLevel)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	if err != nil {
		t.Errorf("#initX2ResetResponseHandlerTest - failed to initialize logger, error: %s", err)
	}
	readerMock := &mocks.RnibReaderMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, nil)

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, log)
	ranStatusChangeManager := managers.NewRanStatusChangeManager(log, rmrSender)

	h := NewX2ResetResponseHandler(log, rnibDataService, ranStatusChangeManager, converters.NewX2ResetResponseExtractor(log))
	return h, readerMock, rmrMessengerMock
}

func TestX2ResetResponseSuccess(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetResponseHandlerTest(t)
	var payload []byte
	_, err := fmt.Sscanf(SuccessfulX2ResetResponsePackedPdu, "%x", &payload)
	if err != nil {
		t.Fatalf("Failed converting packed pdu. Error: %v\n", err)
	}

	var xAction []byte
	notificationRequest := models.NotificationRequest{RanName: RanName, Len: len(payload), Payload: payload, StartTime: time.Now(), TransactionId: xAction}
	nb := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, NodeType: entities.Node_ENB}
	var rnibErr error
	readerMock.On("GetNodeb", RanName).Return(nb, rnibErr)
	ranRestartedMbuf := getRanRestartedMbuf(nb.NodeType, enums.RIC_TO_RAN)
	rmrMessengerMock.On("SendMsg", ranRestartedMbuf, true).Return(&rmrCgo.MBuf{}, err)
	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertCalled(t, "SendMsg", ranRestartedMbuf, true)
}

func TestX2ResetResponseSuccessEmptyIEs(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetResponseHandlerTest(t)
	var payload []byte
	_, err := fmt.Sscanf(SuccessfulX2ResetResponsePackedPduEmptyIEs, "%x", &payload)
	if err != nil {
		t.Fatalf("Failed converting packed pdu. Error: %v\n", err)
	}

	var xAction []byte
	notificationRequest := models.NotificationRequest{RanName: RanName, Len: len(payload), Payload: payload, StartTime: time.Now(), TransactionId: xAction}
	nb := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, NodeType: entities.Node_ENB}
	var rnibErr error
	readerMock.On("GetNodeb", RanName).Return(nb, rnibErr)
	ranRestartedMbuf := getRanRestartedMbuf(nb.NodeType, enums.RIC_TO_RAN)
	rmrMessengerMock.On("SendMsg", ranRestartedMbuf, true).Return(&rmrCgo.MBuf{}, err)
	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertCalled(t, "SendMsg", ranRestartedMbuf, true)
}

func TestX2ResetResponseShuttingDown(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetResponseHandlerTest(t)
	var payload []byte
	_, err := fmt.Sscanf(SuccessfulX2ResetResponsePackedPdu, "%x", &payload)
	if err != nil {
		t.Fatalf("Failed converting packed pdu. Error: %v\n", err)
	}

	var xAction []byte
	notificationRequest := models.NotificationRequest{RanName: RanName, Len: len(payload), Payload: payload, StartTime: time.Now(), TransactionId: xAction}
	nb := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN, NodeType: entities.Node_ENB}
	var rnibErr error
	readerMock.On("GetNodeb", RanName).Return(nb, rnibErr)
	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestX2ResetResponseInvalidConnectionStatus(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetResponseHandlerTest(t)
	var payload []byte
	_, err := fmt.Sscanf(SuccessfulX2ResetResponsePackedPdu, "%x", &payload)
	if err != nil {
		t.Fatalf("Failed converting packed pdu. Error: %v\n", err)
	}

	var xAction []byte
	notificationRequest := models.NotificationRequest{RanName: RanName, Len: len(payload), Payload: payload, StartTime: time.Now(), TransactionId: xAction}
	nb := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, NodeType: entities.Node_ENB}
	var rnibErr error
	readerMock.On("GetNodeb", RanName).Return(nb, rnibErr)
	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestX2ResetResponseError(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetResponseHandlerTest(t)
	var payload []byte
	_, err := fmt.Sscanf(UnsuccessfulX2ResetResponsePackedPdu, "%x", &payload)
	if err != nil {
		t.Fatalf("Failed converting packed pdu. Error: %v\n", err)
	}

	var xAction []byte
	notificationRequest := models.NotificationRequest{RanName: RanName, Len: len(payload), Payload: payload, StartTime: time.Now(), TransactionId: xAction}
	nb := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, NodeType: entities.Node_ENB}
	var rnibErr error
	readerMock.On("GetNodeb", RanName).Return(nb, rnibErr)
	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestX2ResetResponseGetNodebFailure(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetResponseHandlerTest(t)

	var payload []byte
	_, err := fmt.Sscanf(SuccessfulX2ResetResponsePackedPdu, "%x", &payload)
	if err != nil {
		t.Fatalf("Failed converting packed pdu. Error: %v\n", err)
	}

	var xAction []byte
	notificationRequest := models.NotificationRequest{RanName: RanName, Len: len(payload), Payload: payload, StartTime: time.Now(), TransactionId: xAction}

	var nb *entities.NodebInfo
	rnibErr := common.NewResourceNotFoundError("nodeb not found")
	readerMock.On("GetNodeb", RanName).Return(nb, rnibErr)

	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestX2ResetResponseUnpackFailure(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetResponseHandlerTest(t)

	payload := []byte("Invalid payload")
	var xAction []byte
	notificationRequest := models.NotificationRequest{RanName: RanName, Len: len(payload), Payload: payload, StartTime: time.Now(), TransactionId: xAction}
	nb := &entities.NodebInfo{RanName: RanName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, NodeType: entities.Node_ENB}
	var rnibErr error
	readerMock.On("GetNodeb", RanName).Return(nb, rnibErr)

	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}
