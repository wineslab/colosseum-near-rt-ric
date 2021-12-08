# Copyright 2019 AT&T Intellectual Property
# Copyright 2019 Nokia
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This source code is part of the near-RT RIC (RAN Intelligent Controller)
# platform project (RICP).

cmake_minimum_required(VERSION 3.13)
project(e2)

set(CMAKE_VERBOSE_MAKEFILE off)

set(CMAKE_CXX_STANDARD 17)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -DDEV_PKG=1")
set(PROJECT_NAME "e2")
set(PROJECT_TEST_NAME "e2")

if (NOT CMAKE_BUILD_TYPE STREQUAL "Debug")
    #set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -DASN_DISABLE_OER_SUPPORT -DASN_PDU_COLLECTION -L. -LRIC-E2-TERMINATION/tracelibcpp/build -ggdb3 -Wall -W -Wpedantic")
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -DASN_DISABLE_OER_SUPPORT -DASN_PDU_COLLECTION -L. -g -ggdb3 -O3 -L/usr/lib -L/usr/local/lib -Wall -Wpedantic")
    #only c code with -O3
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -DASN_DISABLE_OER_SUPPORT -DASN_PDU_COLLECTION -L. -O3 -L/usr/lib -L/usr/local/lib -Wall -W -Wpedantic")
else ()
    set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -DASN_DISABLE_OER_SUPPORT -DASN_PDU_COLLECTION -L. -ggdb3 --coverage -L/usr/lib  -L/usr/local/lib -Wall -Wpedantic")
    #only c code with -O3
    set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -DASN_DISABLE_OER_SUPPORT -DASN_PDU_COLLECTION -L. -ggdb3 -L/usr/lib  -L/usr/local/lib -Wall -W -Wpedantic")

endif ()

include_directories(RIC-E2-TERMINATION
        RIC-E2-TERMINATION/3rdparty/oranE2
        RIC-E2-TERMINATION/3rdparty/oranE2SM
        RIC-E2-TERMINATION/3rdparty/cxxopts/include
        RIC-E2-TERMINATION/3rdparty
        cmake-modules)

#E2AP library
file(GLOB E2AP_ASN_MODULE_SRCS "RIC-E2-TERMINATION/3rdparty/oranE2/*.c")
file(GLOB E2AP_ASN_MODULE_HDRS "RIC-E2-TERMINATION/3rdparty/oranE2/*.h")

add_library(asn1codec ${E2AP_ASN_MODULE_SRCS} ${E2AP_ASN_MODULE_HDRS})
install(TARGETS asn1codec DESTINATION /usr/lib)
install(FILES ${E2AP_E2AP_ASN_MODULE_HDRS} DESTINATION /usr/include/asn1c)


#E2SM library
file(GLOB E2SM_ASN_MODULE_SRCS "RIC-E2-TERMINATION/3rdparty/oranE2SM/*.c")
file(GLOB E2SM_ASN_MODULE_HDRS "RIC-E2-TERMINATION/3rdparty/oranE2SM/*.h")

add_library(asn1ce2smcodec ${E2SM_ASN_MODULE_SRCS} ${E2SM_ASN_MODULE_HDRS})
#add_custom_command(
#        TARGET asn1ce2smcodec
#        POST_BUILD
#        COMMAND objcopy
#        ARGS --prefix-symbols=e2sm_ libasn1ce2smcodec.a
#)

install(TARGETS asn1ce2smcodec DESTINATION /usr/lib)
install(FILES ${E2SM_ASN_MODULE_HDRS} DESTINATION /usr/include/asn1c)

include_directories(RIC-E2-TERMINATION/TEST)
include_directories(RIC-E2-TERMINATION/TEST/e2smTest)
include_directories(RIC-E2-TERMINATION/TEST/T1)
include_directories(RIC-E2-TERMINATION/TEST/T2)


add_definitions(-DBOOST_LOG_DYN_LINK)

link_libraries(nsl
        sctp
        gcov
        c
        m
        dl
        mdclog
        rmr_si
        asn1codec
        asn1ce2smcodec
        boost_system
        boost_log_setup
        boost_log
        boost_date_time
        boost_thread
        boost_system
        rt
        tbb
        boost_filesystem
        cgreen
        pthread)

add_executable(e2 RIC-E2-TERMINATION/sctpThread.cpp
        RIC-E2-TERMINATION/sctpThread.h
        RIC-E2-TERMINATION/openTracing.h
        RIC-E2-TERMINATION/mapWrapper.h
        RIC-E2-TERMINATION/base64.h
        RIC-E2-TERMINATION/base64.cpp
        RIC-E2-TERMINATION/ReadConfigFile.h
        RIC-E2-TERMINATION/BuildRunName.h
        RIC-E2-TERMINATION/statCollector.h
        RIC-E2-TERMINATION/3rdparty/pugixml/src/pugiconfig.hpp
        RIC-E2-TERMINATION/3rdparty/pugixml/src/pugixml.cpp
        RIC-E2-TERMINATION/3rdparty/pugixml/src/pugixml.hpp
        RIC-E2-TERMINATION/BuildXml.h)
target_link_libraries(e2 libasn1ce2smcodec.a)
target_link_libraries(e2 librmr_si.a)
target_link_libraries(e2 libicui18n.a)
target_link_libraries(e2 libicuuc.a)
target_link_libraries(e2 libicudata.a)

#target_link_libraries(e2 libnng.a)

add_executable(testConfigFile
        RIC-E2-TERMINATION/ReadConfigFile.h
        RIC-E2-TERMINATION/TEST/ConfigurationFileTest/testConfigFile.cpp)

add_executable(b64Test
        RIC-E2-TERMINATION/base64.cpp
        RIC-E2-TERMINATION/base64.h
        RIC-E2-TERMINATION/TEST/base64/testBase64.cpp)


add_executable(sctpClient
        RIC-E2-TERMINATION/TEST/testAsn/sctpClient/sctpClient.cpp
        RIC-E2-TERMINATION/TEST/testAsn/sctpClient/sctpClient.h
        RIC-E2-TERMINATION/TEST/testAsn/rmrClient/rmrClient.h
        RIC-E2-TERMINATION/TEST/testAsn/httpServer/HttpServer.cpp
        RIC-E2-TERMINATION/TEST/testAsn/httpServer/HttpServer.h
        RIC-E2-TERMINATION/base64.h
        RIC-E2-TERMINATION/base64.cpp
        RIC-E2-TERMINATION/TEST/T1/E2Builder.h

        #RIC-E2-TERMINATION/TEST/T1/Test1.cpp
        #RIC-E2-TERMINATION/TEST/T1/
        )
target_link_libraries(sctpClient libpistache.a)
target_link_libraries(sctpClient librmr_si.a)
#target_link_libraries(sctpClient libnng.a)

add_executable(setUpMessages
        RIC-E2-TERMINATION/TEST/testAsn/setUpMessages/SetUpMessages.cpp
        RIC-E2-TERMINATION/BuildRunName.h
        )

#find_package(RapidJSON)

add_executable(teste2
        RIC-E2-TERMINATION/TEST/e2test.h
        RIC-E2-TERMINATION/TEST/e2test.cpp
        )

if (CMAKE_BUILD_TYPE STREQUAL "Debug")
    LIST(APPEND CMAKE_MODULE_PATH "${PROJECT_SOURCE_DIR}/cmake_modules")
    set(CMAKE_BUILD_TYPE "Debug")
    include(cmake-modules/CodeCoverage.cmake)
    target_link_libraries(${PROJECT_TEST_NAME} gcov)

    set(LDFLAGS "--coverage -fprofile-arcs")

    #    setup_target_for_coverage_lcov(${PROJECT_NAME}_coverage ${PROJECT_TEST_NAME} coverage)
    set(COVERAGE_LCOV_EXCLUDES ${PROJECT_SOURCE_DIR}/config'
            ${PROJECT_SOURCE_DIR}/log'
            '${PROJECT_SOURCE_DIR}/docs'
            '${PROJECT_SOURCE_DIR}/RIC-E2-TERMINATION/3rdparty'
            '${PROJECT_SOURCE_DIR}/RIC-E2-TERMINATION/oranE2'
            '${PROJECT_SOURCE_DIR}/RIC-E2-TERMINATION/TEST'
            '${PROJECT_SOURCE_DIR}/RIC-E2-TERMINATION/config')

    set(COVERAGE_EXCLUDES '${PROJECT_SOURCE_DIR}/config'
            ${PROJECT_SOURCE_DIR}/log'
            '${PROJECT_SOURCE_DIR}/docs'
            '${PROJECT_SOURCE_DIR}/RIC-E2-TERMINATION/3rdparty'
            '${PROJECT_SOURCE_DIR}/RIC-E2-TERMINATION/oranE2'
            '${PROJECT_SOURCE_DIR}/RIC-E2-TERMINATION/TEST'
            '${PROJECT_SOURCE_DIR}/RIC-E2-TERMINATION/config')
    append_coverage_compiler_flags()
    SETUP_TARGET_FOR_COVERAGE_LCOV(NAME e2_coverage
            EXECUTABLE e2
            DEPENDENCIES e2)
endif ()
