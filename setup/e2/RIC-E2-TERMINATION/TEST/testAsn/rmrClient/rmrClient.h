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
// Created by adi ENZEL on 2/12/20.
//

#include <iostream>
#include <string>
#include <cstring>
#include <exception>
#include <utility>

#include <unistd.h>

//rmr testing thread
#include <rmr/rmr.h>


using namespace std;
#define MAX_RECEIVED_BUFFER 8192

class RmrException: public std::exception
{
    std::string msg;
public:
    explicit RmrException(std::string  msg) : msg(std::move(msg)){}

    const char* what() const  noexcept override {
        return msg.c_str();
    }
};



class RmrClient {
private:
    void *rmrCtx = nullptr;
    int rmr_fd = 0;

    void getRmrContext(const char *address, int epoll_fd) {
        rmrCtx = rmr_init((char *)address, MAX_RECEIVED_BUFFER, RMRFL_NONE);
        if (rmrCtx == nullptr) {
            cerr << "Failed to initialize RMR. address = " << address << endl;
            throw RmrException("Failed to initialize RMR");
        }

        rmr_set_stimeout(rmrCtx, 0);    // disable retries for any send operation
        cout << "Wait for RMR_Ready" << endl;
        auto rmrReady = 0;
        auto count = 0;
        while (!rmrReady) {
            if ((rmrReady = rmr_ready(rmrCtx)) == 0) {
                usleep(1000000);
            }
            count++;
            if (count % 60 == 0) {
                cout << "Wait for RMR ready state : " << count << " seconds" << endl;
            }
        }
        cout << "RMR running" << endl;


        rmr_init_trace(rmrCtx, 200);
        // get the RMR fd for the epoll
        rmr_fd = rmr_get_rcvfd(rmrCtx);
        struct epoll_event event{};
        // add RMR fd to epoll
        event.events = (EPOLLIN);
        event.data.fd = rmr_fd;
        // add listening RMR FD to epoll
        if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, rmr_fd, &event)) {
            cerr << "Failed to add RMR descriptor to epoll. " << strerror(errno) << endl;
            close(rmr_fd);
            rmr_close(rmrCtx);
            throw RmrException("Failed to add RMR descriptor to epoll.");
        }
    }

public:
    RmrClient(const char *address, int epoll_fd)  {
        try {
            getRmrContext(address, epoll_fd);
        } catch (RmrException &e) {
            cout << e.what() << endl;
            exit(-1);
        }
    }

    RmrClient() = delete;
    RmrClient(const RmrClient &) = delete;
    RmrClient &operator=(const RmrClient &) = delete;
    RmrClient &operator&&(const RmrClient &) = delete;


    inline void *getRmrCtx() const {
        return rmrCtx;
    }

    int getRmrFd() const {
        return rmr_fd;
    }

    inline rmr_mbuf_t *allocateRmrMsg(int size) {return (rmr_alloc_msg(rmrCtx, size));}

    static inline void freeRmrMsg(rmr_mbuf_t *msg) {rmr_free_msg(msg);}
};