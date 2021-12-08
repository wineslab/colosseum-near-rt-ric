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
// Created by adi ENZEL on 3/24/20.
//

#ifndef E2_STATCOLLECTOR_H
#define E2_STATCOLLECTOR_H

#include <unordered_map>
#include <mutex>
#include <shared_mutex>
#include <thread>
#include <string>
#include <iostream>
#include <utility>
#include <chrono>
#include <ctime>
#include <iomanip>
#include <mdclog/mdclog.h>
//#include <tbb/concurrent_unordered_map.h>

//using namespace tbb;

typedef struct statResult {
    std::string ranName;
    uint32_t receivedMessages;
    uint32_t sentMessages;
} statResult_t ;

class StatCollector {

    static std::mutex singltonMutex;
    static std::atomic<StatCollector *> obj;

public:
    static StatCollector* GetInstance() {
        StatCollector* pStatCollector = obj.load(std::memory_order_acquire);
        if (pStatCollector == nullptr) {
            std::lock_guard<std::mutex> lock(singltonMutex);
            pStatCollector = obj.load(std::memory_order_relaxed);
            if (pStatCollector == nullptr) {
                pStatCollector = new StatCollector();
                obj.store(pStatCollector, std::memory_order_release);
            }
        }
        return pStatCollector;
    }

    void incSentMessage(const std::string &key) {
        increment(sentMessages, key);
    }
    void incRecvMessage(const std::string &key) {
        increment(recvMessages, key);
    }

    std::vector<statResult_t> &getCurrentStats() {
        results.clear();

        for (auto const &e : recvMessages) {
            statResult_t result {};
            result.ranName = e.first;
            result.receivedMessages = e.second.load(std::memory_order_acquire);
            auto found = sentMessages.find(result.ranName);
            if (found != sentMessages.end()) {
                result.sentMessages = found->second.load(std::memory_order_acquire);
            } else {
              result.sentMessages = 0;
            }

            results.emplace_back(result);
        }
        return results;
    }

    StatCollector(const StatCollector&)= delete;
    StatCollector& operator=(const StatCollector&)= delete;

private:
    //tbb::concurrent_unordered_map<std::string, int> sentMessages;
    std::unordered_map<std::string, std::atomic<int>> sentMessages;
    std::unordered_map<std::string, std::atomic<int>> recvMessages;
//    tbb::concurrent_unordered_map<std::string, int> recvMessages;
    std::vector<statResult_t> results;


//    StatCollector() = default;
    StatCollector() {
        sentMessages.clear();
        recvMessages.clear();
    }
    ~StatCollector() = default;


    void increment(std::unordered_map<std::string, std::atomic<int>> &map, const std::string &key);

};

void StatCollector::increment(std::unordered_map<std::string, std::atomic<int>> &map, const std::string &key) {
    if (map.empty()) {
        map.emplace(std::piecewise_construct,
                    std::forward_as_tuple(key),
                    std::forward_as_tuple(1));
        return;
    }
    auto found = map.find(key);
    if (found != map.end()) { //inc
        map[key].fetch_add(1, std::memory_order_release);
        //map[key]++;
    } else { //add
        //sentMessages.emplace(std::make_pair(std::string(key), std::atomic<int>(0)));
        map.emplace(std::piecewise_construct,
                    std::forward_as_tuple(key),
                    std::forward_as_tuple(1));
    }

}


// must define this to allow StatCollector private variables to be known to compiler linker
std::mutex StatCollector::singltonMutex;
std::atomic<StatCollector *> StatCollector::obj;


void statColectorThread(void *runtime) {
    bool *stop_loop = (bool *)runtime;
    auto *statCollector = StatCollector::GetInstance();
    std::time_t tt = std::chrono::system_clock::to_time_t (std::chrono::system_clock::now());

    struct std::tm * ptm = std::localtime(&tt);
    std::cout << "Waiting for the next minute to begin...\n";
    ptm->tm_min = ptm->tm_min + (5 - ptm->tm_min % 5);
    ptm->tm_sec=0;

    std::this_thread::sleep_until(std::chrono::system_clock::from_time_t(mktime(ptm)));

// alligned to 5 minutes
    while (true) {
        if (*stop_loop) {
            break;
        }
        for (auto const &e : statCollector->getCurrentStats()) {
            if (mdclog_level_get() >= MDCLOG_INFO) {
                mdclog_write(MDCLOG_INFO, "RAN : %s sent messages : %d recived messages : %d\n",
                             e.ranName.c_str(), e.sentMessages, e.receivedMessages);
            }
        }
        std::this_thread::sleep_for(std::chrono::seconds(300));
    }
}
#endif //E2_STATCOLLECTOR_H
