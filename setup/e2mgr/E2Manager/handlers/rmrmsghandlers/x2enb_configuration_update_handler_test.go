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
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"fmt"
	"testing"
	"time"
	"unsafe"
)

const PackedX2EnbConfigurationUpdateAck = "200800080000010011400100"
const PackedX2EnbConfigurationUpdateFailure = "400800080000010005400142"

func initX2EnbConfigurationUpdateHandlerTest(t *testing.T) (X2EnbConfigurationUpdateHandler, *mocks.RmrMessengerMock) {
	log := initLog(t)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, log)
	h := NewX2EnbConfigurationUpdateHandler(log, rmrSender)
	return h, rmrMessengerMock
}

func TestHandleX2EnbConfigUpdateSuccess(t *testing.T) {
	h, rmrMessengerMock := initX2EnbConfigurationUpdateHandlerTest(t)

	ranName := "test"
	xAction := []byte("123456aa")
	var payload []byte
	_, _ = fmt.Sscanf(PackedX2EnbConfigurationUpdateAck, "%x", &payload)
	var msgSrc unsafe.Pointer

	mBuf := rmrCgo.NewMBuf(rmrCgo.RIC_ENB_CONFIGURATION_UPDATE_ACK, len(payload), ranName, &payload, &xAction, msgSrc)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: mBuf.Len, Payload: *mBuf.Payload,
		StartTime: time.Now(), TransactionId:xAction}
	var err error
	rmrMessengerMock.On("SendMsg", mBuf, true).Return(&rmrCgo.MBuf{}, err)
	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mBuf, true)
}

func TestHandleX2EnbConfigUpdateFailure(t *testing.T) {
	h, rmrMessengerMock := initX2EnbConfigurationUpdateHandlerTest(t)

	ranName := "test"
	xAction := []byte("123456aa")

	var payload []byte
	_, _ = fmt.Sscanf(PackedX2EnbConfigurationUpdateFailure, "%x", &payload)
	var msgSrc unsafe.Pointer

	mBuf := rmrCgo.NewMBuf(rmrCgo.RIC_ENB_CONFIGURATION_UPDATE_FAILURE, len(payload), ranName, &payload, &xAction, msgSrc)
	notificationRequest := models.NotificationRequest{RanName: mBuf.Meid, Len: 0, Payload: []byte{0},
		StartTime: time.Now(), TransactionId:xAction}
	rmrMessengerMock.On("SendMsg", mBuf, true).Return(&rmrCgo.MBuf{}, fmt.Errorf("send failure"))
	h.Handle(&notificationRequest)
	rmrMessengerMock.AssertCalled(t, "SendMsg", mBuf, true)
}
