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

#include "../e2sm.h"
#include <mdclog/mdclog.h>

void init_log()
{
    mdclog_attr_t *attr;
    mdclog_attr_init(&attr);
    mdclog_attr_set_ident(attr, "e2smTests");
    mdclog_init(attr);
    mdclog_attr_destroy(attr);
}

int main(const int argc, char **argv) {
    init_log();
    //mdclog_level_set(MDCLOG_WARN);
    //mdclog_level_set(MDCLOG_INFO);
    mdclog_level_set(MDCLOG_DEBUG);

    unsigned char plmnidData[3] = {0x33, 0xF4, 0x55};

    //mdclog_write(MDCLOG_INFO, "Test PLMN_Identity_t");
    PLMN_Identity_t *plmnid = createPLMN_ID(plmnidData);

    unsigned char enbData[3] = {0x66, 0x77, 0x88};
    ENB_ID_t *enb  = createENB_ID(ENB_ID_PR_macro_eNB_ID, enbData);
    enbData[2] = 0x89;
    ENB_ID_t *enb1  = createENB_ID(ENB_ID_PR_home_eNB_ID, enbData);
    enbData[2] = 0x89;
    ENB_ID_t *enb2  = createENB_ID(ENB_ID_PR_long_Macro_eNB_ID, enbData);
    enbData[2] = 0x89;
    ENB_ID_t *enb3  = createENB_ID(ENB_ID_PR_short_Macro_eNB_ID, enbData);

    unsigned char gnbData[3] = {0x99, 0xaa, 0xbb};
    GNB_ID_t *gnb = createGnb_id(gnbData, 26);

    GlobalENB_ID_t *globalEnb = createGlobalENB_ID(plmnid, enb2);
    GlobalGNB_ID_t *globaGnb = createGlobalGNB_ID(plmnid, gnb);

    Interface_ID_t *gnbInterfaceId = createInterfaceIDForGnb(globaGnb);

    Interface_ID_t *enbInterfaceId = createInterfaceIDForEnb(globalEnb);

    InterfaceMessageType_t *initiatingInterface = createInterfaceMessageInitiating(28);

    InterfaceMessageType_t *succsesfulInterface = createInterfaceMessageSuccsesful(29);

    InterfaceMessageType_t *unSuccsesfulInterface = createInterfaceMessageUnsuccessful(29);

    InterfaceProtocolIE_Value_t *intVal = createInterfaceProtocolValueInt(88);

    InterfaceProtocolIE_Value_t *enumVal = createInterfaceProtocolValueEnum(2);

    InterfaceProtocolIE_Value_t *boolVal = createInterfaceProtocolValueBool(0);

    InterfaceProtocolIE_Value_t *bitStringVal = createInterfaceProtocolValueBitString((unsigned char *)"abcd0987", 60);

    uint8_t octe[6] = {0x11, 0x12, 0x13, 0x14, 0x15, 0x16};
    InterfaceProtocolIE_Value_t *octetsVal = createInterfaceProtocolValueOCTETS(octe);

    InterfaceProtocolIE_Item_t *item1 = createInterfaceProtocolIE_Item(10, 0, intVal);
    InterfaceProtocolIE_Item_t *item2 = createInterfaceProtocolIE_Item(10, 1, enumVal);
    InterfaceProtocolIE_Item_t *item3 = createInterfaceProtocolIE_Item(10, 0, boolVal);
    InterfaceProtocolIE_Item_t *item4 = createInterfaceProtocolIE_Item(10, 3, bitStringVal);
    InterfaceProtocolIE_Item_t *item5 = createInterfaceProtocolIE_Item(10, 4, octetsVal);

    ActionParameter_Item_t *actItem1 = creatActionParameter_Item(17, createActionParameterValue_Int(9));

    ActionParameter_Value_t *actP_enum = createActionParameterValue_Enum(5);
    ActionParameter_Item_t *actItem2 = creatActionParameter_Item(18, actP_enum);

    ActionParameter_Value_t *actP_bool = createActionParameterValue_Bool(0);
    ActionParameter_Item_t *actItem3 = creatActionParameter_Item(18, actP_bool);

    //ActionParameter_Value_t *actP_bitString = createActionParameterValue_Bit_String((unsigned char *)"ABCDEF", 42);
    ActionParameter_Item_t *actItem4 = creatActionParameter_Item(17, createActionParameterValue_Bit_String((unsigned char *)"ABCDEF", 42));

    uint8_t octe1[7] = {0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27};
    ActionParameter_Value_t *actP_octates = createActionParameterValue_OCTETS(octe1);
    ActionParameter_Item_t *actItem5 = creatActionParameter_Item(18, actP_octates);


    char print[10] = {'a', 'b', 'C', 'D', 'e', 'f', 'g', 'H', 'I', '\0'};
    ActionParameter_Value_t *actP_printable = createActionParameterValue_PRINTS(print);
    ActionParameter_Item_t *actItem6 = creatActionParameter_Item(18, actP_printable);

    InterfaceProtocolIE_Item_t *interfaceProtocolItemList[] = {item1, item2, item3, item4, item5};
    uint8_t buffer[4096];
    size_t len;
    if ((len = createEventTrigger(gnbInterfaceId,
                                       0,
                                       initiatingInterface,
                                       *interfaceProtocolItemList,
                                       5,
                                       buffer,
                                       4096))  <= 0) {
        mdclog_write(MDCLOG_ERR, "returned error from createEventTrigger");
    }

    ActionParameter_Item_t *actionParamList[] = {actItem1, actItem2, actItem3, actItem4, actItem5, actItem6};
    len = createActionDefinition(10034, *actionParamList, 6, buffer, 4096);

    len = createE2SM_gNB_X2_indicationHeader(1, enbInterfaceId, nullptr, 0, buffer, 4096);
    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "Check bad direction in header indication");
        len = createE2SM_gNB_X2_indicationHeader(2, enbInterfaceId, nullptr, 0, buffer, 4096);
        if (len == (size_t)-1) {
            mdclog_write(MDCLOG_DEBUG, "successes call function returned NULL please ignore ERROR log");
        }
    }

    uint8_t msg[] = {0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x0};

    int x = 1;
    for (int i = 0; i < x; i++) {
        len = createE2SM_gNB_X2_indicationMessage(msg, strlen((char *)msg), buffer, 4096);
    }

    for (int i = 0; i < x; i++) {
        len = createE2SM_gNB_X2_callProcessID(8, buffer, 4096);
    }

    for (int i = 0; i < x; i++) {
        len = createE2SM_gNB_X2_controlHeader(enbInterfaceId, 1, buffer, 4096);
    }

    for (int i = 0; i < x; i++) {
        len = createE2SM_gNB_X2_controlHeader(gnbInterfaceId, 1, buffer, 4096);
    }

    for (int i = 0; i < x; i++) {
        len = createE2SM_gNB_X2_controlMessage(msg, strlen((char *)msg), buffer, 4096);
    }


    return 0;
}