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
// Created by adi ENZEL on 11/19/19.
//

#include "ReadConfigFile.h"
#include <mdclog/mdclog.h>
#include <cgreen/cgreen.h>

Describe(Cgreen);
BeforeEach(Cgreen) {}
AfterEach(Cgreen) {}

using namespace cgreen;

void init_log() {
    mdclog_attr_t *attr;
    mdclog_attr_init(&attr);
    mdclog_attr_set_ident(attr, "TestConfiguration");
    mdclog_init(attr);
    mdclog_attr_destroy(attr);
}

Ensure(Cgreen, fileNotExist) {
    ReadConfigFile  conf {};
    assert_that( conf.openConfigFile("kuku") == -1);
}

Ensure(Cgreen, fileExists) {
    ReadConfigFile  conf {};
    assert_that(conf.openConfigFile("config/config.conf") == 0);
}

Ensure(Cgreen, goodparams) {
    ReadConfigFile  conf {};
    assert_that(conf.openConfigFile("config/config.conf") == 0);
    assert_that(conf.getIntValue("nano") == 38000);
    assert_that(conf.getStringValue("loglevel") == "info");
    assert_that(conf.getStringValue("volume") == "log");

}

Ensure(Cgreen, badParams) {
    ReadConfigFile  conf {};
    assert_that(conf.openConfigFile("config/config.conf") == 0);
    assert_that(conf.getIntValue("nano") != 38002);
    assert_that(conf.getStringValue("loglevel") != "");
    assert_that(conf.getStringValue("volume") != "bob");
    assert_that(conf.getStringValue("volum") != "bob");
}

Ensure(Cgreen, wrongType) {
    ReadConfigFile  conf {};
    assert_that(conf.openConfigFile("config/config.conf") == 0);
    assert_that(conf.getStringValue("nano") != "debug");
    assert_that(conf.getIntValue("loglevel") != 3);
    assert_that(conf.getDoubleValue("loglevel") != 3.0);
}

Ensure(Cgreen, badValues) {
    ReadConfigFile  conf {};
    assert_that(conf.openConfigFile("config/config.bad") == 0);
}

Ensure(Cgreen, sectionTest) {
    ReadConfigFile  conf {};
    assert_that(conf.openConfigFile("config/config.sec") == 0);
    assert_that(conf.getIntValue("config.nano") == 38000);
}

Ensure(Cgreen, sectionBadTest) {
    ReadConfigFile  conf {};
    assert_that(conf.openConfigFile("config/config.secbad") == -1);
    //assert_that(conf.getIntValue("config.nano") == 38000);
}

int main(const int argc, char **argv) {
    mdclog_severity_t loglevel = MDCLOG_INFO;
    init_log();
    mdclog_level_set(loglevel);

    //TestSuite *suite = create_test_suite();
    TestSuite *suite = create_named_test_suite_(__FUNCTION__, __FILE__, __LINE__);

    add_test_with_context(suite, Cgreen, fileNotExist);
    add_test_with_context(suite, Cgreen, fileExists);
    add_test_with_context(suite, Cgreen, goodparams);
    add_test_with_context(suite, Cgreen, badParams);
    add_test_with_context(suite, Cgreen, wrongType);
    add_test_with_context(suite, Cgreen, badValues);
    add_test_with_context(suite, Cgreen, sectionTest);
    add_test_with_context(suite, Cgreen, sectionBadTest);

    return cgreen::run_test_suite(suite, create_text_reporter());

}