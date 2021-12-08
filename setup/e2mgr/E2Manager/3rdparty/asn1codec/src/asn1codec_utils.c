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


#include <stdio.h>
#include <string.h>
#include <errno.h>
#undef NDEBUG
#include <assert.h>
#include <asn1codec_utils.h>
#include <constr_TYPE.h>
#include <xer_encoder.h>

/*
 * Printer for the e2ap pdu.
 * The string representation of the pdu stored in buf.
 *
 * Input:
 * pdu - the pdu to print.
 * buf_size - the size of the storage buffer.
 * buf - hold the string representation of the pdu.
 */
bool
asn1_pdu_printer(E2AP_PDU_t const *pdu, size_t buf_size, char *buf)
{
	bool rc = true;
	char *bufloc = 0;
	size_t sizeloc = 0;
	buf[0] = 0;
	FILE *stream = open_memstream(&bufloc, &sizeloc);

	errno = 0;
	if (asn_fprint(stream, &asn_DEF_E2AP_PDU, pdu)){
		snprintf(buf, buf_size, "#%s.%s - Failed to print %s, error = %d ", __FILE__, __func__, asn_DEF_E2AP_PDU.name, errno);
		strerror_r(errno, buf+strlen(buf), buf_size - strlen(buf));
		rc = false;
	} else {
		buf_size = buf_size > sizeloc ? sizeloc: buf_size -1;
		memcpy(buf, bufloc, buf_size);
		buf[buf_size] = 0;
	}

	fclose(stream);
	free(bufloc);
	return rc;
}


/*
 * XML Printer for the e2ap pdu.
 * The string representation of the pdu stored in buf.
 *
 * Input:
 * pdu - the pdu to print.
 * buf_size - the size of the storage buffer.
 * buf - hold the string representation of the pdu.
 */
bool
asn1_pdu_xer_printer(E2AP_PDU_t const *pdu, size_t buf_size, char *buf)
{
	bool rc = true;
	char *bufloc = 0;
	size_t sizeloc = 0;
	buf[0] = 0;
	FILE *stream = open_memstream(&bufloc, &sizeloc);

	errno = 0;
	if (xer_fprint(stream, &asn_DEF_E2AP_PDU, pdu)){
		snprintf(buf, buf_size, "#%s.%s - Failed to print %s, error = %d ", __FILE__, __func__, asn_DEF_E2AP_PDU.name, errno);
		strerror_r(errno, buf+strlen(buf), buf_size - strlen(buf));
		rc = false;
	} else {
		buf_size = buf_size > sizeloc ? sizeloc: buf_size -1;
		memcpy(buf, bufloc, buf_size);
		buf[buf_size] = 0;
	}

	fclose(stream);
	free(bufloc);
	return rc;
}

/*
 * Unpack the pdu from ASN.1 PER encoding.
 *
 * Input:
 * pdu - storage for unpacked pdu.
 * packed_buf_size - size of the encoded data.
 * packed_buf - storage of the packed pdu
 * err_buf_size - size of the err_buf which may hold the error string in case of
 * an error. err_buf - storage for the error string
 *
 * Return: true in case of success, false in case of failure.
 */
bool
per_unpack_pdu(E2AP_PDU_t *pdu, size_t packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf)
{
	return unpack_pdu_aux(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf,ATS_ALIGNED_BASIC_PER);
}

bool
unpack_pdu_aux(E2AP_PDU_t *pdu, size_t packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf,enum asn_transfer_syntax syntax)
{
	char spec[256];
	size_t err_msg_size = err_buf_size;

	//ATS_BASIC_XER ATS_ALIGNED_BASIC_PER, ATS_UNALIGNED_BASIC_PER,ATS_ALIGNED_CANONICAL_PER
	errno = 0;
	asn_dec_rval_t rval =
	asn_decode(0,syntax , &asn_DEF_E2AP_PDU, (void **)&pdu, packed_buf, packed_buf_size);
	switch(rval.code) {
	case RC_OK:
		if (asn_check_constraints(&asn_DEF_E2AP_PDU, pdu,err_buf, &err_msg_size)){
			snprintf(spec, sizeof(spec), "#%s.%s - Constraint check failed: ", __FILE__, __func__);
			size_t spec_actual_size = strlen(spec);
			if (spec_actual_size + err_msg_size < err_buf_size){
				memmove(err_buf + spec_actual_size, err_buf, err_msg_size + 1);
				memcpy(err_buf, spec, spec_actual_size);
			}
			return false;
		}
		return true;

	break;
	case RC_WMORE:
	case RC_FAIL:
	default:
		break;
	}

	snprintf(err_buf, err_buf_size, "#%s.%s - Failed to decode %s (consumed %zu), error = %d ", __FILE__, __func__, asn_DEF_E2AP_PDU.name, rval.consumed, errno);
	strerror_r(errno, err_buf+strlen(err_buf), err_buf_size - strlen(err_buf));
	return false;
}

/*
 * Pack the pdu using ASN.1 PER encoding.
 *
 * Input:
 * pdu - the pdu to pack.
 * packed_buf_size - in: size of packed_buf; out: number of chars used.
 * packed_buf - storage for the packed pdu
 * err_buf_size - size of the err_buf which may hold the error string in case of
 * an error. err_buf - storage for the error string
 *
 * Return: true in case of success, false in case of failure.
 */
bool
per_pack_pdu(E2AP_PDU_t *pdu, size_t *packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf)
{
	return pack_pdu_aux(pdu, packed_buf_size, packed_buf,err_buf_size, err_buf,ATS_ALIGNED_BASIC_PER);
}

bool
pack_pdu_aux(E2AP_PDU_t *pdu, size_t *packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf,enum asn_transfer_syntax syntax)
{
	char spec[256];
	size_t err_msg_size = err_buf_size;

    if (asn_check_constraints(&asn_DEF_E2AP_PDU, pdu,err_buf, &err_msg_size)){
		snprintf(spec, sizeof(spec), "#%s.%s - Constraint check failed: ", __FILE__, __func__);
		size_t spec_actual_size = strlen(spec);
		if (spec_actual_size + err_msg_size < err_buf_size){
			memmove(err_buf + spec_actual_size, err_buf, err_msg_size + 1);
			memcpy(err_buf, spec, spec_actual_size);
		}
    	return false;
    }

	errno = 0;
asn_enc_rval_t res =
		asn_encode_to_buffer(0, syntax, &asn_DEF_E2AP_PDU, pdu, packed_buf, *packed_buf_size);
	if(res.encoded == -1) {
		snprintf(err_buf, err_buf_size, "#%s.%s - Failed to encode %s, error = %d ", __FILE__, __func__, asn_DEF_E2AP_PDU.name, errno);
		strerror_r(errno, err_buf+strlen(err_buf), err_buf_size - strlen(err_buf));
		return false;
	} else {
		/* Encoded successfully. */
		if (*packed_buf_size < res.encoded){
			snprintf(err_buf, err_buf_size, "#%s.%s - Encoded output of %s, is too big:%zu", __FILE__, __func__, asn_DEF_E2AP_PDU.name,res.encoded);
			return false;
		} else {
			*packed_buf_size = res.encoded;
		}
	}
	return true;
}

/*
 * Create a new pdu.
 * Abort the process on allocation failure.
 */
E2AP_PDU_t *new_pdu(size_t sz /*ignored (may be used for a custom allocator)*/)
{
	E2AP_PDU_t *pdu = calloc(1, sizeof(E2AP_PDU_t));
	assert(pdu != 0);
	return pdu;
}

void delete_pdu(E2AP_PDU_t *pdu)
{
	ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
}

