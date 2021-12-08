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
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func setupX2ResetRequestHandlerTest(t *testing.T) (*X2ResetRequestHandler, *mocks.RmrMessengerMock, *mocks.RnibReaderMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := getRmrSender(rmrMessengerMock, log)
	handler := NewX2ResetRequestHandler(log, rmrSender, rnibDataService)

	return handler, rmrMessengerMock, readerMock
}
func TestHandleSuccessfulDefaultCause(t *testing.T) {
	handler, rmrMessengerMock, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	// o&m intervention
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	var xAction[]byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xAction, msgSrc)

	rmrMessengerMock.On("SendMsg", msg, true).Return(msg, nil)

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.Nil(t, actual)
}

func TestHandleSuccessfulRequestedCause(t *testing.T) {
	handler, rmrMessengerMock, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x40}
	var xAction[]byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", msg, true).Return(msg, nil)

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName, Cause: "protocol:transfer-syntax-error"})

	assert.Nil(t, actual)
}

func TestHandleFailureUnknownCause(t *testing.T) {
	handler, _, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName, Cause: "XXX"})

	assert.IsType(t, e2managererrors.NewRequestValidationError(), actual)

}

func TestHandleFailureWrongState(t *testing.T) {
	handler, _, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.IsType(t, e2managererrors.NewWrongStateError(X2_RESET_ACTIVITY_NAME, entities.ConnectionStatus_name[int32(nodeb.ConnectionStatus)]), actual)
}

func TestHandleFailureRanNotFound(t *testing.T) {
	handler, _, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"

	readerMock.On("GetNodeb", ranName).Return(&entities.NodebInfo{}, common.NewResourceNotFoundError("nodeb not found"))

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.IsType(t, e2managererrors.NewResourceNotFoundError(), actual)
}

func TestHandleFailureRnibError(t *testing.T) {
	handler, _, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"

	readerMock.On("GetNodeb", ranName).Return(&entities.NodebInfo{}, common.NewInternalError(fmt.Errorf("internal error")))

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.IsType(t, e2managererrors.NewRnibDbError(), actual)
}

func TestHandleFailureRmrError(t *testing.T) {
	handler, rmrMessengerMock, readerMock := setupX2ResetRequestHandlerTest(t)

	ranName := "test1"
	// o&m intervention
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	var xAction[]byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", msg, true).Return(&rmrCgo.MBuf{}, fmt.Errorf("rmr error"))

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	_, actual := handler.Handle(models.ResetRequest{RanName: ranName})

	assert.IsType(t, e2managererrors.NewRmrError(), actual)
}
