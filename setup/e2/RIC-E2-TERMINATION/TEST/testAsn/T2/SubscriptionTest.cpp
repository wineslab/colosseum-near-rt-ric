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
// Created by adi ENZEL on 6/17/19.
//

#include <mdclog/mdclog.h>

#include "asn/type_defs.h"
#include "asn/per/codec.hpp"
#include "asn/printer.hpp"

#include "X2AP-CommonDataTypes.hpp"
#include "X2AP-Containers.hpp"
#include "X2AP-Constants.hpp"
#include "X2AP-IEs.hpp"
#include "X2AP-PDU-Contents.hpp"

#include "E2AP-Constants.hpp"
#include "E2AP-IEs.hpp"
#include "E2AP-PDU-Contents.hpp"
#include "E2AP-PDU-Descriptions.hpp"


#include <iostream>
#include <cstdio>
#include <cctype>
#include <cstring>

#include <pthread.h>
#include <rmr/rmr.h>
#include <rmr/RIC_message_types.h>

#include "logInit.h"

// test X2SetUP request and response
using namespace std;

#define MAXEVENTS 64

int main(const int argc, char **argv) {
    mdclog_severity_t loglevel = MDCLOG_INFO;

    auto buff = new string("Subscription TEST");
    init_log((char *)buff->c_str());

    mdclog_level_set(loglevel);

    if (argc < 5){
        mdclog_mdc_add("app", argv[0]);
        mdclog_write(MDCLOG_ERR, "Usage ran <ran name> rmr <rmr address> [logLevel <debug/warning/info/error]");
        return -1 ;
    }

    char ranName[256] {0};
    char rmrAddress[256] {0};

    char str1[128];
    for (int i = 1; i < argc; i += 2) {
        for (int j = 0; j < strlen(argv[i]); j++) {
            str1[j] = (char)tolower(argv[i][j]);
        }
        str1[strlen(argv[i])] = 0;
        if (strcmp("ran", str1) == 0) {
            strcpy(ranName, argv[i + 1]);
        } else if (strcmp("rmr", str1) == 0) {
            strcpy(rmrAddress, argv[i + 1]);
        } else if (strcmp("loglevel", str1) == 0) {
            if (strcmp("debug", argv[i + 1]) == 0) {
                loglevel = MDCLOG_DEBUG;
            } else if (strcmp("info", argv[i + 1]) == 0) {
                loglevel = MDCLOG_INFO;
            } else if (strcmp("warning", argv[i + 1]) == 0) {
                loglevel = MDCLOG_WARN;
            } else if (strcmp("error", argv[i + 1]) == 0) {
                loglevel = MDCLOG_ERR;
            }
        }
    }

    void *rmrCtx = rmr_init(rmrAddress, RMR_MAX_RCV_BYTES, RMRFL_NONE);
    if (rmrCtx == nullptr ) {
        mdclog_write(MDCLOG_ERR, "RMR failed to initialise : %s", strerror(errno));
        return(-1);
    }

    // get the RMR fd for the epoll
    auto rmrListenFd = rmr_get_rcvfd(rmrCtx);

    auto epoll_fd = epoll_create1(0);
    if (epoll_fd == -1) {
        mdclog_write(MDCLOG_ERR,"failed to open epoll descriptor");
        rmr_close(rmrCtx);
        return -2;
    }

    struct epoll_event event {};
    event.events = EPOLLIN;
    event.data.fd = rmrListenFd;
    // add listening sctpPort to epoll
    if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, rmrListenFd, &event)) {
        mdclog_write(MDCLOG_ERR, "Failed to add RMR descriptor to epoll");
        close(rmrListenFd);
        rmr_close(rmrCtx);
        return -3;
    }


    // we need to find that routing table exist and we can run
    if (mdclog_level_get() >= MDCLOG_INFO) {
        mdclog_write(MDCLOG_INFO, "We are after RMR INIT wait for RMR_Ready");
    }

    int rmrReady = 0;
    int count = 0;
    while (!rmrReady) {
        if ((rmrReady = rmr_ready(rmrCtx)) == 0) {
            sleep(1);
        }
        count++;
        if (count % 60 == 0) {
            mdclog_write(MDCLOG_INFO, "waiting to RMR ready state for %d seconds", count);
        }
        if (count > 180) {
            mdclog_write(MDCLOG_ERR, "RMR not ready tried for 3 minutes ");
            return(-2);
        }
    }
    if (mdclog_level_get() >= MDCLOG_INFO) {
        mdclog_write(MDCLOG_INFO, "RMR running");
    }

    E2AP_PDU cuAckPdu {};
    auto &succ = cuAckPdu.select_successfulOutcome();
    succ.ref_procedureCode().select_id_endcConfigurationUpdate();
    succ.ref_criticality().select_id_endcConfigurationUpdate();
    auto &endcConfAck = succ.ref_value().select_id_endcConfigurationUpdate();
    auto &confAckIes =  endcConfAck.ref_protocolIEs();

    ENDCConfigurationUpdateAcknowledge::protocolIEs_t::value_type endcENB {};
    endcENB.ref_id().select_id_RespondingNodeType_EndcConfigUpdate();
    endcENB.ref_criticality().select_id_RespondingNodeType_EndcConfigUpdate();

    auto &respondingNode = endcENB.ref_value().select_id_RespondingNodeType_EndcConfigUpdate();

    auto &enb = respondingNode.select_respond_eNB();


    confAckIes.push_back(endcENB);



    E2AP_PDU pdu {};
    auto &initiatingMsg = pdu.select_initiatingMessage();
    initiatingMsg.ref_procedureCode().select_id_ricSubscription();
    initiatingMsg.ref_criticality().select_id_ricSubscription();
    auto &subscription = initiatingMsg.ref_value().select_id_ricSubscription();

    auto &ies = subscription.ref_protocolIEs();



    RICsubscriptionRequest::protocolIEs_t::value_type ranFuncId {};
    ranFuncId.ref_id().select_id_RANfunctionID();
    ranFuncId.ref_criticality().select_id_RANfunctionID();
    ranFuncId.ref_value().select_id_RANfunctionID().set(28);
    ies.push_back(ranFuncId);

    RICsubscriptionRequest::protocolIEs_t::value_type ricRequestId {};
    ricRequestId.ref_id().select_id_RICrequestID();
    ricRequestId.ref_criticality().select_id_RICrequestID();
    ricRequestId.ref_value().select_id_RICrequestID().ref_ricRequestorID().set(44);
    ricRequestId.ref_value().select_id_RICrequestID().ref_ricRequestSequenceNumber().set(55);
    ies.push_back(ricRequestId);

    RICsubscriptionRequest::protocolIEs_t::value_type ricSubId {};
    ricSubId.ref_id().select_id_RICsubscription();
    ricSubId.ref_criticality().select_id_RICsubscription();

    //E2SM_gNB_X2_eventTriggerDefinition_t evt;
    uint8_t v1[] = {0x02, 0xf8, 0x29, 0x88};

    RICeventTriggerDefinition eventTriggerDef {};  // octet string in E2AP but struct in E2SM
    eventTriggerDef.set(4, v1);

//    eventTriggerDef.
//
//
//    RICaction_Admitted_List::value_type actbl {};
//    actbl.ref_id().select_id_RICaction_Admitted_Item();
//    actbl.ref_criticality().select_id_RICaction_Admitted_Item();


//    RICaction_ToBeSetup_Item actb1{};
//
//    actbl.ref_value().select_id_RICaction_Admitted_Item().ref_ricActionID().set(actb1);
//    ricSubId.ref_value().select_id_RICsubscription().ref_ricAction_ToBeSetup_List().set(actbl);


    ies.push_back(ricSubId);

/*

    ies.push_back(ranFuncId);

    X2SetupRequest::protocolIEs_t::value_type sc {};
    ies.push_back(sc);

    sc.ref_id().select_id_ServedCells();
    sc.ref_criticality().select_id_ServedCells();

    ServedCells::value_type sce;
    sc.ref_value().select_id_ServedCells().push_back(sce);

    sce.ref_servedCellInfo().ref_pCI().set(0x1F7);
    uint8_t v3[] = {0x1, 0x2};
    sce.ref_servedCellInfo().ref_tAC().set(2,v3);
    sce.ref_servedCellInfo().ref_cellId().ref_pLMN_Identity().set(3, v1);
    uint8_t v4[] = {0x00, 0x07, 0xab, ((unsigned)0x50) >> (unsigned)4};
    sce.ref_servedCellInfo().ref_cellId().ref_eUTRANcellIdentifier().set_buffer(28, v4);

    BroadcastPLMNs_Item::value_type bpe;
    sce.ref_servedCellInfo().ref_broadcastPLMNs().push_back(bpe);
    bpe.set(3, v1);

    sce.ref_servedCellInfo().ref_eUTRA_Mode_Info().select_fDD().ref_uL_EARFCN().set(0x1);
    sce.ref_servedCellInfo().ref_eUTRA_Mode_Info().select_fDD().ref_dL_EARFCN().set(0x1);
    sce.ref_servedCellInfo().ref_eUTRA_Mode_Info().select_fDD().ref_uL_Transmission_Bandwidth().set(Transmission_Bandwidth::bw50);
    sce.ref_servedCellInfo().ref_eUTRA_Mode_Info().select_fDD().ref_dL_Transmission_Bandwidth().set(Transmission_Bandwidth::bw50);

*/

    unsigned char s_buffer[64 * 1024];
    asn::per::EncoderCtx ctx{s_buffer, sizeof(s_buffer)};
    std::cout << asn::get_printed(pdu) << std::endl;
    if (!asn::per::pack(pdu, ctx)) {
        std::cout << ctx.refErrorCtx().toString() << std::endl;
        return -3;
    }
    size_t packed_buf_size;
    packed_buf_size = static_cast<size_t>(ctx.refBuffer().getBytesUsed());

    // build message
    char data[4096] {};
    //auto delimiter = (const char) '|';
    sprintf(data, "%s/0", ctx.refBuffer().getBytes(packed_buf_size));

    rmr_mbuf_t *msg = rmr_alloc_msg(rmrCtx, int(strlen(data)));
    rmr_bytes2meid(msg, (unsigned char const*)ranName, strlen(ranName));
    rmr_bytes2payload(msg, (unsigned char const*)data, strlen(data));
    rmr_bytes2xact(msg, (unsigned char const*)ranName, strlen(ranName));
    msg->mtype = RIC_SUB_REQ;
    msg->state = 0;

    msg = rmr_send_msg(rmrCtx, msg);
    if (msg->state != 0) {
        mdclog_write(MDCLOG_ERR, "Message state %d while sending RIC_X2_SETUP to %s", msg->state, ranName);
        rmr_free_msg(msg);
        rmr_close(rmrCtx);
        return -4;
    }
    rmr_free_msg(msg);


    unsigned char allocBuffer[64*1024] {0};
    auto *events = (struct epoll_event *)calloc(MAXEVENTS, sizeof(event));

    while (true) {

        auto numOfEvents = epoll_wait(epoll_fd, events, MAXEVENTS, -1);
        if (numOfEvents < 0) {
            mdclog_write(MDCLOG_ERR, "Epoll wait failed, errno = %s", strerror(errno));
            rmr_close(rmrCtx);
            return -4;
        }
        for (auto i = 0; i < numOfEvents; i++) {
            if ((events[i].events & EPOLLERR) || (events[i].events & EPOLLHUP) || (!(events[i].events & EPOLLIN))) {
                mdclog_write(MDCLOG_ERR, "epoll error");
            } else if (rmrListenFd == events[i].data.fd) {
                msg = rmr_alloc_msg(rmrCtx, 4096);
                if (msg == nullptr) {
                    mdclog_write(MDCLOG_ERR, "RMR Allocation message, %s", strerror(errno));
                    rmr_close(rmrCtx);
                    return -5;
                }

                msg = rmr_rcv_msg(rmrCtx, msg);
                if (msg == nullptr) {
                    mdclog_write(MDCLOG_ERR, "RMR Receving message, %s", strerror(errno));
                    rmr_close(rmrCtx);
                    return -6;
                }
                memset(allocBuffer, 0, 64*1024);
                switch (msg->mtype) {
                    default: {
                        mdclog_write(MDCLOG_INFO, "RMR receiveing message type %d", msg->mtype);
                        asn::per::DecoderCtx dCtx{msg->payload, (size_t) msg->len, allocBuffer, sizeof(allocBuffer)};
                        E2AP_PDU opdu;
                        if (!asn::per::unpack(opdu, dCtx)) {
                            mdclog_write(MDCLOG_ERR, "Failed to unpack ASN message, %s", dCtx.refErrorCtx().toString());
                            rmr_close(rmrCtx);
                            return -7;
                        }

                        switch (opdu.get_index()) {
                            case 1: { //initiating message
                                mdclog_write(MDCLOG_INFO, "ASN initiating message type %ld",
                                             opdu.get_initiatingMessage()->ref_procedureCode().ref_value().get());
                                break;
                            }
                            case 2: { //successful message
                                mdclog_write(MDCLOG_INFO, "ASN initiating message type %ld",
                                             opdu.get_successfulOutcome()->ref_procedureCode().ref_value().get());
                                break;
                            }
                            case 3: { //unsuccessesful message
                                mdclog_write(MDCLOG_INFO, "ASN initiating message type %ld",
                                             opdu.get_unsuccessfulOutcome()->ref_procedureCode().ref_value().get());
                                break;
                            }

                        }
                        mdclog_write(MDCLOG_INFO, "RMR receiveing message from E2 terminator, %d",
                                     msg->mtype);
                        break;
                    }
                }
            }
        }
    }
}