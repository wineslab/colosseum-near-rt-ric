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


package httpmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupGetNodebRequestHandlerTest(t *testing.T) (*GetNodebRequestHandler, *mocks.RnibReaderMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, nil)
	handler := NewGetNodebRequestHandler(log, rnibDataService)
	return handler, readerMock
}

func TestHandleGetNodebSuccess(t *testing.T) {
	handler, readerMock := setupGetNodebRequestHandlerTest(t)

	ranName := "test1"
	var rnibError error
	readerMock.On("GetNodeb", ranName).Return(&entities.NodebInfo{RanName:ranName}, rnibError)
	response, err := handler.Handle(models.GetNodebRequest{RanName: ranName})
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.IsType(t, &models.GetNodebResponse{}, response)
}

func TestHandleGetNodebFailure(t *testing.T) {
	handler, readerMock := setupGetNodebRequestHandlerTest(t)
	ranName := "test1"
	var nodebInfo *entities.NodebInfo
	readerMock.On("GetNodeb", ranName).Return(nodebInfo, common.NewInternalError(errors.New("#reader.GetNodeb - Internal Error")))
	response, err := handler.Handle(models.GetNodebRequest{RanName: ranName})
	assert.NotNil(t, err)
	assert.Nil(t, response)
}
