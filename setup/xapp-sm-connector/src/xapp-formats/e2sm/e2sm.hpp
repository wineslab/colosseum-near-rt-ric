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

/* Classes to handle E2 service model based on e2sm-gNB-X2-release-1-v040.asn */

#ifndef E2SM_
#define E2SM_


#include <sstream>
#include <mdclog/mdclog.h>
#include <E2N_E2SM-gNB-X2-indicationHeader.h>
#include <E2N_E2SM-gNB-X2-indicationMessage.h>
#include <E2N_E2SM-gNB-X2-controlHeader.h>
#include <E2N_E2SM-gNB-X2-controlMessage.h>
#include <E2N_E2SM-gNB-X2-eventTriggerDefinition.h>

#include <E2N_GlobalGNB-ID.h>
#include <E2N_TypeOfMessage.h>
#include <E2N_InterfaceProtocolIE-Item.h>

#include<E2N_InterfaceProtocolIE-ID.h>
#include<E2N_InterfaceProtocolIE-Value.h>
#include<E2N_InterfaceProtocolIE-Test.h>
#include "../../xapp-formats/e2sm/e2sm_helpers.hpp"

#define INITIAL_LIST_SIZE 4

  

  
/* builder class for E2SM event trigger definition */

class e2sm_event_trigger {
public:
  e2sm_event_trigger(void);
  ~e2sm_event_trigger(void);
    
  bool set_fields(E2N_E2SM_gNB_X2_eventTriggerDefinition_t *, e2sm_event_trigger_helper &);
  bool get_fields(E2N_E2SM_gNB_X2_eventTriggerDefinition_t *, e2sm_event_trigger_helper &);
  bool encode_event_trigger(unsigned char *, size_t *, e2sm_event_trigger_helper &);

  std::string  get_error (void) const {return error_string ;};
  
private:

  E2N_E2SM_gNB_X2_eventTriggerDefinition_t * event_trigger; // used for encoding
  E2N_GlobalGNB_ID_t gNodeB_ID;
  struct E2N_InterfaceProtocolIE_Item * ie_list;
  unsigned int ie_list_size;
    
  //std::vector<struct InterfaceProtocolIE_Item> ie_list;
  E2N_E2SM_gNB_X2_eventTriggerDefinition::E2N_E2SM_gNB_X2_eventTriggerDefinition__interfaceProtocolIE_List *condition_list;
    
  char errbuf[128];
  size_t errbuf_len;
  std::string error_string;
};
  
    
/* builder class for E2SM indication  using ASN1c */
  
class e2sm_indication {
public:
    
  e2sm_indication(void);
  ~e2sm_indication(void);
    
  E2N_E2SM_gNB_X2_indicationHeader_t * get_header(void);
  E2N_E2SM_gNB_X2_indicationMessage_t * get_message(void);

  bool set_header_fields(E2N_E2SM_gNB_X2_indicationHeader_t *, e2sm_header_helper &);
  bool get_header_fields(E2N_E2SM_gNB_X2_indicationHeader_t *, e2sm_header_helper &);
    
  bool set_message_fields(E2N_E2SM_gNB_X2_indicationMessage_t *, e2sm_message_helper &);
  bool get_message_fields(E2N_E2SM_gNB_X2_indicationMessage_t *, e2sm_message_helper &);

  bool encode_indication_header(unsigned char * , size_t * , e2sm_header_helper &); 
  bool encode_indication_message(unsigned char *, size_t *, e2sm_message_helper &);
  std::string  get_error (void) const {return error_string ; };
    
private:
  
  E2N_E2SM_gNB_X2_indicationHeader_t *header; // used for encoding
  E2N_E2SM_gNB_X2_indicationMessage_t *message; // used for encoding
    
  char errbuf[128];
  size_t errbuf_len;
  E2N_GlobalGNB_ID_t gNodeB_ID;
  std::string error_string;

  
};

/* builder class for E2SM control  using ASN1c */
  
class e2sm_control {
public:
    
  e2sm_control(void);
  ~e2sm_control(void);
    
  E2N_E2SM_gNB_X2_controlHeader_t * get_header(void);
  E2N_E2SM_gNB_X2_controlMessage_t * get_message(void);

  bool set_header_fields(E2N_E2SM_gNB_X2_controlHeader_t *, e2sm_header_helper &);
  bool get_header_fields(E2N_E2SM_gNB_X2_controlHeader_t *, e2sm_header_helper &);
    
  bool set_message_fields(E2N_E2SM_gNB_X2_controlMessage_t *, e2sm_message_helper &);
  bool get_message_fields(E2N_E2SM_gNB_X2_controlMessage_t *, e2sm_message_helper &);

  bool encode_control_header(unsigned char * , size_t * , e2sm_header_helper &); 
  bool encode_control_message(unsigned char *, size_t *, e2sm_message_helper &);
  std::string  get_error (void) const {return error_string ; };
    
private:
  
  E2N_E2SM_gNB_X2_controlHeader_t *header; // used for encoding
  E2N_E2SM_gNB_X2_controlMessage_t *message; // used for encoding
    
  char errbuf[128];
  size_t errbuf_len;
  E2N_GlobalGNB_ID_t gNodeB_ID;
  std::string error_string;

  
};

#endif
