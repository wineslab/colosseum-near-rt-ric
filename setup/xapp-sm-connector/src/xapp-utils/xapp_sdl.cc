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
 * xapp_sdl.cc
 *
 *  Created on: Mar, 2020
 *  Author: Shraboni Jana
 */
#include "xapp_sdl.hpp"
/*need to work on the SDL FLow. Currently data hardcoded.
An xApp can use the SDL for two things:
- persisting state for itself (in case it fails and recovers)
- making information available for other xApps. The xApp would typically write using SDL directly.
- The consumer of the data could also use SDL directly or use an access library like in the case of the R-NIB.
*/
bool XappSDL::set_data(shareddatalayer::SyncStorage *sdl){
	try{
		//connecting to the Redis and generating a random key for namespace "hwxapp"
		mdclog_write(MDCLOG_INFO,  "IN SDL Set Data", __FILE__, __LINE__);
		DataMap dmap;
		char key[4]="abc";
		std::cout << "KEY: "<< key << std::endl;
		Key k = key;
		Data d;
		uint8_t num = 101;
		d.push_back(num);
		dmap.insert({k,d});
		Namespace ns(sdl_namespace);
		sdl->set(ns, dmap);
	}
	catch(...){
		mdclog_write(MDCLOG_ERR,  "SDL Error in Set Data for Namespace=%s",sdl_namespace);
		return false;
	}
	return true;
}

void XappSDL::get_data(shareddatalayer::SyncStorage *sdl){
	Namespace ns(sdl_namespace);
	DataMap dmap;
	std::string prefix="";
	Keys K = sdl->findKeys(ns, prefix);	// just the prefix
	DataMap Dk = sdl->get(ns, K);
	for(auto si=K.begin();si!=K.end();++si){
		std::vector<uint8_t> val_v = Dk[(*si)]; // 4 lines to unpack a string
		char val[val_v.size()+1];				// from Data
		int i;
		for(i=0;i<val_v.size();++i) val[i] = (char)(val_v[i]);
		val[i]='\0';
		printf("KEYS and Values %s = %s\n",(*si).c_str(), val);
	}

	mdclog_write(MDCLOG_INFO,  "IN SDL Get Data", __FILE__, __LINE__);
}
