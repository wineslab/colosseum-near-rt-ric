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

func TestPrepareX2ResetResponsePDU(t *testing.T) {
	_,err := logger.InitLogger(logger.InfoLevel)
	if err!=nil{
		t.Errorf("failed to initialize logger, error: %s", err)
	}
	packedPdu := "200700080000010011400100"
	packedX2ResetResponse := PackedX2ResetResponse

	tmp := fmt.Sprintf("%x", packedX2ResetResponse)
	if len(tmp) != len(packedPdu) {
		t.Errorf("want packed len:%d, got: %d\n", len(packedPdu)/2, len(packedX2ResetResponse)/2)
	}

	if strings.Compare(tmp, packedPdu) != 0 {
		t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", packedPdu, tmp)
	}
}

func TestPrepareX2ResetResponsePDUFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	err  = prepareX2ResetResponsePDU(1, 4096)
	if err == nil {
		t.Errorf("want: error, got: success.\n")
	}

	expected:= "#x2_reset_response.prepareX2ResetResponsePDU - failed to build and pack the reset response message #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("want :[%s], got: [%s]\n", expected, err)
	}
}