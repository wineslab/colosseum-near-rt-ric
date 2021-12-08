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

#include "base64.h"

int base64::encode(const unsigned char *src, int srcLen, char unsigned *dst, long &dstLen) {
    unsigned char *pos;
    const unsigned char *end, *in;

    if (dstLen <= 0 || srcLen <= 0) {
        mdclog_write(MDCLOG_ERR, "source or destination length are 0. dst =%ld source = %d",
                     dstLen, srcLen);
        return -1;
    }
    if (dstLen < (srcLen * 4 / 3)) {
        mdclog_write(MDCLOG_ERR, "Destination size %ld must be at least 140 percent from source size %d",
                     dstLen, srcLen);
        return -1;
    }
    if (dst == nullptr) {
        mdclog_write(MDCLOG_ERR, "Destination must be allocated and freed by caller the function not allocate the memory");
        return -1;
    }
    if (src == nullptr) {
        mdclog_write(MDCLOG_ERR, "source is null pointer");
        return -1;
    }

    end = src + srcLen;
    in = src;
    pos = dst;
    while (end - in >= 3) {
        *pos++ = base64_table[in[0] >> (unsigned int)2];
        *pos++ = base64_table[((in[0] & 0x03) << (unsigned int)4) | (in[1] >> (unsigned int)4)];
        *pos++ = base64_table[((in[1] & 0x0f) << (unsigned int)2) | (in[2] >> (unsigned int)6)];
        *pos++ = base64_table[in[2] & (unsigned int)0x3f];
        in += 3;
    }

    if (end - in) {
        *pos++ = base64_table[in[0] >> (unsigned int)2];
        if (end - in == 1) {
            *pos++ = base64_table[(in[0] & 0x03) << (unsigned int)4];
            *pos++ = '=';
        } else {
            *pos++ = base64_table[((in[0] & 0x03) << 4) | (in[1] >> (unsigned int)4)];
            *pos++ = base64_table[(in[1] & 0x0f) << 2];
        }
        *pos++ = '=';
    }

    *pos = '\0';
    dstLen = pos - dst;return 0;
}

int base64::decode(const unsigned char *src, int srcLen, char unsigned *dst, long dstLen) {
    unsigned char inv_table[INVERSE_TABLE_SIZE];
    memset(inv_table, 0x80, INVERSE_TABLE_SIZE);
    for (ulong i = 0; i < sizeof(base64_table) - 1; i++) {
        inv_table[base64_table[i]] = (unsigned char) i;
    }
    inv_table['='] = 0;


    if (dstLen == 0 || dstLen  < (int)(srcLen / 4 * 3)) {
        mdclog_write(MDCLOG_ERR, "Destination size %ld can be up to 40  smaller then source size %d",
                     dstLen, srcLen);
        return -1;
    }
    if (dst == nullptr) {
        mdclog_write(MDCLOG_ERR, "Destination must be allocated and freed by caller the function not allocate the memory");
        return -1;
    }

    unsigned char *pos, block[4], tmp;
    long i;
    int pad = 0;

    size_t count = 0;

    for (i = 0; i < srcLen; i++) {
        if (inv_table[src[i]] != 0x80) {
            count++;
        }
    }

    if (count == 0 || count % 4)
        return -1;

    pos = dst;
    count = 0;
    for (i = 0; i < srcLen; i++) {
        tmp = inv_table[src[i]];
        if (tmp == 0x80) {
            continue;
        }
        block[count] = tmp;

        if (src[i] == '=') {
            pad++;
        }

        count++;
        if (count == 4) {
            *pos++ = (block[0] << 2) | ((unsigned char)block[1] >> (unsigned int)4);
            *pos++ = (block[1] << 4) | ((unsigned char)block[2] >> (unsigned int)2);
            *pos++ = (block[2] << 6) | block[3];
            count = 0;
            if (pad) {
                if (pad == 1) {
                    pos--;
                }
                else if (pad == 2) {
                    pos -= 2;
                }
                else {
                    return -1;
                }
                break;
            }
        }
    }

    dstLen = pos - dst;
    return 0;
}
