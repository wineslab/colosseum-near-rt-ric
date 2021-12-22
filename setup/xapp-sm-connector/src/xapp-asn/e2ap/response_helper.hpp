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

#ifndef S_RESPONSE_HELPER_
#define S_RESPONSE_HELPER_

#include <vector>
#include <memory>

/* Simple structure to store action for RICaction of the Subscription response based on E2 v0.31 */
struct ActionResponse {
public:
  ActionResponse(int id): _is_admit(true), _id(id), _cause(-1), _sub_cause(-1){};
  ActionResponse(int id, int cause, int sub_cause): _is_admit(false), _id(id), _cause(cause), _sub_cause(sub_cause){};
  
  int get_id() const{
    return _id;
  };

  int get_cause() const{
    return _cause;
  };

  int get_sub_cause() const{
    return _sub_cause;
  };

  bool is_admitted(void){
    return _is_admit;
  };
  
private:

  bool _is_admit;
  int _id, _cause, _sub_cause;
  
};


struct subscription_response_helper {
  
public:

  using action_t = std::vector<ActionResponse>;
  
  subscription_response_helper(void){
    _action_admitted_ref = std::make_unique<action_t>();
    _action_not_admitted_ref = std::make_unique<action_t>();
    
  };
  
  // copy operator
  subscription_response_helper(const subscription_response_helper &he ){
    _action_admitted_ref = std::make_unique<action_t>();
    _action_not_admitted_ref = std::make_unique<action_t>();
    
    _req_id = he.get_request_id();
    _req_seq_no = he.get_req_seq();
    _func_id = he.get_function_id();
    
    // Take care of the actions
    for (auto const & e: *(he.get_admitted_list())){
      add_action(e.get_id());
    }
    
    for(auto const  & e: *(he.get_not_admitted_list())){
      add_action(e.get_id(), e.get_cause(), e.get_sub_cause());
    };
  }
  

  // assignment operator
  void operator=(const subscription_response_helper & he){
    _action_admitted_ref = std::make_unique<action_t>();
    _action_not_admitted_ref = std::make_unique<action_t>();
    
    _req_id = he.get_request_id();
    _req_seq_no = he.get_req_seq();
    _func_id = he.get_function_id();
    
    
    // Take care of the actions
    for (auto  const & e: *(he.get_admitted_list())){
      add_action(e.get_id());
    }
  
    for(auto const  & e: *(he.get_not_admitted_list())){
      add_action(e.get_id(), e.get_cause(), e.get_sub_cause());
    };
    
  }
  
  action_t * get_admitted_list (void ) const {return _action_admitted_ref.get();};
  action_t * get_not_admitted_list (void ) const{return _action_not_admitted_ref.get();};
  
  void set_request(int id, int seq_no){
    _req_id = id;
    _req_seq_no = seq_no;
    
  };

  void clear(void){
    _action_admitted_ref.get()->clear();
    _action_not_admitted_ref.get()->clear();
  }

  
  void set_function_id(int id){
    _func_id = id;
  };

  void add_action(int id){
    ActionResponse a(id) ;
    _action_admitted_ref.get()->push_back(a);
  };

  void add_action(int id, int cause, int sub_cause){
    ActionResponse a (id, cause, sub_cause);
    _action_not_admitted_ref.get()->push_back(a);
  };


  int  get_request_id(void) const{
    return _req_id;
  }
  
  int get_req_seq(void) const{
    return _req_seq_no;
  }

  int  get_function_id(void) const{
    return _func_id;
  }
 
  std::string  to_string(void){
    std::string Info;
    Info += "Request ID = " + std::to_string(_req_id) + "\n";
    Info += "Request Sequence No = "  + std::to_string(_req_seq_no) + "\n";
    Info += "RAN Function ID = " + std::to_string(_func_id) + "\n";
    Info += "Actions Admitted =\n";
    int i = 0;
    for(auto & e: *(_action_admitted_ref)){
        Info += std::to_string(i)  + ": ID=" + std::to_string(e.get_id()) + "\n";
        i++;
    }    
    Info += "Actions Not Admitted =\n";
    i = 0;
    for(auto & e: *(_action_not_admitted_ref)){
      Info += std::to_string(i)  + ": ID=" + std::to_string(e.get_id()) +  ": Cause =" + std::to_string(e.get_cause()) + ": Sub-Cause=" + std::to_string(e.get_sub_cause()) + "\n";
      i++;
    }    
  
    return Info;
  } 

private:
  int _req_id, _req_seq_no, _func_id;
  std::unique_ptr<action_t> _action_admitted_ref;
  std::unique_ptr<action_t> _action_not_admitted_ref;
  
};
  

#endif
