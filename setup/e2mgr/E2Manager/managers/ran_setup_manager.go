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
	"e2mgr/e2managererrors"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"unsafe"
)

type RanSetupManager struct {
	logger          *logger.Logger
	rnibDataService services.RNibDataService
	rmrSender       *rmrsender.RmrSender
}

type IRanSetupManager interface {
	ExecuteSetup(nodebInfo *entities.NodebInfo, status entities.ConnectionStatus) error
}

func NewRanSetupManager(logger *logger.Logger, rmrSender *rmrsender.RmrSender, rnibDataService services.RNibDataService) *RanSetupManager {
	return &RanSetupManager{
		logger:          logger,
		rnibDataService: rnibDataService,
		rmrSender:       rmrSender,
	}
}

// Update retries and connection status 
func (m *RanSetupManager) updateConnectionStatus(nodebInfo *entities.NodebInfo, status entities.ConnectionStatus) error {
	// Update retries and connection status
	nodebInfo.ConnectionStatus = status
	err := m.rnibDataService.UpdateNodebInfo(nodebInfo)
	if err != nil {
		m.logger.Errorf("#RanSetupManager.updateConnectionStatus - Ran name: %s - Failed updating RAN's connection status to %v : %s", nodebInfo.RanName, status, err)
	} else {
		m.logger.Infof("#RanSetupManager.updateConnectionStatus - Ran name: %s - Successfully updated rNib. RAN's current connection status: %v", nodebInfo.RanName, status)
	}
	return err
}

// Decrement retries and connection status (disconnected)
func (m *RanSetupManager) updateConnectionStatusDisconnected(nodebInfo *entities.NodebInfo) error {
	// Update retries and connection status
	nodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
	err := m.rnibDataService.UpdateNodebInfo(nodebInfo)
	if err != nil {
		m.logger.Errorf("#RanSetupManager.updateConnectionStatusDisconnected - Ran name: %s - Failed updating RAN's connection status to DISCONNECTED : %s", nodebInfo.RanName, err)
	} else {
		m.logger.Infof("#RanSetupManager.updateConnectionStatusDisconnected - Ran name: %s - Successfully updated rNib. RAN's current connection status: DISCONNECTED", nodebInfo.RanName)
	}
	return err
}

func (m *RanSetupManager) prepareSetupRequest(nodebInfo *entities.NodebInfo) (int, *models.E2RequestMessage, error) {
	// Build the endc/x2 setup request
	switch nodebInfo.E2ApplicationProtocol {
	case entities.E2ApplicationProtocol_X2_SETUP_REQUEST:
		rmrMsgType := rmrCgo.RIC_X2_SETUP_REQ
		request := models.NewE2RequestMessage(nodebInfo.RanName /*tid*/, nodebInfo.Ip, uint16(nodebInfo.Port), nodebInfo.RanName, e2pdus.PackedX2setupRequest)
		return rmrMsgType, request, nil
	case entities.E2ApplicationProtocol_ENDC_X2_SETUP_REQUEST:
		rmrMsgType := rmrCgo.RIC_ENDC_X2_SETUP_REQ
		request := models.NewE2RequestMessage(nodebInfo.RanName /*tid*/, nodebInfo.Ip, uint16(nodebInfo.Port), nodebInfo.RanName, e2pdus.PackedEndcX2setupRequest)
		return rmrMsgType, request, nil
	}

	m.logger.Errorf("#RanSetupManager.prepareSetupRequest - Unsupported nodebInfo.E2ApplicationProtocol %d ", nodebInfo.E2ApplicationProtocol)
	return 0, nil, e2managererrors.NewInternalError()
}

// ExecuteSetup updates the connection status and number of attempts in the nodebInfo and send an endc/x2 setup request to establish a connection with the RAN
func (m *RanSetupManager) ExecuteSetup(nodebInfo *entities.NodebInfo, status entities.ConnectionStatus) error {

	// Update retries and connection status
	if err := m.updateConnectionStatus(nodebInfo, status); err != nil {
		return e2managererrors.NewRnibDbError()
	}

	// Build the endc/x2 setup request
	rmrMsgType, request, err := m.prepareSetupRequest(nodebInfo)
	if err != nil {
		return err
	}

	// Send the endc/x2 setup request
	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := models.NewRmrMessage(rmrMsgType, nodebInfo.RanName, request.GetMessageAsBytes(m.logger), xAction, msgSrc)

	err = m.rmrSender.Send(msg)

	if err != nil {
		m.logger.Errorf("#RanSetupManager.ExecuteSetup - failed sending setup request to RMR: %s", err)

		err := m.updateConnectionStatusDisconnected(nodebInfo)

		// Decrement retries and connection status (disconnected)
		if err != nil {
			return e2managererrors.NewRnibDbError()
		}

		return e2managererrors.NewRmrError()
	}

	return nil
}
