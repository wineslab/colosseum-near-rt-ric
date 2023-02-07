/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "E2SM-KPM-RC"
 * 	found in "e2sm-kpm-rc.asn"
 * 	`asn1c -fcompound-names -fno-include-deps -findirect-choice -pdu=auto -gen-PER -gen-OER -no-gen-example -D .`
 */

#ifndef	_MeasQuantityResultsEUTRA_H_
#define	_MeasQuantityResultsEUTRA_H_


#include <asn_application.h>

/* Including external dependencies */
#include "RSRP-RangeEUTRA.h"
#include "RSRQ-RangeEUTRA.h"
#include "SINR-RangeEUTRA.h"
#include <constr_SEQUENCE.h>

#ifdef __cplusplus
extern "C" {
#endif

/* MeasQuantityResultsEUTRA */
typedef struct MeasQuantityResultsEUTRA {
	RSRP_RangeEUTRA_t	*rsrp;	/* OPTIONAL */
	RSRQ_RangeEUTRA_t	*rsrq;	/* OPTIONAL */
	SINR_RangeEUTRA_t	*sinr;	/* OPTIONAL */
	/*
	 * This type is extensible,
	 * possible extensions are below.
	 */
	
	/* Context for parsing across buffer boundaries */
	asn_struct_ctx_t _asn_ctx;
} MeasQuantityResultsEUTRA_t;

/* Implementation */
extern asn_TYPE_descriptor_t asn_DEF_MeasQuantityResultsEUTRA;
extern asn_SEQUENCE_specifics_t asn_SPC_MeasQuantityResultsEUTRA_specs_1;
extern asn_TYPE_member_t asn_MBR_MeasQuantityResultsEUTRA_1[3];

#ifdef __cplusplus
}
#endif

#endif	/* _MeasQuantityResultsEUTRA_H_ */
#include <asn_internal.h>