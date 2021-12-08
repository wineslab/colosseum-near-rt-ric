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
// Created by adi ENZEL on 11/27/19.
//

#include "base64.h"
#include <mdclog/mdclog.h>
#include <cgreen/cgreen.h>
#include <cstdio>
#include <cerrno>
#include <cstdlib>
#include <cstring>

using namespace std;

Describe(base64);
BeforeEach(base64) {}
AfterEach(base64) {}

using namespace cgreen;

void init_log() {
    mdclog_attr_t *attr;
    mdclog_attr_init(&attr);
    mdclog_attr_set_ident(attr, "TestConfiguration");
    mdclog_init(attr);
    mdclog_attr_destroy(attr);
}

const char *data = "ABC123Test Lets Try this' input and see What \"happens\"";

Ensure(base64, encDec) {
    string str = "ABC123Test Lets Try this' input and see What \"happens\"";
    auto *buf = (unsigned char *)malloc(str.length() * 2);
    auto length = (long)(str.length() * 2);
    base64::encode((unsigned char *)str.c_str(), str.length(), buf, length);
    auto *backBackBuff = (unsigned char *)malloc(length);
    auto length2 = length;
    assert_that(base64::decode(buf, length, backBackBuff, length2) == 0);
    std::string str1( backBackBuff, backBackBuff + sizeof backBackBuff / sizeof backBackBuff[0]);

    assert_that(str.length() == (ulong)length2)
    //auto val = str.compare((const char *)backBackBuff);
    assert_that(str.compare((const char *)backBackBuff) == 0)
    free(backBackBuff);
    free(buf);
}

Ensure(base64, errorsHandling) {
    string str = "ABC123Test Lets Try this' input and see What \"happens\"";
    auto *buf = (unsigned char *)malloc(str.length() * 2);
    auto length = (long)(str.length());
    assert_that(base64::encode((unsigned char *)str.c_str(), str.length(), buf, length) == -1);
    length = (long)(str.length() * 2);
    assert_that(base64::encode((unsigned char *)str.c_str(), str.length(), buf, length) == 0);
    auto *backBackBuff = (unsigned char *)malloc(length);
    auto length2 = length >> 2;
    assert_that(base64::decode(buf, length, backBackBuff, length2) == -1);
    //std::string str1( backBackBuff, backBackBuff + sizeof backBackBuff / sizeof backBackBuff[0]);
    auto length1 = 0l;
    assert_that(base64::encode((unsigned char *)str.c_str(), str.length(), nullptr , length) == -1);
//    assert_that(base64::encode((unsigned char *)str.c_str(), str.length(), nullptr , length) == -1);
    assert_that(base64::encode(nullptr, str.length(), backBackBuff , length) == -1);
    assert_that(base64::encode((unsigned char *)str.c_str(), length1, backBackBuff , length) == -1);
    assert_that(base64::encode(nullptr, str.length(), backBackBuff , length1) == -1);
    length1 = -1;
    assert_that(base64::encode((unsigned char *)str.c_str(), length1, backBackBuff , length) == -1);
    assert_that(base64::encode(nullptr, str.length(), backBackBuff , length1) == -1);

}


int main(const int argc, char **argv) {
    mdclog_severity_t loglevel = MDCLOG_INFO;
    init_log();
    mdclog_level_set(loglevel);

    //TestSuite *suite = create_test_suite();
    TestSuite *suite = create_named_test_suite_(__FUNCTION__, __FILE__, __LINE__);

    add_test_with_context(suite, base64, encDec);
    add_test_with_context(suite, base64, errorsHandling);

    return cgreen::run_test_suite(suite, create_text_reporter());

}