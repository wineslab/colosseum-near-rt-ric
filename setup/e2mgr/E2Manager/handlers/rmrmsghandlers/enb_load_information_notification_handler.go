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
//

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package rmrmsghandlers

/*import (
	"e2mgr/converters"
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"time"
)

type EnbLoadInformationNotificationHandler struct {
	logger          *logger.Logger
	rnibDataService services.RNibDataService
	//extractor       converters.IEnbLoadInformationExtractor
}

func NewEnbLoadInformationNotificationHandler(logger *logger.Logger, rnibDataService services.RNibDataService, extractor converters.IEnbLoadInformationExtractor) EnbLoadInformationNotificationHandler {
	return EnbLoadInformationNotificationHandler{
		logger:          logger,
		rnibDataService: rnibDataService,
		extractor: extractor,
	}
}

func elapsed(startTime time.Time) float64 {
	return float64(time.Since(startTime)) / float64(time.Millisecond)
}

func (h EnbLoadInformationNotificationHandler) Handle(request *models.NotificationRequest) {

	pdu, err := converters.UnpackX2apPdu(h.logger, e2pdus.MaxAsn1CodecAllocationBufferSize, request.Len, request.Payload, e2pdus.MaxAsn1CodecMessageBufferSize)

	if err != nil {
		h.logger.Errorf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Unpack failed. Error: %v", request.RanName, err)
		return
	}

	h.logger.Debugf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Unpacked message successfully", request.RanName)

	ranLoadInformation := &entities.RanLoadInformation{LoadTimestamp: uint64(request.StartTime.UnixNano())}

	err = h.extractor.ExtractAndBuildRanLoadInformation(pdu, ranLoadInformation)

	if err != nil {
		h.logger.Errorf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Failed at ExtractAndBuildRanLoadInformation. Error: %v", request.RanName, err)
		return
	}

	h.logger.Debugf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Successfully done with extracting and building RAN load information. elapsed: %f ms", request.RanName, elapsed(request.StartTime))

	rnibErr := h.rnibDataService.SaveRanLoadInformation(request.RanName, ranLoadInformation)

	if rnibErr != nil {
		h.logger.Errorf("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Failed saving RAN load information. Error: %v", request.RanName, rnibErr)
		return
	}

	h.logger.Infof("#EnbLoadInformationNotificationHandler.Handle - RAN name: %s - Successfully saved RAN load information to RNIB. elapsed: %f ms", request.RanName, elapsed(request.StartTime))
}*/
