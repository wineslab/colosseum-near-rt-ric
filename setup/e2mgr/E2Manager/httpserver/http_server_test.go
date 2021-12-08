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


package httpserver

import (
	"e2mgr/logger"
	"e2mgr/mocks"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupRouterAndMocks() (*mux.Router, *mocks.RootControllerMock, *mocks.NodebControllerMock, *mocks.E2TControllerMock) {
	rootControllerMock := &mocks.RootControllerMock{}
	rootControllerMock.On("HandleHealthCheckRequest").Return(nil)

	nodebControllerMock := &mocks.NodebControllerMock{}
	nodebControllerMock.On("Shutdown").Return(nil)
	nodebControllerMock.On("GetNodeb").Return(nil)
	nodebControllerMock.On("GetNodebIdList").Return(nil)

	e2tControllerMock := &mocks.E2TControllerMock{}

	e2tControllerMock.On("GetE2TInstances").Return(nil)

	router := mux.NewRouter()
	initializeRoutes(router, rootControllerMock, nodebControllerMock, e2tControllerMock)
	return router, rootControllerMock, nodebControllerMock, e2tControllerMock
}

func TestRouteGetNodebIds(t *testing.T) {
	router, _, nodebControllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/nodeb/ids", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	nodebControllerMock.AssertNumberOfCalls(t, "GetNodebIdList", 1)
}

func TestRouteGetNodebRanName(t *testing.T) {
	router, _, nodebControllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/nodeb/ran1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.Equal(t, "ran1", rr.Body.String(), "handler returned wrong body")
	nodebControllerMock.AssertNumberOfCalls(t, "GetNodeb", 1)
}

func TestRouteGetHealth(t *testing.T) {
	router, rootControllerMock, _, _ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	rootControllerMock.AssertNumberOfCalls(t, "HandleHealthCheckRequest", 1)
}

func TestRoutePutNodebShutdown(t *testing.T) {
	router, _, nodebControllerMock, _ := setupRouterAndMocks()

	req, err := http.NewRequest("PUT", "/v1/nodeb/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	nodebControllerMock.AssertNumberOfCalls(t, "Shutdown", 1)
}

func TestRouteNotFound(t *testing.T) {
	router, _, _,_ := setupRouterAndMocks()

	req, err := http.NewRequest("GET", "/v1/no/such/route", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, "handler returned wrong status code")
}

func TestRunError(t *testing.T) {
	log := initLog(t)
	err := Run(log, 1234567, &mocks.RootControllerMock{}, &mocks.NodebControllerMock{}, &mocks.E2TControllerMock{})
	assert.NotNil(t, err)
}

func TestRun(t *testing.T) {
	log := initLog(t)
	_, rootControllerMock, nodebControllerMock, e2tControllerMock := setupRouterAndMocks()
	go Run(log, 11223, rootControllerMock, nodebControllerMock, e2tControllerMock)

	time.Sleep(time.Millisecond * 100)
	resp, err := http.Get("http://localhost:11223/v1/health")
	if err != nil {
		t.Fatalf("failed to perform GET to http://localhost:11223/v1/health")
	}
	assert.Equal(t, 200, resp.StatusCode)
}

func initLog(t *testing.T) *logger.Logger {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#initLog test - failed to initialize logger, error: %s", err)
	}
	return log
}
