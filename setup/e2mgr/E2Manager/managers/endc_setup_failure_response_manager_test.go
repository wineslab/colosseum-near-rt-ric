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

package managers

import (
	"e2mgr/converters"
	"e2mgr/tests"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestPopulateNodebByPduFailure(t *testing.T) {
	logger := tests.InitLog(t)
	nodebInfo := &entities.NodebInfo{}
	nodebIdentity := &entities.NbIdentity{}
	converter := converters.NewEndcSetupFailureResponseConverter(logger)
	handler := NewEndcSetupFailureResponseManager(converter)
	err := handler.PopulateNodebByPdu(logger, nodebIdentity, nodebInfo, createRandomPayload())
	assert.NotNil(t, err)
}

func TestPopulateNodebByPduSuccess(t *testing.T) {
	logger := tests.InitLog(t)
	nodebInfo := &entities.NodebInfo{}
	nodebIdentity := &entities.NbIdentity{}
	converter := converters.NewEndcSetupFailureResponseConverter(logger)
	handler := NewEndcSetupFailureResponseManager(converter)
	err := handler.PopulateNodebByPdu(logger, nodebIdentity, nodebInfo, createSetupFailureResponsePayload(t))
	assert.Nil(t, err)
	assert.Equal(t, entities.ConnectionStatus_CONNECTED_SETUP_FAILED, nodebInfo.ConnectionStatus)
	assert.Equal(t, entities.Failure_ENDC_X2_SETUP_FAILURE, nodebInfo.FailureType)

}

func createSetupFailureResponsePayload(t *testing.T) []byte {
	packedPdu := "4024001a0000030005400200000016400100001140087821a00000008040"
	var payload []byte
	_, err := fmt.Sscanf(packedPdu, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}
	return payload
}

func createRandomPayload() []byte {
	payload := make([]byte, 20)
	rand.Read(payload)
	return payload
}
