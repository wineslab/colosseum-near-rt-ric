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
// #include <x2setup_response_wrapper.h>
import "C"
import (
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"unsafe"
)

type EndcSetupFailureResponseConverter struct {
	logger *logger.Logger
}

type IEndcSetupFailureResponseConverter interface {
	UnpackEndcSetupFailureResponseAndExtract(packedBuf []byte) (*entities.SetupFailure, error)
}

func NewEndcSetupFailureResponseConverter(logger *logger.Logger) *EndcSetupFailureResponseConverter {
	return &EndcSetupFailureResponseConverter{
		logger: logger,
	}
}


// Populate and return the EN-DC/X2 setup response failure structure with data from the pdu
func endcX2SetupFailureResponseToProtobuf(pdu *C.E2AP_PDU_t) (*entities.SetupFailure, error) {
	setupFailure := entities.SetupFailure{}

	if pdu.present == C.E2AP_PDU_PR_unsuccessfulOutcome {
		//dereference a union of pointers (C union is represented as a byte array with the size of the largest member)
		unsuccessfulOutcome := *(**C.UnsuccessfulOutcome_t)(unsafe.Pointer(&pdu.choice[0]))
		if unsuccessfulOutcome != nil && unsuccessfulOutcome.value.present == C.UnsuccessfulOutcome__value_PR_ENDCX2SetupFailure {
			endcX2SetupFailure := (*C.ENDCX2SetupFailure_t)(unsafe.Pointer(&unsuccessfulOutcome.value.choice[0]))
			if endcX2SetupFailure != nil && endcX2SetupFailure.protocolIEs.list.count > 0 {
				count:=int(endcX2SetupFailure.protocolIEs.list.count)
				endcX2SetupFailure_IEs_slice := (*[1 << 30]*C.ENDCX2SetupFailure_IEs_t)(unsafe.Pointer(endcX2SetupFailure.protocolIEs.list.array))[:count:count]
				for _, endcX2SetupFailure_IE := range endcX2SetupFailure_IEs_slice {
					if endcX2SetupFailure_IE != nil {
						switch endcX2SetupFailure_IE.value.present {
						case C.ENDCX2SetupFailure_IEs__value_PR_Cause:
							causeIE := (*C.Cause_t)(unsafe.Pointer(&endcX2SetupFailure_IE.value.choice[0]))
							err := getCause(causeIE, &setupFailure)
							if err != nil {
								return nil, err
							}
						case C.ENDCX2SetupFailure_IEs__value_PR_TimeToWait:
							setupFailure.TimeToWait = entities.TimeToWait(1 + *((*C.TimeToWait_t)(unsafe.Pointer(&endcX2SetupFailure_IE.value.choice[0]))))
						case C.ENDCX2SetupFailure_IEs__value_PR_CriticalityDiagnostics:
							cdIE := (*C.CriticalityDiagnostics_t)(unsafe.Pointer(&endcX2SetupFailure_IE.value.choice[0]))
							if cd, err := getCriticalityDiagnostics(cdIE); err == nil {
								setupFailure.CriticalityDiagnostics = cd
							} else {
								return nil, err
							}
						}
					}
				}
			}
		}
	}

	return &setupFailure, nil
}

func (c *EndcSetupFailureResponseConverter) UnpackEndcSetupFailureResponseAndExtract(packedBuf []byte) (*entities.SetupFailure, error) {
	pdu, err := UnpackX2apPdu(c.logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(packedBuf), packedBuf, e2pdus.MaxAsn1CodecMessageBufferSize)
	if err != nil {
		return nil, err
	}

	defer C.delete_pdu(pdu)

	return endcX2SetupFailureResponseToProtobuf(pdu)
}
