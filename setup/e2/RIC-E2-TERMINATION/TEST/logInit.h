/*
 * Copyright 2019 AT&T Intellectual Property
 * Copyright 2019 Nokia
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * This source code is part of the near-RT RIC (RAN Intelligent Controller)
 * platform project (RICP).
 */


//
// Created by adi ENZEL on 6/11/19.
//

#ifndef __LOG_INIT__
#define __LOG_INIT__

#include <mdclog/mdclog.h>

#ifdef __cplusplus
extern "C"
{
#endif


void init_log(char *name) {
    mdclog_attr_t *attr;
    mdclog_attr_init(&attr);
    mdclog_attr_set_ident(attr, name);
    mdclog_init(attr);
    mdclog_attr_destroy(attr);
}

#ifdef __cplusplus
}
#endif
#endif
