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
 * hw_unit_tests.cc
 *
 *  Created on: Mar, 2020
 *  Author: Shraboni Jana
 */

#include<iostream>
#include<stdlib.h>
#include<gtest/gtest.h>

#include "test_db.h"
#include "test_rmr.h"
#include "test_hc.h"
#include "test_subs.h"
#include "test_e2sm.h"

using namespace std;


int main(int argc, char* argv[])
{
	char *aux;
	aux=getenv("RMR_SEED_RT");
	if (aux==NULL || *aux == '\0'){

		char rmr_seed[80]="RMR_SEED_RT=../init/routes.txt";
		putenv(rmr_seed);
	}
	//get configuration
	XappSettings config;
	//change the priority depending upon application requirement
	config.loadDefaultSettings();
	config.loadEnvVarSettings();

	//initialize rmr
	std::unique_ptr<XappRmr> rmr = std::make_unique<XappRmr>("38000");
	rmr->xapp_rmr_init(true);

	//create a dummy xapp
	std::unique_ptr<Xapp> dm_xapp = std::make_unique<Xapp>(std::ref(config),std::ref(*rmr));
	dm_xapp->Run();

	testing::InitGoogleTest(&argc, argv);
	int res = RUN_ALL_TESTS();



	return res;
}
