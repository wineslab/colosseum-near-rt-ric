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

#ifndef S_RESPONSE_
#define S_RESPONSE_

#include <mdclog/mdclog.h>
#include <vector>
#include <iostream>
#include <sstream>
#include <asn_application.h>
#include <E2N_E2AP-PDU.h>
#include <E2N_SuccessfulOutcome.h>
#include <E2N_UnsuccessfulOutcome.h>
#include <E2N_ProtocolIE-Field.h>
#include <E2N_ProtocolIE-Single-Container.h>
#include <E2N_ProcedureCode.h>
#include "response_helper.hpp"

#define NUM_SUBSCRIPTION_RESPONSE_IES 4
#define NUM_SUBSCRIPTION_FAILURE_IES 3
#define INITIAL_RESPONSE_LIST_SIZE 4
  
class subscription_response{   
public:
    
  subscription_response(void);
  ~subscription_response(void);
    
  bool encode_e2ap_subscription_response(unsigned char *, size_t *,  subscription_response_helper &, bool);
  void get_fields(E2N_SuccessfulOutcome_t *, subscription_response_helper &);    
  void get_fields(E2N_UnsuccessfulOutcome_t *, subscription_response_helper &);
  
  std::string get_error(void) const{
    return error_string;
  }
    
private:

  void set_fields_success( subscription_response_helper &);
  void set_fields_unsuccess( subscription_response_helper &);

  E2N_E2AP_PDU_t * e2ap_pdu_obj;
  E2N_SuccessfulOutcome_t * successMsg;
  E2N_UnsuccessfulOutcome_t * unsuccessMsg;
    

  E2N_RICsubscriptionResponse_IEs_t *IE_array;
  E2N_RICsubscriptionFailure_IEs_t *IE_Failure_array;
  

  E2N_RICaction_Admitted_ItemIEs_t * ie_admitted_list;
  E2N_RICaction_NotAdmitted_ItemIEs_t * ie_not_admitted_list;
  unsigned int ie_admitted_list_size, ie_not_admitted_list_size;
  
  char errbuf[128];
  size_t errbuf_len = 128;
  std::string error_string;
};




#endif
