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

#ifndef X2_SCTP_THREAD_H
#define X2_SCTP_THREAD_H

#include <algorithm>

#include <cstdio>
#include <cerrno>
#include <cstdlib>
#include <cstring>
#include <random>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netinet/in_systm.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <netinet/ip_icmp.h>
#include <netinet/sctp.h>
#include <thread>
#include <atomic>
#include <sys/param.h>
#include <sys/file.h>
#include <ctime>
#include <netdb.h>
#include <sys/epoll.h>
#include <mutex>
#include <shared_mutex>
#include <iterator>
#include <map>
#include <sys/inotify.h>
#include <csignal>

#include <rmr/rmr.h>
#include <rmr/RIC_message_types.h>
#include <mdclog/mdclog.h>
#include <functional>
#include <iostream>

#include <boost/algorithm/string/predicate.hpp>
#include <boost/lexical_cast.hpp>
#include <boost/move/utility.hpp>
#include <boost/log/sources/logger.hpp>
#include <boost/log/sources/record_ostream.hpp>
#include <boost/log/sources/global_logger_storage.hpp>
#include <boost/log/utility/setup/file.hpp>
#include <boost/log/utility/setup/common_attributes.hpp>
#include <boost/filesystem.hpp>

#include <mdclog/mdclog.h>

#include "oranE2/E2AP-PDU.h"
#include "oranE2/ProtocolIE-Container.h"
#include "oranE2/InitiatingMessage.h"
#include "oranE2/SuccessfulOutcome.h"
#include "oranE2/UnsuccessfulOutcome.h"
#include "oranE2/ProtocolIE-Container.h"
#include "oranE2/ProtocolIE-Field.h"
#include "oranE2/GlobalE2node-gNB-ID.h"
#include "oranE2/GlobalE2node-en-gNB-ID.h"
#include "oranE2/GlobalE2node-ng-eNB-ID.h"
#include "oranE2/GlobalE2node-eNB-ID.h"

#include "cxxopts.hpp"
//#include "config-cpp/include/config-cpp/config-cpp.h"


#include "mapWrapper.h"
#include "statCollector.h"

#include "base64.h"

#include "ReadConfigFile.h"

using namespace std;
namespace logging = boost::log;
namespace src = boost::log::sources;
namespace keywords = boost::log::keywords;
namespace sinks = boost::log::sinks;
namespace posix_time = boost::posix_time;
namespace expr = boost::log::expressions;

#define SRC_PORT 36422
#define SA      struct sockaddr
#define MAX_ENODB_NAME_SIZE 64

#define MAXEVENTS 128

#define RECEIVE_SCTP_BUFFER_SIZE (128 * 1024)
#define RECEIVE_XAPP_BUFFER_SIZE RECEIVE_SCTP_BUFFER_SIZE 

typedef mapWrapper Sctp_Map_t;



#define VOLUME_URL_SIZE 256
#define KA_MESSAGE_SIZE 2048

typedef struct sctp_params {
    uint16_t rmrPort = 0;
    uint16_t sctpPort = SRC_PORT;
    int      epoll_fd = 0;
    int      listenFD = 0;
    int      rmrListenFd = 0;
    int      inotifyFD = 0;
    int      inotifyWD = 0;
    void     *rmrCtx = nullptr;
    Sctp_Map_t *sctpMap = nullptr;
    char      ka_message[KA_MESSAGE_SIZE] {};
    int       ka_message_length = 0;
    char       rmrAddress[256] {}; // "tcp:portnumber" "tcp:5566" listen to all address on port 5566
    mdclog_severity_t logLevel = MDCLOG_INFO;
    char volume[VOLUME_URL_SIZE];
    string myIP {};
    string fqdn {};
    string podName {};
    string configFilePath {};
    string configFileName {};
    bool trace = true;
    //shared_timed_mutex fence; // moved to mapWrapper
} sctp_params_t;

typedef struct ConnectedCU {
    int fileDescriptor = 0;
    char hostName[NI_MAXHOST] {};
    char portNumber[NI_MAXSERV] {};
    char enodbName[MAX_ENODB_NAME_SIZE] {};
    char asnData[RECEIVE_SCTP_BUFFER_SIZE] {};
    size_t asnLength = 0;
    int mtype = 0;
    bool isConnected = false;
    bool gotSetup = false;
    sctp_params_t *sctpParams = nullptr;
} ConnectedCU_t ;

#define MAX_RMR_BUFF_ARRY 32
typedef struct RmrMessagesBuffer {
    char ka_message[KA_MESSAGE_SIZE] {};
    int  ka_message_len = 0;
    void *rmrCtx = nullptr;
    rmr_mbuf_t *sendMessage= nullptr;
    //rmr_mbuf_t *sendBufferedMessages[MAX_RMR_BUFF_ARRY] {};
    rmr_mbuf_t *rcvMessage= nullptr;
    //rmr_mbuf_t *rcvBufferedMessages[MAX_RMR_BUFF_ARRY] {};
} RmrMessagesBuffer_t;

typedef struct formatedMessage {
    char enodbName[MAX_ENODB_NAME_SIZE];
    struct timespec time;
    int messageType;
    char direction;
    ssize_t asnLength;
    unsigned char *asndata;
} FormatedMessage_t;

typedef struct ReportingMessages {
    FormatedMessage_t message {};
    ConnectedCU_t *peerInfo = nullptr;
    long outLen = 0;
    unsigned char base64Data[RECEIVE_SCTP_BUFFER_SIZE * 2] {};
    char buffer[RECEIVE_SCTP_BUFFER_SIZE * 8] {};
    StatCollector *statCollector = nullptr;
} ReportingMessages_t;

cxxopts::ParseResult parse(int argc, const char *argv[], sctp_params_t &pSctpParams);

int buildInotify(sctp_params_t &sctpParams);

void handleTermInit(sctp_params_t &sctpParams);

void handleConfigChange(sctp_params_t *sctpParams);

void listener(sctp_params_t *params);

void sendTermInit(sctp_params_t &sctpParams);

int setSocketNoBlocking(int socket);

void handleEinprogressMessages(struct epoll_event &event,
                               ReportingMessages_t &message,
                               RmrMessagesBuffer_t &rmrMessageBuffer,
                               sctp_params_t *params);

void handlepoll_error(struct epoll_event &event,
                      ReportingMessages_t &message,
                      RmrMessagesBuffer_t &rmrMessageBuffer,
                      sctp_params_t *params);


void cleanHashEntry(ConnectedCU_t *peerInfo, Sctp_Map_t *m);


/**
 *
 * @param message
 * @param rmrMessageBuffer
 */
void getRequestMetaData(ReportingMessages_t &message, RmrMessagesBuffer_t &rmrMessageBuffer);

/**
 *
 * @param sctpMap
 * @param messagBuffer
 * @param message
 * @param failedMesgId
 * @return
 */
int sendMessagetoCu(Sctp_Map_t *sctpMap,
                    RmrMessagesBuffer_t &messagBuffer,
                    ReportingMessages_t &message,
                    int failedMesgId);

void sendFailedSendingMessagetoXapp(RmrMessagesBuffer_t &rmrMessageBuffer,
                                    ReportingMessages_t &message,
                                    int failedMesgId);

int sendRequestToXapp(ReportingMessages_t &message,
                      int requestId,
                      RmrMessagesBuffer_t &rmrMmessageBuffer);

/**
 *
 * @param message
 * @param msgType
 * @param requestType
 * @param rmrMessageBuffer
 * @param sctpMap
 * @return
 */
/*
int sendResponseToXapp(ReportingMessages_t &message,
                       int msgType,
                       int requestType,
                       RmrMessagesBuffer_t &rmrMessageBuffer,
                       Sctp_Map_t *sctpMap);
*/

/**
 *
 * @param peerInfo
 * @param message
 * @param m
 * @return
 */
int sendSctpMsg(ConnectedCU_t *peerInfo,
                ReportingMessages_t &message,
                Sctp_Map_t *m);

/**
 *
 * @param events
 * @param sctpMap
 * @param numOfMessages
 * @param rmrMessageBuffer
 * @param ts
 * @return
 */
int receiveDataFromSctp(struct epoll_event *events,
                        Sctp_Map_t *sctpMap,
                        int &numOfMessages,
                        RmrMessagesBuffer_t &rmrMessageBuffer,
                        struct timespec &ts);

/**
 *
 * @param rmrAddress
 * @return
 */
void getRmrContext(sctp_params_t &pSctpParams);

/**
 *
 * @param epoll_fd
 * @param rmrCtx
 * @param sctpMap
 * @param messagBuffer
 * @return
 */
int receiveXappMessages(Sctp_Map_t *sctpMap,
                        RmrMessagesBuffer_t &rmrMessageBuffer,
                        struct timespec &ts);

/**
 *
 * @param messagBuffer
 * @param failedMsgId
 * @param sctpMap
 * @return
 */
int sendDirectionalSctpMsg(RmrMessagesBuffer_t &messagBuffer,
                           ReportingMessages_t &message,
                           int failedMsgId,
                           Sctp_Map_t *sctpMap);
/**
 *
 * @param pdu
 * @param message
 * @param rmrMessageBuffer
 */
void asnInitiatingRequest(E2AP_PDU_t *pdu,
                          Sctp_Map_t *sctpMap,
                          ReportingMessages_t &message,
                          RmrMessagesBuffer_t &rmrMessageBuffer);
/**
 *
 * @param pdu
 * @param message
 * @param sctpMap
 * @param rmrMessageBuffer
 */
void asnSuccsesfulMsg(E2AP_PDU_t *pdu,
                      Sctp_Map_t *sctpMap,
                      ReportingMessages_t &message,
                      RmrMessagesBuffer_t &rmrMessageBuffer);
/**
 *
 * @param pdu
 * @param message
 * @param sctpMap
 * @param rmrMessageBuffer
 */
void asnUnSuccsesfulMsg(E2AP_PDU_t *pdu,
                        Sctp_Map_t *sctpMap,
                        ReportingMessages_t &message,
                        RmrMessagesBuffer_t &rmrMessageBuffer);

/**
 *
 * @param rmrMessageBuffer
 * @param message
 * @return
 */
int sendRmrMessage(RmrMessagesBuffer_t &rmrMessageBuffer, ReportingMessages_t &message);
/**
 *
 * @param epoll_fd
 * @param peerInfo
 * @param events
 * @param sctpMap
 * @param enodbName
 * @param msgType
 * @returnsrc::logger_mt& lg = my_logger::get();
 */
int addToEpoll(int epoll_fd, ConnectedCU_t *peerInfo, uint32_t events, Sctp_Map_t *sctpMap, char *enodbName, int msgType);
/**
 *
 * @param epoll_fd
 * @param peerInfo
 * @param events
 * @param sctpMap
 * @param enodbName
 * @param msgType
 * @return
 */
int modifyToEpoll(int epoll_fd, ConnectedCU_t *peerInfo, uint32_t events, Sctp_Map_t *sctpMap, char *enodbName, int msgType);

/**
 *
 * @param message
 */
void buildJsonMessage(ReportingMessages_t &message);

/**
 *
 *
 * @param state
 * @return
 */
string translateRmrErrorMessages(int state);


static inline uint64_t rdtscp(uint32_t &aux) {
    uint64_t rax,rdx;
    asm volatile ("rdtscp\n" : "=a" (rax), "=d" (rdx), "=c" (aux) : :);
    return (rdx << (unsigned)32) + rax;
}
#ifndef RIC_SCTP_CONNECTION_FAILURE
#define RIC_SCTP_CONNECTION_FAILURE  10080
#endif

#endif //X2_SCTP_THREAD_H
