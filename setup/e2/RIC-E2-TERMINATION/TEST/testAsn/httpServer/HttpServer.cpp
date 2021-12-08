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

//
// Created by adi ENZEL on 2/16/20.
//

#include "HttpServer.h"
#include <algorithm>
#include <random>


#include "../sctpClient/sctpClient.h"

#include "../T1/E2Builder.h"
#include "../base64.h"

using namespace std;
using namespace Pistache;

#define RECEIVE_SCTP_BUFFER_SIZE 8192

namespace Generic {

    void handleReady(const Rest::Request&, Http::ResponseWriter response) {
        response.send(Http::Code::Ok, "1");
    }
}


    HttpServer::HttpServer(Address addr)
            : httpBaseSocket(0), httpEndpoint(std::make_shared<Http::Endpoint>(addr)) { }

    void HttpServer::init(size_t thr) {
        if ((httpBaseSocket = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
            fprintf(stderr, "Socket() error. %s\n", strerror(errno));
            exit(-1);
        }
        auto optval = 1;
        if (setsockopt(httpBaseSocket, SOL_SOCKET, SO_REUSEPORT, &optval, sizeof optval) != 0) {
            fprintf(stderr, "setsockopt SO_REUSEPORT Error, %s %s, %d\n", strerror(errno), __func__, __LINE__);
            close(httpBaseSocket);
            exit(-1);
        }
        optval = 1;
        if (setsockopt(httpBaseSocket, SOL_SOCKET, SO_REUSEADDR, &optval, sizeof optval) != 0) {
            fprintf(stderr, "setsockopt SO_REUSEADDR Error, %s %s, %d\n", strerror(errno), __func__, __LINE__);
            close(httpBaseSocket);
            exit(-1);
        }

        struct sockaddr_in address{};
        address.sin_family = AF_INET;
        if(inet_pton(AF_INET, "127.0.0.1", &address.sin_addr)<=0)
        {
            fprintf(stderr,"Invalid address/Address not supported. %s", strerror(errno));
            exit(-1);
        }


        address.sin_port = htons(9098);
        if (connect(httpBaseSocket, (SA *)(&address), sizeof(address)) < 0) {
            fprintf(stderr, "connect() error. %s\n", strerror(errno));
            exit(-1);
        }
        auto opts = Http::Endpoint::options().threads(thr);
        httpEndpoint->init(opts);
        setupRoutes();
    }

    void HttpServer::start() {
        std::random_device device{};
        std::mt19937 generator(device());
        std::uniform_int_distribution<long> distribution(1, (long) 1e12);
        transactionCounter = distribution(generator);


        httpEndpoint->setHandler(router.handler());
        httpEndpoint->serve();
    }

    void HttpServer::setupRoutes() {
        using namespace Rest;

        Routes::Get(router, "/setup/:ricaddress/:ricPort/:mcc/:mnc", Routes::bind(&HttpServer::sendSetupReq, this));
        //Routes::Post(router, "/ricIndication/:ricid/:subscriptionId/:mcc/:mnc", Routes::bind(&HttpServer::sendSetupReq, this));
        Routes::Get(router, "/ready", Routes::bind(&Generic::handleReady));
    }


    void HttpServer::sendSetupReq(const Rest::Request& request, Http::ResponseWriter response) {
        auto mcc = request.param(":mcc").as<int>();
        auto mnc = request.param(":mnc").as<int>();
        auto ricAdress = request.param(":ricaddress").as<std::string>();
        auto ricPort = request.param(":ricPort").as<int>();
        //TODO  build setup to send to address
        E2AP_PDU_t pdu;

        buildSetupRequest(&pdu,mcc, mnc);
        // encode PDU to PER

        auto buffer_size =  RECEIVE_SCTP_BUFFER_SIZE;
        unsigned char buffer[RECEIVE_SCTP_BUFFER_SIZE] = {};
        // encode to xml
        asn_enc_rval_t er;
        er = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2AP_PDU, &pdu, buffer, buffer_size);
        if (er.encoded == -1) {
           cerr << "encoding of : " <<  asn_DEF_E2AP_PDU.name << " failed, "<< strerror(errno) << endl;
           response.send(Http::Code::Internal_Server_Error, "strerror(errno)");
           return;
        } else if (er.encoded > (ssize_t)buffer_size) {
            cerr << "Buffer of size : " << buffer_size << " is to small for : " << asn_DEF_E2AP_PDU.name << endl;
            response.send(Http::Code::Internal_Server_Error, "Buffer of size is too small");
            return;
        }

        long len = er.encoded * 4 / 3 + 128;
        auto *base64Buff = (unsigned char *)calloc(1,len + 1024);
        char tx[32];
        snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);

        auto sentLen = snprintf((char *)base64Buff, 1024, "%d|%s|%s|%d|", setupRequest_gnb, tx, ricAdress.c_str(), ricPort);

        base64::encode(buffer, er.encoded, &base64Buff[sentLen], len);
        sentLen += len;
        len = send(httpBaseSocket, base64Buff, sentLen, 0);
        if (len < 0) {
            cerr << "failed sending setupRequest_gnb to Other thread. Error : " << strerror(errno) << endl;
            response.send(Http::Code::Internal_Server_Error, "Failed send buffer");
            free(base64Buff);
            return;
        }
        char tx1[128];
        snprintf((char *) tx1, sizeof tx1, "{\"id\": %s}", tx);
        response.send(Http::Code::Ok, tx1);
        free(base64Buff);
    }
