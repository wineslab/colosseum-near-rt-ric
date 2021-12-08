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

package converters

import (
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"fmt"
	"strings"
	"testing"
)

/*
 * Unpack an x2setup response returned from RAN.
 * Verify it matches the want pdu.
 */

func TestUnpackX2apSetupResponse(t *testing.T) {
	logger, _ := logger.InitLogger(logger.DebugLevel)

	wantPduAsStr := `SuccessfulOutcome ::= {
            procedureCode: 6
            criticality: 0 (reject)
            value: X2SetupResponse ::= {
                protocolIEs: ProtocolIE-Container ::= {
                    X2SetupResponse-IEs ::= {
                        id: 21
                        criticality: 0 (reject)
                        value: GlobalENB-ID ::= {
                            pLMN-Identity: 02 F8 29
                            eNB-ID: 00 7A 80 (4 bits unused)
                        }
                    }
                    X2SetupResponse-IEs ::= {
                        id: 20
                        criticality: 0 (reject)
                        value: ServedCells ::= {
                            SEQUENCE ::= {
                                servedCellInfo: ServedCell-Information ::= {
                                    pCI: 99
                                    cellId: ECGI ::= {
                                        pLMN-Identity: 02 F8 29
                                        eUTRANcellIdentifier: 00 07 AB 50 (4 bits unused)
                                    }
                                    tAC: 01 02
                                    broadcastPLMNs: BroadcastPLMNs-Item ::= {
                                        02 F8 29
                                    }
                                    eUTRA-Mode-Info: FDD-Info ::= {
                                        uL-EARFCN: 1
                                        dL-EARFCN: 1
                                        uL-Transmission-Bandwidth: 3 (bw50)
                                        dL-Transmission-Bandwidth: 3 (bw50)
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }`

	inputPayloadAsStr := "2006002a000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829000001000133"
	var payload []byte

	_, err := fmt.Sscanf(inputPayloadAsStr, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}

	response, err := UnpackX2apPduAndRefine(logger, e2pdus.MaxAsn1CodecAllocationBufferSize , len(payload), payload, e2pdus.MaxAsn1CodecMessageBufferSize /*message buffer*/)
	if err != nil {
		t.Errorf("want: success, got: unpack failed. Error: %v\n", err)
	}

	want := strings.Fields(wantPduAsStr)
	got := strings.Fields(response.PduPrint)
	if len(want) != len(got) {
		t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", wantPduAsStr, response.PduPrint)
	}
	for i := 0; i < len(want); i++ {
		if strings.Compare(want[i], got[i]) != 0 {
			t.Errorf("\nwant :\t[%s]\n got: \t\t[%s]\n", wantPduAsStr, strings.TrimSpace(response.PduPrint))
		}

	}
}

/*unpacking error*/

func TestUnpackX2apSetupResponseUnpackError(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)

	wantError := "unpacking error: #src/asn1codec_utils.c.unpack_pdu_aux - Failed to decode E2AP-PDU (consumed 0), error = 0 Success"
	//--------------------2006002a
	inputPayloadAsStr := "2006002b000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829000001000133"
	var payload []byte
	_, err := fmt.Sscanf(inputPayloadAsStr, "%x", &payload)
	if err != nil {
		t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
	}

	_, err = UnpackX2apPduAndRefine(logger, e2pdus.MaxAsn1CodecAllocationBufferSize , len(payload), payload, e2pdus.MaxAsn1CodecMessageBufferSize /*message buffer*/)
	if err != nil {
		if 0 != strings.Compare(fmt.Sprintf("%s", err), wantError) {
			t.Errorf("want failure: %s, got: %s", wantError, err)
		}
	} else {
		t.Errorf("want failure: %s, got: success", wantError)

	}
}
