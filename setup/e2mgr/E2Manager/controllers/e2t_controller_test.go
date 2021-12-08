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

package controllers

import (
	"e2mgr/configuration"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/providers/httpmsghandlerprovider"
	"e2mgr/services"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/magiconair/properties/assert"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const E2TAddress string = "10.0.2.15:38000"
const E2TAddress2 string = "10.0.2.16:38001"

type controllerE2TInstancesTestContext struct {
	e2tAddresses         []string
	e2tInstances         []*entities.E2TInstance
	error                error
	expectedStatusCode   int
	expectedJsonResponse string
}

func setupE2TControllerTest(t *testing.T) (*E2TController, *mocks.RnibReaderMock) {
	log := initLog(t)
	config := configuration.ParseConfiguration()

	readerMock := &mocks.RnibReaderMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, nil)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, log)
	handlerProvider := httpmsghandlerprovider.NewIncomingRequestHandlerProvider(log, nil, config, rnibDataService, nil, e2tInstancesManager, &managers.E2TAssociationManager{}, nil)
	controller := NewE2TController(log, handlerProvider)
	return controller, readerMock
}

func controllerGetE2TInstancesTestExecuter(t *testing.T, context *controllerE2TInstancesTestContext) {
	controller, readerMock := setupE2TControllerTest(t)
	writer := httptest.NewRecorder()
	readerMock.On("GetE2TAddresses").Return(context.e2tAddresses, context.error)

	if context.e2tInstances != nil {
		readerMock.On("GetE2TInstances", context.e2tAddresses).Return(context.e2tInstances, context.error)
	}

	req, _ := http.NewRequest("GET", "/e2t/list", nil)
	controller.GetE2TInstances(writer, req)
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, context.expectedJsonResponse, string(bodyBytes))
	readerMock.AssertExpectations(t)
}

func TestControllerGetE2TInstancesSuccess(t *testing.T) {
	ranNames1 := []string{"test1", "test2", "test3"}
	e2tInstanceResponseModel1 := models.NewE2TInstanceResponseModel(E2TAddress, ranNames1)
	e2tInstanceResponseModel2 := models.NewE2TInstanceResponseModel(E2TAddress2, []string{})
	e2tInstancesResponse := models.E2TInstancesResponse{e2tInstanceResponseModel1, e2tInstanceResponseModel2}
	bytes, _ := json.Marshal(e2tInstancesResponse)

	context := controllerE2TInstancesTestContext{
		e2tAddresses:         []string{E2TAddress, E2TAddress2},
		e2tInstances:         []*entities.E2TInstance{{Address: E2TAddress, AssociatedRanList: ranNames1}, {Address: E2TAddress2, AssociatedRanList: []string{}}},
		error:                nil,
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: string(bytes),
	}

	controllerGetE2TInstancesTestExecuter(t, &context)
}

func TestControllerGetE2TInstancesEmptySuccess(t *testing.T) {
	e2tInstancesResponse := models.E2TInstancesResponse{}
	bytes, _ := json.Marshal(e2tInstancesResponse)

	context := controllerE2TInstancesTestContext{
		e2tAddresses:         []string{},
		e2tInstances:         nil,
		error:                nil,
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: string(bytes),
	}

	controllerGetE2TInstancesTestExecuter(t, &context)
}

func TestControllerGetE2TInstancesInternal(t *testing.T) {
	context := controllerE2TInstancesTestContext{
		e2tAddresses:         nil,
		e2tInstances:         nil,
		error:                common.NewInternalError(errors.New("error")),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: "{\"errorCode\":500,\"errorMessage\":\"RNIB error\"}",
	}

	controllerGetE2TInstancesTestExecuter(t, &context)
}

func TestInvalidRequestName(t *testing.T) {
	controller, _ := setupE2TControllerTest(t)

	writer := httptest.NewRecorder()

	header := &http.Header{}

	controller.handleRequest(writer, header, "", nil, true)

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, 501)
}
