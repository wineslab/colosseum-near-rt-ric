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

package clients

import (
	"bytes"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/models"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

const E2TAddress = "10.0.2.15:38000"
const E2TAddress2 = "10.0.2.15:38001"
const RanName = "test1"

func initRoutingManagerClientTest(t *testing.T) (*RoutingManagerClient, *mocks.HttpClientMock, *configuration.Configuration) {
	logger := initLog(t)
	config := &configuration.Configuration{}
	config.RoutingManager.BaseUrl = "http://iltlv740.intl.att.com:8080/ric/v1/handles/"
	httpClientMock := &mocks.HttpClientMock{}
	rmClient := NewRoutingManagerClient(logger, config, httpClientMock)
	return rmClient, httpClientMock, config
}

func TestDeleteE2TInstanceSuccess(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test1"}, nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + "e2t"
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", url, "application/json", body).Return(&http.Response{StatusCode: http.StatusOK, Body: respBody}, nil)
	err := rmClient.DeleteE2TInstance(E2TAddress, []string{"test1"})
	assert.Nil(t, err)
}

func TestDeleteE2TInstanceFailure(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test1"},nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + "e2t"
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Delete", url, "application/json", body).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: respBody}, nil)
	err := rmClient.DeleteE2TInstance(E2TAddress, []string{"test1"})
	assert.IsType(t, &e2managererrors.RoutingManagerError{}, err)
}

func TestDeleteE2TInstanceDeleteFailure(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.NewRoutingManagerDeleteRequestModel(E2TAddress, []string{"test1"},nil)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + "e2t"
	httpClientMock.On("Delete", url, "application/json", body).Return(&http.Response{}, errors.New("error"))
	err := rmClient.DeleteE2TInstance(E2TAddress, []string{"test1"})
	assert.IsType(t, &e2managererrors.RoutingManagerError{}, err)
}

func TestAddE2TInstanceSuccess(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.NewRoutingManagerE2TData(E2TAddress)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + "e2t"
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)
	err := rmClient.AddE2TInstance(E2TAddress)
	assert.Nil(t, err)
}

func TestAddE2TInstanceHttpPostFailure(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.NewRoutingManagerE2TData(E2TAddress)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + "e2t"
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{}, errors.New("error"))
	err := rmClient.AddE2TInstance(E2TAddress)
	assert.IsType(t, &e2managererrors.RoutingManagerError{}, err)
}

func TestAddE2TInstanceFailure(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.NewRoutingManagerE2TData(E2TAddress)
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + "e2t"
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: respBody}, nil)
	err := rmClient.AddE2TInstance(E2TAddress)
	assert.NotNil(t, err)
}

func TestAssociateRanToE2TInstance_Success(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress, RanName)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + AssociateRanToE2TInstanceApiSuffix
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)
	err := rmClient.AssociateRanToE2TInstance(E2TAddress, RanName)
	assert.Nil(t, err)
}

func TestAssociateRanToE2TInstance_RoutingManagerError(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress, RanName)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + AssociateRanToE2TInstanceApiSuffix
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{}, errors.New("error"))
	err := rmClient.AssociateRanToE2TInstance(E2TAddress, RanName)
	assert.IsType(t, &e2managererrors.RoutingManagerError{}, err)
}

func TestAssociateRanToE2TInstance_RoutingManager_400(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress, RanName)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + AssociateRanToE2TInstanceApiSuffix
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: respBody}, nil)
	err := rmClient.AssociateRanToE2TInstance(E2TAddress, RanName)
	assert.IsType(t, &e2managererrors.RoutingManagerError{}, err)
}

func TestDissociateRanE2TInstance_Success(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress, RanName)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + DissociateRanE2TInstanceApiSuffix
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{StatusCode: http.StatusCreated, Body: respBody}, nil)
	err := rmClient.DissociateRanE2TInstance(E2TAddress, RanName)
	assert.Nil(t, err)
}

func TestDissociateRanE2TInstance_RoutingManagerError(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress, RanName)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + DissociateRanE2TInstanceApiSuffix
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{}, errors.New("error"))
	err := rmClient.DissociateRanE2TInstance(E2TAddress, RanName)
	assert.IsType(t, &e2managererrors.RoutingManagerError{}, err)
}

func TestDissociateRanE2TInstance_RoutingManager_400(t *testing.T) {
	rmClient, httpClientMock, config := initRoutingManagerClientTest(t)

	data := models.RoutingManagerE2TDataList{models.NewRoutingManagerE2TData(E2TAddress, RanName)}
	marshaled, _ := json.Marshal(data)
	body := bytes.NewBuffer(marshaled)
	url := config.RoutingManager.BaseUrl + DissociateRanE2TInstanceApiSuffix
	respBody := ioutil.NopCloser(bytes.NewBufferString(""))
	httpClientMock.On("Post", url, "application/json", body).Return(&http.Response{StatusCode: http.StatusBadRequest, Body: respBody}, nil)
	err := rmClient.DissociateRanE2TInstance(E2TAddress, RanName)
	assert.IsType(t, &e2managererrors.RoutingManagerError{}, err)
}

// TODO: extract to test_utils
func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#delete_all_request_handler_test.TestHandleSuccessFlow - failed to initialize logger, error: %s", err)
	}
	return log
}

//func TestAddE2TInstanceInteg(t *testing.T) {
//	logger := initLog(t)
//	config := configuration.ParseConfiguration()
//	httpClient := &http.Client{}
//	rmClient := NewRoutingManagerClient(logger, config, httpClient)
//	err := rmClient.AddE2TInstance(E2TAddress)
//	assert.Nil(t, err)
//}
