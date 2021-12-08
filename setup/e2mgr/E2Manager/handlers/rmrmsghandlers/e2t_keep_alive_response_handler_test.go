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

package rmrmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/services"
	"testing"
)

func initE2TKeepAliveTest(t *testing.T) (*logger.Logger, E2TKeepAliveResponseHandler, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.E2TInstancesManagerMock) {

	logger := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}

	readerMock := &mocks.RnibReaderMock{}
	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)

	e2tInstancesManagerMock := &mocks.E2TInstancesManagerMock{}
	handler := NewE2TKeepAliveResponseHandler(logger, rnibDataService, e2tInstancesManagerMock)
	return logger, handler, readerMock, writerMock, e2tInstancesManagerMock
}

func TestE2TKeepAliveUnmarshalPayloadFailure(t *testing.T) {
	_, handler, _, _, e2tInstancesManagerMock := initE2TKeepAliveTest(t)
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte("asd")}
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertNotCalled(t, "ResetKeepAliveTimestamp")
}

func TestE2TKeepAliveUnmarshalPayloadSuccess(t *testing.T) {
	_, handler, _, _, e2tInstancesManagerMock := initE2TKeepAliveTest(t)

	jsonRequest := "{\"address\":\"10.10.2.15:9800\"}"
	notificationRequest := &models.NotificationRequest{RanName: RanName, Payload: []byte(jsonRequest)}

	e2tInstancesManagerMock.On("ResetKeepAliveTimestamp", "10.10.2.15:9800").Return(nil)
	handler.Handle(notificationRequest)
	e2tInstancesManagerMock.AssertCalled(t, "ResetKeepAliveTimestamp", "10.10.2.15:9800")
}
