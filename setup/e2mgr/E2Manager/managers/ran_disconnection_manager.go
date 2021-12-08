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
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type IRanDisconnectionManager interface {
	DisconnectRan(inventoryName string) error
}

type RanDisconnectionManager struct {
	logger                *logger.Logger
	config                *configuration.Configuration
	rnibDataService       services.RNibDataService
	ranSetupManager       *RanSetupManager
	e2tAssociationManager *E2TAssociationManager
}

func NewRanDisconnectionManager(logger *logger.Logger, config *configuration.Configuration, rnibDataService services.RNibDataService, e2tAssociationManager *E2TAssociationManager) *RanDisconnectionManager {
	return &RanDisconnectionManager{
		logger:                logger,
		config:                config,
		rnibDataService:       rnibDataService,
		e2tAssociationManager: e2tAssociationManager,
	}
}

func (m *RanDisconnectionManager) DisconnectRan(inventoryName string) error {
	nodebInfo, err := m.rnibDataService.GetNodeb(inventoryName)

	if err != nil {
		m.logger.Errorf("#RanDisconnectionManager.DisconnectRan - RAN name: %s - Failed fetching RAN from rNib. Error: %v", inventoryName, err)
		return err
	}

	connectionStatus := nodebInfo.GetConnectionStatus()
	m.logger.Infof("#RanDisconnectionManager.DisconnectRan - RAN name: %s - RAN's connection status: %s", nodebInfo.RanName, connectionStatus)


	if connectionStatus == entities.ConnectionStatus_SHUT_DOWN {
		m.logger.Warnf("#RanDisconnectionManager.DisconnectRan - RAN name: %s - quit. RAN's connection status is SHUT_DOWN", nodebInfo.RanName)
		return nil
	}

	if connectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
		return m.updateNodebInfo(nodebInfo, entities.ConnectionStatus_SHUT_DOWN)
	}

	err = m.updateNodebInfo(nodebInfo, entities.ConnectionStatus_DISCONNECTED)

	if err != nil {
		return err
	}

	e2tAddress := nodebInfo.AssociatedE2TInstanceAddress
	return m.e2tAssociationManager.DissociateRan(e2tAddress, nodebInfo.RanName)
}

func (m *RanDisconnectionManager) updateNodebInfo(nodebInfo *entities.NodebInfo, connectionStatus entities.ConnectionStatus) error {

	nodebInfo.ConnectionStatus = connectionStatus;
	err := m.rnibDataService.UpdateNodebInfo(nodebInfo)

	if err != nil {
		m.logger.Errorf("#RanDisconnectionManager.updateNodebInfo - RAN name: %s - Failed updating RAN's connection status to %s in rNib. Error: %v", nodebInfo.RanName, connectionStatus, err)
		return err
	}

	m.logger.Infof("#RanDisconnectionManager.updateNodebInfo - RAN name: %s - Successfully updated rNib. RAN's current connection status: %s", nodebInfo.RanName, nodebInfo.ConnectionStatus)
	return nil
}
