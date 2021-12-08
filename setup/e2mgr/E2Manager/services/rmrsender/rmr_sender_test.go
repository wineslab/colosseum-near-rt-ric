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


package rmrsender

import (
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func initRmrSenderTest(t *testing.T) (*logger.Logger, *mocks.RmrMessengerMock) {
	log := initLog(t)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessengerMock.On("IsReady").Return(true)
	rmrMessengerMock.On("Close").Return()
	return log, rmrMessengerMock
}

//func TestRmrSender_CloseContext(t *testing.T) {
//	logger, rmrMessengerMock := initRmrSenderTest(t)
//
//	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
//	rmrSender := NewRmrSender(logger, &rmrMessenger)
//
//	rmrSender.CloseContext()
//	time.Sleep(time.Microsecond * 10)
//}

func TestRmrSenderSendSuccess(t *testing.T) {
	logger, rmrMessengerMock := initRmrSenderTest(t)

	ranName := "test"
	payload := []byte("some payload")
	var xAction []byte
	var msgSrc unsafe.Pointer
	mbuf := rmrCgo.NewMBuf(123, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", mbuf, true).Return(&rmrCgo.MBuf{}, nil)
	rmrMsg := models.NewRmrMessage(123, ranName, payload, xAction, nil)
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrSender := NewRmrSender(logger, rmrMessenger)
	err := rmrSender.Send(rmrMsg)
	assert.Nil(t, err)
	rmrMessengerMock.AssertCalled(t, "SendMsg",mbuf, true)

}

func TestRmrSenderSendFailure(t *testing.T) {
	logger, rmrMessengerMock := initRmrSenderTest(t)

	ranName := "test"
	payload := []byte("some payload")
	var xAction []byte
	var msgSrc unsafe.Pointer
	mbuf := rmrCgo.NewMBuf(123, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", mbuf, true).Return(mbuf, fmt.Errorf("rmr send failure"))
	rmrMsg := models.NewRmrMessage(123, ranName, payload, xAction, nil)
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrSender := NewRmrSender(logger, rmrMessenger)
	err := rmrSender.Send(rmrMsg)
	rmrMessengerMock.AssertCalled(t, "SendMsg",mbuf, true)
	assert.NotNil(t, err)
}

// TODO: extract to test_utils
func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Fatalf("#initLog - failed to initialize logger, error: %s", err)
	}
	return log
}
