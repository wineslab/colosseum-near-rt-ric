/*
 * Copyright 2019 AT&T Intellectual Property
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



/*
 * Generated by asn1c-0.9.29 (http://lionet.info/asn1c)
 * From ASN.1 module "X2AP-IEs"
 * 	found in "../../asnFiles/X2AP-IEs.asn"
 * 	`asn1c -fcompound-names -fincludes-quoted -fno-include-deps -findirect-choice -gen-PER -no-gen-OER -D.`
 */

#include "MDT-Configuration.h"

#include "M1ThresholdEventA2.h"
#include "M1PeriodicReporting.h"
#include "ProtocolExtensionContainer.h"
asn_TYPE_member_t asn_MBR_MDT_Configuration_1[] = {
	{ ATF_NOFLAGS, 0, offsetof(struct MDT_Configuration, mdt_Activation),
		(ASN_TAG_CLASS_CONTEXT | (0 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_MDT_Activation,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"mdt-Activation"
		},
	{ ATF_NOFLAGS, 0, offsetof(struct MDT_Configuration, areaScopeOfMDT),
		(ASN_TAG_CLASS_CONTEXT | (1 << 2)),
		+1,	/* EXPLICIT tag at current level */
		&asn_DEF_AreaScopeOfMDT,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"areaScopeOfMDT"
		},
	{ ATF_NOFLAGS, 0, offsetof(struct MDT_Configuration, measurementsToActivate),
		(ASN_TAG_CLASS_CONTEXT | (2 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_MeasurementsToActivate,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"measurementsToActivate"
		},
	{ ATF_NOFLAGS, 0, offsetof(struct MDT_Configuration, m1reportingTrigger),
		(ASN_TAG_CLASS_CONTEXT | (3 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_M1ReportingTrigger,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"m1reportingTrigger"
		},
	{ ATF_POINTER, 3, offsetof(struct MDT_Configuration, m1thresholdeventA2),
		(ASN_TAG_CLASS_CONTEXT | (4 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_M1ThresholdEventA2,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"m1thresholdeventA2"
		},
	{ ATF_POINTER, 2, offsetof(struct MDT_Configuration, m1periodicReporting),
		(ASN_TAG_CLASS_CONTEXT | (5 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_M1PeriodicReporting,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"m1periodicReporting"
		},
	{ ATF_POINTER, 1, offsetof(struct MDT_Configuration, iE_Extensions),
		(ASN_TAG_CLASS_CONTEXT | (6 << 2)),
		-1,	/* IMPLICIT tag at current level */
		&asn_DEF_ProtocolExtensionContainer_170P166,
		0,
		{ 0, 0, 0 },
		0, 0, /* No default value */
		"iE-Extensions"
		},
};
static const int asn_MAP_MDT_Configuration_oms_1[] = { 4, 5, 6 };
static const ber_tlv_tag_t asn_DEF_MDT_Configuration_tags_1[] = {
	(ASN_TAG_CLASS_UNIVERSAL | (16 << 2))
};
static const asn_TYPE_tag2member_t asn_MAP_MDT_Configuration_tag2el_1[] = {
    { (ASN_TAG_CLASS_CONTEXT | (0 << 2)), 0, 0, 0 }, /* mdt-Activation */
    { (ASN_TAG_CLASS_CONTEXT | (1 << 2)), 1, 0, 0 }, /* areaScopeOfMDT */
    { (ASN_TAG_CLASS_CONTEXT | (2 << 2)), 2, 0, 0 }, /* measurementsToActivate */
    { (ASN_TAG_CLASS_CONTEXT | (3 << 2)), 3, 0, 0 }, /* m1reportingTrigger */
    { (ASN_TAG_CLASS_CONTEXT | (4 << 2)), 4, 0, 0 }, /* m1thresholdeventA2 */
    { (ASN_TAG_CLASS_CONTEXT | (5 << 2)), 5, 0, 0 }, /* m1periodicReporting */
    { (ASN_TAG_CLASS_CONTEXT | (6 << 2)), 6, 0, 0 } /* iE-Extensions */
};
asn_SEQUENCE_specifics_t asn_SPC_MDT_Configuration_specs_1 = {
	sizeof(struct MDT_Configuration),
	offsetof(struct MDT_Configuration, _asn_ctx),
	asn_MAP_MDT_Configuration_tag2el_1,
	7,	/* Count of tags in the map */
	asn_MAP_MDT_Configuration_oms_1,	/* Optional members */
	3, 0,	/* Root/Additions */
	7,	/* First extension addition */
};
asn_TYPE_descriptor_t asn_DEF_MDT_Configuration = {
	"MDT-Configuration",
	"MDT-Configuration",
	&asn_OP_SEQUENCE,
	asn_DEF_MDT_Configuration_tags_1,
	sizeof(asn_DEF_MDT_Configuration_tags_1)
		/sizeof(asn_DEF_MDT_Configuration_tags_1[0]), /* 1 */
	asn_DEF_MDT_Configuration_tags_1,	/* Same as above */
	sizeof(asn_DEF_MDT_Configuration_tags_1)
		/sizeof(asn_DEF_MDT_Configuration_tags_1[0]), /* 1 */
	{ 0, 0, SEQUENCE_constraint },
	asn_MBR_MDT_Configuration_1,
	7,	/* Elements count */
	&asn_SPC_MDT_Configuration_specs_1	/* Additional specs */
};

