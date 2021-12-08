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
	"e2mgr/logger"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"strings"
	"testing"
)

/*
Test permutations of x2 setup response to protobuf enb
*/

func TestUnpackX2SetupFailureResponseAndExtract(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)

	var testCases = []struct {
		response  string
		packedPdu string
		failure   error
	}{
		{
			response: "CONNECTED_SETUP_FAILED network_layer_cause:HANDOVER_DESIRABLE_FOR_RADIO_REASONS time_to_wait:V1S criticality_diagnostics:<procedure_code:33 triggering_message:UNSUCCESSFUL_OUTCOME procedure_criticality:NOTIFY information_element_criticality_diagnostics:<ie_criticality:REJECT ie_id:128 type_of_error:MISSING > > ",
			/*
				E2AP-PDU:
				 unsuccessfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupFailure
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x5
				     criticality_t = 0x1
				     Cause:
				      radioNetwork_t = 0
				    ProtocolIE_Container_elm
				     id_t = 0x16
				     criticality_t = 0x1
				     TimeToWait = 0
				    ProtocolIE_Container_elm
				     id_t = 0x11
				     criticality_t = 0x1
				     CriticalityDiagnostics
				      procedureCode_t = 0x21
				      triggeringMessage_t = 0x2
				      procedureCriticality_t = 0x2
				      iEsCriticalityDiagnostics_t:
				       CriticalityDiagnostics_IE_List_elm
				        iECriticality_t = 0
				        iE_ID_t = 0x80
				        typeOfError_t = 0x1
			*/
			packedPdu: "4006001a0000030005400200000016400100001140087821a00000008040"},
		{
			response: "CONNECTED_SETUP_FAILED transport_layer_cause:TRANSPORT_RESOURCE_UNAVAILABLE criticality_diagnostics:<procedure_code:33 triggering_message:UNSUCCESSFUL_OUTCOME procedure_criticality:NOTIFY information_element_criticality_diagnostics:<ie_criticality:REJECT ie_id:128 type_of_error:MISSING > > ",
			/*
				E2AP-PDU:
				 unsuccessfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupFailure
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x5
				     criticality_t = 0x1
				     Cause:
				      transport_t = 0
				    ProtocolIE_Container_elm
				     id_t = 0x11
				     criticality_t = 0x1
				     CriticalityDiagnostics
				      procedureCode_t = 0x21
				      triggeringMessage_t = 0x2
				      procedureCriticality_t = 0x2
				      iEsCriticalityDiagnostics_t:
				       CriticalityDiagnostics_IE_List_elm
				        iECriticality_t = 0
				        iE_ID_t = 0x80
				        typeOfError_t = 0x1
			*/
			packedPdu: "400600140000020005400120001140087821a00000008040"},
		{
			response: "CONNECTED_SETUP_FAILED protocol_cause:ABSTRACT_SYNTAX_ERROR_IGNORE_AND_NOTIFY criticality_diagnostics:<triggering_message:UNSUCCESSFUL_OUTCOME procedure_criticality:NOTIFY information_element_criticality_diagnostics:<ie_criticality:REJECT ie_id:128 type_of_error:MISSING > > ",
			/*
				E2AP-PDU:
				 unsuccessfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupFailure
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x5
				     criticality_t = 0x1
				     Cause:
				      protocol_t = 0x2
				    ProtocolIE_Container_elm
				     id_t = 0x11
				     criticality_t = 0x1
				     CriticalityDiagnostics
				      triggeringMessage_t = 0x2
				      procedureCriticality_t = 0x2
				      iEsCriticalityDiagnostics_t:
				       CriticalityDiagnostics_IE_List_elm
				        iECriticality_t = 0
				        iE_ID_t = 0x80
				        typeOfError_t = 0x1
			*/
			packedPdu: "400600130000020005400144001140073a800000008040"},

		{
			response: "CONNECTED_SETUP_FAILED miscellaneous_cause:UNSPECIFIED criticality_diagnostics:<procedure_criticality:NOTIFY information_element_criticality_diagnostics:<ie_criticality:REJECT ie_id:128 type_of_error:MISSING > > ",
			/*
				E2AP-PDU:
				 unsuccessfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupFailure
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x5
				     criticality_t = 0x1
				     Cause:
				      misc_t = 0x4
				    ProtocolIE_Container_elm
				     id_t = 0x11
				     criticality_t = 0x1
				     CriticalityDiagnostics
				      procedureCriticality_t = 0x2
				      iEsCriticalityDiagnostics_t:
				       CriticalityDiagnostics_IE_List_elm
				        iECriticality_t = 0
				        iE_ID_t = 0x80
				        typeOfError_t = 0x1
			*/
			packedPdu: "400600120000020005400168001140061a0000008040"},

		{
			response: "CONNECTED_SETUP_FAILED miscellaneous_cause:UNSPECIFIED criticality_diagnostics:<information_element_criticality_diagnostics:<ie_criticality:REJECT ie_id:128 type_of_error:MISSING > information_element_criticality_diagnostics:<ie_criticality:NOTIFY ie_id:255 type_of_error:NOT_UNDERSTOOD > > ",
			/*
				E2AP-PDU:
				 unsuccessfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupFailure
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x5
				     criticality_t = 0x1
				     Cause:
				      misc_t = 0x4
				    ProtocolIE_Container_elm
				     id_t = 0x11
				     criticality_t = 0x1
				     CriticalityDiagnostics
				      iEsCriticalityDiagnostics_t:
				       CriticalityDiagnostics_IE_List_elm
				        iECriticality_t = 0
				        iE_ID_t = 0x80
				        typeOfError_t = 0x1
				       CriticalityDiagnostics_IE_List_elm
				        iECriticality_t = 0x2
				        iE_ID_t = 0xff
				        typeOfError_t = 0
			*/
			packedPdu: "4006001500000200054001680011400908010000804800ff00"},


		{
			response: "CONNECTED_SETUP_FAILED miscellaneous_cause:UNSPECIFIED criticality_diagnostics:<procedure_code:33 > ",
			/*
				E2AP-PDU:
				 unsuccessfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupFailure
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x5
				     criticality_t = 0x1
				     Cause:
				      misc_t = 0x4
				    ProtocolIE_Container_elm
				     id_t = 0x11
				     criticality_t = 0x1
				     CriticalityDiagnostics
				      procedureCode_t = 0x21
			*/
			packedPdu: "4006000e0000020005400168001140024021"},

		{
			response: "CONNECTED_SETUP_FAILED miscellaneous_cause:UNSPECIFIED ",
			/*
				E2AP-PDU:
				 unsuccessfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupFailure
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x5
				     criticality_t = 0x1
				     Cause:
				      misc_t = 0x4
			*/
			packedPdu: "400600080000010005400168"},
		{
			response: "CONNECTED_SETUP_FAILED network_layer_cause:HANDOVER_DESIRABLE_FOR_RADIO_REASONS time_to_wait:V1S criticality_diagnostics:<procedure_code:33 triggering_message:UNSUCCESSFUL_OUTCOME procedure_criticality:NOTIFY information_element_criticality_diagnostics:<ie_criticality:REJECT ie_id:128 type_of_error:MISSING > > ",
			/*
				E2AP-PDU:
				 unsuccessfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupFailure
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x5
				     criticality_t = 0x1
				     Cause:
				      radioNetwork_t = 0
				    ProtocolIE_Container_elm
				     id_t = 0x16
				     criticality_t = 0x1
				     TimeToWait = 0
				    ProtocolIE_Container_elm
				     id_t = 0x11
				     criticality_t = 0x1
				     CriticalityDiagnostics
				      procedureCode_t = 0x21
				      triggeringMessage_t = 0x2
				      procedureCriticality_t = 0x2
				      iEsCriticalityDiagnostics_t:
				       CriticalityDiagnostics_IE_List_elm
				        iECriticality_t = 0
				        iE_ID_t = 0x80
				        typeOfError_t = 0x1
			*/
			packedPdu: "4006001a0000030005400200000016400100001140087821a00000008040",
			//failure: fmt.Errorf("getAtom for path [unsuccessfulOutcome_t X2SetupFailure protocolIEs_t ProtocolIE_Container_elm Cause radioNetwork_t] failed, rc = 2" /*NO_SPACE_LEFT*/),
		},
	}

	converter := NewX2SetupFailureResponseConverter(logger)

	for _, tc := range testCases {
		t.Run(tc.packedPdu, func(t *testing.T) {

			var payload []byte
			_, err := fmt.Sscanf(tc.packedPdu, "%x", &payload)
			if err != nil {
				t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
			}

			response, err := converter.UnpackX2SetupFailureResponseAndExtract(payload)

			if err != nil {
				if tc.failure == nil {
					t.Errorf("want: success, got: error: %v\n", err)
				} else {
					if strings.Compare(err.Error(), tc.failure.Error()) != 0 {
						t.Errorf("want: %s, got: %s", tc.failure, err)
					}
				}
			}

			if response == nil {
				if tc.failure == nil {
					t.Errorf("want: response=%s, got: empty response", tc.response)
				}
			} else {
				nb := &entities.NodebInfo{}
				nb.ConnectionStatus = entities.ConnectionStatus_CONNECTED_SETUP_FAILED
				nb.SetupFailure = response
				nb.FailureType = entities.Failure_X2_SETUP_FAILURE
				respStr := fmt.Sprintf("%s %s", nb.ConnectionStatus, response)
				if !strings.EqualFold(respStr, tc.response) {
					t.Errorf("want: response=[%s], got: [%s]", tc.response, respStr)
				}

			}
		})
	}
}
