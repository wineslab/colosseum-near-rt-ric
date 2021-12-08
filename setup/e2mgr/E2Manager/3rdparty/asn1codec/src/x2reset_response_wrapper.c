/*
 * Copyright 2019 AT&T Intellectual Property
 * Copyright 2019 Nokia
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * This source code is part of the near-RT RIC (RAN Intelligent Controller)
 * platform project (RICP).
 */


#include <string.h>
#include <errno.h>
#undef NDEBUG
#include <assert.h>
#include <asn_application.h>
#include <E2AP-PDU.h>
#include <ProcedureCode.h>
#include <SuccessfulOutcome.h>
#include <ProtocolIE-ID.h>
#include <ProtocolIE-Field.h>
#include <x2reset_response_wrapper.h>

/*
 * Build and pack a reset response.
 * Abort the process on allocation failure.
 *  packed_buf_size - in: size of packed_buf; out: number of chars used.
 */

bool
build_pack_x2reset_response(size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf)
{
	bool rc = true;
    E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));

    ResetResponse_t *resetResponse;
    SuccessfulOutcome_t *successfulOutcome = calloc(1, sizeof(SuccessfulOutcome_t));
    ResetResponse_IEs_t *resetResponse_ie = calloc(1, sizeof(ResetResponse_IEs_t));

    assert(pdu != 0);
    assert(successfulOutcome != 0);
    assert(resetResponse_ie != 0);

    pdu->present = E2AP_PDU_PR_successfulOutcome;
    pdu->choice.successfulOutcome = successfulOutcome;

    successfulOutcome->procedureCode = ProcedureCode_id_reset;
    successfulOutcome->criticality = Criticality_reject;
    successfulOutcome->value.present = SuccessfulOutcome__value_PR_ResetResponse;
    resetResponse = &successfulOutcome->value.choice.ResetResponse;

    resetResponse_ie->id = ProtocolIE_ID_id_CriticalityDiagnostics;
    resetResponse_ie->criticality = Criticality_ignore;
    resetResponse_ie->value.present =  ResetResponse_IEs__value_PR_CriticalityDiagnostics;

    ASN_SEQUENCE_ADD(&resetResponse->protocolIEs, resetResponse_ie);

    CriticalityDiagnostics_IE_List_t *critList = calloc(1, sizeof(CriticalityDiagnostics_IE_List_t));
    assert(critList != 0);

    CriticalityDiagnostics_IE_List__Member *member= calloc(1, sizeof(CriticalityDiagnostics_IE_List__Member));
    assert(member != 0);

    ASN_SEQUENCE_ADD(critList ,member);
    ASN_SEQUENCE_ADD(resetResponse_ie->value.choice.CriticalityDiagnostics.iEsCriticalityDiagnostics, critList);

    rc = per_pack_pdu(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf);
    ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);

    return rc;
}

