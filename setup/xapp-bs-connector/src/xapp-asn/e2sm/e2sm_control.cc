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
#include "e2sm_control.hpp"

 //initialize
 e2sm_control::e2sm_control(void){

	memset(&head_fmt1, 0, sizeof(E2SM_HelloWorld_ControlHeader_Format1_t));

	memset(&msg_fmt1, 0, sizeof(E2SM_HelloWorld_ControlMessage_Format1_t));



    control_head = 0;
    control_head = ( E2SM_HelloWorld_ControlHeader_t *)calloc(1, sizeof( E2SM_HelloWorld_ControlHeader_t));
    assert(control_head != 0);

    control_msg = 0;
    control_msg = (E2SM_HelloWorld_ControlMessage_t*)calloc(1, sizeof(E2SM_HelloWorld_ControlMessage_t));
    assert(control_msg !=0);

    errbuf_len = 128;
  };

 e2sm_control::~e2sm_control(void){

  mdclog_write(MDCLOG_DEBUG, "Freeing event trigger object memory");

  control_head->choice.controlHeader_Format1 = 0;

  control_msg->choice.controlMessage_Format1 = 0;

  ASN_STRUCT_FREE(asn_DEF_E2SM_HelloWorld_ControlHeader, control_head);
  ASN_STRUCT_FREE(asn_DEF_E2SM_HelloWorld_ControlMessage, control_msg);


};

bool e2sm_control::encode_control_header(unsigned char *buf, size_t *size, e2sm_control_helper &helper){

  ASN_STRUCT_RESET(asn_DEF_E2SM_HelloWorld_ControlHeader, control_head);

  bool res;
  res = set_fields(control_head, helper);
  if (!res){

    return false;
  }

  int ret_constr = asn_check_constraints(&asn_DEF_E2SM_HelloWorld_ControlHeader, control_head, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    return false;
  }

  xer_fprint(stdout, &asn_DEF_E2SM_HelloWorld_ControlHeader, control_head);

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_HelloWorld_ControlHeader, control_head, buf, *size);

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

bool e2sm_control::encode_control_message(unsigned char *buf, size_t *size, e2sm_control_helper &helper){

  bool res;
  res = set_fields(control_msg, helper);
  if (!res){
    return false;
  }


  int ret_constr = asn_check_constraints(&asn_DEF_E2SM_HelloWorld_ControlMessage, control_msg, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    return false;
  }

  xer_fprint(stdout, &asn_DEF_E2SM_HelloWorld_ControlMessage, control_msg);

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_HelloWorld_ControlMessage, control_msg, buf, *size);

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

bool e2sm_control::set_fields(E2SM_HelloWorld_ControlHeader_t * ref_control_head, e2sm_control_helper & helper){

 if(ref_control_head == 0){
    error_string = "Invalid reference for Event Trigger Definition set fields";
    return false;
  }

  ref_control_head->present = E2SM_HelloWorld_ControlHeader_PR_controlHeader_Format1;

  head_fmt1.controlHeaderParam = helper.header;

  ref_control_head->choice.controlHeader_Format1 = &head_fmt1;

  return true;
};

bool e2sm_control::set_fields(E2SM_HelloWorld_ControlMessage_t * ref_control_msg, e2sm_control_helper & helper){

 if(ref_control_msg == 0){
    error_string = "Invalid reference for Event Action Definition set fields";
    return false;
  }
  ref_control_msg->present = E2SM_HelloWorld_ControlMessage_PR_controlMessage_Format1;

  msg_fmt1.controlMsgParam.buf = helper.message;
  msg_fmt1.controlMsgParam.size = helper.message_len;


  ref_control_msg->choice.controlMessage_Format1 = &msg_fmt1;


  return true;
};

bool e2sm_control::get_fields(E2SM_HelloWorld_ControlHeader_t * ref_indictaion_header, e2sm_control_helper & helper){

	if (ref_indictaion_header == 0){
	    error_string = "Invalid reference for Control Header get fields";
	    return false;
	  }

	helper.header = ref_indictaion_header->choice.controlHeader_Format1->controlHeaderParam;
	return true;
}

bool e2sm_control::get_fields(E2SM_HelloWorld_ControlMessage_t * ref_control_message, e2sm_control_helper & helper){

	  if (ref_control_message == 0){
	  	    error_string = "Invalid reference for Control Message get fields";
	  	    return false;
	  	  }
	  helper.message = ref_control_message->choice.controlMessage_Format1->controlMsgParam.buf;
	  helper.message_len = ref_control_message->choice.controlMessage_Format1->controlMsgParam.size;

	  return true;
  }




