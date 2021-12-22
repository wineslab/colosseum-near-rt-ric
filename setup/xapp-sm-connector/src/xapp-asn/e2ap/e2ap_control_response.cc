/*
==================================================================================

        Copyright (c) 2019-2020 AT&T Intellectual Property.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
 */
/*
 * ric_control_response.c
 *
 *  Created on: Jul 11, 2019
 *      Author: sjana, Ashwin Sridharan
 */

#include "e2ap_control_response.hpp"

// Set up the initiating message and also allocate protocolIEs in container
// Note : this bypasses requirement to use ASN_SEQUENCE_ADD. We can directly
// assign pointers to the array in ProtocolIE. However, this also leaves us on the
// hook to manually clear the memory

ric_control_response::ric_control_response(void){

	e2ap_pdu_obj = 0;
	e2ap_pdu_obj = (E2AP_PDU_t * )calloc(1, sizeof(E2AP_PDU_t));
	assert(e2ap_pdu_obj != 0);

	successMsg = 0;
	successMsg = (SuccessfulOutcome_t * )calloc(1, sizeof(SuccessfulOutcome_t));
	assert(successMsg != 0);

	successMsg->procedureCode = ProcedureCode_id_RICcontrol;
	successMsg->criticality = Criticality_reject;
	successMsg->value.present = SuccessfulOutcome__value_PR_RICcontrolAcknowledge;


	unsuccessMsg = 0;
	unsuccessMsg = (UnsuccessfulOutcome_t * )calloc(1, sizeof(UnsuccessfulOutcome_t));
	assert(unsuccessMsg != 0);


	unsuccessMsg->procedureCode = ProcedureCode_id_RICcontrol;
	unsuccessMsg->criticality = Criticality_reject;
	unsuccessMsg->value.present = UnsuccessfulOutcome__value_PR_RICcontrolFailure;

	IE_array = 0;
	IE_array = (RICcontrolAcknowledge_IEs_t *)calloc(NUM_CONTROL_ACKNOWLEDGE_IES, sizeof(RICcontrolAcknowledge_IEs_t));
	assert(IE_array != 0);

	RICcontrolAcknowledge_t * ric_acknowledge = &(successMsg->value.choice.RICcontrolAcknowledge);
	for(int i = 0; i < NUM_CONTROL_ACKNOWLEDGE_IES; i++){
		ASN_SEQUENCE_ADD(&(ric_acknowledge->protocolIEs), &(IE_array[i]));
	}


	IE_failure_array = 0;
	IE_failure_array = (RICcontrolFailure_IEs_t *)calloc(NUM_CONTROL_FAILURE_IES, sizeof(RICcontrolFailure_IEs_t));
	assert(IE_failure_array != 0);

	RICcontrolFailure_t * ric_failure = &(unsuccessMsg->value.choice.RICcontrolFailure);
	for(int i = 0; i < NUM_CONTROL_FAILURE_IES; i++){
		ASN_SEQUENCE_ADD(&(ric_failure->protocolIEs), &(IE_failure_array[i]));
	}

};


// Clear assigned protocolIE list from RIC control_request IE container
ric_control_response::~ric_control_response(void){

	mdclog_write(MDCLOG_DEBUG, "Freeing E2AP Control Response object memory");

	RICcontrolAcknowledge_t * ric_acknowledge = &(successMsg->value.choice.RICcontrolAcknowledge);
	for(int i  = 0; i < ric_acknowledge->protocolIEs.list.size; i++){
		ric_acknowledge->protocolIEs.list.array[i] = 0;
	}
	if (ric_acknowledge->protocolIEs.list.size > 0){
		free(ric_acknowledge->protocolIEs.list.array);
		ric_acknowledge->protocolIEs.list.array = 0;
		ric_acknowledge->protocolIEs.list.count = 0;
	}

	RICcontrolFailure_t * ric_failure = &(unsuccessMsg->value.choice.RICcontrolFailure);
	for(int i  = 0; i < ric_failure->protocolIEs.list.size; i++){
		ric_failure->protocolIEs.list.array[i] = 0;
	}
	if (ric_failure->protocolIEs.list.size > 0){
		free(ric_failure->protocolIEs.list.array);
		ric_failure->protocolIEs.list.array = 0;
		ric_failure->protocolIEs.list.count = 0;
	}

	free(IE_array);
	free(IE_failure_array);
	free(successMsg);
	free(unsuccessMsg);

	e2ap_pdu_obj->choice.initiatingMessage = 0;
	e2ap_pdu_obj->present = E2AP_PDU_PR_initiatingMessage;

	ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, e2ap_pdu_obj);
	mdclog_write(MDCLOG_DEBUG, "Freed E2AP Control Response object mempory");
}


bool ric_control_response::encode_e2ap_control_response(unsigned char *buf, size_t *size, ric_control_helper & dinput, bool is_success){

	bool res;
	if (is_success){
		res = set_fields(successMsg, dinput);
	}
	else{
		res = set_fields(unsuccessMsg, dinput);
	}

	if (!res){
		return false;
	}


	if (is_success){
		e2ap_pdu_obj->choice.successfulOutcome = successMsg;
		e2ap_pdu_obj->present = E2AP_PDU_PR_successfulOutcome ;
	}
	else{
		e2ap_pdu_obj->choice.unsuccessfulOutcome = unsuccessMsg;
		e2ap_pdu_obj->present = E2AP_PDU_PR_unsuccessfulOutcome ;

	}

	//xer_fprint(stdout, &asn_DEF_E2AP_PDU, e2ap_pdu_obj);

	int ret_constr = asn_check_constraints(&asn_DEF_E2AP_PDU, (void *) e2ap_pdu_obj, errbuf, &errbuf_len);
	if(ret_constr){
		error_string.assign(errbuf, errbuf_len);
		error_string = "Constraints failed for encoding control response. Reason = " + error_string;
		return false;
	}

	asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2AP_PDU, e2ap_pdu_obj, buf, *size);

	if(retval.encoded == -1){
		error_string.assign(strerror(errno));
		return false;
	}
	else {
		if(*size < retval.encoded){
			std::stringstream ss;
			ss  <<"Error encoding E2AP Control response . Reason =  encoded pdu size " << retval.encoded << " exceeds buffer size " << *size << std::endl;
			error_string = ss.str();
			return false;
		}
	}

	*size = retval.encoded;
	return true;

}

bool ric_control_response::set_fields(SuccessfulOutcome_t *successMsg, ric_control_helper &dinput){
	unsigned int ie_index;

	if (successMsg == 0){
		error_string = "Invalid reference for E2AP Control Acknowledge in set_fields";
		return false;
	}

	// for(i = 0; i < NUM_CONTROL_ACKNOWLEDGE_IES;i++){
	//   memset(&(IE_array[i]), 0, sizeof(RICcontrolAcknowledge_IEs_t));
	// }

	//RICcontrolAcknowledge_t * ric_acknowledge = &(successMsg->value.choice.RICcontrolAcknowledge);
	//ric_acknowledge->protocolIEs.list.count = 0;

	ie_index = 0;
	RICcontrolAcknowledge_IEs_t *ies_ricreq = &IE_array[ie_index];
	ies_ricreq->criticality = Criticality_reject;
	ies_ricreq->id = ProtocolIE_ID_id_RICrequestID;
	ies_ricreq->value.present = RICcontrolAcknowledge_IEs__value_PR_RICrequestID;
	RICrequestID_t *ricrequest_ie = &ies_ricreq->value.choice.RICrequestID;
	ricrequest_ie->ricRequestorID = dinput.req_id;
	//ricrequest_ie->ricRequestSequenceNumber = dinput.req_seq_no;
	//ASN_SEQUENCE_ADD(&(ric_acknowledge->protocolIEs), ies_ricreq);

	ie_index = 1;
	RICcontrolAcknowledge_IEs_t *ies_ranfunc = &IE_array[ie_index];
	ies_ranfunc->criticality = Criticality_reject;
	ies_ranfunc->id = ProtocolIE_ID_id_RANfunctionID;
	ies_ranfunc->value.present = RICcontrolAcknowledge_IEs__value_PR_RANfunctionID;
	RANfunctionID_t *ranfunction_ie = &ies_ranfunc->value.choice.RANfunctionID;
	*ranfunction_ie = dinput.func_id;
	//ASN_SEQUENCE_ADD(&(ric_acknowledge->protocolIEs), ies_ranfunc);

	// ie_index = 2;
	// RICcontrolAcknowledge_IEs_t *ies_riccallprocessid = &IE_array[ie_index];
	// ies_riccallprocessid->criticality = Criticality_reject;
	// ies_riccallprocessid->id = ProtocolIE_ID_id_RICcallProcessID;
	// ies_riccallprocessid->value.present = RICcontrolAcknowledge_IEs__value_PR_RICcallProcessID;
	// RICcallProcessID_t *riccallprocessid_ie = &ies_riccallprocessid->value.choice.RICcallProcessID;
	// riccallprocessid_ie->buf = dinput.call_process_id;
	// riccallprocessid_ie->size = dinput.call_process_id_size;
	// ASN_SEQUENCE_ADD(&(ric_acknowledge->protocolIEs), ies_riccallprocessid);

	ie_index = 2;
	RICcontrolAcknowledge_IEs_t *ies_ric_cause = &IE_array[ie_index];
	ies_ric_cause->criticality = Criticality_reject;
	ies_ric_cause->id = ProtocolIE_ID_id_RICcontrolStatus;
	ies_ric_cause->value.present = RICcontrolAcknowledge_IEs__value_PR_RICcontrolStatus;
	ies_ric_cause->value.choice.RICcontrolStatus = dinput.control_status;
	//ASN_SEQUENCE_ADD(&(ric_acknowledge->protocolIEs), ies_ric_cause);

	return true;

};

bool ric_control_response::set_fields(UnsuccessfulOutcome_t *unsuccessMsg, ric_control_helper &dinput){
	unsigned int ie_index;

	if (unsuccessMsg == 0){
		error_string = "Invalid reference for E2AP Control Failure in set_fields";
		return false;
	}

	// for(i = 0; i < NUM_CONTROL_FAILURE_IES;i++){
	//   memset(&(IE_failure_array[i]), 0, sizeof(RICcontrolFailure_IEs_t));
	// }

	//RICcontrolFailure_t * ric_failure = &(unsuccessMsg->value.choice.RICcontrolFailure);
	//ric_failure->protocolIEs.list.count = 0;

	ie_index = 0;
	RICcontrolFailure_IEs_t *ies_ricreq = &IE_failure_array[ie_index];
	ies_ricreq->criticality = Criticality_reject;
	ies_ricreq->id = ProtocolIE_ID_id_RICrequestID;
	ies_ricreq->value.present = RICcontrolFailure_IEs__value_PR_RICrequestID;
	RICrequestID_t *ricrequest_ie = &(ies_ricreq->value.choice.RICrequestID);
	ricrequest_ie->ricRequestorID = dinput.req_id;
	//ricrequest_ie->ricRequestSequenceNumber = dinput.req_seq_no;
	//ASN_SEQUENCE_ADD(&(ric_failure->protocolIEs), ies_ricreq);

	ie_index = 1;
	RICcontrolFailure_IEs_t *ies_ranfunc = &IE_failure_array[ie_index];
	ies_ranfunc->criticality = Criticality_reject;
	ies_ranfunc->id = ProtocolIE_ID_id_RANfunctionID;
	ies_ranfunc->value.present = RICcontrolFailure_IEs__value_PR_RANfunctionID;
	RANfunctionID_t *ranfunction_ie = &(ies_ranfunc->value.choice.RANfunctionID);
	*ranfunction_ie = dinput.func_id;
	//ASN_SEQUENCE_ADD(&(ric_failure->protocolIEs), ies_ranfunc);

	// ie_index = 2;
	// RICcontrolFailure_IEs_t *ies_riccallprocessid = &IE_failure_array[i];
	// ies_riccallprocessid->criticality = Criticality_reject;
	// ies_riccallprocessid->id = ProtocolIE_ID_id_RICcallProcessID;
	// ies_riccallprocessid->value.present = RICcontrolFailure_IEs__value_PR_RICcallProcessID;
	// RICcallProcessID_t *riccallprocessid_ie = &(ies_riccallprocessid->value.choice.RICcallProcessID);
	// riccallprocessid_ie->buf = dinput.call_process_id;
	// riccallprocessid_ie->size = dinput.call_process_id_size;
	// ASN_SEQUENCE_ADD(&(ric_failure->protocolIEs), ies_riccallprocessid);

	ie_index = 2;
	RICcontrolFailure_IEs_t *ies_ric_cause = &IE_failure_array[ie_index];
	ies_ric_cause->criticality = Criticality_ignore;
	ies_ric_cause->id = ProtocolIE_ID_id_Cause;
	ies_ric_cause->value.present = RICcontrolFailure_IEs__value_PR_Cause;
	Cause_t * ric_cause = &(ies_ric_cause->value.choice.Cause);
	ric_cause->present = (Cause_PR)dinput.cause;

	switch(dinput.cause){
	case Cause_PR_ricService:
		ric_cause->choice.ricService = dinput.sub_cause;
		break;
	case Cause_PR_transport:
		ric_cause->choice.transport = dinput.sub_cause;
		break;
	case Cause_PR_protocol:
		ric_cause->choice.protocol= dinput.sub_cause;
		break;
	case Cause_PR_misc:
		ric_cause->choice.misc = dinput.sub_cause;
		break;
	case Cause_PR_ricRequest:
		ric_cause->choice.ricRequest = dinput.sub_cause;
		break;
	default:
		std::cout <<"Error ! Illegal cause enum" << dinput.cause << std::endl;
		return false;
	}

	//ASN_SEQUENCE_ADD(&(ric_failure->protocolIEs), ies_ric_cause);
	return true;

};




bool ric_control_response:: get_fields(SuccessfulOutcome_t * success_msg,  ric_control_helper &dout)
{
	if (success_msg == 0){
		error_string = "Invalid reference for E2AP Control Acknowledge message in get_fields";
		return false;
	}


	for(int edx = 0; edx < success_msg->value.choice.RICcontrolAcknowledge.protocolIEs.list.count; edx++) {
		RICcontrolAcknowledge_IEs_t *memb_ptr = success_msg->value.choice.RICcontrolAcknowledge.protocolIEs.list.array[edx];

		switch(memb_ptr->id)
		{

		case (ProtocolIE_ID_id_RICcallProcessID):
  			dout.call_process_id =  memb_ptr->value.choice.RICcallProcessID.buf;
		dout.call_process_id_size = memb_ptr->value.choice.RICcallProcessID.size;
		break;

		case (ProtocolIE_ID_id_RICrequestID):
  			dout.req_id = memb_ptr->value.choice.RICrequestID.ricRequestorID;
		//dout.req_seq_no = memb_ptr->value.choice.RICrequestID.ricRequestSequenceNumber;
		break;

		case (ProtocolIE_ID_id_RANfunctionID):
  			dout.func_id = memb_ptr->value.choice.RANfunctionID;
		break;

		case (ProtocolIE_ID_id_Cause):
  			dout.control_status = memb_ptr->value.choice.RICcontrolStatus;
		break;

		}

	}

	return true;

}


bool ric_control_response:: get_fields(UnsuccessfulOutcome_t * unsuccess_msg,  ric_control_helper &dout)
{
	if (unsuccess_msg == 0){
		error_string = "Invalid reference for E2AP Control Failure message in get_fields";
		return false;
	}


	for(int edx = 0; edx < unsuccess_msg->value.choice.RICcontrolFailure.protocolIEs.list.count; edx++) {
		RICcontrolFailure_IEs_t *memb_ptr = unsuccess_msg->value.choice.RICcontrolFailure.protocolIEs.list.array[edx];

		switch(memb_ptr->id)
		{

		case (ProtocolIE_ID_id_RICcallProcessID):
  			dout.call_process_id =  memb_ptr->value.choice.RICcallProcessID.buf;
		dout.call_process_id_size = memb_ptr->value.choice.RICcallProcessID.size;
		break;

		case (ProtocolIE_ID_id_RICrequestID):
  			dout.req_id = memb_ptr->value.choice.RICrequestID.ricRequestorID;
		//dout.req_seq_no = memb_ptr->value.choice.RICrequestID.ricRequestSequenceNumber;
		break;

		case (ProtocolIE_ID_id_RANfunctionID):
  			dout.func_id = memb_ptr->value.choice.RANfunctionID;
		break;


		case (ProtocolIE_ID_id_Cause):
  			dout.cause = memb_ptr->value.choice.Cause.present;
		switch(dout.cause){
		case  Cause_PR_ricService :
			dout.sub_cause = memb_ptr->value.choice.Cause.choice.ricService;
			break;

		case Cause_PR_transport :
			dout.sub_cause = memb_ptr->value.choice.Cause.choice.transport;
			break;

		case  Cause_PR_protocol :
			dout.sub_cause = memb_ptr->value.choice.Cause.choice.protocol;
			break;

		case Cause_PR_misc :
			dout.sub_cause = memb_ptr->value.choice.Cause.choice.misc;
			break;

		case Cause_PR_ricRequest :
			dout.sub_cause = memb_ptr->value.choice.Cause.choice.ricRequest;
			break;

		default:
			dout.sub_cause = -1;
			break;
		}

		default:
			break;
		}

	}

	return true;

}

