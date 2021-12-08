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

#ifndef E2_SCTPCLIENT_H
#define E2_SCTPCLIENT_H

#include <cstdio>
#include <cerrno>
#include <cstdlib>
#include <cstring>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netinet/in_systm.h>
#include <netinet/in.h>
#include <unistd.h>
#include <sys/file.h>
#include <netdb.h>
#include <sys/epoll.h>
#include <map>
#include "oranE2/E2AP-PDU.h"

#include "../httpServer/HttpServer.h"

#include "../rmrClient/rmrClient.h"
#include "cxxopts/include/cxxopts.hpp"
#include "../base64.h"
#include "../mapWrapper.h"




using namespace std;
using namespace Pistache;

#define SA struct sockaddr

#define MAXEVENTS 128
#define SCTP_BUFFER_SIZE (64*1024)

typedef enum messages {
    setupRequest_gnb,
    setupRequest_en_gNB,
    setupRequest_ng_eNB,
    setupRequest_eNB,

    nothing
} messages_t;


typedef struct SctpClient {
    string host {};
    int rmrPort{};
    int epoll_fd{};
    int sctpSock{};
    int httpBaseSocket {};
    int httpSocket {};
    mapWrapper mapKey;

    char delimiter[2] {'|', 0};

} SctpClient_t;

void createHttpLocalSocket(SctpClient_t &sctpClient);

int createEpoll(SctpClient_t &sctpClient);

inline static uint64_t rdtscp(uint32_t &aux) {
    uint64_t rax, rdx;
    asm volatile ("rdtscp\n" : "=a" (rax), "=d" (rdx), "=c" (aux) : :);
    return (rdx << (unsigned) 32) + rax;
}

int createSctpConnction(SctpClient_t *sctpClient, const char *address, int port, bool local = true);

__attribute_warn_unused_result__ int createListeningTcpConnection(SctpClient_t *sctpClient);

int modifyEpollToRead(SctpClient_t *sctpClient, int modifiedSocket);

cxxopts::ParseResult parse(SctpClient_t &sctpClient, int argc, const char *argv[]);

#endif //E2_SCTPCLIENT_H
