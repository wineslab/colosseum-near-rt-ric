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

// #cgo CFLAGS: -I../../3rdparty/asn1codec/inc/ -I../../3rdparty/asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../../3rdparty/asn1codec/lib/ -L../../3rdparty/asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <configuration_update_wrapper.h>
import "C"
import (
	"e2mgr/converters"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services/rmrsender"
	"e2mgr/utils"
)

type EndcConfigurationUpdateHandler struct {
	logger *logger.Logger
	rmrSender *rmrsender.RmrSender
}

func NewEndcConfigurationUpdateHandler(logger *logger.Logger, rmrSender *rmrsender.RmrSender) EndcConfigurationUpdateHandler {
	return EndcConfigurationUpdateHandler{
		logger: logger,
		rmrSender: rmrSender,
	}
}

func (h EndcConfigurationUpdateHandler) Handle(request *models.NotificationRequest) {

	refinedMessage, err := converters.UnpackX2apPduAndRefine(h.logger, e2pdus.MaxAsn1CodecAllocationBufferSize /*allocation buffer*/, request.Len, request.Payload, e2pdus.MaxAsn1CodecMessageBufferSize /*message buffer*/)

	if err != nil {
		h.logger.Errorf("#endc_configuration_update_handler.Handle - unpack failed. Error: %v", err)

		msg := models.NewRmrMessage(rmrCgo.RIC_ENDC_CONF_UPDATE_FAILURE, request.RanName, e2pdus.PackedEndcConfigurationUpdateFailure, request.TransactionId, request.GetMsgSrc())
		_ = h.rmrSender.Send(msg)

		h.logger.Infof("#EndcConfigurationUpdateHandler.Handle - Summary: elapsed time for receiving and handling endc configuration update initiating message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))
		return
	}

	h.logger.Infof("#endc_configuration_update_handler.Handle - Endc configuration update initiating message received")
	h.logger.Debugf("#endc_configuration_update_handler.Handle - Endc configuration update initiating message payload: %s", refinedMessage.PduPrint)
	msg := models.NewRmrMessage(rmrCgo.RIC_ENDC_CONF_UPDATE_ACK, request.RanName, e2pdus.PackedEndcConfigurationUpdateAck, request.TransactionId, request.GetMsgSrc())
	_ = h.rmrSender.Send(msg)

	h.logger.Infof("#EndcConfigurationUpdateHandler.Handle - Summary: elapsed time for receiving and handling endc configuration update initiating message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))
}
