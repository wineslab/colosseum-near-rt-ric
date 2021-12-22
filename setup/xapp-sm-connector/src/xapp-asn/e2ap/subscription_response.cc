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


#include "subscription_response.hpp"

/* The xAPP need only worry about the get_fields from a response, since it does
not generate a response. Generating response however is included to support testing. 
*/


// Primarly for generation
subscription_response::subscription_response(void){

  e2ap_pdu_obj = 0;
  e2ap_pdu_obj = (E2AP_PDU_t *)calloc(1, sizeof(E2AP_PDU_t));
  assert(e2ap_pdu_obj != 0);

  successMsg = 0;
  successMsg = (SuccessfulOutcome_t *)calloc(1, sizeof(SuccessfulOutcome_t));
  assert(successMsg != 0);

  unsuccessMsg = 0;
  unsuccessMsg = (UnsuccessfulOutcome_t *)calloc(1, sizeof(UnsuccessfulOutcome_t));
  assert(unsuccessMsg != 0);

  IE_array = 0;
  IE_array = (RICsubscriptionResponse_IEs_t *)calloc(NUM_SUBSCRIPTION_RESPONSE_IES, sizeof(RICsubscriptionResponse_IEs_t));
  assert(IE_array != 0);

  IE_Failure_array = 0;
  IE_Failure_array = (RICsubscriptionFailure_IEs_t *)calloc(NUM_SUBSCRIPTION_FAILURE_IES, sizeof(RICsubscriptionFailure_IEs_t));
  assert(IE_Failure_array != 0);

  ie_admitted_list = 0;
  ie_admitted_list = (RICaction_Admitted_ItemIEs_t *)calloc(INITIAL_RESPONSE_LIST_SIZE, sizeof(RICaction_Admitted_ItemIEs_t));
  assert(ie_admitted_list != 0);
  ie_admitted_list_size = INITIAL_RESPONSE_LIST_SIZE;
  
  ie_not_admitted_list = 0;
  ie_not_admitted_list = (RICaction_NotAdmitted_ItemIEs_t *)calloc(INITIAL_RESPONSE_LIST_SIZE, sizeof(RICaction_NotAdmitted_ItemIEs_t));
  assert(ie_not_admitted_list != 0);
  ie_not_admitted_list_size = INITIAL_RESPONSE_LIST_SIZE;




  
};

  

// Clear assigned protocolIE list from RIC indication IE container
subscription_response::~subscription_response(void){

  mdclog_write(MDCLOG_DEBUG, "Freeing subscription response memory");
  RICaction_Admitted_List_t * response_admitted_list = (RICaction_Admitted_List_t *) &(IE_array[2].value.choice.RICaction_Admitted_List);
  
  for(int i = 0; i < response_admitted_list->list.size; i++){
    response_admitted_list->list.array[i] = 0;
  }

  
  if (response_admitted_list->list.size > 0){
    free(response_admitted_list->list.array);
    response_admitted_list->list.array = 0;
    response_admitted_list->list.size = 0;
    response_admitted_list->list.count = 0;
  }

  
  RICaction_NotAdmitted_List_t * response_not_admitted_list = &(IE_array[3].value.choice.RICaction_NotAdmitted_List);
  for(int i = 0; i < response_not_admitted_list->list.size; i++){
    response_not_admitted_list->list.array[i] = 0;
  }
  
  if (response_not_admitted_list->list.size > 0){
    free(response_not_admitted_list->list.array);
    response_not_admitted_list->list.array = 0;
    response_not_admitted_list->list.size = 0;
    response_not_admitted_list->list.count = 0;
  }
     
  RICsubscriptionResponse_t * ric_subscription_response = &(successMsg->value.choice.RICsubscriptionResponse);
  for(int i = 0; i < ric_subscription_response->protocolIEs.list.size ; i++){
    ric_subscription_response->protocolIEs.list.array[i] = 0;
  }
  
  if (ric_subscription_response->protocolIEs.list.size > 0){
    free(ric_subscription_response->protocolIEs.list.array);
    ric_subscription_response->protocolIEs.list.array = 0;
    ric_subscription_response->protocolIEs.list.size = 0;
    ric_subscription_response->protocolIEs.list.count = 0;
  }
  
  
  RICaction_NotAdmitted_List_t * failure_not_admitted_list = &(IE_Failure_array[2].value.choice.RICaction_NotAdmitted_List);
  for(int i = 0; i < failure_not_admitted_list->list.size; i++){
    failure_not_admitted_list->list.array[i] = 0;
  }

  if ( failure_not_admitted_list->list.size > 0){
    free( failure_not_admitted_list->list.array);
    failure_not_admitted_list->list.array = 0;
    failure_not_admitted_list->list.size = 0;
    failure_not_admitted_list->list.count = 0;
  }
  
     
  RICsubscriptionFailure_t * ric_subscription_failure = &(unsuccessMsg->value.choice.RICsubscriptionFailure);
  for(int i = 0; i < ric_subscription_failure->protocolIEs.list.size; i++){
    ric_subscription_failure->protocolIEs.list.array[i] = 0;
  }
  
  if (  ric_subscription_failure->protocolIEs.list.size > 0){
    free(ric_subscription_failure->protocolIEs.list.array);
    ric_subscription_failure->protocolIEs.list.array = 0;
    ric_subscription_failure->protocolIEs.list.size = 0;
    ric_subscription_failure->protocolIEs.list.count = 0;
  }


  free(ie_admitted_list);  
  free(ie_not_admitted_list);
  free(IE_Failure_array);
  free(IE_array);

  
  ASN_STRUCT_FREE(asn_DEF_SuccessfulOutcome, successMsg);
  ASN_STRUCT_FREE(asn_DEF_UnsuccessfulOutcome, unsuccessMsg);

  
  e2ap_pdu_obj->choice.initiatingMessage = NULL;
  e2ap_pdu_obj->present = E2AP_PDU_PR_NOTHING;
  ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, e2ap_pdu_obj);

  mdclog_write(MDCLOG_DEBUG, "Freed subscription response memory ");

  
};


bool subscription_response::encode_e2ap_subscription_response(unsigned char *buf, size_t *size, subscription_response_helper &dinput, bool is_success){
  

  if(is_success){
    set_fields_success(dinput);
    e2ap_pdu_obj->present =  E2AP_PDU_PR_successfulOutcome;
    e2ap_pdu_obj->choice.successfulOutcome = successMsg;

    successMsg->procedureCode = ProcedureCode_id_RICsubscription;
    successMsg->criticality = Criticality_reject;
    successMsg->value.present = SuccessfulOutcome__value_PR_RICsubscriptionResponse;
 
  }
  else{
    set_fields_unsuccess(dinput);
    e2ap_pdu_obj->present = E2AP_PDU_PR_unsuccessfulOutcome;
    e2ap_pdu_obj->choice.unsuccessfulOutcome = unsuccessMsg;

    unsuccessMsg->procedureCode = ProcedureCode_id_RICsubscription;
    unsuccessMsg->criticality = Criticality_reject;
    unsuccessMsg->value.present = UnsuccessfulOutcome__value_PR_RICsubscriptionFailure;

  }
    
  //xer_fprint(stdout, &asn_DEF_E2AP_PDU, e2ap_pdu_obj);

  int ret_constr = asn_check_constraints(&asn_DEF_E2AP_PDU, (void *) e2ap_pdu_obj, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(errbuf, errbuf_len);
    error_string = "Constraints failed for encoding subscription response. Reason = " + error_string;
    return false;
  }
    
  asn_enc_rval_t res = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2AP_PDU, e2ap_pdu_obj, buf, *size);
    
  if(res.encoded == -1){
    std::cout <<"Error encoding PDU. Reason =" << strerror(errno) << std::endl;
    return false;
  }
  else {
    if(*size < res.encoded){
      fprintf(stderr,  "Buffer assigned too small to encode: %s",(char *)(asn_DEF_E2AP_PDU.name));
      res.encoded = -1;
      return false;
    }
  }
    
  *size = res.encoded;
  return true;
    
}
  
void subscription_response::set_fields_success(subscription_response_helper &helper){

  int ie_index;

  RICsubscriptionResponse_t * subscription_response = &(successMsg->value.choice.RICsubscriptionResponse);
  //reset list count ..
  subscription_response->protocolIEs.list.count = 0;
    
  ie_index = 0;
  RICsubscriptionResponse_IEs_t *ies_ricreq = &IE_array[ie_index];
  
  ies_ricreq->criticality = Criticality_reject;
  ies_ricreq->id = ProtocolIE_ID_id_RICrequestID;
  ies_ricreq->value.present = RICsubscriptionResponse_IEs__value_PR_RICrequestID;
  RICrequestID_t *ricrequest_ie = &ies_ricreq->value.choice.RICrequestID;
  ricrequest_ie->ricRequestorID = helper.get_request_id();
 // ricrequest_ie->ricRequestSequenceNumber = helper.get_req_seq();
  ASN_SEQUENCE_ADD(&subscription_response->protocolIEs, &(IE_array[ie_index]));  

  
  ie_index = 1;
  RICsubscriptionResponse_IEs_t *ies_ranfunc = &IE_array[ie_index];
  ies_ranfunc->criticality = Criticality_reject;
  ies_ranfunc->id = ProtocolIE_ID_id_RANfunctionID;
  ies_ranfunc->value.present = RICsubscriptionResponse_IEs__value_PR_RANfunctionID;
  RANfunctionID_t *ranfunction_ie = &ies_ranfunc->value.choice.RANfunctionID;
  *ranfunction_ie = helper.get_function_id();
  ASN_SEQUENCE_ADD(&subscription_response->protocolIEs, &(IE_array[ie_index]));

  
  ie_index = 2;
  RICsubscriptionResponse_IEs_t *ies_admitted_actid = &IE_array[ie_index];
  ies_admitted_actid->criticality = Criticality_reject;
  ies_admitted_actid->id = ProtocolIE_ID_id_RICactions_Admitted;
  RICaction_Admitted_List_t *ric_admitted_actions_ie = &ies_admitted_actid->value.choice.RICaction_Admitted_List;
  ric_admitted_actions_ie->list.count = 0;
  std::vector<ActionResponse> * ref_admitted_action_array = helper.get_admitted_list();  

  if(ref_admitted_action_array->size() ==  0){
    ies_admitted_actid->value.present =  RICsubscriptionResponse_IEs__value_PR_NOTHING;
  }
  else{  
    ies_admitted_actid->value.present =  RICsubscriptionResponse_IEs__value_PR_RICaction_Admitted_List;

    // resize memory ?
    if (ref_admitted_action_array->size() >= ie_admitted_list_size){
      ie_admitted_list_size = 2 * ref_admitted_action_array->size();
      free(ie_admitted_list);
      ie_admitted_list = (RICaction_Admitted_ItemIEs_t *)calloc(ie_admitted_list_size, sizeof(RICaction_Admitted_ItemIEs_t));
      assert(ie_admitted_list != 0);
    };
  

    for(unsigned int i = 0; i < ref_admitted_action_array->size(); i ++){
      ie_admitted_list[i].criticality = Criticality_ignore;
      ie_admitted_list[i].id = ProtocolIE_ID_id_RICaction_Admitted_Item ;
      ie_admitted_list[i].value.present = RICaction_Admitted_ItemIEs__value_PR_RICaction_Admitted_Item;
      ie_admitted_list[i].value.choice.RICaction_Admitted_Item.ricActionID = (*ref_admitted_action_array)[i].get_id();
      ASN_SEQUENCE_ADD(ric_admitted_actions_ie, &(ie_admitted_list[i]));
    }
  }  
  ASN_SEQUENCE_ADD(&subscription_response->protocolIEs, &(IE_array[ie_index]));

  // optional IE : add only if non-zero list 
  ie_index = 3;
  std::vector<ActionResponse> * ref_notadmitted_action_array = helper.get_not_admitted_list();
  if (ref_notadmitted_action_array->size() > 0){
    
    RICsubscriptionResponse_IEs_t *ies_notadmitted_actid = &IE_array[ie_index];
    ies_notadmitted_actid->criticality = Criticality_reject;
    ies_notadmitted_actid->id = ProtocolIE_ID_id_RICactions_NotAdmitted;

    RICaction_NotAdmitted_List_t *ric_not_admitted_actions_ie = &ies_notadmitted_actid->value.choice.RICaction_NotAdmitted_List;
    ric_not_admitted_actions_ie->list.count = 0;
  
  
    ies_notadmitted_actid->value.present =  RICsubscriptionResponse_IEs__value_PR_RICaction_NotAdmitted_List;
  
    // resize memory ?
    if (ref_notadmitted_action_array->size() >= ie_not_admitted_list_size){
      ie_not_admitted_list_size = 2 * ref_notadmitted_action_array->size();
      free(ie_not_admitted_list);
      ie_not_admitted_list = (RICaction_NotAdmitted_ItemIEs_t *)calloc(ie_not_admitted_list_size, sizeof(RICaction_NotAdmitted_ItemIEs_t));
      assert(ie_not_admitted_list != 0);
    
    };
  
  
    for(unsigned int i = 0; i < ref_notadmitted_action_array->size(); i ++){
      ie_not_admitted_list[i].criticality = Criticality_ignore;
      ie_not_admitted_list[i].id = ProtocolIE_ID_id_RICaction_NotAdmitted_Item ;
      ie_not_admitted_list[i].value.present = RICaction_NotAdmitted_ItemIEs__value_PR_RICaction_NotAdmitted_Item;;
      ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.ricActionID = (*ref_notadmitted_action_array)[i].get_id();

      int cause = (*ref_notadmitted_action_array)[i].get_cause();
      switch(cause){
      case Cause_PR_ricService:
	ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.ricService = (*ref_notadmitted_action_array)[i].get_sub_cause();
	break;
      case Cause_PR_transport:
	ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.transport = (*ref_notadmitted_action_array)[i].get_sub_cause();
	break;
      case Cause_PR_protocol:
	ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.protocol= (*ref_notadmitted_action_array)[i].get_sub_cause();
	break;
      case Cause_PR_misc:
	ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.misc = (*ref_notadmitted_action_array)[i].get_sub_cause();
	break;
      case Cause_PR_ricRequest:
	ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.ricRequest = (*ref_notadmitted_action_array)[i].get_sub_cause();
	break;
      default:
	mdclog_write(MDCLOG_ERR, "Error :: %s, %d : Unknown RIC cause %d\n", __FILE__, __LINE__, cause);
	return;
      }

      ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.present = (Cause_PR)cause;
      ASN_SEQUENCE_ADD(ric_not_admitted_actions_ie, &(ie_not_admitted_list[i]));
    }

    ASN_SEQUENCE_ADD(&subscription_response->protocolIEs, &(IE_array[ie_index]));
  }
  
}

void subscription_response:: get_fields(SuccessfulOutcome_t * success_msg,  subscription_response_helper & dout)
{

  assert(success_msg != NULL);
  
  RICrequestID_t *requestid;
  RANfunctionID_t * ranfunctionid;
  RICaction_Admitted_List_t  * ric_admitted_action_list;
  RICaction_NotAdmitted_List_t * ric_not_admitted_action_list;
    
  for(int edx = 0; edx < success_msg->value.choice.RICsubscriptionResponse.protocolIEs.list.count; edx++) {
    RICsubscriptionResponse_IEs_t *memb_ptr = success_msg->value.choice.RICsubscriptionResponse.protocolIEs.list.array[edx];
    
    switch(memb_ptr->id)
      {
      case (ProtocolIE_ID_id_RICrequestID):
	requestid = &memb_ptr->value.choice.RICrequestID;
	//dout.set_request(requestid->ricRequestorID, requestid->ricRequestSequenceNumber);
	break;
	  
      case (ProtocolIE_ID_id_RANfunctionID):
	ranfunctionid = &memb_ptr->value.choice.RANfunctionID;
	dout.set_function_id(*ranfunctionid);
	break;
	  
      case (ProtocolIE_ID_id_RICactions_Admitted):
	ric_admitted_action_list = &memb_ptr->value.choice.RICaction_Admitted_List;
	  
	// admitted actions
	for(int index = 0; index < ric_admitted_action_list->list.count; index ++){
	  RICaction_Admitted_ItemIEs_t * item = (RICaction_Admitted_ItemIEs_t *)ric_admitted_action_list->list.array[index];
	  long int id = item->value.choice.RICaction_Admitted_Item.ricActionID;
	  dout.get_admitted_list()->push_back(ActionResponse(id));
	};

	break;

      case (ProtocolIE_ID_id_RICactions_NotAdmitted):
	ric_not_admitted_action_list = &memb_ptr->value.choice.RICaction_NotAdmitted_List;
	  
	for(int index = 0; index < ric_not_admitted_action_list->list.count; index ++){
	  RICaction_NotAdmitted_ItemIEs_t * item = (RICaction_NotAdmitted_ItemIEs_t *)ric_not_admitted_action_list->list.array[index];
	  long int id = item->value.choice.RICaction_NotAdmitted_Item.ricActionID;
	  int cause = item->value.choice.RICaction_NotAdmitted_Item.cause.present;
	  int sub_cause;
	  switch(cause){
	      
	  case  Cause_PR_ricService :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.ricService;
	    break;
	      
	  case Cause_PR_transport :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.transport;
	    break;
	      
	  case  Cause_PR_protocol :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.protocol;
	    break;
	      
	  case Cause_PR_misc :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.misc;
	    break;
	      
	  case Cause_PR_ricRequest :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.ricRequest;
	    break;
	      
	  default:
	    std::cout <<"Error ! Illegal cause enum" << cause << std::endl;
	    return;
	  }  
	  dout.get_not_admitted_list()->push_back(ActionResponse(id, cause, sub_cause));
	}
	break;
      }
      
  }
    
  //asn_fprint(stdout, &asn_DEF_E2AP_PDU, e2pdu);
}


void subscription_response::set_fields_unsuccess( subscription_response_helper &helper){

  int ie_index;
  RICsubscriptionFailure_t * ric_subscription_failure = &(unsuccessMsg->value.choice.RICsubscriptionFailure);
  // reset list count
  ric_subscription_failure->protocolIEs.list.count = 0;
  
  ie_index = 0;
  RICsubscriptionFailure_IEs_t *ies_ricreq = &IE_Failure_array[ie_index];
  ies_ricreq->criticality = Criticality_reject;
  ies_ricreq->id = ProtocolIE_ID_id_RICrequestID;
  ies_ricreq->value.present = RICsubscriptionFailure_IEs__value_PR_RICrequestID;
  RICrequestID_t *ricrequest_ie = &ies_ricreq->value.choice.RICrequestID;
  ricrequest_ie->ricRequestorID = helper.get_request_id();
//  ricrequest_ie->ricRequestSequenceNumber = helper.get_req_seq();
  ASN_SEQUENCE_ADD(&ric_subscription_failure->protocolIEs, &(IE_Failure_array[ie_index]));  
  
  ie_index = 1;
  RICsubscriptionFailure_IEs_t *ies_ranfunc = &IE_Failure_array[ie_index];
  ies_ranfunc->criticality = Criticality_reject;
  ies_ranfunc->id = ProtocolIE_ID_id_RANfunctionID;
  ies_ranfunc->value.present = RICsubscriptionFailure_IEs__value_PR_RANfunctionID;
  RANfunctionID_t *ranfunction_ie = &ies_ranfunc->value.choice.RANfunctionID;
  *ranfunction_ie = helper.get_function_id();
  ASN_SEQUENCE_ADD(&ric_subscription_failure->protocolIEs, &(IE_Failure_array[ie_index]));  

  ie_index = 2;
  RICsubscriptionFailure_IEs_t *ies_notadmitted_actid = &IE_Failure_array[ie_index];
  ies_notadmitted_actid->criticality = Criticality_reject;
  ies_notadmitted_actid->id = ProtocolIE_ID_id_RICactions_NotAdmitted;
  RICaction_NotAdmitted_List_t *ric_not_admitted_actions_ie = &ies_notadmitted_actid->value.choice.RICaction_NotAdmitted_List;
  ric_not_admitted_actions_ie->list.count = 0;   
  std::vector<ActionResponse> * ref_notadmitted_action_array = helper.get_not_admitted_list();  
  if(ref_notadmitted_action_array->size() == 0){
    ies_notadmitted_actid->value.present =  RICsubscriptionFailure_IEs__value_PR_NOTHING;
  }
  else{
    ies_notadmitted_actid->value.present =  RICsubscriptionFailure_IEs__value_PR_RICaction_NotAdmitted_List;

    // resize memory  ?
    if (ref_notadmitted_action_array->size() >= ie_not_admitted_list_size){
      ie_not_admitted_list_size = 2 * ref_notadmitted_action_array->size();
      free(ie_not_admitted_list);
      ie_not_admitted_list = (RICaction_NotAdmitted_ItemIEs_t *)calloc(ie_not_admitted_list_size, sizeof(RICaction_NotAdmitted_ItemIEs_t));
      assert(ie_not_admitted_list != 0);
    };
      
  
    // reset the list count on ricAction_ToBeSetup_List;
    for(unsigned int i = 0; i < ref_notadmitted_action_array->size(); i ++){
      ie_not_admitted_list[i].criticality = Criticality_ignore;
      ie_not_admitted_list[i].id = ProtocolIE_ID_id_RICaction_NotAdmitted_Item ;
      ie_not_admitted_list[i].value.present = RICaction_NotAdmitted_ItemIEs__value_PR_RICaction_NotAdmitted_Item;;
      ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.ricActionID = (*ref_notadmitted_action_array)[i].get_id();

      int cause = (*ref_notadmitted_action_array)[i].get_cause();
      switch(cause){
      case Cause_PR_ricService:
    	  ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.ricService = (*ref_notadmitted_action_array)[i].get_sub_cause();
    	  break;
      case Cause_PR_transport:
    	  ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.transport = (*ref_notadmitted_action_array)[i].get_sub_cause();
    	  break;
      case Cause_PR_protocol:
    	  ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.protocol= (*ref_notadmitted_action_array)[i].get_sub_cause();
    	  break;
      case Cause_PR_misc:
    	  ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.misc = (*ref_notadmitted_action_array)[i].get_sub_cause();
    	  break;
      case Cause_PR_ricRequest:
    	  ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.choice.ricRequest = (*ref_notadmitted_action_array)[i].get_sub_cause();
    	  break;
      default:
	mdclog_write(MDCLOG_ERR, "Error :: %s, %d : Unknown RIC cause %d\n", __FILE__, __LINE__, cause);
	return ;
      }

      ie_not_admitted_list[i].value.choice.RICaction_NotAdmitted_Item.cause.present = (Cause_PR)cause;
      
      ASN_SEQUENCE_ADD(ric_not_admitted_actions_ie, &(ie_not_admitted_list[i]));
    }

  }
  ASN_SEQUENCE_ADD(&ric_subscription_failure->protocolIEs, &(IE_Failure_array[ie_index]));  


  // // criticality diagnostics is not generated/parsed currently since optional
  // i = 3;
  // RICsubscriptionFailure_IEs_t *ies_criticality_diagnostics= &IE_Failure_array[i];
  // ies_criticality_diagnostics->criticality = Criticality_ignore;
  // ies_criticality_diagnostics->id = ProtocolIE_ID_id_CriticalityDiagnostics ;
  // ies_criticality_diagnostics->value.present = RICsubscriptionFailure_IEs__value_PR_NOTHING;

    
}

void subscription_response:: get_fields(UnsuccessfulOutcome_t * unsuccess_msg,  subscription_response_helper & dout)
{

  assert(unsuccess_msg != NULL);
  
  RICrequestID_t *requestid;
  RANfunctionID_t * ranfunctionid;
  RICaction_NotAdmitted_List_t * ric_not_admitted_action_list;
    
  for(int edx = 0; edx < unsuccess_msg->value.choice.RICsubscriptionFailure.protocolIEs.list.count; edx++) {
    RICsubscriptionFailure_IEs_t *memb_ptr = unsuccess_msg->value.choice.RICsubscriptionFailure.protocolIEs.list.array[edx];
    
    switch(memb_ptr->id)
      {
      case (ProtocolIE_ID_id_RICrequestID):
	requestid = &memb_ptr->value.choice.RICrequestID;
	//dout.set_request(requestid->ricRequestorID, requestid->ricRequestSequenceNumber);
	break;
	  
      case (ProtocolIE_ID_id_RANfunctionID):
	ranfunctionid = &memb_ptr->value.choice.RANfunctionID;
	dout.set_function_id(*ranfunctionid);
	break;
	  

      case (ProtocolIE_ID_id_RICactions_NotAdmitted):
	ric_not_admitted_action_list = &memb_ptr->value.choice.RICaction_NotAdmitted_List;
	  
	for(int index = 0; index < ric_not_admitted_action_list->list.count; index ++){
	  RICaction_NotAdmitted_ItemIEs_t * item = (RICaction_NotAdmitted_ItemIEs_t *)ric_not_admitted_action_list->list.array[index];
	  long int id = item->value.choice.RICaction_NotAdmitted_Item.ricActionID;
	  int cause = item->value.choice.RICaction_NotAdmitted_Item.cause.present;
	  int sub_cause;
	  switch(cause){
	      
	  case  Cause_PR_ricService :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.ricService;
	    break;
	      
	  case Cause_PR_transport :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.transport;
	    break;
	      
	  case  Cause_PR_protocol :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.protocol;
	    break;
	      
	  case Cause_PR_misc :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.misc;
	    break;
	      
	  case Cause_PR_ricRequest :
	    sub_cause = item->value.choice.RICaction_NotAdmitted_Item.cause.choice.ricRequest;
	    break;
	      
	  default:
	    mdclog_write(MDCLOG_ERR, "Error :: %s, %d : Unknown RIC cause %d\n", __FILE__, __LINE__, cause);
	    return;
	  }  
	  dout.get_not_admitted_list()->push_back(ActionResponse(id, cause, sub_cause));
	}
	break;
      }
      
  }
    
  //asn_fprint(stdout, &asn_DEF_E2AP_PDU, e2pdu);
}



