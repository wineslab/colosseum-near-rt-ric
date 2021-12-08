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

#include "E2Builder.h"

#include "asn1cFiles/ProtocolIE-Field.h"

template<typename T>
X2SetupRequest_IEs_t *buildX2SetupIE(ProtocolIE_ID_t id,
        Criticality_t criticality,
        X2SetupRequest_IEs__value_PR present,
        T *value) {
    auto *x2SetupIE = (X2SetupRequest_IEs_t *)calloc(1, sizeof(X2SetupRequest_IEs_t));
    x2SetupIE->id = id;
    x2SetupIE->criticality = criticality;
    x2SetupIE->value.present = present;

    switch (present) {
        case X2SetupRequest_IEs__value_PR_GlobalENB_ID: {
            memcpy(&x2SetupIE->value.choice.GlobalENB_ID, value, sizeof(GlobalENB_ID_t));
            break;
        }
        case X2SetupRequest_IEs__value_PR_ServedCells: {
            memcpy(&x2SetupIE->value.choice.ServedCells, value, sizeof(ServedCells_t));
            break;
        }
        case X2SetupRequest_IEs__value_PR_GUGroupIDList: {
            memcpy(&x2SetupIE->value.choice.GUGroupIDList, value, sizeof(GUGroupIDList_t));
            break;
        }
        case X2SetupRequest_IEs__value_PR_LHN_ID: {
            memcpy(&x2SetupIE->value.choice.LHN_ID, value, sizeof(LHN_ID_t));
            break;
        }
        case X2SetupRequest_IEs__value_PR_NOTHING:
        default:
            free(x2SetupIE);
            x2SetupIE = nullptr;
            break;
    }
    return x2SetupIE;
}

/**
 *
 * @param x2Setup
 * @param member
 */
void buildE2SetupRequest(X2SetupRequest_t *x2Setup, vector<X2SetupRequest_IEs_t> &member) {
    for (auto v : member) {
        ASN_SEQUENCE_ADD(&x2Setup->protocolIEs.list, &v);
    }
}

void init_log() {
    mdclog_attr_t *attr;
    mdclog_attr_init(&attr);
    mdclog_attr_set_ident(attr, "setup Request");
    mdclog_init(attr);
    mdclog_attr_destroy(attr);
}

int main(const int argc, char **argv) {
    init_log();
    //mdclog_level_set(MDCLOG_WARN);
    //mdclog_level_set(MDCLOG_INFO);
    mdclog_level_set(MDCLOG_DEBUG);

//    x2Setup	X2AP-ELEMENTARY-PROCEDURE ::= {
//            INITIATING MESSAGE		X2SetupRequest
//            SUCCESSFUL OUTCOME		X2SetupResponse
//            UNSUCCESSFUL OUTCOME	X2SetupFailure
//            PROCEDURE CODE			id-x2Setup
//            CRITICALITY				reject
//    }
//
//

//    X2SetupRequest ::= SEQUENCE {
//            protocolIEs		ProtocolIE-Container	{{X2SetupRequest-IEs}},
//            ...
//    }
//
//    X2SetupRequest-IEs X2AP-PROTOCOL-IES ::= {
//            { ID id-GlobalENB-ID			CRITICALITY reject	TYPE GlobalENB-ID			PRESENCE mandatory}|
//            { ID id-ServedCells				CRITICALITY reject	TYPE ServedCells			PRESENCE mandatory}|
//            { ID id-GUGroupIDList			CRITICALITY reject	TYPE GUGroupIDList			PRESENCE optional}|
//            { ID id-LHN-ID					CRITICALITY ignore	TYPE LHN-ID					PRESENCE optional},
//            ...
//    }



}