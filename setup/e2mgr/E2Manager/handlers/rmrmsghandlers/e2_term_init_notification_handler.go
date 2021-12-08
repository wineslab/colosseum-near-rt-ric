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
	"e2mgr/clients"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type E2TermInitNotificationHandler struct {
	logger                  *logger.Logger
	ranDisconnectionManager *managers.RanDisconnectionManager
	e2tInstancesManager     managers.IE2TInstancesManager
	routingManagerClient    clients.IRoutingManagerClient
}

func NewE2TermInitNotificationHandler(logger *logger.Logger, ranDisconnectionManager *managers.RanDisconnectionManager, e2tInstancesManager managers.IE2TInstancesManager, routingManagerClient clients.IRoutingManagerClient) E2TermInitNotificationHandler {
	return E2TermInitNotificationHandler{
		logger:                  logger,
		ranDisconnectionManager: ranDisconnectionManager,
		e2tInstancesManager:     e2tInstancesManager,
		routingManagerClient:    routingManagerClient,
	}
}

func (h E2TermInitNotificationHandler) Handle(request *models.NotificationRequest) {
	unmarshalledPayload := models.E2TermInitPayload{}
	err := json.Unmarshal(request.Payload, &unmarshalledPayload)

	if err != nil {
		h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Error unmarshaling E2 Term Init payload: %s", err)
		return
	}

	e2tAddress := unmarshalledPayload.Address

	if len(e2tAddress) == 0 {
		h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Empty E2T address received")
		return
	}

	h.logger.Infof("#E2TermInitNotificationHandler.Handle - E2T payload: %s - handling E2_TERM_INIT", unmarshalledPayload)

	e2tInstance, err := h.e2tInstancesManager.GetE2TInstance(e2tAddress)

	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)

		if !ok {
			h.logger.Errorf("#E2TermInitNotificationHandler.Handle - Failed retrieving E2TInstance. error: %s", err)
			return
		}

		h.HandleNewE2TInstance(e2tAddress, unmarshalledPayload.PodName)
		return
	}

	if len(e2tInstance.AssociatedRanList) == 0 {
		h.logger.Infof("#E2TermInitNotificationHandler.Handle - E2T Address: %s - E2T instance has no associated RANs", e2tInstance.Address)
		return
	}

	if e2tInstance.State == entities.ToBeDeleted{
		h.logger.Infof("#E2TermInitNotificationHandler.Handle - E2T Address: %s - E2T instance status is: %s, ignore", e2tInstance.Address, e2tInstance.State)
		return
	}

	h.HandleExistingE2TInstance(e2tInstance)

	h.logger.Infof("#E2TermInitNotificationHandler.Handle - Completed handling of E2_TERM_INIT")
}

func (h E2TermInitNotificationHandler) HandleExistingE2TInstance(e2tInstance *entities.E2TInstance) {

	for _, ranName := range e2tInstance.AssociatedRanList {

		if err := h.ranDisconnectionManager.DisconnectRan(ranName); err != nil {
			if _, ok := err.(*common.ResourceNotFoundError); !ok{
				break
			}
		}
	}
}

func (h E2TermInitNotificationHandler) HandleNewE2TInstance(e2tAddress string, podName string) {

	err := h.routingManagerClient.AddE2TInstance(e2tAddress)

	if err != nil{
		h.logger.Errorf("#E2TermInitNotificationHandler.HandleNewE2TInstance - e2t address: %s - routing manager failure", e2tAddress)
		return
	}

	_ = h.e2tInstancesManager.AddE2TInstance(e2tAddress, podName)
}