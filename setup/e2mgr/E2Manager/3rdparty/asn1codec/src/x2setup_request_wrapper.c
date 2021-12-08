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
#include <InitiatingMessage.h>
#include <X2SetupRequest.h>
#include <GlobalENB-ID.h>
#include <PLMN-Identity.h>
#include <ENB-ID.h>
#include <FDD-Info.h>
#include <ServedCells.h>
#include <ProtocolIE-ID.h>
#include <ProtocolIE-Field.h>
#include <x2setup_request_wrapper.h>

static void assignPLMN_Identity (PLMN_Identity_t *pLMN_Identity, uint8_t const* pLMNId);
static void assignENB_ID(GlobalENB_ID_t *globalENB_ID,uint8_t const* eNBId, unsigned int bitqty);
static void assignServedCell_Information(ServedCell_Information_t *servedCell_Information,uint8_t const* pLMN_Identity, uint8_t const* eNBId, unsigned int bitqty,uint8_t const *ric_flag);

/*
 * Build and pack X2 setup request.
 * Abort the process on allocation failure.
 *  packed_buf_size - in: size of packed_buf; out: number of chars used.
 */

bool
build_pack_x2setup_request(
		uint8_t const* pLMN_Identity, uint8_t const* eNBId, unsigned int bitqty /*18, 20, 21, 28*/, uint8_t const *ric_flag,
		size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf
)
{
	return build_pack_x2setup_request_aux(
			pLMN_Identity, eNBId, bitqty, ric_flag,
			packed_buf_size, packed_buf,err_buf_size,err_buf,ATS_ALIGNED_BASIC_PER);

}

bool
build_pack_x2setup_request_aux(
		uint8_t const* pLMN_Identity, uint8_t const* eNBId, unsigned int bitqty /*18, 20, 21, 28*/, uint8_t const *ric_flag,
		size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf,enum asn_transfer_syntax syntax
)
{
	bool rc = true;
	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	InitiatingMessage_t *initiatingMessage = calloc(1, sizeof(InitiatingMessage_t));
	X2SetupRequest_t *x2SetupRequest;

    assert(pdu != 0);
    assert(initiatingMessage != 0);


    pdu->present = E2AP_PDU_PR_initiatingMessage;
    pdu->choice.initiatingMessage = initiatingMessage;

    initiatingMessage->procedureCode = ProcedureCode_id_x2Setup;
    initiatingMessage->criticality = Criticality_reject;
    initiatingMessage->value.present = InitiatingMessage__value_PR_X2SetupRequest;
    x2SetupRequest = &initiatingMessage->value.choice.X2SetupRequest;

    X2SetupRequest_IEs_t *globalENB_ID_ie = calloc(1, sizeof(X2SetupRequest_IEs_t));
    assert(globalENB_ID_ie != 0);
    ASN_SEQUENCE_ADD(&x2SetupRequest->protocolIEs, globalENB_ID_ie);

    globalENB_ID_ie->id = ProtocolIE_ID_id_GlobalENB_ID;
    globalENB_ID_ie->criticality = Criticality_reject;
    globalENB_ID_ie->value.present = X2SetupRequest_IEs__value_PR_GlobalENB_ID;
	GlobalENB_ID_t *globalENB_ID = &globalENB_ID_ie->value.choice.GlobalENB_ID;

	assignPLMN_Identity(&globalENB_ID->pLMN_Identity, pLMN_Identity);
	assignENB_ID(globalENB_ID, eNBId, bitqty);

    X2SetupRequest_IEs_t *servedCells_ie = calloc(1, sizeof(X2SetupRequest_IEs_t));
    assert(servedCells_ie != 0);
    ASN_SEQUENCE_ADD(&x2SetupRequest->protocolIEs, servedCells_ie);

    servedCells_ie->id = ProtocolIE_ID_id_ServedCells;
    servedCells_ie->criticality = Criticality_reject;
    servedCells_ie->value.present = X2SetupRequest_IEs__value_PR_ServedCells;

    ServedCells__Member *servedCells__Member = calloc(1,sizeof(ServedCells__Member));
    assert(servedCells__Member !=0);
    ASN_SEQUENCE_ADD(&servedCells_ie->value.choice.ServedCells, servedCells__Member);

    assignServedCell_Information(&servedCells__Member->servedCellInfo, pLMN_Identity,eNBId, bitqty,ric_flag);

    rc = pack_pdu_aux(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf,syntax);

    ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
    return rc;
}

static void assignPLMN_Identity (PLMN_Identity_t *pLMN_Identity, uint8_t const* pLMNId)
{
	pLMN_Identity->size = pLMN_Identity_size;
	pLMN_Identity->buf = calloc(1,pLMN_Identity->size);
	assert(pLMN_Identity->buf != 0);
	memcpy(pLMN_Identity->buf, pLMNId, pLMN_Identity->size);
}

/*
 * Calculate and assign the value of ENB_ID.
 * Abort the process on allocation failure.
 */
static void assignENB_ID(GlobalENB_ID_t *globalENB_ID,uint8_t const* eNBId, unsigned int bitqty)
{
	size_t size_in_bytes = (bitqty / 8) + ((bitqty % 8) > 0);
	int unused_bits = 8 - (bitqty % 8);
	uint8_t *tbuf;
	switch (bitqty){
	case shortMacro_eNB_ID_size:
		globalENB_ID->eNB_ID.present = ENB_ID_PR_short_Macro_eNB_ID;
		globalENB_ID->eNB_ID.choice.short_Macro_eNB_ID.size = size_in_bytes;
		globalENB_ID->eNB_ID.choice.short_Macro_eNB_ID.bits_unused = unused_bits;
		tbuf = globalENB_ID->eNB_ID.choice.short_Macro_eNB_ID.buf = calloc(1, size_in_bytes);
		assert(globalENB_ID->eNB_ID.choice.short_Macro_eNB_ID.buf  != 0);
		memcpy(globalENB_ID->eNB_ID.choice.short_Macro_eNB_ID.buf,eNBId, size_in_bytes) ;
		tbuf[size_in_bytes - 1] <<= unused_bits;
		break;
	case macro_eNB_ID_size:
		globalENB_ID->eNB_ID.present =ENB_ID_PR_macro_eNB_ID;
		globalENB_ID->eNB_ID.choice.macro_eNB_ID.size = size_in_bytes;
		globalENB_ID->eNB_ID.choice.macro_eNB_ID.bits_unused = unused_bits;
		tbuf = globalENB_ID->eNB_ID.choice.macro_eNB_ID.buf = calloc(1, size_in_bytes);
		assert(globalENB_ID->eNB_ID.choice.macro_eNB_ID.buf != 0);
		memcpy(globalENB_ID->eNB_ID.choice.macro_eNB_ID.buf,eNBId,size_in_bytes);
		tbuf[size_in_bytes - 1] <<= unused_bits;
		break;
	case longMacro_eNB_ID_size:
		globalENB_ID->eNB_ID.present =ENB_ID_PR_long_Macro_eNB_ID;
		globalENB_ID->eNB_ID.choice.long_Macro_eNB_ID.size = size_in_bytes;
		globalENB_ID->eNB_ID.choice.long_Macro_eNB_ID.bits_unused = unused_bits;
		tbuf = globalENB_ID->eNB_ID.choice.long_Macro_eNB_ID.buf = calloc(1, size_in_bytes);
		assert(globalENB_ID->eNB_ID.choice.long_Macro_eNB_ID.buf != 0);
		memcpy(globalENB_ID->eNB_ID.choice.long_Macro_eNB_ID.buf,eNBId,size_in_bytes);
		tbuf[size_in_bytes - 1] <<= unused_bits;
		break;
	case home_eNB_ID_size:
		globalENB_ID->eNB_ID.present = ENB_ID_PR_home_eNB_ID;
		globalENB_ID->eNB_ID.choice.home_eNB_ID.size = size_in_bytes;
		globalENB_ID->eNB_ID.choice.home_eNB_ID.bits_unused =unused_bits;
		tbuf = globalENB_ID->eNB_ID.choice.home_eNB_ID.buf = calloc(1,size_in_bytes);
		assert(globalENB_ID->eNB_ID.choice.home_eNB_ID.buf != 0);
		memcpy(globalENB_ID->eNB_ID.choice.home_eNB_ID.buf,eNBId,size_in_bytes);
		tbuf[size_in_bytes - 1] <<= unused_bits;
		break;
	default:
		break;
	}

}

/*
 * Calculate and assign the value of ServedCell_Information.
 * Abort the process on allocation failure.
 */
static void assignServedCell_Information(
		ServedCell_Information_t *servedCell_Information,
		uint8_t const* pLMN_Identity,
		uint8_t const* eNBId,
		unsigned int bitqty,
		uint8_t const *ric_flag)
{
	size_t size_in_bytes = 	(eUTRANcellIdentifier_size / 8) + ((eUTRANcellIdentifier_size % 8) > 0);
	int unused_bits = 8 - (eUTRANcellIdentifier_size % 8);
	size_t bitqty_size_in_bytes = (bitqty / 8) + ((bitqty % 8) > 0);
	int bitqty_unused_bits = 8 - (bitqty % 8);

	servedCell_Information->pCI = 503;
	assignPLMN_Identity(&servedCell_Information->cellId.pLMN_Identity, pLMN_Identity);

	servedCell_Information->cellId.eUTRANcellIdentifier.size = size_in_bytes;
	servedCell_Information->cellId.eUTRANcellIdentifier.bits_unused = unused_bits;
	servedCell_Information->cellId.eUTRANcellIdentifier.buf = calloc(1,servedCell_Information->cellId.eUTRANcellIdentifier.size);
	assert(servedCell_Information->cellId.eUTRANcellIdentifier.buf != 0);
	memcpy(servedCell_Information->cellId.eUTRANcellIdentifier.buf, eNBId, bitqty_size_in_bytes);
	if (bitqty < eUTRANcellIdentifier_size)	{
		servedCell_Information->cellId.eUTRANcellIdentifier.buf[bitqty_size_in_bytes - 1] <<= bitqty_unused_bits;
	} else {
		servedCell_Information->cellId.eUTRANcellIdentifier.buf[size_in_bytes - 1] <<= unused_bits;
	}

	servedCell_Information->tAC.size = 2;
	servedCell_Information->tAC.buf = calloc(1,servedCell_Information->tAC.size);
	assert(servedCell_Information->tAC.buf != 0);


	PLMN_Identity_t *broadcastPLMN_Identity = calloc(1, sizeof(PLMN_Identity_t));
	assert(broadcastPLMN_Identity != 0);
	ASN_SEQUENCE_ADD(&servedCell_Information->broadcastPLMNs, broadcastPLMN_Identity);

	assignPLMN_Identity(broadcastPLMN_Identity, pLMN_Identity); //ric_flag: disabled because a real eNB rejects the message

	servedCell_Information->eUTRA_Mode_Info.present= EUTRA_Mode_Info_PR_fDD;
	servedCell_Information->eUTRA_Mode_Info.choice.fDD = calloc(1, sizeof(FDD_Info_t));
	assert(servedCell_Information->eUTRA_Mode_Info.choice.fDD != 0);
	servedCell_Information->eUTRA_Mode_Info.choice.fDD->uL_EARFCN = 0;
	servedCell_Information->eUTRA_Mode_Info.choice.fDD->dL_EARFCN = 0;
	servedCell_Information->eUTRA_Mode_Info.choice.fDD->uL_Transmission_Bandwidth = Transmission_Bandwidth_bw6;
	servedCell_Information->eUTRA_Mode_Info.choice.fDD->dL_Transmission_Bandwidth = Transmission_Bandwidth_bw15;
}

/* Build and pack X2 setup request.
 * Abort the process on allocation failure.
 * packed_buf_size - in: size of packed_buf; out: number of chars used.
 */
bool
build_pack_endc_x2setup_request(
		uint8_t const* pLMN_Identity, uint8_t const* eNBId, unsigned int bitqty /*18, 20, 21, 28*/, uint8_t const *ric_flag,
		size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf
)
{
	return build_pack_endc_x2setup_request_aux(
			pLMN_Identity, eNBId, bitqty, ric_flag,
			packed_buf_size, packed_buf,err_buf_size,  err_buf,ATS_ALIGNED_BASIC_PER
	);
}

bool
build_pack_endc_x2setup_request_aux(
		uint8_t const* pLMN_Identity, uint8_t const* eNBId, unsigned int bitqty /*18, 20, 21, 28*/, uint8_t const *ric_flag,
		size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf,enum asn_transfer_syntax syntax
)
{
	bool rc = true;
	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	InitiatingMessage_t *initiatingMessage = calloc(1, sizeof(InitiatingMessage_t));
	ENDCX2SetupRequest_t *endcX2SetupRequest;

    assert(pdu != 0);
    assert(initiatingMessage != 0);

    pdu->present = E2AP_PDU_PR_initiatingMessage;
    pdu->choice.initiatingMessage = initiatingMessage;

    initiatingMessage->procedureCode = ProcedureCode_id_endcX2Setup;
    initiatingMessage->criticality = Criticality_reject;
    initiatingMessage->value.present = InitiatingMessage__value_PR_ENDCX2SetupRequest;
    endcX2SetupRequest = &initiatingMessage->value.choice.ENDCX2SetupRequest;
    ENDCX2SetupRequest_IEs_t *endcX2SetupRequest_IEs = calloc(1, sizeof(ENDCX2SetupRequest_IEs_t));
    assert(endcX2SetupRequest_IEs != 0);
	ASN_SEQUENCE_ADD(&endcX2SetupRequest->protocolIEs, endcX2SetupRequest_IEs);
	endcX2SetupRequest_IEs->id = ProtocolIE_ID_id_InitiatingNodeType_EndcX2Setup;
	endcX2SetupRequest_IEs->criticality = Criticality_reject;
	endcX2SetupRequest_IEs->value.present = ENDCX2SetupRequest_IEs__value_PR_InitiatingNodeType_EndcX2Setup;
	endcX2SetupRequest_IEs->value.choice.InitiatingNodeType_EndcX2Setup.present = InitiatingNodeType_EndcX2Setup_PR_init_eNB;

	ProtocolIE_Container_119P85_t *enb_ENDCX2SetupReqIE_Container = calloc(1, sizeof(ProtocolIE_Container_119P85_t));
	assert(enb_ENDCX2SetupReqIE_Container != 0);
	endcX2SetupRequest_IEs->value.choice.InitiatingNodeType_EndcX2Setup.choice.init_eNB = (struct ProtocolIE_Container*)enb_ENDCX2SetupReqIE_Container;
	ENB_ENDCX2SetupReqIEs_t *globalENB_ID_ie = calloc(1, sizeof(ENB_ENDCX2SetupReqIEs_t));
	assert(globalENB_ID_ie != 0);
	ASN_SEQUENCE_ADD(enb_ENDCX2SetupReqIE_Container,globalENB_ID_ie);
	globalENB_ID_ie->id = ProtocolIE_ID_id_GlobalENB_ID;
	globalENB_ID_ie->criticality = Criticality_reject;
	globalENB_ID_ie->value.present = ENB_ENDCX2SetupReqIEs__value_PR_GlobalENB_ID;

	GlobalENB_ID_t *globalENB_ID = &globalENB_ID_ie->value.choice.GlobalENB_ID;
	assignPLMN_Identity(&globalENB_ID->pLMN_Identity, pLMN_Identity);
	assignENB_ID(globalENB_ID, eNBId, bitqty);


	ENB_ENDCX2SetupReqIEs_t *ServedEUTRAcellsENDCX2ManagementList_ie = calloc(1, sizeof(ENB_ENDCX2SetupReqIEs_t));
	assert(ServedEUTRAcellsENDCX2ManagementList_ie != 0);
    ASN_SEQUENCE_ADD(enb_ENDCX2SetupReqIE_Container, ServedEUTRAcellsENDCX2ManagementList_ie);

    ServedEUTRAcellsENDCX2ManagementList_ie->id = ProtocolIE_ID_id_ServedEUTRAcellsENDCX2ManagementList;
	ServedEUTRAcellsENDCX2ManagementList_ie->criticality = Criticality_reject;
	ServedEUTRAcellsENDCX2ManagementList_ie->value.present = ENB_ENDCX2SetupReqIEs__value_PR_ServedEUTRAcellsENDCX2ManagementList;


	ServedEUTRAcellsENDCX2ManagementList__Member *servedEUTRAcellsENDCX2ManagementList__Member = calloc(1, sizeof(ServedEUTRAcellsENDCX2ManagementList__Member));
	assert(servedEUTRAcellsENDCX2ManagementList__Member != 0);
	ASN_SEQUENCE_ADD(&ServedEUTRAcellsENDCX2ManagementList_ie->value.choice.ServedEUTRAcellsENDCX2ManagementList, servedEUTRAcellsENDCX2ManagementList__Member);

	assignServedCell_Information(&servedEUTRAcellsENDCX2ManagementList__Member->servedEUTRACellInfo, pLMN_Identity, eNBId, bitqty,ric_flag);

    rc = pack_pdu_aux(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf, syntax);

    ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
    return rc;
}
