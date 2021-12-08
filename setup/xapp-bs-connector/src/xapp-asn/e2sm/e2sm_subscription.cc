/*
  ==================================================================================

  Copyright (c) 2019-2020 AT&T Intellectual Property.

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
/*
 * e2sm_control.cc
 *
 *  Created on: Apr 30, 2020
 *      Author: Shraboni Jana
 */
/* Classes to handle E2 service model based on e2sm-HelloWorld-v001.asn */

#include "e2sm_subscription.hpp"

 //initialize
 e2sm_subscription::e2sm_subscription(void){

	memset(&event_fmt1, 0, sizeof(E2SM_HelloWorld_EventTriggerDefinition_Format1_t));

	memset(&actn_fmt1, 0, sizeof(E2SM_HelloWorld_ActionDefinition_Format1_t));


	ran_param = 0;
	ran_param = (RANparameter_Item_t*)calloc(1, sizeof(RANparameter_Item_t));
	assert(ran_param != 0);

    event_trigger = 0;
    event_trigger = ( E2SM_HelloWorld_EventTriggerDefinition_t *)calloc(1, sizeof( E2SM_HelloWorld_EventTriggerDefinition_t));
    assert(event_trigger != 0);

    action_defn = 0;
    action_defn = (E2SM_HelloWorld_ActionDefinition_t*)calloc(1, sizeof(E2SM_HelloWorld_ActionDefinition_t));
    assert(action_defn !=0);

    errbuf_len = 128;
  };

 e2sm_subscription::~e2sm_subscription(void){

  mdclog_write(MDCLOG_DEBUG, "Freeing event trigger object memory");

  event_trigger->choice.eventDefinition_Format1 = 0;

  action_defn->choice.actionDefinition_Format1 = 0;

  free(ran_param);

  ASN_STRUCT_FREE(asn_DEF_E2SM_HelloWorld_EventTriggerDefinition, event_trigger);
  ASN_STRUCT_FREE(asn_DEF_E2SM_HelloWorld_ActionDefinition, action_defn);


};

bool e2sm_subscription::encode_event_trigger(unsigned char *buf, size_t *size, e2sm_subscription_helper &helper){

  ASN_STRUCT_RESET(asn_DEF_E2SM_HelloWorld_EventTriggerDefinition, event_trigger);

  bool res;
  res = set_fields(event_trigger, helper);
  if (!res){

    return false;
  }

  int ret_constr = asn_check_constraints(&asn_DEF_E2SM_HelloWorld_EventTriggerDefinition, event_trigger, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    return false;
  }

  xer_fprint(stdout, &asn_DEF_E2SM_HelloWorld_EventTriggerDefinition, event_trigger);

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_HelloWorld_EventTriggerDefinition, event_trigger, buf, *size);

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

bool e2sm_subscription::encode_action_defn(unsigned char *buf, size_t *size, e2sm_subscription_helper &helper){

  bool res;
  res = set_fields(action_defn, helper);
  if (!res){
    return false;
  }


  int ret_constr = asn_check_constraints(&asn_DEF_E2SM_HelloWorld_ActionDefinition, action_defn, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    return false;
  }

  xer_fprint(stdout, &asn_DEF_E2SM_HelloWorld_ActionDefinition, action_defn);

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_HelloWorld_ActionDefinition, action_defn, buf, *size);

  if(retval.encoded == -1){
    error_string.assign(strerror(errno));
    return false;
  }
  else if (retval.encoded > *size){
    std::stringstream ss;
    ss  <<"Error encoding action definition. Reason =  encoded pdu size " << retval.encoded << " exceeds buffer size " << *size << std::endl;
    error_string = ss.str();
    return false;
  }
  else{
    *size = retval.encoded;
  }

  return true;
}

bool e2sm_subscription::set_fields(E2SM_HelloWorld_EventTriggerDefinition_t * ref_event_trigger, e2sm_subscription_helper & helper){

 if(ref_event_trigger == 0){
    error_string = "Invalid reference for Event Trigger Definition set fields";
    return false;
  }

  ref_event_trigger->present = E2SM_HelloWorld_EventTriggerDefinition_PR_eventDefinition_Format1;

  event_fmt1.triggerNature = helper.triger_nature;

  ref_event_trigger->choice.eventDefinition_Format1 = &event_fmt1;

  return true;
};

bool e2sm_subscription::set_fields(E2SM_HelloWorld_ActionDefinition_t * ref_action_defn, e2sm_subscription_helper & helper){

 if(ref_action_defn == 0){
    error_string = "Invalid reference for Event Action Definition set fields";
    return false;
  }
  ref_action_defn->present = E2SM_HelloWorld_ActionDefinition_PR_actionDefinition_Format1;


  ranparam_helper_t paramlst = helper.get_paramlist();

  for(RANParam_Helper item:paramlst){
	  ran_param->ranParameter_ID = item.getran_helper()._param_id;
	  ran_param->ranParameter_Name.buf = item.getran_helper()._param_name;
	  ran_param->ranParameter_Name.size = item.getran_helper()._param_name_len;
	  ran_param->ranParameter_Test = item.getran_helper()._param_test;
	  ran_param->ranParameter_Value.buf = item.getran_helper()._param_value;
	  ran_param->ranParameter_Value.size = item.getran_helper()._param_value_len;
	  ASN_SEQUENCE_ADD(&(actn_fmt1.ranParameter_List->list.array), ran_param);
  }


  ref_action_defn->choice.actionDefinition_Format1 = &actn_fmt1;


  return true;
};

