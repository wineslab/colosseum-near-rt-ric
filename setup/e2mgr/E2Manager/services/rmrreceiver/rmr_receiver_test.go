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


package rmrreceiver

import (
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/managers/notificationmanager"
	"e2mgr/mocks"
	"e2mgr/providers/rmrmsghandlerprovider"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"fmt"
	"testing"
	"time"
)

func TestListenAndHandle(t *testing.T) {
	log, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#rmr_service_test.TestListenAndHandle - failed to initialize logger, error: %s", err)
	}
	rmrReceiver := initRmrReceiver(log)
	go rmrReceiver.ListenAndHandle()
	time.Sleep(time.Microsecond * 10)
}

func initRmrMessenger(log *logger.Logger) rmrCgo.RmrMessenger {
	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)

	// TODO: that's not good since we don't actually test anything. if the error is populated then the loop will just continue and it's sort of a "workaround" for that method to be called
	var buf *rmrCgo.MBuf
	e := fmt.Errorf("test error")
	rmrMessengerMock.On("RecvMsg").Return(buf, e)
	return rmrMessenger
}

func initRmrReceiver(logger *logger.Logger) *RmrReceiver {
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}
	httpClient := &mocks.HttpClientMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	rmrMessenger := initRmrMessenger(logger)
	rmrSender := rmrsender.NewRmrSender(logger, rmrMessenger)
	ranSetupManager := managers.NewRanSetupManager(logger, rmrSender, rnibDataService)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, logger)
	routingManagerClient := clients.NewRoutingManagerClient(logger, config, httpClient)
	e2tAssociationManager := managers.NewE2TAssociationManager(logger, rnibDataService, e2tInstancesManager, routingManagerClient)
	rmrNotificationHandlerProvider := rmrmsghandlerprovider.NewNotificationHandlerProvider()
	rmrNotificationHandlerProvider.Init(logger, config, rnibDataService, rmrSender, ranSetupManager, e2tInstancesManager, routingManagerClient, e2tAssociationManager)
	notificationManager := notificationmanager.NewNotificationManager(logger, rmrNotificationHandlerProvider)
	return NewRmrReceiver(logger, rmrMessenger, notificationManager)
}
