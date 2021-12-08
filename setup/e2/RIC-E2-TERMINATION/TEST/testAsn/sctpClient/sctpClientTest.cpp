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
// Created by adi ENZEL on 2/10/20.
//
#include <thread>
#include <vector>
#include <cgreen/cgreen.h>

//#include "sctpClient.cpp"

int epoll_fd = 0;

Describe(Cgreen);
BeforeEach(Cgreen) {}
AfterEach(Cgreen) {
    close(epoll_fd);
}

using namespace cgreen;
using namespace std;


Ensure(Cgreen, createEpoll) {
    epoll_fd = createEpoll();
    assert_that(epoll_fd != -1);
    assert_that(epoll_fd > 0);
}



Ensure(Cgreen, createConnectionIpV6) {
    auto epoll_fd = createEpoll();
    assert_that(epoll_fd != -1);
    assert_that(epoll_fd > 0);
    unsigned num_cpus = std::thread::hardware_concurrency();
    std::vector<std::thread> threads(num_cpus);
//    int i = 0;

//    threads[i] = std::thread(listener, &epoll_fd);
//    auto port = 36422;
//    auto fd = createSctpConnction("::1", port, epoll_fd);
//    assert_that(fd != -1);
//    assert_that(fd == 0);
//    threads[i].join();
//
//    close(fd);
}

Ensure(Cgreen, createConnectionIpV4) {
    auto epoll_fd = createEpoll();
    assert_that(epoll_fd != -1);
    assert_that(epoll_fd > 0);
    unsigned num_cpus = std::thread::hardware_concurrency();
    std::vector<std::thread> threads(num_cpus);
//    int i = 0;

//    threads[i] = std::thread(listener, &epoll_fd);
//    auto port = 36422;
//    auto fd = createSctpConnction("127.0.0.1", port, epoll_fd);
//    assert_that(fd != -1);
//    assert_that(fd == 0);
//
//    threads[i].join();
//    close(fd);
}


//int main(const int argc, char **argv) {
//    TestSuite *suite = create_named_test_suite_(__FUNCTION__, __FILE__, __LINE__);
//
//    add_test_with_context(suite, Cgreen, createEpoll);
//    //add_test_with_context(suite, Cgreen, createConnectionIpV6);
//    add_test_with_context(suite, Cgreen, createConnectionIpV4);
//    return cgreen::run_test_suite(suite, create_text_reporter());
//}