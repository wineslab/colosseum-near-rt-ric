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


package converters

// #cgo CFLAGS: -I../3rdparty/asn1codec/inc/ -I../3rdparty/asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../3rdparty/asn1codec/lib/ -L../3rdparty/asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <SuccessfulOutcome.h>
import "C"
import (
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"unsafe"
)

type X2ResetResponseExtractor struct {
	logger *logger.Logger
}

func NewX2ResetResponseExtractor(logger *logger.Logger) *X2ResetResponseExtractor {
	return &X2ResetResponseExtractor{
		logger: logger,
	}
}

type IX2ResetResponseExtractor interface {
	ExtractCriticalityDiagnosticsFromPdu(packedBuffer []byte) (*entities.CriticalityDiagnostics, error)
}

func (e *X2ResetResponseExtractor) ExtractCriticalityDiagnosticsFromPdu(packedBuffer []byte) (*entities.CriticalityDiagnostics, error) {
	pdu, err := UnpackX2apPdu(e.logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(packedBuffer), packedBuffer, e2pdus.MaxAsn1CodecMessageBufferSize)

	if err != nil {
		return nil, err
	}

	if pdu.present != C.E2AP_PDU_PR_successfulOutcome {
		return nil, fmt.Errorf("Invalid E2AP_PDU value")
	}

	successfulOutcome := *(**C.SuccessfulOutcome_t)(unsafe.Pointer(&pdu.choice[0]))

	if successfulOutcome == nil || successfulOutcome.value.present != C.SuccessfulOutcome__value_PR_ResetResponse {
		return nil, fmt.Errorf("Unexpected SuccessfulOutcome value")
	}

	resetResponse := (*C.ResetResponse_t)(unsafe.Pointer(&successfulOutcome.value.choice[0]))

	protocolIEsListCount := resetResponse.protocolIEs.list.count

	if protocolIEsListCount == 0 {
		return nil, nil
	}

	if protocolIEsListCount != 1 {
		return nil, fmt.Errorf("Invalid protocolIEs list count")
	}

	resetResponseIEs := (*[1 << 30]*C.ResetResponse_IEs_t)(unsafe.Pointer(resetResponse.protocolIEs.list.array))[:int(protocolIEsListCount):int(protocolIEsListCount)]

	resetResponseIE := resetResponseIEs[0]

	if resetResponseIE.value.present != C.ResetResponse_IEs__value_PR_CriticalityDiagnostics {
		return nil, fmt.Errorf("Invalid protocolIEs value")
	}

	cd := (*C.CriticalityDiagnostics_t)(unsafe.Pointer(&resetResponseIE.value.choice[0]))

	return getCriticalityDiagnostics(cd)
}
