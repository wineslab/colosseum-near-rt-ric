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

package httpmsghandlers

import "C"
import (
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"time"
)

type DeleteAllRequestHandler struct {
	rnibDataService     services.RNibDataService
	rmrSender           *rmrsender.RmrSender
	config              *configuration.Configuration
	logger              *logger.Logger
	e2tInstancesManager managers.IE2TInstancesManager
	rmClient            clients.IRoutingManagerClient
}

const PartialSuccessDueToRmErrorMessage = "Operation succeeded except for routing manager outbound call"

func NewDeleteAllRequestHandler(logger *logger.Logger, rmrSender *rmrsender.RmrSender, config *configuration.Configuration, rnibDataService services.RNibDataService, e2tInstancesManager managers.IE2TInstancesManager, rmClient clients.IRoutingManagerClient) *DeleteAllRequestHandler {
	return &DeleteAllRequestHandler{
		logger:              logger,
		rnibDataService:     rnibDataService,
		rmrSender:           rmrSender,
		config:              config,
		e2tInstancesManager: e2tInstancesManager,
		rmClient:            rmClient,
	}
}

func (h *DeleteAllRequestHandler) Handle(request models.Request) (models.IResponse, error) {
	h.logger.Infof("#DeleteAllRequestHandler.Handle - handling shutdown request")

	e2tAddresses, err := h.e2tInstancesManager.GetE2TAddresses()

	if err != nil {
		return nil, err
	}

	if len(e2tAddresses) == 0 {
		err, _ = h.updateNodebs(h.updateNodebInfoForceShutdown)
		return nil, err
	}

	dissocErr := h.rmClient.DissociateAllRans(e2tAddresses)

	if dissocErr != nil {
		h.logger.Warnf("#DeleteAllRequestHandler.Handle - routing manager failure. continue flow.")
	}

	err, updatedAtLeastOnce := h.updateNodebs(h.updateNodebInfoShuttingDown)

	if err != nil {
		return nil, err
	}

	err = h.e2tInstancesManager.ClearRansOfAllE2TInstances()

	if err != nil {
		return nil, err
	}

	rmrMessage := models.RmrMessage{MsgType: rmrCgo.RIC_SCTP_CLEAR_ALL}

	err = h.rmrSender.Send(&rmrMessage)

	if err != nil {
		h.logger.Errorf("#DeleteAllRequestHandler.Handle - failed to send sctp clear all message to RMR: %s", err)
		return nil, e2managererrors.NewRmrError()
	}

	if !updatedAtLeastOnce {
		h.logger.Infof("#DeleteAllRequestHandler.Handle - DB wasn't updated, not activating timer")

		if dissocErr != nil {
			return models.NewRedButtonPartialSuccessResponseModel(PartialSuccessDueToRmErrorMessage), nil
		}

		return nil, nil
	}

	time.Sleep(time.Duration(h.config.BigRedButtonTimeoutSec) * time.Second)
	h.logger.Infof("#DeleteAllRequestHandler.Handle - timer expired")

	err, _ = h.updateNodebs(h.updateNodebInfoShutDown)

	if err != nil {
		return nil, err
	}

	if dissocErr != nil {
		return models.NewRedButtonPartialSuccessResponseModel(PartialSuccessDueToRmErrorMessage), nil
	}

	return nil, nil
}

func (h *DeleteAllRequestHandler) updateNodebs(updateCb func(node *entities.NodebInfo) (error, bool)) (error, bool) {
	nbIdentityList, err := h.rnibDataService.GetListNodebIds()

	if err != nil {
		h.logger.Errorf("#DeleteAllRequestHandler.updateNodebs - failed to get nodes list from rNib. Error: %s", err)
		return e2managererrors.NewRnibDbError(), false
	}

	updatedAtLeastOnce := false

	for _, nbIdentity := range nbIdentityList {
		node, err := h.rnibDataService.GetNodeb(nbIdentity.InventoryName)

		if err != nil {
			_, ok := err.(*common.ResourceNotFoundError)

			if !ok {
				h.logger.Errorf("#DeleteAllRequestHandler.updateNodebs - failed to get nodeB entity for ran name: %s from rNib. error: %s", nbIdentity.InventoryName, err)
				return e2managererrors.NewRnibDbError(), false
			}
			continue
		}

		err, updated := updateCb(node)

		if err != nil {
			return err, false
		}

		if updated {
			updatedAtLeastOnce = true
		}
	}

	return nil, updatedAtLeastOnce
}

func (h *DeleteAllRequestHandler) updateNodebInfoForceShutdown(node *entities.NodebInfo) (error, bool) {
	err := h.updateNodebInfo(node, entities.ConnectionStatus_SHUT_DOWN, true)

	if err != nil {
		return err, false
	}

	return nil, true
}

func (h *DeleteAllRequestHandler) updateNodebInfoShuttingDown(node *entities.NodebInfo) (error, bool) {
	if node.ConnectionStatus == entities.ConnectionStatus_SHUT_DOWN {
		return nil, false
	}

	err := h.updateNodebInfo(node, entities.ConnectionStatus_SHUTTING_DOWN, true)

	if err != nil {
		return err, false
	}

	return nil, true
}

func (h *DeleteAllRequestHandler) updateNodebInfoShutDown(node *entities.NodebInfo) (error, bool) {
	if node.ConnectionStatus == entities.ConnectionStatus_SHUT_DOWN {
		return nil, false
	}

	if node.ConnectionStatus != entities.ConnectionStatus_SHUTTING_DOWN {
		h.logger.Warnf("#DeleteAllRequestHandler.updateNodebInfoShutDown - RAN name: %s - ignore, status is not Shutting Down", node.RanName)
		return nil, false
	}

	err :=  h.updateNodebInfo(node, entities.ConnectionStatus_SHUT_DOWN, false)

	if err != nil {
		return err, false
	}

	return nil, true
}

func (h *DeleteAllRequestHandler) updateNodebInfo(node *entities.NodebInfo, connectionStatus entities.ConnectionStatus, resetAssociatedE2TAddress bool) error {
	node.ConnectionStatus = connectionStatus

	if resetAssociatedE2TAddress {
		node.AssociatedE2TInstanceAddress = ""
	}

	err := h.rnibDataService.UpdateNodebInfo(node)

	if err != nil {
		h.logger.Errorf("#DeleteAllRequestHandler.updateNodebInfo - RAN name: %s - failed updating nodeB entity in rNib. error: %s", node.RanName, err)
		return e2managererrors.NewRnibDbError()
	}

	h.logger.Infof("#DeleteAllRequestHandler.updateNodebInfo - RAN name: %s, connection status: %s", node.RanName, connectionStatus)
	return nil

}
