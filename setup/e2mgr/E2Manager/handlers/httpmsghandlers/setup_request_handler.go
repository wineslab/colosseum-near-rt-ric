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

import (
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	X2SetupActivityName   = "X2_SETUP"
	EndcSetupActivityName = "ENDC_SETUP"
)

type SetupRequestHandler struct {
	rNibDataService       services.RNibDataService
	logger                *logger.Logger
	ranSetupManager       managers.IRanSetupManager
	protocol              entities.E2ApplicationProtocol
	e2tAssociationManager *managers.E2TAssociationManager
	e2tInstancesManager   managers.IE2TInstancesManager
}

func NewSetupRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService,
	ranSetupManager managers.IRanSetupManager, protocol entities.E2ApplicationProtocol, e2tInstancesManager managers.IE2TInstancesManager, e2tAssociationManager *managers.E2TAssociationManager) *SetupRequestHandler {
	return &SetupRequestHandler{
		logger:                logger,
		rNibDataService:       rNibDataService,
		ranSetupManager:       ranSetupManager,
		protocol:              protocol,
		e2tAssociationManager: e2tAssociationManager,
		e2tInstancesManager:   e2tInstancesManager,
	}
}

func (h *SetupRequestHandler) Handle(request models.Request) (models.IResponse, error) {

	setupRequest := request.(models.SetupRequest)

	err := h.validateRequestDetails(setupRequest)
	if err != nil {
		return nil, err
	}

	nodebInfo, err := h.rNibDataService.GetNodeb(setupRequest.RanName)

	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)
		if !ok {
			h.logger.Errorf("#SetupRequestHandler.Handle - failed to get nodeB entity for ran name: %v from RNIB. Error: %s", setupRequest.RanName, err)
			return nil, e2managererrors.NewRnibDbError()
		}

		result := h.connectNewRan(&setupRequest, h.protocol)
		return nil, result
	}

	if nodebInfo.ConnectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
		h.logger.Errorf("#SetupRequestHandler.connectExistingRanWithAssociatedE2TAddress - RAN: %s in wrong state (%s)", nodebInfo.RanName, entities.ConnectionStatus_name[int32(nodebInfo.ConnectionStatus)])
		result := e2managererrors.NewWrongStateError(h.getActivityName(h.protocol), entities.ConnectionStatus_name[int32(nodebInfo.ConnectionStatus)])
		return nil, result
	}

	if len(nodebInfo.AssociatedE2TInstanceAddress) != 0 {
		result := h.connectExistingRanWithAssociatedE2TAddress(nodebInfo)
		return nil, result
	}

	result := h.connectExistingRanWithoutAssociatedE2TAddress(nodebInfo)
	return nil, result
}

func createInitialNodeInfo(requestDetails *models.SetupRequest, protocol entities.E2ApplicationProtocol) (*entities.NodebInfo, *entities.NbIdentity) {

	nodebInfo := &entities.NodebInfo{
		Ip:                    requestDetails.RanIp,
		Port:                  uint32(requestDetails.RanPort),
		ConnectionStatus:      entities.ConnectionStatus_CONNECTING,
		E2ApplicationProtocol: protocol,
		RanName:               requestDetails.RanName,
	}

	nbIdentity := &entities.NbIdentity{
		InventoryName: requestDetails.RanName,
	}

	return nodebInfo, nbIdentity
}

func (h *SetupRequestHandler) connectExistingRanWithoutAssociatedE2TAddress(nodebInfo *entities.NodebInfo) error {
	e2tAddress, err := h.e2tInstancesManager.SelectE2TInstance()

	if err != nil {
		h.logger.Errorf("#SetupRequestHandler.connectExistingRanWithoutAssociatedE2TAddress - RAN name: %s - failed selecting E2T instance", nodebInfo.RanName)

		if nodebInfo.ConnectionStatus == entities.ConnectionStatus_DISCONNECTED{
			return err
		}

		nodebInfo.ConnectionStatus = entities.ConnectionStatus_DISCONNECTED
		updateError := h.rNibDataService.UpdateNodebInfo(nodebInfo)

		if updateError != nil {
			h.logger.Errorf("#SetupRequestHandler.connectExistingRanWithoutAssociatedE2TAddress - RAN name: %s - failed updating nodeb. error: %s", nodebInfo.RanName, updateError)
		}

		return err
	}

	err = h.e2tAssociationManager.AssociateRan(e2tAddress, nodebInfo)

	if err != nil {
		h.logger.Errorf("#SetupRequestHandler.connectExistingRanWithoutAssociatedE2TAddress - RAN name: %s - failed associating ran to e2t address %s. error: %s", nodebInfo.RanName, e2tAddress, err)
		return err
	}

	h.logger.Infof("#SetupRequestHandler.connectExistingRanWithoutAssociatedE2TAddress - RAN name: %s - successfully updated nodeb in rNib", nodebInfo.RanName)

	result := h.ranSetupManager.ExecuteSetup(nodebInfo, entities.ConnectionStatus_CONNECTING)
	return result
}

func (h *SetupRequestHandler) connectExistingRanWithAssociatedE2TAddress(nodebInfo *entities.NodebInfo) error {
	status := entities.ConnectionStatus_CONNECTING
	if nodebInfo.ConnectionStatus == entities.ConnectionStatus_CONNECTED {
		status = nodebInfo.ConnectionStatus
	}
	err := h.rNibDataService.UpdateNodebInfo(nodebInfo)

	if err != nil {
		h.logger.Errorf("#SetupRequestHandler.connectExistingRanWithAssociatedE2TAddress - RAN name: %s - failed resetting connection attempts of RAN. error: %s", nodebInfo.RanName, err)
		return e2managererrors.NewRnibDbError()
	}

	h.logger.Infof("#SetupRequestHandler.connectExistingRanWithAssociatedE2TAddress - RAN name: %s - successfully reset connection attempts of RAN", nodebInfo.RanName)

	result := h.ranSetupManager.ExecuteSetup(nodebInfo, status)
	return result
}

func (h *SetupRequestHandler) connectNewRan(request *models.SetupRequest, protocol entities.E2ApplicationProtocol) error {

	e2tAddress, err := h.e2tInstancesManager.SelectE2TInstance()

	if err != nil {
		h.logger.Errorf("#SetupRequestHandler.connectNewRan - RAN name: %s - failed selecting E2T instance", request.RanName)
		return err
	}

	nodebInfo, nodebIdentity := createInitialNodeInfo(request, protocol)

	err = h.rNibDataService.SaveNodeb(nodebIdentity, nodebInfo)

	if err != nil {
		h.logger.Errorf("#SetupRequestHandler.connectNewRan - RAN name: %s - failed to save initial nodeb entity in RNIB. error: %s", request.RanName, err)
		return e2managererrors.NewRnibDbError()
	}

	h.logger.Infof("#SetupRequestHandler.connectNewRan - RAN name: %s - initial nodeb entity was saved to rNib", request.RanName)

	err = h.e2tAssociationManager.AssociateRan(e2tAddress, nodebInfo)

	if err != nil {
		h.logger.Errorf("#SetupRequestHandler.connectNewRan - RAN name: %s - failed associating ran to e2t address %s. error: %s", request.RanName, e2tAddress, err)
		return err
	}

	result := h.ranSetupManager.ExecuteSetup(nodebInfo, entities.ConnectionStatus_CONNECTING)

	return result
}

func (handler *SetupRequestHandler) validateRequestDetails(request models.SetupRequest) error {

	if request.RanPort == 0 {
		handler.logger.Errorf("#SetupRequestHandler.validateRequestDetails - validation failure: port cannot be zero")
		return e2managererrors.NewRequestValidationError()
	}
	err := validation.ValidateStruct(&request,
		validation.Field(&request.RanIp, validation.Required, is.IP),
		validation.Field(&request.RanName, validation.Required),
	)

	if err != nil {
		handler.logger.Errorf("#SetupRequestHandler.validateRequestDetails - validation failure, error: %v", err)
		return e2managererrors.NewRequestValidationError()
	}

	return nil
}

func (handler *SetupRequestHandler) getActivityName(protocol entities.E2ApplicationProtocol) string {
	if protocol == entities.E2ApplicationProtocol_X2_SETUP_REQUEST {
		return X2SetupActivityName
	}
	return EndcSetupActivityName
}
