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



#include <stddef.h>
#include <stdbool.h>
#include <stdint.h>
#include <asn1codec_utils.h>

#ifndef INC_X2RESET_RESPONSE_WRAPPER_H
#define INC_X2RESET_RESPONSE_WRAPPER_H

#ifdef __cplusplus
extern "C"
{
#endif
bool
build_pack_x2reset_response(size_t* packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf);
#ifdef __cplusplus
}
#endif

#endif /* INC_X2RESET_RESPONSE_WRAPPER_H */

