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
 * ric_control_request.c
 *
 *  Created on: Jul 11, 2019
 *      Author: sjana, Ashwin Sridharan
 */

#include "e2ap_control.hpp"

// Set up memory allocations for each IE for encoding
// We are responsible for memory management for each IE for encoding
// Hence destructor should clear out memory
// When decoding, we rely on asn1c macro (ASN_STRUCT_FREE to be called
// for releasing memory by external calling function)
ric_control_request::ric_control_request(void){

  e2ap_pdu_obj = 0;
  e2ap_pdu_obj = (E2AP_PDU_t * )calloc(1, sizeof(E2AP_PDU_t));
  assert(e2ap_pdu_obj != 0);

  initMsg = 0;
  initMsg = (InitiatingMessage_t * )calloc(1, sizeof(InitiatingMessage_t));
  assert(initMsg != 0);

  IE_array = 0;
  IE_array = (RICcontrolRequest_IEs_t *)calloc(NUM_CONTROL_REQUEST_IES, sizeof(RICcontrolRequest_IEs_t));
  assert(IE_array != 0);

  e2ap_pdu_obj->present = E2AP_PDU_PR_initiatingMessage;
  e2ap_pdu_obj->choice.initiatingMessage = initMsg;

  
};


// Clear assigned protocolIE list from RIC control_request IE container
ric_control_request::~ric_control_request(void){

  mdclog_write(MDCLOG_DEBUG, "Freeing E2AP Control Request object memory");
  
  RICcontrolRequest_t *ricControl_Request  = &(initMsg->value.choice.RICcontrolRequest);
  for(int i = 0; i < ricControl_Request->protocolIEs.list.size; i++){
    ricControl_Request->protocolIEs.list.array[i] = 0;
  }
  
  if (ricControl_Request->protocolIEs.list.size > 0){
    free(ricControl_Request->protocolIEs.list.array);
    ricControl_Request->protocolIEs.list.size = 0;
    ricControl_Request->protocolIEs.list.count = 0;
  }
  
  free(IE_array);
  free(initMsg);
  e2ap_pdu_obj->choice.initiatingMessage = 0;
  
  ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, e2ap_pdu_obj);
  mdclog_write(MDCLOG_DEBUG, "Freed E2AP Control Request object mempory");
  
}


bool ric_control_request::encode_e2ap_control_request(unsigned char *buf, size_t *size, ric_control_helper & dinput){

  initMsg->procedureCode = ProcedureCode_id_RICcontrol;
  initMsg->criticality = Criticality_ignore;
  initMsg->value.present = InitiatingMessage__value_PR_RICcontrolRequest;

  bool res;
  
  res = set_fields(initMsg, dinput);
  if (!res){
    return false;
  }

  int ret_constr = asn_check_constraints(&asn_DEF_E2AP_PDU, (void *) e2ap_pdu_obj, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(errbuf, errbuf_len);
    error_string = "Constraints failed for encoding control . Reason = " + error_string;
    return false;
  }

  //xer_fprint(stdout, &asn_DEF_E2AP_PDU, e2ap_pdu_obj);
  
  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2AP_PDU, e2ap_pdu_obj, buf, *size);
  
  if(retval.encoded == -1){
    error_string.assign(strerror(errno));
    return false;
  }
  else {
    if(*size < retval.encoded){
      std::stringstream ss;
      ss  <<"Error encoding event trigger definition. Reason =  encoded pdu size " << retval.encoded << " exceeds buffer size " << *size << std::endl;
      error_string = ss.str();
      return false;
    }
  }

  *size = retval.encoded;
  return true;
  
}

bool ric_control_request::set_fields(InitiatingMessage_t *initMsg, ric_control_helper &dinput){
  unsigned int ie_index;

  if (initMsg == 0){
    error_string = "Invalid reference for E2AP Control_Request message in set_fields";
    return false;
  }

  RICcontrolRequest_t * ric_control_request = &(initMsg->value.choice.RICcontrolRequest);
  ric_control_request->protocolIEs.list.count = 0; // reset 
  
  // for(i = 0; i < NUM_CONTROL_REQUEST_IES;i++){
  //   memset(&(IE_array[i]), 0, sizeof(RICcontrolRequest_IEs_t));
  // }
 
  // Mandatory IE
  ie_index = 0;
  RICcontrolRequest_IEs_t *ies_ricreq = &IE_array[ie_index];
  ies_ricreq->criticality = Criticality_reject;
  ies_ricreq->id = ProtocolIE_ID_id_RICrequestID;
  ies_ricreq->value.present = RICcontrolRequest_IEs__value_PR_RICrequestID;
  RICrequestID_t *ricrequest_ie = &ies_ricreq->value.choice.RICrequestID;
  ricrequest_ie->ricRequestorID = dinput.req_id;
  //ricrequest_ie->ricRequestSequenceNumber = dinput.req_seq_no;
  ASN_SEQUENCE_ADD(&(ric_control_request->protocolIEs), &(IE_array[ie_index]));

  // Mandatory IE
  ie_index = 1;
  RICcontrolRequest_IEs_t *ies_ranfunc = &IE_array[ie_index];
  ies_ranfunc->criticality = Criticality_reject;
  ies_ranfunc->id = ProtocolIE_ID_id_RANfunctionID;
  ies_ranfunc->value.present = RICcontrolRequest_IEs__value_PR_RANfunctionID;
  RANfunctionID_t *ranfunction_ie = &ies_ranfunc->value.choice.RANfunctionID;
  *ranfunction_ie = dinput.func_id;
  ASN_SEQUENCE_ADD(&(ric_control_request->protocolIEs), &(IE_array[ie_index]));


  // Mandatory IE
  ie_index = 2;
  RICcontrolRequest_IEs_t *ies_richead = &IE_array[ie_index];
  ies_richead->criticality = Criticality_reject;
  ies_richead->id = ProtocolIE_ID_id_RICcontrolHeader;
  ies_richead->value.present = RICcontrolRequest_IEs__value_PR_RICcontrolHeader;
  RICcontrolHeader_t *richeader_ie = &ies_richead->value.choice.RICcontrolHeader;
  richeader_ie->buf = dinput.control_header;
  richeader_ie->size = dinput.control_header_size;
  ASN_SEQUENCE_ADD(&(ric_control_request->protocolIEs), &(IE_array[ie_index]));

  // Mandatory IE
  ie_index = 3;
  RICcontrolRequest_IEs_t *ies_indmsg = &IE_array[ie_index];
  ies_indmsg->criticality = Criticality_reject;
  ies_indmsg->id = ProtocolIE_ID_id_RICcontrolMessage;
  ies_indmsg->value.present = RICcontrolRequest_IEs__value_PR_RICcontrolMessage;
  RICcontrolMessage_t *ricmsg_ie = &ies_indmsg->value.choice.RICcontrolMessage;
  ricmsg_ie->buf = dinput.control_msg;
  ricmsg_ie->size = dinput.control_msg_size;
  ASN_SEQUENCE_ADD(&(ric_control_request->protocolIEs), &(IE_array[ie_index]));

  // Optional IE
  ie_index = 4;
  if (dinput.control_ack >= 0){
    RICcontrolRequest_IEs_t *ies_indtyp = &IE_array[ie_index];
    ies_indtyp->criticality = Criticality_reject;
    ies_indtyp->id = ProtocolIE_ID_id_RICcontrolAckRequest;
    ies_indtyp->value.present = RICcontrolRequest_IEs__value_PR_RICcontrolAckRequest;
    RICcontrolAckRequest_t *ricackreq_ie = &ies_indtyp->value.choice.RICcontrolAckRequest;
    *ricackreq_ie = dinput.control_ack;
    ASN_SEQUENCE_ADD(&(ric_control_request->protocolIEs), &(IE_array[ie_index]));
  }

  // Optional IE
  ie_index = 5;
  if(dinput.call_process_id_size > 0){
    RICcontrolRequest_IEs_t *ies_callprocid = &IE_array[ie_index];
    ies_callprocid->criticality = Criticality_reject;
    ies_callprocid->id = ProtocolIE_ID_id_RICcallProcessID;
    ies_callprocid->value.present = RICcontrolRequest_IEs__value_PR_RICcallProcessID;
    RICcallProcessID_t *riccallprocessid_ie = &ies_callprocid->value.choice.RICcallProcessID;
    riccallprocessid_ie->buf = dinput.call_process_id;
    riccallprocessid_ie->size = dinput.call_process_id_size;
    ASN_SEQUENCE_ADD(&(ric_control_request->protocolIEs), &(IE_array[ie_index]));

  }
  return true;

};

  


bool ric_control_request:: get_fields(InitiatingMessage_t * init_msg,  ric_control_helper &dout)
{
  if (init_msg == 0){
    error_string = "Invalid reference for E2AP Control_Request message in get_fields";
    return false;
  }
  
 
  for(int edx = 0; edx < init_msg->value.choice.RICcontrolRequest.protocolIEs.list.count; edx++) {
    RICcontrolRequest_IEs_t *memb_ptr = init_msg->value.choice.RICcontrolRequest.protocolIEs.list.array[edx];
    
    switch(memb_ptr->id)
      {
      case (ProtocolIE_ID_id_RICcontrolHeader):
  	dout.control_header = memb_ptr->value.choice.RICcontrolHeader.buf;
  	dout.control_header_size = memb_ptr->value.choice.RICcontrolHeader.size;
  	break;
	
      case (ProtocolIE_ID_id_RICcontrolMessage):
  	dout.control_msg =  memb_ptr->value.choice.RICcontrolMessage.buf;
  	dout.control_msg_size = memb_ptr->value.choice.RICcontrolMessage.size;
  	break;

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
	
      case (ProtocolIE_ID_id_RICcontrolAckRequest):
  	dout.control_ack = memb_ptr->value.choice.RICcontrolAckRequest;
  	break;
	
      default:
  	break;
      }
    
  }
  
  return true;

}

InitiatingMessage_t * ric_control_request::get_message(void)  {
    return initMsg;
}
