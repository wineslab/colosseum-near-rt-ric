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

#include <string.h>

#include "CppUTest/TestHarness.h"
#include "CppUTestExt/MockSupport.h"

#include "ut_helpers.hpp"

TEST_GROUP(exstrings_nget)
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

void threadDetachedSuccess()
{
    mock().expectOneCall("pthread_detach")
        .andReturnValue(0);
}

void nKeysFoundMget(long keys)
{
    for (long i = 0 ; i < keys ; i++) {
        mock().expectOneCall("RedisModule_CreateStringFromCallReply")
              .andReturnValue(malloc(UT_DUMMY_BUFFER_SIZE));
        mock().expectNCalls(2, "RedisModule_ReplyWithString");
    }
}

void nKeysNotFoundMget(long keys)
{
    void* ptr = NULL;
    mock().expectNCalls(keys, "RedisModule_CreateStringFromCallReply")
          .andReturnValue(ptr);
    mock().expectNoCall("RedisModule_ReplyWithString");
}

void expectNReplies(long count)
{
    mock().expectOneCall("RedisModule_ReplySetArrayLength")
          .withParameter("len", 2*count);
}

void threadSafeContextLockedAndUnlockedEqualTimes()
{
    int locked = mock().getData("TimesThreadSafeContextWasLocked").getIntValue();
    int unlocked = mock().getData("TimesThreadSafeContextWasUnlocked").getIntValue();
    CHECK_EQUAL(locked, unlocked);
}

TEST(exstrings_nget, nget_atomic_automemory_enabled)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);
    long keys_found_with_scan = 0;

    mock().expectOneCall("RedisModule_AutoMemory");
    mock().expectOneCall("RedisModule_CallReplyLength")
          .andReturnValue((int)keys_found_with_scan);
    int ret = NGet_Atomic_RedisCommand(&ctx, redisStrVec,  2);
    mock().checkExpectations();
    CHECK_EQUAL(ret, REDISMODULE_OK);

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_atomic_command_parameter_number_incorrect)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().expectOneCall("RedisModule_WrongArity");
    int ret = NGet_Atomic_RedisCommand(&ctx, redisStrVec,  3);
    CHECK_EQUAL(ret, REDISMODULE_ERR);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_atomic_command_3rd_parameter_was_not_equal_to_COUNT)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(4);
    const char * not_count_literal = "NOT_COUNT";
    size_t not_count_len = strlen(not_count_literal);

    mock().expectOneCall("RedisModule_StringPtrLen")
          .withOutputParameterReturning("len", &not_count_len, sizeof(size_t))
          .andReturnValue((void*)not_count_literal);
    mock().expectOneCall("RedisModule_ReplyWithError");
    int ret = NGet_Atomic_RedisCommand(&ctx, redisStrVec,  4);
    CHECK_EQUAL(ret, REDISMODULE_ERR);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_atomic_command_4th_parameter_was_not_integer)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(4);
    const char * count_literal = "COUNT";
    size_t count_len = strlen(count_literal);
    size_t count_number = 123;

    mock().expectOneCall("RedisModule_StringPtrLen")
          .withOutputParameterReturning("len", &count_len, sizeof(size_t))
          .andReturnValue((void*)count_literal);
    mock().expectOneCall("RedisModule_StringToLongLong")
          .withOutputParameterReturning("ll", &count_number, sizeof(size_t))
          .andReturnValue(REDISMODULE_ERR);
    mock().expectOneCall("RedisModule_ReplyWithError");
    int ret = NGet_Atomic_RedisCommand(&ctx, redisStrVec,  4);
    CHECK_EQUAL(ret, REDISMODULE_ERR);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_atomic_command_4th_parameter_was_negative)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(4);
    const char * count_literal = "COUNT";
    size_t count_len = strlen(count_literal);
    size_t count_number = -123;

    mock().expectOneCall("RedisModule_StringPtrLen")
          .withOutputParameterReturning("len", &count_len, sizeof(size_t))
          .andReturnValue((void*)count_literal);
    mock().expectOneCall("RedisModule_StringToLongLong")
          .withOutputParameterReturning("ll", &count_number, sizeof(size_t))
          .andReturnValue(REDISMODULE_OK);
    mock().expectOneCall("RedisModule_ReplyWithError");
    int ret = NGet_Atomic_RedisCommand(&ctx, redisStrVec,  4);
    CHECK_EQUAL(ret, REDISMODULE_ERR);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_atomic_command_scan_returned_zero_keys)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().expectOneCall("RedisModule_ReplyWithArray")
          .withParameter("len", (long)REDISMODULE_POSTPONED_ARRAY_LEN);
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(0);
    expectNReplies(0);
    mock().expectNoCall("RedisModule_Call");

    int ret = NGet_Atomic_RedisCommand(&ctx, redisStrVec,  2);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_atomic_command_3_keys_scanned_0_keys_mget)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_ReplyWithArray")
          .withParameter("len", (long)REDISMODULE_POSTPONED_ARRAY_LEN);

    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "MGET");
    nKeysNotFoundMget(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    expectNReplies(0);
    int ret = NGet_Atomic_RedisCommand(&ctx, redisStrVec,  2);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_atomic_command_3_keys_scanned_3_keys_mget)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_ReplyWithArray")
          .withParameter("len", (long)REDISMODULE_POSTPONED_ARRAY_LEN);
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "MGET");
    nKeysFoundMget(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    expectNReplies(3);
    int ret = NGet_Atomic_RedisCommand(&ctx, redisStrVec,  2);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_atomic_command_3_keys_scanned_2_keys_mget)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_ReplyWithArray")
          .withParameter("len", (long)REDISMODULE_POSTPONED_ARRAY_LEN);
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "MGET");
    nKeysFoundMget(2);
    nKeysNotFoundMget(1);
    mock().expectOneCall("RedisModule_FreeCallReply");
    expectNReplies(2);
    int ret = NGet_Atomic_RedisCommand(&ctx, redisStrVec,  2);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_automemory_enabled)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().setData("pthread_create_free_block_client_args", 1);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_AutoMemory");

    int ret = NGet_NoAtomic_RedisCommand(&ctx, redisStrVec,  2);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_thread_create_success)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().setData("pthread_create_free_block_client_args", 1);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_BlockClient");
    mock().expectOneCall("pthread_create");
    mock().expectNoCall("RedisModule_AbortBlock");

    int ret = NGet_NoAtomic_RedisCommand(&ctx, redisStrVec,  2);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_thread_create_fail)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    mock().ignoreOtherCalls();
    mock().expectOneCall("RedisModule_BlockClient");
    mock().expectOneCall("pthread_create")
          .andReturnValue(1);
    mock().expectOneCall("RedisModule_AbortBlock");

    int ret = NGet_NoAtomic_RedisCommand(&ctx, redisStrVec,  2);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_parameter_number_incorrect)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(4);

    mock().expectOneCall("RedisModule_WrongArity");
    mock().expectNoCall("RedisModule_BlockClient");

    int ret = NGet_NoAtomic_RedisCommand(&ctx, redisStrVec,  3);

    CHECK_EQUAL(ret, REDISMODULE_ERR);
    mock().checkExpectations();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_threadmain_3rd_parameter_was_not_equal_to_COUNT)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = createRedisStrVec(4);
    const char * not_count_literal = "NOT_COUNT";
    size_t not_count_len = strlen(not_count_literal);

    mock().expectOneCall("RedisModule_StringPtrLen")
          .withOutputParameterReturning("len", &not_count_len, sizeof(size_t))
          .andReturnValue((void*)not_count_literal);
    mock().expectOneCall("RedisModule_ReplyWithError");
    mock().expectNoCall("RedisModule_BlockClient");

    int ret = NGet_NoAtomic_RedisCommand(&ctx, redisStrVec,  4);

    CHECK_EQUAL(ret, REDISMODULE_ERR);

    mock().checkExpectations();
    threadSafeContextLockedAndUnlockedEqualTimes();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_4th_parameter_was_not_integer)
{
    RedisModuleCtx ctx;
    const char * count_literal = "COUNT";
    size_t count_len = strlen(count_literal);
    size_t count_number = -123;
    RedisModuleString ** redisStrVec = createRedisStrVec(4);

    mock().expectOneCall("RedisModule_StringPtrLen")
          .withOutputParameterReturning("len", &count_len, sizeof(size_t))
          .andReturnValue((void*)count_literal);
    mock().expectOneCall("RedisModule_StringToLongLong")
          .withOutputParameterReturning("ll", &count_number, sizeof(size_t))
          .andReturnValue(REDISMODULE_OK);
    mock().expectOneCall("RedisModule_ReplyWithError");

    int ret = NGet_NoAtomic_RedisCommand(&ctx, redisStrVec,  4);

    CHECK_EQUAL(ret, REDISMODULE_ERR);
    mock().checkExpectations();
    threadSafeContextLockedAndUnlockedEqualTimes();

    delete []redisStrVec;
}

typedef struct RedisModuleBlockedClientArgs {
    RedisModuleBlockedClient *bc;
    RedisModuleString **argv;
    int argc;
} RedisModuleBlockedClientArgs;

TEST(exstrings_nget, nget_noatomic_threadmain_3_keys_scanned_3_keys_mget)
{
    RedisModuleCtx ctx;
    RedisModuleBlockedClientArgs *bca =
        (RedisModuleBlockedClientArgs*)RedisModule_Alloc(sizeof(RedisModuleBlockedClientArgs));
    RedisModuleBlockedClient *bc = RedisModule_BlockClient(&ctx,NULL,NULL,NULL,0);
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    bca->bc = bc;
    bca->argv = redisStrVec;
    bca->argc = 2;

    mock().ignoreOtherCalls();
    threadDetachedSuccess();
    mock().expectOneCall("RedisModule_ReplyWithArray")
          .withParameter("len", (long)REDISMODULE_POSTPONED_ARRAY_LEN);
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "MGET");
    nKeysFoundMget(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    expectNReplies(3);
    mock().expectOneCall("RedisModule_FreeThreadSafeContext");
    mock().expectOneCall("RedisModule_UnblockClient");

    NGet_NoAtomic_ThreadMain((void*)bca);

    mock().checkExpectations();
    threadSafeContextLockedAndUnlockedEqualTimes();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_threadmain_3_keys_scanned_0_keys_mget)
{
    RedisModuleCtx ctx;
    RedisModuleBlockedClientArgs *bca = (RedisModuleBlockedClientArgs*)malloc(sizeof(RedisModuleBlockedClientArgs));
    RedisModuleBlockedClient *bc = RedisModule_BlockClient(&ctx,NULL,NULL,NULL,0);
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    bca->bc = bc;
    bca->argv = redisStrVec;
    bca->argc = 2;

    mock().ignoreOtherCalls();
    threadDetachedSuccess();
    mock().expectOneCall("RedisModule_ReplyWithArray")
          .withParameter("len", (long)REDISMODULE_POSTPONED_ARRAY_LEN);
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "MGET");
    nKeysNotFoundMget(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    expectNReplies(0);
    mock().expectOneCall("RedisModule_FreeThreadSafeContext");
    mock().expectOneCall("RedisModule_UnblockClient");

    NGet_NoAtomic_ThreadMain((void*)bca);

    mock().checkExpectations();
    threadSafeContextLockedAndUnlockedEqualTimes();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_threadmain_3_keys_scanned_2_keys_mget)
{
    RedisModuleCtx ctx;
    RedisModuleBlockedClientArgs *bca = (RedisModuleBlockedClientArgs*)malloc(sizeof(RedisModuleBlockedClientArgs));
    RedisModuleBlockedClient *bc = RedisModule_BlockClient(&ctx,NULL,NULL,NULL,0);
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    bca->bc = bc;
    bca->argv = redisStrVec;
    bca->argc = 2;

    mock().ignoreOtherCalls();
    threadDetachedSuccess();
    mock().expectOneCall("RedisModule_ReplyWithArray")
          .withParameter("len", (long)REDISMODULE_POSTPONED_ARRAY_LEN);
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(3);
    mock().expectOneCall("RedisModule_FreeCallReply");
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "MGET");
    nKeysNotFoundMget(1);
    nKeysFoundMget(2);
    mock().expectOneCall("RedisModule_FreeCallReply");
    expectNReplies(2);
    mock().expectOneCall("RedisModule_FreeThreadSafeContext");
    mock().expectOneCall("RedisModule_UnblockClient");

    NGet_NoAtomic_ThreadMain((void*)bca);

    mock().checkExpectations();
    threadSafeContextLockedAndUnlockedEqualTimes();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_threadmain_scan_returned_zero_keys)
{
    RedisModuleCtx ctx;
    RedisModuleBlockedClientArgs *bca = (RedisModuleBlockedClientArgs*)malloc(sizeof(RedisModuleBlockedClientArgs));
    RedisModuleBlockedClient *bc = RedisModule_BlockClient(&ctx,NULL,NULL,NULL,0);
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    bca->bc = bc;
    bca->argv = redisStrVec;
    bca->argc = 2;

    mock().ignoreOtherCalls();
    threadDetachedSuccess();
    mock().expectOneCall("RedisModule_ReplyWithArray")
          .withParameter("len", (long)REDISMODULE_POSTPONED_ARRAY_LEN);
    mock().expectOneCall("RedisModule_Call")
          .withParameter("cmdname", "SCAN");
    returnNKeysFromScanSome(0);
    mock().expectOneCall("RedisModule_FreeCallReply");
    mock().expectNoCall("RedisModule_Call");
    mock().expectOneCall("RedisModule_FreeThreadSafeContext");
    mock().expectOneCall("RedisModule_UnblockClient");

    NGet_NoAtomic_ThreadMain((void*)bca);

    mock().checkExpectations();
    threadSafeContextLockedAndUnlockedEqualTimes();

    delete []redisStrVec;
}

TEST(exstrings_nget, nget_noatomic_threadmain_thread_detached)
{
    RedisModuleCtx ctx;
    RedisModuleBlockedClientArgs *bca = (RedisModuleBlockedClientArgs*)malloc(sizeof(RedisModuleBlockedClientArgs));
    RedisModuleBlockedClient *bc = RedisModule_BlockClient(&ctx,NULL,NULL,NULL,0);
    RedisModuleString ** redisStrVec = createRedisStrVec(2);

    bca->bc = bc;
    bca->argv = redisStrVec;
    bca->argc = 2;

    mock().ignoreOtherCalls();
    threadDetachedSuccess();

    NGet_NoAtomic_ThreadMain((void*)bca);

    mock().checkExpectations();

    delete []redisStrVec;
}
