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
 *//*
 * xapp.hpp
 *
 *  Mar, 2020 (Shraboni Jana)
 *
 */

#pragma once

#ifndef SRC_XAPP_HPP_
#define SRC_XAPP_HPP_

#include <iostream>
#include <string>
#include <memory>
#include <csignal>
#include <stdio.h>
#include <pthread.h>
#include <unordered_map>
#include "xapp_rmr.hpp"
#include "xapp_sdl.hpp"
#include "rapidjson/writer.h"
#include "rapidjson/document.h"
#include "rapidjson/error/error.h"
#include <thread>
#include <netinet/tcp.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include "agent_connector.hpp"
#include <vector>

#include "msgs_proc.hpp"
#include "subs_mgmt.hpp"
#include "xapp_config.hpp"
extern "C" {
#include "rnib/rnibreader.h"
}
using namespace std;
using namespace std::placeholders;
using namespace rapidjson;

// used to identify xApp requests with DUs
// Note: this value is updated at Docker build to use the last octet of the xApp IP
#define XAPP_REQ_ID 0

#define SOCKET_PORT_EXT 7000
#define XAPP_TERMINATE "terminate"

// id of gnb to control as seen from the ric
#define GNB_ID ""

#define DEBUG 0

class Xapp{
public:

  Xapp(XappSettings &, XappRmr &);

  ~Xapp(void);

  void stop(void);

  void startup(SubscriptionHandler &);
  void shutdown(void);

  void start_xapp_receiver(XappMsgHandler &);
  void Run();

  void sdl_data(void);

  Xapp(Xapp const &)=delete;
  Xapp& operator=(Xapp const &) = delete;

  void register_handler(XappMsgHandler &fn){
    _callbacks.emplace_back(fn);
  }

  //getters/setters.
  void set_rnib_gnblist(void);
  std::vector<std::string> get_rnib_gnblist(){ return rnib_gnblist; }

private:
  void startup_subscribe_requests(void );
  void shutdown_subscribe_deletes(void);
  void send_ric_control_request(char* payload, std::string gnb_id);
  void startup_get_policies(void );

  void handle_rx_msg(void);
  void handle_rx_msg_agent(std::string agent_ip);
  void handle_external_control_message(int port);
  void terminate_du_reporting(void);

  XappRmr * rmr_ref;
  XappSettings * config_ref;
  SubscriptionHandler *subhandler_ref;

  std::mutex *xapp_mutex;
  std::vector<std::thread> xapp_rcv_thread;
  std::vector<std::string> rnib_gnblist;
  std::vector<XappMsgHandler> _callbacks;

  std::vector<std::unique_ptr<std::thread>> control_thr_rx;
  std::unique_ptr<std::thread> ext_control_thr_rx;
};


#endif /* SRC_XAPP_HPP_ */
