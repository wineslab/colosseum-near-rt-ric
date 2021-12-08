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
	"testing"
)

func TestSetupResponsePopulateNodebByPduFailure(t *testing.T) {
	logger := tests.InitLog(t)
	nodebInfo := &entities.NodebInfo{}
	nodebIdentity := &entities.NbIdentity{}
	converter:= converters.NewEndcSetupResponseConverter(logger)
	handler := NewEndcSetupResponseManager(converter)
	err := handler.PopulateNodebByPdu(logger, nodebIdentity, nodebInfo, createRandomPayload())
	assert.NotNil(t, err)
}

func TestSetupResponsePopulateNodebByPduSuccess(t *testing.T) {
	logger := tests.InitLog(t)
	nodebInfo := &entities.NodebInfo{}
	nodebIdentity := &entities.NbIdentity{}
	converter:= converters.NewEndcSetupResponseConverter(logger)
	handler := NewEndcSetupResponseManager(converter)
	err := handler.PopulateNodebByPdu(logger, nodebIdentity, nodebInfo, createSetupResponsePayload(t))
	assert.Nil(t, err)
	assert.Equal(t, entities.ConnectionStatus_CONNECTED, nodebInfo.ConnectionStatus)
	assert.Equal(t, entities.Node_GNB, nodebInfo.NodeType)

}

func createSetupResponsePayload(t *testing.T) []byte {
	packedPdu := "202400808e00000100f600808640000200fc00090002f829504a952a0a00fd007200010c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a000211e148033e4e5e4c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a00021a0044033e4e5e000000002c001e3f271f2e3d4ff0031e3f274400010000150400000a00020000"
	var payload []byte
	_, err := fmt.Sscanf(packedPdu, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}
	return payload
}
