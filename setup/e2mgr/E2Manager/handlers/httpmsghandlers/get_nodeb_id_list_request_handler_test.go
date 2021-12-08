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

func setupGetNodebIdListRequestHandlerTest(t *testing.T) (*GetNodebIdListRequestHandler, *mocks.RnibReaderMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}

	rnibDataService := services.NewRnibDataService(log, config, readerMock, nil)
	handler := NewGetNodebIdListRequestHandler(log, rnibDataService)
	return handler, readerMock
}

func TestHandleGetNodebIdListSuccess(t *testing.T) {
	handler, readerMock := setupGetNodebIdListRequestHandlerTest(t)
	var rnibError error
	readerMock.On("GetListNodebIds").Return([]*entities.NbIdentity{}, rnibError)
	response, err := handler.Handle(nil)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.IsType(t, &models.GetNodebIdListResponse{}, response)
}

func TestHandleGetNodebIdListFailure(t *testing.T) {
	handler, readerMock := setupGetNodebIdListRequestHandlerTest(t)
	var nodebIdList []*entities.NbIdentity
	readerMock.On("GetListNodebIds").Return(nodebIdList, common.NewInternalError(errors.New("#reader.GetListNodebIds - Internal Error")))
	response, err := handler.Handle(nil)
	assert.NotNil(t, err)
	assert.Nil(t, response)
}
