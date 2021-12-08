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
 * subs_mgmt.cc
 * Created on: 2019
 * Author: Ashwin Shridharan, Shraboni Jana
 */
#include "subs_mgmt.hpp"

#include <errno.h>
#include <stdio.h>
#include <string.h>

SubscriptionHandler::SubscriptionHandler(unsigned int timeout_seconds):_time_out(std::chrono::seconds(timeout_seconds)){
	  _data_lock = std::make_unique<std::mutex>();
	  _cv = std::make_unique<std::condition_variable>();
};

void SubscriptionHandler::clear(void){
  {
    std::lock_guard<std::mutex> lock(*(_data_lock).get());
    status_table.clear();
  }
  
};


bool SubscriptionHandler::add_request_entry(transaction_identifier id, transaction_status status){

  // add entry in hash table if it does not exist
  //auto search = status_table.find(id);
  //if(search != status_table.end()){
  if(transaction_present(status_table, id)){
    return false;
  }
  std::cout << "add_request_entry " << id << std::endl;
  status_table[id] = status;
  return true;

};



bool SubscriptionHandler::delete_request_entry(transaction_identifier id){

  auto search = status_table.find(id);

  if (!trans_table.empty()) {
	  auto search2 = trans_table.find(id);
	  if(search2 !=trans_table.end()){
		  trans_table.erase(search2);
	  }
  }

  if (search != status_table.end()){
    status_table.erase(search);
    mdclog_write(MDCLOG_INFO,"Entry for Transaction ID deleted: %s", id);

    return true;
  }
  mdclog_write(MDCLOG_INFO,"Entry not found in SubscriptionHandler for Transaction ID: %d",id);

  return false;
};


bool SubscriptionHandler::set_request_status(transaction_identifier id, transaction_status status){

  // change status of a request only if it exists.
  //auto search = status_table.find(id);
  if(transaction_present(status_table, id)){
    status_table[id] = status;
    std::cout << "set_request_status " << id << " to status " << status << std::endl;
    return true;
  }

  return false;
  
};


int const SubscriptionHandler::get_request_status(transaction_identifier id){
  //auto search = status_table.find(id);
  auto search = find_transaction(status_table, id);
  std::cout << "get_request_status " << id << std::endl;
  if (search == status_table.end()){
    return -1;
  }

  return search->second;
}
				   


bool SubscriptionHandler::is_request_entry(transaction_identifier id){

  std::cout << "is_request_entry, looking for key: " << id << std::endl;

  if (transaction_present(status_table, id)) {
    std::cout << "Key found" << std::endl;
    return true;
  }
  else {
    std::cout << "Key NOT found" << std::endl;
    return false;
  }
}

// Handles subscription responses
void SubscriptionHandler::manage_subscription_response(int message_type, transaction_identifier id, const void *message_payload, size_t message_len){

  bool res;
  std::cout << "In Manage subscription" << std::endl;

  // wake up all waiting users ...
  if(is_request_entry(id)){
    std::cout << "In Manage subscription, inside if loop" << std::endl;
    std::cout << "Setting to request_success: " << request_success << std::endl;

    set_request_status(id, request_success);
     _cv.get()->notify_all();
  }

  // decode received message payload
  E2AP_PDU_t *pdu = nullptr;
  auto retval = asn_decode(nullptr, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2AP_PDU, (void **) &pdu, message_payload, message_len);
  
  // print decoded payload
  if (retval.code == RC_OK) {
    char *printBuffer;
    size_t size;
    FILE *stream = open_memstream(&printBuffer, &size);
    asn_fprint(stream, &asn_DEF_E2AP_PDU, pdu);
    mdclog_write(MDCLOG_DEBUG, "Decoded E2AP PDU: %s", printBuffer);
  }
}

std::unordered_map<transaction_identifier, transaction_status>::iterator find_transaction(std::unordered_map<transaction_identifier, transaction_status> map,
		transaction_identifier id) {
  auto iter = map.begin();
  while (iter != map.end()) {
    if (strcmp((const char*) iter->first, (const char*) id) == 0) {
      break;
    }
    ++iter;
  }

  return iter;
}

bool transaction_present(std::unordered_map<transaction_identifier, transaction_status> map, transaction_identifier id) {
  auto iter = map.begin();
  while (iter != map.end()) {
    if (strcmp((const char*) iter->first, (const char*) id) == 0) {
      return true;
    }
    ++iter;
  }

  return false;
}

