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

#include <stdlib.h>

#include "redismodule.h"
#include "ut_helpers.hpp"

#include <CppUTest/TestHarness.h>
#include <CppUTestExt/MockSupport.h>

RedisModuleString **createRedisStrVec(size_t size)
{
    RedisModuleString ** redisStrVec = new RedisModuleString*[size];
    for (size_t i = 0 ; i < size ; i++) {
        redisStrVec[i] = (RedisModuleString *)UT_DUMMY_PTR_ADDRESS;
    }
    return redisStrVec;
}

void returnNKeysFromScanSome(long keys)
{
    mock().expectOneCall("RedisModule_CallReplyLength")
          .andReturnValue((int)keys);
    for (long i = 0 ; i < keys ; i++) {
        mock().expectOneCall("RedisModule_CreateStringFromCallReply")
              .andReturnValue(malloc(UT_DUMMY_BUFFER_SIZE));
    }
}

