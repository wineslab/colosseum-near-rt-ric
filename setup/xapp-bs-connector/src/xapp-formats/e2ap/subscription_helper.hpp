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


#ifndef SUB_HELPER_
#define SUB_HELPER_

/* 
   Simple structure to store action related information based on E2 v0.22
   Used for subscription request, response etc
   
   ricActionID					RICactionID,
   ricActionType		  		RICactionType,
   ricActionDefinition			RICactionDefinition 	OPTIONAL,
   ricSubsequentAction			RICsubsequentAction 	OPTIONAL,
   ricCause
*/

#include <iostream>
#include <vector>
#include <memory>

#include "generic_helpers.hpp"


// Note : if no action definition specified, octet length of action definition  is NULL
// If no subsequent action specified, default is subsequent_action = 0, time to wait is 0
struct Action {

public:
  
  Action(int id, int type): _is_def(false), _is_subs_act(false), _id(id), _type(type), _next_action(0), _wait(0){};
  Action(int id, int type, const void *def, size_t def_size, int next, int wait): _is_def(false), _is_subs_act(false), _id(id), _type(type){
    
    if (def_size > 0){
      _is_def = true;
      _action_definition.set_ref(def);
      _action_definition.set_size(def_size);
    }
    
    if(next >= 0 && wait >= 0){
      _is_subs_act = true;
      _next_action = next;
      _wait = wait;
    }
  };

  
  int get_id() const{
    return _id;
  }

  int get_type() const {
    return _type;
  }


  const void * get_definition(void )  {
    return _action_definition.get_ref();
  }

  int get_definition_size(void) const {
    return _action_definition.get_size();
  };
  

  int get_subsequent_action() const {
    return _next_action;
  };

  int get_wait() const {
    return _wait;
  }

  bool is_definition() const{

    return _is_def;
  }

  bool is_subsequent_action() const{
    return _is_subs_act;
  }
    
private:

  bool _is_def;
  bool _is_subs_act;
  int _id, _type, _next_action, _wait, _cause, _sub_cause;
  bool _is_admit;
  octet_helper _action_definition;

};


/*
 Helper class that stores subscription data 
*/


struct subscription_helper {

public:

  using action_t = std::vector<Action>;
  
  subscription_helper(){
    _action_ref = std::make_unique<action_t>();
    curr_index = 0;    
  };
  
  action_t * get_list() const {return _action_ref.get();};

  void clear(void){
    _action_ref.get()->clear();
  }
  
  void set_request(int id, int seq_no){
    _req_id = id;
    _req_seq_no = seq_no;
    
  };

  void set_function_id(int id){
    _func_id = id;
  };

  void set_event_def(const void *ref, size_t size){
    _event_def.set_ref(ref);
    _event_def.set_size(size);
   };

 
  void add_action(int id, int type){
    Action a(id, type) ;
    _action_ref.get()->push_back(a);
  };

  void add_action(int id, int type, std::string action_def, int next_action, int wait_time){
    Action a (id, type, action_def.c_str(), action_def.length(), next_action, wait_time);
    _action_ref.get()->push_back(a);
  };


  int  get_request_id(void) const{
    return _req_id;
  }

  int  get_req_seq(void) const {
    return _req_seq_no;
  }

  int  get_function_id(void) const{
    return _func_id;
  }
  
  const void * get_event_def(void)  {
    return _event_def.get_ref();
  }

  int get_event_def_size(void) const {
    return _event_def.get_size();
  }

  void print_sub_info(void){
    std::cout <<"Request ID = " << _req_id << std::endl;
    std::cout <<"Request Sequence Number = " << _req_seq_no << std::endl;
    std::cout <<"RAN Function ID = " << _func_id << std::endl;
    for(auto const & e: *(_action_ref.get())){
      std::cout <<"Action ID = " << e.get_id() << " Action Type = " << e.get_type() << std::endl;
    }
  };
  
private:
  
  std::unique_ptr<action_t> _action_ref;
  int curr_index;
  int _req_id, _req_seq_no, _func_id;
  octet_helper _event_def;
};

#endif
