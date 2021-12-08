/*
 *   Copyright (c) 2019 AT&T Intellectual Property.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

/*
 * This source code is part of the near-RT RIC (RAN Intelligent Controller)
 * platform project (RICP).
 */

package e2pdus

import (
	"e2mgr/logger"
	"fmt"
	"strings"
	"testing"
)

func TestKnownCausesToX2ResetPDU(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}
	var testCases = []struct {
		cause     string
		packedPdu string
	}{
		{
			cause:     OmInterventionCause,
			packedPdu: "000700080000010005400164",
		},
		{
			cause:     "PROTOCOL:transfer-syntax-error",
			packedPdu: "000700080000010005400140",
		},
		{
			cause:     "transport:transport-RESOURCE-unavailable",
			packedPdu: "000700080000010005400120",
		},

		{
			cause:     "radioNetwork:invalid-MME-groupid",
			packedPdu: "00070009000001000540020680",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.packedPdu, func(t *testing.T) {

			payload, ok := KnownCausesToX2ResetPDU(tc.cause)
			if !ok {
				t.Errorf("want: success, got: not found.\n")
			} else {
				tmp := fmt.Sprintf("%x", payload)
				if len(tmp) != len(tc.packedPdu) {
					t.Errorf("want packed len:%d, got: %d\n", len(tc.packedPdu)/2, len(payload)/2)
				}

				if strings.Compare(tmp, tc.packedPdu) != 0 {
					t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", tc.packedPdu, tmp)
				}
			}
		})
	}
}

func TestKnownCausesToX2ResetPDUFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	_, ok := KnownCausesToX2ResetPDU("xxxx")
	if ok {
		t.Errorf("want: not found, got: success.\n")
	}
}

func TestPrepareX2ResetPDUsFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	err = prepareX2ResetPDUs(1, 4096)
	if err == nil {
		t.Errorf("want: error, got: success.\n")
	}

	expected := "failed to build and pack the reset message #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big:"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("want :[%s], got: [%s]\n", expected, err)
	}
}
