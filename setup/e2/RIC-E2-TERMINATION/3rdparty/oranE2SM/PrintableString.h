/*
 *
 * Copyright 2020 AT&T Intellectual Property
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
 *
 */

/*-
 * Copyright (c) 2003-2017 Lev Walkin <vlm@lionet.info>. All rights reserved.
 * Redistribution and modifications are permitted subject to BSD license.
 */
#ifndef	_PrintableString_H_
#define	_PrintableString_H_

#include <OCTET_STRING.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef OCTET_STRING_t PrintableString_t;  /* Implemented via OCTET STRING */

extern asn_TYPE_descriptor_t asn_DEF_PrintableString;
extern asn_TYPE_operation_t asn_OP_PrintableString;

asn_constr_check_f PrintableString_constraint;

#define PrintableString_free            OCTET_STRING_free
#define PrintableString_print           OCTET_STRING_print_utf8
#define PrintableString_compare         OCTET_STRING_compare
#define PrintableString_decode_ber      OCTET_STRING_decode_ber
#define PrintableString_encode_der      OCTET_STRING_encode_der
#define PrintableString_decode_xer      OCTET_STRING_decode_xer_utf8
#define PrintableString_encode_xer      OCTET_STRING_encode_xer_utf8
#define PrintableString_decode_uper     OCTET_STRING_decode_uper
#define PrintableString_encode_uper     OCTET_STRING_encode_uper
#define PrintableString_decode_aper     OCTET_STRING_decode_aper
#define PrintableString_encode_aper     OCTET_STRING_encode_aper

#ifdef __cplusplus
}
#endif

#endif	/* _PrintableString_H_ */
