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

#ifndef S_DEL_RESPONSE_
#define S_DEL_RESPONSE_

#include <mdclog/mdclog.h>
#include <vector>
#include <iostream>
#include <sstream>
#include <mdclog/mdclog.h>
#include <asn_application.h>
#include <E2AP-PDU.h>
#include <SuccessfulOutcome.h>
#include <UnsuccessfulOutcome.h>
#include <RICsubscriptionDeleteResponse.h>
#include <RICsubscriptionDeleteFailure.h>
#include <ProtocolIE-Field.h>
#include <ProcedureCode.h>
#include "response_helper.hpp"

#define NUM_SUBSCRIPTION_DELETE_RESPONSE_IES 2
#define NUM_SUBSCRIPTION_DELETE_FAILURE_IES 2
  
class subscription_delete_response {   
public:
    
  subscription_delete_response(void);
  ~subscription_delete_response(void);
    
  bool encode_e2ap_subscription_delete_response(unsigned char *, size_t *,  subscription_response_helper &, bool);
  bool set_fields(SuccessfulOutcome_t *, subscription_response_helper &);
  bool get_fields(SuccessfulOutcome_t *, subscription_response_helper &);
    
  bool set_fields(UnsuccessfulOutcome_t *, subscription_response_helper &);
  bool get_fields(UnsuccessfulOutcome_t *, subscription_response_helper &);
  
  std::string get_error_string(void) const {
    return error_string;
  }
    
private:

  E2AP_PDU_t * e2ap_pdu_obj;
  SuccessfulOutcome_t * successMsg;
  UnsuccessfulOutcome_t * unsuccessMsg;
    

  RICsubscriptionDeleteResponse_IEs_t *IE_array;
  RICsubscriptionDeleteFailure_IEs_t *IE_Failure_array;
  
  
  char errbuf[128];
  size_t errbuf_len = 128;
  std::string error_string;
  
};




#endif
