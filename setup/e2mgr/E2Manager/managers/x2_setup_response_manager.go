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
	"e2mgr/converters"
	"e2mgr/logger"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type X2SetupResponseManager struct {
	converter converters.IX2SetupResponseConverter
}

func NewX2SetupResponseManager(converter converters.IX2SetupResponseConverter) *X2SetupResponseManager {
	return &X2SetupResponseManager{
		converter: converter,
	}
}

func (m *X2SetupResponseManager) PopulateNodebByPdu(logger *logger.Logger, nbIdentity *entities.NbIdentity, nodebInfo *entities.NodebInfo, payload []byte) error {

	enbId, enb, err := m.converter.UnpackX2SetupResponseAndExtract(payload)

	if err != nil {
		logger.Errorf("#X2SetupResponseManager.PopulateNodebByPdu - RAN name: %s - Unpack and extract failed. %v", nodebInfo.RanName, err)
		return err
	}

	logger.Infof("#X2SetupResponseManager.PopulateNodebByPdu - RAN name: %s - Unpacked payload and extracted protobuf successfully", nodebInfo.RanName)

	nbIdentity.GlobalNbId = enbId
	nodebInfo.GlobalNbId = enbId
	nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	nodebInfo.NodeType = entities.Node_ENB
	nodebInfo.Configuration = &entities.NodebInfo_Enb{Enb: enb}

	return nil
}
