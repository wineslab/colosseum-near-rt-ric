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
 * xapp_config.cc
 * Created on: 2019
 * Author: Ashwin Shridharan, Shraboni Jana
 */

#include "xapp_config.hpp"

string& XappSettings::operator[](const SettingName& theName){
    return theSettings[theName];
}

void XappSettings::loadCmdlineSettings(int argc, char **argv){

	   // Parse command line options to over ride
	  static struct option long_options[] =
	    {
	    		{"xappname", required_argument, 0, 'n'},
				{"xappid", required_argument, 0, 'x'},
				{"port", required_argument, 0, 'p'},
				{"threads", required_argument,    0, 't'},
				{"ves-interval", required_argument, 0, 'i'},
				{"gNodeB", required_argument, 0, 'g'}

	    };


	   while(1) {

		int option_index = 0;
		char c = getopt_long(argc, argv, "n:p:t:s:g:a:v:u:i:c:x:", long_options, &option_index);

	        if(c == -1){
		    break;
	         }

		switch(c)
		  {

		  case 'n':
		    theSettings[XAPP_NAME].assign(optarg);
		    break;

		  case 'p':
		    theSettings[HW_PORT].assign(optarg);
		    break;

		  case 't':
			theSettings[THREADS].assign(optarg);
		    mdclog_write(MDCLOG_INFO, "Number of threads set to %s from command line e\n", theSettings[THREADS].c_str());
		    break;


		  case 'x':
		    theSettings[XAPP_ID].assign(optarg);
		    mdclog_write(MDCLOG_INFO, "XAPP ID set to  %s from command line ", theSettings[XAPP_ID].c_str());
		    break;

		  case 'h':
		    usage(argv[0]);
		    exit(0);

		  default:
		    usage(argv[0]);
		    exit(1);
		  }
	   };

}

void XappSettings::loadDefaultSettings(){


		 if(theSettings[XAPP_NAME].empty()){
		  theSettings[XAPP_NAME] = DEFAULT_XAPP_NAME;
		  }

	  	  if(theSettings[XAPP_ID].empty()){
	  		  theSettings[XAPP_ID] = DEFAULT_XAPP_NAME; //for now xapp_id is same as xapp_name since single xapp instance.
	  	  }
	  	  if(theSettings[LOG_LEVEL].empty()){
	  		  theSettings[LOG_LEVEL] = DEFAULT_LOG_LEVEL;
	  	  }
	  	  if(theSettings[HW_PORT].empty()){
	  		  theSettings[HW_PORT] = DEFAULT_PORT;
	  	  }
	  	  if(theSettings[MSG_MAX_BUFFER].empty()){
	  		  theSettings[MSG_MAX_BUFFER] = DEFAULT_MSG_MAX_BUFFER;
	  	  }

	  	 if(theSettings[THREADS].empty()){
	  		  		  theSettings[THREADS] = DEFAULT_THREADS;
	  		  	  }


}

void XappSettings::loadEnvVarSettings(){

	  if (const char *env_xname = std::getenv("XAPP_NAME")){
		  theSettings[XAPP_NAME].assign(env_xname);
		  mdclog_write(MDCLOG_INFO,"Xapp Name set to %s from environment variable", theSettings[XAPP_NAME].c_str());
	  }
	  if (const char *env_xid = std::getenv("XAPP_NAME")){
		   theSettings[XAPP_ID].assign(env_xid);
		   mdclog_write(MDCLOG_INFO,"Xapp ID set to %s from environment variable", theSettings[XAPP_ID].c_str());
	  }

	  if (const char *env_ports = std::getenv("HW_PORT")){
		  theSettings[HW_PORT].assign(env_ports);
	 	  mdclog_write(MDCLOG_INFO,"Ports set to %s from environment variable", theSettings[HW_PORT].c_str());
	  }
	  if (const char *env_ports = std::getenv("MSG_MAX_BUFFER")){
	 		  theSettings[MSG_MAX_BUFFER].assign(env_ports);
	 	 	  mdclog_write(MDCLOG_INFO,"Ports set to %s from environment variable", theSettings[MSG_MAX_BUFFER].c_str());
	 	  }

}

void XappSettings::usage(char *command){
	std::cout <<"Usage : " << command << " " << std::endl;
	std::cout <<" --name[-n] xapp_instance_name "<< std::endl;
    std::cout <<" --port[-p] port to listen on e.g tcp:4561  "<< std::endl;
    std::cout << "--threads[-t] number of listener threads "<< std::endl ;

}
