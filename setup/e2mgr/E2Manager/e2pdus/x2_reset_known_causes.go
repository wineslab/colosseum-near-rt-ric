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

// #cgo CFLAGS: -I../3rdparty/asn1codec/inc/  -I../3rdparty/asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../3rdparty/asn1codec/lib/ -L../3rdparty/asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <x2reset_request_wrapper.h>
import "C"
import (
	"fmt"
	"strings"
	"unsafe"
)

const (
	MaxAsn1PackedBufferSize       = 4096
	MaxAsn1CodecMessageBufferSize = 4096
	MaxAsn1CodecAllocationBufferSize = 4096 // TODO: remove later
)

// Used as default by the x2_reset_request
const (
	OmInterventionCause = "misc:om-intervention"
)

type cause struct {
	causeGroup uint32
	cause      int
}

var knownCauses = map[string]cause{
	"misc:control-processing-overload":                {causeGroup: C.Cause_PR_misc, cause: C.CauseMisc_control_processing_overload},
	"misc:hardware-failure":                           {causeGroup: C.Cause_PR_misc, cause: C.CauseMisc_hardware_failure},
	OmInterventionCause:                               {causeGroup: C.Cause_PR_misc, cause: C.CauseMisc_om_intervention},
	"misc:not-enough-user-plane-processing-resources": {causeGroup: C.Cause_PR_misc, cause: C.CauseMisc_not_enough_user_plane_processing_resources},
	"misc:unspecified":                                {causeGroup: C.Cause_PR_misc, cause: C.CauseMisc_unspecified},

	"protocol:transfer-syntax-error":                             {causeGroup: C.Cause_PR_protocol, cause: C.CauseProtocol_transfer_syntax_error},
	"protocol:abstract-syntax-error-reject":                      {causeGroup: C.Cause_PR_protocol, cause: C.CauseProtocol_abstract_syntax_error_reject},
	"protocol:abstract-syntax-error-ignore-and-notify":           {causeGroup: C.Cause_PR_protocol, cause: C.CauseProtocol_abstract_syntax_error_ignore_and_notify},
	"protocol:message-not-compatible-with-receiver-state":        {causeGroup: C.Cause_PR_protocol, cause: C.CauseProtocol_message_not_compatible_with_receiver_state},
	"protocol:semantic-error":                                    {causeGroup: C.Cause_PR_protocol, cause: C.CauseProtocol_semantic_error},
	"protocol:unspecified":                                       {causeGroup: C.Cause_PR_protocol, cause: C.CauseProtocol_unspecified},
	"protocol:abstract-syntax-error-falsely-constructed-message": {causeGroup: C.Cause_PR_protocol, cause: C.CauseProtocol_abstract_syntax_error_falsely_constructed_message},

	"transport:transport-resource-unavailable": {causeGroup: C.Cause_PR_transport, cause: C.CauseTransport_transport_resource_unavailable},
	"transport:unspecified":                    {causeGroup: C.Cause_PR_transport, cause: C.CauseTransport_unspecified},

	"radioNetwork:handover-desirable-for-radio-reasons":                            {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_handover_desirable_for_radio_reasons},
	"radioNetwork:time-critical-handover":                                          {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_time_critical_handover},
	"radioNetwork:resource-optimisation-handover":                                  {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_resource_optimisation_handover},
	"radioNetwork:reduce-load-in-serving-cell":                                     {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_reduce_load_in_serving_cell},
	"radioNetwork:partial-handover":                                                {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_partial_handover},
	"radioNetwork:unknown-new-enb-ue-x2ap-id":                                      {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_unknown_new_eNB_UE_X2AP_ID},
	"radioNetwork:unknown-old-enb-ue-x2ap-id":                                      {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_unknown_old_eNB_UE_X2AP_ID},
	"radioNetwork:unknown-pair-of-ue-x2ap-id":                                      {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_unknown_pair_of_UE_X2AP_ID},
	"radioNetwork:ho-target-not-allowed":                                           {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_ho_target_not_allowed},
	"radioNetwork:tx2relocoverall-expiry":                                          {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_tx2relocoverall_expiry},
	"radioNetwork:trelocprep-expiry":                                               {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_trelocprep_expiry},
	"radioNetwork:cell-not-available":                                              {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_cell_not_available},
	"radioNetwork:no-radio-resources-available-in-target-cell":                     {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_no_radio_resources_available_in_target_cell},
	"radioNetwork:invalid-mme-groupid":                                             {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_invalid_MME_GroupID},
	"radioNetwork:unknown-mme-code":                                                {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_unknown_MME_Code},
	"radioNetwork:encryption-and-or-integrity-protection-algorithms-not-supported": {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_encryption_and_or_integrity_protection_algorithms_not_supported},
	"radioNetwork:reportcharacteristicsempty":                                      {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_reportCharacteristicsEmpty},
	"radioNetwork:noreportperiodicity":                                             {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_noReportPeriodicity},
	"radioNetwork:existingMeasurementID":                                           {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_existingMeasurementID},
	"radioNetwork:unknown-enb-measurement-id":                                      {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_unknown_eNB_Measurement_ID},
	"radioNetwork:measurement-temporarily-not-available":                           {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_measurement_temporarily_not_available},
	"radioNetwork:unspecified":                                                     {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_unspecified},
	"radioNetwork:load-balancing":                                                  {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_load_balancing},
	"radioNetwork:handover-optimisation":                                           {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_handover_optimisation},
	"radioNetwork:value-out-of-allowed-range":                                      {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_value_out_of_allowed_range},
	"radioNetwork:multiple-E-RAB-ID-instances":                                     {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_multiple_E_RAB_ID_instances},
	"radioNetwork:switch-off-ongoing":                                              {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_switch_off_ongoing},
	"radioNetwork:not-supported-qci-value":                                         {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_not_supported_QCI_value},
	"radioNetwork:measurement-not-supported-for-the-object":                        {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_measurement_not_supported_for_the_object},
	"radioNetwork:tdcoverall-expiry":                                               {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_tDCoverall_expiry},
	"radioNetwork:tdcprep-expiry":                                                  {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_tDCprep_expiry},
	"radioNetwork:action-desirable-for-radio-reasons":                              {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_action_desirable_for_radio_reasons},
	"radioNetwork:reduce-load":                                                     {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_reduce_load},
	"radioNetwork:resource-optimisation":                                           {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_resource_optimisation},
	"radioNetwork:time-critical-action":                                            {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_time_critical_action},
	"radioNetwork:target-not-allowed":                                              {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_target_not_allowed},
	"radioNetwork:no-radio-resources-available":                                    {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_no_radio_resources_available},
	"radioNetwork:invalid-qos-combination":                                         {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_invalid_QoS_combination},
	"radioNetwork:encryption-algorithms-not-aupported":                             {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_encryption_algorithms_not_aupported},
	"radioNetwork:procedure-cancelled":                                             {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_procedure_cancelled},
	"radioNetwork:rrm-purpose":                                                     {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_rRM_purpose},
	"radioNetwork:improve-user-bit-rate":                                           {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_improve_user_bit_rate},
	"radioNetwork:user-inactivity":                                                 {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_user_inactivity},
	"radioNetwork:radio-connection-with-ue-lost":                                   {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_radio_connection_with_UE_lost},
	"radioNetwork:failure-in-the-radio-interface-procedure":                        {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_failure_in_the_radio_interface_procedure},
	"radioNetwork:bearer-option-not-supported":                                     {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_bearer_option_not_supported},
	"radioNetwork:mcg-mobility":                                                    {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_mCG_Mobility},
	"radioNetwork:scg-mobility":                                                    {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_sCG_Mobility},
	"radioNetwork:count-reaches-max-value":                                         {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_count_reaches_max_value},
	"radioNetwork:unknown-old-en-gnb-ue-x2ap-id":                                   {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_unknown_old_en_gNB_UE_X2AP_ID},
	"radioNetwork:pdcp-Overload":                                                   {causeGroup: C.Cause_PR_radioNetwork, cause: C.CauseRadioNetwork_pDCP_Overload},
}

var knownCausesToX2ResetPDUs = map[string][]byte{}

func prepareX2ResetPDUs(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int) error {
	packedBuffer := make([]C.uchar, maxAsn1PackedBufferSize)
	errorBuffer := make([]C.char, maxAsn1CodecMessageBufferSize)

	for k, cause := range knownCauses {
		var payloadSize = C.ulong(maxAsn1PackedBufferSize)
		if status := C.build_pack_x2reset_request(cause.causeGroup, C.int(cause.cause), &payloadSize, &packedBuffer[0], C.ulong(maxAsn1CodecMessageBufferSize), &errorBuffer[0]); !status {
			return fmt.Errorf("#x2_reset_known_causes.prepareX2ResetPDUs - failed to build and pack the reset message %s ", C.GoString(&errorBuffer[0]))
		}
		knownCausesToX2ResetPDUs[strings.ToLower(k)] = C.GoBytes(unsafe.Pointer(&packedBuffer[0]), C.int(payloadSize))
	}
	return nil
}

// KnownCausesToX2ResetPDU returns a packed x2 reset pdu with the specified cause (case insensitive match).
func KnownCausesToX2ResetPDU(cause string) ([]byte, bool) {
	v, ok := knownCausesToX2ResetPDUs[strings.ToLower(cause)]
	return v, ok
}

func init() {
	if err := prepareX2ResetPDUs(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize); err != nil {
		panic(err)
	}
}
