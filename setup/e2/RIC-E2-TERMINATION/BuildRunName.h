/*
 * Copyright 2020 AT&T Intellectual Property
 * Copyright 2020 Nokia
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


#ifndef E2_BUILDRUNNAME_H
#define E2_BUILDRUNNAME_H

#include <3rdparty/oranE2/ProtocolIE-Field.h>
#include "oranE2/ProtocolIE-Container.h"
#include "oranE2/ProtocolIE-Field.h"
#include "oranE2/GlobalE2node-gNB-ID.h"
#include "oranE2/GlobalE2node-en-gNB-ID.h"
#include "oranE2/GlobalE2node-ng-eNB-ID.h"
#include "oranE2/GlobalE2node-eNB-ID.h"

/**    02 F8 29
 * return the size of the string //
 */
static int translatePlmnId(char * plmnId, const unsigned char *data, const char* type) {
    auto mcc1 = (unsigned char)((unsigned char)data[0] & (unsigned char)0x0F);
    auto mcc2 = (unsigned char)(((unsigned char)((unsigned char)data[0] & (unsigned char)0xF0)) >> (unsigned char)4);
    ///auto mcc3 = (unsigned char)((data[1] & (unsigned char)0xF0) >> (unsigned char)4);
    auto mcc3 = (unsigned char)((unsigned char)(data[1] & (unsigned char)0x0F));

    auto mnc1 = (unsigned char)(data[2] & (unsigned char)0x0F);
    auto mnc2 =  (unsigned char)(((unsigned char)(data[2] & (unsigned char)0xF0) >> (unsigned char)4));
    //auto mnc3 = (unsigned char)(((unsigned char)(data[1] & (unsigned char)0x0F) >> (unsigned char)4) );
    auto mnc3 = (unsigned char)((data[1] & (unsigned char)0xF0) >> (unsigned char)4);

    int j = 0;
    if (mnc3 != 15) {
        j = snprintf(plmnId, 20, "%s%1d%1d%1d-%1d%1d%1d", type, mcc1, mcc2, mcc3, mnc1, mnc2, mnc3);
    }
    else {
        j = snprintf(plmnId, 20, "%s%1d%1d%1d-0%1d%1d", type, mcc1, mcc2, mcc3, mnc1, mnc2);
    }

    return j;
}

static int translateBitStringToChar(char *ranName, BIT_STRING_t &data) {
    // dont care of last unused bits
    char buffer[256] {};
    auto j = snprintf(buffer, 256, "%s-", ranName);
    memcpy(ranName, buffer, j);

    unsigned b1 = 0;
    unsigned b2 = 0;
    for (auto i = 0; i < (int)data.size; i++) {
        b1 = data.buf[i] & (unsigned)0xF0;
        b1 = b1 >> (unsigned)4;
        j = snprintf(buffer, 256, "%s%1x", ranName, b1);
        memcpy(ranName, buffer, j);
        b2 = data.buf[i] & (unsigned)0x0F;
        j = snprintf(buffer, 256, "%s%1x", ranName, b2);
        memcpy(ranName, buffer, j);
    }
    return j;
}


int buildRanName(char *ranName, E2setupRequestIEs_t *ie) {
    switch (ie->value.choice.GlobalE2node_ID.present) {
        case GlobalE2node_ID_PR_gNB: {
            auto *gnb = ie->value.choice.GlobalE2node_ID.choice.gNB;
            translatePlmnId(ranName, (const unsigned char *)gnb->global_gNB_ID.plmn_id.buf, (const char *)"gnb:");
            if (gnb->global_gNB_ID.gnb_id.present == GNB_ID_Choice_PR_gnb_ID) {
                translateBitStringToChar(ranName, gnb->global_gNB_ID.gnb_id.choice.gnb_ID);
            }
            break;
        }
        case GlobalE2node_ID_PR_en_gNB: {
            auto *enGnb = ie->value.choice.GlobalE2node_ID.choice.en_gNB;
            translatePlmnId(ranName,
                            (const unsigned char *)enGnb->global_gNB_ID.pLMN_Identity.buf,
                            (const char *)"en-gnb:");
            if (enGnb->global_gNB_ID.gNB_ID.present == ENGNB_ID_PR_gNB_ID) {
                translateBitStringToChar(ranName, enGnb->global_gNB_ID.gNB_ID.choice.gNB_ID);
            }
            break;
        }
        case GlobalE2node_ID_PR_ng_eNB: {
            auto *ngEnb = ie->value.choice.GlobalE2node_ID.choice.ng_eNB;
            char *buf = (char *)ngEnb->global_ng_eNB_ID.plmn_id.buf;
            char str[20] = {};
            BIT_STRING_t *data = nullptr;
            switch (ngEnb->global_ng_eNB_ID.enb_id.present) {
                case ENB_ID_Choice_PR_enb_ID_macro: {
                    strncpy(str, (const char *)"ng-enB-macro:", 13);
                    data = &ngEnb->global_ng_eNB_ID.enb_id.choice.enb_ID_macro;
                    break;
                }
                case ENB_ID_Choice_PR_enb_ID_shortmacro: {
                    strncpy(str, (const char *)"ng-enB-shortmacro:", 18);
                    data = &ngEnb->global_ng_eNB_ID.enb_id.choice.enb_ID_shortmacro;
                    break;
                }
                case ENB_ID_Choice_PR_enb_ID_longmacro: {
                    strncpy(str, (const char *)"ng-enB-longmacro:", 17);
                    data = &ngEnb->global_ng_eNB_ID.enb_id.choice.enb_ID_longmacro;
                }
                case ENB_ID_Choice_PR_NOTHING: {
                    break;
                }
                default:
                    break;
            }
            translatePlmnId(ranName, (const unsigned char *)buf, (const char *)str);
            translateBitStringToChar(ranName, *data);
            break;
        }
        case GlobalE2node_ID_PR_eNB: {
            auto *enb = ie->value.choice.GlobalE2node_ID.choice.eNB;
            char *buf = (char *)enb->global_eNB_ID.pLMN_Identity.buf;
            char str[20] = {};
            BIT_STRING_t *data = nullptr;

            switch (enb->global_eNB_ID.eNB_ID.present) {
                case ENB_ID_PR_macro_eNB_ID: {
                    strncpy(str, (const char *)"enB-macro:", 10);
                    data = &enb->global_eNB_ID.eNB_ID.choice.macro_eNB_ID;
                    break;
                }
                case ENB_ID_PR_home_eNB_ID: {
                    strncpy(str, (const char *)"enB-home:", 9);
                    data = &enb->global_eNB_ID.eNB_ID.choice.home_eNB_ID;
                    break;
                }
                case ENB_ID_PR_short_Macro_eNB_ID: {
                    strncpy(str, (const char *)"enB-shortmacro:", 15);
                    data = &enb->global_eNB_ID.eNB_ID.choice.short_Macro_eNB_ID;
                    break;
                }
                case ENB_ID_PR_long_Macro_eNB_ID: {
                    strncpy(str, (const char *)"enB-longmacro:", 14);
                    data = &enb->global_eNB_ID.eNB_ID.choice.long_Macro_eNB_ID;
                    break;
                }
                case ENB_ID_PR_NOTHING:
                default: {
                    break;
                }
            }
            translatePlmnId(ranName, (const unsigned char *)buf, (const char *)str);
            translateBitStringToChar(ranName, *data);
            break;
        }
        case GlobalE2node_ID_PR_NOTHING:
        default:
            return -1;
    }
    return 0;
}


#endif //E2_BUILDRUNNAME_H
