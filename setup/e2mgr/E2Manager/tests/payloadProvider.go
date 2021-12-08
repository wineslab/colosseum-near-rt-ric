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


package tests

// #cgo CFLAGS: -I../3rdparty/asn1codec/inc/ -I../3rdparty/asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../3rdparty/asn1codec/lib/ -L../3rdparty/asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <SuccessfulOutcome.h>
//
// bool
// build_pack_x2_reset_response(size_t* packed_buf_size, unsigned char* packed_buf, size_t err_buf_size, char* err_buf){
//     bool rc = true;
//     E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
//     SuccessfulOutcome_t *successfulOutcome = calloc(1, sizeof(SuccessfulOutcome_t));
//     ResetResponse_t *resetResponse;
//     ResetResponse_IEs_t *resetResponse_IEs = calloc(1, sizeof(ResetResponse_IEs_t));
//
//     assert(pdu != 0);
//     assert(successfulOutcome != 0);
//     assert(resetResponse_IEs != 0);
//
//     pdu->present = E2AP_PDU_PR_successfulOutcome;
//     pdu->choice.successfulOutcome = successfulOutcome;
//
//     successfulOutcome->procedureCode = ProcedureCode_id_reset;
//     successfulOutcome->criticality = Criticality_reject;
//     successfulOutcome->value.present = SuccessfulOutcome__value_PR_ResetResponse;
//     resetResponse = &successfulOutcome->value.choice.ResetResponse;
//
//     CriticalityDiagnostics_IE_List_t	*critList = calloc(1, sizeof(CriticalityDiagnostics_IE_List_t));
//     assert(critList != 0);
//     resetResponse_IEs->id = ProtocolIE_ID_id_CriticalityDiagnostics;
//     resetResponse_IEs->criticality = Criticality_ignore;
//     resetResponse_IEs->value.present =  ResetResponse_IEs__value_PR_CriticalityDiagnostics;
//     ASN_SEQUENCE_ADD(resetResponse_IEs->value.choice.CriticalityDiagnostics.iEsCriticalityDiagnostics,critList);
//
//     CriticalityDiagnostics_IE_List__Member *member= calloc(1, sizeof(CriticalityDiagnostics_IE_List__Member));
//     assert(member != 0);
//     ASN_SEQUENCE_ADD(critList ,member);
//
//     ASN_SEQUENCE_ADD(&resetResponse->protocolIEs, resetResponse_IEs);
//
//     rc = per_pack_pdu(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf);
//
//     ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
//     return rc;
// }
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)
const PackedBufferSize = 4096

func BuildPackedX2ResetResponse()([]byte, error){
	payloadSize := C.ulong(PackedBufferSize)
	packedBuffer := [PackedBufferSize]C.uchar{}
	errorBuffer := [PackedBufferSize]C.char{}
	res := bool(C.build_pack_x2_reset_response(&payloadSize, &packedBuffer[0], PackedBufferSize, &errorBuffer[0]))
	if !res {
		return nil, errors.New(fmt.Sprintf("packing error: %s", C.GoString(&errorBuffer[0])))
	}
	return C.GoBytes(unsafe.Pointer(&packedBuffer), C.int(payloadSize)), nil
}
