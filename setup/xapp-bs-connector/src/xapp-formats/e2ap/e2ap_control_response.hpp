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
/*
 * ric_indication.h
 *
 *  Created on: Jul 11, 2019
 *      Author: sjana, Ashwin Sridharan
 */

#ifndef E2AP_RIC_CONTROL_RESPONSE_H_
#define E2AP_RIC_CONTROL_RESPONSE_H_

  
#include <iostream>
#include <errno.h>
#include <mdclog/mdclog.h>
#include <sstream>
#include <E2N_E2AP-PDU.h>
#include <E2N_SuccessfulOutcome.h>
#include <E2N_UnsuccessfulOutcome.h>
#include <E2N_RICcontrolAcknowledge.h>
#include <E2N_RICcontrolFailure.h>
#include <E2N_ProtocolIE-Field.h>
#include "e2ap_control_helper.hpp"

#define NUM_CONTROL_ACKNOWLEDGE_IES 3
#define NUM_CONTROL_FAILURE_IES 3

  
class ric_control_response{
    
public:
  ric_control_response(void);
  ~ric_control_response(void);
  
  bool encode_e2ap_control_response(unsigned char *, size_t *,  ric_control_helper &, bool);


  bool set_fields(E2N_SuccessfulOutcome_t *, ric_control_helper &);
  bool get_fields(E2N_SuccessfulOutcome_t *, ric_control_helper &);

  bool set_fields(E2N_UnsuccessfulOutcome_t *, ric_control_helper &);
  bool get_fields(E2N_UnsuccessfulOutcome_t *, ric_control_helper &);
  
  std::string get_error(void) const {return error_string ; };

private:
  
  E2N_E2AP_PDU_t * e2ap_pdu_obj;
  E2N_SuccessfulOutcome_t * successMsg;
  E2N_UnsuccessfulOutcome_t * unsuccessMsg;
  
  E2N_RICcontrolAcknowledge_IEs_t *IE_array;
  E2N_RICcontrolFailure_IEs_t *IE_failure_array;
  
  std::string error_string;
  
  char errbuf[128];
  size_t errbuf_len = 128;
};


#endif /* E2AP_RIC_CONTROL_RESPONSE_H_ */
