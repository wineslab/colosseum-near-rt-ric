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


#include <unistd.h>
#include <string.h>

extern "C" {
#include "redismodule.h"
#include <stdio.h>
#include <stdlib.h>
#include <ctype.h>
}

#include <CppUTest/TestHarness.h>
#include <CppUTestExt/MockSupport.h>

#include "ut_helpers.hpp"

RedisModuleCallReply *RedisModule_Call(RedisModuleCtx *ctx, const char *cmdname, const char *fmt, ...)
{
    (void)ctx;
    (void)fmt;
    return (RedisModuleCallReply *)mock().actualCall("RedisModule_Call")
                                         .withParameter("cmdname", cmdname)
                                         .returnPointerValueOrDefault(malloc(UT_DUMMY_BUFFER_SIZE));
}

int RedisModule_ReplyWithString(RedisModuleCtx *ctx, RedisModuleString *str)
{
    (void)ctx;
    (void)str;
    return mock()
        .actualCall("RedisModule_ReplyWithString")
        .returnIntValueOrDefault(REDISMODULE_OK);
}

RedisModuleString *RedisModule_CreateStringFromCallReply(RedisModuleCallReply *reply)
{
    (void)reply;
    return (RedisModuleString *)mock().actualCall("RedisModule_CreateStringFromCallReply")
        .returnPointerValue();
}

void RedisModule_CloseKey(RedisModuleKey *kp)
{
    (void)kp;
    mock().actualCall("RedisModule_CloseKey");
}

size_t RedisModule_CallReplyLength(RedisModuleCallReply *reply)
{
    (void)reply;
    return (size_t)mock().actualCall("RedisModule_CallReplyLength")
                         .returnIntValue();
}

int RedisModule_ReplyWithArray(RedisModuleCtx *ctx, long len)
{
    (void)ctx;
    return (int)mock().actualCall("RedisModule_ReplyWithArray")
                      .withParameter("len", len)
                      .returnIntValueOrDefault(REDISMODULE_OK);
}

void RedisModule_ReplySetArrayLength(RedisModuleCtx *ctx, long len)
{
    (void)ctx;
    mock().actualCall("RedisModule_ReplySetArrayLength")
                      .withParameter("len", len);
}

RedisModuleString *RedisModule_CreateString(RedisModuleCtx *ctx, const char *ptr, size_t len)
{
    (void)ctx;
    (void)ptr;
    (void)len;
    void* buf = malloc(UT_DUMMY_BUFFER_SIZE);
    return (RedisModuleString *) mock()
        .actualCall("RedisModule_CreateString")
        .returnPointerValueOrDefault(buf);
}

RedisModuleString *RedisModule_CreateStringFromLongLong(RedisModuleCtx *ctx, long long ll)
{
    (void)ctx;
    (void)ll;
    void* buf = malloc(UT_DUMMY_BUFFER_SIZE);
    return (RedisModuleString *)mock()
        .actualCall("RedisModule_CreateStringFromLongLong")
        .returnPointerValueOrDefault(buf);
}

void RedisModule_AutoMemory(RedisModuleCtx *ctx)
{
    (void)ctx;
    mock().actualCall("RedisModule_AutoMemory");
}

void RedisModule_FreeString(RedisModuleCtx *ctx, RedisModuleString *str)
{
    (void)ctx;
    free(str);
    mock().actualCall("RedisModule_FreeString");
}

int RedisModule_StringToLongLong(const RedisModuleString *str, long long *ll)
{
    (void)str;
    return (int)mock().actualCall("RedisModule_StringToLongLong")
                      .withOutputParameter("ll", ll)
                      .returnIntValueOrDefault(REDISMODULE_OK);
}

void RedisModule_FreeCallReply(RedisModuleCallReply *reply)
{
    free(reply);
    mock().actualCall("RedisModule_FreeCallReply");
}

RedisModuleCallReply *RedisModule_CallReplyArrayElement(RedisModuleCallReply *reply, size_t idx)
{
    (void)reply;
    (void)idx;
    return (RedisModuleCallReply *)mock()
        .actualCall("RedisModule_CallReplyArrayElement")
        .returnPointerValueOrDefault(NULL);
}

int RedisModule_ReplyWithLongLong(RedisModuleCtx *ctx, long long ll)
{
    (void)ctx;
    return (int)mock()
        .actualCall("RedisModule_ReplyWithLongLong")
        .withParameter("ll", (int)ll)
        .returnIntValueOrDefault(REDISMODULE_OK);
}

long long RedisModule_CallReplyInteger(RedisModuleCallReply *reply)
{
    (void)reply;
    return (long long)mock()
        .actualCall("RedisModule_CallReplyInteger")
        .returnIntValue();
}

int RedisModule_CallReplyType(RedisModuleCallReply *reply)
{
    (void)reply;
    return (int)mock()
        .actualCall("RedisModule_CallReplyType")
        .returnIntValue();
}

int RedisModule_WrongArity(RedisModuleCtx *ctx)
{
    (void)ctx;
    return (int)mock()
        .actualCall("RedisModule_WrongArity")
        .returnIntValueOrDefault(REDISMODULE_ERR);
}

int RedisModule_ReplyWithError(RedisModuleCtx *ctx, const char *err)
{
    (void)ctx;
    (void)err;
    return (int)mock()
        .actualCall("RedisModule_ReplyWithError")
        .returnIntValueOrDefault(REDISMODULE_OK);
}

int RedisModule_ReplyWithCallReply(RedisModuleCtx *ctx, RedisModuleCallReply *reply)
{
    (void)ctx;
    (void)reply;
    return (int)mock()
        .actualCall("RedisModule_ReplyWithCallReply")
        .returnIntValueOrDefault(REDISMODULE_OK);
}

void *RedisModule_OpenKey(RedisModuleCtx *ctx, RedisModuleString *keyname, int mode)
{
    (void)ctx;
    (void)keyname;
    (void)mode;
    return (void *)mock()
        .actualCall("RedisModule_OpenKey")
        .returnPointerValue();
}

int RedisModule_KeyType(RedisModuleKey *kp)
{
    (void)kp;
    return (int)mock()
        .actualCall("RedisModule_KeyType")
        .returnIntValue();
}

const char *RedisModule_StringPtrLen(const RedisModuleString *str, size_t *len)
{
    (void)str;
    if (len != NULL) {
        return (const char *)mock()
            .actualCall("RedisModule_StringPtrLen")
            .withOutputParameter("len", len)
            .returnPointerValue();
    } else {
        return (const char *)mock()
            .actualCall("RedisModule_StringPtrLen")
            .returnPointerValue();
    }
}

RedisModuleBlockedClient *RedisModule_BlockClient(RedisModuleCtx *ctx, RedisModuleCmdFunc reply_callback,
                                                  RedisModuleCmdFunc timeout_callback,
                                                  void (*free_privdata)(RedisModuleCtx*,void*),
                                                  long long timeout_ms)
{
    (void)ctx;
    (void)reply_callback;
    (void)timeout_callback;
    (void)free_privdata;
    (void)timeout_ms;

    void *buf = malloc(UT_DUMMY_BUFFER_SIZE);
    return (RedisModuleBlockedClient *)mock()
        .actualCall("RedisModule_BlockClient")
        .returnPointerValueOrDefault(buf);
}

int RedisModule_UnblockClient(RedisModuleBlockedClient *bc, void *privdata)
{
    (void)privdata;

    free(bc);
    return (int)mock()
        .actualCall("RedisModule_UnblockClient")
        .returnIntValueOrDefault(REDISMODULE_OK);
}

const char *RedisModule_CallReplyStringPtr(RedisModuleCallReply *reply, size_t *len)
{
    (void)reply;
    (void)len;

    static char cursor_zero_literal[] = "0";
    return (const char *)mock()
        .actualCall("RedisModule_CallReplyStringPtr")
        .returnPointerValueOrDefault(cursor_zero_literal);
}

int RedisModule_AbortBlock(RedisModuleBlockedClient *bc)
{
    free(bc);
    return mock()
        .actualCall("RedisModule_AbortBlock")
        .returnIntValueOrDefault(REDISMODULE_OK);
}

int RedisModule_ReplyWithNull(RedisModuleCtx *ctx)
{
    (void)ctx;
    return mock()
        .actualCall("RedisModule_ReplyWithNull")
        .returnIntValueOrDefault(REDISMODULE_OK);
}

void RedisModule_ThreadSafeContextUnlock(RedisModuleCtx *ctx)
{
    (void)ctx;
    int tmp = mock().getData("TimesThreadSafeContextWasUnlocked").getIntValue();
    mock().setData("TimesThreadSafeContextWasUnlocked", tmp + 1);
    mock()
        .actualCall("RedisModule_ThreadSafeContextUnlock");
}

void RedisModule_ThreadSafeContextLock(RedisModuleCtx *ctx)
{
    (void)ctx;
    int tmp = mock().getData("TimesThreadSafeContextWasLocked").getIntValue();
    mock().setData("TimesThreadSafeContextWasLocked", tmp + 1);
    mock()
        .actualCall("RedisModule_ThreadSafeContextLock");
}

RedisModuleCtx *RedisModule_GetThreadSafeContext(RedisModuleBlockedClient *bc)
{
    (void)bc;
    return (RedisModuleCtx *)mock()
        .actualCall("RedisModule_GetThreadSafeContext")
        .returnPointerValueOrDefault(0);
}

void RedisModule_FreeThreadSafeContext(RedisModuleCtx *ctx)
{
    (void)ctx;
    mock()
        .actualCall("RedisModule_FreeThreadSafeContext");
}

/* This is included inline inside each Redis module. */
int RedisModule_Init(RedisModuleCtx *ctx, const char *name, int ver, int apiver)
{
    (void)ctx;
    (void)name;
    (void)ver;
    (void)apiver;

    return mock()
        .actualCall("RedisModule_Init")
        .returnIntValueOrDefault(REDISMODULE_OK);
}

int RedisModule_CreateCommand(RedisModuleCtx *ctx, const char *name, RedisModuleCmdFunc cmdfunc, const char *strflags, int firstkey, int lastkey, int keystep)
{
    (void)ctx;
    (void)name;
    (void)cmdfunc;
    (void)strflags;
    (void)firstkey;
    (void)lastkey;
    (void)keystep;

    return mock()
        .actualCall("RedisModule_CreateCommand")
        .returnIntValueOrDefault(REDISMODULE_OK);
}

void *RedisModule_Alloc(size_t bytes)
{
    void *buf = malloc(bytes);
    return mock()
        .actualCall("RedisModule_Alloc")
        .returnPointerValueOrDefault(buf);
}

void RedisModule_Free(void *ptr)
{
    free(ptr);
    mock()
        .actualCall("RedisModule_Free");
}
