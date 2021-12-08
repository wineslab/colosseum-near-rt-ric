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
#include <x2setup_request_wrapper.h>

void test_build_pack_x2setup_request();
void test_build_pack_endc_x2setup_request();
void test_unpack(void);

int
main(int argc, char* argv[])
{
    test_build_pack_x2setup_request();
    test_build_pack_endc_x2setup_request();
    test_unpack();
    exit(0);
}

void test_build_pack_x2setup_request(){
    size_t error_buf_size = 8192;
    size_t packed_buf_size = 4096;
    unsigned char responseDataBuf[packed_buf_size];
    char responseErrorBuf[error_buf_size];
    uint8_t pLMN_Identity[] = {0xa,0xb,0xc}; 
    uint8_t ric_flag[] = {0xa,0xd,0xe};
    uint8_t eNBId[] = {0xab, 0xcd, 0x7/*0xf,0x7,0x2*/};
    bool result;
    E2AP_PDU_t *pdu;
    unsigned int bitqty = 21;
    /**********************************************************************************/

    printf("\n----- ATS_ALIGNED_BASIC_PER ----\n");
    packed_buf_size = 4096;
    result = build_pack_x2setup_request_aux(
		    pLMN_Identity, eNBId, bitqty , ric_flag,
		    &packed_buf_size, responseDataBuf, error_buf_size, responseErrorBuf,ATS_ALIGNED_BASIC_PER);
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
    printf("#%s: 21%s\n", __func__, responseErrorBuf);

    printf("\n----- ATS_UNALIGNED_BASIC_PER ----\n");
    packed_buf_size = 4096;
    result = build_pack_x2setup_request_aux(
		    pLMN_Identity, eNBId, bitqty , ric_flag,
		    &packed_buf_size, responseDataBuf, error_buf_size, responseErrorBuf,ATS_UNALIGNED_BASIC_PER);
    if (!result) {
        printf("#%s failed. Packing error %s\n", __func__, responseErrorBuf);
        return;
    }
    printf("#%s packed size:%lu\nPayload:\n", __func__, packed_buf_size);
    for (size_t i = 0; i < packed_buf_size; ++i)
        printf("%02x",responseDataBuf[i]);
    printf("\n");

    pdu =calloc(1, sizeof(E2AP_PDU_t));
    if (!unpack_pdu_aux(pdu, packed_buf_size, responseDataBuf,error_buf_size, responseErrorBuf,ATS_UNALIGNED_BASIC_PER)){
    	printf("#%s failed. Packing error %s\n", __func__, responseErrorBuf);
    }
    responseErrorBuf[0] = 0;
    asn1_pdu_printer(pdu, sizeof(responseErrorBuf), responseErrorBuf);
    printf("#%s: 21%s\n", __func__, responseErrorBuf);
}

void test_build_pack_endc_x2setup_request(){
	size_t error_buf_size = 8192;
	size_t packed_buf_size = 4096;
	unsigned char responseDataBuf[packed_buf_size];
    uint8_t pLMN_Identity[] = {0xa,0xb,0xc};
    uint8_t ric_flag[] = {0xa,0xd,0xe};
    uint8_t eNBId[] = {0xf,0x7,0x2};
    unsigned int bitqty=18;

	char responseErrorBuf[error_buf_size];
	bool result = build_pack_endc_x2setup_request(
			pLMN_Identity, eNBId, bitqty , ric_flag,
			&packed_buf_size, responseDataBuf, error_buf_size, responseErrorBuf);
    if (!result) {
        printf("#%s. Packing error %s\n", __func__, responseErrorBuf);
        return;
    }
    printf("#%s packed size:%lu\nPayload:\n", __func__, packed_buf_size);
    for (size_t i = 0; i < packed_buf_size; ++i)
        printf("%02x",responseDataBuf[i]);
    printf("\n");
}

void test_unpack(void)
{
	return; // No need for now.
	char responseErrorBuf[8192];
	printf("\n--------------- case #1\n\n");
	{
		uint8_t buf[] = {0x00,0x24,0x00,0x32,0x00,0x00,0x01,0x00,0xf4,0x00,0x2b,0x00,0x00,0x02,0x00,0x15,0x00,0x09,0x00,0xbb,0xbc,0xcc,0x80,0x03,0xab,0xcd,0x80,0x00,0xfa,0x00,0x17,0x00,0x00,0x01,0xf7,0x00,0xbb,0xbc,0xcc,0xab,0xcd,0x80,0x00,0x00,0x00,0xbb,0xbc,0xcc,0x00,0x00,0x00,0x00,0x00,0x01};
		E2AP_PDU_t *pdu =calloc(1, sizeof(E2AP_PDU_t));
		if (!unpack_pdu_aux(pdu, sizeof(buf), buf ,sizeof(responseErrorBuf), responseErrorBuf,ATS_ALIGNED_BASIC_PER)){
			printf("#%s failed. Packing error %s\n", __func__, responseErrorBuf);
		}

		responseErrorBuf[0] = 0;
		asn1_pdu_printer(pdu, sizeof(responseErrorBuf), responseErrorBuf);
		printf("#%s: %s\n", __func__, responseErrorBuf);
	}

	printf("\n--------------- case #2\n\n");
	{
		uint8_t buf[] = {0x00,0x06,0x00,0x2b,0x00,0x00,0x02,0x00,0x15,0x00,0x09,0x00,0x0a,0x0b,0x0c,0x81,0x03,0xab,0xcd,0xc0,0x00,0x14,0x00,0x17,0x00,0x00,0x01,0xf7,0x00,0x0a,0x0b,0x0c,0xab,0xcd,0xc0,0x00,0x00,0x00,0x0a,0x0d,0x0e,0x00,0x00,0x00,0x00,0x00,0x01};
		E2AP_PDU_t *pdu =calloc(1, sizeof(E2AP_PDU_t));
		if (!unpack_pdu_aux(pdu, sizeof(buf), buf ,sizeof(responseErrorBuf), responseErrorBuf,ATS_ALIGNED_BASIC_PER)){
			printf("#%s failed. Packing error %s\n", __func__, responseErrorBuf);
		}

		responseErrorBuf[0] = 0;
		asn1_pdu_printer(pdu, sizeof(responseErrorBuf), responseErrorBuf);
		printf("#%s: %s\n", __func__, responseErrorBuf);
	}
}
