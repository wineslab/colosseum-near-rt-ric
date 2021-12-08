/*
 * Copyright (c) 2018-2020 Nokia.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

/*
 * This source code is part of the near-RT RIC (RAN Intelligent Controller)
 * platform project (RICP).
 */

extern "C" {
#include "exstringsStub.h"
#include "redismodule.h"
}

#include "CppUTest/TestHarness.h"
#include "CppUTestExt/MockSupport.h"

#include "ut_helpers.hpp"

void nDelReturnNKeysFromUnlink(int count)
{
    mock()
        .expectOneCall("RedisModule_CallReplyInteger")
        .andReturnValue(count);

}

TEST_GROUP(exstrings_ndel)
{
    void setup()
    {
        mock().enable();
        mock().ignoreOtherCalls();
    }

    void teardown()
    {
        mock().clear();
        mock().disable();
    }

};

TEST(exstrings_ndel, ndel_atomic_automemory_enabled)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_AutoMemory");
    int ret = NDel_Atomic_RedisCommand(&ctx, redisStrVec,  3);
    mock().checkExpectations();
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    delete []redisStrVec;
}

TEST(exstrings_ndel, ndel_atomic_command_parameter_parameter_number_incorrect)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_WrongArity");
    int ret = NDel_Atomic_RedisCommand(&ctx, redisStrVec,  3);
    mock().checkExpectations();
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    delete []redisStrVec;
}

TEST(exstrings_ndel, ndel_atomic_command_scan_0_keys_found)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    returnNKeysFromScanSome(0);
    mock().expectOneCall("RedisModule_ReplyWithLongLong")
          .withParameter("ll", 0);
    int ret = NDel_Atomic_RedisCommand(&ctx, redisStrVec,  2);
    mock().checkExpectations();
    CHECK_EQUAL(ret, REDISMODULE_OK);

    delete []redisStrVec;
}

TEST(exstrings_ndel, ndel_atomic_command_scan_3_keys_found_3_keys_deleted)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(3);
    mock()
        .expectOneCall("RedisModule_Call")
        .withParameter("cmdname", "UNLINK");
    nDelReturnNKeysFromUnlink(3);
    mock().expectOneCall("RedisModule_ReplyWithLongLong")
          .withParameter("ll", 3);
    int ret = NDel_Atomic_RedisCommand(&ctx, redisStrVec,  2);
    mock().checkExpectations();
    CHECK_EQUAL(ret, REDISMODULE_OK);

    delete []redisStrVec;
}

TEST(exstrings_ndel, ndel_atomic_command_scan_3_keys_found_0_keys_deleted)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(3);
    mock()
        .expectOneCall("RedisModule_Call")
        .withParameter("cmdname", "UNLINK");
    nDelReturnNKeysFromUnlink(0);
    mock().expectOneCall("RedisModule_ReplyWithLongLong")
          .withParameter("ll", 0);
    int ret = NDel_Atomic_RedisCommand(&ctx, redisStrVec,  2);
    mock().checkExpectations();
    CHECK_EQUAL(ret, REDISMODULE_OK);

    delete []redisStrVec;
}

TEST(exstrings_ndel, ndel_atomic_command_scan_3_keys_found_1_keys_deleted)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(3);
    mock()
        .expectOneCall("RedisModule_Call")
        .withParameter("cmdname", "UNLINK");
    nDelReturnNKeysFromUnlink(1);
    mock().expectOneCall("RedisModule_ReplyWithLongLong")
          .withParameter("ll", 1);
    int ret = NDel_Atomic_RedisCommand(&ctx, redisStrVec,  2);
    mock().checkExpectations();
    CHECK_EQUAL(ret, REDISMODULE_OK);

    delete []redisStrVec;
}
