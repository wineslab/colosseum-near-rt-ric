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


#ifndef INC_ASN1CODEC_UTILS_H_
#define INC_ASN1CODEC_UTILS_H_

#ifndef ASN_DISABLE_OER_SUPPORT
#define ASN_DISABLE_OER_SUPPORT
#endif

#ifndef ASN_PDU_COLLECTION
#define ASN_PDU_COLLECTION
#endif

#include <stdbool.h>
#include <E2AP-PDU.h>
#include <ProtocolIE-Field.h>
#include <ProtocolExtensionContainer.h>
#include <ProtocolExtensionField.h>
#include <CriticalityDiagnostics-IE-List.h>

#define pLMN_Identity_size 3
#define shortMacro_eNB_ID_size  18
#define macro_eNB_ID_size       20
#define longMacro_eNB_ID_size   21
#define home_eNB_ID_size        28
#define eUTRANcellIdentifier_size 28

#ifdef __cplusplus
extern "C"
{
#endif

bool asn1_pdu_printer(const E2AP_PDU_t *pdu, size_t obufsz, char *buf);
bool asn1_pdu_xer_printer(const E2AP_PDU_t *pdu, size_t obufsz, char *buf);
bool per_unpack_pdu(E2AP_PDU_t *pdu, size_t packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf);
bool per_pack_pdu(E2AP_PDU_t *pdu, size_t *packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf);
bool unpack_pdu_aux(E2AP_PDU_t *pdu, size_t packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf,enum asn_transfer_syntax syntax);
bool pack_pdu_aux(E2AP_PDU_t *pdu, size_t *packed_buf_size, unsigned char* packed_buf,size_t err_buf_size, char* err_buf,enum asn_transfer_syntax syntax);

E2AP_PDU_t *new_pdu(size_t sz);
void delete_pdu(E2AP_PDU_t *pdu);

#ifdef __cplusplus
}
#endif

#endif /* INC_ASN1CODEC_UTILS_H_ */
