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

#include <errno.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <asn1codec_utils.h>

int
main(int argc, char* argv[])
{
	char responseErrorBuf[1<< 16];
	uint8_t buf[1 << 16];
	size_t count = fread(buf, 1, sizeof(buf), stdin);
	if (count == sizeof(buf)) {
		printf("#%s failed. Input is too big\n", __func__);
		exit(-1);
	}
	if (!feof(stdin)){
		printf("#%s failed. Error while reading input: %s\n", __func__, strerror(errno));
		exit(-1);
	}

	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	if (!unpack_pdu_aux(pdu, count, buf ,sizeof(responseErrorBuf), responseErrorBuf,ATS_BASIC_XER)){
		printf("#%s failed. Unpacking error %s\n", __func__, responseErrorBuf);
		exit(-1);
	}

	responseErrorBuf[0] = 0;
	asn1_pdu_printer(pdu, sizeof(responseErrorBuf), responseErrorBuf);
	printf("#%s: %s\n", __func__, responseErrorBuf);

	{
	size_t per_packed_buf_size;
	uint8_t per_packed_buf[1 << 16];

	responseErrorBuf[0] = 0;
	per_packed_buf_size = sizeof(per_packed_buf);

	if (!pack_pdu_aux(pdu,&per_packed_buf_size, per_packed_buf,sizeof(responseErrorBuf), responseErrorBuf, ATS_ALIGNED_BASIC_PER)){
		printf("#%s failed. Packing error %s\n", __func__, responseErrorBuf);
		exit(-1);
	}
	ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
	printf("#%s packed size:%zu\nPayload:\n", __func__, per_packed_buf_size);
	for (size_t i= 0; i < per_packed_buf_size; i++)
		printf("%02x",per_packed_buf[i]);
	printf("\n");

	pdu =calloc(1, sizeof(E2AP_PDU_t));
	if (!unpack_pdu_aux(pdu, per_packed_buf_size, per_packed_buf ,sizeof(responseErrorBuf), responseErrorBuf,ATS_ALIGNED_BASIC_PER)){
		printf("#%s failed. Packing error %s\n", __func__, responseErrorBuf);
	}

	responseErrorBuf[0] = 0;
	asn1_pdu_printer(pdu, sizeof(responseErrorBuf), responseErrorBuf);
	ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
	printf("#%s: %s\n", __func__, responseErrorBuf);
	}
    exit(0);
}


