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
import "C"
import (
	"e2mgr/converters"
	"e2mgr/enums"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/models"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/utils"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
)

type X2ResetResponseHandler struct {
	logger                 *logger.Logger
	rnibDataService        services.RNibDataService
	ranStatusChangeManager managers.IRanStatusChangeManager
	extractor              converters.IX2ResetResponseExtractor
}

func NewX2ResetResponseHandler(logger *logger.Logger, rnibDataService services.RNibDataService, ranStatusChangeManager managers.IRanStatusChangeManager, x2ResetResponseExtractor converters.IX2ResetResponseExtractor) X2ResetResponseHandler {
	return X2ResetResponseHandler{
		logger:                 logger,
		rnibDataService:        rnibDataService,
		ranStatusChangeManager: ranStatusChangeManager,
		extractor:              x2ResetResponseExtractor,
	}
}

func (h X2ResetResponseHandler) Handle(request *models.NotificationRequest) {
	ranName := request.RanName
	h.logger.Infof("#X2ResetResponseHandler.Handle - RAN name: %s - received reset response. Payload: %x", ranName, request.Payload)

	nodebInfo, err := h.rnibDataService.GetNodeb(ranName)
	if err != nil {
		h.logger.Errorf("#x2ResetResponseHandler.Handle - RAN name: %s - failed to retrieve nodebInfo entity. Error: %s", ranName, err)
		return
	}

	if nodebInfo.ConnectionStatus == entities.ConnectionStatus_SHUTTING_DOWN {
		h.logger.Warnf("#X2ResetResponseHandler.Handle - RAN name: %s, connection status: %s - nodeB entity in incorrect state", nodebInfo.RanName, nodebInfo.ConnectionStatus)
		h.logger.Infof("#X2ResetResponseHandler.Handle - Summary: elapsed time for receiving and handling reset request message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))
		return
	}

	if nodebInfo.ConnectionStatus != entities.ConnectionStatus_CONNECTED {
		h.logger.Errorf("#X2ResetResponseHandler.Handle - RAN name: %s, connection status: %s - nodeB entity in incorrect state", nodebInfo.RanName, nodebInfo.ConnectionStatus)
		h.logger.Infof("#X2ResetResponseHandler.Handle - Summary: elapsed time for receiving and handling reset request message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))
		return
	}

	isSuccessfulResetResponse, err := h.isSuccessfulResetResponse(ranName, request.Payload)

	h.logger.Infof("#X2ResetResponseHandler.Handle - Summary: elapsed time for receiving and handling reset request message from E2 terminator: %f ms", utils.ElapsedTime(request.StartTime))

	if err != nil || !isSuccessfulResetResponse {
		return
	}

	_ = h.ranStatusChangeManager.Execute(rmrCgo.RAN_RESTARTED, enums.RIC_TO_RAN, nodebInfo)
}

func (h X2ResetResponseHandler) isSuccessfulResetResponse(ranName string, packedBuffer []byte) (bool, error) {

	criticalityDiagnostics, err := h.extractor.ExtractCriticalityDiagnosticsFromPdu(packedBuffer)

	if err != nil {
		h.logger.Errorf("#X2ResetResponseHandler.isSuccessfulResetResponse - RAN name: %s - Failed extracting pdu: %s", ranName, err)
		return false, err
	}

	if criticalityDiagnostics != nil {
		h.logger.Errorf("#X2ResetResponseHandler.isSuccessfulResetResponse - RAN name: %s - Unsuccessful RESET response message. Criticality diagnostics: %s", ranName, criticalityDiagnostics)
		return false, nil
	}

	h.logger.Infof("#X2ResetResponseHandler.isSuccessfulResetResponse - RAN name: %s - Successful RESET response message", ranName)
	return true, nil
}
