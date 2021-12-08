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
	"e2mgr/e2pdus"
	"e2mgr/enums"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/tests"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"testing"
	"time"
	"unsafe"
)

func initX2ResetRequestNotificationHandlerTest(t *testing.T) (X2ResetRequestNotificationHandler, *mocks.RnibReaderMock, *mocks.RmrMessengerMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, nil)

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, log)
	ranStatusChangeManager := managers.NewRanStatusChangeManager(log, rmrSender)
	h := NewX2ResetRequestNotificationHandler(log, rnibDataService, ranStatusChangeManager, rmrSender)
	return h, readerMock, rmrMessengerMock
}

func getRanRestartedMbuf(nodeType entities.Node_Type, messageDirection enums.MessageDirection) *rmrCgo.MBuf {
	var xAction []byte
	resourceStatusPayload := models.NewResourceStatusPayload(nodeType, messageDirection)
	resourceStatusJson, _ := json.Marshal(resourceStatusPayload)
	var msgSrc unsafe.Pointer
	return rmrCgo.NewMBuf(rmrCgo.RAN_RESTARTED, len(resourceStatusJson), RanName, &resourceStatusJson, &xAction, msgSrc)
}

func TestHandleX2ResetRequestNotificationSuccess(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetRequestNotificationHandlerTest(t)
	ranName := "test"
	xAction := []byte("123456aa")
	notificationRequest := models.NewNotificationRequest(ranName, []byte{}, time.Now(), xAction, nil)

	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTED, NodeType: entities.Node_ENB}
	var err error
	readerMock.On("GetNodeb", ranName).Return(nb, err)
	var msgSrc unsafe.Pointer
	resetResponseMbuf := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET_RESP, len(e2pdus.PackedX2ResetResponse), ranName, &e2pdus.PackedX2ResetResponse, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", resetResponseMbuf, true).Return(&rmrCgo.MBuf{}, err)
	ranRestartedMbuf := getRanRestartedMbuf(nb.NodeType, enums.RAN_TO_RIC)
	rmrMessengerMock.On("SendMsg", ranRestartedMbuf, true).Return(&rmrCgo.MBuf{}, err)
	h.Handle(notificationRequest)
	rmrMessengerMock.AssertCalled(t, "SendMsg", resetResponseMbuf, true)
	rmrMessengerMock.AssertCalled(t, "SendMsg", ranRestartedMbuf, true)
}

func TestHandleX2ResetRequestNotificationShuttingDownStatus(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetRequestNotificationHandlerTest(t)
	var payload []byte

	xAction := []byte("123456aa")
	var msgSrc unsafe.Pointer
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload), "RanName", &payload, &xAction, msgSrc)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId: xAction}

	nb := &entities.NodebInfo{RanName: mBuf.Meid, ConnectionStatus: entities.ConnectionStatus_SHUTTING_DOWN,}
	var rnibErr error

	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestHandleX2ResetRequestNotificationDisconnectStatus(t *testing.T) {
	h, readerMock, rmrMessengerMock := initX2ResetRequestNotificationHandlerTest(t)
	var payload []byte
	xAction := []byte("123456aa")
	var msgSrc unsafe.Pointer
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload), "RanName", &payload, &xAction, msgSrc)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload, StartTime: time.Now(), TransactionId: xAction}
	nb := &entities.NodebInfo{RanName: mBuf.Meid, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED,}
	var rnibErr error
	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}

func TestHandleX2ResetRequestNotificationGetNodebFailed(t *testing.T) {

	 h, readerMock, rmrMessengerMock := initX2ResetRequestNotificationHandlerTest(t)
	var payload []byte
	var xAction []byte
	var msgSrc unsafe.Pointer
	mBuf := rmrCgo.NewMBuf(tests.MessageType, len(payload), "RanName", &payload, &xAction, msgSrc)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId: xAction}

	var nb *entities.NodebInfo
	rnibErr := &common.ResourceNotFoundError{}
	readerMock.On("GetNodeb", mBuf.Meid).Return(nb, rnibErr)

	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertNotCalled(t, "SendMsg")
}
