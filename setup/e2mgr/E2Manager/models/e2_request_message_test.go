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


package models_test

import (
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/tests"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

const transactionId = "transactionId"
const expectedMessageAsBytesHex = "31302e302e302e337c333830317c746573747c347c01020304"

func TestNewE2RequestMessage(t *testing.T){
	e2 :=models.NewE2RequestMessage(transactionId, tests.RanIp, uint16(tests.Port), tests.RanName, tests.DummyPayload)
	assert.NotNil(t, e2)
	assert.IsType(t, *e2, models.E2RequestMessage{})
	assert.Equal(t, tests.RanName, e2.RanName())
	assert.Equal(t, transactionId, e2.TransactionId())
}

func TestGetMessageAsBytes(t *testing.T){
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("#nodeb_controller_test.TestHandleRequestSuccess - failed to initialize logger, error: %s", err)
	}

	e2 := models.NewE2RequestMessage(transactionId, tests.RanIp, uint16(tests.Port), tests.RanName, tests.DummyPayload)
	bytes := e2.GetMessageAsBytes(log)
	assert.Equal(t, expectedMessageAsBytesHex, hex.EncodeToString(bytes))
}