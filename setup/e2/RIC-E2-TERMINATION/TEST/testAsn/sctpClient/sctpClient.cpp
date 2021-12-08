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
 * Unless required by applicable law fprintfor agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

//
// Created by adi ENZEL on 2/10/20.
//

#include "sctpClient.h"

#define READ_BUFFER_SIZE 64 * 1024


using namespace std;


void createHttpLocalSocket(SctpClient_t *sctpClient) {
    struct sockaddr_in address{};
    int addrlen = sizeof(address);
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(9098);
    sctpClient->httpSocket = accept(sctpClient->httpBaseSocket, (struct sockaddr *) &address, (socklen_t *) &addrlen) < 0;
    if (sctpClient->httpSocket) {
        fprintf(stderr, "Accept() error. %s\n", strerror(errno));
        exit(-1);
    }
    struct epoll_event event{};
    event.data.fd = sctpClient->httpSocket;
    event.events = (EPOLLIN | EPOLLET);
    if (epoll_ctl(sctpClient->epoll_fd, EPOLL_CTL_ADD, sctpClient->httpSocket, &event) < 0) {
        fprintf(stderr, "epoll_ctl EPOLL_CTL_ADD, %s\n", strerror(errno));
        close(sctpClient->httpSocket);
        exit(-1);
    }
}

int createEpoll(SctpClient &sctpClient) {
    sctpClient.epoll_fd = epoll_create1(0);
    if (sctpClient.epoll_fd == -1) {
        fprintf(stderr, "failed to open epoll descriptor. %s\n", strerror(errno));
        return -1;
    }
    return sctpClient.epoll_fd;
}



int createSctpConnction(SctpClient *sctpClient, const char *address, int port, bool local) {
    sctpClient->sctpSock = socket(AF_INET6, SOCK_STREAM, IPPROTO_SCTP);
    if (sctpClient->sctpSock < 0) {
        fprintf(stderr, "Socket Error, %s %s, %d\n", strerror(errno), __func__, __LINE__);
        return -1;
    }
    auto optval = 1;
    if (setsockopt(sctpClient->sctpSock, SOL_SOCKET, SO_REUSEPORT, &optval, sizeof optval) != 0) {
        fprintf(stderr, "setsockopt SO_REUSEPORT Error, %s %s, %d\n", strerror(errno), __func__, __LINE__);
        close(sctpClient->sctpSock);
        return -1;
    }
    optval = 1;
    if (setsockopt(sctpClient->sctpSock, SOL_SOCKET, SO_REUSEADDR, &optval, sizeof optval) != 0) {
        fprintf(stderr, "setsockopt SO_REUSEADDR Error, %s %s, %d\n", strerror(errno), __func__, __LINE__);
        close(sctpClient->sctpSock);
        return -1;
    }
    struct sockaddr_in6 servaddr = {};
//    struct addrinfo hints = {};
//    struct addrinfo *result;

    servaddr.sin6_family = AF_INET6;
    servaddr.sin6_port = htons(port);      /* daytime server */
    inet_pton(AF_INET6, address, &servaddr.sin6_addr);
    // the bind here is to maintain the client port this is only if the test is not on the same IP as the tested system
    if (!local) {
        struct sockaddr_in6 localAddr{};
        localAddr.sin6_family = AF_INET6;
        localAddr.sin6_addr = in6addr_any;
        localAddr.sin6_port = htons(port);
        if (bind(sctpClient->sctpSock, (struct sockaddr *) &localAddr, sizeof(struct sockaddr_in6)) < 0) {
            fprintf(stderr, "bind Socket Error, %s %s, %d\n", strerror(errno), __func__, __LINE__);
            return -1;
        }//Ends the binding.
    }

    // Add to Epol
    struct epoll_event event{};
    event.data.fd = sctpClient->sctpSock;
    event.events = (EPOLLOUT | EPOLLIN | EPOLLET);
    if (epoll_ctl(sctpClient->epoll_fd, EPOLL_CTL_ADD, sctpClient->sctpSock, &event) < 0) {
        fprintf(stderr, "epoll_ctl EPOLL_CTL_ADD, %s\n", strerror(errno));
        close(sctpClient->sctpSock);
        return -1;
    }

    char hostBuff[NI_MAXHOST];
    char portBuff[NI_MAXHOST];

    if (getnameinfo((SA *) &servaddr, sizeof(servaddr),
                    hostBuff, sizeof(hostBuff),
                    portBuff, sizeof(portBuff),
                    (uint) (NI_NUMERICHOST) | (uint) (NI_NUMERICSERV)) != 0) {
        fprintf(stderr, "getnameinfo() Error, %s  %s %d\n", strerror(errno), __func__, __LINE__);
        return -1;
    }

    auto flags = fcntl(sctpClient->sctpSock, F_GETFL, 0);
    if (flags == -1) {
        fprintf(stderr, "fcntl error. %s\n", strerror(errno));
        close(sctpClient->sctpSock);
        return -1;
    }

    flags = (unsigned) flags | (unsigned) O_NONBLOCK;
    if (fcntl(sctpClient->sctpSock, F_SETFL, flags) == -1) {
        fprintf(stderr, "fcntl set O_NONBLOCK fail. %s\n", strerror(errno));
        close(sctpClient->sctpSock);
        return -1;
    }

    if (connect(sctpClient->sctpSock, (SA *) &servaddr, sizeof(servaddr)) < 0) {
        if (errno != EINPROGRESS) {
            fprintf(stderr, "connect FD %d to host : %s port %d, %s\n", sctpClient->sctpSock, address, port,
                    strerror(errno));
            close(sctpClient->sctpSock);
            return -1;
        }
        fprintf(stdout, "Connect to FD %d returned with EINPROGRESS : %s\n", sctpClient->sctpSock, strerror(errno));
    }
    return sctpClient->sctpSock;
}

__attribute_warn_unused_result__ int createListeningTcpConnection(SctpClient *sctpClient) {
    if ((sctpClient->httpBaseSocket = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
        fprintf(stderr, "socket failed. %s", strerror(errno));
        return -1;
    }

    struct sockaddr_in address{};
    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(9098);

    if (bind(sctpClient->httpBaseSocket, (struct sockaddr *)&address, sizeof(address)) < 0) {
        fprintf(stderr, "Bind failed , %s %s, %d\n", strerror(errno), __func__, __LINE__);
        return -1;
    }

    struct epoll_event event{};
    event.data.fd = sctpClient->httpBaseSocket;
    event.events = (EPOLLIN | EPOLLET);
    if (epoll_ctl(sctpClient->epoll_fd, EPOLL_CTL_ADD, sctpClient->httpBaseSocket, &event) < 0) {
        fprintf(stderr, "epoll_ctl EPOLL_CTL_ADD, %s\n", strerror(errno));
        close(sctpClient->httpBaseSocket);
        return -1;
    }
    if (listen(sctpClient->httpBaseSocket, 128) < 0)
    {
        fprintf(stderr,"listen() error. %s", strerror(errno));
        return -1;
    }
    return 0;
}

__attribute_warn_unused_result__ int modifyEpollToRead(SctpClient *sctpClient, int modifiedSocket) {
    struct epoll_event event{};
    event.data.fd = modifiedSocket;
    event.events = (EPOLLIN | EPOLLET);
    if (epoll_ctl(sctpClient->epoll_fd, EPOLL_CTL_MOD, modifiedSocket, &event) < 0) {
        fprintf(stderr, "failed to open epoll descriptor. %s\n", strerror(errno));
        return -1;
    }
    return 0;
}


__attribute_warn_unused_result__ cxxopts::ParseResult parse(SctpClient &sctpClient, int argc, const char *argv[]) {
    cxxopts::Options options(argv[0], "sctp client test application");
    options.positional_help("[optional args]").show_positional_help();
    options.allow_unrecognised_options().add_options()
            ("a,host", "Host address", cxxopts::value<std::string>(sctpClient.host)->default_value("127.0.0.1"))
            ("p,port", "port number", cxxopts::value<int>(sctpClient.rmrPort)->default_value("38200"))
            ("h,help", "Print help");

    auto result = options.parse(argc, argv);

    if (result.count("help")) {
        std::cout << options.help({""}) << std::endl;
        exit(0);
    }
    return result;
}

void run(SctpClient_t *sctpClient) {
    cout << "in theread" << endl;
    sleep(10);
    cout << "in theread after sleep" << endl;
    sleep(10);
}

void runFunc(SctpClient_t *sctpClient) {
    cout << "in theread 1" << endl;

    char rmrAddress[128] {};
    cout << "in theread 2" << endl;

    snprintf(rmrAddress, 128, "%d", sctpClient->rmrPort);
    cout << "in theread 3" << endl;

    RmrClient rmrClient = {rmrAddress, sctpClient->epoll_fd};
    cout << "in theread 4" << endl;


    auto *events = (struct epoll_event *) calloc(MAXEVENTS, sizeof(struct epoll_event));
//    auto counter = 1000;
    //    uint64_t st = 0;
//    uint32_t aux1 = 0;
//    st = rdtscp(aux1);
    E2AP_PDU_t *pdu = nullptr;

    auto *msg = rmrClient.allocateRmrMsg(8192);
    while (true) {
        auto numOfEvents = epoll_wait(sctpClient->epoll_fd, events, MAXEVENTS, 1000);
        if (numOfEvents < 0) {
            if (errno == EINTR) {
                fprintf(stderr, "got EINTR : %s\n", strerror(errno));
                continue;
            }
            fprintf(stderr, "Epoll wait failed, errno = %s\n", strerror(errno));
            break;
        }
        if (numOfEvents == 0) { // timeout
//            if (--counter <= 0) {
//                fprintf(stdout, "Finish waiting for epoll. going out of the thread\n");
//                continue;
//            }
        }

        auto done = 0;

        for (auto i = 0; i < numOfEvents; i++) {
            uint32_t aux1 = 0;
            auto start = rdtscp(aux1);

            if ((events[i].events & EPOLLERR) || (events[i].events & EPOLLHUP)) {
                fprintf(stderr, "Got EPOLLERR or EPOLLHUP on fd = %d, errno = %s\n", events[i].data.fd,
                        strerror(errno));
                close(events[i].data.fd);
            } else if (events[i].events & EPOLLOUT) {  // AFTER EINPROGRESS
                if (modifyEpollToRead(sctpClient, events[i].data.fd) < 0) {
                    fprintf(stderr, "failed modify FD %d after got EINPROGRESS\n", events[i].data.fd);
                    close(events[i].data.fd);
                    continue;
                }
                fprintf(stdout, "Connected to server after EinProgreevents[i].data.fdss FD %d\n", events[i].data.fd);
                //TODO need to define RmrClient class
            } else if (events[i].data.fd == sctpClient->httpBaseSocket) {
                createHttpLocalSocket(sctpClient);
            } else if (events[i].data.fd == sctpClient->httpSocket) {
                //TODO handle messages from the http server
                char buffer[READ_BUFFER_SIZE] {};
                while (true) {
                    auto size = read(sctpClient->httpSocket, buffer, READ_BUFFER_SIZE);
                    if (size < 0) {
                        if (errno == EINTR) {
                            continue;
                        }
                        if (errno != EAGAIN) {
                            fprintf(stderr, "Read error, %s\n", strerror(errno));
                            done = 1;
                        }
                        break; // EAGAIN exit from loop on read normal way or on read error
                    }
                    // we got message get the id of message
                    char *tmp;

                    // get mesage type
                    char *val = strtok_r(buffer, sctpClient->delimiter, &tmp);
                    messages_t messageType;
                    char *dummy;
                    if (val != nullptr) {
                        messageType = (decltype(messageType))strtol(val, &dummy, 10);
                    } else {
                        fprintf(stderr,"wrong message %s", buffer);
                        break;
                    }
                    char sctpLinkId[128] {};
                    val = strtok_r(nullptr, sctpClient->delimiter, &tmp);
                    if (val != nullptr) {
                        memcpy(sctpLinkId, val, tmp - val);
                    } else {
                        fprintf(stderr,"wriong id %s", buffer);
                        break;
                    }

                    char *values[128] {};
                    int index = 0;
                    while ((val = strtok_r(nullptr, sctpClient->delimiter, &tmp)) != nullptr) {
                        auto valueLen = tmp - val;
                        values[index] = (char *)calloc(1, valueLen);
                        memcpy(values[index], val, valueLen);
                        index++;
                    }
                    values[i] = (char *)calloc(1, strlen(tmp));

                    switch ((int)messageType) {
                        case setupRequest_gnb:
                        case setupRequest_en_gNB:
                        case setupRequest_ng_eNB:
                        case setupRequest_eNB: {

                            char *ricAddress = nullptr;
                            if (values[0] != nullptr) {
                               ricAddress = values[0];
                            } else {
                                fprintf(stderr,"wrong address %s", buffer);
                                break;
                            }
                            //ric port
                            int ricPort = 0;

                            if (values[1] != nullptr) {
                                ricPort = (decltype(ricPort))strtol(values[1], &dummy, 10);
                            } else {
                                fprintf(stderr,"wrong port %s", buffer);
                                for (auto e : values) {
                                    if (e != nullptr) {
                                        free(e);
                                    }
                                }
                                break;
                            }

                            // need to send message to E2Term
                            // build connection
                            auto fd = createSctpConnction(sctpClient, (const char *)ricAddress, ricPort);
                            if (fd < 0) {
                                fprintf(stderr,"Failed to create connection to %s:%d\n", ricAddress, ricPort);
                                for (auto e : values) {
                                    if (e != nullptr) {
                                        free(e);
                                    }
                                }
                                break;
                            }

                            auto len = strlen(values[index]);
                            auto *b64Decoded = (unsigned char *)calloc(1, len);
                            base64::decode((const unsigned char *)values[index], len, b64Decoded, (long)len);

                            for (auto e : values) {
                                if (e != nullptr) {
                                    free(e);
                                }
                            }
                            // send data
                            while (true) {
                                if (send(fd, b64Decoded, len, MSG_NOSIGNAL) < 0) {
                                    if (errno == EINTR) {
                                        continue;
                                    }
                                    cerr << "Error sendingdata to e2Term. " << strerror(errno) << endl;
                                    break;
                                }
                                cout << "Message sent" << endl;
                                break;
                            }
                            free(b64Decoded);
                            char key[128] {};
                            char *value = (char *)calloc(1,256);
                            snprintf(key, 128, "id:%s", sctpLinkId);
                            snprintf(value, 16, "%d", fd);
                            sctpClient->mapKey.setkey(key, (void *)value);

                            snprintf(key, 128, "fd:%d", fd);
                            snprintf(&value[128], 128, "%s", sctpLinkId);
                            sctpClient->mapKey.setkey(key, (void *)&value[128]);
                            break;
                        }
                        case nothing:
                        default: {
                            break;
                        }
                    }
                 }
            } else if (events[i].data.fd == rmrClient.getRmrFd()) {
                msg->state = 0;
                msg = rmr_rcv_msg(rmrClient.getRmrCtx(), msg);
                if (msg == nullptr) {
                    cerr << "rmr_rcv_msg return with null pointer" << endl;
                    exit(-1);
                } else if (msg->state != 0) {
                    cerr << "rmr_rcv_msg return with error status number : " << msg->state << endl;
                    msg->state = 0;
                    continue;
                }
         sleep(100);       cout << "Got RMR message number : " << msg->mtype << endl;

            } else { // got data from server
                /* We RMR_ERR_RETRY have data on the fd waiting to be read. Read and display it.
                * We must read whatever data is available completely, as we are running
                *  in edge-triggered mode and won't get a notification again for the same data. */
                //TODO build a callback function to support many tests
                if (pdu != nullptr) {
                    ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, pdu);
                }
                unsigned char buffer[SCTP_BUFFER_SIZE]{};
                while (true) {
                    auto len = read(events[i].data.fd, buffer, SCTP_BUFFER_SIZE);
                    if (len < 0) {
                        if (errno == EINTR) {
                            continue;
                        }
                        /* If errno == EAGAIN, that means we have read all
                        data. So go back to the main loop. */
                        if (errno != EAGAIN) {
                            fprintf(stderr, "Read error, %s\n", strerror(errno));
                            done = 1;
                        }
                        break; // EAGAIN exit from loop on read normal way or on read error
                    } else if (len == 0) {
                        /* End of file. The remote has closed the connection. */
                        fprintf(stdout, "EOF Closed connection - descriptor = %d", events[i].data.fd);
                        done = 1;
                        break;
                    }

                    asn_dec_rval_t rval;
                    rval = asn_decode(nullptr,
                                      ATS_ALIGNED_BASIC_PER,
                                      &asn_DEF_E2AP_PDU,
                                      (void **) &pdu,
                                      buffer,
                                      len);
                    if (rval.code != RC_OK) {
                        fprintf(stderr, "Error %d Decoding E2AP PDU from E2TERM\n", rval.code);
                        break;
                    }
                    //TODO handle messages
                    //                    switch (pdu->present) {
                    //                        case E2AP_PDU_PR_initiatingMessage: {//initiating message
                    //                            asnInitiatingRequest(pdu, message, rmrMessageBuffer);
                    //                            break;
                    //                        }
                    //                        case E2AP_PDU_PR_successfulOutcome: { //successful outcome
                    //                            asnSuccsesfulMsg(pdu, message, sctpMap, rmrMessageBuffer);
                    //                            break;
                    //                        }
                    //                        case E2AP_PDU_PR_unsuccessfulOutcome: { //Unsuccessful Outcome
                    //                            asnUnSuccsesfulMsg(pdu, message, sctpMap, rmrMessageBuffer);
                    //                            break;
                    //            ipv6 client server c program            }
                    //                        case E2AP_PDU_PR_NOTHING:
                    //                        default:
                    //                            fprintf(stderr, "Unknown index %d in E2AP PDU\n", pdu->present);
                    //                            break;
                    //                    }
                }
            }
            aux1 = 0;
            fprintf(stdout, "one loop took  %ld clocks\n", rdtscp(aux1) - start);
        }
        if (done) {
            //TODO report to RMR on closed connection
        }
    }
    return;// nullptr;
}

auto main(const int argc, const char **argv) -> int {
    SctpClient_t sctpClient;
    //unsigned num_cpus = std::thread::hardware_concurrency();

    auto result = parse(sctpClient, argc, argv);
    auto epoll_fd = createEpoll(sctpClient);
    if (epoll_fd <= 0) {
        exit(-1);
    }
    if (createListeningTcpConnection(&sctpClient) < 0) {
        exit(-1);
    }
    std::thread th(runFunc, &sctpClient);
//    std::thread th(run, &sctpClient);


    sleep(29);
    //start the http server
    Port port(9080);

    Address addr(Ipv4::any(), port);
    HttpServer server(addr);

    server.init(1);


    server.start();

    th.join();

}
