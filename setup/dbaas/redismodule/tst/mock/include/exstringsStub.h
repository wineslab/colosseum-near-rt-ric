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

#ifndef EXSTRINGSTUB_H_
#define EXSTRINGSTUB_H_


#include "redismodule.h"

int setStringGenericCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc, const int flag);
int SetIE_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int SetNE_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int RedisModule_OnLoad(RedisModuleCtx *ctx, RedisModuleString **argv, int argc) ;
int delStringGenericCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc, const int flag);
int DelIE_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int DelNE_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int SetPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int SetMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int SetIEPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int SetIEMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int SetNEPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int SetNXPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int SetNXMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int SetXXPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int DelPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int DelMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int DelIEPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int DelIEMPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int DelNEPub_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int NDel_Atomic_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int NGet_Atomic_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
int NGet_NoAtomic_RedisCommand(RedisModuleCtx *ctx, RedisModuleString **argv, int argc);
void *NGet_NoAtomic_ThreadMain(void *arg);

#endif
