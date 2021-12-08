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

#ifndef ASN_DISABLE_OER_SUPPORT
#define ASN_DISABLE_OER_SUPPORT // this is to remove the OER do not remove it
#endif


#ifdef __cplusplus
extern "C"
{
#endif


#ifndef E2_E2SM_H
#define E2_E2SM_H

#include <stdlib.h>
#include <sys/types.h>
#include <string.h>
#include <error.h>
#include <errno.h>
#include <mdclog/mdclog.h>
#include <time.h>
#include <math.h>


#include "asn1cFiles/ENB-ID.h"

#include "asn1cFiles/E2SM-gNB-X2-actionDefinition.h"

#include "asn1cFiles/E2SM-gNB-X2-callProcessID.h"
#include "asn1cFiles/E2SM-gNB-X2-controlHeader.h"
#include "asn1cFiles/E2SM-gNB-X2-controlMessage.h"
#include "asn1cFiles/E2SM-gNB-X2-indicationHeader.h"
#include "asn1cFiles/E2SM-gNB-X2-indicationMessage.h"
#include "asn1cFiles/E2SM-gNB-X2-eventTriggerDefinition.h"


#include "asn1cFiles/ActionParameter-Item.h"
#include "asn1cFiles/ActionParameter-Value.h"
#include "asn1cFiles/PLMN-Identity.h"
#include "asn1cFiles/GlobalENB-ID.h"
#include "asn1cFiles/GlobalGNB-ID.h"
#include "asn1cFiles/Interface-ID.h"
#include "asn1cFiles/InterfaceMessageType.h"
#include "asn1cFiles/InterfaceProtocolIE-Item.h"

/**
 *
 * @param data
 * @return
 */
PLMN_Identity_t *createPLMN_ID(const unsigned char *data);
/**
 *
 * @param enbType
 * @param data
 * @return
 */
ENB_ID_t *createENB_ID(ENB_ID_PR enbType, unsigned char *data);
/**
 *
 * @param data
 * @param numOfBits
 * @return
 */
GNB_ID_t *createGnb_id(const unsigned char *data, int numOfBits);
/**
 *
 * @param plmnIdentity
 * @param enbId
 * @return
 */
GlobalENB_ID_t *createGlobalENB_ID(PLMN_Identity_t *plmnIdentity, ENB_ID_t *enbId);
/**
 *
 * @param plmnIdent#ifdef __cplusplus
}
#endif
ity
 * @param gnb
 * @return
 */
GlobalGNB_ID_t *createGlobalGNB_ID(PLMN_Identity_t *plmnIdentity, GNB_ID_t *gnb);

/**
 *
 * @param gnb
 * @return
 */
Interface_ID_t *createInterfaceIDForGnb(GlobalGNB_ID_t *gnb);
/**
 *
 * @param enb
 * @return
 */
Interface_ID_t *createInterfaceIDForEnb(GlobalENB_ID_t *enb);

/**
 *
 * @param procedureCode
 * @return
 */
InterfaceMessageType_t *createInterfaceMessageInitiating(ProcedureCode_t procedureCode);
/**
 *
 * @param procedureCode
 * @return
 */
InterfaceMessageType_t *createInterfaceMessageSuccsesful(ProcedureCode_t procedureCode);

/**
 *
 * @param procedureCode
 * @return
 */
InterfaceMessageType_t *createInterfaceMessageUnsuccessful(ProcedureCode_t procedureCode);


/**
 *
 * @param number
 * @return
 */
InterfaceProtocolIE_Value_t *createInterfaceProtocolValueInt(long number);
/**
 *
 * @param number
 * @return
 */
InterfaceProtocolIE_Value_t *createInterfaceProtocolValueEnum(long number);
/**
 *
 * @param val
 * @return
 */
InterfaceProtocolIE_Value_t *createInterfaceProtocolValueBool(int val);

/**
 *
 * @param buf
 * @param numOfBits
 * @return
 */
InterfaceProtocolIE_Value_t *createInterfaceProtocolValueBitString(unsigned char *buf, int numOfBits);

/**
 *
 * @param buf
 * @param size
 * @return
 */
InterfaceProtocolIE_Value_t *createInterfaceProtocolValueOCTETS(uint8_t *buf);
/**
 *
 * @param id
 * @param test
 * @param value
 * @return
 */
InterfaceProtocolIE_Item_t *createInterfaceProtocolIE_Item(long id, long test, InterfaceProtocolIE_Value_t *value);

/**
 *
 * @param number
 * @return
 */
ActionParameter_Value_t *createActionParameterValue_Int(long number);

/**
 *
 * @param number
 * @return
 */
ActionParameter_Value_t *createActionParameterValue_Enum(long number);

/**
 *
 * @param val
 * @return
 */
ActionParameter_Value_t *createActionParameterValue_Bool(int val);

/**
 *
 * @param buf
 * @param numOfBits
 * @return
 */
ActionParameter_Value_t *createActionParameterValue_Bit_String(unsigned char *buf, int numOfBits);

/**
 *
 * @param buf
 * @param size
 * @return
 */
ActionParameter_Value_t *createActionParameterValue_OCTETS(uint8_t *buf);

/**
 *
 * @param buf
 * @param size
 * @return
 */
ActionParameter_Value_t *createActionParameterValue_PRINTS(char *buf);

/**
 *
 * @param id
 * @param val
 * @return
 */
ActionParameter_Item_t *creatActionParameter_Item(long id, ActionParameter_Value_t *val);

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
                          size_t buffer_size);


/**
 *
 * @param styleId
 * @param actionParamList
 * @param listSize
 * @param buffer
 * @param buffer_size
 * @return
 */
size_t createActionDefinition(long styleId, ActionParameter_Item_t actionParamList[], int listSize,
                              uint8_t *buffer,
                              size_t buffer_size);
/**
 *
 * @param interfaceDirection
 * @param interfaceId
 * @param timestamp
 * @param size
 * @param buffer
 * @param buffer_size
 * @return
 */
size_t createE2SM_gNB_X2_indicationHeader(long interfaceDirection,
                                                                   Interface_ID_t *interfaceId,
                                                                   uint8_t *timestamp, //can put NULL if size == 0
                                                                   int size,
                                                                   uint8_t *buffer,
                                                                   size_t buffer_size);

/**
 *
 * @param message
 * @param msgSize
 * @param buffer
 * @param buffer_size
 * @return
 */
size_t createE2SM_gNB_X2_indicationMessage(uint8_t *message, uint msgSize,
                                                                     uint8_t *buffer,
                                                                     size_t buffer_size);

/**
 *
 * @param callProcess_Id
 * @param buffer
 * @param buffer_size
 * @return
 */
size_t createE2SM_gNB_X2_callProcessID(long callProcess_Id,
                                                             uint8_t *buffer,
                                                             size_t buffer_size);

size_t createE2SM_gNB_X2_controlHeader(Interface_ID_t *interfaceId, long direction,
                                                             uint8_t *buffer,
                                                             size_t buffer_size);

size_t createE2SM_gNB_X2_controlMessage(uint8_t *message, uint msgSize,
                                        uint8_t *buffer,
                                        size_t buffer_size);
#endif //E2_E2SM_H

#ifdef __cplusplus
}
#endif
