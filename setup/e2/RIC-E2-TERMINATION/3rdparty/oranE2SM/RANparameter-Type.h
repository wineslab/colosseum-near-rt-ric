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

/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "E2SM-gNB-NRT-IEs"
 * 	found in "../asnTextFiles/e2sm-gNB-NRT-v401.asn"
 * 	`asn1c -fcompound-names -fincludes-quoted -fno-include-deps -findirect-choice -gen-PER -no-gen-OER -D.`
 */

#ifndef	_RANparameter_Type_H_
#define	_RANparameter_Type_H_


#include "asn_application.h"

/* Including external dependencies */
#include "NativeEnumerated.h"

#ifdef __cplusplus
extern "C" {
#endif

/* Dependencies */
typedef enum RANparameter_Type {
	RANparameter_Type_integer	= 0,
	RANparameter_Type_enumerated	= 1,
	RANparameter_Type_boolean	= 2,
	RANparameter_Type_bit_string	= 3,
	RANparameter_Type_octet_string	= 4,
	RANparameter_Type_printable_string	= 5
	/*
	 * Enumeration is extensible
	 */
} e_RANparameter_Type;

/* RANparameter-Type */
typedef long	 RANparameter_Type_t;

/* Implementation */
extern asn_per_constraints_t asn_PER_type_RANparameter_Type_constr_1;
extern asn_TYPE_descriptor_t asn_DEF_RANparameter_Type;
extern const asn_INTEGER_specifics_t asn_SPC_RANparameter_Type_specs_1;
asn_struct_free_f RANparameter_Type_free;
asn_struct_print_f RANparameter_Type_print;
asn_constr_check_f RANparameter_Type_constraint;
ber_type_decoder_f RANparameter_Type_decode_ber;
der_type_encoder_f RANparameter_Type_encode_der;
xer_type_decoder_f RANparameter_Type_decode_xer;
xer_type_encoder_f RANparameter_Type_encode_xer;
per_type_decoder_f RANparameter_Type_decode_uper;
per_type_encoder_f RANparameter_Type_encode_uper;
per_type_decoder_f RANparameter_Type_decode_aper;
per_type_encoder_f RANparameter_Type_encode_aper;

#ifdef __cplusplus
}
#endif

#endif	/* _RANparameter_Type_H_ */
#include "asn_internal.h"
