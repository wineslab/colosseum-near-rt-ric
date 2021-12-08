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

//
// Created by adi ENZEL on 12/10/19.
//

#ifndef E2_E2BUILDER_H
#define E2_E2BUILDER_H

#include <cstring>
#include <cstdio>
#include <cerrno>
#include <cstdlib>
#include <sys/types.h>
#include <error.h>
#include <algorithm>
#include <3rdparty/oranE2SM/E2SM-gNB-NRT-RANfunction-Definition.h>
#include <3rdparty/oranE2SM/RIC-InsertStyle-List.h>
#include <3rdparty/oranE2SM/RANparameterDef-Item.h>
#include <3rdparty/oranE2/GlobalE2node-en-gNB-ID.h>


//#include <mdclog/mdclog.h>


#include "oranE2/E2AP-PDU.h"
#include "oranE2/InitiatingMessage.h"
#include "oranE2/SuccessfulOutcome.h"
#include "oranE2/UnsuccessfulOutcome.h"

#include "oranE2/ProtocolIE-Field.h"
#include "oranE2/ENB-ID.h"
#include "oranE2/GlobalENB-ID.h"
#include "oranE2/GlobalE2node-gNB-ID.h"
#include "oranE2/constr_TYPE.h"
#include "oranE2/asn_constant.h"

using namespace std;

#define printEntry(type, function);  fprintf(stdout, "start Test %s , %s", type, function);


static void checkAndPrint(asn_TYPE_descriptor_t *typeDescriptor, void *data, char *dataType, const char *function) {
    char errbuf[128]; /* Buffer for error message */
    size_t errlen = sizeof(errbuf); /* Size of the buffer */
    if (asn_check_constraints(typeDescriptor, data, errbuf, &errlen) != 0) {
        fprintf(stderr, "%s Constraint validation failed: %s", dataType, errbuf);
    }
    fprintf(stdout, "%s successes function %s", dataType, function);
}

void createPLMN_IDByMCCandMNC(PLMN_Identity_t *plmnId, int mcc, int mnc) {

    //printEntry("PLMN_Identity_t", __func__)

    ASN_STRUCT_RESET(asn_DEF_PLMN_Identity, plmnId);
    plmnId->size = 3;
    plmnId->buf = (uint8_t *) calloc(1, 3);
    volatile auto mcc1 = (unsigned int) (mcc / 100);
    volatile auto mcc2 = (unsigned int) (mcc / 10 % 10);
    volatile auto mcc3 = (unsigned int) (mcc % 10);
    plmnId->buf[0] = mcc2 << 4 | mcc1;

    volatile auto mnc1 = (unsigned int)0;
    volatile auto mnc2 = (unsigned int)0;
    volatile auto mnc3 = (unsigned int)0;

    if (mnc >= 100) {
        mnc1 = (unsigned int) (mnc / 100);
        mnc2 = (unsigned int) (mnc / 10 % 10);
        mnc3 = (unsigned int) (mnc % 10);
    } else {
        mnc1 = (unsigned int) (mnc / 10);
        mnc2 = (unsigned int) (mnc % 10);
        mnc3 = 15;
    }
    plmnId->buf[1] = mcc3 << 4 | mnc3 ;
    plmnId->buf[2] = mnc2 << 4 | mnc1 ;

    //checkAndPrint(&asn_DEF_PLMN_Identity, plmnId, (char *) "PLMN_Identity_t", __func__);

}


PLMN_Identity_t *createPLMN_ID(const unsigned char *data) {
    printEntry("PLMN_Identity_t", __func__)
    auto *plmnId = (PLMN_Identity_t *) calloc(1, sizeof(PLMN_Identity_t));
    ASN_STRUCT_RESET(asn_DEF_PLMN_Identity, plmnId);
    plmnId->size = 3;
    plmnId->buf = (uint8_t *) calloc(1, 3);
    memcpy(plmnId->buf, data, 3);

    checkAndPrint(&asn_DEF_PLMN_Identity, plmnId, (char *) "PLMN_Identity_t", __func__);

    return plmnId;
}

BIT_STRING_t *createBIT_STRING(int size, int unusedBits, uint8_t *data) {
    printEntry("BIT_STRING_t", __func__)
    auto *bitString = (BIT_STRING_t *) calloc(1, sizeof(BIT_STRING_t));
    ASN_STRUCT_RESET(asn_DEF_BIT_STRING, bitString);
    bitString->size = size;
    bitString->bits_unused = unusedBits;
    bitString->buf = (uint8_t *) calloc(1, size);
    // set bits to zero
    data[bitString->size - 1] = ((unsigned) (data[bitString->size - 1] >>
                                                                       (unsigned) bitString->bits_unused)
            << (unsigned) bitString->bits_unused);
    memcpy(bitString->buf, data, size);

    checkAndPrint(&asn_DEF_BIT_STRING, bitString, (char *) "BIT_STRING_t", __func__);

    return bitString;
}


OCTET_STRING_t *createOCTET_STRING(const unsigned char *data, int size) {
    printEntry("OCTET_STRING_t", __func__)
    auto *octs = (PLMN_Identity_t *) calloc(1, sizeof(PLMN_Identity_t));
    ASN_STRUCT_RESET(asn_DEF_OCTET_STRING, octs);
    octs->size = size;
    octs->buf = (uint8_t *) calloc(1, size);
    memcpy(octs->buf, data, size);

    checkAndPrint(&asn_DEF_OCTET_STRING, octs, (char *) "OCTET_STRING_t", __func__);
    return octs;
}


ENB_ID_t *createENB_ID(ENB_ID_PR enbType, unsigned char *data) {
    printEntry("ENB_ID_t", __func__)
    auto *enb = (ENB_ID_t *) calloc(1, sizeof(ENB_ID_t));
    ASN_STRUCT_RESET(asn_DEF_ENB_ID, enb);

    enb->present = enbType;

    switch (enbType) {
        case ENB_ID_PR_macro_eNB_ID: { // 20 bit 3 bytes
            enb->choice.macro_eNB_ID.size = 3;
            enb->choice.macro_eNB_ID.bits_unused = 4;

            enb->present = ENB_ID_PR_macro_eNB_ID;

            enb->choice.macro_eNB_ID.buf = (uint8_t *) calloc(1, enb->choice.macro_eNB_ID.size);
            data[enb->choice.macro_eNB_ID.size - 1] = ((unsigned) (data[enb->choice.macro_eNB_ID.size - 1] >>
                                                                                                           (unsigned) enb->choice.macro_eNB_ID.bits_unused)
                    << (unsigned) enb->choice.macro_eNB_ID.bits_unused);
            memcpy(enb->choice.macro_eNB_ID.buf, data, enb->choice.macro_eNB_ID.size);

            break;
        }
        case ENB_ID_PR_home_eNB_ID: { // 28 bit 4 bytes
            enb->choice.home_eNB_ID.size = 4;
            enb->choice.home_eNB_ID.bits_unused = 4;
            enb->present = ENB_ID_PR_home_eNB_ID;

            enb->choice.home_eNB_ID.buf = (uint8_t *) calloc(1, enb->choice.home_eNB_ID.size);
            data[enb->choice.home_eNB_ID.size - 1] = ((unsigned) (data[enb->choice.home_eNB_ID.size - 1] >>
                                                                                                         (unsigned) enb->choice.home_eNB_ID.bits_unused)
                    << (unsigned) enb->choice.home_eNB_ID.bits_unused);
            memcpy(enb->choice.home_eNB_ID.buf, data, enb->choice.home_eNB_ID.size);
            break;
        }
        case ENB_ID_PR_short_Macro_eNB_ID: { // 18 bit - 3 bytes
            enb->choice.short_Macro_eNB_ID.size = 3;
            enb->choice.short_Macro_eNB_ID.bits_unused = 6;
            enb->present = ENB_ID_PR_short_Macro_eNB_ID;

            enb->choice.short_Macro_eNB_ID.buf = (uint8_t *) calloc(1, enb->choice.short_Macro_eNB_ID.size);
            data[enb->choice.short_Macro_eNB_ID.size - 1] = (
                    (unsigned) (data[enb->choice.short_Macro_eNB_ID.size - 1] >>
                                                                              (unsigned) enb->choice.short_Macro_eNB_ID.bits_unused)
                            << (unsigned) enb->choice.short_Macro_eNB_ID.bits_unused);
            memcpy(enb->choice.short_Macro_eNB_ID.buf, data, enb->choice.short_Macro_eNB_ID.size);
            break;
        }
        case ENB_ID_PR_long_Macro_eNB_ID: { // 21
            enb->choice.long_Macro_eNB_ID.size = 3;
            enb->choice.long_Macro_eNB_ID.bits_unused = 3;
            enb->present = ENB_ID_PR_long_Macro_eNB_ID;

            enb->choice.long_Macro_eNB_ID.buf = (uint8_t *) calloc(1, enb->choice.long_Macro_eNB_ID.size);
            data[enb->choice.long_Macro_eNB_ID.size - 1] =
                    ((unsigned) (data[enb->choice.long_Macro_eNB_ID.size - 1] >> (unsigned) enb->choice.long_Macro_eNB_ID.bits_unused)
                            << (unsigned) enb->choice.long_Macro_eNB_ID.bits_unused);
            memcpy(enb->choice.long_Macro_eNB_ID.buf, data, enb->choice.long_Macro_eNB_ID.size);
            break;
        }
        default:
            free(enb);
            return nullptr;
    }

    checkAndPrint(&asn_DEF_ENB_ID, enb, (char *) "ENB_ID_t", __func__);
    return enb;
}

GlobalENB_ID_t *createGlobalENB_ID(PLMN_Identity_t *plmnIdentity, ENB_ID_t *enbId) {
    printEntry("GlobalENB_ID_t", __func__)
    auto *genbId = (GlobalENB_ID_t *) calloc(1, sizeof(GlobalENB_ID_t));
    ASN_STRUCT_RESET(asn_DEF_GlobalENB_ID, genbId);
    memcpy(&genbId->pLMN_Identity, plmnIdentity, sizeof(PLMN_Identity_t));
    memcpy(&genbId->eNB_ID, enbId, sizeof(ENB_ID_t));

    checkAndPrint(&asn_DEF_GlobalENB_ID, genbId, (char *) "GlobalENB_ID_t", __func__);
    return genbId;
}


//static void buildInitiatingMessagePDU(E2AP_PDU_t &pdu, InitiatingMessage_t *initMsg) {
//    pdu.present = E2AP_PDU_PR_initiatingMessage;
//    pdu.choice.initiatingMessage = initMsg;
//}

template<typename T>
static void buildInitMsg(InitiatingMessage_t &initMsg,
                         InitiatingMessage__value_PR present,
                         ProcedureCode_t procedureCode,
                         Criticality_t criticality,
                         T *value) {
    initMsg.value.present = present;
    initMsg.procedureCode = procedureCode;
    initMsg.criticality = criticality;

    switch (present) {
        case InitiatingMessage__value_PR_RICsubscriptionRequest: {
            memcpy(&initMsg.value.choice.RICsubscriptionRequest, value, sizeof(*value));
            break;
        }
        case InitiatingMessage__value_PR_RICsubscriptionDeleteRequest: {
            memcpy(&initMsg.value.choice.RICsubscriptionDeleteRequest, value, sizeof(*value));
            break;
        }
        case InitiatingMessage__value_PR_RICserviceUpdate: {
            memcpy(&initMsg.value.choice.RICserviceUpdate, value, sizeof(*value));
            break;
        }
        case InitiatingMessage__value_PR_RICcontrolRequest: {
            memcpy(&initMsg.value.choice.RICcontrolRequest, value, sizeof(*value));
            break;
        }
        case InitiatingMessage__value_PR_RICindication: {
            memcpy(&initMsg.value.choice.RICindication, value, sizeof(*value));
            break;
        }
        case InitiatingMessage__value_PR_RICserviceQuery: {
            memcpy(&initMsg.value.choice.RICserviceQuery, value, sizeof(*value));
            break;
        }
        case InitiatingMessage__value_PR_NOTHING:
        default : {
            break;
        }
    }
}

//static void buildSuccsesfulMessagePDU(E2AP_PDU_t &pdu, SuccessfulOutcome_t *succMsg) {
//    pdu.present = E2AP_PDU_PR_successfulOutcome;
//    pdu.choice.successfulOutcome = succMsg;
//}

template<typename T>
static void buildSuccMsg(SuccessfulOutcome_t &succMsg,
                         SuccessfulOutcome__value_PR present,
                         ProcedureCode_t procedureCode,
                         Criticality_t criticality,
                         T *value) {
    succMsg.value.present = present;
    succMsg.procedureCode = procedureCode;
    succMsg.criticality = criticality;

    switch (present) {
        case SuccessfulOutcome__value_PR_RICsubscriptionResponse: {
            memcpy(&succMsg.value.choice.RICsubscriptionResponse, value, sizeof(*value));
            break;
        }
        case SuccessfulOutcome__value_PR_RICsubscriptionDeleteResponse: {
            memcpy(&succMsg.value.choice.RICsubscriptionDeleteResponse, value, sizeof(*value));
            break;
        }
        case SuccessfulOutcome__value_PR_RICserviceUpdateAcknowledge: {
            memcpy(&succMsg.value.choice.RICserviceUpdateAcknowledge, value, sizeof(*value));
            break;
        }
        case SuccessfulOutcome__value_PR_RICcontrolAcknowledge: {
            memcpy(&succMsg.value.choice.RICcontrolAcknowledge, value, sizeof(*value));
            break;
        }
        case SuccessfulOutcome__value_PR_ResetResponse: {
            memcpy(&succMsg.value.choice.ResetResponse, value, sizeof(*value));
            break;
        }
        case SuccessfulOutcome__value_PR_NOTHING:
        default:
            break;
    }
}


//static void buildUnSucssesfullMessagePDU(E2AP_PDU_t &pdu, UnsuccessfulOutcome_t *unSuccMsg) {
//    pdu.present = E2AP_PDU_PR_unsuccessfulOutcome;
//    pdu.choice.unsuccessfulOutcome = unSuccMsg;
//}

template<typename T>
static void buildUnSuccMsg(UnsuccessfulOutcome_t &unSuccMsg,
                           UnsuccessfulOutcome__value_PR present,
                           ProcedureCode_t procedureCode,
                           Criticality_t criticality,
                           T *value) {
    unSuccMsg.value.present = present;
    unSuccMsg.procedureCode = procedureCode;
    unSuccMsg.criticality = criticality;

    switch (present) {
        case UnsuccessfulOutcome__value_PR_RICsubscriptionFailure: {
            memcpy(&unSuccMsg.value.choice.RICsubscriptionFailure, value, sizeof(*value));
            break;
        }
        case UnsuccessfulOutcome__value_PR_RICsubscriptionDeleteFailure: {
            memcpy(&unSuccMsg.value.choice.RICsubscriptionDeleteFailure, value, sizeof(*value));
            break;
        }
        case UnsuccessfulOutcome__value_PR_RICserviceUpdateFailure: {
            memcpy(&unSuccMsg.value.choice.RICserviceUpdateFailure, value, sizeof(*value));
            break;
        }
        case UnsuccessfulOutcome__value_PR_RICcontrolFailure: {
            memcpy(&unSuccMsg.value.choice.RICcontrolFailure, value, sizeof(*value));
            break;

        }
        case UnsuccessfulOutcome__value_PR_NOTHING:
        default:
            break;
    }
}


//static void createPLMN_ID(PLMN_Identity_t &plmnId, const unsigned char *data) {
//    //printEntry("PLMN_Identity_t", __func__)
//    //PLMN_Identity_t *plmnId = calloc(1, sizeof(PLMN_Identity_t));
//    ASN_STRUCT_RESET(asn_DEF_PLMN_Identity, &plmnId);
//    plmnId.size = 3;//    uint64_t st = 0;
////    uint32_t aux1 = 0;
////    st = rdtscp(aux1);
//
//    plmnId.buf = (uint8_t *) calloc(1, 3);
//    memcpy(plmnId.buf, data, 3);
//
//    checkAndPrint(&asn_DEF_PLMN_Identity, &plmnId, (char *) "PLMN_Identity_t", __func__);
//
//}

static void createENB_ID(ENB_ID_t &enb, ENB_ID_PR enbType, unsigned char *data) {
    //printEntry("ENB_ID_t", __func__)
    ASN_STRUCT_RESET(asn_DEF_ENB_ID, &enb);
    enb.present = enbType;
    switch (enbType) {
        case ENB_ID_PR_macro_eNB_ID: { // 20 bit 3 bytes
            enb.choice.macro_eNB_ID.size = 3;
            enb.choice.macro_eNB_ID.bits_unused = 4;

            enb.present = ENB_ID_PR_macro_eNB_ID;

            enb.choice.macro_eNB_ID.buf = (uint8_t *) calloc(1, enb.choice.macro_eNB_ID.size);
            data[enb.choice.macro_eNB_ID.size - 1] = ((unsigned) (data[enb.choice.macro_eNB_ID.size - 1]
                    >> (unsigned) enb.choice.macro_eNB_ID.bits_unused)
                    << (unsigned) enb.choice.macro_eNB_ID.bits_unused);
            memcpy(enb.choice.macro_eNB_ID.buf, data, enb.choice.macro_eNB_ID.size);

            break;
        }
        case ENB_ID_PR_home_eNB_ID: { // 28 bit 4 bytes
            enb.choice.home_eNB_ID.size = 4;
            enb.choice.home_eNB_ID.bits_unused = 4;
            enb.present = ENB_ID_PR_home_eNB_ID;

            enb.choice.home_eNB_ID.buf = (uint8_t *) calloc(1, enb.choice.home_eNB_ID.size);
            data[enb.choice.home_eNB_ID.size - 1] = ((unsigned) (data[enb.choice.home_eNB_ID.size - 1]
                    >> (unsigned) enb.choice.home_eNB_ID.bits_unused)
                    << (unsigned) enb.choice.home_eNB_ID.bits_unused);
            memcpy(enb.choice.home_eNB_ID.buf, data, enb.choice.home_eNB_ID.size);
            break;
        }
        case ENB_ID_PR_short_Macro_eNB_ID: { // 18 bit - 3 bytes
            enb.choice.short_Macro_eNB_ID.size = 3;
            enb.choice.short_Macro_eNB_ID.bits_unused = 6;
            enb.present = ENB_ID_PR_short_Macro_eNB_ID;

            enb.choice.short_Macro_eNB_ID.buf = (uint8_t *) calloc(1, enb.choice.short_Macro_eNB_ID.size);
            data[enb.choice.short_Macro_eNB_ID.size - 1] = ((unsigned) (data[enb.choice.short_Macro_eNB_ID.size - 1]
                    >> (unsigned) enb.choice.short_Macro_eNB_ID.bits_unused)
                    << (unsigned) enb.choice.short_Macro_eNB_ID.bits_unused);
            memcpy(enb.choice.short_Macro_eNB_ID.buf, data, enb.choice.short_Macro_eNB_ID.size);
            break;
        }
        case ENB_ID_PR_long_Macro_eNB_ID: { // 21
            enb.choice.long_Macro_eNB_ID.size = 3;
            enb.choice.long_Macro_eNB_ID.bits_unused = 3;
            enb.present = ENB_ID_PR_long_Macro_eNB_ID;

            enb.choice.long_Macro_eNB_ID.buf = (uint8_t *) calloc(1, enb.choice.long_Macro_eNB_ID.size);
            data[enb.choice.long_Macro_eNB_ID.size - 1] = ((unsigned) (data[enb.choice.long_Macro_eNB_ID.size - 1]
                    >> (unsigned) enb.choice.long_Macro_eNB_ID.bits_unused)
                    << (unsigned) enb.choice.long_Macro_eNB_ID.bits_unused);
            memcpy(enb.choice.long_Macro_eNB_ID.buf, data, enb.choice.long_Macro_eNB_ID.size);
            break;
        }
        default:
            break;
    }

    checkAndPrint(&asn_DEF_ENB_ID, &enb, (char *) "ENB_ID_t", __func__);
}


static void buildGlobalENB_ID(GlobalENB_ID_t *gnbId,
                              const unsigned char *gnbData,
                              ENB_ID_PR enbType,
                              unsigned char *enbData) {
    auto *plmnID = createPLMN_ID(gnbData);
    memcpy(&gnbId->pLMN_Identity, plmnID, sizeof(PLMN_Identity_t));
    createENB_ID(gnbId->eNB_ID, enbType, enbData);
    checkAndPrint(&asn_DEF_GlobalENB_ID, gnbId, (char *) "GlobalENB_ID_t", __func__);
}


void buildSetupRequest(E2AP_PDU_t *pdu, int mcc, int mnc) {
    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, pdu);

    pdu->choice.initiatingMessage = (InitiatingMessage_t *)calloc(1, sizeof(InitiatingMessage_t));
    auto *initiatingMessage = pdu->choice.initiatingMessage;
    ASN_STRUCT_RESET(asn_DEF_InitiatingMessage, initiatingMessage);
    initiatingMessage->procedureCode = ProcedureCode_id_E2setup;
    initiatingMessage->criticality = Criticality_reject;
    initiatingMessage->value.present = InitiatingMessage__value_PR_E2setupRequest;

    auto *e2SetupRequestIEs = (E2setupRequestIEs_t *)calloc(1, sizeof(E2setupRequestIEs_t));
    ASN_STRUCT_RESET(asn_DEF_E2setupRequestIEs, e2SetupRequestIEs);

    e2SetupRequestIEs->value.choice.GlobalE2node_ID.present = GlobalE2node_ID_PR_gNB;
    e2SetupRequestIEs->value.choice.GlobalE2node_ID.choice.gNB = (GlobalE2node_gNB_ID_t *)calloc(1, sizeof(GlobalE2node_gNB_ID_t));
    auto *globalE2NodeGNbId = e2SetupRequestIEs->value.choice.GlobalE2node_ID.choice.gNB;
    ASN_STRUCT_RESET(asn_DEF_GlobalE2node_gNB_ID, globalE2NodeGNbId);

    createPLMN_IDByMCCandMNC(&globalE2NodeGNbId->global_gNB_ID.plmn_id, mcc, mnc);
    globalE2NodeGNbId->global_gNB_ID.gnb_id.present = GNB_ID_Choice_PR_gnb_ID;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.size = 4;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf =
            (uint8_t *) calloc(1, globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.size); //22..32 bits
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.bits_unused = 0;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf[0] = 0xB5;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf[1] = 0xC6;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf[2] = 0x77;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf[3] = 0x88;

    e2SetupRequestIEs->criticality = Criticality_reject;
    e2SetupRequestIEs->id = ProtocolIE_ID_id_GlobalE2node_ID;
    e2SetupRequestIEs->value.present = E2setupRequestIEs__value_PR_GlobalE2node_ID;


    auto *e2SetupRequest = &initiatingMessage->value.choice.E2setupRequest;

    ASN_STRUCT_RESET(asn_DEF_E2setupRequest, e2SetupRequest);
    ASN_SEQUENCE_ADD(&e2SetupRequest->protocolIEs.list, e2SetupRequestIEs);

    pdu->present = E2AP_PDU_PR_initiatingMessage;
}

void buildSetupRequesteenGNB(E2AP_PDU_t *pdu, int mcc, int mnc) {
    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, pdu);

    pdu->choice.initiatingMessage = (InitiatingMessage_t *)calloc(1, sizeof(InitiatingMessage_t));
    auto *initiatingMessage = pdu->choice.initiatingMessage;
    ASN_STRUCT_RESET(asn_DEF_InitiatingMessage, initiatingMessage);
    initiatingMessage->procedureCode = ProcedureCode_id_E2setup;
    initiatingMessage->criticality = Criticality_reject;
    initiatingMessage->value.present = InitiatingMessage__value_PR_E2setupRequest;

    auto *e2SetupRequestIEs = (E2setupRequestIEs_t *)calloc(1, sizeof(E2setupRequestIEs_t));
    ASN_STRUCT_RESET(asn_DEF_E2setupRequestIEs, e2SetupRequestIEs);

    e2SetupRequestIEs->value.choice.GlobalE2node_ID.present = GlobalE2node_ID_PR_en_gNB;
    e2SetupRequestIEs->value.choice.GlobalE2node_ID.choice.en_gNB = (GlobalE2node_en_gNB_ID_t *)calloc(1, sizeof(GlobalE2node_en_gNB_ID_t));
    auto *globalE2NodeEN_GNb = e2SetupRequestIEs->value.choice.GlobalE2node_ID.choice.en_gNB;
    ASN_STRUCT_RESET(asn_DEF_GlobalE2node_en_gNB_ID, globalE2NodeEN_GNb);

    globalE2NodeEN_GNb->global_gNB_ID.gNB_ID.present = ENGNB_ID_PR_gNB_ID;
    createPLMN_IDByMCCandMNC(&globalE2NodeEN_GNb->global_gNB_ID.pLMN_Identity, mcc, mnc);
    globalE2NodeEN_GNb->global_gNB_ID.gNB_ID.choice.gNB_ID.size = 4;
    globalE2NodeEN_GNb->global_gNB_ID.gNB_ID.choice.gNB_ID.buf =
            (uint8_t *) calloc(1, globalE2NodeEN_GNb->global_gNB_ID.gNB_ID.choice.gNB_ID.size); //22..32 bits
    globalE2NodeEN_GNb->global_gNB_ID.gNB_ID.choice.gNB_ID.buf[0] = 0xC5;
    globalE2NodeEN_GNb->global_gNB_ID.gNB_ID.choice.gNB_ID.buf[1] = 0xC6;
    globalE2NodeEN_GNb->global_gNB_ID.gNB_ID.choice.gNB_ID.buf[2] = 0xC7;
    globalE2NodeEN_GNb->global_gNB_ID.gNB_ID.choice.gNB_ID.buf[3] = 0xF8;
    globalE2NodeEN_GNb->global_gNB_ID.gNB_ID.choice.gNB_ID.bits_unused = 0;
    e2SetupRequestIEs->criticality = Criticality_reject;
    e2SetupRequestIEs->id = ProtocolIE_ID_id_GlobalE2node_ID;
    e2SetupRequestIEs->value.present = E2setupRequestIEs__value_PR_GlobalE2node_ID;


    auto *e2SetupRequest = &initiatingMessage->value.choice.E2setupRequest;

    ASN_STRUCT_RESET(asn_DEF_E2setupRequest, e2SetupRequest);
    ASN_SEQUENCE_ADD(&e2SetupRequest->protocolIEs.list, e2SetupRequestIEs);

    pdu->present = E2AP_PDU_PR_initiatingMessage;
}



void buildSetupRequestWithFunc(E2AP_PDU_t *pdu, int mcc, int mnc) {
    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, pdu);

    pdu->choice.initiatingMessage = (InitiatingMessage_t *)calloc(1, sizeof(InitiatingMessage_t));
    auto *initiatingMessage = pdu->choice.initiatingMessage;
    ASN_STRUCT_RESET(asn_DEF_InitiatingMessage, initiatingMessage);
    initiatingMessage->procedureCode = ProcedureCode_id_E2setup;
    initiatingMessage->criticality = Criticality_reject;
    initiatingMessage->value.present = InitiatingMessage__value_PR_E2setupRequest;

    auto *e2SetupRequestIEs = (E2setupRequestIEs_t *)calloc(1, sizeof(E2setupRequestIEs_t));
    ASN_STRUCT_RESET(asn_DEF_E2setupRequestIEs, e2SetupRequestIEs);

    e2SetupRequestIEs->value.choice.GlobalE2node_ID.choice.gNB = (GlobalE2node_gNB_ID_t *)calloc(1, sizeof(GlobalE2node_gNB_ID_t));
    auto *globalE2NodeGNbId = e2SetupRequestIEs->value.choice.GlobalE2node_ID.choice.gNB;
    ASN_STRUCT_RESET(asn_DEF_GlobalE2node_gNB_ID, globalE2NodeGNbId);

    createPLMN_IDByMCCandMNC(&globalE2NodeGNbId->global_gNB_ID.plmn_id, mcc, mnc);
    globalE2NodeGNbId->global_gNB_ID.gnb_id.present = GNB_ID_Choice_PR_gnb_ID;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.size = 4;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf =
            (uint8_t *) calloc(1, globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.size); //22..32 bits
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.bits_unused = 0;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf[0] = 0xB5;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf[1] = 0xC6;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf[2] = 0x77;
    globalE2NodeGNbId->global_gNB_ID.gnb_id.choice.gnb_ID.buf[3] = 0x88;

    e2SetupRequestIEs->criticality = Criticality_reject;
    e2SetupRequestIEs->id = ProtocolIE_ID_id_GlobalE2node_ID;
    e2SetupRequestIEs->value.present = E2setupRequestIEs__value_PR_GlobalE2node_ID;
    e2SetupRequestIEs->value.choice.GlobalE2node_ID.present = GlobalE2node_ID_PR_gNB;


    auto *e2SetupRequest = &initiatingMessage->value.choice.E2setupRequest;

    ASN_STRUCT_RESET(asn_DEF_E2setupRequest, e2SetupRequest);
    ASN_SEQUENCE_ADD(&e2SetupRequest->protocolIEs.list, e2SetupRequestIEs);

    auto *ranFlistIEs = (E2setupRequestIEs_t *)calloc(1, sizeof(E2setupRequestIEs_t));
    ASN_STRUCT_RESET(asn_DEF_E2setupRequestIEs, ranFlistIEs);
    ranFlistIEs->criticality = Criticality_reject;
    ranFlistIEs->id = ProtocolIE_ID_id_RANfunctionsAdded;
    ranFlistIEs->value.present = E2setupRequestIEs__value_PR_RANfunctions_List;

    auto *itemIes = (RANfunction_ItemIEs_t *)calloc(1, sizeof(RANfunction_ItemIEs_t));
    ASN_STRUCT_RESET(asn_DEF_RANfunction_ItemIEs, itemIes);


    E2SM_gNB_NRT_RANfunction_Definition_t ranFunDef;
    uint8_t funcDes[] = "asdfghjklpoiuytrewq\0";
    ranFunDef.ranFunction_Name.ranFunction_Description.buf = (uint8_t *)calloc(1, strlen((char *)funcDes));
    ranFunDef.ranFunction_Name.ranFunction_Description.size = strlen((char *)funcDes);
    memcpy(ranFunDef.ranFunction_Name.ranFunction_Description.buf, funcDes, strlen((char *)funcDes));

    uint8_t funcOID[] = "ABCDEFGHIJ1234567890\0";
    ranFunDef.ranFunction_Name.ranFunction_E2SM_OID.buf = (uint8_t *)calloc(1, strlen((char *)funcOID));
    ranFunDef.ranFunction_Name.ranFunction_E2SM_OID.size = strlen((char *)funcOID);
    memcpy(ranFunDef.ranFunction_Name.ranFunction_E2SM_OID.buf, funcOID, strlen((char *)funcOID));

    uint8_t shortName[] = "Nothing to declare\0";
    ranFunDef.ranFunction_Name.ranFunction_ShortName.buf = (uint8_t *)calloc(1, strlen((char *)shortName));
    ranFunDef.ranFunction_Name.ranFunction_ShortName.size = strlen((char *)shortName);
    memcpy(ranFunDef.ranFunction_Name.ranFunction_ShortName.buf, shortName, strlen((char *)shortName));


    RIC_InsertStyle_List_t insertStyleList;
    insertStyleList.ric_CallProcessIDFormat_Type = 28l;
    insertStyleList.ric_IndicationHeaderFormat_Type = 29;
    insertStyleList.ric_IndicationMessageFormat_Type = 30;
    insertStyleList.ric_InsertActionFormat_Type = 31l;

    uint8_t styleName[] = "What a style\0";

    insertStyleList.ric_InsertStyle_Name.buf = (uint8_t *)calloc(1, strlen((char *)styleName));
    insertStyleList.ric_InsertStyle_Name.size = strlen((char *)styleName);
    memcpy(insertStyleList.ric_InsertStyle_Name.buf, styleName, strlen((char *)styleName));


    insertStyleList.ric_InsertStyle_Type = 23;

    RANparameterDef_Item_t raNparameterDefItem;
    raNparameterDefItem.ranParameter_ID = 8;
    raNparameterDefItem.ranParameter_Type = 12;

    uint8_t ItemName[] = "What a style\0";
    raNparameterDefItem.ranParameter_Name.buf = (uint8_t *)calloc(1, strlen((char *)ItemName));
    raNparameterDefItem.ranParameter_Name.size = strlen((char *)ItemName);
    memcpy(raNparameterDefItem.ranParameter_Name.buf, ItemName, strlen((char *)ItemName));

    ASN_SEQUENCE_ADD(&insertStyleList.ric_InsertRanParameterDef_List.list, &raNparameterDefItem);

    ASN_SEQUENCE_ADD(&ranFunDef.ric_InsertStyle_List->list, &insertStyleList);
    //ranFunDef.ric_InsertStyle_List.

    uint8_t buffer[8192];
    size_t buffer_size = 8192;
    auto *ranDef = &itemIes->value.choice.RANfunction_Item.ranFunctionDefinition;

    auto er = asn_encode_to_buffer(nullptr, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_gNB_NRT_RANfunction_Definition, &ranFunDef, buffer, buffer_size);
    if (er.encoded == -1) {
        cerr << "encoding of " << asn_DEF_E2SM_gNB_NRT_RANfunction_Definition.name << " failed, " << strerror(errno) << endl;
        exit(-1);
    } else if (er.encoded > (ssize_t) buffer_size) {
        cerr << "Buffer of size " << buffer_size << " is to small for " << asn_DEF_E2SM_gNB_NRT_RANfunction_Definition.name << endl;
        exit(-1);
    } else {
        ranDef->buf = (uint8_t *)calloc(1, er.encoded);
        ranDef->size = er.encoded;
        memcpy(ranDef->buf, buffer, ranDef->size);
    }



    itemIes->id = ProtocolIE_ID_id_RANfunction_Item;
    itemIes->criticality = Criticality_reject;
    itemIes->value.present = RANfunction_ItemIEs__value_PR_RANfunction_Item;
    itemIes->value.choice.RANfunction_Item.ranFunctionID = 1;
//    auto *ranDef = &itemIes->value.choice.RANfunction_Item.ranFunctionDefinition;
//    ranDef->size = 3;
//    ranDef->buf = (uint8_t *)calloc(1, ranDef->size);
//    memcpy(ranDef->buf, buf, ranDef->size);

    ASN_SEQUENCE_ADD(&ranFlistIEs->value.choice.RANfunctions_List.list, itemIes);

    auto *itemIes1 = (RANfunction_ItemIEs_t *)calloc(1, sizeof(RANfunction_ItemIEs_t));
    ASN_STRUCT_RESET(asn_DEF_RANfunction_ItemIEs, itemIes1);
    itemIes1->id = ProtocolIE_ID_id_RANfunction_Item;
    itemIes1->criticality = Criticality_reject;
    itemIes1->value.present = RANfunction_ItemIEs__value_PR_RANfunction_Item;
    itemIes1->value.choice.RANfunction_Item.ranFunctionID = 7;
    ranDef = &itemIes1->value.choice.RANfunction_Item.ranFunctionDefinition;

    ranDef->buf = (uint8_t *)calloc(1, er.encoded);
    ranDef->size = er.encoded;
    memcpy(ranDef->buf, buffer, ranDef->size);

    ASN_SEQUENCE_ADD(&ranFlistIEs->value.choice.RANfunctions_List.list, itemIes1);


    ASN_SEQUENCE_ADD(&e2SetupRequest->protocolIEs.list, ranFlistIEs);

    pdu->present = E2AP_PDU_PR_initiatingMessage;
}





void buildSetupSuccsessfulResponse(E2AP_PDU_t *pdu, int mcc, int mnc, uint8_t *data) {
    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, pdu);

    pdu->choice.successfulOutcome = (SuccessfulOutcome_t *)calloc(1, sizeof(SuccessfulOutcome_t));
    SuccessfulOutcome_t *successfulOutcome = pdu->choice.successfulOutcome;
    ASN_STRUCT_RESET(asn_DEF_E2setupResponse, &successfulOutcome->value.choice.E2setupResponse);
    successfulOutcome->procedureCode = ProcedureCode_id_E2setup;
    successfulOutcome->criticality = Criticality_reject;
    successfulOutcome->value.present = SuccessfulOutcome__value_PR_E2setupResponse;


    auto *globalRicidIE = (E2setupResponseIEs_t *)calloc(1, sizeof(E2setupResponseIEs_t));
    ASN_STRUCT_RESET(asn_DEF_E2setupResponseIEs, globalRicidIE);

    globalRicidIE->criticality = Criticality_reject;
    globalRicidIE->id = ProtocolIE_ID_id_GlobalRIC_ID;
    globalRicidIE->value.present = E2setupResponseIEs__value_PR_GlobalRIC_ID;
    createPLMN_IDByMCCandMNC(&globalRicidIE->value.choice.GlobalRIC_ID.pLMN_Identity, mcc, mnc);

    globalRicidIE->value.choice.GlobalRIC_ID.ric_ID = {nullptr, 3, 4};
    globalRicidIE->value.choice.GlobalRIC_ID.ric_ID.buf = (uint8_t *)calloc(1, globalRicidIE->value.choice.GlobalRIC_ID.ric_ID.size);
    memcpy(globalRicidIE->value.choice.GlobalRIC_ID.ric_ID.buf, data, globalRicidIE->value.choice.GlobalRIC_ID.ric_ID.size);
    globalRicidIE->value.choice.GlobalRIC_ID.ric_ID.buf[2] &= (unsigned)0xF0;


    ASN_STRUCT_RESET(asn_DEF_E2setupResponse, &successfulOutcome->value.choice.E2setupResponse);
    ASN_SEQUENCE_ADD(&successfulOutcome->value.choice.E2setupResponse.protocolIEs.list, globalRicidIE);

    auto *ranFunctionAdd = (E2setupResponseIEs_t *)calloc(1, sizeof(E2setupResponseIEs_t));
    ASN_STRUCT_RESET(asn_DEF_E2setupResponseIEs, ranFunctionAdd);
    ranFunctionAdd->criticality = Criticality_reject;
    ranFunctionAdd->id = ProtocolIE_ID_id_RANfunctionsAccepted;
    ranFunctionAdd->value.present = E2setupResponseIEs__value_PR_RANfunctionsID_List;

    auto *ranFuncIdItemIEs = (RANfunctionID_ItemIEs_t *)calloc(1, sizeof(RANfunctionID_ItemIEs_t));

    ranFuncIdItemIEs->criticality = Criticality_ignore;
    ranFuncIdItemIEs->id = ProtocolIE_ID_id_RANfunctionID_Item;
    ranFuncIdItemIEs->value.present = RANfunctionID_ItemIEs__value_PR_RANfunctionID_Item;
    ranFuncIdItemIEs->value.choice.RANfunctionID_Item.ranFunctionID = 10;
    ranFuncIdItemIEs->value.choice.RANfunctionID_Item.ranFunctionRevision = 1;

    ASN_SEQUENCE_ADD(&ranFunctionAdd->value.choice.RANfunctionsID_List.list, ranFuncIdItemIEs);
    ASN_SEQUENCE_ADD(&successfulOutcome->value.choice.E2setupResponse.protocolIEs.list, ranFunctionAdd);





    pdu->present = E2AP_PDU_PR_successfulOutcome;
}


void buildSetupUnSuccsessfulResponse(E2AP_PDU_t *pdu) {
    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, pdu);

    pdu->choice.unsuccessfulOutcome = (UnsuccessfulOutcome_t *)calloc(1, sizeof(UnsuccessfulOutcome_t));
    UnsuccessfulOutcome_t *uns = pdu->choice.unsuccessfulOutcome;
    uns->procedureCode = ProcedureCode_id_E2setup;
    uns->criticality = Criticality_reject;
    uns->value.present = UnsuccessfulOutcome__value_PR_E2setupFailure;

    ASN_STRUCT_RESET(asn_DEF_E2setupFailure, &uns->value.choice.E2setupFailure);


    {
        auto *e2SetupFIE = (E2setupFailureIEs_t *) calloc(1, sizeof(E2setupFailureIEs_t));
        ASN_STRUCT_RESET(asn_DEF_E2setupFailureIEs, e2SetupFIE);

        e2SetupFIE->criticality = Criticality_ignore;
        e2SetupFIE->id = ProtocolIE_ID_id_Cause;
        e2SetupFIE->value.present = E2setupFailureIEs__value_PR_Cause;
        e2SetupFIE->value.choice.Cause.present = Cause_PR_transport;
        e2SetupFIE->value.choice.Cause.choice.transport = CauseTransport_transport_resource_unavailable;


        ASN_SEQUENCE_ADD(&uns->value.choice.E2setupFailure.protocolIEs.list, e2SetupFIE);
    }

    {
        auto *e2SetupFIE = (E2setupFailureIEs_t *) calloc(1, sizeof(E2setupFailureIEs_t));
        ASN_STRUCT_RESET(asn_DEF_E2setupFailureIEs, e2SetupFIE);

        e2SetupFIE->criticality = Criticality_ignore;
        e2SetupFIE->id = ProtocolIE_ID_id_TimeToWait;
        e2SetupFIE->value.present = E2setupFailureIEs__value_PR_TimeToWait;
        e2SetupFIE->value.choice.TimeToWait = TimeToWait_v60s;

        ASN_SEQUENCE_ADD(&uns->value.choice.E2setupFailure.protocolIEs.list, e2SetupFIE);
    }
    {
        auto *e2SetupFIE = (E2setupFailureIEs_t *) calloc(1, sizeof(E2setupFailureIEs_t));
        ASN_STRUCT_RESET(asn_DEF_E2setupFailureIEs, e2SetupFIE);

        e2SetupFIE->criticality = Criticality_ignore;
        e2SetupFIE->id = ProtocolIE_ID_id_CriticalityDiagnostics;
        e2SetupFIE->value.present = E2setupFailureIEs__value_PR_CriticalityDiagnostics;
        e2SetupFIE->value.choice.CriticalityDiagnostics.procedureCode = (ProcedureCode_t *)calloc(1,sizeof(ProcedureCode_t));
        *e2SetupFIE->value.choice.CriticalityDiagnostics.procedureCode = ProcedureCode_id_E2setup;
        e2SetupFIE->value.choice.CriticalityDiagnostics.triggeringMessage = (TriggeringMessage_t *)calloc(1,sizeof(TriggeringMessage_t));
        *e2SetupFIE->value.choice.CriticalityDiagnostics.triggeringMessage = TriggeringMessage_initiating_message;
        e2SetupFIE->value.choice.CriticalityDiagnostics.procedureCriticality = (Criticality_t *)calloc(1, sizeof(Criticality_t));
        *e2SetupFIE->value.choice.CriticalityDiagnostics.procedureCriticality = Criticality_reject;
        ASN_SEQUENCE_ADD(&uns->value.choice.E2setupFailure.protocolIEs.list, e2SetupFIE);
    }

    pdu->present = E2AP_PDU_PR_unsuccessfulOutcome;
}

#endif //E2_E2BUILDER_H
