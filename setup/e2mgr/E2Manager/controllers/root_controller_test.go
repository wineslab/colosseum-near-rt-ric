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


package controllers

import (
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/services"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupNodebControllerTest(t *testing.T) (services.RNibDataService, *mocks.RnibReaderMock){
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Errorf("#... - failed to initialize logger, error: %s", err)
	}
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(logger, config, readerMock, writerMock)
	return rnibDataService, readerMock
}

func TestNewRequestController(t *testing.T) {
	rnibDataService, _ := setupNodebControllerTest(t)
	assert.NotNil(t, NewRootController(rnibDataService))
}

func TestHandleHealthCheckRequestGood(t *testing.T) {
	rnibDataService, rnibReaderMock := setupNodebControllerTest(t)

	var nbList []*entities.NbIdentity
	rnibReaderMock.On("GetListNodebIds").Return(nbList, nil)

	rc := NewRootController(rnibDataService)
	writer := httptest.NewRecorder()
	rc.HandleHealthCheckRequest(writer, nil)
	assert.Equal(t, http.StatusOK, writer.Result().StatusCode)
}

func TestHandleHealthCheckRequestOtherError(t *testing.T) {
	rnibDataService, rnibReaderMock := setupNodebControllerTest(t)

	mockOtherErr := &common.InternalError{Err: fmt.Errorf("non connection error")}
	var nbList []*entities.NbIdentity
	rnibReaderMock.On("GetListNodebIds").Return(nbList, mockOtherErr)

	rc := NewRootController(rnibDataService)
	writer := httptest.NewRecorder()
	rc.HandleHealthCheckRequest(writer, nil)
	assert.Equal(t, http.StatusOK, writer.Result().StatusCode)
}

func TestHandleHealthCheckRequestConnError(t *testing.T) {
	rnibDataService, rnibReaderMock := setupNodebControllerTest(t)

	mockConnErr := &common.InternalError{Err: &net.OpError{Err: fmt.Errorf("connection error")}}
	var nbList []*entities.NbIdentity
	rnibReaderMock.On("GetListNodebIds").Return(nbList, mockConnErr)


	rc := NewRootController(rnibDataService)
	writer := httptest.NewRecorder()
	rc.HandleHealthCheckRequest(writer, nil)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
}