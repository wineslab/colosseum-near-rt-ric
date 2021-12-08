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

// #cgo CFLAGS: -I../../3rdparty/asn1codec/inc/  -I../../3rdparty/asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../../3rdparty/asn1codec/lib/ -L../../3rdparty/asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <x2reset_response_wrapper.h>
import "C"
import (
	"e2mgr/e2pdus"
	"e2mgr/enums"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
	"e2mgr/utils"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type X2ResetRequestNotificationHandler struct {
	logger                 *logger.Logger
	rnibDataService        services.RNibDataService
	ranStatusChangeManager managers.IRanStatusChangeManager
	rmrSender              *rmrsender.RmrSender
}

func NewX2ResetRequestNotificationHandler(logger *logger.Logger, rnibDataService services.RNibDataService, ranStatusChangeManager managers.IRanStatusChangeManager, rmrSender *rmrsender.RmrSender) X2ResetRequestNotificationHandler {
	return X2ResetRequestNotificationHandler{
		logger:                 logger,
		rnibDataService:        rnibDataService,
		ranStatusChangeManager: ranStatusChangeManager,
		rmrSender:              rmrSender,
	}
}

func (h X2ResetRequestNotificationHandler) Handle(request *models.NotificationRequest) {

	h.logger.Infof("#X2ResetRequestNotificationHandler.Handle - Ran name: %s", request.RanName)

	nb, rNibErr := h.rnibDataService.GetNodeb(request.RanName)
	if rNibErr != nil {
		h.logger.Errorf("#X2ResetRequestNotificationHandler.Handle - failed to retrieve nodeB entity. RanName: %s. Error: %s", request.RanName, rNibErr.Error())
		h.logger.Infof("#X2ResetRequestNotificationHandler.Handle - Summary: elapsed time for receiving and handling reset request message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))
		return
	}

	h.logger.Debugf("#X2ResetRequestNotificationHandler.Handle - nodeB entity retrieved. RanName %s, ConnectionStatus %s", nb.RanName, nb.ConnectionStatus)

	if nb.ConnectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
		h.logger.Warnf("#X2ResetRequestNotificationHandler.Handle - nodeB entity in incorrect state. RanName %s, ConnectionStatus %s", nb.RanName, nb.ConnectionStatus)
		h.logger.Infof("#X2ResetRequestNotificationHandler.Handle - Summary: elapsed time for receiving and handling reset request message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))
		return
	}

	if nb.ConnectionStatus != entities.ConnectionStatus_CONNECTED {
		h.logger.Errorf("#X2ResetRequestNotificationHandler.Handle - nodeB entity in incorrect state. RanName %s, ConnectionStatus %s", nb.RanName, nb.ConnectionStatus)
		h.logger.Infof("#X2ResetRequestNotificationHandler.Handle - Summary: elapsed time for receiving and handling reset request message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))
		return
	}

	msg := models.NewRmrMessage(rmrCgo.RIC_X2_RESET_RESP, request.RanName, e2pdus.PackedX2ResetResponse, request.TransactionId, request.GetMsgSrc())

	_ = h.rmrSender.Send(msg)
	h.logger.Infof("#X2ResetRequestNotificationHandler.Handle - Summary: elapsed time for receiving and handling reset request message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))
	_ = h.ranStatusChangeManager.Execute(rmrCgo.RAN_RESTARTED, enums.RAN_TO_RIC, nb)
}
