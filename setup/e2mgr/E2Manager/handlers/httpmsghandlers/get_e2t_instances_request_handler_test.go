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

package httpmsghandlers

import (
	"e2mgr/configuration"
	"e2mgr/managers"
	"e2mgr/mocks"
	"e2mgr/models"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"testing"
)

const E2TAddress2 = "10.0.2.15:3213"

func setupGetE2TInstancesListRequestHandlerTest(t *testing.T) (*GetE2TInstancesRequestHandler, *mocks.RnibReaderMock) {
	log := initLog(t)
	config := &configuration.Configuration{RnibRetryIntervalMs: 10, MaxRnibConnectionAttempts: 3}
	readerMock := &mocks.RnibReaderMock{}
	rnibDataService := services.NewRnibDataService(log, config, readerMock, nil)
	e2tInstancesManager := managers.NewE2TInstancesManager(rnibDataService, log)
	handler := NewGetE2TInstancesRequestHandler(log, e2tInstancesManager)
	return handler, readerMock
}

func TestGetE2TInstancesFailure(t *testing.T) {
	handler, rnibReaderMock := setupGetE2TInstancesListRequestHandlerTest(t)
	rnibReaderMock.On("GetE2TAddresses").Return([]string{}, common.NewInternalError(errors.New("error")))
	_, err := handler.Handle(nil)
	assert.NotNil(t, err)
}

func TestGetE2TInstancesNoInstances(t *testing.T) {
	handler, rnibReaderMock := setupGetE2TInstancesListRequestHandlerTest(t)
	rnibReaderMock.On("GetE2TAddresses").Return([]string{}, nil)
	resp, err := handler.Handle(nil)
	assert.Nil(t, err)
	assert.IsType(t, models.E2TInstancesResponse{}, resp)
	assert.Len(t, resp, 0)
}

func TestGetE2TInstancesSuccess(t *testing.T) {
	handler, rnibReaderMock := setupGetE2TInstancesListRequestHandlerTest(t)

	e2tAddresses := []string{E2TAddress, E2TAddress2}
	rnibReaderMock.On("GetE2TAddresses").Return(e2tAddresses, nil)
	e2tInstance := entities.E2TInstance{Address: E2TAddress, AssociatedRanList: []string{"test1", "test2"}}
	e2tInstance2 := entities.E2TInstance{Address: E2TAddress2, AssociatedRanList: []string{"test3", "test4", "test5"}}

	rnibReaderMock.On("GetE2TInstances", e2tAddresses).Return([]*entities.E2TInstance{&e2tInstance, &e2tInstance2}, nil)
	resp, err := handler.Handle(nil)
	assert.Nil(t, err)
	assert.IsType(t, models.E2TInstancesResponse{}, resp)
	assert.Len(t, resp, 2)
}
