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
#include <ProcedureCode.h>
#include <InitiatingMessage.h>
#include <ProtocolIE-ID.h>
#include <x2reset_request_wrapper.h>

/*
 * Build and pack a reset request.
 * Abort the process on allocation failure.
 *  packed_buf_size - in: size of packed_buf; out: number of chars used.
 */

bool
build_pack_x2reset_request(enum Cause_PR cause_group, int cause_value, size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf)
{
	return build_pack_x2reset_request_aux(cause_group, cause_value, packed_buf_size, packed_buf,err_buf_size,err_buf,ATS_ALIGNED_BASIC_PER);

}

bool
build_pack_x2reset_request_aux(enum Cause_PR cause_group, int cause_value, size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf,enum asn_transfer_syntax syntax)
{
	bool rc = true;
	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	InitiatingMessage_t *initiatingMessage = calloc(1, sizeof(InitiatingMessage_t));
	ResetRequest_t	 *resetRequest;

    assert(pdu != 0);
    assert(initiatingMessage != 0);


    pdu->present = E2AP_PDU_PR_initiatingMessage;
    pdu->choice.initiatingMessage = initiatingMessage;

    initiatingMessage->procedureCode = ProcedureCode_id_reset;
    initiatingMessage->criticality = Criticality_reject;
    initiatingMessage->value.present = InitiatingMessage__value_PR_ResetRequest;
    resetRequest = &initiatingMessage->value.choice.ResetRequest;

    ResetRequest_IEs_t *cause_ie = calloc(1, sizeof(ResetRequest_IEs_t));
    assert(cause_ie != 0);
    ASN_SEQUENCE_ADD(&resetRequest->protocolIEs, cause_ie);

    cause_ie->id = ProtocolIE_ID_id_Cause;
    cause_ie->criticality = Criticality_ignore;
    cause_ie->value.present = ResetRequest_IEs__value_PR_Cause;
    Cause_t *cause = &cause_ie->value.choice.Cause;
    cause->present = cause_group;
    switch (cause->present) {
    case Cause_PR_radioNetwork:
    	cause->choice.radioNetwork = cause_value;
    	break;
    case Cause_PR_transport:
    	cause->choice.transport = cause_value;
    	break;
    case Cause_PR_protocol:
    	cause->choice.protocol = cause_value;
    	break;
    case Cause_PR_misc:
    	cause->choice.misc = cause_value;
    	break;
    default:
    	cause->choice.misc = CauseMisc_om_intervention;
    	break;
    }

    rc = pack_pdu_aux(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf,syntax);

    ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
    return rc;
}

