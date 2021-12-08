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

#include <cstring>
#include <cstdio>
#include <cerrno>
#include <cstdlib>
#include <iostream>


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

#include "E2Builder.h"

using namespace std;


#include "BuildRunName.h"

void buildRanName(E2AP_PDU_t *pdu, unsigned char *buffer) {
    for (auto i = 0; i < pdu->choice.initiatingMessage->value.choice.E2setupRequest.protocolIEs.list.count; i++) {
        auto *ie = pdu->choice.initiatingMessage->value.choice.E2setupRequest.protocolIEs.list.array[i];
        if (ie->id == ProtocolIE_ID_id_GlobalE2node_ID) {
            if (ie->value.present == E2setupRequestIEs__value_PR_GlobalE2node_ID) {
                memset(buffer, 0, 128);
                buildRanName( (char *) buffer, ie);
            }
        }
    }

}

void extractPdu(E2AP_PDU_t *pdu, unsigned char *buffer, int buffer_size) {
    asn_enc_rval_t er;
    er = asn_encode_to_buffer(nullptr, ATS_BASIC_XER, &asn_DEF_E2AP_PDU, pdu, buffer, buffer_size);
    if (er.encoded == -1) {
        cerr << "encoding of " << asn_DEF_E2AP_PDU.name << " failed, " << strerror(errno) << endl;
        exit(-1);
    } else if (er.encoded > (ssize_t) buffer_size) {
        cerr << "Buffer of size " << buffer_size << " is to small for " << asn_DEF_E2AP_PDU.name << endl;
        exit(-1);
    } else {
        cout << "XML result = " << buffer << endl;
    }
}


std::string otherXml = "<E2AP-PDU>\n"
                       "    <successfulOutcome>\n"
                       "        <procedureCode>1</procedureCode>\n"
                       "        <criticality><reject/></criticality>\n"
                       "        <value>\n"
                       "            <E2setupResponse>\n"
                       "                <protocolIEs>\n"
                       "                    <E2setupResponseIEs>\n"
                       "                        <id>4</id>\n"
                       "                        <criticality><reject/></criticality>\n"
                       "                        <value>\n"
                       "                            <GlobalRIC-ID>\n"
                       "                                <pLMN-Identity>13 10 14</pLMN-Identity>\n"
                       "                                <ric-ID>\n"
                       "                                    10011001101010101011"
                       "                                </ric-ID>\n"
                       "                            </GlobalRIC-ID>\n"
                       "                        </value>\n"
                       "                    </E2setupResponseIEs>\n"
                       "                    <E2setupResponseIEs>\n"
                       "                        <id>9</id>\n"
                       "                        <criticality><reject/></criticality>\n"
                       "                        <value>\n"
                       "                            <RANfunctionsID-List>\n"
                       "                                <ProtocolIE-SingleContainer>\n"
                       "                                    <id>6</id>\n"
                       "                                    <criticality><ignore/></criticality>\n"
                       "                                    <value>\n"
                       "                                        <RANfunctionID-Item>\n"
                       "                                            <ranFunctionID>1</ranFunctionID>\n"
                       "                                            <ranFunctionRevision>1</ranFunctionRevision>\n"
                       "                                        </RANfunctionID-Item>\n"
                       "                                    </value>\n"
                       "                                </ProtocolIE-SingleContainer>\n"
                       "                                <ProtocolIE-SingleContainer>\n"
                       "                                    <id>6</id>\n"
                       "                                    <criticality><ignore/></criticality>\n"
                       "                                    <value>\n"
                       "                                        <RANfunctionID-Item>\n"
                       "                                            <ranFunctionID>2</ranFunctionID>\n"
                       "                                            <ranFunctionRevision>1</ranFunctionRevision>\n"
                       "                                        </RANfunctionID-Item>\n"
                       "                                    </value>\n"
                       "                                </ProtocolIE-SingleContainer>\n"
                       "                                <ProtocolIE-SingleContainer>\n"
                       "                                    <id>6</id>\n"
                       "                                    <criticality><ignore/></criticality>\n"
                       "                                    <value>\n"
                       "                                        <RANfunctionID-Item>\n"
                       "                                            <ranFunctionID>3</ranFunctionID>\n"
                       "                                            <ranFunctionRevision>1</ranFunctionRevision>\n"
                       "                                        </RANfunctionID-Item>\n"
                       "                                    </value>\n"
                       "                                </ProtocolIE-SingleContainer>\n"
                       "                            </RANfunctionsID-List>\n"
                       "                        </value>\n"
                       "                    </E2setupResponseIEs>\n"
                       "                </protocolIEs>\n"
                       "            </E2setupResponse>\n"
                       "        </value>\n"
                       "    </successfulOutcome>\n"
                       "</E2AP-PDU>\n";



std::string newXml =
        "<E2AP-PDU><successfulOutcome><procedureCode>1</procedureCode><criticality><reject/></criticality><value><E2setupResponse><protocolIEs><E2setupResponseIEs><id>4</id><criticality><reject/></criticality><value><GlobalRIC-ID><pLMN-Identity>13 10 14</pLMN-Identity><ric-ID>10101010110011001110</ric-ID></GlobalRIC-ID></value></E2setupResponseIEs><E2setupResponseIEs><id>9</id><criticality><reject/></criticality><value><RANfunctionsID-List><ProtocolIE-SingleContainer><id>6</id><criticality><ignore/></criticality><value><RANfunctionID-Item><ranFunctionID>1</ranFunctionID><ranFunctionRevision>1</ranFunctionRevision></RANfunctionID-Item></value></ProtocolIE-SingleContainer><ProtocolIE-SingleContainer><id>6</id><criticality><ignore/></criticality><value><RANfunctionID-Item><ranFunctionID>2</ranFunctionID><ranFunctionRevision>1</ranFunctionRevision></RANfunctionID-Item></value></ProtocolIE-SingleContainer><ProtocolIE-SingleContainer><id>6</id><criticality><ignore/></criticality><value><RANfunctionID-Item><ranFunctionID>3</ranFunctionID><ranFunctionRevision>1</ranFunctionRevision></RANfunctionID-Item></value></ProtocolIE-SingleContainer></RANfunctionsID-List></value></E2setupResponseIEs></protocolIEs></E2setupResponse></value></successfulOutcome></E2AP-PDU>";
std::string setupFailure = "<E2AP-PDU>"
                             "<unsuccessfulOutcome>"
                               "<procedureCode>1</procedureCode>"
                               "<criticality><reject/></criticality>"
                               "<value>"
                                 "<E2setupFailure>"
                                   "<protocolIEs>"
                                     "<E2setupFailureIEs>"
                                       "<id>1</id>"
                                       "<criticality><reject/></criticality>"
                                       "<value>"
                                         "<Cause>"
                                           "<transport>"
                                             "<transport-resource-unavailable/>"
                                           "</transport>"
                                         "</Cause>"
                                       "</value>"
                                     "</E2setupFailureIEs>"
                                   "</protocolIEs>"
                                 "</E2setupFailure>"
                               "</value>"
                             "</unsuccessfulOutcome>"
                           "</E2AP-PDU>";



std::string otherSucc = "  <E2AP-PDU>"
                        "<successfulOutcome>"
                        "<procedureCode>1</procedureCode>"
                        "<criticality><reject/></criticality>"
                        "<value>"
                        "<E2setupResponse>"
                        "<protocolIEs>"
                        "<E2setupResponseIEs>"
                        "<id>4</id>"
                        "<criticality><reject/></criticality>"
                        "<value>"
                        "<GlobalRIC-ID>"
                        "<pLMN-Identity>131014</pLMN-Identity>"
                        "<ric-ID>10101010110011001110</ric-ID>"
                        "</GlobalRIC-ID>"
                        "</value>"
                        "</E2setupResponseIEs><E2setupResponseIEs><id>9</id><criticality><reject/></criticality><value><RANfunctionsID-List><ProtocolIE-SingleContainer><id>6</id><criticality><ignore/></criticality><value><RANfunctionID-Item><ranFunctionID>1</ranFunctionID><ranFunctionRevision>1</ranFunctionRevision></RANfunctionID-Item></value></ProtocolIE-SingleContainer><ProtocolIE-SingleContainer><id>6</id><criticality><ignore/></criticality><value><RANfunctionID-Item><ranFunctionID>2</ranFunctionID><ranFunctionRevision>1</ranFunctionRevision></RANfunctionID-Item></value></ProtocolIE-SingleContainer><ProtocolIE-SingleContainer><id>6</id><criticality><ignore/></criticality><value><RANfunctionID-Item><ranFunctionID>3</ranFunctionID><ranFunctionRevision>1</ranFunctionRevision></RANfunctionID-Item></value></ProtocolIE-SingleContainer></RANfunctionsID-List></value></E2setupResponseIEs></protocolIEs></E2setupResponse></value></successfulOutcome></E2AP-PDU>";

auto main(const int argc, char **argv) -> int {
    E2AP_PDU_t pdu;
    char *printBuffer;
    size_t size;
    FILE *stream = open_memstream(&printBuffer, &size);
    auto buffer_size =  8192;
    unsigned char buffer[8192] = {};

    E2AP_PDU_t *XERpdu  = nullptr;
    cout << "message size = " <<  otherSucc.length() << endl;
    auto rval = asn_decode(nullptr, ATS_BASIC_XER, &asn_DEF_E2AP_PDU, (void **) &XERpdu,
                           otherSucc.c_str(), otherSucc.length());
    if (rval.code != RC_OK) {
        cout <<  "Error " << rval.code << " (unpack) setup response " << endl;
        //return -1;
    }

    asn_fprint(stream, &asn_DEF_E2AP_PDU, XERpdu);
    cout << "Encoding E2AP PDU of size  " << size << endl << printBuffer << endl;
    fseek(stream,0,SEEK_SET);

//    cout << "=========================" << endl << otherXml << endl << "========================" << endl;

    buildSetupRequest(&pdu, 311, 410);
    asn_fprint(stream, &asn_DEF_E2AP_PDU, &pdu);
    cout << "Encoding E2AP PDU of size  " << size << endl << printBuffer << endl;
    fseek(stream,0,SEEK_SET);

    extractPdu(&pdu, buffer, buffer_size);
    buildRanName(&pdu, buffer);
    cout << "Ran name = " << buffer << endl;

    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, &pdu);
    memset(buffer, 0, buffer_size);

    cout << "========== en-gnb ============" << endl;
    buildSetupRequesteenGNB(&pdu, 320, 512);
    asn_fprint(stream, &asn_DEF_E2AP_PDU, &pdu);
    cout << "Encoding E2AP PDU of size  " << size << endl << printBuffer << endl;
    fseek(stream,0,SEEK_SET);
    cout << "========== en-gnb ============" << endl;

    extractPdu(&pdu, buffer, buffer_size);
    buildRanName(&pdu, buffer);
    cout << "Ran name = " << buffer << endl;

    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, &pdu);
    memset(buffer, 0, buffer_size);

    buildSetupRequestWithFunc(&pdu, 311, 410);
    extractPdu(&pdu, buffer, buffer_size);

    buildRanName(&pdu, buffer);
    cout << "Ran name = " << buffer << endl;

    cout << "Sucessesfull outcome" << endl;
    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, &pdu);
    memset(buffer, 0, buffer_size);
    uint8_t data[4] = {0x99, 0xAA, 0xBB, 0};

    buildSetupSuccsessfulResponse(&pdu, 311, 410, data);

    asn_fprint(stream, &asn_DEF_E2AP_PDU, &pdu);
    cout << "Encoding E2AP PDU of size  " << size << endl << printBuffer << endl;
    fseek(stream,0,SEEK_SET);

    extractPdu(&pdu, buffer, buffer_size);


    cout << "Failure outcome" << endl;
    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, &pdu);
    memset(buffer, 0, buffer_size);

    buildSetupUnSuccsessfulResponse(&pdu);
    asn_fprint(stream, &asn_DEF_E2AP_PDU, &pdu);
    cout << "Encoding E2AP PDU of size  " << size << endl << printBuffer << endl;
    fseek(stream,0,SEEK_SET);

    extractPdu(&pdu, buffer, buffer_size);
}
