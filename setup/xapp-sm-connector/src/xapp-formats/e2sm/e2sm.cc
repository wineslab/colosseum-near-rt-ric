/*
  ==================================================================================

  Copyright (c) 2018-2019 AT&T Intellectual Property.
  
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
  
  http://www.apache.org/licenses/LICENSE-2.0
  
  Unless required by applicable law or agreed to in writing, softwares
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  ==================================================================================
*/

/* Classes to handle E2 service model based on e2sm-gNB-X2-release-1-v040.asn */

#include "e2sm.hpp"



  //initialize
  e2sm_event_trigger::e2sm_event_trigger(void){

    memset(&gNodeB_ID, 0, sizeof(E2N_GlobalGNB_ID_t));

    event_trigger = 0;
    event_trigger = ( E2N_E2SM_gNB_X2_eventTriggerDefinition_t *)calloc(1, sizeof( E2N_E2SM_gNB_X2_eventTriggerDefinition_t));
    assert(event_trigger != 0);
    
    // allocate space for gNodeB id  (used for encoding)
    gNodeB_ID.gNB_ID.choice.gNB_ID.buf = 0;
    gNodeB_ID.gNB_ID.choice.gNB_ID.buf = (uint8_t *)calloc(4, sizeof(uint8_t));
    assert(gNodeB_ID.gNB_ID.choice.gNB_ID.buf != 0);
    
    // allocate space for plmn identity  (used for encoding)
    gNodeB_ID.pLMN_Identity.buf = 0;
    gNodeB_ID.pLMN_Identity.buf = (uint8_t *) calloc(4, sizeof(uint8_t));
    assert(gNodeB_ID.pLMN_Identity.buf != 0);

    ie_list = 0;
    ie_list = ( struct E2N_InterfaceProtocolIE_Item *) calloc(INITIAL_LIST_SIZE, sizeof( struct E2N_InterfaceProtocolIE_Item));
    assert(ie_list != 0);
    ie_list_size = INITIAL_LIST_SIZE;

    condition_list = 0;
    condition_list = (E2N_E2SM_gNB_X2_eventTriggerDefinition::E2N_E2SM_gNB_X2_eventTriggerDefinition__interfaceProtocolIE_List *) calloc(1, sizeof(E2N_E2SM_gNB_X2_eventTriggerDefinition::E2N_E2SM_gNB_X2_eventTriggerDefinition__interfaceProtocolIE_List ));
    assert(condition_list != 0);

 
    
  };
  
e2sm_event_trigger::~e2sm_event_trigger(void){

  mdclog_write(MDCLOG_DEBUG, "Freeing event trigger object memory");
  for(int i = 0; i < condition_list->list.size; i++){
    condition_list->list.array[i] = 0;
  }

  if (condition_list->list.size > 0){
    free(condition_list->list.array);
    condition_list->list.array = 0;
    condition_list->list.size = 0;
    condition_list->list.count = 0;
  }

  free(condition_list);
  condition_list = 0;
  
  free(gNodeB_ID.gNB_ID.choice.gNB_ID.buf);
  gNodeB_ID.gNB_ID.choice.gNB_ID.buf = 0;
  
  free(gNodeB_ID.pLMN_Identity.buf);
  gNodeB_ID.pLMN_Identity.buf = 0;
  
  free(ie_list);
  ie_list = 0;
  
  event_trigger->interface_ID.choice.global_gNB_ID = 0;
  event_trigger->interfaceProtocolIE_List = 0;
  
  ASN_STRUCT_FREE(asn_DEF_E2N_E2SM_gNB_X2_eventTriggerDefinition, event_trigger);
  mdclog_write(MDCLOG_DEBUG, "Freed event trigger object memory");

 
};

bool e2sm_event_trigger::encode_event_trigger(unsigned char *buf, size_t *size, e2sm_event_trigger_helper &helper){
  
  bool res;
  res = set_fields(event_trigger, helper);
  if (!res){
    return false;
  }
  
  int ret_constr = asn_check_constraints(&asn_DEF_E2N_E2SM_gNB_X2_eventTriggerDefinition, event_trigger, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    return false;
  }

  //xer_fprint(stdout, &asn_DEF_E2N_E2SM_gNB_X2_eventTriggerDefinition, event_trigger);
  
  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2N_E2SM_gNB_X2_eventTriggerDefinition, event_trigger, buf, *size);
  
  if(retval.encoded == -1){
    error_string.assign(strerror(errno));
    return false;
  }
  else if (retval.encoded > *size){
    std::stringstream ss;
    ss  <<"Error encoding event trigger definition. Reason =  encoded pdu size " << retval.encoded << " exceeds buffer size " << *size << std::endl;
    error_string = ss.str();
    return false;
  }
  else{
    *size = retval.encoded;
  }
  
  return true;
}


bool e2sm_event_trigger::set_fields(E2N_E2SM_gNB_X2_eventTriggerDefinition_t * ref_event_trigger, e2sm_event_trigger_helper & helper){
  if(ref_event_trigger == 0){
    error_string = "Invalid reference for Event Trigger Definition set fields";
    return false;
  }
      
  // set the message type
  ref_event_trigger->interfaceMessageType.procedureCode = helper.procedure_code;
  ref_event_trigger->interfaceMessageType.typeOfMessage = helper.message_type;
  
  ref_event_trigger->interfaceDirection = helper.interface_direction; 
  ref_event_trigger->interface_ID.present = E2N_Interface_ID_PR_global_gNB_ID;
  
  ref_event_trigger->interface_ID.choice.global_gNB_ID = &gNodeB_ID;

  // to do : need to put correct code here for upding plmn id and gNodeB
  // for now just place holders :
  //================================================================
  memcpy(gNodeB_ID.pLMN_Identity.buf, helper.plmn_id.c_str(), 3);
  gNodeB_ID.pLMN_Identity.size = 3;
  
  memcpy(gNodeB_ID.gNB_ID.choice.gNB_ID.buf, helper.egNB_id.c_str(), 3);
  gNodeB_ID.gNB_ID.choice.gNB_ID.size = 3;
  
  // we only do global gNodeB id for now, not eNodeB
  gNodeB_ID.gNB_ID.present = E2N_GNB_ID_PR_gNB_ID;
  //================================================================
  
  
  // Add in any requested IE items
  std::vector<Item> * ref_ie_array = helper.get_list();

  if (ref_ie_array->size() == 0){
    ref_event_trigger->interfaceProtocolIE_List = 0;
    
  }
  else{
    ref_event_trigger->interfaceProtocolIE_List = condition_list;
    
    //resize memory ? 
    if(ref_ie_array->size() > ie_list_size){
      ie_list_size = 2 * ref_ie_array->size();
      free(ie_list);
      ie_list = (struct E2N_InterfaceProtocolIE_Item *)calloc(ie_list_size, sizeof(struct E2N_InterfaceProtocolIE_Item));
      assert(ie_list != 0);
    }
    
    // reset the count so that adds start from the beginning
    ref_event_trigger->interfaceProtocolIE_List->list.count = 0;
    
    for(unsigned int i = 0; i < ref_ie_array->size(); i++){

      ie_list[i].interfaceProtocolIE_ID = (*ref_ie_array)[i].interface_id;
      ie_list[i].interfaceProtocolIE_Test = (*ref_ie_array)[i].test;
      
      //switch(ie_list[i].interfaceProtocolIE_Value.present){
      switch((*ref_ie_array)[i].val_type){
	
      case (E2N_InterfaceProtocolIE_Value_PR_valueInt):
	ie_list[i].interfaceProtocolIE_Value.present = E2N_InterfaceProtocolIE_Value_PR_valueInt;
	ie_list[i].interfaceProtocolIE_Value.choice.valueInt = (*ref_ie_array)[i].value_n;
	break;
	
      case (E2N_InterfaceProtocolIE_Value_PR_valueEnum):
	ie_list[i].interfaceProtocolIE_Value.present = E2N_InterfaceProtocolIE_Value_PR_valueEnum;
	ie_list[i].interfaceProtocolIE_Value.choice.valueEnum = (*ref_ie_array)[i].value_n;
	break;
	
      case (E2N_InterfaceProtocolIE_Value_PR_valueBool):
	ie_list[i].interfaceProtocolIE_Value.present = E2N_InterfaceProtocolIE_Value_PR_valueBool;
	ie_list[i].interfaceProtocolIE_Value.choice.valueBool = (*ref_ie_array)[i].value_n;
	break;
	
      case (E2N_InterfaceProtocolIE_Value_PR_valueBitS):
	ie_list[i].interfaceProtocolIE_Value.present = E2N_InterfaceProtocolIE_Value_PR_valueBitS;
	ie_list[i].interfaceProtocolIE_Value.choice.valueBitS.buf = (uint8_t *)(*ref_ie_array)[i].value_s.c_str();
	ie_list[i].interfaceProtocolIE_Value.choice.valueBitS.size = (*ref_ie_array)[i].value_s.length();
	break;

      case (E2N_InterfaceProtocolIE_Value_PR_valueOctS):
	ie_list[i].interfaceProtocolIE_Value.present = E2N_InterfaceProtocolIE_Value_PR_valueOctS;
	ie_list[i].interfaceProtocolIE_Value.choice.valueOctS.buf = (uint8_t *)(*ref_ie_array)[i].value_s.c_str();
	ie_list[i].interfaceProtocolIE_Value.choice.valueOctS.size = (*ref_ie_array)[i].value_s.length();
	break;

      default:
	{
	  std::stringstream ss;
	  ss <<"Error ! " << __FILE__ << "," << __LINE__ << " illegal enum " << (*ref_ie_array)[i].val_type << " for interface Protocol IE value" << std::endl;
	  std::string error_string = ss.str();
	  return false;
	}
      }
      
      ASN_SEQUENCE_ADD(ref_event_trigger->interfaceProtocolIE_List, &ie_list[i]);
    }
  }

  return true;
};
  

bool e2sm_event_trigger::get_fields(E2N_E2SM_gNB_X2_eventTriggerDefinition_t * ref_event_trigger, e2sm_event_trigger_helper & helper){

  if (ref_event_trigger == 0){
    error_string = "Invalid reference for Event Trigger definition get fields";
    return false;
  }
  
  helper.procedure_code = ref_event_trigger->interfaceMessageType.procedureCode;
  helper.message_type   = ref_event_trigger->interfaceMessageType.typeOfMessage;
  helper.interface_direction = ref_event_trigger->interfaceDirection;
  
  helper.plmn_id.assign((const char *)ref_event_trigger->interface_ID.choice.global_gNB_ID->pLMN_Identity.buf, ref_event_trigger->interface_ID.choice.global_gNB_ID->pLMN_Identity.size);
  helper.egNB_id.assign((const char *)ref_event_trigger->interface_ID.choice.global_gNB_ID->gNB_ID.choice.gNB_ID.buf, ref_event_trigger->interface_ID.choice.global_gNB_ID->gNB_ID.choice.gNB_ID.size);
  for(int i = 0; i < ref_event_trigger->interfaceProtocolIE_List->list.count; i++){
    struct E2N_InterfaceProtocolIE_Item * ie_item = ref_event_trigger->interfaceProtocolIE_List->list.array[i];
    switch(ie_item->interfaceProtocolIE_Value.present){
    case (E2N_InterfaceProtocolIE_Value_PR_valueInt):
      helper.add_protocol_ie_item(ie_item->interfaceProtocolIE_ID, ie_item->interfaceProtocolIE_Test, ie_item->interfaceProtocolIE_Value.present, ie_item->interfaceProtocolIE_Value.choice.valueInt);
      break;
    case (E2N_InterfaceProtocolIE_Value_PR_valueEnum):
      helper.add_protocol_ie_item(ie_item->interfaceProtocolIE_ID, ie_item->interfaceProtocolIE_Test, ie_item->interfaceProtocolIE_Value.present, ie_item->interfaceProtocolIE_Value.choice.valueEnum);
      break;
    case (E2N_InterfaceProtocolIE_Value_PR_valueBool):
      helper.add_protocol_ie_item(ie_item->interfaceProtocolIE_ID, ie_item->interfaceProtocolIE_Test, ie_item->interfaceProtocolIE_Value.present, ie_item->interfaceProtocolIE_Value.choice.valueBool);	    
      break;
    case (E2N_InterfaceProtocolIE_Value_PR_valueBitS):
      helper.add_protocol_ie_item(ie_item->interfaceProtocolIE_ID, ie_item->interfaceProtocolIE_Test, ie_item->interfaceProtocolIE_Value.present, std::string((const char *)ie_item->interfaceProtocolIE_Value.choice.valueBitS.buf,ie_item->interfaceProtocolIE_Value.choice.valueBitS.size) );
      break;
    case (E2N_InterfaceProtocolIE_Value_PR_valueOctS):
      helper.add_protocol_ie_item(ie_item->interfaceProtocolIE_ID, ie_item->interfaceProtocolIE_Test, ie_item->interfaceProtocolIE_Value.present, std::string((const char *)ie_item->interfaceProtocolIE_Value.choice.valueOctS.buf,ie_item->interfaceProtocolIE_Value.choice.valueOctS.size) );
      break;
    default:
      mdclog_write(MDCLOG_ERR, "Error : %s, %d: Unkown interface protocol IE type %d in event trigger definition\n", __FILE__, __LINE__, ie_item->interfaceProtocolIE_Value.present);
      return false;
    }
  }
  
  return true;
};
    

  
   
// initialize
e2sm_indication::e2sm_indication(void) {
  
  memset(&gNodeB_ID, 0, sizeof(E2N_GlobalGNB_ID_t));
    
  // allocate space for gNodeB id  (used for encoding)
  gNodeB_ID.gNB_ID.choice.gNB_ID.buf = (uint8_t *)calloc(4, sizeof(uint8_t));
  assert(gNodeB_ID.gNB_ID.choice.gNB_ID.buf != 0);
    
  // allocate space for plmn identity  (used for encoding)
  gNodeB_ID.pLMN_Identity.buf = (uint8_t *) calloc(4, sizeof(uint8_t));
  assert(gNodeB_ID.pLMN_Identity.buf != 0);

  header = 0;
  header = (E2N_E2SM_gNB_X2_indicationHeader_t *)calloc(1, sizeof(E2N_E2SM_gNB_X2_indicationHeader_t));
  assert(header != 0);

  message = 0;
  message = (E2N_E2SM_gNB_X2_indicationMessage_t *)calloc(1, sizeof(E2N_E2SM_gNB_X2_indicationMessage_t));
  assert(message != 0);
}
  
e2sm_indication::~e2sm_indication(void){
  mdclog_write(MDCLOG_DEBUG, "Freeing E2N_E2SM Indication  object memory");

  free(gNodeB_ID.gNB_ID.choice.gNB_ID.buf);
  free(gNodeB_ID.pLMN_Identity.buf);
  
  header->interface_ID.choice.global_gNB_ID = 0;

  ASN_STRUCT_FREE(asn_DEF_E2N_E2SM_gNB_X2_indicationHeader, header);

  message->interfaceMessage.buf = 0;
  message->interfaceMessage.size = 0;

  ASN_STRUCT_FREE(asn_DEF_E2N_E2SM_gNB_X2_indicationMessage, message);
  mdclog_write(MDCLOG_DEBUG, "Freed E2SM Indication  object memory");
    
}
  
  

bool e2sm_indication::encode_indication_header(unsigned char *buf, size_t *size, e2sm_header_helper &helper){
    
  bool res;
  res = set_header_fields(header, helper);
  if (!res){
    return false;
  }

  int ret_constr = asn_check_constraints(&asn_DEF_E2N_E2SM_gNB_X2_indicationHeader, header, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    error_string = "E2SM Indication Header Constraint failed : " + error_string;

    return false;
  }

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2N_E2SM_gNB_X2_indicationHeader, header, buf, *size);

  if(retval.encoded == -1){
    error_string.assign(strerror(errno));
    error_string = "Error encoding E2N_E2SM Indication Header. Reason = " + error_string;
    return false;
  }
  else if (retval.encoded > *size){
    std::stringstream ss;
    ss  <<"Error encoding E2SM Indication Header . Reason =  encoded pdu size " << retval.encoded << " exceeds buffer size " << *size << std::endl;
    error_string = ss.str();
    return false;
  }
  else{
    *size = retval.encoded;
  }
    
  return true;
}


bool e2sm_indication::encode_indication_message(unsigned char *buf, size_t *size, e2sm_message_helper &helper){

  set_message_fields(message, helper); 

  int ret_constr = asn_check_constraints(&asn_DEF_E2N_E2SM_gNB_X2_indicationMessage, message, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    error_string = "E2SM Indication Message Constraint failed : " + error_string;
    return false;
  }

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2N_E2SM_gNB_X2_indicationMessage, message, buf, *size);
  if(retval.encoded == -1){
    error_string.assign(strerror(errno));
    error_string = "Error encoding E2SM Indication Header. Reason = " + error_string;
    return false;
  }
  else if (retval.encoded > *size){
    std::stringstream ss;
    ss  <<"Error encoding E2N_E2SM Indication Message . Reason =  encoded pdu size " << retval.encoded << " exceeds buffer size " << *size << std::endl;
    error_string = ss.str();
    
    return false;
  }
  else{
    *size = retval.encoded;
  }
  
  return true;
}



// Used when generating an indication header 
bool e2sm_indication::set_header_fields(E2N_E2SM_gNB_X2_indicationHeader_t *header,  e2sm_header_helper &helper){

  if (header == 0){
    error_string = "Invalid reference for E2SM Indication Header set fields";
    return false;
  }
  
  
  header->interfaceDirection = helper.interface_direction;
  header->interface_ID.present = E2N_Interface_ID_PR_global_gNB_ID;
  header->interface_ID.choice.global_gNB_ID = &gNodeB_ID;


  // to do : need to put correct code here for upding plmn id and gNodeB
  // for now just place holders :
  memcpy(gNodeB_ID.pLMN_Identity.buf, helper.plmn_id.c_str(), 3);
  gNodeB_ID.pLMN_Identity.size = 3;
  
  memcpy(gNodeB_ID.gNB_ID.choice.gNB_ID.buf, helper.egNB_id.c_str(), 3);
  gNodeB_ID.gNB_ID.choice.gNB_ID.size = 3;
  
  // we only do global gNodeB id for now, not eNodeB
  gNodeB_ID.gNB_ID.present = E2N_GNB_ID_PR_gNB_ID;

  return true;
  
};


// used when decoding an indication header
bool e2sm_indication::get_header_fields(E2N_E2SM_gNB_X2_indicationHeader_t *header,  e2sm_header_helper &helper){

  if (header == 0){
    error_string = "Invalid reference for E2SM Indication header get fields";
    return false;
  }
  
  helper.interface_direction = header->interfaceDirection;
  helper.plmn_id.assign((const char *)header->interface_ID.choice.global_gNB_ID->pLMN_Identity.buf, header->interface_ID.choice.global_gNB_ID->pLMN_Identity.size);
  helper.egNB_id.assign((const char *)header->interface_ID.choice.global_gNB_ID->gNB_ID.choice.gNB_ID.buf, header->interface_ID.choice.global_gNB_ID->gNB_ID.choice.gNB_ID.size);
  
  // to do : add code to decipher plmn and global gnodeb from ints (since that is likely the convention for packing)

  return true;
}



// Used when generating an indication message 
bool   e2sm_indication::set_message_fields(E2N_E2SM_gNB_X2_indicationMessage_t *interface_message,  e2sm_message_helper &helper){

  if(interface_message == 0){
    error_string = "Invalid reference for E2SM Indication Message set fields";
    return false;
  }

  // interface-message is an octet string. just point it to the buffer
  interface_message->interfaceMessage.buf = &(helper.x2ap_pdu[0]);
  interface_message->interfaceMessage.size = helper.x2ap_pdu_size;

  return true;
  
};

// used when decoding an indication message
bool e2sm_indication::get_message_fields( E2N_E2SM_gNB_X2_indicationMessage_t *interface_message, e2sm_message_helper &helper){

  
  if(interface_message == 0){
    error_string = "Invalid reference for E2SM Indication Message get fields";
    return false;
  }

  // interface message is an octet string
  helper.x2ap_pdu = interface_message->interfaceMessage.buf;;
  helper.x2ap_pdu_size = interface_message->interfaceMessage.size;

  return true;
  
}
  

   
// initialize
e2sm_control::e2sm_control(void) {
  
  memset(&gNodeB_ID, 0, sizeof(E2N_GlobalGNB_ID_t));
    
  // allocate space for gNodeB id  (used for encoding)
  gNodeB_ID.gNB_ID.choice.gNB_ID.buf = (uint8_t *)calloc(4, sizeof(uint8_t));
  assert(gNodeB_ID.gNB_ID.choice.gNB_ID.buf != 0);
    
  // allocate space for plmn identity  (used for encoding)
  gNodeB_ID.pLMN_Identity.buf = (uint8_t *) calloc(4, sizeof(uint8_t));
  assert(gNodeB_ID.pLMN_Identity.buf != 0);

  header = 0;
  header = (E2N_E2SM_gNB_X2_controlHeader_t *)calloc(1, sizeof(E2N_E2SM_gNB_X2_controlHeader_t));
  assert(header != 0);

  message = 0;
  message = (E2N_E2SM_gNB_X2_controlMessage_t *)calloc(1, sizeof(E2N_E2SM_gNB_X2_controlMessage_t));
  assert(message != 0);
}
  
e2sm_control::~e2sm_control(void){
  mdclog_write(MDCLOG_DEBUG, "Freeing E2SM Control  object memory");

  free(gNodeB_ID.gNB_ID.choice.gNB_ID.buf);
  free(gNodeB_ID.pLMN_Identity.buf);
  header->interface_ID.choice.global_gNB_ID = 0;
  ASN_STRUCT_FREE(asn_DEF_E2N_E2SM_gNB_X2_controlHeader, header);

  message->interfaceMessage.buf = 0;
  ASN_STRUCT_FREE(asn_DEF_E2N_E2SM_gNB_X2_controlMessage, message);

  mdclog_write(MDCLOG_DEBUG, "Freed E2SM Control  object memory");
    
}
  
  

bool e2sm_control::encode_control_header(unsigned char *buf, size_t *size, e2sm_header_helper &helper){
    
  bool res;
  res = set_header_fields(header, helper);
  if (!res){
    return false;
  }

  int ret_constr = asn_check_constraints(&asn_DEF_E2N_E2SM_gNB_X2_controlHeader, header, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    error_string = "E2SM Control Header Constraint failed : " + error_string;

    return false;
  }

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2N_E2SM_gNB_X2_controlHeader, header, buf, *size);

  if(retval.encoded == -1){
    error_string.assign(strerror(errno));
    error_string = "Error encoding E2SM Control Header. Reason = " + error_string;
    return false;
  }
  else if (retval.encoded > *size){
    std::stringstream ss;
    ss  <<"Error encoding E2N_E2SM Control Header . Reason =  encoded pdu size " << retval.encoded << " exceeds buffer size " << *size << std::endl;
    error_string = ss.str();
    return false;
  }
  else{
    *size = retval.encoded;
  }
    
  return true;
}


bool e2sm_control::encode_control_message(unsigned char *buf, size_t *size, e2sm_message_helper &helper){

  set_message_fields(message, helper); 

  int ret_constr = asn_check_constraints(&asn_DEF_E2N_E2SM_gNB_X2_controlMessage, message, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    error_string = "E2SM Control Message Constraint failed : " + error_string;
    return false;
  }

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2N_E2SM_gNB_X2_controlMessage, message, buf, *size);
  if(retval.encoded == -1){
    error_string.assign(strerror(errno));
    error_string = "Error encoding E2SM Control Message. Reason = " + error_string;
    return false;
  }
  else if (retval.encoded > *size){
    std::stringstream ss;
    ss  <<"Error encoding E2SM Control Message . Reason =  encoded pdu size " << retval.encoded << " exceeds buffer size " << *size << std::endl;
    error_string = ss.str();
    
    return false;
  }
  else{
    *size = retval.encoded;
  }
  
  return true;
}



// Used when generating an indication header 
bool e2sm_control::set_header_fields(E2N_E2SM_gNB_X2_controlHeader_t *header,  e2sm_header_helper &helper){

  if (header == 0){
    error_string = "Invalid reference for E2SM Control Header set fields";
    return false;
  }
  
  
  header->interfaceDirection = helper.interface_direction;
  header->interface_ID.present = E2N_Interface_ID_PR_global_gNB_ID;
  header->interface_ID.choice.global_gNB_ID = &gNodeB_ID;


  // to do : need to put correct code here for upding plmn id and gNodeB
  // for now just place holders :
  memcpy(gNodeB_ID.pLMN_Identity.buf, helper.plmn_id.c_str(), 3);
  gNodeB_ID.pLMN_Identity.size = 3;
  
  memcpy(gNodeB_ID.gNB_ID.choice.gNB_ID.buf, helper.egNB_id.c_str(), 3);
  gNodeB_ID.gNB_ID.choice.gNB_ID.size = 3;
  
  // we only do global gNodeB id for now, not eNodeB
  gNodeB_ID.gNB_ID.present = E2N_GNB_ID_PR_gNB_ID;

  return true;
  
};


// used when decoding an indication header
bool e2sm_control::get_header_fields(E2N_E2SM_gNB_X2_controlHeader_t *header,  e2sm_header_helper &helper){

  if (header == 0){
    error_string = "Invalid reference for E2SM Control header get fields";
    return false;
  }
  
  helper.interface_direction = header->interfaceDirection;
  helper.plmn_id.assign((const char *)header->interface_ID.choice.global_gNB_ID->pLMN_Identity.buf, header->interface_ID.choice.global_gNB_ID->pLMN_Identity.size);
  helper.egNB_id.assign((const char *)header->interface_ID.choice.global_gNB_ID->gNB_ID.choice.gNB_ID.buf, header->interface_ID.choice.global_gNB_ID->gNB_ID.choice.gNB_ID.size);
  
  // to do : add code to decipher plmn and global gnodeb from ints (since that is likely the convention for packing)

  return true;
}



// Used when generating an indication message 
bool   e2sm_control::set_message_fields(E2N_E2SM_gNB_X2_controlMessage_t *interface_message,  e2sm_message_helper &helper){

  if(interface_message == 0){
    error_string = "Invalid reference for E2SM Control Message set fields";
    return false;
  }

  // interface-message is an octet string. just point it to the buffer
  interface_message->interfaceMessage.buf = &(helper.x2ap_pdu[0]);
  interface_message->interfaceMessage.size = helper.x2ap_pdu_size;

  return true;
  
};

// used when decoding an indication message
bool e2sm_control::get_message_fields( E2N_E2SM_gNB_X2_controlMessage_t *interface_message, e2sm_message_helper &helper){

  
  if(interface_message == 0){
    error_string = "Invalid reference for E2SM Control Message get fields";
    return false;
  }

  // interface message is an octet string
  helper.x2ap_pdu = interface_message->interfaceMessage.buf;;
  helper.x2ap_pdu_size = interface_message->interfaceMessage.size;

  return true;
  
}
  
