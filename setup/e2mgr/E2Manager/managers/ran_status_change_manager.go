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
	"e2mgr/models"
	"e2mgr/services/rmrsender"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"unsafe"
)

type RanStatusChangeManager struct {
	logger    *logger.Logger
	rmrSender *rmrsender.RmrSender
}

func NewRanStatusChangeManager(logger *logger.Logger, rmrSender *rmrsender.RmrSender) *RanStatusChangeManager {
	return &RanStatusChangeManager{
		logger:    logger,
		rmrSender: rmrSender,
	}
}

type IRanStatusChangeManager interface {
	Execute(msgType int, msgDirection enums.MessageDirection, nodebInfo *entities.NodebInfo) error
}

func (m *RanStatusChangeManager) Execute(msgType int, msgDirection enums.MessageDirection, nodebInfo *entities.NodebInfo) error {

	resourceStatusPayload := models.NewResourceStatusPayload(nodebInfo.NodeType, msgDirection)
	resourceStatusJson, err := json.Marshal(resourceStatusPayload)

	if err != nil {
		m.logger.Errorf("#RanStatusChangeManager.Execute - RAN name: %s - Error marshaling resource status payload: %v", nodebInfo.RanName, err)
		return err
	}

	var xAction []byte
	var msgSrc unsafe.Pointer
	rmrMessage := models.NewRmrMessage(msgType, nodebInfo.RanName, resourceStatusJson, xAction, msgSrc)
	return m.rmrSender.Send(rmrMessage)
}
