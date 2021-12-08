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
//

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package managers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
	"unsafe"
)

func initE2TKeepAliveTest(t *testing.T) (*mocks.RmrMessengerMock, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.E2TShutdownManagerMock, *E2TKeepAliveWorker) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3, KeepAliveResponseTimeoutMs: 400, KeepAliveDelayMs: 100}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	e2tShutdownManagerMock := &mocks.E2TShutdownManagerMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	e2tInstancesManager := NewE2TInstancesManager(rnibDataService, logger)

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, logger)

	e2tKeepAliveWorker := NewE2TKeepAliveWorker(logger, rmrSender, e2tInstancesManager, e2tShutdownManagerMock, config)

	return rmrMessengerMock, readerMock, writerMock, e2tShutdownManagerMock, &e2tKeepAliveWorker
}

func TestSendKeepAliveRequest(t *testing.T) {
	rmrMessengerMock, _, _, _, e2tKeepAliveWorker := initE2TKeepAliveTest(t)

	rmrMessengerMock.On("SendMsg", mock.Anything, false).Return(&rmrCgo.MBuf{}, nil)

	e2tKeepAliveWorker.SendKeepAliveRequest()

	var payload, xAction []byte
	var msgSrc unsafe.Pointer
	req := rmrCgo.NewMBuf(rmrCgo.E2_TERM_KEEP_ALIVE_REQ, 0, "", &payload, &xAction, msgSrc)

	rmrMessengerMock.AssertCalled(t, "SendMsg", req, false)
}

func TestShutdownExpiredE2T_InternalError(t *testing.T) {
	rmrMessengerMock, readerMock, _, _, e2tKeepAliveWorker := initE2TKeepAliveTest(t)

	readerMock.On("GetE2TAddresses").Return([]string{}, common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")))

	e2tKeepAliveWorker.E2TKeepAliveExpired()

	rmrMessengerMock.AssertNotCalled(t, "Shutdown")
}

func TestShutdownExpiredE2T_NoAddresses(t *testing.T) {
	rmrMessengerMock, readerMock, _, _, e2tKeepAliveWorker := initE2TKeepAliveTest(t)

	addresses := []string{}

	readerMock.On("GetE2TAddresses").Return(addresses, nil)

	e2tKeepAliveWorker.E2TKeepAliveExpired()

	rmrMessengerMock.AssertNotCalled(t, "Shutdown")
}

func TestShutdownExpiredE2T_NotExpired_InternalError(t *testing.T) {
	rmrMessengerMock, readerMock, _, _, e2tKeepAliveWorker := initE2TKeepAliveTest(t)

	addresses := []string{E2TAddress,E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.AssociatedRanList = []string{"test1","test2","test3"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.AssociatedRanList = []string{"test4","test5","test6", "test7"}

	readerMock.On("GetE2TAddresses").Return(addresses, nil)
	readerMock.On("GetE2TInstances",addresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")))

	e2tKeepAliveWorker.E2TKeepAliveExpired()

	rmrMessengerMock.AssertNotCalled(t, "Shutdown")
}

func TestShutdownExpiredE2T_NoE2T(t *testing.T) {
	rmrMessengerMock, readerMock, _, _, e2tKeepAliveWorker := initE2TKeepAliveTest(t)

	readerMock.On("GetE2TAddresses").Return([]string{}, common.NewResourceNotFoundError("not found"))

	e2tKeepAliveWorker.E2TKeepAliveExpired()

	rmrMessengerMock.AssertNotCalled(t, "Shutdown")
}

func TestShutdownExpiredE2T_NotExpired(t *testing.T) {
	rmrMessengerMock, readerMock, _, _, e2tKeepAliveWorker := initE2TKeepAliveTest(t)

	addresses := []string{E2TAddress,E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.AssociatedRanList = []string{"test1","test2","test3"}
	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.AssociatedRanList = []string{"test4","test5","test6", "test7"}

	readerMock.On("GetE2TAddresses").Return(addresses, nil)
	readerMock.On("GetE2TInstances",addresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)

	e2tKeepAliveWorker.E2TKeepAliveExpired()

	rmrMessengerMock.AssertNotCalled(t, "Shutdown")
}

func TestShutdownExpiredE2T_One_E2TExpired(t *testing.T) {
	_, readerMock, _, e2tShutdownManagerMock, e2tKeepAliveWorker := initE2TKeepAliveTest(t)

	addresses := []string{E2TAddress,E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.AssociatedRanList = []string{"test1","test2","test3"}

	time.Sleep(time.Duration(400) * time.Millisecond)

	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.AssociatedRanList = []string{"test4","test5","test6", "test7"}

	readerMock.On("GetE2TAddresses").Return(addresses, nil)
	readerMock.On("GetE2TInstances",addresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	e2tShutdownManagerMock.On("Shutdown", e2tInstance1).Return(nil)

	e2tKeepAliveWorker.E2TKeepAliveExpired()

	e2tShutdownManagerMock.AssertNumberOfCalls(t, "Shutdown", 1)
}

func TestShutdownExpiredE2T_Two_E2TExpired(t *testing.T) {
	_, readerMock, _, e2tShutdownManagerMock, e2tKeepAliveWorker := initE2TKeepAliveTest(t)

	addresses := []string{E2TAddress,E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.AssociatedRanList = []string{"test1","test2","test3"}

	e2tInstance2 := entities.NewE2TInstance(E2TAddress2, PodName)
	e2tInstance2.AssociatedRanList = []string{"test4","test5","test6", "test7"}

	time.Sleep(time.Duration(400) * time.Millisecond)

	readerMock.On("GetE2TAddresses").Return(addresses, nil)
	readerMock.On("GetE2TInstances",addresses).Return([]*entities.E2TInstance{e2tInstance1, e2tInstance2}, nil)
	e2tShutdownManagerMock.On("Shutdown", e2tInstance1).Return(nil)
	e2tShutdownManagerMock.On("Shutdown", e2tInstance2).Return(nil)

	e2tKeepAliveWorker.E2TKeepAliveExpired()

	e2tShutdownManagerMock.AssertNumberOfCalls(t, "Shutdown", 2)
}

func TestExecute_Two_E2TExpired(t *testing.T) {
	rmrMessengerMock, readerMock, _, e2tShutdownManagerMock, e2tKeepAliveWorker := initE2TKeepAliveTest(t)

	addresses := []string{E2TAddress,E2TAddress2}
	e2tInstance1 := entities.NewE2TInstance(E2TAddress, PodName)
	e2tInstance1.AssociatedRanList = []string{"test1","test2","test3"}

	readerMock.On("GetE2TAddresses").Return(addresses, nil)
	readerMock.On("GetE2TInstances",addresses).Return([]*entities.E2TInstance{e2tInstance1}, nil)
	e2tShutdownManagerMock.On("Shutdown", e2tInstance1).Return(nil)
	rmrMessengerMock.On("SendMsg", mock.Anything, false).Return(&rmrCgo.MBuf{}, nil)

	go e2tKeepAliveWorker.Execute()

	time.Sleep(time.Duration(500) * time.Millisecond)

	var payload, xAction []byte
	var msgSrc unsafe.Pointer
	req := rmrCgo.NewMBuf(rmrCgo.E2_TERM_KEEP_ALIVE_REQ, 0, "", &payload, &xAction, msgSrc)

	rmrMessengerMock.AssertCalled(t, "SendMsg", req, false)
	e2tShutdownManagerMock.AssertCalled(t, "Shutdown", e2tInstance1)
}