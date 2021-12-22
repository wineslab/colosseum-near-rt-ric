

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


#include "subscription_delete_request.hpp"
  
subscription_delete::subscription_delete(void){

  _name = "default";
  
  e2ap_pdu_obj = (E2N_E2AP_PDU_t * )calloc(1, sizeof(E2N_E2AP_PDU_t));
  assert(e2ap_pdu_obj != 0);

  initMsg = (E2N_InitiatingMessage_t * )calloc(1, sizeof(E2N_InitiatingMessage_t));
  assert(initMsg != 0);
  
  IE_array = (E2N_RICsubscriptionDeleteRequest_IEs_t *)calloc(NUM_SUBSCRIPTION_DELETE_IES, sizeof(E2N_RICsubscriptionDeleteRequest_IEs_t));
  assert(IE_array != 0);
  
  E2N_RICsubscriptionDeleteRequest_t * subscription_delete = &(initMsg->value.choice.RICsubscriptionDeleteRequest);
  for(int i = 0; i < NUM_SUBSCRIPTION_DELETE_IES; i++){
    ASN_SEQUENCE_ADD(&subscription_delete->protocolIEs, &(IE_array[i]));
  }
  
};



// Clear assigned protocolIE list from RIC indication IE container
subscription_delete::~subscription_delete(void){
    
  mdclog_write(MDCLOG_DEBUG, "Freeing subscription delete request object memory");
  E2N_RICsubscriptionDeleteRequest_t * subscription_delete = &(initMsg->value.choice.RICsubscriptionDeleteRequest);
  
  for(int i = 0; i < subscription_delete->protocolIEs.list.size; i++){
    subscription_delete->protocolIEs.list.array[i] = 0;
  }

  if (subscription_delete->protocolIEs.list.size > 0){
    free(subscription_delete->protocolIEs.list.array);
    subscription_delete->protocolIEs.list.count = 0;
    subscription_delete->protocolIEs.list.size = 0;
    subscription_delete->protocolIEs.list.array = 0;
  }
  
  free(IE_array);
  free(initMsg);
  e2ap_pdu_obj->choice.initiatingMessage = 0;

  ASN_STRUCT_FREE(asn_DEF_E2N_E2AP_PDU, e2ap_pdu_obj);
  mdclog_write(MDCLOG_DEBUG, "Freed subscription delete request object memory");
  

};


bool subscription_delete::encode_e2ap_subscription(unsigned char *buf, size_t *size,  subscription_helper &dinput){

  e2ap_pdu_obj->choice.initiatingMessage = initMsg;
  e2ap_pdu_obj->present = E2N_E2AP_PDU_PR_initiatingMessage;
  set_fields( dinput);

  initMsg->procedureCode = E2N_ProcedureCode_id_ricSubscriptionDelete;
  initMsg->criticality = E2N_Criticality_reject;
  initMsg->value.present = E2N_InitiatingMessage__value_PR_RICsubscriptionDeleteRequest;

  //xer_fprint(stdout, &asn_DEF_E2N_E2AP_PDU, e2ap_pdu_obj);
  
  int ret_constr = asn_check_constraints(&asn_DEF_E2N_E2AP_PDU, (void *) e2ap_pdu_obj, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(errbuf, errbuf_len);
    error_string = "Constraints failed for encoding subscription delete request. Reason = " + error_string;
    return false;
  }
  
  asn_enc_rval_t res = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2N_E2AP_PDU, e2ap_pdu_obj, buf, *size);
    
  if(res.encoded == -1){
    error_string.assign(strerror(errno));
    error_string = "Error encoding Subscription Delete Request. Reason = " + error_string;
    return false;
  }
  else {
    if(*size < res.encoded){
      std::stringstream ss;
      ss  <<"Error encoding Subscription Delete Request . Reason =  encoded pdu size " << res.encoded << " exceeds buffer size " << *size << std::endl;
      error_string = ss.str();
      res.encoded = -1;
      return false;
    }
  }
    
  *size = res.encoded;
  return true;
    
}


bool  subscription_delete::set_fields( subscription_helper &helper){
  unsigned int ie_index;
  
  ie_index = 0;
  E2N_RICsubscriptionDeleteRequest_IEs_t *ies_ricreq = &IE_array[ie_index];
  ies_ricreq->criticality = E2N_Criticality_reject;
  ies_ricreq->id = E2N_ProtocolIE_ID_id_RICrequestID;
  ies_ricreq->value.present = E2N_RICsubscriptionDeleteRequest_IEs__value_PR_RICrequestID;
  E2N_RICrequestID_t *ricrequest_ie = &ies_ricreq->value.choice.RICrequestID;
  ricrequest_ie->ricRequestorID = helper.get_request_id();
  ricrequest_ie->ricRequestSequenceNumber = helper.get_req_seq();


  
  ie_index = 1;
  E2N_RICsubscriptionDeleteRequest_IEs_t *ies_ranfunc = &IE_array[ie_index];
  ies_ranfunc->criticality = E2N_Criticality_reject;
  ies_ranfunc->id = E2N_ProtocolIE_ID_id_RANfunctionID;
  ies_ranfunc->value.present = E2N_RICsubscriptionDeleteRequest_IEs__value_PR_RANfunctionID;
  E2N_RANfunctionID_t *ranfunction_ie = &ies_ranfunc->value.choice.RANfunctionID;
  *ranfunction_ie = helper.get_function_id();

  
  return true;
};


   

bool  subscription_delete:: get_fields(E2N_InitiatingMessage_t * init_msg,  subscription_helper & dout)
{

  if (init_msg == 0){
    error_string = "Invalid reference for initiating message for get string";
    return false;
  }
  
  E2N_RICrequestID_t *requestid;
  E2N_RANfunctionID_t * ranfunctionid;
    
  for(int edx = 0; edx < init_msg->value.choice.RICsubscriptionDeleteRequest.protocolIEs.list.count; edx++) {
    E2N_RICsubscriptionDeleteRequest_IEs_t *memb_ptr = init_msg->value.choice.RICsubscriptionDeleteRequest.protocolIEs.list.array[edx];
    
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
    
  //asn_fprint(stdout, &asn_DEF_E2N_E2AP_PDU, e2pdu);
  }

  return true;
}



