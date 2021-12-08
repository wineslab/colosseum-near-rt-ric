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

package managers

import (
	"e2mgr/enums"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rmrCgo"
	"e2mgr/services/rmrsender"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func initRanStatusChangeManagerTest(t *testing.T) (*logger.Logger, *mocks.RmrMessengerMock, *rmrsender.RmrSender) {
	logger, err := logger.InitLogger(logger.DebugLevel)
	if err != nil {
		t.Fatalf("#initStatusChangeManagerTest - failed to initialize logger, error: %s", err)
	}

	rmrMessengerMock := &mocks.RmrMessengerMock{}
	rmrSender := initRmrSender(rmrMessengerMock, logger)

	return logger, rmrMessengerMock, rmrSender
}

func TestMarshalFailure(t *testing.T) {
	logger, _, rmrSender := initRanStatusChangeManagerTest(t)
	m := NewRanStatusChangeManager(logger, rmrSender)

	nodebInfo := entities.NodebInfo{}
	err := m.Execute(123, 4, &nodebInfo)

	assert.NotNil(t, err)
}

func TestMarshalSuccess(t *testing.T) {
	logger, rmrMessengerMock, rmrSender := initRanStatusChangeManagerTest(t)
	m := NewRanStatusChangeManager(logger, rmrSender)

	nodebInfo := entities.NodebInfo{NodeType: entities.Node_ENB}
	var err error
	rmrMessengerMock.On("SendMsg", mock.Anything, true).Return(&rmrCgo.MBuf{}, err)
	err  = m.Execute(rmrCgo.RAN_CONNECTED, enums.RIC_TO_RAN, &nodebInfo)

	assert.Nil(t, err)
}
