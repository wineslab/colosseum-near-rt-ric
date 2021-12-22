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
/*
 * test_asn.h
 *
 *  Created on: Apr, 2020
 *      Author: Shraboni Jana
 */

#ifndef TEST_TEST_ASN_H_
#define TEST_TEST_ASN_H_
#include<iostream>
#include<gtest/gtest.h>

#include "subscription_request.hpp"
#include "xapp.hpp"
#include "e2sm_helpers.hpp"
#include "e2sm_subscription.hpp"
#include "e2sm_indication.hpp"
#include "e2sm_control.hpp"

using namespace std;
TEST(E2SM, SubscriptionRequest)
{

	unsigned char event_buf[128];
	size_t event_buf_len = 128;

	unsigned char act_buf[128];
	size_t act_buf_len = 128;

	bool res;


	e2sm_subscription_helper e2sm_subsdata;
	std::unique_ptr<ranparam_helper> *ranhelp;
	e2sm_subscription e2sm_subs;


	e2sm_subsdata.triger_nature = 0;

	int param_id = 1;
	unsigned char param_name[20];
	strcpy((char*)param_name,"ParamName");
	int param_name_len = strlen((const char*)param_name);

	int param_test = 0;
	unsigned char param_value[20];
	strcpy((char*)param_value,"ParamValue");
	int param_value_len = strlen((const char*)param_value);

	e2sm_subsdata.add_param(param_id, param_name, param_name_len, param_test, param_value, param_value_len);


	// Encode the event trigger definition
	res = e2sm_subs.encode_event_trigger(&event_buf[0], &event_buf_len, e2sm_subsdata);
	if(!res)
		std::cout << e2sm_subs.get_error() << std::endl;

	ASSERT_TRUE(res);

	// Encode the action defintion
	res = e2sm_subs.encode_action_defn(&act_buf[0], &act_buf_len, e2sm_subsdata);
	if(!res)
		std::cout << e2sm_subs.get_error() << std::endl;
	ASSERT_TRUE(res);

}
TEST(E2SM, IndicationMessage)
{

	unsigned char header_buf[128];
	size_t header_buf_len = 128;

	unsigned char msg_buf[128];
	size_t msg_buf_len = 128;

	bool res;
	asn_dec_rval_t retval;


	e2sm_indication_helper e2sm_inddata;
	e2sm_indication e2sm_inds;

	unsigned char msg[20] = "HelloWorld";

	e2sm_inddata.header = 1001;
	e2sm_inddata.message = msg;
	e2sm_inddata.message_len = strlen((const char*)e2sm_inddata.message);


	// Encode the indication header
	res = e2sm_inds.encode_indication_header(&header_buf[0], &header_buf_len, e2sm_inddata);
	if(!res)
		std::cout << e2sm_inds.get_error() << std::endl;

	ASSERT_TRUE(res);

	// Encode the indication message
	res = e2sm_inds.encode_indication_message(&msg_buf[0], &msg_buf_len, e2sm_inddata);
	if(!res)
		std::cout << e2sm_inds.get_error() << std::endl;
	ASSERT_TRUE(res);

	//decode the indication header
	e2sm_indication_helper e2sm_decodedata;


	E2SM_HelloWorld_IndicationHeader_t *header = 0; // used for decoding
	retval = asn_decode(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_HelloWorld_IndicationHeader, (void**)&(header), &header_buf[0], header_buf_len);

	ASSERT_TRUE(retval.code == RC_OK);
	res = e2sm_inds.get_fields(header, e2sm_decodedata);

	//decode the indication message

	E2SM_HelloWorld_IndicationMessage_t *mesg = 0; // used for decoding
	retval = asn_decode(0, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_HelloWorld_IndicationMessage, (void**)&(mesg), &msg_buf[0], msg_buf_len);

	ASSERT_TRUE(retval.code == RC_OK);
	res = e2sm_inds.get_fields(mesg, e2sm_decodedata);


	std::cout << "Indication Header:" << e2sm_decodedata.header << std::endl;
	std::cout << "Indication Message:" << e2sm_decodedata.message << std::endl;
	std::cout << "Indication Message Len:" << e2sm_decodedata.message_len << std::endl;

	ASSERT_EQ(e2sm_inddata.header, e2sm_decodedata.header);
	ASSERT_EQ(e2sm_inddata.message_len, e2sm_decodedata.message_len);
	for (int i = 0; i < e2sm_inddata.message_len; ++i) {
	  EXPECT_EQ(e2sm_inddata.message[i], e2sm_decodedata.message[i]) << "Encoded and Decoded Msg differ at index " << i;
	}


}

TEST(E2SM, ControlMessage)
{

	unsigned char header_buf[128];
	size_t header_buf_len = 128;

	unsigned char msg_buf[128];
	size_t msg_buf_len = 128;

	bool res;
	asn_dec_rval_t retval;


	e2sm_control_helper e2sm_cntrldata;
	e2sm_control e2sm_cntrl;

	unsigned char msg[20] = "HelloWorld";

	e2sm_cntrldata.header = 1001;
	e2sm_cntrldata.message = msg;
	e2sm_cntrldata.message_len = strlen((const char*)e2sm_cntrldata.message);


	// Encode the indication header
	res = e2sm_cntrl.encode_control_header(&header_buf[0], &header_buf_len, e2sm_cntrldata);
	if(!res)
		std::cout << e2sm_cntrl.get_error() << std::endl;

	ASSERT_TRUE(res);

	// Encode the indication message
	res = e2sm_cntrl.encode_control_message(&msg_buf[0], &msg_buf_len, e2sm_cntrldata);
	if(!res)
		std::cout << e2sm_cntrl.get_error() << std::endl;
	ASSERT_TRUE(res);
}

#endif /* TEST_TEST_ASN_H_ */
