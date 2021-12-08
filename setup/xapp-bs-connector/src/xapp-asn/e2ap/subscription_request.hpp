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

#pragma once

#ifndef S_REQUEST_
#define S_REQUEST_

#include <mdclog/mdclog.h>
#include <vector>
#include <sstream>

#include <asn_application.h>
#include <E2AP-PDU.h>
#include <InitiatingMessage.h>
#include <RICsubscriptionRequest.h>
#include <RICsubscriptionRequest.h>
#include <ProtocolIE-Field.h>
#include <ProtocolIE-SingleContainer.h>
#include <RICactions-ToBeSetup-List.h>
#include <RICsubsequentAction.h>
#include "subscription_helper.hpp"

#define NUM_SUBSCRIPTION_REQUEST_IES 3
#define INITIAL_REQUEST_LIST_SIZE 4
  
class subscription_request{   
public:

  subscription_request(std::string name);
  subscription_request(void);
  ~subscription_request(void);
  
  bool encode_e2ap_subscription(unsigned char *, size_t *,  subscription_helper &);
  bool set_fields(InitiatingMessage_t *, subscription_helper &);
  bool get_fields(InitiatingMessage_t *, subscription_helper &);
    
  std::string get_error(void) const{
    return error_string;
  }
    
private:
    
  InitiatingMessage_t *initMsg;
  E2AP_PDU_t * e2ap_pdu_obj;

  RICsubscriptionRequest_IEs_t * IE_array;
  RICaction_ToBeSetup_ItemIEs_t * action_array;
  unsigned int action_array_size;  
  char errbuf[128];
  size_t errbuf_len = 128;
  std::string _name;
  std::string error_string;
};



#endif
