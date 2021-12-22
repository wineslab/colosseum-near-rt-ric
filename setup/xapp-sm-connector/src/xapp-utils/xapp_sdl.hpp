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
 * xapp_sdl.hpp
 *
 *  Created on: Mar, 2020
 *  Author: Shraboni Jana
 */
#pragma once

#ifndef SRC_XAPP_UTILS_XAPP_SDL_HPP_
#define SRC_XAPP_UTILS_XAPP_SDL_HPP_

#include <iostream>
#include <string>
#include <memory>
#include <vector>
#include <map>
#include <set>
#include <sdl/syncstorage.hpp>
#include <mdclog/mdclog.h>

using namespace std;
using Namespace = std::string;
using Key = std::string;
using Data = std::vector<uint8_t>;
using DataMap = std::map<Key, Data>;
using Keys = std::set<Key>;

class XappSDL{
private:
	std::string sdl_namespace;

public:
	XappSDL(std::string ns) { sdl_namespace=ns; }
	void get_data(shareddatalayer::SyncStorage *);
	bool set_data(shareddatalayer::SyncStorage *);
};

#endif /* SRC_XAPP_UTILS_XAPP_SDL_HPP_ */
