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
#include <UnsuccessfulOutcome.h>
#include <ProtocolIE-ID.h>
#include <ProtocolIE-Field.h>
#include <configuration_update_wrapper.h>

/*
 * Build and pack ENB Configuration Update Acknowledge (successful outcome response).
 * Abort the process on allocation failure.
 *  packed_buf_size - in: size of packed_buf; out: number of chars used.
 */
bool
build_pack_x2enb_configuration_update_ack(
		size_t* packed_buf_size,
        unsigned char* packed_buf,
		size_t err_buf_size,
		char* err_buf)
{
	bool rc = true;
	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	SuccessfulOutcome_t *successfulOutcome = calloc(1, sizeof(SuccessfulOutcome_t));
	ENBConfigurationUpdateAcknowledge_t *enbConfigurationUpdateAcknowledge;
	ENBConfigurationUpdateAcknowledge_IEs_t *enbConfigurationUpdateAcknowledge_IEs = calloc(1, sizeof(ENBConfigurationUpdateAcknowledge_IEs_t));

    assert(pdu != 0);
    assert(successfulOutcome != 0);
    assert(enbConfigurationUpdateAcknowledge_IEs != 0);

    pdu->present = E2AP_PDU_PR_successfulOutcome;
    pdu->choice.successfulOutcome = successfulOutcome;

    successfulOutcome->procedureCode = ProcedureCode_id_eNBConfigurationUpdate;
    successfulOutcome->criticality = Criticality_reject;
    successfulOutcome->value.present = SuccessfulOutcome__value_PR_ENBConfigurationUpdateAcknowledge;
    enbConfigurationUpdateAcknowledge = &successfulOutcome->value.choice.ENBConfigurationUpdateAcknowledge;

    enbConfigurationUpdateAcknowledge_IEs->id = ProtocolIE_ID_id_CriticalityDiagnostics;
	enbConfigurationUpdateAcknowledge_IEs->criticality = Criticality_ignore;
	enbConfigurationUpdateAcknowledge_IEs->value.present =  ENBConfigurationUpdateAcknowledge_IEs__value_PR_CriticalityDiagnostics;

    ASN_SEQUENCE_ADD(&enbConfigurationUpdateAcknowledge->protocolIEs, enbConfigurationUpdateAcknowledge_IEs);

  	rc = per_pack_pdu(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf);

    ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
    return rc;
}

/*
 * Build and pack ENB Configuration Update Failure (unsuccessful outcome message).
 * Abort the process on allocation failure.
 *  packed_buf_size - in: size of packed_buf; out: number of chars used.
 */
bool
build_pack_x2enb_configuration_update_failure(
		size_t* packed_buf_size,
        unsigned char* packed_buf,
        size_t err_buf_size,
		char* err_buf)
{
	bool rc = true;
	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	UnsuccessfulOutcome_t *unsuccessfulOutcome = calloc(1, sizeof(UnsuccessfulOutcome_t));
	ENBConfigurationUpdateFailure_t *enbConfigurationUpdateFailure;
	ENBConfigurationUpdateFailure_IEs_t *enbConfigurationUpdateFailure_IEs = calloc(1, sizeof(ENBConfigurationUpdateFailure_IEs_t));

    assert(pdu != 0);
    assert(unsuccessfulOutcome != 0);
    assert(enbConfigurationUpdateFailure_IEs != 0);


    pdu->present = E2AP_PDU_PR_unsuccessfulOutcome;
    pdu->choice.unsuccessfulOutcome = unsuccessfulOutcome;

    unsuccessfulOutcome->procedureCode = ProcedureCode_id_eNBConfigurationUpdate;
    unsuccessfulOutcome->criticality = Criticality_reject;
    unsuccessfulOutcome->value.present = UnsuccessfulOutcome__value_PR_ENBConfigurationUpdateFailure;
    enbConfigurationUpdateFailure = &unsuccessfulOutcome->value.choice.ENBConfigurationUpdateFailure;

    enbConfigurationUpdateFailure_IEs->id = ProtocolIE_ID_id_Cause;
	enbConfigurationUpdateFailure_IEs->criticality = Criticality_ignore;
	enbConfigurationUpdateFailure_IEs->value.present = ENBConfigurationUpdateFailure_IEs__value_PR_Cause;
	enbConfigurationUpdateFailure_IEs->value.choice.Cause.present = Cause_PR_protocol;
	enbConfigurationUpdateFailure_IEs->value.choice.Cause.choice.protocol= CauseProtocol_abstract_syntax_error_reject;
    ASN_SEQUENCE_ADD(&enbConfigurationUpdateFailure->protocolIEs, enbConfigurationUpdateFailure_IEs);


    rc = per_pack_pdu(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf);

    ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
    return rc;
}

/*
 * Build and pack ENDC Configuration Update Acknowledge (successful outcome response).
 * Abort the process on allocation failure.
 *  packed_buf_size - in: size of packed_buf; out: number of chars used.
 */
bool
build_pack_endc_configuration_update_ack(
		size_t* packed_buf_size,
		unsigned char* packed_buf,
		size_t err_buf_size,
		char* err_buf)
{
	bool rc = true;
	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	SuccessfulOutcome_t *successfulOutcome = calloc(1, sizeof(SuccessfulOutcome_t));
	ENDCConfigurationUpdateAcknowledge_t *endcConfigurationUpdateAcknowledge;
	ENDCConfigurationUpdateAcknowledge_IEs_t *endcConfigurationUpdateAcknowledge_IEs = calloc(1, sizeof(ENDCConfigurationUpdateAcknowledge_IEs_t));

    assert(pdu != 0);
    assert(successfulOutcome != 0);
    assert(endcConfigurationUpdateAcknowledge_IEs != 0);

    pdu->present = E2AP_PDU_PR_successfulOutcome;
    pdu->choice.successfulOutcome = successfulOutcome;

    successfulOutcome->procedureCode = ProcedureCode_id_endcConfigurationUpdate;
    successfulOutcome->criticality = Criticality_reject;
    successfulOutcome->value.present = SuccessfulOutcome__value_PR_ENDCConfigurationUpdateAcknowledge;
    endcConfigurationUpdateAcknowledge = &successfulOutcome->value.choice.ENDCConfigurationUpdateAcknowledge;
    ASN_SEQUENCE_ADD(&endcConfigurationUpdateAcknowledge->protocolIEs, endcConfigurationUpdateAcknowledge_IEs);

    endcConfigurationUpdateAcknowledge_IEs->id = ProtocolIE_ID_id_RespondingNodeType_EndcConfigUpdate;
	endcConfigurationUpdateAcknowledge_IEs->criticality = Criticality_reject;
	endcConfigurationUpdateAcknowledge_IEs->value.present = ENDCConfigurationUpdateAcknowledge_IEs__value_PR_RespondingNodeType_EndcConfigUpdate;
	endcConfigurationUpdateAcknowledge_IEs->value.choice.RespondingNodeType_EndcConfigUpdate.present = RespondingNodeType_EndcConfigUpdate_PR_respond_eNB;

	ProtocolIE_Container_119P95_t *enb_ENDCConfigUpdateAckIEs_Container = calloc(1, sizeof(ProtocolIE_Container_119P95_t));
	assert(enb_ENDCConfigUpdateAckIEs_Container != 0);
	endcConfigurationUpdateAcknowledge_IEs->value.choice.RespondingNodeType_EndcConfigUpdate.choice.respond_eNB = (struct ProtocolIE_Container*)enb_ENDCConfigUpdateAckIEs_Container;

	//Leave the respond_eNB container empty (ENB_ENDCConfigUpdateAckIEs_t is an empty element).

    rc = per_pack_pdu(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf);

    ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
    return rc;
}

/*
 * Build and pack ENDC Configuration Update Failure (unsuccessful outcome message).
 * Abort the process on allocation failure.
 *  packed_buf_size - in: size of packed_buf; out: number of chars used.
 */
bool
build_pack_endc_configuration_update_failure(
		size_t* packed_buf_size,
        unsigned char* packed_buf,
        size_t err_buf_size,
		char* err_buf)
{
	bool rc = true;
	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	UnsuccessfulOutcome_t *unsuccessfulOutcome = calloc(1, sizeof(UnsuccessfulOutcome_t));
	ENDCConfigurationUpdateFailure_t *endcConfigurationUpdateFailure;
	ENDCConfigurationUpdateFailure_IEs_t *endcConfigurationUpdateFailure_IEs = calloc(1, sizeof(ENDCConfigurationUpdateFailure_IEs_t));

    assert(pdu != 0);
    assert(unsuccessfulOutcome != 0);
    assert(endcConfigurationUpdateFailure_IEs != 0);


    pdu->present = E2AP_PDU_PR_unsuccessfulOutcome;
    pdu->choice.unsuccessfulOutcome = unsuccessfulOutcome;

    unsuccessfulOutcome->procedureCode = ProcedureCode_id_endcConfigurationUpdate;
    unsuccessfulOutcome->criticality = Criticality_reject;
    unsuccessfulOutcome->value.present = UnsuccessfulOutcome__value_PR_ENDCConfigurationUpdateFailure;
    endcConfigurationUpdateFailure = &unsuccessfulOutcome->value.choice.ENDCConfigurationUpdateFailure;

    endcConfigurationUpdateFailure_IEs->id = ProtocolIE_ID_id_Cause;
    endcConfigurationUpdateFailure_IEs->criticality = Criticality_ignore;
    endcConfigurationUpdateFailure_IEs->value.present = ENDCConfigurationUpdateFailure_IEs__value_PR_Cause;
    endcConfigurationUpdateFailure_IEs->value.choice.Cause.present = Cause_PR_protocol;
    endcConfigurationUpdateFailure_IEs->value.choice.Cause.choice.protocol= CauseProtocol_abstract_syntax_error_reject;
    ASN_SEQUENCE_ADD(&endcConfigurationUpdateFailure->protocolIEs, endcConfigurationUpdateFailure_IEs);

    rc = per_pack_pdu(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf);

    ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
    return rc;
}

