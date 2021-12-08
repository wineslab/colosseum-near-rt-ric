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
	"e2mgr/clients"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type E2TAssociationManager struct {
	logger             *logger.Logger
	rnibDataService    services.RNibDataService
	e2tInstanceManager IE2TInstancesManager
	rmClient           clients.IRoutingManagerClient
}

func NewE2TAssociationManager(logger *logger.Logger, rnibDataService services.RNibDataService, e2tInstanceManager IE2TInstancesManager, rmClient clients.IRoutingManagerClient) *E2TAssociationManager {
	return &E2TAssociationManager{
		logger:             logger,
		rnibDataService:    rnibDataService,
		e2tInstanceManager: e2tInstanceManager,
		rmClient:           rmClient,
	}
}

func (m *E2TAssociationManager) AssociateRan(e2tAddress string, nodebInfo *entities.NodebInfo) error {
	ranName := nodebInfo.RanName
	m.logger.Infof("#E2TAssociationManager.AssociateRan - Associating RAN %s to E2T Instance address: %s", ranName, e2tAddress)

	err := m.associateRanAndUpdateNodeb(e2tAddress, nodebInfo)
	if err != nil {
		m.logger.Errorf("#E2TAssociationManager.AssociateRan - RoutingManager failure: Failed to associate RAN %s to E2T %s. Error: %s", nodebInfo, e2tAddress, err)
		return err
	}
	err = m.e2tInstanceManager.AddRansToInstance(e2tAddress, []string{ranName})
	if err != nil {
		m.logger.Errorf("#E2TAssociationManager.AssociateRan - RAN name: %s - Failed to add RAN to E2T instance %s. Error: %s", ranName, e2tAddress, err)
		return e2managererrors.NewRnibDbError()
	}
	m.logger.Infof("#E2TAssociationManager.AssociateRan - successfully associated RAN %s with E2T %s", ranName, e2tAddress)
	return nil
}

func (m *E2TAssociationManager) associateRanAndUpdateNodeb(e2tAddress string, nodebInfo *entities.NodebInfo) error {

	rmErr := m.rmClient.AssociateRanToE2TInstance(e2tAddress, nodebInfo.RanName)
	if rmErr != nil {
		nodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	} else {
		nodebInfo.ConnectionStatus = entities.ConnectionStatus_CONNECTED
		nodebInfo.AssociatedE2TInstanceAddress = e2tAddress
	}
	rNibErr := m.rnibDataService.UpdateNodebInfo(nodebInfo)
	if rNibErr != nil {
		m.logger.Errorf("#E2TAssociationManager.associateRanAndUpdateNodeb - RAN name: %s - Failed to update nodeb entity in rNib. Error: %s", nodebInfo.RanName, rNibErr)
	}
	var err error
	if rmErr != nil {
		err = e2managererrors.NewRoutingManagerError()
	} else if rNibErr != nil{
		err = e2managererrors.NewRnibDbError()
	}
	return err
}

func (m *E2TAssociationManager) DissociateRan(e2tAddress string, ranName string) error {
	m.logger.Infof("#E2TAssociationManager.DissociateRan - Dissociating RAN %s from E2T Instance address: %s", ranName, e2tAddress)

	nodebInfo, rnibErr := m.rnibDataService.GetNodeb(ranName)
	if rnibErr != nil {
		m.logger.Errorf("#E2TAssociationManager.DissociateRan - RAN name: %s - Failed fetching RAN from rNib. Error: %s", ranName, rnibErr)
		return rnibErr
	}

	nodebInfo.AssociatedE2TInstanceAddress = ""
	rnibErr = m.rnibDataService.UpdateNodebInfo(nodebInfo)
	if rnibErr != nil {
		m.logger.Errorf("#E2TAssociationManager.DissociateRan - RAN name: %s - Failed to update RAN.AssociatedE2TInstanceAddress in rNib. Error: %s", ranName, rnibErr)
		return rnibErr
	}

	err := m.e2tInstanceManager.RemoveRanFromInstance(ranName, e2tAddress)
	if err != nil {
		m.logger.Errorf("#E2TAssociationManager.DissociateRan - RAN name: %s - Failed to remove RAN from E2T instance %s. Error: %s", ranName, e2tAddress, err)
		return err
	}

	err = m.rmClient.DissociateRanE2TInstance(e2tAddress, ranName)
	if err != nil {
		m.logger.Errorf("#E2TAssociationManager.DissociateRan - RoutingManager failure: Failed to dissociate RAN %s from E2T %s. Error: %s", ranName, e2tAddress, err)
	} else {
		m.logger.Infof("#E2TAssociationManager.DissociateRan - successfully dissociated RAN %s from E2T %s", ranName, e2tAddress)
	}
	return nil
}

func (m *E2TAssociationManager) RemoveE2tInstance(e2tInstance *entities.E2TInstance) error {
	m.logger.Infof("#E2TAssociationManager.RemoveE2tInstance -  Removing E2T %s and dessociating its associated RANs.", e2tInstance.Address)

	err := m.rmClient.DeleteE2TInstance(e2tInstance.Address, e2tInstance.AssociatedRanList)
	if err != nil {
		m.logger.Warnf("#E2TAssociationManager.RemoveE2tInstance - RoutingManager failure: Failed to delete E2T %s. Error: %s", e2tInstance.Address, err)
		// log and continue
	}

	err = m.e2tInstanceManager.RemoveE2TInstance(e2tInstance.Address)
	if err != nil {
		m.logger.Errorf("#E2TAssociationManager.RemoveE2tInstance - Failed to remove E2T %s. Error: %s", e2tInstance.Address, err)
		return err
	}

	m.logger.Infof("#E2TAssociationManager.RemoveE2tInstance -  E2T %s successfully removed.", e2tInstance.Address)
	return nil
}
