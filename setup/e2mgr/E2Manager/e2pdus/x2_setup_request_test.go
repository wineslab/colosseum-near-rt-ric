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
	"fmt"
	"strings"
	"testing"
)

/*
 * Create and pack an x2ap setup request.
 * Verify the packed representation matches the want value.
 */
func TestPackX2apSetupRequest(t *testing.T) {
	pLMNId := []byte{0xbb, 0xbc, 0xcc}
	ricFlag := []byte{0xbb, 0xbc, 0xcc} /*pLMNId [3]bytes*/

	var testCases = []struct {
		eNBId       []byte
		eNBIdBitqty uint
		packedPdu   string
	}{
		{
			eNBId:       []byte{0xab, 0xcd, 0x2}, /*00000010 -> 10000000*/
			eNBIdBitqty: ShortMacro_eNB_ID,
			packedPdu:   "0006002b0000020015000900bbbccc8003abcd8000140017000001f700bbbcccabcd80000000bbbccc000000000001",
		},

		{
			eNBId:       []byte{0xab, 0xcd, 0xe},
			eNBIdBitqty: Macro_eNB_ID,
			packedPdu:   "0006002a0000020015000800bbbccc00abcde000140017000001f700bbbcccabcde0000000bbbccc000000000001",
		},
		{
			eNBId:       []byte{0xab, 0xcd, 0x7}, /*00000111 -> 00111000*/
			eNBIdBitqty: LongMacro_eNB_ID,
			//packedPdu:   "0006002b0000020015000900bbbccc8103abcd3800140017000001f700bbbcccabcd38000000bbbccc000000000001",
			packedPdu:   "0006002b0000020015000900bbbcccc003abcd3800140017000001f700bbbcccabcd38000000bbbccc000000000001",
		},
		{
			eNBId:       []byte{0xab, 0xcd, 0xef, 0x8},
			eNBIdBitqty: Home_eNB_ID,
			packedPdu:   "0006002b0000020015000900bbbccc40abcdef8000140017000001f700bbbcccabcdef800000bbbccc000000000001",
		},
	}

	// TODO: Consider using testify's assert/require
	// testing/quick to input random value
	for _, tc := range testCases {
		t.Run(tc.packedPdu, func(t *testing.T) {

			payload, _, err := preparePackedX2SetupRequest(MaxAsn1PackedBufferSize /*max packed buffer*/, MaxAsn1CodecMessageBufferSize /*max message buffer*/, pLMNId, tc.eNBId, tc.eNBIdBitqty,ricFlag)
			if err != nil {
				t.Errorf("want: success, got: pack failed. Error: %v\n", err)
			} else {
				t.Logf("packed X2AP setup request(size=%d): %x\n", len(payload), payload)
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

/*Packing error*/

func TestPackX2apSetupRequestPackError(t *testing.T) {

	wantError := "packing error: #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big:46"
	pLMNId := []byte{0xbb, 0xbc, 0xcc}
	ricFlag := []byte{0xbb, 0xbc, 0xcc} /*pLMNId [3]bytes*/
	eNBId := []byte{0xab, 0xcd, 0xe}
	eNBIdBitqty := uint(Macro_eNB_ID)
	_, _, err := preparePackedX2SetupRequest(40 /*max packed buffer*/, MaxAsn1CodecMessageBufferSize /*max message buffer*/, pLMNId, eNBId, eNBIdBitqty, ricFlag)
	if err != nil {
		if 0 != strings.Compare(fmt.Sprintf("%s", err), wantError) {
			t.Errorf("want failure: %s, got: %s", wantError, err)
		}
	} else {
		t.Errorf("want failure: %s, got: success", wantError)

	}
}
