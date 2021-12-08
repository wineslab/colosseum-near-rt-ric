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
	"e2mgr/configuration"
	"e2mgr/logger"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"time"
)

type IE2TShutdownManager interface {
	Shutdown(e2tInstance *entities.E2TInstance) error
}

type E2TShutdownManager struct {
	logger                *logger.Logger
	config                *configuration.Configuration
	rnibDataService       services.RNibDataService
	e2TInstancesManager   IE2TInstancesManager
	e2tAssociationManager *E2TAssociationManager
	kubernetesManager     *KubernetesManager
}

func NewE2TShutdownManager(logger *logger.Logger, config *configuration.Configuration, rnibDataService services.RNibDataService, e2TInstancesManager IE2TInstancesManager, e2tAssociationManager *E2TAssociationManager, kubernetes *KubernetesManager) *E2TShutdownManager {
	return &E2TShutdownManager{
		logger:                logger,
		config:                config,
		rnibDataService:       rnibDataService,
		e2TInstancesManager:   e2TInstancesManager,
		e2tAssociationManager: e2tAssociationManager,
		kubernetesManager:     kubernetes,
	}
}

func (m E2TShutdownManager) Shutdown(e2tInstance *entities.E2TInstance) error {
	m.logger.Infof("#E2TShutdownManager.Shutdown - E2T %s is Dead, RIP", e2tInstance.Address)

	isE2tInstanceBeingDeleted := m.isE2tInstanceAlreadyBeingDeleted(e2tInstance)
	if isE2tInstanceBeingDeleted {
		m.logger.Infof("#E2TShutdownManager.Shutdown - E2T %s is already being deleted", e2tInstance.Address)
		return nil
	}

	//go m.kubernetesManager.DeletePod(e2tInstance.PodName)

	err := m.markE2tInstanceToBeDeleted(e2tInstance)
	if err != nil {
		m.logger.Errorf("#E2TShutdownManager.Shutdown - Failed to mark E2T %s as 'ToBeDeleted'.", e2tInstance.Address)
		return err
	}

	err = m.clearNodebsAssociation(e2tInstance.AssociatedRanList)
	if err != nil {
		m.logger.Errorf("#E2TShutdownManager.Shutdown - Failed to clear nodebs association to E2T %s.", e2tInstance.Address)
		return err
	}

	err = m.e2tAssociationManager.RemoveE2tInstance(e2tInstance)
	if err != nil {
		m.logger.Errorf("#E2TShutdownManager.Shutdown - Failed to remove E2T %s.", e2tInstance.Address)
		return err
	}

	m.logger.Infof("#E2TShutdownManager.Shutdown - E2T %s was shutdown successfully.", e2tInstance.Address)
	return nil
}

func (m E2TShutdownManager) clearNodebsAssociation(ranNamesToBeDissociated []string) error {
	for _, ranName := range ranNamesToBeDissociated {
		nodeb, err := m.rnibDataService.GetNodeb(ranName)
		if err != nil {
			m.logger.Warnf("#E2TShutdownManager.associateAndSetupNodebs - Failed to get nodeb %s from db.", ranName)
			_, ok := err.(*common.ResourceNotFoundError)
			if !ok {
				continue
			}
			return err
		}
		nodeb.AssociatedE2TInstanceAddress = ""
		nodeb.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED

		err = m.rnibDataService.UpdateNodebInfo(nodeb)
		if err != nil {
			m.logger.Errorf("#E2TShutdownManager.associateAndSetupNodebs - Failed to save nodeb %s from db.", ranName)
			return err
		}
	}
	return nil
}

func (m E2TShutdownManager) markE2tInstanceToBeDeleted(e2tInstance *entities.E2TInstance) error {
	e2tInstance.State = entities.ToBeDeleted
	e2tInstance.DeletionTimestamp = time.Now().UnixNano()

	return m.rnibDataService.SaveE2TInstance(e2tInstance)
}

func (m E2TShutdownManager) isE2tInstanceAlreadyBeingDeleted(e2tInstance *entities.E2TInstance) bool {
	delta := time.Now().UnixNano() - e2tInstance.DeletionTimestamp
	timestampNanosec := int64(time.Duration(m.config.E2TInstanceDeletionTimeoutMs) * time.Millisecond)

	return e2tInstance.State == entities.ToBeDeleted && delta <= timestampNanosec
}
