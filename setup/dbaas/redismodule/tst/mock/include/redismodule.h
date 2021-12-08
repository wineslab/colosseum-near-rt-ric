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

#ifndef REDISMODULE_H
#define REDISMODULE_H

#include <sys/types.h>
#include <stdint.h>
#include <stdio.h>

/* Error status return values. */
#define REDISMODULE_OK 0
#define REDISMODULE_ERR 1

/* API versions. */
#define REDISMODULE_APIVER_1 1

/* API flags and constants */
#define REDISMODULE_READ (1<<0)
#define REDISMODULE_WRITE (1<<1)

/* Key types. */
#define REDISMODULE_KEYTYPE_EMPTY 0
#define REDISMODULE_KEYTYPE_STRING 1
#define REDISMODULE_KEYTYPE_LIST 2
#define REDISMODULE_KEYTYPE_HASH 3
#define REDISMODULE_KEYTYPE_SET 4
#define REDISMODULE_KEYTYPE_ZSET 5
#define REDISMODULE_KEYTYPE_MODULE 6

/* Reply types. */
#define REDISMODULE_REPLY_UNKNOWN -1
#define REDISMODULE_REPLY_STRING 0
#define REDISMODULE_REPLY_ERROR 1
#define REDISMODULE_REPLY_INTEGER 2
#define REDISMODULE_REPLY_ARRAY 3
#define REDISMODULE_REPLY_NULL 4

/* Postponed array length. */
#define REDISMODULE_POSTPONED_ARRAY_LEN -1

/* Error messages. */
#define REDISMODULE_ERRORMSG_WRONGTYPE "WRONGTYPE Operation against a key holding the wrong kind of value"

#define REDISMODULE_NOT_USED(V) ((void) V)

typedef long long mstime_t;

/* UT dummy definitions for opaque redis types */
typedef struct { int dummy; } RedisModuleCtx;
typedef struct { int dummy; } RedisModuleKey;
typedef struct { int dummy; } RedisModuleString;
typedef struct { int dummy; } RedisModuleCallReply;
typedef struct { int dummy; } RedisModuleIO;
typedef struct { int dummy; } RedisModuleType;
typedef struct { int dummy; } RedisModuleDigest;
typedef struct { int dummy; } RedisModuleBlockedClient;

typedef void *(*RedisModuleTypeLoadFunc)(RedisModuleIO *rdb, int encver);
typedef void (*RedisModuleTypeSaveFunc)(RedisModuleIO *rdb, void *value);
typedef void (*RedisModuleTypeRewriteFunc)(RedisModuleIO *aof, RedisModuleString *key, void *value);
typedef size_t (*RedisModuleTypeMemUsageFunc)(const void *value);
typedef void (*RedisModuleTypeDigestFunc)(RedisModuleDigest *digest, void *value);
typedef void (*RedisModuleTypeFreeFunc)(void *value);

typedef int (*RedisModuleCmdFunc) (RedisModuleCtx *ctx, RedisModuleString **argv, int argc);

int RedisModule_CreateCommand(RedisModuleCtx *ctx, const char *name, RedisModuleCmdFunc cmdfunc, const char *strflags, int firstkey, int lastkey, int keystep);
int RedisModule_WrongArity(RedisModuleCtx *ctx);
int RedisModule_ReplyWithLongLong(RedisModuleCtx *ctx, long long ll);
void *RedisModule_OpenKey(RedisModuleCtx *ctx, RedisModuleString *keyname, int mode);
RedisModuleCallReply *RedisModule_Call(RedisModuleCtx *ctx, const char *cmdname, const char *fmt, ...);
void RedisModule_FreeCallReply(RedisModuleCallReply *reply);
int RedisModule_CallReplyType(RedisModuleCallReply *reply);
long long RedisModule_CallReplyInteger(RedisModuleCallReply *reply);
const char *RedisModule_StringPtrLen(const RedisModuleString *str, size_t *len);
int RedisModule_ReplyWithError(RedisModuleCtx *ctx, const char *err);
int RedisModule_ReplyWithString(RedisModuleCtx *ctx, RedisModuleString *str);
int RedisModule_ReplyWithNull(RedisModuleCtx *ctx);
int RedisModule_ReplyWithCallReply(RedisModuleCtx *ctx, RedisModuleCallReply *reply);
const char *RedisModule_CallReplyStringPtr(RedisModuleCallReply *reply, size_t *len);
RedisModuleString *RedisModule_CreateStringFromCallReply(RedisModuleCallReply *reply);

int RedisModule_KeyType(RedisModuleKey *kp);
void RedisModule_CloseKey(RedisModuleKey *kp);

int RedisModule_Init(RedisModuleCtx *ctx, const char *name, int ver, int apiver);

size_t RedisModule_CallReplyLength(RedisModuleCallReply *reply);
RedisModuleCallReply *RedisModule_CallReplyArrayElement(RedisModuleCallReply *reply, size_t idx);
int RedisModule_ReplyWithArray(RedisModuleCtx *ctx, long len);
void RedisModule_FreeString(RedisModuleCtx *ctx, RedisModuleString *str);
RedisModuleBlockedClient *RedisModule_BlockClient(RedisModuleCtx *ctx, RedisModuleCmdFunc reply_callback, RedisModuleCmdFunc timeout_callback, void (*free_privdata)(RedisModuleCtx*,void*), long long timeout_ms);
int RedisModule_UnblockClient(RedisModuleBlockedClient *bc, void *privdata);
int RedisModule_AbortBlock(RedisModuleBlockedClient *bc);
RedisModuleString *RedisModule_CreateString(RedisModuleCtx *ctx, const char *ptr, size_t len);
void RedisModule_FreeThreadSafeContext(RedisModuleCtx *ctx);
int RedisModule_StringToLongLong(const RedisModuleString *str, long long *ll);
void RedisModule_ThreadSafeContextLock(RedisModuleCtx *ctx);
void RedisModule_ThreadSafeContextUnlock(RedisModuleCtx *ctx);
void RedisModule_ReplySetArrayLength(RedisModuleCtx *ctx, long len);
RedisModuleCtx *RedisModule_GetThreadSafeContext(RedisModuleBlockedClient *bc);
RedisModuleString *RedisModule_CreateStringFromLongLong(RedisModuleCtx *ctx, long long ll);
void RedisModule_AutoMemory(RedisModuleCtx *ctx);
void *RedisModule_Alloc(size_t bytes);
void RedisModule_Free(void *ptr);

#endif /* REDISMODULE_H */
