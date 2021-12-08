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
// Created by adi ENZEL on 2/18/20.
//

#ifndef E2_HTTPSERVER_H
#define E2_HTTPSERVER_H

#include <pistache/http.h>
#include <pistache/router.h>
#include <pistache/endpoint.h>
#include <pistache/net.h>


using namespace std;
using namespace Pistache;

#define SA struct sockaddr

class HttpServer {
public:
    explicit HttpServer(Address addr);

    void init(size_t thr = 2);

    void start();

private:

    long transactionCounter = 0;
    int httpBaseSocket;
    void setupRoutes();


    void sendSetupReq(const Rest::Request& request, Http::ResponseWriter response);

    std::shared_ptr<Http::Endpoint> httpEndpoint;
    Rest::Router router;
};

#endif //E2_HTTPSERVER_H
