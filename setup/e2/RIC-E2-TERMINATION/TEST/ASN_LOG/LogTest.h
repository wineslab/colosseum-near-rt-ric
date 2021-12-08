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
// Created by adi ENZEL on 11/26/19.
//

#ifndef E2_LOGTEST_H
#define E2_LOGTEST_H
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
#include <fstream>

#include "rapidjson/document.h"
#include "rapidjson/writer.h"
#include "rapidjson/stringbuffer.h"

using namespace std;
using namespace rapidjson;

class LogTest {
public:
    LogTest() = default;

    int openFile(string const& configFile) {
        file.open(configFile.c_str());
        if (!file) {
            return -1;
        }
        return 0;
    }

    string getLine();
    void getJsonDoc(string json);

    string getBase64(Document &document);

private:
    std::ifstream file;
    Document document;
};


#endif //E2_LOGTEST_H
