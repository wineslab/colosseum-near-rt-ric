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

type EndcSetupResponseManager struct{
	converter converters.IEndcSetupResponseConverter
}

func NewEndcSetupResponseManager(converter converters.IEndcSetupResponseConverter) *EndcSetupResponseManager {
	return &EndcSetupResponseManager{
		converter: converter,
	}
}

func (m *EndcSetupResponseManager) PopulateNodebByPdu(logger *logger.Logger, nbIdentity *entities.NbIdentity, nodebInfo *entities.NodebInfo, payload []byte) error {

	gnbId, gnb, err := m.converter.UnpackEndcSetupResponseAndExtract(payload)

	if err != nil {
		logger.Errorf("#EndcSetupResponseManager.PopulateNodebByPdu - RAN name: %s - Unpack and extract failed. Error: %v", nodebInfo.RanName, err)
		return err
	}

	logger.Infof("#EndcSetupResponseManager.PopulateNodebByPdu - RAN name: %s - Unpacked payload and extracted protobuf successfully", nodebInfo.RanName)

	nbIdentity.GlobalNbId = gnbId
	nodebInfo.GlobalNbId = gnbId
	nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
	nodebInfo.NodeType = entities.Node_GNB
	nodebInfo.Configuration = &entities.NodebInfo_Gnb{Gnb: gnb}

	return nil
}
