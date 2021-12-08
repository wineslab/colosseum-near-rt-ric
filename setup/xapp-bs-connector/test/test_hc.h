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
 * test_a1.h
 *
 *  Created on: Mar, 2020
 *  Author: Shraboni Jana
 */

#include<iostream>
#include<gtest/gtest.h>
#include<rapidjson/stringbuffer.h>
#include<rapidjson/writer.h>
#include<string.h>
#include "xapp.hpp"
#define HC_MSG_SIZE 512


using namespace std;
TEST(Xapp, RMRHealthCheck){

	 int total_num_msgs = 2;
	 int num_attempts = 10;

	 std::unique_ptr<XappRmr> rmr;
	 rmr = std::make_unique<XappRmr>("4560",num_attempts);
	 rmr->xapp_rmr_init(true);

	 XappSettings config;

	 std::unique_ptr<Xapp> hw_xapp = std::make_unique<Xapp>(std::ref(config),std::ref(*rmr));

	 std::unique_ptr<XappMsgHandler> mp_handler = std::make_unique<XappMsgHandler>("HW-Xapp-id");

	 hw_xapp->start_xapp_receiver(std::ref(*mp_handler));
	 sleep(5);

	 xapp_rmr_header hdr;
	 hdr.message_type = RIC_HEALTH_CHECK_REQ;
	 char strMsg[HC_MSG_SIZE];

	 for(int i = 0; i < total_num_msgs; i++){
		 snprintf(strMsg,HC_MSG_SIZE, "HelloWorld: RMR Health Check %d", i);
		 clock_gettime(CLOCK_REALTIME, &(hdr.ts));
		 hdr.payload_length = strlen(strMsg);

		 bool res = rmr->xapp_rmr_send(&hdr,(void*)strMsg);
		 usleep(1);
	 }
	 sleep(2);
	 hw_xapp->stop();

};

TEST(Xapp, A1HealthCheck){

	//Read the json file and send it using rmr.
	//string json = "{\"policy_type_id\": \"1\",\"policy_instance_id\":\"3d2157af-6a8f-4a7c-810f-38c2f824bf12\",\"operation\": \"CREATE\"}";
	string json = "{\"operation\": \"CREATE\", \"policy_type_id\": 1, \"policy_instance_id\": \"hwpolicy321\", \"payload\": {\"threshold\": 5}}";
	int n = json.length();
	char strMsg[n + 1];
	strcpy(strMsg, json.c_str());
	Document d;
	d.Parse(strMsg);

	int num_attempts = 5;

	std::unique_ptr<XappRmr> rmr;
	rmr = std::make_unique<XappRmr>("4560",num_attempts);
	rmr->xapp_rmr_init(true);

	XappSettings config;

	std::unique_ptr<Xapp> hw_xapp = std::make_unique<Xapp>(std::ref(config),std::ref(*rmr));

	std::unique_ptr<XappMsgHandler> mp_handler = std::make_unique<XappMsgHandler>("HW-Xapp-id");

	hw_xapp->start_xapp_receiver(std::ref(*mp_handler));
	sleep(5);

	xapp_rmr_header hdr;
	hdr.message_type = A1_POLICY_REQ;
    clock_gettime(CLOCK_REALTIME, &(hdr.ts));


	hdr.payload_length = strlen(strMsg);

	bool res_msg1 = rmr->xapp_rmr_send(&hdr,(void*)strMsg);
	ASSERT_TRUE(res_msg1);

	usleep(1);

	bool res_msg2 = rmr->xapp_rmr_send(&hdr,(void*)strMsg);
	ASSERT_TRUE(res_msg2);

	sleep(2);
	hw_xapp->stop();
}
