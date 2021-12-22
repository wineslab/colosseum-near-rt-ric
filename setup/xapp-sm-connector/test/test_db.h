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
/*
 * test_rnib.h
 *
 *  Created on: Mar 23, 2020
 *  Author: Shraboni Jana
 */

#ifndef TEST_TEST_DB_H_
#define TEST_TEST_DB_H_

#include<iostream>
#include<gtest/gtest.h>
#include "xapp.hpp"

using namespace std;

TEST(Xapp, getGNBlist)
{
	XappSettings config;
	XappRmr rmr("7001");

	Xapp hw_xapp(std::ref(config),rmr);
	hw_xapp.set_rnib_gnblist();
	auto gnblist = hw_xapp.get_rnib_gnblist();
	int sz = gnblist.size();
	EXPECT_GE(sz,0);
	std::cout << "************gnb ids retrieved using R-NIB**************" << std::endl;
	for(int i = 0; i<sz; i++){
		std::cout << gnblist[i] << std::endl;
	}

}

TEST(Xapp, SDLData){

	//Xapp's SDL namespace.
    std::string nmspace = "hw-xapp";
	XappSDL xappsdl(nmspace);

	std::unique_ptr<shareddatalayer::SyncStorage> sdl(shareddatalayer::SyncStorage::create());
	bool res = xappsdl.set_data(sdl.get());
	ASSERT_TRUE(res);

	xappsdl.get_data(sdl.get());

}



#endif /* TEST_TEST_DB_H_ */
