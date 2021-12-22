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
 * e2sm_indication.cc
 *
 *  Created on: Apr, 2020
 *      Author: Shraboni Jana
 */
/* Classes to handle E2 service model based on e2sm-HelloWorld-v001.asn */
#include "e2sm_indication.hpp"

 //initialize
 e2sm_indication::e2sm_indication(void){

	memset(&head_fmt1, 0, sizeof(E2SM_HelloWorld_IndicationHeader_Format1_t));

	memset(&msg_fmt1, 0, sizeof(E2SM_HelloWorld_IndicationMessage_Format1_t));



    indication_head = 0;
    indication_head = ( E2SM_HelloWorld_IndicationHeader_t *)calloc(1, sizeof( E2SM_HelloWorld_IndicationHeader_t));
    assert(indication_head != 0);

    indication_msg = 0;
    indication_msg = (E2SM_HelloWorld_IndicationMessage_t*)calloc(1, sizeof(E2SM_HelloWorld_IndicationMessage_t));
    assert(indication_msg !=0);

    errbuf_len = 128;
  };

 e2sm_indication::~e2sm_indication(void){

  mdclog_write(MDCLOG_DEBUG, "Freeing event trigger object memory");

  indication_head->choice.indicationHeader_Format1 = 0;

  indication_msg->choice.indicationMessage_Format1 = 0;

  ASN_STRUCT_FREE(asn_DEF_E2SM_HelloWorld_IndicationHeader, indication_head);
  ASN_STRUCT_FREE(asn_DEF_E2SM_HelloWorld_IndicationMessage, indication_msg);


};

bool e2sm_indication::encode_indication_header(unsigned char *buf, size_t *size, e2sm_indication_helper &helper){

  ASN_STRUCT_RESET(asn_DEF_E2SM_HelloWorld_IndicationHeader, indication_head);

  bool res;
  res = set_fields(indication_head, helper);
  if (!res){

    return false;
  }

  int ret_constr = asn_check_constraints(&asn_DEF_E2SM_HelloWorld_IndicationHeader, indication_head, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    return false;
  }

  xer_fprint(stdout, &asn_DEF_E2SM_HelloWorld_IndicationHeader, indication_head);

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_HelloWorld_IndicationHeader, indication_head, buf, *size);

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

bool e2sm_indication::encode_indication_message(unsigned char *buf, size_t *size, e2sm_indication_helper &helper){

  bool res;
  res = set_fields(indication_msg, helper);
  if (!res){
    return false;
  }


  int ret_constr = asn_check_constraints(&asn_DEF_E2SM_HelloWorld_IndicationMessage, indication_msg, errbuf, &errbuf_len);
  if(ret_constr){
    error_string.assign(&errbuf[0], errbuf_len);
    return false;
  }

  xer_fprint(stdout, &asn_DEF_E2SM_HelloWorld_IndicationMessage, indication_msg);

  asn_enc_rval_t retval = asn_encode_to_buffer(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_HelloWorld_IndicationMessage, indication_msg, buf, *size);

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

bool e2sm_indication::set_fields(E2SM_HelloWorld_IndicationHeader_t * ref_indication_head, e2sm_indication_helper & helper){

 if(ref_indication_head == 0){
    error_string = "Invalid reference for Event Trigger Definition set fields";
    return false;
  }

  ref_indication_head->present = E2SM_HelloWorld_IndicationHeader_PR_indicationHeader_Format1;

  head_fmt1.indicationHeaderParam = helper.header;

  ref_indication_head->choice.indicationHeader_Format1 = &head_fmt1;

  return true;
};

bool e2sm_indication::set_fields(E2SM_HelloWorld_IndicationMessage_t * ref_indication_msg, e2sm_indication_helper & helper){

 if(ref_indication_msg == 0){
    error_string = "Invalid reference for Event Action Definition set fields";
    return false;
  }
  ref_indication_msg->present = E2SM_HelloWorld_IndicationMessage_PR_indicationMessage_Format1;

  msg_fmt1.indicationMsgParam.buf = helper.message;
  msg_fmt1.indicationMsgParam.size = helper.message_len;


  ref_indication_msg->choice.indicationMessage_Format1 = &msg_fmt1;


  return true;
};

bool e2sm_indication::get_fields(E2SM_HelloWorld_IndicationHeader_t * ref_indictaion_header, e2sm_indication_helper & helper){

	if (ref_indictaion_header == 0){
	    error_string = "Invalid reference for Indication Header get fields";
	    return false;
	  }

	helper.header = ref_indictaion_header->choice.indicationHeader_Format1->indicationHeaderParam;
	return true;
}

bool e2sm_indication::get_fields(E2SM_HelloWorld_IndicationMessage_t * ref_indication_message, e2sm_indication_helper & helper){

	  if (ref_indication_message == 0){
	  	    error_string = "Invalid reference for Indication Message get fields";
	  	    return false;
	  	  }
	  helper.message = ref_indication_message->choice.indicationMessage_Format1->indicationMsgParam.buf;
	  helper.message_len = ref_indication_message->choice.indicationMessage_Format1->indicationMsgParam.size;

	  return true;
  }

