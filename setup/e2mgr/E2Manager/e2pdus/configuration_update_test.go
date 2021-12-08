/*******************************************************************************
 *
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
 *
 *******************************************************************************/

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

func TestPrepareEndcConfigurationUpdateFailurePDU(t *testing.T) {
	_,err := logger.InitLogger(logger.InfoLevel)
	if err!=nil{
		t.Errorf("failed to initialize logger, error: %s", err)
	}
	packedPdu := "402500080000010005400142"
	packedEndcConfigurationUpdateFailure := PackedEndcConfigurationUpdateFailure

	tmp := fmt.Sprintf("%x", packedEndcConfigurationUpdateFailure)
	if len(tmp) != len(packedPdu) {
		t.Errorf("want packed len:%d, got: %d\n", len(packedPdu)/2, len(packedEndcConfigurationUpdateFailure)/2)
	}

	if strings.Compare(tmp, packedPdu) != 0 {
		t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", packedPdu, tmp)
	}
}

func TestPrepareEndcConfigurationUpdateFailurePDUFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	err  = prepareEndcConfigurationUpdateFailurePDU(1, 4096)
	if err == nil {
		t.Errorf("want: error, got: success.\n")
	}

	expected:= "#configuration_update.prepareEndcConfigurationUpdateFailurePDU - failed to build and pack the endc configuration update failure message #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("want :[%s], got: [%s]\n", expected, err)
	}
}

func TestPrepareX2EnbConfigurationUpdateFailurePDU(t *testing.T) {
	_,err := logger.InitLogger(logger.InfoLevel)
	if err!=nil{
		t.Errorf("failed to initialize logger, error: %s", err)
	}
	packedPdu := "400800080000010005400142"
	packedEndcX2ConfigurationUpdateFailure := PackedX2EnbConfigurationUpdateFailure

	tmp := fmt.Sprintf("%x", packedEndcX2ConfigurationUpdateFailure)
	if len(tmp) != len(packedPdu) {
		t.Errorf("want packed len:%d, got: %d\n", len(packedPdu)/2, len(packedEndcX2ConfigurationUpdateFailure)/2)
	}

	if strings.Compare(tmp, packedPdu) != 0 {
		t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", packedPdu, tmp)
	}
}

func TestPrepareX2EnbConfigurationUpdateFailurePDUFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	err  = prepareX2EnbConfigurationUpdateFailurePDU(1, 4096)
	if err == nil {
		t.Errorf("want: error, got: success.\n")
	}

	expected:= "#configuration_update.prepareX2EnbConfigurationUpdateFailurePDU - failed to build and pack the x2 configuration update failure message #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("want :[%s], got: [%s]\n", expected, err)
	}
}

func TestPrepareEndcConfigurationUpdateAckPDU(t *testing.T) {
	_,err := logger.InitLogger(logger.InfoLevel)
	if err!=nil{
		t.Errorf("failed to initialize logger, error: %s", err)
	}
	packedPdu := "2025000a00000100f70003000000"
	packedEndcConfigurationUpdateAck := PackedEndcConfigurationUpdateAck

	tmp := fmt.Sprintf("%x", packedEndcConfigurationUpdateAck)
	if len(tmp) != len(packedPdu) {
		t.Errorf("want packed len:%d, got: %d\n", len(packedPdu)/2, len(packedEndcConfigurationUpdateAck)/2)
	}

	if strings.Compare(tmp, packedPdu) != 0 {
		t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", packedPdu, tmp)
	}
}

func TestPrepareEndcConfigurationUpdateAckPDUFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	err  = prepareEndcConfigurationUpdateAckPDU(1, 4096)
	if err == nil {
		t.Errorf("want: error, got: success.\n")
	}

	expected:= "#configuration_update.prepareEndcConfigurationUpdateAckPDU - failed to build and pack the endc configuration update ack message #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("want :[%s], got: [%s]\n", expected, err)
	}
}

func TestPrepareX2EnbConfigurationUpdateAckPDU(t *testing.T) {
	_,err := logger.InitLogger(logger.InfoLevel)
	if err!=nil{
		t.Errorf("failed to initialize logger, error: %s", err)
	}
	packedPdu := "200800080000010011400100"
	packedEndcX2ConfigurationUpdateAck := PackedX2EnbConfigurationUpdateAck

	tmp := fmt.Sprintf("%x", packedEndcX2ConfigurationUpdateAck)
	if len(tmp) != len(packedPdu) {
		t.Errorf("want packed len:%d, got: %d\n", len(packedPdu)/2, len(packedEndcX2ConfigurationUpdateAck)/2)
	}

	if strings.Compare(tmp, packedPdu) != 0 {
		t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", packedPdu, tmp)
	}
}

func TestPrepareX2EnbConfigurationUpdateAckPDUFailure(t *testing.T) {
	_, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	err  = prepareX2EnbConfigurationUpdateAckPDU(1, 4096)
	if err == nil {
		t.Errorf("want: error, got: success.\n")
	}

	expected:= "#configuration_update.prepareX2EnbConfigurationUpdateAckPDU - failed to build and pack the x2 configuration update ack message #src/asn1codec_utils.c.pack_pdu_aux - Encoded output of E2AP-PDU, is too big"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("want :[%s], got: [%s]\n", expected, err)
	}
}