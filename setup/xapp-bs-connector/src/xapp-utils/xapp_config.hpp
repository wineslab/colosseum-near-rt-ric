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
 * xapp_config.hpp
 * Created on: 2019
 * Author: Ashwin Shridharan, Shraboni Jana
 */

#ifndef SRC_XAPP_CONFIG_XAPP_CONFIG_HPP_
#define SRC_XAPP_CONFIG_XAPP_CONFIG_HPP_

#include <getopt.h>
#include <map>
#include <iostream>
#include <cstdlib>
#include <mdclog/mdclog.h>

#define DEFAULT_XAPP_NAME "hwxapp"
#define DEFAULT_PORT "4560"
#define DEFAULT_MSG_MAX_BUFFER "2072"
#define DEFAULT_THREADS "1"

#define DEFAULT_LOG_LEVEL	MDCLOG_WARN

#define ASN_BUFF_MAX_SIZE		4096
#define MAX_SUBSCRIPTION_ATTEMPTS	10
#define HELLOWORLD_POLICY_ID 2

using namespace std;

struct XappSettings{

public:
	typedef enum{
		  XAPP_ID,
		  XAPP_NAME,
		  HW_PORT,
		  MSG_MAX_BUFFER,
		  THREADS,
		  LOG_LEVEL
	}SettingName;

	void loadDefaultSettings();
	void loadCmdlineSettings(int, char **);
	void loadEnvVarSettings();
	void usage(char*);
	string& operator[](const SettingName& theName);
private:
	typedef map<SettingName, std::string> SettingCollection;
	SettingCollection theSettings;

};



#endif /* SRC_XAPP_CONFIG_XAPP_CONFIG_HPP_ */
