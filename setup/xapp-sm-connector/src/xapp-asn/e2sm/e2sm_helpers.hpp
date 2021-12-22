/*
  ==================================================================================

  Copyright (c) 2019-2020 AT&T Intellectual Property.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, softwares
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  ==================================================================================
*/
/*
 * e2sm_control.cc
 *
 *  Created on: Apr 30, 2020
 *      Author: Shraboni Jana
 */
/* Classes to handle E2 service model based on e2sm-HelloWorld-v001.asn */
#ifndef E2SM_HELPER_
#define E2SM_HELPER_

#include <errno.h>
#include <iostream>
#include <vector>
#include <sstream>
#include <memory>

typedef struct ranparam_helper ranparam_helper;
struct ranparam_helper {
	  long int _param_id;
	  unsigned char* _param_name;
	  size_t _param_name_len;
	  int _param_test;
	  unsigned char* _param_value;
	  size_t _param_value_len;

};
class RANParam_Helper{
private:
	ranparam_helper _ranparam_helper;

public:

	RANParam_Helper(int id, unsigned char *param_name, size_t param_name_len, int param_test, unsigned char* param_value, size_t param_value_len){
		_ranparam_helper._param_id = id;
		_ranparam_helper._param_name = param_name;
		_ranparam_helper._param_name_len = param_name_len;
		_ranparam_helper._param_test = param_test;
		_ranparam_helper._param_value = param_value;
		_ranparam_helper._param_value_len = param_value_len;
	  }

	const ranparam_helper & getran_helper() const {
		return _ranparam_helper;
	}
	void print_ranparam_info(void){
	    std::cout <<"Param ID = " << _ranparam_helper._param_id << std::endl;
	    std::cout << "Parame Name =" << _ranparam_helper._param_name << std::endl;
	    std::cout <<"Param Test = " << _ranparam_helper._param_test << std::endl;
	    std::cout <<"Param Value = " << _ranparam_helper._param_value << std::endl;
	}
};


using ranparam_helper_t = std::vector<RANParam_Helper>;

typedef struct e2sm_subscription_helper e2sm_subscription_helper;
struct e2sm_subscription_helper {
public:


  int triger_nature;
  ranparam_helper_t param;
  void add_param(int id, unsigned char *param_name, size_t param_name_len, int param_test, unsigned char* param_value, size_t param_value_len){
	  RANParam_Helper rparam(id,param_name,param_name_len,param_test,param_value,param_value_len);
      param.push_back(rparam);
    };
  ranparam_helper_t get_paramlist() const {return param;};



};


typedef struct e2sm_indication_helper e2sm_indication_helper;


struct e2sm_indication_helper {
	long int header;
	unsigned char* message;
	size_t message_len;
};

typedef struct e2sm_control_helper e2sm_control_helper;

struct e2sm_control_helper {
	long int header;
	unsigned char* message;
	size_t message_len;
};

#endif
