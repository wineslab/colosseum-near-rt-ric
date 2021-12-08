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

//
// Created by adi ENZEL on 6/19/19.
//

#include "e2sm.h"

#define printEntry(type, function) \
    if (mdclog_level_get() >= MDCLOG_DEBUG) { \
        mdclog_write(MDCLOG_DEBUG, "start Test %s , %s", type, function); \
    }


static void checkAndPrint(asn_TYPE_descriptor_t *typeDescriptor, void *data, char *dataType, const char *function) {
    char errbuf[128]; /* Buffer for error message */
    size_t errlen = sizeof(errbuf); /* Size of the buffer */
    if (asn_check_constraints(typeDescriptor, data, errbuf, &errlen) != 0) {
        mdclog_write(MDCLOG_ERR, "%s Constraint validation failed: %s", dataType, errbuf);
    } else if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "%s successes function %s", dataType, function);
    }
}

static size_t encodebuff(int codingType,
                         asn_TYPE_descriptor_t *typeDescriptor,
                         void *objectData,
                         uint8_t *buffer,
                         size_t buffer_size) {
    asn_enc_rval_t er;
    struct timespec start = {0,0};
    struct timespec end   = {0,0};
    clock_gettime(CLOCK_MONOTONIC, &start);
    er = asn_encode_to_buffer(0, codingType, typeDescriptor, objectData, buffer, buffer_size);
    clock_gettime(CLOCK_MONOTONIC, &end);
    if (er.encoded == -1) {
        mdclog_write(MDCLOG_ERR, "encoding of %s failed, %s", asn_DEF_E2SM_gNB_X2_eventTriggerDefinition.name, strerror(errno));
    } else if (er.encoded > (ssize_t)buffer_size) {
        mdclog_write(MDCLOG_ERR, "Buffer of size %d is to small for %s", (int) buffer_size,
                     typeDescriptor->name);
    } else if (mdclog_level_get() >= MDCLOG_DEBUG) {
        if (codingType == ATS_BASIC_XER) {
            mdclog_write(MDCLOG_DEBUG, "Buffer of size %d, data = %s", (int) er.encoded, buffer);
        }
        else {
            if (mdclog_level_get() >= MDCLOG_INFO) {
                char *printBuffer;
                size_t size;
                FILE *stream = open_memstream(&printBuffer, &size);
                asn_fprint(stream, typeDescriptor, objectData);
                mdclog_write(MDCLOG_DEBUG, "Encoding E2SM PDU past : %s", printBuffer);
            }


            mdclog_write(MDCLOG_DEBUG, "Buffer of size %d", (int) er.encoded);
        }
    }
    mdclog_write(MDCLOG_INFO, "Encoding type %d, time is %ld seconds, %ld nanoseconds", codingType, end.tv_sec - start.tv_sec, end.tv_nsec - start.tv_nsec);
    //mdclog_write(MDCLOG_INFO, "Encoding time is %3.9f seconds", ((double)end.tv_sec + 1.0e-9 * end.tv_nsec) - ((double)start.tv_sec + 1.0e-9 * start.tv_nsec));
    return er.encoded;
}

PLMN_Identity_t *createPLMN_ID(const unsigned char *data) {
    printEntry("PLMN_Identity_t", __func__)
    PLMN_Identity_t *plmnId = calloc(1, sizeof(PLMN_Identity_t));
    ASN_STRUCT_RESET(asn_DEF_PLMN_Identity, plmnId);
    plmnId->size = 3;
    plmnId->buf = calloc(1, 3);
    memcpy(plmnId->buf, data, 3);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_PLMN_Identity, plmnId, "PLMN_Identity_t", __func__);
    }

    return plmnId;
}

ENB_ID_t *createENB_ID(ENB_ID_PR enbType, unsigned char *data) {
    printEntry("ENB_ID_t", __func__)
    ENB_ID_t *enb = calloc(1, sizeof(ENB_ID_t));
    ASN_STRUCT_RESET(asn_DEF_ENB_ID, enb);

    enb->present = enbType;

    switch (enbType) {
        case ENB_ID_PR_macro_eNB_ID: { // 20 bit 3 bytes
            enb->choice.macro_eNB_ID.size = 3;
            enb->choice.macro_eNB_ID.bits_unused = 4;

            enb->present = ENB_ID_PR_macro_eNB_ID;

            enb->choice.macro_eNB_ID.buf = calloc(1, enb->choice.macro_eNB_ID.size);
            data[enb->choice.macro_eNB_ID.size - 1] = ((unsigned)(data[enb->choice.macro_eNB_ID.size - 1] >>
                    (unsigned)enb->choice.macro_eNB_ID.bits_unused) << (unsigned)enb->choice.macro_eNB_ID.bits_unused);
            memcpy(enb->choice.macro_eNB_ID.buf, data, enb->choice.macro_eNB_ID.size);

            break;
        }
        case ENB_ID_PR_home_eNB_ID: { // 28 bit 4 bytes
            enb->choice.home_eNB_ID.size = 4;
            enb->choice.home_eNB_ID.bits_unused = 4;
            enb->present = ENB_ID_PR_home_eNB_ID;

            enb->choice.home_eNB_ID.buf = calloc(1, enb->choice.home_eNB_ID.size);
            data[enb->choice.home_eNB_ID.size - 1] = ((unsigned)(data[enb->choice.home_eNB_ID.size - 1] >>
                    (unsigned)enb->choice.home_eNB_ID.bits_unused) << (unsigned)enb->choice.home_eNB_ID.bits_unused);
            memcpy(enb->choice.home_eNB_ID.buf, data, enb->choice.home_eNB_ID.size);
            break;
        }
        case ENB_ID_PR_short_Macro_eNB_ID: { // 18 bit - 3 bytes
            enb->choice.short_Macro_eNB_ID.size = 3;
            enb->choice.short_Macro_eNB_ID.bits_unused = 6;
            enb->present = ENB_ID_PR_short_Macro_eNB_ID;

            enb->choice.short_Macro_eNB_ID.buf = calloc(1, enb->choice.short_Macro_eNB_ID.size);
            data[enb->choice.short_Macro_eNB_ID.size - 1] = ((unsigned)(data[enb->choice.short_Macro_eNB_ID.size - 1] >>
                    (unsigned)enb->choice.short_Macro_eNB_ID.bits_unused) << (unsigned)enb->choice.short_Macro_eNB_ID.bits_unused);
            memcpy(enb->choice.short_Macro_eNB_ID.buf, data, enb->choice.short_Macro_eNB_ID.size);
            break;
        }
        case ENB_ID_PR_long_Macro_eNB_ID: { // 21
            enb->choice.long_Macro_eNB_ID.size = 3;
            enb->choice.long_Macro_eNB_ID.bits_unused = 3;
            enb->present = ENB_ID_PR_long_Macro_eNB_ID;

            enb->choice.long_Macro_eNB_ID.buf = calloc(1, enb->choice.long_Macro_eNB_ID.size);
            data[enb->choice.long_Macro_eNB_ID.size - 1] = ((unsigned)(data[enb->choice.long_Macro_eNB_ID.size - 1] >>
                    (unsigned)enb->choice.long_Macro_eNB_ID.bits_unused) << (unsigned)enb->choice.long_Macro_eNB_ID.bits_unused);
            memcpy(enb->choice.long_Macro_eNB_ID.buf, data, enb->choice.long_Macro_eNB_ID.size);
            break;
        }
        default:
            free(enb);
            return NULL;

    }

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_ENB_ID, enb, "ENB_ID_t", __func__);
    }
    return enb;
}

GNB_ID_t *createGnb_id(const unsigned char *data, int numOfBits) {
    printEntry("GNB_ID_t", __func__)
    if (numOfBits < 22 || numOfBits > 32) {
        mdclog_write(MDCLOG_ERR, "GNB_ID_t number of bits = %d, needs to be 22 .. 32", numOfBits);
        return NULL;
    }
    GNB_ID_t *gnb = calloc(1, sizeof(GNB_ID_t));
    ASN_STRUCT_RESET(asn_DEF_GNB_ID, gnb);

    gnb->present = GNB_ID_PR_gNB_ID;
    gnb->choice.gNB_ID.size = numOfBits % 8 == 0 ? (unsigned int)(numOfBits / 8) : (unsigned int)(numOfBits / 8 + 1);
    gnb->choice.gNB_ID.bits_unused = (int)gnb->choice.gNB_ID.size * 8 - numOfBits;
    gnb->choice.gNB_ID.buf = calloc(1, gnb->choice.gNB_ID.size);
    memcpy(gnb->choice.gNB_ID.buf, data, gnb->choice.gNB_ID.size);
    gnb->choice.gNB_ID.buf[gnb->choice.gNB_ID.size - 1] =
            ((unsigned)(gnb->choice.gNB_ID.buf[gnb->choice.gNB_ID.size - 1] >> (unsigned)gnb->choice.gNB_ID.bits_unused)
                  << (unsigned)gnb->choice.gNB_ID.bits_unused);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_GNB_ID, gnb, "GNB_ID_t", __func__);
    }

    return gnb;

}

GlobalENB_ID_t *createGlobalENB_ID(PLMN_Identity_t *plmnIdentity, ENB_ID_t *enbId) {
    printEntry("GlobalENB_ID_t", __func__)
    GlobalENB_ID_t *genbId = calloc(1, sizeof(GlobalENB_ID_t));
    ASN_STRUCT_RESET(asn_DEF_GlobalENB_ID, genbId);
    memcpy(&genbId->pLMN_Identity, plmnIdentity, sizeof(PLMN_Identity_t));
    memcpy(&genbId->eNB_ID, enbId, sizeof(ENB_ID_t));

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_GlobalENB_ID, genbId, "GlobalENB_ID_t", __func__);
    }
    return genbId;
}

GlobalGNB_ID_t *createGlobalGNB_ID(PLMN_Identity_t *plmnIdentity, GNB_ID_t *gnb) {
    printEntry("GlobalGNB_ID_t", __func__)
    GlobalGNB_ID_t *ggnbId = calloc(1, sizeof(GlobalGNB_ID_t));
    ASN_STRUCT_RESET(asn_DEF_GlobalGNB_ID, ggnbId);

    memcpy(&ggnbId->pLMN_Identity, plmnIdentity, sizeof(PLMN_Identity_t));
    memcpy(&ggnbId->gNB_ID, gnb, sizeof(GNB_ID_t));

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_GlobalGNB_ID, ggnbId, "GlobalGNB_ID_t", __func__);
    }

    return ggnbId;
}


Interface_ID_t *createInterfaceIDForGnb(GlobalGNB_ID_t *gnb) {
    printEntry("Interface_ID_t", __func__)
    Interface_ID_t *interfaceId = calloc(1, sizeof(Interface_ID_t));
    ASN_STRUCT_RESET(asn_DEF_Interface_ID, interfaceId);

    interfaceId->present = Interface_ID_PR_global_gNB_ID;
    //memcpy(&interfaceId->choice.global_gNB_ID, gnb, sizeof(GlobalGNB_ID_t));
    interfaceId->choice.global_gNB_ID = gnb;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_Interface_ID, interfaceId, "Interface_ID_t", __func__);
    }

    return interfaceId;
}

Interface_ID_t *createInterfaceIDForEnb(GlobalENB_ID_t *enb) {
    printEntry("Interface_ID_t", __func__)
    Interface_ID_t *interfaceId = calloc(1, sizeof(Interface_ID_t));
    ASN_STRUCT_RESET(asn_DEF_Interface_ID, interfaceId);

    interfaceId->present = Interface_ID_PR_global_eNB_ID;
    //memcpy(&interfaceId->choice.global_eNB_ID, enb, sizeof(GlobalENB_ID_t));
    interfaceId->choice.global_eNB_ID = enb;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_Interface_ID, interfaceId, "Interface_ID_t", __func__);
    }

    return interfaceId;
}



InterfaceMessageType_t *createInterfaceMessageInitiating(long procedureCode) {
    printEntry("InterfaceMessageType_t", __func__)
    InterfaceMessageType_t *intMsgT = calloc(1, sizeof(InterfaceMessageType_t));
    ASN_STRUCT_RESET(asn_DEF_InterfaceMessageType, intMsgT);

    intMsgT->procedureCode = procedureCode;
    intMsgT->typeOfMessage = TypeOfMessage_initiating_message;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_InterfaceMessageType, intMsgT, "InterfaceMessageType_t", __func__);
    }

    return intMsgT;
}

InterfaceMessageType_t *createInterfaceMessageSuccsesful(long procedureCode) {
    printEntry("InterfaceMessageType_t", __func__)
    InterfaceMessageType_t *intMsgT = calloc(1, sizeof(InterfaceMessageType_t));
    ASN_STRUCT_RESET(asn_DEF_InterfaceMessageType, intMsgT);

    intMsgT->procedureCode = procedureCode;
    intMsgT->typeOfMessage = TypeOfMessage_successful_outcome;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_InterfaceMessageType, intMsgT, "InterfaceMessageType_t", __func__);
    }

    return intMsgT;
}

InterfaceMessageType_t *createInterfaceMessageUnsuccessful(long procedureCode) {
    printEntry("InterfaceMessageType_t", __func__)
    InterfaceMessageType_t *intMsgT = calloc(1, sizeof(InterfaceMessageType_t));
    ASN_STRUCT_RESET(asn_DEF_InterfaceMessageType, intMsgT);

    intMsgT->procedureCode = procedureCode;
    intMsgT->typeOfMessage = TypeOfMessage_unsuccessful_outcome;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_InterfaceMessageType, intMsgT, "InterfaceMessageType_t", __func__);
    }

    return intMsgT;
}

InterfaceProtocolIE_Value_t *createInterfaceProtocolValueInt(long number) {
    printEntry("InterfaceProtocolIE_Value_t", __func__)
    InterfaceProtocolIE_Value_t *value = calloc(1, sizeof(InterfaceProtocolIE_Value_t));
    ASN_STRUCT_RESET(asn_DEF_InterfaceProtocolIE_Value, value);

    value->present = InterfaceProtocolIE_Value_PR_valueInt;
    value->choice.valueInt = number;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_InterfaceProtocolIE_Value, value, "InterfaceProtocolIE_Value_t", __func__);
    }

    return value;
}

InterfaceProtocolIE_Value_t *createInterfaceProtocolValueEnum(long number) {
    printEntry("InterfaceProtocolIE_Value_t", __func__)
    InterfaceProtocolIE_Value_t *value = calloc(1, sizeof(InterfaceProtocolIE_Value_t));
    ASN_STRUCT_RESET(asn_DEF_InterfaceProtocolIE_Value, value);

    value->present = InterfaceProtocolIE_Value_PR_valueEnum;
    value->choice.valueEnum = number;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_InterfaceProtocolIE_Value, value, "InterfaceProtocolIE_Value_t", __func__);
    }

    return value;
}

InterfaceProtocolIE_Value_t *createInterfaceProtocolValueBool(int val) {
    printEntry("InterfaceProtocolIE_Value_t", __func__)
    InterfaceProtocolIE_Value_t *value = calloc(1, sizeof(InterfaceProtocolIE_Value_t));
    ASN_STRUCT_RESET(asn_DEF_InterfaceProtocolIE_Value, value);

    value->present = InterfaceProtocolIE_Value_PR_valueBool;
    value->choice.valueBool = val == 0 ? 0 : 1;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_InterfaceProtocolIE_Value, value, "InterfaceProtocolIE_Value_t", __func__);
    }

    return value;
}


InterfaceProtocolIE_Value_t *createInterfaceProtocolValueBitString(unsigned char *buf, int numOfBits) {
    printEntry("InterfaceProtocolIE_Value_t", __func__)
    size_t size = numOfBits % 8 == 0 ? (unsigned int)(numOfBits / 8) : (unsigned int)(numOfBits / 8 + 1);
    if (strlen((const char *)buf) < size) {
        mdclog_write(MDCLOG_ERR, "size of buffer is small : %d needs to be %d in %s", (int)strlen((const char *)buf), (int)size, __func__);
    }
    InterfaceProtocolIE_Value_t *value = calloc(1, sizeof(InterfaceProtocolIE_Value_t));
    ASN_STRUCT_RESET(asn_DEF_InterfaceProtocolIE_Value, value);

    value->present = InterfaceProtocolIE_Value_PR_valueBitS;
    value->choice.valueBitS.size = numOfBits % 8 == 0 ? (unsigned int)(numOfBits / 8) : (unsigned int)(numOfBits / 8 + 1);
    value->choice.valueBitS.buf = calloc(1, value->choice.valueBitS.size);
    int bits_unused = (int)value->choice.valueBitS.size * 8 - numOfBits;
    value->choice.valueBitS.bits_unused = bits_unused;

    memcpy(value->choice.valueBitS.buf, buf, value->choice.valueBitS.size);
    value->choice.valueBitS.buf[size -1] = ((unsigned)(buf[size - 1]) >> (unsigned)bits_unused << (unsigned)bits_unused);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_InterfaceProtocolIE_Value, value, "InterfaceProtocolIE_Value_t", __func__);
    }

    return value;
}


InterfaceProtocolIE_Value_t *createInterfaceProtocolValueOCTETS(uint8_t *buf) {
    printEntry("InterfaceProtocolIE_Value_t", __func__)
    size_t size = strlen((const char *)buf);
    InterfaceProtocolIE_Value_t *value = calloc(1, sizeof(InterfaceProtocolIE_Value_t));
    ASN_STRUCT_RESET(asn_DEF_InterfaceProtocolIE_Value, value);

    value->present = InterfaceProtocolIE_Value_PR_valueOctS;
    value->choice.valueOctS.size = size;
    value->choice.valueOctS.buf = calloc(1, value->choice.valueOctS.size);
    memcpy(value->choice.valueOctS.buf, buf, value->choice.valueOctS.size);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_InterfaceProtocolIE_Value, value, "InterfaceProtocolIE_Value_t", __func__);
    }

    return value;
}


InterfaceProtocolIE_Item_t *createInterfaceProtocolIE_Item(long id, long test, InterfaceProtocolIE_Value_t *value) {
    printEntry("InterfaceProtocolIE_Item_t", __func__)
    if (test < InterfaceProtocolIE_Test_equal || test > InterfaceProtocolIE_Test_present) {
        mdclog_write(MDCLOG_ERR, "InterfaceProtocolIE_Item_t test value is %ld,  out of scope %d .. %d ",
                test, InterfaceProtocolIE_Test_equal, InterfaceProtocolIE_Test_present);
        return NULL;
    }
    InterfaceProtocolIE_Item_t *intProtIt = calloc(1, sizeof(InterfaceProtocolIE_Item_t));
    ASN_STRUCT_RESET(asn_DEF_InterfaceProtocolIE_Item, intProtIt);


    intProtIt->interfaceProtocolIE_ID = id;

    intProtIt->interfaceProtocolIE_Test = test;

    memcpy(&intProtIt->interfaceProtocolIE_Value, value, sizeof(InterfaceProtocolIE_Value_t));

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_InterfaceProtocolIE_Item, intProtIt, "InterfaceProtocolIE_Item_t", __func__);
    }

    return intProtIt;

}



ActionParameter_Value_t *createActionParameterValue_Int(long number) {
    printEntry("ActionParameter_Value_t", __func__)
    ActionParameter_Value_t *value = calloc(1, sizeof(ActionParameter_Value_t));
    ASN_STRUCT_RESET(asn_DEF_ActionParameter_Value, value);

    value->present = ActionParameter_Value_PR_valueInt;
    value->choice.valueInt = number;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_ActionParameter_Value, value, "ActionParameter_Value_t", __func__);
    }

    return value;
}

ActionParameter_Value_t *createActionParameterValue_Enum(long number) {
    printEntry("ActionParameter_Value_t", __func__)
    ActionParameter_Value_t *value = calloc(1, sizeof(ActionParameter_Value_t));
    ASN_STRUCT_RESET(asn_DEF_ActionParameter_Value, value);

    value->present = ActionParameter_Value_PR_valueEnum;
    value->choice.valueEnum = number;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_ActionParameter_Value, value, "ActionParameter_Value_t", __func__);
    }

    return value;
}

ActionParameter_Value_t *createActionParameterValue_Bool(int val) {
    printEntry("ActionParameter_Value_t", __func__)
    ActionParameter_Value_t *value = calloc(1, sizeof(ActionParameter_Value_t));
    ASN_STRUCT_RESET(asn_DEF_ActionParameter_Value, value);

    value->present = ActionParameter_Value_PR_valueBool;
    value->choice.valueBool = val == 0 ? 0 : 1;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_ActionParameter_Value, value, "ActionParameter_Value_t", __func__);
    }

    return value;
}


ActionParameter_Value_t *createActionParameterValue_Bit_String(unsigned char *buf, int numOfBits) {
    printEntry("ActionParameter_Value_t", __func__)
    size_t size = numOfBits % 8 == 0 ? (unsigned int)(numOfBits / 8) : (unsigned int)(numOfBits / 8 + 1);
    if (strlen((const char *)buf) < size) {
        mdclog_write(MDCLOG_ERR, "size of buffer is small : %d needs to be %d in %s", (int)strlen((const char *)buf), (int)size, __func__);
    }
    int bits_unused = (int)size * 8 - numOfBits;

    ActionParameter_Value_t *value = calloc(1, sizeof(ActionParameter_Value_t));
    ASN_STRUCT_RESET(asn_DEF_ActionParameter_Value, value);

    value->present = ActionParameter_Value_PR_valueBitS;
    value->choice.valueBitS.size = size;
    value->choice.valueBitS.buf = calloc(1, size);
    value->choice.valueBitS.bits_unused = bits_unused;

    memcpy(value->choice.valueBitS.buf, buf, size);
    value->choice.valueBitS.buf[size -1] = ((unsigned)(buf[size - 1]) >> (unsigned)bits_unused << (unsigned)bits_unused);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_ActionParameter_Value, value, "ActionParameter_Value_t", __func__);
    }

    return value;
}


ActionParameter_Value_t *createActionParameterValue_OCTETS(uint8_t *buf) {
    printEntry("ActionParameter_Value_t", __func__)
    size_t size = strlen((const char *)buf);
    ActionParameter_Value_t *value = calloc(1, sizeof(ActionParameter_Value_t));
    ASN_STRUCT_RESET(asn_DEF_ActionParameter_Value, value);

    value->present = ActionParameter_Value_PR_valueOctS;
    value->choice.valueOctS.size = size;
    value->choice.valueOctS.buf = calloc(1, value->choice.valueOctS.size);
    memcpy(value->choice.valueOctS.buf, buf, value->choice.valueOctS.size);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_ActionParameter_Value, value, "ActionParameter_Value_t", __func__);
    }

    return value;
}

/**
 *
 * @param buf buffer that must be null terminated
 * @return ActionParameter_Value_t *
 */
ActionParameter_Value_t *createActionParameterValue_PRINTS(char *buf) {
    printEntry("ActionParameter_Value_t", __func__)
    size_t size = strlen((const char *)buf);
    ActionParameter_Value_t *value = calloc(1, sizeof(ActionParameter_Value_t));
    ASN_STRUCT_RESET(asn_DEF_ActionParameter_Value, value);

    value->present = ActionParameter_Value_PR_valuePrtS;
    value->choice.valuePrtS.size = size;
    value->choice.valuePrtS.buf = calloc(1, value->choice.valuePrtS.size);
    memcpy(value->choice.valuePrtS.buf, buf, value->choice.valuePrtS.size);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_ActionParameter_Value, value, "ActionParameter_Value_t", __func__);
    }

    return value;
}

ActionParameter_Item_t *creatActionParameter_Item(long id, ActionParameter_Value_t *val) {
    printEntry("ActionParameter_Item_t", __func__)
    if (id < 0 || id > 255) {
        mdclog_write(MDCLOG_ERR, "ActionParameter_Item_t id = %ld, values are 0 .. 255", id);
        return NULL;
    }
    ActionParameter_Item_t *actionParameterItem = calloc(1, sizeof(ActionParameter_Item_t));
    ASN_STRUCT_RESET(asn_DEF_ActionParameter_Item, actionParameterItem);

    actionParameterItem->actionParameter_ID = id;
    memcpy(&actionParameterItem->actionParameter_Value, val, sizeof(ActionParameter_Value_t));

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_ActionParameter_Item, actionParameterItem, "ActionParameter_Item_t", __func__);
    }

    return actionParameterItem;
}

/**
 *
 * @param interfaceId
 * @param direction
 * @param messageType
 * @param interfaceProtocolItemList
 * @param listSize
 * @param buffer
 * @param buffer_size
 * @return
 */
size_t createEventTrigger(Interface_ID_t *interfaceId, long direction,
                          InterfaceMessageType_t *messageType,
                          InterfaceProtocolIE_Item_t interfaceProtocolItemList[],
                          int listSize,
                          uint8_t *buffer,
                          size_t buffer_size) {
    printEntry("E2SM_gNB_X2_eventTriggerDefinition_t", __func__)
    if (direction < InterfaceDirection_incoming || direction > InterfaceDirection_outgoing) {
        mdclog_write(MDCLOG_ERR, "E2SM_gNB_X2_eventTriggerDefinition_t direction = %ld, values are %d .. %d",
                     direction, InterfaceDirection_incoming, InterfaceDirection_outgoing);
        return -1;
    }

    E2SM_gNB_X2_eventTriggerDefinition_t *eventTrigger = calloc(1, sizeof(E2SM_gNB_X2_eventTriggerDefinition_t));
    ASN_STRUCT_RESET(asn_DEF_E2SM_gNB_X2_eventTriggerDefinition, eventTrigger);

    memcpy(&eventTrigger->interface_ID , interfaceId, sizeof(Interface_ID_t));

    eventTrigger->interfaceDirection = direction;
    memcpy(&eventTrigger->interfaceMessageType, messageType, sizeof(InterfaceMessageType_t));

    for (int i = 0; i < listSize; i++) {
        ASN_SEQUENCE_ADD(eventTrigger->interfaceProtocolIE_List, &interfaceProtocolItemList[i]);
    }

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_E2SM_gNB_X2_eventTriggerDefinition, eventTrigger, "E2SM_gNB_X2_eventTriggerDefinition_t", __func__);
    }

    size_t  len = encodebuff(ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_gNB_X2_eventTriggerDefinition,
                             eventTrigger,
                             buffer,
                             buffer_size);


    if (mdclog_level_get() >= MDCLOG_INFO) {
        uint8_t buf1[4096];
        //asn_enc_rval_t er1;
        encodebuff(ATS_BASIC_XER, &asn_DEF_E2SM_gNB_X2_eventTriggerDefinition,
                                 eventTrigger,
                                 buf1,
                                 4096);

    }

    return len;
}


size_t createActionDefinition(long styleId, ActionParameter_Item_t actionParamList[], int listSize,
                                                       uint8_t *buffer,
                                                       size_t buffer_size) {
    printEntry("E2SM_gNB_X2_actionDefinition_t", __func__)
    E2SM_gNB_X2_actionDefinition_t *actDef = calloc(1, sizeof(E2SM_gNB_X2_actionDefinition_t));
    ASN_STRUCT_RESET(asn_DEF_E2SM_gNB_X2_actionDefinition, actDef);

    actDef->style_ID = styleId;
    for (int i = 0; i < listSize; i++) {
        ASN_SEQUENCE_ADD(actDef->actionParameter_List, &actionParamList[i]);
    }

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_E2SM_gNB_X2_actionDefinition, actDef, "E2SM_gNB_X2_actionDefinition_t", __func__);
    }

    size_t  len = encodebuff(ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_gNB_X2_actionDefinition,
                             actDef,
                             buffer,
                             buffer_size);


    if (mdclog_level_get() >= MDCLOG_INFO) {
        uint8_t buf1[4096];
        //asn_enc_rval_t er1;
        encodebuff(ATS_BASIC_XER, &asn_DEF_E2SM_gNB_X2_actionDefinition,
                   actDef,
                   buf1,
                   4096);

    }

    return len;
}

size_t createE2SM_gNB_X2_indicationHeader(long direction,
                                          Interface_ID_t *interfaceId,
                                          uint8_t *timestamp, //can put NULL if size == 0
                                          int size,
                                          uint8_t *buffer,
                                          size_t buffer_size) {
    printEntry("E2SM_gNB_X2_indicationHeader_t", __func__)
    if (direction < InterfaceDirection_incoming || direction > InterfaceDirection_outgoing) {
        mdclog_write(MDCLOG_ERR, "E2SM_gNB_X2_indicationHeader_t direction = %ld, values are %d .. %d",
                     direction, InterfaceDirection_incoming, InterfaceDirection_outgoing);
        return -1;
    }

    E2SM_gNB_X2_indicationHeader_t *indiHead = calloc(1, sizeof(E2SM_gNB_X2_indicationHeader_t));
    ASN_STRUCT_RESET(asn_DEF_E2SM_gNB_X2_indicationHeader, indiHead);

    indiHead->interfaceDirection = direction;
    memcpy(&indiHead->interface_ID, interfaceId, sizeof(Interface_ID_t));
    if (size > 0) {
        indiHead->timestamp->size = size;
        indiHead->timestamp->buf = calloc(1, sizeof(uint8_t) * size);
        memcpy(indiHead->timestamp->buf, timestamp, size);
    }

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_E2SM_gNB_X2_indicationHeader, indiHead, "E2SM_gNB_X2_indicationHeader_t", __func__);
    }

    size_t  len = encodebuff(ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_gNB_X2_indicationHeader,
                             indiHead,
                             buffer,
                            buffer_size);

    if (mdclog_level_get() >= MDCLOG_INFO) {
        uint8_t buf1[4096];
        //asn_enc_rval_t er1;
        encodebuff(ATS_BASIC_XER, &asn_DEF_E2SM_gNB_X2_indicationHeader,
                   indiHead,
                   buf1,
                   4096);

    }

    return len;
}

size_t createE2SM_gNB_X2_indicationMessage(uint8_t *message, uint msgSize,
                                                                     uint8_t *buffer,
                                                                     size_t buffer_size) {
    printEntry("E2SM_gNB_X2_indicationMessage_t", __func__)
    if (msgSize <= 0) {
        mdclog_write(MDCLOG_ERR, "E2SM_gNB_X2_indicationMessage_t failed messsage size =  %d", msgSize);
        return -1;
    }

    E2SM_gNB_X2_indicationMessage_t *indicationMessage = calloc(1, sizeof(E2SM_gNB_X2_indicationMessage_t));
    ASN_STRUCT_RESET(asn_DEF_E2SM_gNB_X2_indicationMessage, indicationMessage);

    indicationMessage->interfaceMessage.size = msgSize;
    indicationMessage->interfaceMessage.buf = calloc(1, sizeof(uint8_t) * msgSize);
    memcpy(indicationMessage->interfaceMessage.buf, message, msgSize);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_E2SM_gNB_X2_indicationMessage, indicationMessage, "E2SM_gNB_X2_indicationMessage_t", __func__);
    }

    size_t  len = encodebuff(ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_gNB_X2_indicationMessage,
                             indicationMessage,
                             buffer,
                             buffer_size);


    if (mdclog_level_get() >= MDCLOG_INFO) {
        uint8_t buf1[4096];
        //asn_enc_rval_t er1;
        encodebuff(ATS_BASIC_XER, &asn_DEF_E2SM_gNB_X2_indicationMessage,
                   indicationMessage,
                   buf1,
                   4096);

    }

    return len;
}


size_t createE2SM_gNB_X2_callProcessID(long callProcess_Id,
                                                             uint8_t *buffer,
                                                             size_t buffer_size) {
    printEntry("E2SM_gNB_X2_callProcessID_t", __func__)
    E2SM_gNB_X2_callProcessID_t *callProcessId = calloc(1, sizeof(E2SM_gNB_X2_callProcessID_t));
    ASN_STRUCT_RESET(asn_DEF_E2SM_gNB_X2_callProcessID, callProcessId);

    callProcessId->callProcess_ID = callProcess_Id;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_E2SM_gNB_X2_callProcessID, callProcessId, "E2SM_gNB_X2_indicationMessage_t", __func__);
    }

    size_t  len = encodebuff(ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_gNB_X2_callProcessID,
                             callProcessId,
                             buffer,
                             buffer_size);


    if (mdclog_level_get() >= MDCLOG_INFO) {
        uint8_t buf1[4096];
        //asn_enc_rval_t er1;
        encodebuff(ATS_BASIC_XER, &asn_DEF_E2SM_gNB_X2_callProcessID,
                   callProcessId,
                   buf1,
                   4096);

    }

    return len;
}

size_t createE2SM_gNB_X2_controlHeader(Interface_ID_t *interfaceId, long direction,
                                                             uint8_t *buffer,
                                                             size_t buffer_size) {
    printEntry("E2SM_gNB_X2_controlHeader_t", __func__)
    if (direction < InterfaceDirection_incoming || direction > InterfaceDirection_outgoing) {
        mdclog_write(MDCLOG_ERR, "E2SM_gNB_X2_controlHeader_t direction = %ld, values are %d .. %d",
                     direction, InterfaceDirection_incoming, InterfaceDirection_outgoing);
        return -1;
    }
    E2SM_gNB_X2_controlHeader_t *controlHeader = calloc(1, sizeof(E2SM_gNB_X2_controlHeader_t));
    ASN_STRUCT_RESET(asn_DEF_E2SM_gNB_X2_controlHeader, controlHeader);

    memcpy(&controlHeader->interface_ID, interfaceId, sizeof(Interface_ID_t));
    controlHeader->interfaceDirection = direction;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_E2SM_gNB_X2_controlHeader, controlHeader, "E2SM_gNB_X2_controlHeader_t", __func__);
    }

    size_t  len = encodebuff(ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_gNB_X2_controlHeader,
                             controlHeader,
                             buffer,
                             buffer_size);


    if (mdclog_level_get() >= MDCLOG_INFO) {
        uint8_t buf1[4096];
        //asn_enc_rval_t er1;
        encodebuff(ATS_BASIC_XER, &asn_DEF_E2SM_gNB_X2_controlHeader,
                   controlHeader,
                   buf1,
                   4096);

    }

    return len;
}


size_t createE2SM_gNB_X2_controlMessage(uint8_t *message, uint msgSize,
                                                               uint8_t *buffer,
                                                               size_t buffer_size) {
    printEntry("E2SM_gNB_X2_controlMessage_t", __func__)
    E2SM_gNB_X2_controlMessage_t *controlMsg = calloc(1, sizeof(E2SM_gNB_X2_controlMessage_t));
    ASN_STRUCT_RESET(asn_DEF_E2SM_gNB_X2_controlMessage, controlMsg);

    controlMsg->interfaceMessage.size = msgSize;
    controlMsg->interfaceMessage.buf = calloc(1, sizeof(uint8_t) * msgSize);
    memcpy(controlMsg->interfaceMessage.buf, message, msgSize);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        checkAndPrint(&asn_DEF_E2SM_gNB_X2_controlMessage, controlMsg, "E2SM_gNB_X2_controlMessage_t", __func__);
    }

    size_t  len = encodebuff(ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_gNB_X2_controlMessage,
                             controlMsg,
                             buffer,
                             buffer_size);


    if (mdclog_level_get() >= MDCLOG_INFO) {
        uint8_t buf1[4096];
        //asn_enc_rval_t er1;
        encodebuff(ATS_BASIC_XER, &asn_DEF_E2SM_gNB_X2_controlMessage,
                   controlMsg,
                   buf1,
                   4096);

    }

    return len;
}
