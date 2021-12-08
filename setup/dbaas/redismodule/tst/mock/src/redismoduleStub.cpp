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


int RedisModule_CreateCommand(RedisModuleCtx *ctx, const char *name, RedisModuleCmdFunc cmdfunc, const char *strflags, int firstkey, int lastkey, int keystep)
{
    (void)ctx;
    (void)name;
    (void)cmdfunc;
    (void)strflags;
    (void)firstkey;
    (void)lastkey;
    (void)keystep;
    return REDISMODULE_OK;

}

int RedisModule_WrongArity(RedisModuleCtx *ctx)
{
    (void)ctx;
    mock().setData("RedisModule_WrongArity", 1);
    return REDISMODULE_ERR;
}

int RedisModule_ReplyWithLongLong(RedisModuleCtx *ctx, long long ll)
{

    (void)ctx;
    mock().setData("RedisModule_ReplyWithLongLong", (int)ll);
    return REDISMODULE_OK;
}

void *RedisModule_OpenKey(RedisModuleCtx *ctx, RedisModuleString *keyname, int mode)
{
    (void)ctx;
    (void)keyname;
    (void)mode;

    if (mock().hasData("RedisModule_OpenKey_no"))
    {
        return (void*)(0);
    }

    if (mock().hasData("RedisModule_OpenKey_have"))
    {
        return (void*)(111111);
    }


    return (void*)(0);
}

RedisModuleCallReply *RedisModule_Call(RedisModuleCtx *ctx, const char *cmdname, const char *fmt, ...)
{
    (void)ctx;
    (void)cmdname;
    (void)fmt;

    if (!strcmp(cmdname, "GET"))
        mock().setData("GET", 1);
    else if (!strcmp(cmdname, "SET"))
        mock().setData("SET", 1);
    else if (!strcmp(cmdname, "MSET"))
        mock().setData("MSET", 1);
    else if (!strcmp(cmdname, "DEL"))
        mock().setData("DEL", 1);
    else if (!strcmp(cmdname, "UNLINK"))
        mock().setData("UNLINK", 1);
    else if (!strcmp(cmdname, "PUBLISH"))
        mock().setData("PUBLISH", mock().getData("PUBLISH").getIntValue() + 1);
    else if (!strcmp(cmdname, "KEYS"))
        mock().setData("KEYS", 1);
    else if (!strcmp(cmdname, "MGET"))
        mock().setData("MGET", 1);
    else if (!strcmp(cmdname, "SCAN"))
        mock().setData("SCAN", 1);

    if (mock().hasData("RedisModule_Call_Return_Null"))
        return NULL;
    else
        return (RedisModuleCallReply *)1;
}

void RedisModule_FreeCallReply(RedisModuleCallReply *reply)
{
    (void)reply;
    mock().setData("RedisModule_FreeCallReply", mock().getData("RedisModule_FreeCallReply").getIntValue()+1);
}

int RedisModule_CallReplyType(RedisModuleCallReply *reply)
{

    (void)reply;
    if (mock().hasData("RedisModule_CallReplyType_null"))
    {
        return REDISMODULE_REPLY_NULL;
    }

    if (mock().hasData("RedisModule_CallReplyType_inter"))
    {
        return REDISMODULE_REPLY_INTEGER;
    }

    if (mock().hasData("RedisModule_CallReplyType_str"))
    {
        return REDISMODULE_REPLY_STRING;
    }

    if (mock().hasData("RedisModule_CallReplyType_err"))
    {
        return REDISMODULE_REPLY_ERROR;
    }

    return REDISMODULE_REPLY_NULL;;

}

long long RedisModule_CallReplyInteger(RedisModuleCallReply *reply)
{

    (void)reply;
    return mock().getData("RedisModule_CallReplyInteger").getIntValue();
}

const char *RedisModule_StringPtrLen(const RedisModuleString *str, size_t *len)
{

    (void)str;
    if (len) *len = 5;
    if (mock().hasData("RedisModule_String_same"))
    {
        return "11111";
    }

    if (mock().hasData("RedisModule_String_nosame"))
    {
        return "22222";
    }

    if (mock().hasData("RedisModule_String_count"))
    {
        return "COUNT";
    }

    if (mock().hasData("RedisModule_String_count1"))
    {
        if (len) *len = 6;
        return "COUNT1";
    }

    return "11111";
}

int RedisModule_ReplyWithError(RedisModuleCtx *ctx, const char *err)
{
    (void)ctx;
    (void)err;
    mock().setData("RedisModule_ReplyWithError", 1);
    return REDISMODULE_OK;
}

int RedisModule_ReplyWithString(RedisModuleCtx *ctx, RedisModuleString *str)
{
    (void)ctx;
    (void)str;
    mock().setData("RedisModule_ReplyWithString", mock().getData("RedisModule_ReplyWithString").getIntValue()+1);
    return REDISMODULE_OK;
}

int RedisModule_ReplyWithNull(RedisModuleCtx *ctx)
{

    (void)ctx;
    mock().setData("RedisModule_ReplyWithNull", 1);
    return REDISMODULE_OK;
}

int RedisModule_ReplyWithCallReply(RedisModuleCtx *ctx, RedisModuleCallReply *reply)
{
    (void)ctx;
    (void)reply;
    mock().setData("RedisModule_ReplyWithCallReply", 1);
    return REDISMODULE_OK;
}

const char *RedisModule_CallReplyStringPtr(RedisModuleCallReply *reply, size_t *len)
{
    (void)reply;

    if (mock().hasData("RedisModule_String_same"))
    {
        if (len)
            *len = 5;
        return "11111";
    }

    if (mock().hasData("RedisModule_String_nosame"))
    {
        if (len)
            *len = 6;
        return "333333";
    }

    return "11111";
}

RedisModuleString *RedisModule_CreateStringFromCallReply(RedisModuleCallReply *reply)
{
    (void)reply;
    return (RedisModuleString *)1;
}


int RedisModule_KeyType(RedisModuleKey *kp)
{


    (void)kp;
    if (mock().hasData("RedisModule_KeyType_empty"))
    {
        return REDISMODULE_KEYTYPE_EMPTY;
    }

    if (mock().hasData("RedisModule_KeyType_str"))
    {
        return REDISMODULE_KEYTYPE_STRING;
    }

    if (mock().hasData("RedisModule_KeyType_set"))
    {

        return REDISMODULE_KEYTYPE_SET;
    }

    return REDISMODULE_KEYTYPE_EMPTY;


}

void RedisModule_CloseKey(RedisModuleKey *kp)
{
    (void)kp;
    mock().actualCall("RedisModule_CloseKey");
}

/* This is included inline inside each Redis module. */
int RedisModule_Init(RedisModuleCtx *ctx, const char *name, int ver, int apiver)
{
    (void)ctx;
    (void)name;
    (void)ver;
    (void)apiver;
    return REDISMODULE_OK;
}

size_t RedisModule_CallReplyLength(RedisModuleCallReply *reply)
{
    (void)reply;
    return mock().getData("RedisModule_CallReplyLength").getIntValue();
}


RedisModuleCallReply *RedisModule_CallReplyArrayElement(RedisModuleCallReply *reply, size_t idx)
{
    (void)reply;
    (void)idx;
    return (RedisModuleCallReply *)1;
}

int RedisModule_ReplyWithArray(RedisModuleCtx *ctx, long len)
{
    (void)ctx;
    mock().setData("RedisModule_ReplyWithArray", (int)len);
    return REDISMODULE_OK;
}

void RedisModule_FreeString(RedisModuleCtx *ctx, RedisModuleString *str)
{
    (void)ctx;
    (void)str;
    mock().setData("RedisModule_FreeString", mock().getData("RedisModule_FreeString").getIntValue()+1);
    return;
}

RedisModuleBlockedClient *RedisModule_BlockClient(RedisModuleCtx *ctx, RedisModuleCmdFunc reply_callback, RedisModuleCmdFunc timeout_callback, void (*free_privdata)(RedisModuleCtx*,void*), long long timeout_ms)
{
    (void)ctx;
    (void)reply_callback;
    (void)timeout_callback;
    (void)free_privdata;
    (void)timeout_ms;
    RedisModuleBlockedClient *bc = (RedisModuleBlockedClient*)malloc(sizeof(RedisModuleBlockedClient));
    mock().setData("RedisModule_BlockClient", 1);
    return bc;
}

int RedisModule_UnblockClient(RedisModuleBlockedClient *bc, void *privdata)
{
    (void)privdata;
    free(bc);
    mock().setData("RedisModule_UnblockClient", mock().getData("RedisModule_UnblockClient").getIntValue()+1);
    return REDISMODULE_OK;
}

int RedisModule_AbortBlock(RedisModuleBlockedClient *bc)
{
    free(bc);
    mock().setData("RedisModule_AbortBlock", 1);
    return REDISMODULE_OK;
}

RedisModuleString *RedisModule_CreateString(RedisModuleCtx *ctx, const char *ptr, size_t len)
{
    (void)ctx;
    (void)ptr;
    (void)len;
    RedisModuleString *rms = (RedisModuleString*)malloc(sizeof(RedisModuleString));
    mock().setData("RedisModule_CreateString", mock().getData("RedisModule_CreateString").getIntValue()+1);
    return rms;
}

void RedisModule_FreeThreadSafeContext(RedisModuleCtx *ctx)
{
    (void)ctx;
    mock().setData("RedisModule_FreeThreadSafeContext", 1);
    return;
}

int RedisModule_StringToLongLong(const RedisModuleString *str, long long *ll)
{
    (void)str;

    int call_no = mock().getData("RedisModule_StringToLongLongCallCount").getIntValue();
    switch(call_no) {
        case 0:
            *ll = mock().getData("RedisModule_StringToLongLongCall_1").getIntValue();
            break;
        case 1:
            *ll = mock().getData("RedisModule_StringToLongLongCall_2").getIntValue();
            break;
        default:
            *ll = mock().getData("RedisModule_StringToLongLongCallDefault").getIntValue();
    }
    mock().setData("RedisModule_StringToLongLongCallCount", call_no + 1);
    return REDISMODULE_OK;
}

void RedisModule_ThreadSafeContextLock(RedisModuleCtx *ctx)
{
    (void)ctx;
    mock().setData("RedisModule_ThreadSafeContextLock", 1);
    return;
}

void RedisModule_ThreadSafeContextUnlock(RedisModuleCtx *ctx)
{
    (void)ctx;
    mock().setData("RedisModule_ThreadSafeContextUnlock", 1);
    return;
}

void RedisModule_ReplySetArrayLength(RedisModuleCtx *ctx, long len)
{
    (void)ctx;
    mock().setData("RedisModule_ReplySetArrayLength", (int)len);
    return;
}

RedisModuleCtx *RedisModule_GetThreadSafeContext(RedisModuleBlockedClient *bc)
{
    (void) bc;
    mock().setData("RedisModule_GetThreadSafeContext", 1);
    return NULL;
}

RedisModuleString *RedisModule_CreateStringFromLongLong(RedisModuleCtx *ctx, long long ll)
{
    (void)ctx;
    (void)ll;
    RedisModuleString *rms = (RedisModuleString*)malloc(sizeof(RedisModuleString));
    mock().setData("RedisModule_CreateStringFromLongLong", mock().getData("RedisModule_CreateStringFromLongLong").getIntValue()+1);
    return rms;
}

void RedisModule_AutoMemory(RedisModuleCtx *ctx)
{
    (void)ctx;
    int old = mock().getData("RedisModule_AutoMemory").getIntValue();
    mock().setData("RedisModule_AutoMemory", old + 1);
    return;
}

void *RedisModule_Alloc(size_t bytes)
{
    mock()
        .actualCall("RedisModule_Alloc");
    return malloc(bytes);
}

void RedisModule_Free(void *ptr)
{
    mock()
        .actualCall("RedisModule_Free");
    free(ptr);
}
