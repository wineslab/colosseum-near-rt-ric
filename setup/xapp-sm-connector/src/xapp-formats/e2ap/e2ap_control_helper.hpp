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

#ifndef CONTROL_HELPER_H
#define CONTROL_HELPER_H

// control and indication helper objects are very similar and can be merged into one
// currently leaving them as two distnict entities till final design becomes clear

typedef struct ric_control_helper ric_control_helper;

struct ric_control_helper{
  ric_control_helper(void):req_id(1), req_seq_no(1), func_id(0), action_id(1), control_ack(-1), cause(0), sub_cause(0), control_status(1), control_msg(0), control_msg_size(0), control_header(0), control_header_size(0), call_process_id(0), call_process_id_size(0){};
  
  long int req_id, req_seq_no, func_id, action_id,  control_ack, cause, sub_cause, control_status;
  
  unsigned char* control_msg;
  size_t control_msg_size;
  
  unsigned char* control_header;
  size_t control_header_size;
  
  unsigned char *call_process_id;
  size_t call_process_id_size;
  
};

#endif
