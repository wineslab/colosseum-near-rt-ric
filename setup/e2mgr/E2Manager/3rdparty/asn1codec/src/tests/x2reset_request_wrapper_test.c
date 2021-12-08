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



#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <x2reset_request_wrapper.h>

void test_build_pack_x2reset_request();
void test_unpack(void);

int
main(int argc, char* argv[])
{
    test_build_pack_x2reset_request();
    exit(0);
}

void test_build_pack_x2reset_request(){
    size_t error_buf_size = 8192;
    size_t packed_buf_size = 4096;
    unsigned char responseDataBuf[packed_buf_size];
    char responseErrorBuf[error_buf_size];
    bool result;
    E2AP_PDU_t *pdu;
    /**********************************************************************************/

    packed_buf_size = 4096;
    result = build_pack_x2reset_request(Cause_PR_radioNetwork,CauseRadioNetwork_time_critical_handover,
		    &packed_buf_size, responseDataBuf, error_buf_size, responseErrorBuf);
    if (!result) {
        printf("#%s failed. Packing error %s\n", __func__, responseErrorBuf);
        return;
    }
    printf("#%s packed size:%lu\nPayload:\n", __func__, packed_buf_size);
    for (size_t i = 0; i < packed_buf_size; ++i)
        printf("%02x",responseDataBuf[i]);
    printf("\n");

    pdu =calloc(1, sizeof(E2AP_PDU_t));
    if (!unpack_pdu_aux(pdu, packed_buf_size, responseDataBuf,error_buf_size, responseErrorBuf,ATS_ALIGNED_BASIC_PER)){
    	printf("#%s failed. Packing error %s\n", __func__, responseErrorBuf);
    }
    responseErrorBuf[0] = 0;
    asn1_pdu_printer(pdu, sizeof(responseErrorBuf), responseErrorBuf);
    printf("#%s: %s\n", __func__, responseErrorBuf);

}

