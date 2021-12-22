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

#pragma once

#ifndef S_REQUEST_
#define S_REQUEST_

#include <mdclog/mdclog.h>
#include <vector>
#include <sstream>

#include <asn_application.h>
#include <E2N_E2AP-PDU.h>
#include <E2N_InitiatingMessage.h>
#include <E2N_RICsubscriptionRequest.h>
#include <E2N_RICsubscription.h>
#include <E2N_ProtocolIE-Field.h>
#include <E2N_ProtocolIE-Single-Container.h>
#include <E2N_RICactions-ToBeSetup-List.h>
#include <E2N_RICsubsequentAction.h>
#include "subscription_helper.hpp"

#define NUM_SUBSCRIPTION_REQUEST_IES 3
#define INITIAL_REQUEST_LIST_SIZE 4
  
class subscription_request{   
public:

  subscription_request(std::string name);
  subscription_request(void);
  ~subscription_request(void);
  
  bool encode_e2ap_subscription(unsigned char *, size_t *,  subscription_helper &);
  bool set_fields(E2N_InitiatingMessage_t *, subscription_helper &);
  bool get_fields(E2N_InitiatingMessage_t *, subscription_helper &);
    
  std::string get_error(void) const{
    return error_string;
  }
    
private:
    
  E2N_InitiatingMessage_t *initMsg;
  E2N_E2AP_PDU_t * e2ap_pdu_obj;

  E2N_RICsubscriptionRequest_IEs_t * IE_array;
  E2N_RICaction_ToBeSetup_ItemIEs_t * action_array;
  unsigned int action_array_size;  
  char errbuf[128];
  size_t errbuf_len = 128;
  std::string _name;
  std::string error_string;
};



#endif
