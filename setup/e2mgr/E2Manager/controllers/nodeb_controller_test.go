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
	"bytes"
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/providers/httpmsghandlerprovider"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/tests"
	"encoding/json"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unsafe"
)

const (
	RanName                      = "test"
	AssociatedE2TInstanceAddress = "10.0.2.15:38000"
	ValidationFailureJson        = "{\"errorCode\":402,\"errorMessage\":\"Validation error\"}"
	ResourceNotFoundJson         = "{\"errorCode\":404,\"errorMessage\":\"Resource not found\"}"
	RnibErrorJson                = "{\"errorCode\":500,\"errorMessage\":\"RNIB error\"}"
	InternalErrorJson            = "{\"errorCode\":501,\"errorMessage\":\"Internal Server Error. Please try again later\"}"
)

var (
	ServedNrCellInformationRequiredFields = []string{"cellId", "choiceNrMode", "nrMode", "nrPci", "servedPlmns"}
	NrNeighbourInformationRequiredFields  = []string{"nrCgi", "choiceNrMode", "nrMode", "nrPci"}
)

type controllerGetNodebTestContext struct {
	ranName              string
	nodebInfo            *entities.NodebInfo
	rnibError            error
	expectedStatusCode   int
	expectedJsonResponse string
}

type controllerGetNodebIdListTestContext struct {
	nodebIdList          []*entities.NbIdentity
	rnibError            error
	expectedStatusCode   int
	expectedJsonResponse string
}

type getNodebInfoResult struct {
	nodebInfo *entities.NodebInfo
	rnibError error
}

type updateGnbCellsParams struct {
	err error
}

type removeServedNrCellsParams struct {
	servedNrCells []*entities.ServedNRCell
	err           error
}

type controllerUpdateGnbTestContext struct {
	getNodebInfoResult        *getNodebInfoResult
	removeServedNrCellsParams *removeServedNrCellsParams
	updateGnbCellsParams      *updateGnbCellsParams
	requestBody               map[string]interface{}
	expectedStatusCode        int
	expectedJsonResponse      string
}

func generateServedNrCells(cellIds ...string) []*entities.ServedNRCell {

	servedNrCells := []*entities.ServedNRCell{}

	for _, v := range cellIds {
		servedNrCells = append(servedNrCells, &entities.ServedNRCell{ServedNrCellInformation: &entities.ServedNRCellInformation{
			CellId: v,
			ChoiceNrMode: &entities.ServedNRCellInformation_ChoiceNRMode{
				Fdd: &entities.ServedNRCellInformation_ChoiceNRMode_FddInfo{

				},
			},
			NrMode:      entities.Nr_FDD,
			NrPci:       5,
			ServedPlmns: []string{"whatever"},
		}})
	}

	return servedNrCells
}

func buildNrNeighbourInformation(propToOmit string) map[string]interface{} {
	ret := map[string]interface{}{
		"nrCgi": "whatever",
		"choiceNrMode": map[string]interface{}{
			"tdd": map[string]interface{}{},
		},
		"nrMode": 1,
		"nrPci":  1,
	}

	if len(propToOmit) != 0 {
		delete(ret, propToOmit)
	}

	return ret
}

func buildServedNrCellInformation(propToOmit string) map[string]interface{} {
	ret := map[string]interface{}{
		"cellId": "whatever",
		"choiceNrMode": map[string]interface{}{
			"fdd": map[string]interface{}{},
		},
		"nrMode": 1,
		"nrPci":  1,
		"servedPlmns": []interface{}{
			"whatever",
		},
	}

	if len(propToOmit) != 0 {
		delete(ret, propToOmit)
	}

	return ret
}

func setupControllerTest(t *testing.T) (*NodebController, *mocks.RnibReaderMock, *mocks.RnibWriterMock, *mocks.RmrMessengerMock, *mocks.E2TInstancesManagerMock) {
	log := initLog(t)
	config := configuration.ParseConfiguration()

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	readerMock := &mocks.RnibReaderMock{}

	writerMock := &mocks.RnibWriterMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, writerMock)
	rmrSender := getRmrSender(rmrMessengerMock, log)
	ranSetupManager := managers.NewRanSetupManager(log, rmrSender, rnibDataService)
	e2tInstancesManager := &mocks.E2TInstancesManagerMock{}
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := clients.NewRoutingManagerClient(log, config, httpClientMock)
	e2tAssociationManager := managers.NewE2TAssociationManager(log, rnibDataService, e2tInstancesManager, rmClient)
	handlerProvider := httpmsghandlerprovider.NewIncomingRequestHandlerProvider(log, rmrSender, config, rnibDataService, ranSetupManager, e2tInstancesManager, e2tAssociationManager, rmClient)
	controller := NewNodebController(log, handlerProvider)
	return controller, readerMock, writerMock, rmrMessengerMock, e2tInstancesManager
}

func TestX2SetupInvalidBody(t *testing.T) {

	controller, _, _, _, _ := setupControllerTest(t)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	httpRequest, _ := http.NewRequest("POST", "http://localhost:3800/v1/nodeb/x2-setup", strings.NewReader("{}{}"))
	httpRequest.Header = header

	writer := httptest.NewRecorder()
	controller.X2Setup(writer, httpRequest)

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)
	assert.Equal(t, e2managererrors.NewInvalidJsonError().Code, errorResponse.Code)
}

func TestX2SetupSuccess(t *testing.T) {

	controller, readerMock, writerMock, rmrMessengerMock, _ := setupControllerTest(t)

	ranName := "test"
	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, AssociatedE2TInstanceAddress: "10.0.2.15:8989"}
	readerMock.On("GetNodeb", ranName).Return(nb, nil)
	var nbUpdated = *nb
	writerMock.On("UpdateNodebInfo", &nbUpdated).Return(nil)

	var nbUpdated2 = &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_X2_SETUP_REQUEST, AssociatedE2TInstanceAddress: "10.0.2.15:8989"}
	writerMock.On("UpdateNodebInfo", nbUpdated2).Return(nil)

	payload := e2pdus.PackedX2setupRequest
	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_SETUP_REQ, len(payload), ranName, &payload, &xAction, msgSrc)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	httpRequest := tests.GetHttpRequest()
	httpRequest.Header = header

	writer := httptest.NewRecorder()
	controller.X2Setup(writer, httpRequest)

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestEndcSetupSuccess(t *testing.T) {

	controller, readerMock, writerMock, rmrMessengerMock, _ := setupControllerTest(t)

	ranName := "test"
	nb := &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_DISCONNECTED, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, AssociatedE2TInstanceAddress: "10.0.2.15:8989"}
	readerMock.On("GetNodeb", ranName).Return(nb, nil)
	var nbUpdated = *nb
	writerMock.On("UpdateNodebInfo", &nbUpdated).Return(nil)

	var nbUpdated2 = &entities.NodebInfo{RanName: ranName, ConnectionStatus: entities.ConnectionStatus_CONNECTING, E2ApplicationProtocol: entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST, AssociatedE2TInstanceAddress: "10.0.2.15:8989"}
	writerMock.On("UpdateNodebInfo", nbUpdated2).Return(nil)

	payload := e2pdus.PackedEndcX2setupRequest
	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_ENDC_X2_SETUP_REQ, len(payload), ranName, &payload, &xAction, msgSrc)

	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(msg, nil)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	httpRequest := tests.GetHttpRequest()
	httpRequest.Header = header

	writer := httptest.NewRecorder()
	controller.EndcSetup(writer, httpRequest)

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestShutdownHandlerRnibError(t *testing.T) {
	controller, _, _, _, e2tInstancesManagerMock := setupControllerTest(t)
	e2tInstancesManagerMock.On("GetE2TAddresses").Return([]string{}, e2managererrors.NewRnibDbError())

	writer := httptest.NewRecorder()

	controller.Shutdown(writer, tests.GetHttpRequest())

	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, e2managererrors.NewRnibDbError().Code)
}

func controllerGetNodebTestExecuter(t *testing.T, context *controllerGetNodebTestContext) {
	controller, readerMock, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()
	readerMock.On("GetNodeb", context.ranName).Return(context.nodebInfo, context.rnibError)
	req, _ := http.NewRequest(http.MethodGet, "/nodeb", nil)
	req = mux.SetURLVars(req, map[string]string{"ranName": context.ranName})
	controller.GetNodeb(writer, req)
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, context.expectedJsonResponse, string(bodyBytes))
}

func controllerGetNodebIdListTestExecuter(t *testing.T, context *controllerGetNodebIdListTestContext) {
	controller, readerMock, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()
	readerMock.On("GetListNodebIds").Return(context.nodebIdList, context.rnibError)
	req, _ := http.NewRequest(http.MethodGet, "/nodeb/ids", nil)
	controller.GetNodebIdList(writer, req)
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, context.expectedJsonResponse, string(bodyBytes))
}

func activateControllerUpdateGnbMocks(context *controllerUpdateGnbTestContext, readerMock *mocks.RnibReaderMock, writerMock *mocks.RnibWriterMock) {
	if context.getNodebInfoResult != nil {
		readerMock.On("GetNodeb", RanName).Return(context.getNodebInfoResult.nodebInfo, context.getNodebInfoResult.rnibError)
	}

	if context.removeServedNrCellsParams != nil {
		writerMock.On("RemoveServedNrCells", RanName, context.removeServedNrCellsParams.servedNrCells).Return(context.removeServedNrCellsParams.err)
	}

	if context.updateGnbCellsParams != nil {
		updatedNodebInfo := *context.getNodebInfoResult.nodebInfo
		gnb := entities.Gnb{}
		_ = jsonpb.Unmarshal(getJsonRequestAsBuffer(context.requestBody), &gnb)
		updatedGnb := *updatedNodebInfo.GetGnb()
		updatedGnb.ServedNrCells = gnb.ServedNrCells
		writerMock.On("UpdateGnbCells", &updatedNodebInfo, gnb.ServedNrCells).Return(context.updateGnbCellsParams.err)
	}
}

func assertControllerUpdateGnb(t *testing.T, context *controllerUpdateGnbTestContext, writer *httptest.ResponseRecorder, readerMock *mocks.RnibReaderMock, writerMock *mocks.RnibWriterMock) {
	assert.Equal(t, context.expectedStatusCode, writer.Result().StatusCode)
	bodyBytes, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, context.expectedJsonResponse, string(bodyBytes))
	readerMock.AssertExpectations(t)
	writerMock.AssertExpectations(t)

	if context.getNodebInfoResult != nil {
		readerMock.AssertNotCalled(t, "GetNodeb")
	}

	if context.updateGnbCellsParams != nil {
		writerMock.AssertNotCalled(t, "UpdateGnb")
	}

	if context.removeServedNrCellsParams != nil {
		writerMock.AssertNotCalled(t, "RemoveServedNrCells")
	}
}

func buildUpdateGnbRequest(context *controllerUpdateGnbTestContext) *http.Request {
	updateGnbUrl := fmt.Sprintf("/nodeb/%s/update", RanName)
	requestBody := getJsonRequestAsBuffer(context.requestBody)
	req, _ := http.NewRequest(http.MethodGet, updateGnbUrl, requestBody)
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"ranName": RanName})
	return req
}

func controllerUpdateGnbTestExecuter(t *testing.T, context *controllerUpdateGnbTestContext) {
	controller, readerMock, writerMock, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()

	activateControllerUpdateGnbMocks(context, readerMock, writerMock)
	req := buildUpdateGnbRequest(context)
	controller.UpdateGnb(writer, req)
	assertControllerUpdateGnb(t, context, writer, readerMock, writerMock)
}

func TestControllerUpdateGnbEmptyServedNrCells(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: nil,
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
			},
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbMissingServedNrCellInformation(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: nil,
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": nil,
				},
			},
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbMissingServedNrCellRequiredProp(t *testing.T) {

	for _, v := range ServedNrCellInformationRequiredFields {
		context := controllerUpdateGnbTestContext{
			getNodebInfoResult: nil,
			requestBody: map[string]interface{}{
				"servedNrCells": []interface{}{
					map[string]interface{}{
						"servedNrCellInformation": buildServedNrCellInformation(v),
					},
				},
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: ValidationFailureJson,
		}

		controllerUpdateGnbTestExecuter(t, &context)
	}
}

func TestControllerUpdateGnbMissingServedNrCellFddOrTdd(t *testing.T) {

	servedNrCellInformation := buildServedNrCellInformation("")
	servedNrCellInformation["choiceNrMode"] = map[string]interface{}{}

	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: nil,
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": servedNrCellInformation,
				},
			},
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbMissingNeighbourInfoFddOrTdd(t *testing.T) {

	nrNeighbourInfo := buildNrNeighbourInformation("")
	nrNeighbourInfo["choiceNrMode"] = map[string]interface{}{}

	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: nil,
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						nrNeighbourInfo,
					},
				},
			},
		},
		expectedStatusCode:   http.StatusBadRequest,
		expectedJsonResponse: ValidationFailureJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbMissingNrNeighbourInformationRequiredProp(t *testing.T) {

	for _, v := range NrNeighbourInformationRequiredFields {
		context := controllerUpdateGnbTestContext{
			getNodebInfoResult: nil,
			requestBody: map[string]interface{}{
				"servedNrCells": []interface{}{
					map[string]interface{}{
						"servedNrCellInformation": buildServedNrCellInformation(""),
						"nrNeighbourInfos": []interface{}{
							buildNrNeighbourInformation(v),
						},
					},
				},
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedJsonResponse: ValidationFailureJson,
		}

		controllerUpdateGnbTestExecuter(t, &context)
	}
}

func TestControllerUpdateGnbValidServedNrCellInformationGetNodebNotFound(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewResourceNotFoundError("#reader.GetNodeb - Not found Error"),
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
				},
			},
		},
		expectedStatusCode:   http.StatusNotFound,
		expectedJsonResponse: ResourceNotFoundJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbValidServedNrCellInformationGetNodebInternalError(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: nil,
			rnibError: common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")),
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
				},
			},
		},
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbGetNodebSuccessInvalidGnbConfiguration(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
			},
			rnibError: nil,
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						buildNrNeighbourInformation(""),
					},
				},
			},
		},
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: InternalErrorJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbGetNodebSuccessRemoveServedNrCellsFailure(t *testing.T) {
	oldServedNrCells := generateServedNrCells("whatever1", "whatever2")
	context := controllerUpdateGnbTestContext{
		removeServedNrCellsParams: &removeServedNrCellsParams{
			err:           common.NewInternalError(errors.New("#writer.UpdateGnbCells - Internal Error")),
			servedNrCells: oldServedNrCells,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{ServedNrCells: oldServedNrCells}},
			},
			rnibError: nil,
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						buildNrNeighbourInformation(""),
					},
				},
			},
		},
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbGetNodebSuccessUpdateGnbCellsFailure(t *testing.T) {
	oldServedNrCells := generateServedNrCells("whatever1", "whatever2")
	context := controllerUpdateGnbTestContext{
		removeServedNrCellsParams: &removeServedNrCellsParams{
			err:           nil,
			servedNrCells: oldServedNrCells,
		},
		updateGnbCellsParams: &updateGnbCellsParams{
			err: common.NewInternalError(errors.New("#writer.UpdateGnbCells - Internal Error")),
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{ServedNrCells: oldServedNrCells}},
			},
			rnibError: nil,
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						buildNrNeighbourInformation(""),
					},
				},
			},
		},
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func TestControllerUpdateGnbSuccess(t *testing.T) {
	context := controllerUpdateGnbTestContext{
		updateGnbCellsParams: &updateGnbCellsParams{
			err: nil,
		},
		getNodebInfoResult: &getNodebInfoResult{
			nodebInfo: &entities.NodebInfo{
				RanName:                      RanName,
				ConnectionStatus:             entities.ConnectionStatus_CONNECTED,
				AssociatedE2TInstanceAddress: AssociatedE2TInstanceAddress,
				Configuration:                &entities.NodebInfo_Gnb{Gnb: &entities.Gnb{}},
			},
			rnibError: nil,
		},
		requestBody: map[string]interface{}{
			"servedNrCells": []interface{}{
				map[string]interface{}{
					"servedNrCellInformation": buildServedNrCellInformation(""),
					"nrNeighbourInfos": []interface{}{
						buildNrNeighbourInformation(""),
					},
				},
			},
		},
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "{\"ranName\":\"test\",\"connectionStatus\":\"CONNECTED\",\"gnb\":{\"servedNrCells\":[{\"servedNrCellInformation\":{\"nrPci\":1,\"cellId\":\"whatever\",\"servedPlmns\":[\"whatever\"],\"nrMode\":\"FDD\",\"choiceNrMode\":{\"fdd\":{}}},\"nrNeighbourInfos\":[{\"nrPci\":1,\"nrCgi\":\"whatever\",\"nrMode\":\"FDD\",\"choiceNrMode\":{\"tdd\":{}}}]}]},\"associatedE2tInstanceAddress\":\"10.0.2.15:38000\"}",
	}

	controllerUpdateGnbTestExecuter(t, &context)
}

func getJsonRequestAsBuffer(requestJson map[string]interface{}) *bytes.Buffer {
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(requestJson)
	return b
}

func TestControllerGetNodebSuccess(t *testing.T) {
	ranName := "test"
	var rnibError error
	context := controllerGetNodebTestContext{
		ranName:              ranName,
		nodebInfo:            &entities.NodebInfo{RanName: ranName, Ip: "10.0.2.15", Port: 1234},
		rnibError:            rnibError,
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: fmt.Sprintf("{\"ranName\":\"%s\",\"ip\":\"10.0.2.15\",\"port\":1234}", ranName),
	}

	controllerGetNodebTestExecuter(t, &context)
}

func TestControllerGetNodebNotFound(t *testing.T) {

	ranName := "test"
	var nodebInfo *entities.NodebInfo
	context := controllerGetNodebTestContext{
		ranName:              ranName,
		nodebInfo:            nodebInfo,
		rnibError:            common.NewResourceNotFoundError("#reader.GetNodeb - Not found Error"),
		expectedStatusCode:   http.StatusNotFound,
		expectedJsonResponse: ResourceNotFoundJson,
	}

	controllerGetNodebTestExecuter(t, &context)
}

func TestControllerGetNodebInternal(t *testing.T) {
	ranName := "test"
	var nodebInfo *entities.NodebInfo
	context := controllerGetNodebTestContext{
		ranName:              ranName,
		nodebInfo:            nodebInfo,
		rnibError:            common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerGetNodebTestExecuter(t, &context)
}

func TestControllerGetNodebIdListSuccess(t *testing.T) {
	var rnibError error
	nodebIdList := []*entities.NbIdentity{
		{InventoryName: "test1", GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId1", NbId: "nbId1"}},
		{InventoryName: "test2", GlobalNbId: &entities.GlobalNbId{PlmnId: "plmnId2", NbId: "nbId2"}},
	}

	context := controllerGetNodebIdListTestContext{
		nodebIdList:          nodebIdList,
		rnibError:            rnibError,
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "[{\"inventoryName\":\"test1\",\"globalNbId\":{\"plmnId\":\"plmnId1\",\"nbId\":\"nbId1\"}},{\"inventoryName\":\"test2\",\"globalNbId\":{\"plmnId\":\"plmnId2\",\"nbId\":\"nbId2\"}}]",
	}

	controllerGetNodebIdListTestExecuter(t, &context)
}

func TestControllerGetNodebIdListEmptySuccess(t *testing.T) {
	var rnibError error
	nodebIdList := []*entities.NbIdentity{}

	context := controllerGetNodebIdListTestContext{
		nodebIdList:          nodebIdList,
		rnibError:            rnibError,
		expectedStatusCode:   http.StatusOK,
		expectedJsonResponse: "[]",
	}

	controllerGetNodebIdListTestExecuter(t, &context)
}

func TestControllerGetNodebIdListInternal(t *testing.T) {
	var nodebIdList []*entities.NbIdentity
	context := controllerGetNodebIdListTestContext{
		nodebIdList:          nodebIdList,
		rnibError:            common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")),
		expectedStatusCode:   http.StatusInternalServerError,
		expectedJsonResponse: RnibErrorJson,
	}

	controllerGetNodebIdListTestExecuter(t, &context)
}

func TestHeaderValidationFailed(t *testing.T) {
	controller, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()

	header := &http.Header{}

	controller.handleRequest(writer, header, httpmsghandlerprovider.ShutdownRequest, nil, true)

	var errorResponse = parseJsonRequest(t, writer.Body)
	err := e2managererrors.NewHeaderValidationError()

	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestShutdownStatusNoContent(t *testing.T) {
	controller, readerMock, _, _, e2tInstancesManagerMock := setupControllerTest(t)
	e2tInstancesManagerMock.On("GetE2TAddresses").Return([]string{}, nil)
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{}, nil)

	writer := httptest.NewRecorder()
	controller.Shutdown(writer, tests.GetHttpRequest())

	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)
}

func TestHandleInternalError(t *testing.T) {
	controller, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	err := e2managererrors.NewInternalError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestHandleCommandAlreadyInProgressError(t *testing.T) {
	controller, _, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()
	err := e2managererrors.NewCommandAlreadyInProgressError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusMethodNotAllowed, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestHandleRoutingManagerError(t *testing.T) {
	controller, _, _, _, _ := setupControllerTest(t)
	writer := httptest.NewRecorder()
	err := e2managererrors.NewRoutingManagerError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusServiceUnavailable, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestHandleE2TInstanceAbsenceError(t *testing.T) {
	controller, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	err := e2managererrors.NewE2TInstanceAbsenceError()

	controller.handleErrorResponse(err, writer)
	var errorResponse = parseJsonRequest(t, writer.Body)

	assert.Equal(t, http.StatusServiceUnavailable, writer.Result().StatusCode)
	assert.Equal(t, errorResponse.Code, err.Code)
	assert.Equal(t, errorResponse.Message, err.Message)
}

func TestValidateHeaders(t *testing.T) {
	controller, _, _, _, _ := setupControllerTest(t)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	result := controller.validateRequestHeader(&header)

	assert.Nil(t, result)
}

func parseJsonRequest(t *testing.T, r io.Reader) models.ErrorResponse {

	var errorResponse models.ErrorResponse
	body, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Error cannot deserialize json request")
	}
	_ = json.Unmarshal(body, &errorResponse)

	return errorResponse
}

func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}

func TestX2ResetHandleSuccessfulRequestedCause(t *testing.T) {
	controller, readerMock, _, rmrMessengerMock, _ := setupControllerTest(t)

	ranName := "test1"
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x40}
	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", msg, mock.Anything).Return(msg, nil)

	writer := httptest.NewRecorder()

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	data4Req := map[string]interface{}{"cause": "protocol:transfer-syntax-error"}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(data4Req)
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2Reset(writer, req)
	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)

}

func TestX2ResetHandleSuccessfulRequestedDefault(t *testing.T) {
	controller, readerMock, _, rmrMessengerMock, _ := setupControllerTest(t)

	ranName := "test1"
	// o&m intervention
	payload := []byte{0x00, 0x07, 0x00, 0x08, 0x00, 0x00, 0x01, 0x00, 0x05, 0x40, 0x01, 0x64}
	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := rmrCgo.NewMBuf(rmrCgo.RIC_X2_RESET, len(payload), ranName, &payload, &xAction, msgSrc)
	rmrMessengerMock.On("SendMsg", msg, true).Return(msg, nil)

	writer := httptest.NewRecorder()

	var nodeb = &entities.NodebInfo{ConnectionStatus: entities.ConnectionStatus_CONNECTED}
	readerMock.On("GetNodeb", ranName).Return(nodeb, nil)

	// no body
	b := new(bytes.Buffer)
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2Reset(writer, req)
	assert.Equal(t, http.StatusNoContent, writer.Result().StatusCode)

}

func TestX2ResetHandleFailureInvalidBody(t *testing.T) {
	controller, _, _, _, _ := setupControllerTest(t)

	ranName := "test1"

	writer := httptest.NewRecorder()

	// Invalid json: attribute name without quotes (should be "cause":).
	b := strings.NewReader("{cause:\"protocol:transfer-syntax-error\"")
	req, _ := http.NewRequest("PUT", "https://localhost:3800/nodeb-reset", b)
	req = mux.SetURLVars(req, map[string]string{"ranName": ranName})

	controller.X2Reset(writer, req)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

}

func TestHandleErrorResponse(t *testing.T) {
	controller, _, _, _, _ := setupControllerTest(t)

	writer := httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewRnibDbError(), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewCommandAlreadyInProgressError(), writer)
	assert.Equal(t, http.StatusMethodNotAllowed, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewHeaderValidationError(), writer)
	assert.Equal(t, http.StatusUnsupportedMediaType, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewWrongStateError("", ""), writer)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewRequestValidationError(), writer)
	assert.Equal(t, http.StatusBadRequest, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewRmrError(), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(e2managererrors.NewResourceNotFoundError(), writer)
	assert.Equal(t, http.StatusNotFound, writer.Result().StatusCode)

	writer = httptest.NewRecorder()
	controller.handleErrorResponse(fmt.Errorf("ErrorError"), writer)
	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
}

func getRmrSender(rmrMessengerMock *mocks.RmrMessengerMock, log *logger.Logger) *rmrsender.RmrSender {
	rmrMessenger := rmrCgo.RmrMessenger(rmrMessengerMock)
	rmrMessengerMock.On("Init", tests.GetPort(), tests.MaxMsgSize, tests.Flags, log).Return(&rmrMessenger)
	return rmrsender.NewRmrSender(log, rmrMessenger)
}
