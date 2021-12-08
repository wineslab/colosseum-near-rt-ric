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
// Created by adi ENZEL on 9/26/19.
//

#ifndef E2_BASE64_H
#define E2_BASE64_H

#include <mdclog/mdclog.h>
#include <cstring>
#include <zconf.h>

static const unsigned char base64_table[65] =
        "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

#define INVERSE_TABLE_SIZE 256

class base64 {
public:
    /**
     *
     * @param src
     * @param srcLen
     * @param dst
     * @param dstLen
     * @return 0 = OK -1 fault
     */
    static int encode(const unsigned char *src, int srcLen, char unsigned *dst, long &dstLen);
    /**
     *
     * @param src
     * @param srcLen
     * @param dst
     * @param dstLen
     * @return 0 = OK -1 fault
     */
    static int decode(const unsigned char *src, int srcLen, char unsigned *dst, long dstLen);

};


#endif //E2_BASE64_H
