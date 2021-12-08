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
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"unsafe"
)

func initRanSetupManagerTest(t *testing.T) (*mocks.RmrMessengerMock, *mocks.RnibWriterMock, *RanSetupManager) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, logger)

	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	ranSetupManager := NewRanSetupManager(logger, rmrSender, rnibDataService)
	return rmrMessengerMock, writerMock, ranSetupManager
}

func TestExecuteSetupConnectingX2Setup(t *testing.T) {
	rmrMessengerMock, writerMock, mgr := initRanSetupManagerTest(t)

	ranName := "test1"

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var rnibErr error
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xAction := []byte(ranName)
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	if err := mgr.ExecuteSetup(initialNodeb, entities.ConnectionStatus_CONNECTING); err != nil {
		t.Errorf("want: success, got: error: %s", err)
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestExecuteSetupConnectingEndcX2Setup(t *testing.T) {
	rmrMessengerMock, writerMock, mgr := initRanSetupManagerTest(t)

	ranName := "test1"

		var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST}
	var rnibErr error
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedEndcX2setupRequest
	xAction := []byte(ranName)
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_ENDC_X2_SETUP_REQ, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	if err := mgr.ExecuteSetup(initialNodeb, entities.ConnectionStatus_CONNECTING); err != nil {
		t.Errorf("want: success, got: error: %s", err)
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestExecuteSetupDisconnected(t *testing.T) {
	rmrMessengerMock, writerMock, mgr := initRanSetupManagerTest(t)

	ranName := "test1"

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodebDisconnected = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var rnibErr error
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodebDisconnected).Return(rnibErr)

	payload := []byte{0}
	xAction := []byte(ranName)
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, fmt.Errorf("send failure"))

	if err := mgr.ExecuteSetup(initialNodeb, entities.ConnectionStatus_CONNECTING); err == nil {
		t.Errorf("want: failure, got: success")
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestExecuteSetupConnectingRnibError(t *testing.T) {
	rmrMessengerMock, writerMock, mgr := initRanSetupManagerTest(t)

	ranName := "test1"

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodebDisconnected = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var rnibErr = common.NewInternalError(fmt.Errorf("DB error"))
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodebDisconnected).Return(rnibErr)

	payload := []byte{0}
	xAction := []byte(ranName)
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, fmt.Errorf("send failure"))

	if err := mgr.ExecuteSetup(initialNodeb, entities.ConnectionStatus_CONNECTING); err == nil {
		t.Errorf("want: failure, got: success")
	} else {
		assert.IsType(t, e2managererrors.NewRnibDbError(), err)
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func TestExecuteSetupDisconnectedRnibError(t *testing.T) {
	rmrMessengerMock, writerMock, mgr := initRanSetupManagerTest(t)

	ranName := "test1"

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var argNodebDisconnected = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST}
	var rnibErr error
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)
	writerMock.On("UpdateNodebInfo", argNodebDisconnected).Return(common.NewInternalError(fmt.Errorf("DB error")))

	payload := []byte{0}
	xAction := []byte(ranName)
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, fmt.Errorf("send failure"))

	if err := mgr.ExecuteSetup(initialNodeb, entities.ConnectionStatus_CONNECTING); err == nil {
		t.Errorf("want: failure, got: success")
	} else {
		assert.IsType(t, e2managererrors.NewRnibDbError(), err)
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 2)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 1)
}

func TestExecuteSetupUnsupportedProtocol(t *testing.T) {
	rmrMessengerMock, writerMock, mgr := initRanSetupManagerTest(t)

	ranName := "test1"

	var initialNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_UNKNOWN_E2_APPLICATION_PROTOCOL}
	var argNodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_UNKNOWN_E2_APPLICATION_PROTOCOL}
	var rnibErr error
	writerMock.On("UpdateNodebInfo", argNodeb).Return(rnibErr)

	payload := e2pdus.PackedX2setupRequest
	xAction := []byte(ranName)
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	if err := mgr.ExecuteSetup(initialNodeb, entities.ConnectionStatus_CONNECTING); err == nil {
		t.Errorf("want: error, got: success")
	}

	writerMock.AssertNumberOfCalls(t, "UpdateNodebInfo", 1)
	rmrMessengerMock.AssertNumberOfCalls(t, "SendMsg", 0)
}

func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#initLog test - failed to initialize logger, error: %s", err)
	}
	return log
}
