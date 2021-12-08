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
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"unsafe"
)

const (
	X2_RESET_ACTIVITY_NAME = "X2_RESET"
)

type X2ResetRequestHandler struct {
	rNibDataService services.RNibDataService
	rmrSender       *rmrsender.RmrSender
	logger          *logger.Logger
}

func NewX2ResetRequestHandler(logger *logger.Logger, rmrSender *rmrsender.RmrSender, rNibDataService services.RNibDataService) *X2ResetRequestHandler {
	return &X2ResetRequestHandler{
		rNibDataService: rNibDataService,
		rmrSender:       rmrSender,
		logger:          logger,
	}
}

func (handler *X2ResetRequestHandler) Handle(request models.Request) (models.IResponse, error) {

	resetRequest := request.(models.ResetRequest)
	handler.logger.Infof("#X2ResetRequestHandler.Handle - Ran name: %s", resetRequest.RanName)

	if len(resetRequest.Cause) == 0 {
		resetRequest.Cause = e2pdus.OmInterventionCause
	}

	payload, ok := e2pdus.KnownCausesToX2ResetPDU(resetRequest.Cause)

	if !ok {
		handler.logger.Errorf("#X2ResetRequestHandler.Handle - Unknown cause (%s)", resetRequest.Cause)
		return nil, e2managererrors.NewRequestValidationError()
	}

	nodeb, err := handler.rNibDataService.GetNodeb(resetRequest.RanName)
	if err != nil {
		handler.logger.Errorf("#X2ResetRequestHandler.Handle - failed to get status of RAN: %s from RNIB. Error: %s", resetRequest.RanName, err.Error())
		_, ok := err.(*common.ResourceNotFoundError)
		if ok {
			return nil, e2managererrors.NewResourceNotFoundError()
		}
		return nil, e2managererrors.NewRnibDbError()
	}

	if nodeb.ConnectionStatus != entities.ConnectionStatus_CONNECTED {
		handler.logger.Errorf("#X2ResetRequestHandler.Handle - RAN: %s in wrong state (%s)", resetRequest.RanName, entities.ConnectionStatus_name[int32(nodeb.ConnectionStatus)])
		return nil, e2managererrors.NewWrongStateError(X2_RESET_ACTIVITY_NAME, entities.ConnectionStatus_name[int32(nodeb.ConnectionStatus)])
	}

	var xAction []byte
	var msgSrc unsafe.Pointer
	msg := models.NewRmrMessage(rmrCgo.RIC_X2_RESET, resetRequest.RanName, payload, xAction, msgSrc)

	err = handler.rmrSender.Send(msg)

	if err != nil {
		handler.logger.Errorf("#X2ResetRequestHandler.Handle - failed to send reset message to RMR: %s", err)
		return nil, e2managererrors.NewRmrError()
	}

	handler.logger.Infof("#X2ResetRequestHandler.Handle - sent x2 reset to RAN: %s with cause: %s", resetRequest.RanName, resetRequest.Cause)

	return nil, nil
}
