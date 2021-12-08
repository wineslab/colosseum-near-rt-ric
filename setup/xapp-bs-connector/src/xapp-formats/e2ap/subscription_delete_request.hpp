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

#ifndef S_DELETE_
#define S_DELETE_

#include <mdclog/mdclog.h>
#include <vector>
#include <sstream>
#include <mdclog/mdclog.h>
#include <asn_application.h>
#include <E2N_E2AP-PDU.h>
#include <E2N_InitiatingMessage.h>
#include <E2N_RICsubscriptionDeleteRequest.h>
#include <E2N_ProtocolIE-Field.h>
#include "subscription_helper.hpp"

#define NUM_SUBSCRIPTION_DELETE_IES 2

class subscription_delete{   
public:

  subscription_delete(void);
  ~subscription_delete(void);
  
  bool encode_e2ap_subscription(unsigned char *, size_t *,  subscription_helper &);
  bool set_fields(subscription_helper &);
  bool get_fields(E2N_InitiatingMessage_t *, subscription_helper &);
    
  std::string get_error(void) const {
    return error_string ;
  }
    
private:
    
  E2N_InitiatingMessage_t *initMsg;
  E2N_E2AP_PDU_t * e2ap_pdu_obj;

  E2N_RICsubscriptionDeleteRequest_IEs_t * IE_array;

  
  char errbuf[128];
  size_t errbuf_len = 128;
  std::string _name;
  std::string error_string;
};



#endif
