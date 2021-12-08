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
	"e2mgr/models"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
)

const VALIDATION_FAILURE_MESSAGE = "#UpdateGnbRequestHandler.Handle - validation failure: %s is a mandatory field"

type UpdateGnbRequestHandler struct {
	logger          *logger.Logger
	rNibDataService services.RNibDataService
}

func NewUpdateGnbRequestHandler(logger *logger.Logger, rNibDataService services.RNibDataService) *UpdateGnbRequestHandler {
	return &UpdateGnbRequestHandler{
		logger:          logger,
		rNibDataService: rNibDataService,
	}
}

func (h *UpdateGnbRequestHandler) Handle(request models.Request) (models.IResponse, error) {

	updateGnbRequest := request.(models.UpdateGnbRequest)

	h.logger.Infof("#UpdateGnbRequestHandler.Handle - Ran name: %s", updateGnbRequest.RanName)

	err := h.validateRequestBody(updateGnbRequest)

	if err != nil {
		return nil, err
	}

	nodebInfo, err := h.rNibDataService.GetNodeb(updateGnbRequest.RanName)

	if err != nil {
		_, ok := err.(*common.ResourceNotFoundError)
		if !ok {
			h.logger.Errorf("#UpdateGnbRequestHandler.Handle - RAN name: %s - failed to get nodeb entity from RNIB. Error: %s", updateGnbRequest.RanName, err)
			return nil, e2managererrors.NewRnibDbError()
		}

		h.logger.Errorf("#UpdateGnbRequestHandler.Handle - RAN name: %s - RAN not found on RNIB. Error: %s", updateGnbRequest.RanName, err)
		return nil, e2managererrors.NewResourceNotFoundError()
	}

	err = h.updateGnbCells(nodebInfo, updateGnbRequest)

	if err != nil {
		return nil, err
	}

	return models.NewUpdateGnbResponse(nodebInfo), nil
}

func (h *UpdateGnbRequestHandler) updateGnbCells(nodebInfo *entities.NodebInfo, updateGnbRequest models.UpdateGnbRequest) error {

	ranName := nodebInfo.RanName
	gnb := nodebInfo.GetGnb()

	if gnb == nil {
		h.logger.Errorf("#UpdateGnbRequestHandler.updateGnbCells - RAN name: %s - nodeb missing gnb configuration", ranName)
		return e2managererrors.NewInternalError()
	}

	if len(gnb.ServedNrCells) != 0 {
		err := h.rNibDataService.RemoveServedNrCells(ranName, gnb.ServedNrCells)

		if err != nil {
			h.logger.Errorf("#UpdateGnbRequestHandler.updateGnbCells - RAN name: %s - Failed removing served nr cells", ranName)
			return e2managererrors.NewRnibDbError()
		}
	}

	gnb.ServedNrCells = updateGnbRequest.ServedNrCells

	err := h.rNibDataService.UpdateGnbCells(nodebInfo, updateGnbRequest.ServedNrCells)

	if err != nil {
		h.logger.Errorf("#UpdateGnbRequestHandler.updateGnbCells - RAN name: %s - Failed updating GNB cells. Error: %s", ranName, err)
		return e2managererrors.NewRnibDbError()
	}

	h.logger.Infof("#UpdateGnbRequestHandler.updateGnbCells - RAN name: %s - Successfully updated GNB cells", ranName)
	return nil
}

func (h *UpdateGnbRequestHandler) validateRequestBody(updateGnbRequest models.UpdateGnbRequest) error {

	if len(updateGnbRequest.ServedNrCells) == 0 {
		h.logger.Errorf(VALIDATION_FAILURE_MESSAGE+" and cannot be empty", "servedCells")
		return e2managererrors.NewRequestValidationError()
	}

	for _, servedNrCell := range updateGnbRequest.ServedNrCells {
		if servedNrCell.ServedNrCellInformation == nil {
			h.logger.Errorf(VALIDATION_FAILURE_MESSAGE+" and cannot be empty", "servedNrCellInformation")
			return e2managererrors.NewRequestValidationError()
		}

		err := isServedNrCellInformationValid(servedNrCell.ServedNrCellInformation)

		if err != nil {
			h.logger.Errorf(VALIDATION_FAILURE_MESSAGE, err)
			return e2managererrors.NewRequestValidationError()
		}

		if len(servedNrCell.NrNeighbourInfos) == 0 {
			continue
		}

		for _, nrNeighbourInformation := range servedNrCell.NrNeighbourInfos {

			err := isNrNeighbourInformationValid(nrNeighbourInformation)

			if err != nil {
				h.logger.Errorf(VALIDATION_FAILURE_MESSAGE, err)
				return e2managererrors.NewRequestValidationError()
			}

		}
	}

	return nil
}

func isServedNrCellInformationValid(servedNrCellInformation *entities.ServedNRCellInformation) error {
	if servedNrCellInformation.CellId == "" {
		return errors.New("cellId")
	}

	if servedNrCellInformation.ChoiceNrMode == nil {
		return errors.New("choiceNrMode")
	}

	if servedNrCellInformation.NrMode == entities.Nr_UNKNOWN {
		return errors.New("nrMode")
	}

	if servedNrCellInformation.NrPci == 0 {
		return errors.New("nrPci")
	}

	if len(servedNrCellInformation.ServedPlmns) == 0 {
		return errors.New("servedPlmns")
	}

	return isServedNrCellInfoChoiceNrModeValid(servedNrCellInformation.ChoiceNrMode)
}

func isServedNrCellInfoChoiceNrModeValid(choiceNrMode *entities.ServedNRCellInformation_ChoiceNRMode) error {
	if choiceNrMode.Fdd != nil {
		return isServedNrCellInfoFddValid(choiceNrMode.Fdd)
	}

	if choiceNrMode.Tdd != nil {
		return isServedNrCellInfoTddValid(choiceNrMode.Tdd)
	}

	return errors.New("served nr cell fdd / tdd")
}

func isServedNrCellInfoTddValid(tdd *entities.ServedNRCellInformation_ChoiceNRMode_TddInfo) error {
	return nil
}

func isServedNrCellInfoFddValid(fdd *entities.ServedNRCellInformation_ChoiceNRMode_FddInfo) error {
	return nil
}

func isNrNeighbourInformationValid(nrNeighbourInformation *entities.NrNeighbourInformation) error {
	if nrNeighbourInformation.NrCgi == "" {
		return errors.New("nrCgi")
	}

	if nrNeighbourInformation.ChoiceNrMode == nil {
		return errors.New("choiceNrMode")
	}

	if nrNeighbourInformation.NrMode == entities.Nr_UNKNOWN {
		return errors.New("nrMode")
	}

	if nrNeighbourInformation.NrPci == 0 {
		return errors.New("nrPci")
	}

	return isNrNeighbourInfoChoiceNrModeValid(nrNeighbourInformation.ChoiceNrMode)
}

func isNrNeighbourInfoChoiceNrModeValid(choiceNrMode *entities.NrNeighbourInformation_ChoiceNRMode) error {
	if choiceNrMode.Fdd != nil {
		return isNrNeighbourInfoFddValid(choiceNrMode.Fdd)
	}

	if choiceNrMode.Tdd != nil {
		return isNrNeighbourInfoTddValid(choiceNrMode.Tdd)
	}

	return errors.New("nr neighbour fdd / tdd")
}

func isNrNeighbourInfoTddValid(tdd *entities.NrNeighbourInformation_ChoiceNRMode_TddInfo) error {
	return nil
}

func isNrNeighbourInfoFddValid(fdd *entities.NrNeighbourInformation_ChoiceNRMode_FddInfo) error {
	return nil
}
