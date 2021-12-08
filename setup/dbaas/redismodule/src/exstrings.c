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

#include "redismodule.h"
#include <pthread.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>

#ifdef __UT__
#include "exstringsStub.h"
#include "commonStub.h"
#endif


/* make sure the response is not NULL or an error.
sends the error to the client and exit the current function if its */
#define  ASSERT_NOERROR(r) \
    if (r == NULL) { \
        return RedisModule_ReplyWithError(ctx,"ERR reply is NULL"); \
    } else if (RedisModule_CallReplyType(r) == REDISMODULE_REPLY_ERROR) { \
        return RedisModule_ReplyWithCallReply(ctx,r); \
    }

#define OBJ_OP_NO 0
#define OBJ_OP_XX (1<<1)     /* OP if key exist */
#define OBJ_OP_NX (1<<2)     /* OP if key not exist */
#define OBJ_OP_IE (1<<4)     /* OP if equal old value */
#define OBJ_OP_NE (1<<5)     /* OP if not equal old value */

#define DEF_COUNT     50
#define ZERO          0
#define MATCH_STR     "MATCH"
#define COUNT_STR     "COUNT"
#define SCANARGC      5

RedisModuleString *def_count_str = NULL, *match_str = NULL, *count_str = NULL, *zero_str = NULL;

typedef struct _NgetArgs {
    RedisModuleString *key;
    RedisModuleString *count;
} NgetArgs;

typedef struct RedisModuleBlockedClientArgs {
    RedisModuleBlockedClient *bc;
    NgetArgs nget_args;
} RedisModuleBlockedClientArgs;

void InitStaticVariable()
{
    if (def_count_str == NULL)
        def_count_str = RedisModule_CreateStringFromLongLong(NULL, DEF_COUNT);
    if (match_str == NULL)
        match_str = RedisModule_CreateString(NULL, MATCH_STR, sizeof(MATCH_STR));
    if (count_str == NULL)
        count_str = RedisModule_CreateString(NULL, COUNT_STR, sizeof(COUNT_STR));
    if (zero_str == NULL)
        zero_str = RedisModule_CreateStringFromLongLong(NULL, ZERO);

    return;
}

int getKeyType(RedisModuleCtx *ctx, RedisModuleString *key_str)
{
    RedisModuleKey *key = RedisModule_OpenKey(ctx, key_str, REDISMODULE_READ);
    int type = RedisModule_KeyType(key);
    RedisModule_CloseKey(key);
    return type;
}

bool replyContentsEqualString(RedisModuleCallReply *reply, RedisModuleString *expected_value)
{
    size_t replylen = 0, expectedlen = 0;
    const char *expectedval = RedisModule_StringPtrLen(expected_value, &expectedlen);
    const char *replyval = RedisModule_CallReplyStringPtr(reply, &replylen);
    return replyval &&
           expectedlen == replylen &&
           !strncmp(expectedval, replyval, replylen);
}

typedef struct _SetParams {
    RedisModuleString **key_val_pairs;
    size_t length;
} SetParams;

typedef struct _PubParams {
    RedisModuleString **channel_msg_pairs;
    size_t length;
} PubParams;

typedef struct _DelParams {
    RedisModuleString **keys;
    size_t length;
} DelParams;

typedef enum _ExstringsStatus {
    EXSTRINGS_STATUS_NO_ERRORS = 0,
    EXSTRINGS_STATUS_ERROR_AND_REPLY_SENT,
    EXSTRINGS_STATUS_NOT_SET
} ExstringsStatus;

void readNgetArgs(RedisModuleCtx *ctx, RedisModuleString **argv, int argc,
                  NgetArgs* nget_args, ExstringsStatus* status)
{
    size_t str_len;
    long long number;

    if(argc == 2) {
        nget_args->key = argv[1];
        nget_args->count = def_count_str;
    } else if (argc == 4) {
        if (strcasecmp(RedisModule_StringPtrLen(argv[2], &str_len), "count")) {
            RedisModule_ReplyWithError(ctx,"-ERR syntax error");
            *status = EXSTRINGS_STATUS_ERROR_AND_REPLY_SENT;
            return;
        }

        int ret = RedisModule_StringToLongLong(argv[3], &number) != REDISMODULE_OK;
        if (ret != REDISMODULE_OK || number < 1) {
            RedisModule_ReplyWithError(ctx,"-ERR value is not an integer or out of range");
            *status = EXSTRINGS_STATUS_ERROR_AND_REPLY_SENT;
            return;
        }

        nget_args->key = argv[1];
        nget_args->count = argv[3];
    } else {
        /* In redis there is a bug (or undocumented feature see link)
         * where calling 'RedisModule_WrongArity'
         * within a blocked client will crash redis.
         *
         * Therefore we need to call this function to validate args
         * before putting the client into blocking mode.
         *
         * Link to issue:
         * https://github.com/antirez/redis/issues/6382
         * 'If any thread tries to access the command arguments from
         *  within the ThreadSafeContext they will crash redis' */
        RedisModule_WrongArity(ctx);
        *status = EXSTRINGS_STATUS_ERROR_AND_REPLY_SENT;
        return;
    }

    *status = EXSTRINGS_STATUS_NO_ERRORS;
    return;
}

long long callReplyLongLong(RedisModuleCallReply* reply)
{
    const char* cursor_str_ptr = RedisModule_CallReplyStringPtr(reply, NULL);
    return strtoll(cursor_str_ptr, NULL, 10);
}

void forwardIfError(RedisModuleCtx *ctx, RedisModuleCallReply *reply, ExstringsStatus* status)
{
    if (RedisModule_CallReplyType(reply) == REDISMODULE_REPLY_ERROR) {
        RedisModule_ReplyWithCallReply(ctx, reply);
        RedisModule_FreeCallReply(reply);
        *status = EXSTRINGS_STATUS_ERROR_AND_REPLY_SENT;
    }
    *status = EXSTRINGS_STATUS_NO_ERRORS;
}

typedef struct _ScannedKeys {
    RedisModuleString **keys;
    size_t len;
} ScannedKeys;

ScannedKeys* allocScannedKeys(size_t len)
{
    ScannedKeys *sk = RedisModule_Alloc(sizeof(ScannedKeys));
    if (sk) {
        sk->len = len;
        sk->keys = RedisModule_Alloc(sizeof(RedisModuleString *)*len);
    }
    return sk;
}

void freeScannedKeys(RedisModuleCtx *ctx, ScannedKeys* sk)
{
    if (sk) {
        size_t j;
        for (j = 0; j < sk->len; j++)
            RedisModule_FreeString(ctx, sk->keys[j]);
        RedisModule_Free(sk->keys);
    }
    RedisModule_Free(sk);
}

typedef struct _ScanSomeState {
    RedisModuleString *key;
    RedisModuleString *count;
    long long cursor;
} ScanSomeState;

ScannedKeys *scanSome(RedisModuleCtx* ctx, ScanSomeState* state, ExstringsStatus* status)
{
    RedisModuleString *scanargv[SCANARGC] = {NULL};

    scanargv[0] = RedisModule_CreateStringFromLongLong(ctx, state->cursor);
    scanargv[1] = match_str;
    scanargv[2] = state->key;
    scanargv[3] = count_str;
    scanargv[4] = state->count;

    RedisModuleCallReply *reply;
    reply = RedisModule_Call(ctx, "SCAN", "v", scanargv, SCANARGC);
    RedisModule_FreeString(ctx, scanargv[0]);
    forwardIfError(ctx, reply, status);
    if (*status == EXSTRINGS_STATUS_ERROR_AND_REPLY_SENT)
        return NULL;

    state->cursor = callReplyLongLong(RedisModule_CallReplyArrayElement(reply, 0));
    RedisModuleCallReply *cr_keys =
        RedisModule_CallReplyArrayElement(reply, 1);

    size_t scanned_keys_len = RedisModule_CallReplyLength(cr_keys);
    if (scanned_keys_len == 0) {
        RedisModule_FreeCallReply(reply);
        *status = EXSTRINGS_STATUS_NO_ERRORS;
        return NULL;
    }

    ScannedKeys *scanned_keys = allocScannedKeys(scanned_keys_len);
    if (scanned_keys == NULL) {
        RedisModule_FreeCallReply(reply);
        RedisModule_ReplyWithError(ctx,"-ERR Out of memory");
        *status = EXSTRINGS_STATUS_ERROR_AND_REPLY_SENT;
        return NULL;
    }

    scanned_keys->len = scanned_keys_len;
    size_t j;
    for (j = 0; j < scanned_keys_len; j++) {
        RedisModuleString *rms = RedisModule_CreateStringFromCallReply(RedisModule_CallReplyArrayElement(cr_keys,j));
        scanned_keys->keys[j] = rms;
    }
    RedisModule_FreeCallReply(reply);
    *status = EXSTRINGS_STATUS_NO_ERRORS;
    return scanned_keys;
}

inline void unlockThreadsafeContext(RedisModuleCtx *ctx, bool using_threadsafe_context)
{
    if (using_threadsafe_context)
        RedisModule_ThreadSafeContextUnlock(ctx);
}

inline void lockThreadsafeContext(RedisModuleCtx *ctx, bool using_threadsafe_context)
{
    if (using_threadsafe_context)
        RedisModule_ThreadSafeContextLock(ctx);
}

void multiPubCommand(RedisModuleCtx *ctx, PubParams* pubParams)
{
    RedisModuleCallReply *reply = NULL;
    for (unsigned int i = 0 ; i < pubParams->length ; i += 2) {
        reply = RedisModule_Call(ctx, "PUBLISH", "v", pubParams->channel_msg_pairs + i, 2);
        RedisModule_FreeCallReply(reply);
    }
}

int setStringGenericCommand(RedisModuleCtx *ctx, RedisModuleString **argv,
                                       int argc, const int flag)
{
    RedisModuleString *oldvalstr = NULL;
    RedisModuleCallReply *reply = NULL;

    if (argc < 4)
        return RedisModule_WrongArity(ctx);
    else
        oldvalstr = argv[3];

    /*Check if key type is string*/
    RedisModuleKey *key = RedisModule_OpenKey(ctx,argv[1],
        REDISMODULE_READ);
    int type = RedisModule_KeyType(key);
    RedisModule_CloseKey(key);

    if (type == REDISMODULE_KEYTYPE_EMPTY) {
        if (flag == OBJ_OP_IE){
            RedisModule_ReplyWithNull(ctx);
            return REDISMODULE_OK;
        }
    } else if (type != REDISMODULE_KEYTYPE_STRING) {
        return RedisModule_ReplyWithError(ctx,REDISMODULE_ERRORMSG_WRONGTYPE);
    }

    /*Get the value*/
    reply = RedisModule_Call(ctx, "GET", "s", argv[1]);
    ASSERT_NOERROR(reply)
    size_t curlen=0, oldvallen=0;
    const char *oldval = RedisModule_StringPtrLen(oldvalstr, &oldvallen);
    const char *curval = RedisModule_CallReplyStringPtr(reply, &curlen);
    if (((flag == OBJ_OP_IE) &&
        (!curval || (oldvallen != curlen) || strncmp(oldval, curval, curlen)))
        ||
        ((flag == OBJ_OP_NE) && curval && (oldvallen == curlen) &&
          !strncmp(oldval, curval, curlen))) {
        RedisModule_FreeCallReply(reply);
        return RedisModule_ReplyWithNull(ctx);
    }
    RedisModule_FreeCallReply(reply);

    /* Prepare the arguments for the command. */
    int i, j=0, cmdargc=argc-2;
    RedisModuleString *cmdargv[cmdargc];
    for (i = 1; i < argc; i++) {
        if (i == 3)
            continue;
        cmdargv[j++] = argv[i];
    }

    /* Call the command and pass back the reply. */
    reply = RedisModule_Call(ctx, "SET", "v!", cmdargv, cmdargc);
    ASSERT_NOERROR(reply)
    RedisModule_ReplyWithCallReply(ctx, reply);

    RedisModule_FreeCallReply(reply);
    return REDISMODULE_OK;
}

int SetIE_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    RedisModule_AutoMemory(ctx);
    return setStringGenericCommand(ctx, argv, argc, OBJ_OP_IE);
}

int SetNE_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    RedisModule_AutoMemory(ctx);
    return setStringGenericCommand(ctx, argv, argc, OBJ_OP_NE);
}

int delStringGenericCommand(RedisModuleCtx *ctx, RedisModuleString **argv,
                                       int argc, const int flag)
{
    RedisModuleString *oldvalstr = NULL;
    RedisModuleCallReply *reply = NULL;

    if (argc == 3)
        oldvalstr = argv[2];
    else
        return RedisModule_WrongArity(ctx);

    /*Check if key type is string*/
    RedisModuleKey *key = RedisModule_OpenKey(ctx,argv[1],
        REDISMODULE_READ);
    int type = RedisModule_KeyType(key);
    RedisModule_CloseKey(key);

    if (type == REDISMODULE_KEYTYPE_EMPTY) {
        return RedisModule_ReplyWithLongLong(ctx, 0);
    } else if (type != REDISMODULE_KEYTYPE_STRING) {
        return RedisModule_ReplyWithError(ctx,REDISMODULE_ERRORMSG_WRONGTYPE);
    }

    /*Get the value*/
    reply = RedisModule_Call(ctx, "GET", "s", argv[1]);
    ASSERT_NOERROR(reply)
    size_t curlen = 0, oldvallen = 0;
    const char *oldval = RedisModule_StringPtrLen(oldvalstr, &oldvallen);
    const char *curval = RedisModule_CallReplyStringPtr(reply, &curlen);
    if (((flag == OBJ_OP_IE) &&
        (!curval || (oldvallen != curlen) || strncmp(oldval, curval, curlen)))
        ||
        ((flag == OBJ_OP_NE) && curval && (oldvallen == curlen) &&
          !strncmp(oldval, curval, curlen))) {
        RedisModule_FreeCallReply(reply);
        return RedisModule_ReplyWithLongLong(ctx, 0);
    }
    RedisModule_FreeCallReply(reply);

    /* Prepare the arguments for the command. */
    int cmdargc=1;
    RedisModuleString *cmdargv[1];
    cmdargv[0] = argv[1];

    /* Call the command and pass back the reply. */
    reply = RedisModule_Call(ctx, "UNLINK", "v!", cmdargv, cmdargc);
    ASSERT_NOERROR(reply)
    RedisModule_ReplyWithCallReply(ctx, reply);

    RedisModule_FreeCallReply(reply);
    return REDISMODULE_OK;
}

int DelIE_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    RedisModule_AutoMemory(ctx);
    return delStringGenericCommand(ctx, argv, argc, OBJ_OP_IE);
}

int DelNE_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    RedisModule_AutoMemory(ctx);
    return delStringGenericCommand(ctx, argv, argc, OBJ_OP_NE);
}
int setPubStringCommon(RedisModuleCtx *ctx, SetParams* setParamsPtr, PubParams* pubParamsPtr)
{
    RedisModuleCallReply *setReply;
    setReply = RedisModule_Call(ctx, "MSET", "v!", setParamsPtr->key_val_pairs, setParamsPtr->length);
    ASSERT_NOERROR(setReply)
    multiPubCommand(ctx, pubParamsPtr);
    RedisModule_ReplyWithCallReply(ctx, setReply);
    RedisModule_FreeCallReply(setReply);
    return REDISMODULE_OK;
}

int SetPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc < 5 || (argc % 2) == 0)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    SetParams setParams = {
                           .key_val_pairs = argv + 1,
                           .length = argc - 3
                          };
    PubParams pubParams = {
                           .channel_msg_pairs = argv + 1 + setParams.length,
                           .length = 2
                          };

    return setPubStringCommon(ctx, &setParams, &pubParams);
}

int SetMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc < 7 || (argc % 2) == 0)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    long long setPairsCount, pubPairsCount;
    RedisModule_StringToLongLong(argv[1], &setPairsCount);
    RedisModule_StringToLongLong(argv[2], &pubPairsCount);
    if (setPairsCount < 1 || pubPairsCount < 1)
        return RedisModule_ReplyWithError(ctx, "ERR SET_PAIR_COUNT and PUB_PAIR_COUNT must be greater than zero");

    long long setLen, pubLen;
    setLen = 2*setPairsCount;
    pubLen = 2*pubPairsCount;

    if (setLen + pubLen + 3 != argc)
        return RedisModule_ReplyWithError(ctx, "ERR SET_PAIR_COUNT or PUB_PAIR_COUNT do not match the total pair count");

    SetParams setParams = {
                           .key_val_pairs = argv + 3,
                           .length = setLen
                          };
    PubParams pubParams = {
                           .channel_msg_pairs = argv + 3 + setParams.length,
                           .length = pubLen
                          };

    return setPubStringCommon(ctx, &setParams, &pubParams);
}

int setIENEPubStringCommon(RedisModuleCtx *ctx, RedisModuleString **argv, int argc, int flag)
{
    SetParams setParams = {
                           .key_val_pairs = argv + 1,
                           .length = 2
                          };
    PubParams pubParams = {
                           .channel_msg_pairs = argv + 4,
                           .length = argc - 4
                          };
    RedisModuleString *key = setParams.key_val_pairs[0];
    RedisModuleString *oldvalstr = argv[3];

    int type = getKeyType(ctx, key);
    if (flag == OBJ_OP_IE && type == REDISMODULE_KEYTYPE_EMPTY) {
        return RedisModule_ReplyWithNull(ctx);
    } else if (type != REDISMODULE_KEYTYPE_STRING && type != REDISMODULE_KEYTYPE_EMPTY) {
        return RedisModule_ReplyWithError(ctx, REDISMODULE_ERRORMSG_WRONGTYPE);
    }

    RedisModuleCallReply *reply = RedisModule_Call(ctx, "GET", "s", key);
    ASSERT_NOERROR(reply)
    bool is_equal = replyContentsEqualString(reply, oldvalstr);
    RedisModule_FreeCallReply(reply);
    if ((flag == OBJ_OP_IE && !is_equal) ||
        (flag == OBJ_OP_NE && is_equal)) {
        return RedisModule_ReplyWithNull(ctx);
    }

    return setPubStringCommon(ctx, &setParams, &pubParams);
}

int SetIEPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc != 6)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    return setIENEPubStringCommon(ctx, argv, argc, OBJ_OP_IE);
}

int SetIEMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc < 6 || (argc % 2) != 0)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    return setIENEPubStringCommon(ctx, argv, argc, OBJ_OP_IE);
}

int SetNEPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc != 6)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    return setIENEPubStringCommon(ctx, argv, argc, OBJ_OP_NE);
}

int setXXNXPubStringCommon(RedisModuleCtx *ctx, RedisModuleString **argv, int argc, int flag)
{
    SetParams setParams = {
                           .key_val_pairs = argv + 1,
                           .length = 2
                          };
    PubParams pubParams = {
                           .channel_msg_pairs = argv + 3,
                           .length = argc - 3
                          };
    RedisModuleString *key = setParams.key_val_pairs[0];

    int type = getKeyType(ctx, key);
    if ((flag == OBJ_OP_XX && type == REDISMODULE_KEYTYPE_EMPTY) ||
        (flag == OBJ_OP_NX && type == REDISMODULE_KEYTYPE_STRING)) {
        return RedisModule_ReplyWithNull(ctx);
    } else if (type != REDISMODULE_KEYTYPE_STRING && type != REDISMODULE_KEYTYPE_EMPTY) {
        RedisModule_ReplyWithError(ctx, REDISMODULE_ERRORMSG_WRONGTYPE);
        return REDISMODULE_OK;
    }

    return setPubStringCommon(ctx, &setParams, &pubParams);
}

int SetNXPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc != 5)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    return setXXNXPubStringCommon(ctx, argv, argc, OBJ_OP_NX);
}

int SetNXMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc < 5 || (argc % 2) == 0)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    return setXXNXPubStringCommon(ctx, argv, argc, OBJ_OP_NX);
}

int SetXXPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc != 5)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    return setXXNXPubStringCommon(ctx, argv, argc, OBJ_OP_XX);
}

int delPubStringCommon(RedisModuleCtx *ctx, DelParams *delParamsPtr, PubParams *pubParamsPtr)
{
    RedisModuleCallReply *reply = RedisModule_Call(ctx, "UNLINK", "v!", delParamsPtr->keys, delParamsPtr->length);
    ASSERT_NOERROR(reply)
    int replytype = RedisModule_CallReplyType(reply);
    if (replytype == REDISMODULE_REPLY_NULL) {
        RedisModule_ReplyWithNull(ctx);
    } else if (RedisModule_CallReplyInteger(reply) == 0) {
        RedisModule_ReplyWithCallReply(ctx, reply);
    } else {
        RedisModule_ReplyWithCallReply(ctx, reply);
        multiPubCommand(ctx, pubParamsPtr);
    }
    RedisModule_FreeCallReply(reply);
    return REDISMODULE_OK;
}

int DelPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc < 4)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    DelParams delParams = {
                           .keys = argv + 1,
                           .length = argc - 3
                          };
    PubParams pubParams = {
                           .channel_msg_pairs = argv + 1 + delParams.length,
                           .length = 2
                          };

    return delPubStringCommon(ctx, &delParams, &pubParams);
}

int DelMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc < 6)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    long long delCount, pubPairsCount;
    RedisModule_StringToLongLong(argv[1], &delCount);
    RedisModule_StringToLongLong(argv[2], &pubPairsCount);
    if (delCount < 1 || pubPairsCount < 1)
        return RedisModule_ReplyWithError(ctx, "ERR DEL_COUNT and PUB_PAIR_COUNT must be greater than zero");

    long long delLen, pubLen;
    delLen = delCount;
    pubLen = 2*pubPairsCount;
    if (delLen + pubLen + 3 != argc)
        return RedisModule_ReplyWithError(ctx, "ERR DEL_COUNT or PUB_PAIR_COUNT do not match the total pair count");

    DelParams delParams = {
                           .keys = argv + 3,
                           .length = delLen
                          };
    PubParams pubParams = {
                           .channel_msg_pairs = argv + 3 + delParams.length,
                           .length = pubLen
                          };

    return delPubStringCommon(ctx, &delParams, &pubParams);
}

int delIENEPubStringCommon(RedisModuleCtx *ctx, RedisModuleString **argv, int argc, int flag)
{
    DelParams delParams = {
                           .keys = argv + 1,
                           .length = 1
                          };
    PubParams pubParams = {
                           .channel_msg_pairs = argv + 3,
                           .length = argc - 3
                          };
    RedisModuleString *key = argv[1];
    RedisModuleString *oldvalstr = argv[2];

    int type = getKeyType(ctx, key);
    if (type == REDISMODULE_KEYTYPE_EMPTY) {
        return RedisModule_ReplyWithLongLong(ctx, 0);
    } else if (type != REDISMODULE_KEYTYPE_STRING) {
        return RedisModule_ReplyWithError(ctx, REDISMODULE_ERRORMSG_WRONGTYPE);
    }

    RedisModuleCallReply *reply = RedisModule_Call(ctx, "GET", "s", key);
    ASSERT_NOERROR(reply)
    bool is_equal = replyContentsEqualString(reply, oldvalstr);
    RedisModule_FreeCallReply(reply);
    if ((flag == OBJ_OP_IE && !is_equal) ||
        (flag == OBJ_OP_NE && is_equal)) {
        return RedisModule_ReplyWithLongLong(ctx, 0);
    }

    return delPubStringCommon(ctx, &delParams, &pubParams);
}

int DelIEPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc != 5)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    return delIENEPubStringCommon(ctx, argv, argc, OBJ_OP_IE);
}

int DelIEMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc < 5 || (argc % 2) == 0)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    return delIENEPubStringCommon(ctx, argv, argc, OBJ_OP_IE);
}

int DelNEPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    if (argc != 5)
        return RedisModule_WrongArity(ctx);

    RedisModule_AutoMemory(ctx);
    return delIENEPubStringCommon(ctx, argv, argc, OBJ_OP_NE);
}

int Nget_RedisCommand(RedisModuleCtx *ctx, NgetArgs* nget_args, bool using_threadsafe_context)
{
    int ret = REDISMODULE_OK;
    size_t replylen = 0;
    RedisModuleCallReply *reply = NULL;
    ExstringsStatus status = EXSTRINGS_STATUS_NOT_SET;
    ScanSomeState scan_state;
    ScannedKeys *scanned_keys;

    scan_state.key = nget_args->key;
    scan_state.count = nget_args->count;
    scan_state.cursor = 0;

    RedisModule_ReplyWithArray(ctx, REDISMODULE_POSTPONED_ARRAY_LEN);
    do {
        lockThreadsafeContext(ctx, using_threadsafe_context);

        status = EXSTRINGS_STATUS_NOT_SET;
        scanned_keys = scanSome(ctx, &scan_state, &status);

        if (status != EXSTRINGS_STATUS_NO_ERRORS) {
            unlockThreadsafeContext(ctx, using_threadsafe_context);
            ret = REDISMODULE_ERR;
            break;
        } else if (scanned_keys == NULL) {
            unlockThreadsafeContext(ctx, using_threadsafe_context);
            continue;
        }

        reply = RedisModule_Call(ctx, "MGET", "v", scanned_keys->keys, scanned_keys->len);

        unlockThreadsafeContext(ctx, using_threadsafe_context);

        status = EXSTRINGS_STATUS_NOT_SET;
        forwardIfError(ctx, reply, &status);
        if (status != EXSTRINGS_STATUS_NO_ERRORS) {
            freeScannedKeys(ctx, scanned_keys);
            ret = REDISMODULE_ERR;
            break;
        }

        size_t i;
        for (i = 0; i < scanned_keys->len; i++) {
            RedisModuleString *rms = RedisModule_CreateStringFromCallReply(RedisModule_CallReplyArrayElement(reply, i));
            if (rms) {
                RedisModule_ReplyWithString(ctx, scanned_keys->keys[i]);
                RedisModule_ReplyWithString(ctx, rms);
                RedisModule_FreeString(ctx, rms);
                replylen += 2;
            }
        }
        RedisModule_FreeCallReply(reply);
        freeScannedKeys(ctx, scanned_keys);
    } while (scan_state.cursor != 0);

    RedisModule_ReplySetArrayLength(ctx,replylen);
    return ret;
}

/* The thread entry point that actually executes the blocking part
 * of the command nget.noatomic
 */
void *NGet_NoAtomic_ThreadMain(void *arg)
{
    pthread_detach(pthread_self());

    RedisModuleBlockedClientArgs *bca = arg;
    RedisModuleBlockedClient *bc = bca->bc;
    RedisModuleCtx *ctx = RedisModule_GetThreadSafeContext(bc);

    Nget_RedisCommand(ctx, &bca->nget_args, true);
    RedisModule_FreeThreadSafeContext(ctx);
    RedisModule_UnblockClient(bc, NULL);
    RedisModule_Free(bca);
    return NULL;
}

int NGet_NoAtomic_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    RedisModule_AutoMemory(ctx);
    pthread_t tid;

    InitStaticVariable();

    RedisModuleBlockedClientArgs *bca = RedisModule_Alloc(sizeof(RedisModuleBlockedClientArgs));
    if (bca == NULL) {
        RedisModule_ReplyWithError(ctx,"-ERR Out of memory");
        return REDISMODULE_ERR;
    }

    ExstringsStatus status = EXSTRINGS_STATUS_NOT_SET;
    readNgetArgs(ctx, argv, argc, &bca->nget_args, &status);
    if (status != EXSTRINGS_STATUS_NO_ERRORS) {
        RedisModule_Free(bca);
        return REDISMODULE_ERR;
    }

    /* Note that when blocking the client we do not set any callback: no
     * timeout is possible since we passed '0', nor we need a reply callback
     * because we'll use the thread safe context to accumulate a reply. */
    RedisModuleBlockedClient *bc = RedisModule_BlockClient(ctx,NULL,NULL,NULL,0);

    bca->bc = bc;

    /* Now that we setup a blocking client, we need to pass the control
     * to the thread. However we need to pass arguments to the thread:
     * the reference to the blocked client handle. */
    if (pthread_create(&tid,NULL,NGet_NoAtomic_ThreadMain,bca) != 0) {
        RedisModule_AbortBlock(bc);
        RedisModule_Free(bca);
        return RedisModule_ReplyWithError(ctx,"-ERR Can't start thread");
    }

    return REDISMODULE_OK;
}

int NGet_Atomic_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    RedisModule_AutoMemory(ctx);
    NgetArgs nget_args;
    ExstringsStatus status = EXSTRINGS_STATUS_NOT_SET;

    InitStaticVariable();

    readNgetArgs(ctx, argv, argc, &nget_args, &status);
    if (status != EXSTRINGS_STATUS_NO_ERRORS) {
        return REDISMODULE_ERR;
    }

    return Nget_RedisCommand(ctx, &nget_args, false);
}

int NDel_Atomic_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc)
{
    RedisModule_AutoMemory(ctx);
    int ret = REDISMODULE_OK;
    long long replylen = 0;
    RedisModuleCallReply *reply = NULL;
    ExstringsStatus status = EXSTRINGS_STATUS_NOT_SET;
    ScanSomeState scan_state;
    ScannedKeys *scanned_keys = NULL;

    InitStaticVariable();
    if (argc != 2)
        return RedisModule_WrongArity(ctx);

    scan_state.key = argv[1];
    scan_state.count = def_count_str;
    scan_state.cursor = 0;

    do {
        status = EXSTRINGS_STATUS_NOT_SET;
        scanned_keys = scanSome(ctx, &scan_state, &status);

        if (status != EXSTRINGS_STATUS_NO_ERRORS) {
            ret = REDISMODULE_ERR;
            break;
        } else if (scanned_keys == NULL) {
            continue;
        }

        reply = RedisModule_Call(ctx, "UNLINK", "v!", scanned_keys->keys, scanned_keys->len);

        status = EXSTRINGS_STATUS_NOT_SET;
        forwardIfError(ctx, reply, &status);
        if (status != EXSTRINGS_STATUS_NO_ERRORS) {
            freeScannedKeys(ctx, scanned_keys);
            ret = REDISMODULE_ERR;
            break;
        }

        replylen += RedisModule_CallReplyInteger(reply);
        RedisModule_FreeCallReply(reply);
        freeScannedKeys(ctx, scanned_keys);
    } while (scan_state.cursor != 0);

    if (ret == REDISMODULE_OK) {
        RedisModule_ReplyWithLongLong(ctx, replylen);
    }

    return ret;
}

/* This function must be present on each Redis module. It is used in order to
 * register the commands into the Redis server. */
int RedisModule_OnLoad(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) {
    REDISMODULE_NOT_USED(argv);
    REDISMODULE_NOT_USED(argc);

    if (RedisModule_Init(ctx,"exstrings",1,REDISMODULE_APIVER_1)
        == REDISMODULE_ERR) return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"setie",
        SetIE_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"setne",
        SetNE_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"delie",
        DelIE_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"delne",
        DelNE_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"nget.atomic",
        NGet_Atomic_RedisCommand,"readonly",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"nget.noatomic",
        NGet_NoAtomic_RedisCommand,"readonly",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"ndel.atomic",
        NDel_Atomic_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"msetpub",
        SetPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"msetmpub",
        SetMPub_RedisCommand,"write deny-oom pubsub",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"setiepub",
        SetIEPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"setiempub",
        SetIEMPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"setnepub",
        SetNEPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"setxxpub",
        SetXXPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"setnxpub",
        SetNXPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"setnxmpub",
        SetNXMPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"delpub",
        DelPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"delmpub",
        DelMPub_RedisCommand,"write deny-oom pubsub",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"deliepub",
        DelIEPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"deliempub",
        DelIEMPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    if (RedisModule_CreateCommand(ctx,"delnepub",
        DelNEPub_RedisCommand,"write deny-oom",1,1,1) == REDISMODULE_ERR)
        return REDISMODULE_ERR;

    return REDISMODULE_OK;
}
