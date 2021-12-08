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


package controllers

import (
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/providers/httpmsghandlerprovider"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"
)

type IE2TController interface {
	GetE2TInstances(writer http.ResponseWriter, r *http.Request)
}

type E2TController struct {
	logger          *logger.Logger
	handlerProvider *httpmsghandlerprovider.IncomingRequestHandlerProvider
}

func NewE2TController(logger *logger.Logger, handlerProvider *httpmsghandlerprovider.IncomingRequestHandlerProvider) *E2TController {
	return &E2TController{
		logger:          logger,
		handlerProvider: handlerProvider,
	}
}

func (c *E2TController) GetE2TInstances(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #E2TController.GetE2TInstances - request: %v", c.prettifyRequest(r))
	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.GetE2TInstancesRequest, nil, false)
}

func (c *E2TController) handleRequest(writer http.ResponseWriter, header *http.Header, requestName httpmsghandlerprovider.IncomingRequest, request models.Request, validateHeader bool) {

	handler, err := c.handlerProvider.GetHandler(requestName)

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	response, err := handler.Handle(request)

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	result, err := response.Marshal()

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	c.logger.Infof("[E2 Manager -> Client] #E2TController.handleRequest - response: %s", result)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(result)
}

func (c *E2TController) handleErrorResponse(err error, writer http.ResponseWriter) {

	var errorResponseDetails models.ErrorResponse
	var httpError int

	if err != nil {
		switch err.(type) {
		case *e2managererrors.RnibDbError:
			e2Error, _ := err.(*e2managererrors.RnibDbError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusInternalServerError
		default:
			e2Error := e2managererrors.NewInternalError()
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusInternalServerError
		}
	}

	errorResponse, _ := json.Marshal(errorResponseDetails)

	c.logger.Errorf("[E2 Manager -> Client] #E2TController.handleErrorResponse - http status: %d, error response: %+v", httpError, errorResponseDetails)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpError)
	_, err = writer.Write(errorResponse)
}

func (c *E2TController) prettifyRequest(request *http.Request) string {
	dump, _ := httputil.DumpRequest(request, true)
	requestPrettyPrint := strings.Replace(string(dump), "\r\n", " ", -1)
	return strings.Replace(requestPrettyPrint, "\n", "", -1)
}
