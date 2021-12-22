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
 * ric_indication.h
 *
 *  Created on: Jul 11, 2019
 *      Author: sjana, Ashwin Sridharan
 */


#ifndef E2AP_INDICATION_HELPER_
#define E2AP_INDICATION_HELPER_

typedef struct ric_indication_helper ric_indication_helper;

struct ric_indication_helper{
  ric_indication_helper(void) : req_id(1), req_seq_no(1), func_id(0), action_id(1), indication_type(0), indication_sn(0), indication_msg(0), indication_msg_size(0), indication_header(0), indication_header_size(0), call_process_id(0), call_process_id_size(0) {};
  long int req_id, req_seq_no, func_id, action_id, indication_type, indication_sn;
  
  unsigned char* indication_msg;
  size_t indication_msg_size;
  
  unsigned char* indication_header;
  size_t indication_header_size;
  
  unsigned char *call_process_id;
  size_t call_process_id_size;
  
};

#endif
