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

package rmrmsghandlers

import (
	"e2mgr/enums"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/utils"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type SetupResponseNotificationHandler struct {
	logger                 *logger.Logger
	rnibDataService        services.RNibDataService
	setupResponseManager   managers.ISetupResponseManager
	ranStatusChangeManager managers.IRanStatusChangeManager
	msgType                int
}

var msgTypeToMsgName = map[int]string{
	rmrCgo.RIC_X2_SETUP_RESP:         "X2 Setup Response",
	rmrCgo.RIC_X2_SETUP_FAILURE:      "X2 Setup Failure Response",
	rmrCgo.RIC_ENDC_X2_SETUP_RESP:    "ENDC Setup Response",
	rmrCgo.RIC_ENDC_X2_SETUP_FAILURE: "ENDC Setup Failure Response",
}

func NewSetupResponseNotificationHandler(logger *logger.Logger, rnibDataService services.RNibDataService, setupResponseManager managers.ISetupResponseManager, ranStatusChangeManager managers.IRanStatusChangeManager, msgType int) SetupResponseNotificationHandler {
	return SetupResponseNotificationHandler{
		logger: logger,
		rnibDataService:        rnibDataService,
		setupResponseManager:   setupResponseManager,
		ranStatusChangeManager: ranStatusChangeManager,
		msgType:                msgType,
	}
}

func (h SetupResponseNotificationHandler) Handle(request *models.NotificationRequest) {
	msgName := msgTypeToMsgName[h.msgType]
	h.logger.Infof("#SetupResponseNotificationHandler - RAN name: %s - Received %s notification", request.RanName, msgName)

	inventoryName := request.RanName

	nodebInfo, rnibErr := h.rnibDataService.GetNodeb(inventoryName)

	if rnibErr != nil {
		h.logger.Errorf("#SetupResponseNotificationHandler - RAN name: %s - Error fetching RAN from rNib: %v", request.RanName, rnibErr)
		return
	}

	if !isConnectionStatusValid(nodebInfo.ConnectionStatus) {
		h.logger.Errorf("#SetupResponseNotificationHandler - RAN name: %s - Invalid RAN connection status: %s", request.RanName, nodebInfo.ConnectionStatus)
		return
	}

	nbIdentity := &entities.NbIdentity{InventoryName: inventoryName}
	err := h.setupResponseManager.PopulateNodebByPdu(h.logger, nbIdentity, nodebInfo, request.Payload)

	if err != nil {
		return
	}

	rnibErr = h.rnibDataService.SaveNodeb(nbIdentity, nodebInfo)

	if rnibErr != nil {
		h.logger.Errorf("#SetupResponseNotificationHandler - RAN name: %s - Error saving RAN to rNib: %v", request.RanName, rnibErr)
		return
	}

	h.logger.Infof("#SetupResponseNotificationHandler - RAN name: %s - Successfully saved RAN to rNib", request.RanName)
	h.logger.Infof("#SetupResponseNotificationHandler - Summary: elapsed time for receiving and handling setup response message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))

	if !isSuccessSetupResponseMessage(h.msgType) {
		return
	}

	_ = h.ranStatusChangeManager.Execute(rmrCgo.RAN_CONNECTED, enums.RIC_TO_RAN, nodebInfo)
}

func isConnectionStatusValid(connectionStatus entities.ConnectionStatus) bool {
	return connectionStatus == entities.ConnectionStatus_CONNECTING || connectionStatus == entities.ConnectionStatus_CONNECTED
}

func isSuccessSetupResponseMessage(msgType int) bool {
	return msgType == rmrCgo.RIC_X2_SETUP_RESP || msgType == rmrCgo.RIC_ENDC_X2_SETUP_RESP
}
