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
#include <stdio.h>
#include <stdlib.h>
#include <ctype.h>

#include "redismodule.h"
#include "commonStub.h"
}

#include <CppUTest/TestHarness.h>
#include <CppUTestExt/MockSupport.h>
#include <CppUTest/MemoryLeakDetectorMallocMacros.h>

typedef struct RedisModuleBlockedClientArgs {
    RedisModuleBlockedClient *bc;
    RedisModuleString **argv;
    int argc;
} RedisModuleBlockedClientArgs;

int pthread_create(pthread_t *thread, const pthread_attr_t *attr,
                   void *(*start_routine) (void *), void *arg)
{
    (void)thread;
    (void)attr;
    (void)start_routine;
    if (mock().getData("pthread_create_free_block_client_args").getIntValue()) {
        RedisModuleBlockedClientArgs* bca = (RedisModuleBlockedClientArgs*)arg;
        free(bca->bc);
        free(bca);
    }

    return mock()
        .actualCall("pthread_create")
        .returnIntValueOrDefault(0);
}

int pthread_detach(pthread_t thread)
{
    (void)thread;

    return mock()
        .actualCall("pthread_detach")
        .returnIntValueOrDefault(0);
}

pthread_t pthread_self(void)
{
    return mock()
        .actualCall("pthread_self")
        .returnIntValueOrDefault(UT_DUMMY_THREAD_ID);
}
