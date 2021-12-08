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

const (
	maxNrOfErrors = 256
)

type X2SetupFailureResponseConverter struct {
	logger *logger.Logger
}

type IX2SetupFailureResponseConverter interface {
	UnpackX2SetupFailureResponseAndExtract(packedBuf []byte) (*entities.SetupFailure, error)
}

func NewX2SetupFailureResponseConverter(logger *logger.Logger) *X2SetupFailureResponseConverter {
	return &X2SetupFailureResponseConverter{
		logger: logger,
	}
}

// The following are possible values of a choice field, find which the pdu contains.
func getCause(causeIE *C.Cause_t, setupFailure *entities.SetupFailure) error {
	switch causeIE.present {
	case C.Cause_PR_radioNetwork:
		v := (*C.CauseRadioNetwork_t)(unsafe.Pointer(&causeIE.choice[0]))
		setupFailure.CauseGroup = &entities.SetupFailure_NetworkLayerCause{NetworkLayerCause: entities.RadioNetworkLayer_Cause(1 + *v)}
	case C.Cause_PR_transport:
		v := (*C.CauseTransport_t)(unsafe.Pointer(&causeIE.choice[0]))
		setupFailure.CauseGroup = &entities.SetupFailure_TransportLayerCause{TransportLayerCause: entities.TransportLayer_Cause(1 + *v)}
	case C.Cause_PR_protocol:
		v := (*C.CauseProtocol_t)(unsafe.Pointer(&causeIE.choice[0]))
		setupFailure.CauseGroup = &entities.SetupFailure_ProtocolCause{ProtocolCause: entities.Protocol_Cause(1 + *v)}
	case C.Cause_PR_misc:
		v := (*C.CauseMisc_t)(unsafe.Pointer(&causeIE.choice[0]))
		setupFailure.CauseGroup = &entities.SetupFailure_MiscellaneousCause{MiscellaneousCause: entities.Miscellaneous_Cause(1 + *v)}
	}
	return nil
}

func getCriticalityDiagnostics(critDiagIE *C.CriticalityDiagnostics_t) (*entities.CriticalityDiagnostics, error) {
	var critDiag *entities.CriticalityDiagnostics

	if critDiagIE.procedureCode != nil {
		critDiag = &entities.CriticalityDiagnostics{}
		critDiag.ProcedureCode = uint32(*critDiagIE.procedureCode)

	}

	if critDiagIE.triggeringMessage != nil {
		if critDiag == nil {
			critDiag = &entities.CriticalityDiagnostics{}
		}
		critDiag.TriggeringMessage = entities.TriggeringMessage(1 + *critDiagIE.triggeringMessage)

	}

	if critDiagIE.procedureCriticality != nil {
		if critDiag == nil {
			critDiag = &entities.CriticalityDiagnostics{}
		}
		critDiag.ProcedureCriticality = entities.Criticality(1 + *critDiagIE.procedureCriticality)

	}

	if critDiagIE.iEsCriticalityDiagnostics != nil && critDiagIE.iEsCriticalityDiagnostics.list.count > 0 && critDiagIE.iEsCriticalityDiagnostics.list.count < maxNrOfErrors {
		if critDiag == nil {
			critDiag = &entities.CriticalityDiagnostics{}
		}
		var infoElements []*entities.InformationElementCriticalityDiagnostic
		iEsCriticalityDiagnostics := (*C.CriticalityDiagnostics_IE_List_t)(critDiagIE.iEsCriticalityDiagnostics)
		count:=int(iEsCriticalityDiagnostics.list.count)
		iEsCriticalityDiagnostics_slice := (*[1 << 30]*C.CriticalityDiagnostics_IE_List__Member)(unsafe.Pointer(iEsCriticalityDiagnostics.list.array))[:count:count]
		for _, criticalityDiagnostics_IE_List__Member := range  iEsCriticalityDiagnostics_slice {
			infoElement := &entities.InformationElementCriticalityDiagnostic{IeCriticality: entities.Criticality(1 + criticalityDiagnostics_IE_List__Member.iECriticality)}
			infoElement.IeId = uint32(criticalityDiagnostics_IE_List__Member.iE_ID)
			infoElement.TypeOfError = entities.TypeOfError(1 + criticalityDiagnostics_IE_List__Member.typeOfError)

			infoElements = append(infoElements, infoElement)

		}
		critDiag.InformationElementCriticalityDiagnostics = infoElements
	}

	return critDiag, nil
}

// Populate and return the EN-DC/X2 setup response failure structure with data from the pdu
func x2SetupFailureResponseToProtobuf(pdu *C.E2AP_PDU_t) (*entities.SetupFailure, error) {
	setupFailure := entities.SetupFailure{}

	if pdu.present == C.E2AP_PDU_PR_unsuccessfulOutcome {
		//dereference a union of pointers (C union is represented as a byte array with the size of the largest member)
		unsuccessfulOutcome := *(**C.UnsuccessfulOutcome_t)(unsafe.Pointer(&pdu.choice[0]))
		if unsuccessfulOutcome != nil && unsuccessfulOutcome.value.present == C.UnsuccessfulOutcome__value_PR_X2SetupFailure {
			x2SetupFailure := (*C.X2SetupFailure_t)(unsafe.Pointer(&unsuccessfulOutcome.value.choice[0]))
			if x2SetupFailure != nil && x2SetupFailure.protocolIEs.list.count > 0 {
				count:=int(x2SetupFailure.protocolIEs.list.count)
				x2SetupFailure_IEs_slice := (*[1 << 30]*C.X2SetupFailure_IEs_t)(unsafe.Pointer(x2SetupFailure.protocolIEs.list.array))[:count:count]
				for _, x2SetupFailure_IE := range x2SetupFailure_IEs_slice {
					if x2SetupFailure_IE != nil {
						switch x2SetupFailure_IE.value.present {
						case C.X2SetupFailure_IEs__value_PR_Cause:
							causeIE := (*C.Cause_t)(unsafe.Pointer(&x2SetupFailure_IE.value.choice[0]))
							err := getCause(causeIE, &setupFailure)
							if err != nil {
								return nil, err
							}
						case C.X2SetupFailure_IEs__value_PR_TimeToWait:
							setupFailure.TimeToWait = entities.TimeToWait(1 + *((*C.TimeToWait_t)(unsafe.Pointer(&x2SetupFailure_IE.value.choice[0]))))
						case C.X2SetupFailure_IEs__value_PR_CriticalityDiagnostics:
							cdIE := (*C.CriticalityDiagnostics_t)(unsafe.Pointer(&x2SetupFailure_IE.value.choice[0]))
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

func (c *X2SetupFailureResponseConverter) UnpackX2SetupFailureResponseAndExtract(packedBuf []byte) (*entities.SetupFailure, error) {
	pdu, err := UnpackX2apPdu(c.logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(packedBuf), packedBuf, e2pdus.MaxAsn1CodecMessageBufferSize)
	if err != nil {
		return nil, err
	}

	defer C.delete_pdu(pdu)

	return x2SetupFailureResponseToProtobuf(pdu)
}