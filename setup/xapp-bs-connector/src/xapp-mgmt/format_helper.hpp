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
 * a1_policy.hpp
 *
 *  Created on: Mar, 2020
 *  Author: Shraboni Jana
 */

#ifndef SRC_XAPP_MGMT_A1MSG_A1_POLICY_HELPER_HPP_
#define SRC_XAPP_MGMT_A1MSG_A1_POLICY_HELPER_HPP_

#include <rapidjson/document.h>
#include <rapidjson/writer.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/schema.h>

using namespace rapidjson;

typedef struct a1_policy_helper a1_policy_helper;

struct a1_policy_helper{

	std::string operation;
	std::string policy_type_id;
	std::string policy_instance_id;
	std::string handler_id;
	std::string status;

};


#endif /* SRC_XAPP_FORMATS_A1MSG_A1_POLICY_HELPER_HPP_ */
