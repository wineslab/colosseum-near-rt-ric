/*
==================================================================================
        Copyright (c) 2018-2019 AT&T Intellectual Property.

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


#include "subscription_delete_response.hpp"

/* The xAPP need only worry about the get_fields from a response, since it does
not generate a response. Generating response however is included to support testing. 
*/


// Primarly for generation
subscription_delete_response::subscription_delete_response(void){

  e2ap_pdu_obj = 0;
  e2ap_pdu_obj = (E2N_E2AP_PDU_t *)calloc(1, sizeof(E2N_E2AP_PDU_t));
  assert(e2ap_pdu_obj != 0);

  successMsg = 0;
  successMsg = (E2N_SuccessfulOutcome_t *)calloc(1, sizeof(E2N_SuccessfulOutcome_t));
  assert(successMsg != 0);

  unsuccessMsg = 0;
  unsuccessMsg = (E2N_UnsuccessfulOutcome_t *)calloc(1, sizeof(E2N_UnsuccessfulOutcome_t));
  assert(unsuccessMsg != 0);

  IE_array = 0;
  IE_array = (E2N_RICsubscriptionDeleteResponse_IEs_t *)calloc(NUM_SUBSCRIPTION_DELETE_RESPONSE_IES, sizeof(E2N_RICsubscriptionDeleteResponse_IEs_t));
  assert(IE_array != 0);

  IE_Failure_array = 0;
  IE_Failure_array = (E2N_RICsubscriptionDeleteFailure_IEs_t *)calloc(NUM_SUBSCRIPTION_DELETE_FAILURE_IES, sizeof(E2N_RICsubscriptionDeleteFailure_IEs_t));
  assert(IE_Failure_array != 0);

  
   
};

  

// Clear assigned protocolIE list from E2N_RIC indication IE container
subscription_delete_response::~subscription_delete_response(void){

  mdclog_write(MDCLOG_DEBUG, "Freeing subscription delete response memory");
  E2N_RICsubscriptionDeleteResponse_t * ric_subscription_delete_response = &(successMsg->value.choice.RICsubscriptionDeleteResponse);
  
  for(unsigned int i = 0; i < ric_subscription_delete_response->protocolIEs.list.size ; i++){
    ric_subscription_delete_response->protocolIEs.list.array[i] = 0;
  }

  
  E2N_RICsubscriptionDeleteFailure_t * ric_subscription_failure = &(unsuccessMsg->value.choice.RICsubscriptionDeleteFailure);
  for(unsigned int i = 0; i < ric_subscription_failure->protocolIEs.list.size; i++){
    ric_subscription_failure->protocolIEs.list.array[i] = 0;
  }

  free(IE_array);
  free(IE_Failure_array);

  ASN_STRUCT_FREE(asn_DEF_E2N_SuccessfulOutcome, successMsg);

  ASN_STRUCT_FREE(asn_DEF_E2N_UnsuccessfulOutcome, unsuccessMsg);
  
  e2ap_pdu_obj->choice.successfulOutcome = NULL;
  e2ap_pdu_obj->choice.unsuccessfulOutcome = NULL;

  ASN_STRUCT_FREE(asn_DEF_E2N_E2AP_PDU, e2ap_pdu_obj);
  mdclog_write(MDCLOG_DEBUG, "Freed subscription delete response memory");

};


bool subscription_delete_response::encode_e2ap_subscription_delete_response(unsigned char *buf, size_t *size,  subscription_response_helper &dinput, bool is_success){

  bool res;
 
  if(is_success){
    res = set_fields(successMsg, dinput);
    if (!res){
      return false;
    }
    e2ap_pdu_obj->present =  E2N_E2AP_PDU_PR_successfulOutcome;
    e2ap_pdu_obj->choice.successfulOutcome = successMsg;
  }
  else{
    res = set_fields(unsuccessMsg, dinput);
    if(! res){
      return false;
    }
    e2ap_pdu_obj->present = E2N_E2AP_PDU_PR_unsuccessfulOutcome;
    e2ap_pdu_obj->choice.unsuccessfulOutcome = unsuccessMsg;
  }
    

  int ret_constr = asn_check_constraints(&asn_DEF_E2N_E2AP_PDU, (void *) e2ap_pdu_obj, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(errbuf, errbuf_len);
    return false;
  }

  //xer_fprint(stdout, &asn_DEF_E2AP_PDU, e2ap_pdu_obj);
  
  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2N_E2AP_PDU, e2ap_pdu_obj, buf, *size);
    
  if(retval.encoded == -1){
    error_string.assign(strerror(errno));
    error_string = "Error encoding subcription delete response. Reason = " + error_string;
    return false;
  }
  else {
    if(*size < retval.encoded){
      std::stringstream ss;
      ss  <<"Error encoding Subscription Delete Response . Reason =  encoded pdu size " << retval.encoded << " exceeds buffer size " << *size << std::endl;
      error_string = ss.str();
      retval.encoded = -1;
      return false;
    }
  }
    
  *size = retval.encoded;
  return true;
    
}
  
bool  subscription_delete_response::set_fields(E2N_SuccessfulOutcome_t *success, subscription_response_helper &helper){

  if (success == 0){
    error_string = "Invalid reference to success message in set fields  subscription delete response";
    return false;
  }
  
  unsigned int ie_index;

  success->procedureCode = E2N_ProcedureCode_id_ricSubscriptionDelete;
  success->criticality = E2N_Criticality_reject;
  success->value.present = E2N_SuccessfulOutcome__value_PR_RICsubscriptionDeleteResponse;   
 
  E2N_RICsubscriptionDeleteResponse_t * subscription_delete_response = &(success->value.choice.RICsubscriptionDeleteResponse);
  subscription_delete_response->protocolIEs.list.count = 0;
  
  ie_index = 0;
  E2N_RICsubscriptionDeleteResponse_IEs_t *ies_ricreq = &IE_array[ie_index];
  
  ies_ricreq->criticality = E2N_Criticality_reject;
  ies_ricreq->id = E2N_ProtocolIE_ID_id_RICrequestID;
  ies_ricreq->value.present = E2N_RICsubscriptionDeleteResponse_IEs__value_PR_RICrequestID;
  E2N_RICrequestID_t *ricrequest_ie = &ies_ricreq->value.choice.RICrequestID;
  ricrequest_ie->ricRequestorID = helper.get_request_id();
  ricrequest_ie->ricRequestSequenceNumber = helper.get_req_seq();
  ASN_SEQUENCE_ADD(&subscription_delete_response->protocolIEs, ies_ricreq);

  
  ie_index = 1;
  E2N_RICsubscriptionDeleteResponse_IEs_t *ies_ranfunc = &IE_array[ie_index];
  ies_ranfunc->criticality = E2N_Criticality_reject;
  ies_ranfunc->id = E2N_ProtocolIE_ID_id_RANfunctionID;
  ies_ranfunc->value.present = E2N_RICsubscriptionDeleteResponse_IEs__value_PR_RANfunctionID;
  E2N_RANfunctionID_t *ranfunction_ie = &ies_ranfunc->value.choice.RANfunctionID;
  *ranfunction_ie = helper.get_function_id();
  ASN_SEQUENCE_ADD(&subscription_delete_response->protocolIEs, ies_ranfunc);

  return true;
 
	
}

bool subscription_delete_response:: get_fields(E2N_SuccessfulOutcome_t * success_msg,  subscription_response_helper & dout)
{

  if (success_msg == 0){
    error_string = "Invalid reference to success message inn get fields subscription delete response";
    return false;
  }
  
  E2N_RICrequestID_t *requestid;
  E2N_RANfunctionID_t * ranfunctionid;
  
  for(int edx = 0; edx < success_msg->value.choice.RICsubscriptionDeleteResponse.protocolIEs.list.count; edx++) {
    E2N_RICsubscriptionDeleteResponse_IEs_t *memb_ptr = success_msg->value.choice.RICsubscriptionDeleteResponse.protocolIEs.list.array[edx];
    
    switch(memb_ptr->id)
      {
      case (E2N_ProtocolIE_ID_id_RICrequestID):
	requestid = &memb_ptr->value.choice.RICrequestID;
	dout.set_request(requestid->ricRequestorID, requestid->ricRequestSequenceNumber);
	break;
	  
      case (E2N_ProtocolIE_ID_id_RANfunctionID):
	ranfunctionid = &memb_ptr->value.choice.RANfunctionID;
	dout.set_function_id(*ranfunctionid);
	break;
      }
    
  }
  
  return true;
  //asn_fprint(stdout, &asn_DEF_E2AP_PDU, e2pdu);
}


bool subscription_delete_response::set_fields(E2N_UnsuccessfulOutcome_t *unsuccess, subscription_response_helper &helper){

  if (unsuccess == 0){
    error_string = "Invalid reference to unsuccess message in set fields  subscription delete response";
    return false;
  }
  
  unsigned int ie_index;

  unsuccess->procedureCode = E2N_ProcedureCode_id_ricSubscriptionDelete;
  unsuccess->criticality = E2N_Criticality_reject;
  unsuccess->value.present = E2N_UnsuccessfulOutcome__value_PR_RICsubscriptionDeleteFailure;

  E2N_RICsubscriptionDeleteFailure_t * ric_subscription_failure = &(unsuccess->value.choice.RICsubscriptionDeleteFailure);
  ric_subscription_failure->protocolIEs.list.count = 0;
  
  ie_index = 0;
  E2N_RICsubscriptionDeleteFailure_IEs_t *ies_ricreq = &IE_Failure_array[ie_index];
    
  ies_ricreq->criticality = E2N_Criticality_reject;
  ies_ricreq->id = E2N_ProtocolIE_ID_id_RICrequestID;
  ies_ricreq->value.present = E2N_RICsubscriptionDeleteFailure_IEs__value_PR_RICrequestID;
  E2N_RICrequestID_t *ricrequest_ie = &ies_ricreq->value.choice.RICrequestID;
  ricrequest_ie->ricRequestorID = helper.get_request_id();
  ricrequest_ie->ricRequestSequenceNumber = helper.get_req_seq();
  ASN_SEQUENCE_ADD(&ric_subscription_failure->protocolIEs, ies_ricreq);
  
  ie_index = 1;
  E2N_RICsubscriptionDeleteFailure_IEs_t *ies_ranfunc = &IE_Failure_array[ie_index];
  ies_ranfunc->criticality = E2N_Criticality_reject;
  ies_ranfunc->id = E2N_ProtocolIE_ID_id_RANfunctionID;
  ies_ranfunc->value.present = E2N_RICsubscriptionDeleteFailure_IEs__value_PR_RANfunctionID;
  E2N_RANfunctionID_t *ranfunction_ie = &ies_ranfunc->value.choice.RANfunctionID;
  *ranfunction_ie = helper.get_function_id();
  ASN_SEQUENCE_ADD(&ric_subscription_failure->protocolIEs, ies_ranfunc);
    

  return true;
    
}

bool  subscription_delete_response:: get_fields(E2N_UnsuccessfulOutcome_t * unsuccess_msg,  subscription_response_helper & dout)
{

  if (unsuccess_msg == 0){
    error_string = "Invalid reference to unsuccess message in get fields  subscription delete response";
    return false;
  }
  
  E2N_RICrequestID_t *requestid;
  E2N_RANfunctionID_t * ranfunctionid;
    
  for(int edx = 0; edx < unsuccess_msg->value.choice.RICsubscriptionDeleteFailure.protocolIEs.list.count; edx++) {
    E2N_RICsubscriptionDeleteFailure_IEs_t *memb_ptr = unsuccess_msg->value.choice.RICsubscriptionDeleteFailure.protocolIEs.list.array[edx];
    
    switch(memb_ptr->id)
      {
      case (E2N_ProtocolIE_ID_id_RICrequestID):
	requestid = &memb_ptr->value.choice.RICrequestID;
	dout.set_request(requestid->ricRequestorID, requestid->ricRequestSequenceNumber);
	break;
	  
      case (E2N_ProtocolIE_ID_id_RANfunctionID):
	ranfunctionid = &memb_ptr->value.choice.RANfunctionID;
	dout.set_function_id(*ranfunctionid);
	break;
	
      }
    
  }

  return true;
  //asn_fprint(stdout, &asn_DEF_E2AP_PDU, e2pdu);
}



