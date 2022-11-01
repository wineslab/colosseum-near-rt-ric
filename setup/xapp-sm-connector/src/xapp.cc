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
 * xapp.cc
 *
 *  Mar, 2020 (Shraboni Jana)
 */

#include "xapp.hpp"

#define BUFFER_SIZE 1024

std::map<string, int> agentIp_socket;
std::map<std::string, std::string> agentIp_gnbId;
std::vector<std::string> drl_agent_ip{AGENT_0};

 Xapp::Xapp(XappSettings &config, XappRmr &rmr){

	rmr_ref = &rmr;
	config_ref = &config;
	xapp_mutex = NULL;
	subhandler_ref = NULL;
	return;
}

Xapp::~Xapp(void){

    //Joining the threads
	int threadcnt = xapp_rcv_thread.size();
	for(int i=0; i<threadcnt; i++){
		if(xapp_rcv_thread[i].joinable())
			xapp_rcv_thread[i].join();
	}
	xapp_rcv_thread.clear();

	if(xapp_mutex!=NULL){
		xapp_mutex->~mutex();
		delete xapp_mutex;
	}

    // close sockets
    close_control_socket_agent();

    // join threads
    for (int i = 0; i < control_thr_rx.size(); ++i) {
        if(control_thr_rx[i] && control_thr_rx[i]->joinable()) {
            control_thr_rx[i]->join();
        }
    }

    if(ext_control_thr_rx && ext_control_thr_rx->joinable()) {
        ext_control_thr_rx->join();
    }
};

//Stop the xapp. Note- To be run only from unit test scripts.
void Xapp::stop(void){
    // Get the mutex lock
	std::lock_guard<std::mutex> guard(*xapp_mutex);
	rmr_ref->set_listen(false);
	rmr_ref->~XappRmr();

	//Detaching the threads....not sure if this is the right way to stop the receiver threads.
	//Hence function should be called only in Unit Tests
	int threadcnt = xapp_rcv_thread.size();
	for(int i=0; i<threadcnt; i++){
		xapp_rcv_thread[i].detach();
	}
	sleep(10);
}

void Xapp::startup(SubscriptionHandler &sub_ref) {

	subhandler_ref = &sub_ref;

    if (GNB_ID == "") {
        // get list of gnbs from ric
        std::cout << "Getting gNB list from RIC" << std::endl;
        set_rnib_gnblist();
    } else if (strcmp(GNB_ID, "file") == 0) {
        // Get gNB list from file
        std::cout << "Getting gNB list from file" << std::endl;
        std::ifstream id_file("gnb_list.txt", ios_base::in);
        if(!id_file) {
            std::cerr << "Error in opening file, gonna crash!" << std::endl;
            exit(-1);
        }
        std::string id_gnb;
        while(getline(id_file,id_gnb)){
            std::cout << "gNB read: " << id_gnb << std::endl;
            rnib_gnblist.push_back(id_gnb);
        }
    } else if (strcmp(GNB_ID, "ns-o-ran") == 0) {
        std::vector <std::string> gnb_ids{"gnb:131-133-31000000", "gnb:131-133-32000000", "gnb:131-133-33000000",
                                          "gnb:131-133-34000000", "gnb:131-133-35000000"};
        for (vector<string>::iterator it = gnb_ids.begin(); it != gnb_ids.end(); it++) {
            std::cout << "gNB read: " << *it << std::endl;
            rnib_gnblist.push_back(*it);
        }
    } else {
        // only insert target gnb
        std::cout << "Querying target gNB" << std::endl;
        rnib_gnblist.push_back(GNB_ID);
    }

    // open external control socket in thread and wait for message
    ext_control_thr_rx = std::unique_ptr<std::thread>(new std::thread{&Xapp::handle_external_control_message, this, SOCKET_PORT_EXT});

    for (int i = 0; i < drl_agent_ip.size(); ++i) {
        // open control socket with agent
        if (open_control_socket_agent(const_cast<char*>(drl_agent_ip[i].c_str()), 4200) == 0) {
            // start receive thread
            std::unique_ptr<std::thread> tmp_thr = std::unique_ptr<std::thread>(new std::thread{&Xapp::handle_rx_msg_agent, this, drl_agent_ip[i]});
            control_thr_rx.push_back(std::move(tmp_thr));
        }
    }

    // send test message
    // send_socket("Hello, Server!", AGENT_1);

	//send subscriptions.
	startup_subscribe_requests();

	//read A1 policies
	// startup_get_policies();

	// send control
	// send_ric_control_request();

	return;
}

void Xapp::Run(){
	rmr_ref->set_listen(true);
	if(xapp_mutex == NULL){
		xapp_mutex = new std::mutex();
	}
	std::lock_guard<std::mutex> guard(*xapp_mutex);

	for(size_t j=0; j < _callbacks.size(); j++){
        std::thread th_recv([&](){ rmr_ref->xapp_rmr_receive(std::move(_callbacks[j]), rmr_ref);});
		xapp_rcv_thread.push_back(std::move(th_recv));
	}

	return;
}

//Starting a seperate single receiver
void Xapp::start_xapp_receiver(XappMsgHandler& mp_handler){
	//start a receiver thread. Can be multiple receiver threads for more than 1 listening port.
	rmr_ref->set_listen(true);
	if(xapp_mutex == NULL){
		xapp_mutex = new std::mutex();
	}

	mdclog_write(MDCLOG_INFO,"Receiver Thread file= %s, line=%d",__FILE__,__LINE__);
	std::lock_guard<std::mutex> guard(*xapp_mutex);
	std::thread th_recv([&](){ rmr_ref->xapp_rmr_receive(std::move(mp_handler), rmr_ref);});
	xapp_rcv_thread.push_back(std::move(th_recv));
	return;
}

void Xapp::shutdown(){
	return;
}

// handle received message from DRL agent
void Xapp::handle_rx_msg(void) {
    std::cout << "handle_rx_msg" << std::endl;

    const size_t max_size = 256;
    char buf[max_size] = {0};

    // listen to control from agent
    while (true) {
        // iterate through map
        std::map<std::string, int>::iterator it;
        for (it = agentIp_socket.begin(); it != agentIp_socket.end(); ++it) {
            std::string agent_ip = it->first;
            int control_sckfd = it->second;

            int rcv_size = recv(control_sckfd, buf, max_size, 0);
            if (rcv_size > 0) {
                std::cout << "Message from agent " << agent_ip << std::endl;
                std::cout << buf << std::endl;

                // get gnb_id from agent IP
                std::map<std::string, std::string>::iterator it_gnb;
                it_gnb = agentIp_gnbId.find(agent_ip);

                // send RIC control request
                if (it_gnb != agentIp_gnbId.end()) {
                    send_ric_control_request(buf, it_gnb->second);
                }
                else {
                    std::cout << "ERROR: No gNB ID found for agent " << agent_ip << std::endl;
                }

                memset(buf, 0, max_size);
            }
        }
    }
}

// handle received message from a specific DRL agent
void Xapp::handle_rx_msg_agent(std::string agent_ip) {
    std::cout << "Opening RX thread with agent " << agent_ip << std::endl;

    const size_t max_size = 256;
    char buf[max_size] = {0};

    // listen to control from agent
    while (true) {
        // get control_sckfd from agent IP
        std::map<std::string, int>::iterator it;
        it = agentIp_socket.find(agent_ip);

        if (it != agentIp_socket.end()) {
            int control_sckfd = it->second;

            int rcv_size = recv(control_sckfd, buf, max_size, 0);
            if (rcv_size > 0) {
                std::cout << "Message from agent " << agent_ip << std::endl;
                std::cout << buf << std::endl;

                // get gnb_id from agent IP
                std::map<std::string, std::string>::iterator it_gnb;
                it_gnb = agentIp_gnbId.find(agent_ip);

                // send RIC control request
                if (it_gnb != agentIp_gnbId.end()) {
                send_ric_control_request(buf, it_gnb->second);
                }
                else {
                std::cout << "ERROR: No gNB ID found for agent " << agent_ip << std::endl;
                }

                memset(buf, 0, max_size);
            }
        }
    }
}

// handle external control socket message
void Xapp::handle_external_control_message(int port) {

    // Create a socket (IPv4, TCP)
    int sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd == -1) {
        std::cout << "Failed to create socket. errno: " << errno << std::endl;
        return;
    }

    // Listen on given port on any address
    sockaddr_in sockaddr;
    sockaddr.sin_family = AF_INET;
    sockaddr.sin_addr.s_addr = INADDR_ANY;
    sockaddr.sin_port = htons(port);

    if (bind(sockfd, (struct sockaddr*)&sockaddr, sizeof(sockaddr)) < 0) {
        std::cout << "Failed to bind to port. Errno: " << errno << std::endl;
        return;
    }

    // Start listening. Hold at most 10 connections in the queue
    if (listen(sockfd, 10) < 0) {
        std::cout << "Failed to listen on socket. Errno: " << errno << std::endl;
        return;
    }

    std::cout << "Opened control socket server on port " << port << std::endl;

    while (true) {
        auto addrlen = sizeof(sockaddr);
        int connection = accept(sockfd, (struct sockaddr*)&sockaddr, (socklen_t*)&addrlen);

        if (connection < 0) {
            continue;
        }

        // Read from the connection
        const size_t max_size = 256;
        char buffer[max_size] = {0};
        auto bytes_read = read(connection, buffer, 100);

        if (bytes_read > 0) {
            std::cout << "External control socket. Message received: " << buffer << std::endl;

            // TODO: check if message is termination, send to DU and shutdown xApp
            if (strcmp(buffer, XAPP_TERMINATE) == 0) {
            	terminate_du_reporting();
            }
        }

        memset(buffer, 0, max_size);
        close(connection);
  }

  close(sockfd);

  return;
}


// terminate all DU reportings and shutdown xApp
void Xapp::terminate_du_reporting(void) {

	std::map<std::string, int>::iterator it;
    for (it = agentIp_socket.begin(); it != agentIp_socket.end(); ++it) {
        std::string agent_ip = it->first;
        int control_sckfd = it->second;

        // get gnb_id from agent IP
        std::map<std::string, std::string>::iterator it_gnb;
        it_gnb = agentIp_gnbId.find(agent_ip);

        std::cout << "Terminating reporting gNB " << it_gnb->second << std::endl;
        send_ric_control_request(XAPP_TERMINATE, it_gnb->second);
    }

    // stop xapp docker container with SIGTERM (15)
    if (DEBUG) {
        std::cout << "Debug mode, only echoing" << std::endl;
        system("echo kill -s 15 1");
    }
    else {
        system("kill -s 15 1");
    }
}

void Xapp::send_ric_control_request(char* payload, std::string gnb_id) {

    std::cout << "Sending RIC Control Request" << std::endl;

	bool res;
	size_t data_size = ASN_BUFF_MAX_SIZE;
	unsigned char	data[data_size];
	unsigned char meid[RMR_MAX_MEID];
	std::string xapp_id = config_ref->operator [](XappSettings::SettingName::XAPP_ID);

	mdclog_write(MDCLOG_INFO, "Preparing to send control in file= %s, line=%d", __FILE__, __LINE__);

    auto gnblist = get_rnib_gnblist();
    int sz = gnblist.size();

    if(sz <= 0) {
	   mdclog_write(MDCLOG_INFO, "Subscriptions cannot be sent as GNBList in RNIB is NULL");
        return;
    }

	// give the message to subscription handler, along with the transmitter.
 	strcpy((char*)meid, gnb_id.c_str());
	std::cout << "RIC Control Request, gNB " << gnb_id << std::endl;

 	// helpers
 	// set fields randomly
 	ric_control_helper din {};
       	//= {
 	//	1,
 	//	1,
 	//	0,
 	//	1,
 	//	-1,
 	//	0,
 	//	0,
 	//	1,
 	//	"test", // control_msg
 	//	5, // control_msg_size
 	//	"testh", // control header
 	//	6,
 	//	"testp", // call process id
 	//	6
 	//};
	const char* msg = payload;
	din.control_msg_size = strlen(msg) + 1;
	mdclog_write(MDCLOG_INFO, "Size of msg %d", din.control_msg_size);
	din.control_msg = (uint8_t*) calloc(din.control_msg_size, sizeof(uint8_t));
	std::memcpy(din.control_msg, msg, din.control_msg_size);
 	ric_control_helper dout {};

 	// control request object
 	ric_control_request ctrl_req {};
 	ric_control_request ctrl_recv {};

 	unsigned char buf[BUFFER_SIZE];
    size_t buf_size = BUFFER_SIZE;

 	res = ctrl_req.encode_e2ap_control_request(&buf[0], &buf_size, din);

 	xapp_rmr_header rmr_header;
	rmr_header.message_type = RIC_CONTROL_REQ;
	rmr_header.payload_length = buf_size; //data_size
	strcpy((char*)rmr_header.meid, gnb_id.c_str());

	mdclog_write(MDCLOG_INFO, "Sending CTRL REQ in file= %s, line=%d for MEID %s", __FILE__, __LINE__, meid);

    int result = rmr_ref->xapp_rmr_send(&rmr_header, (void*)buf);
    if(result) {
   	    mdclog_write(MDCLOG_INFO, "CTRL REQ SUCCESSFUL in file= %s, line=%d for MEID %s",__FILE__,__LINE__, meid);
    }
}

void Xapp::startup_subscribe_requests(void ){

    bool res;
    size_t data_size = ASN_BUFF_MAX_SIZE;
    unsigned char	data[data_size];
    unsigned char meid[RMR_MAX_MEID];
    std::string xapp_id = config_ref->operator [](XappSettings::SettingName::XAPP_ID);

    mdclog_write(MDCLOG_INFO,"Preparing to send subscription in file= %s, line=%d", __FILE__, __LINE__);

    auto gnblist = get_rnib_gnblist();

    int sz = gnblist.size();

    if(sz <= 0)
       mdclog_write(MDCLOG_INFO,"Subscriptions cannot be sent as GNBList in RNIB is NULL");

    for(int i = 0; i<sz; i++){
        std::cout << "Sending subscriptions to: " << gnblist[i] << std::endl;

        // give the message to subscription handler, along with the transmitter.
        strcpy((char*)meid,gnblist[i].c_str());

        // char *strMsg = "Subscription Request from HelloWorld XApp\0";
        // strncpy((char *)data,strMsg,strlen(strMsg));
        // data_size = strlen(strMsg);

        subscription_helper  din;
        subscription_helper  dout;

        subscription_request sub_req;
        subscription_request sub_recv;

        unsigned char buf[BUFFER_SIZE];
        size_t buf_size = BUFFER_SIZE;
        bool res;

        //Random Data  for request
        int request_id = XAPP_REQ_ID;
        int function_id = 200;

        // DU report timer in ms
        std::string event_def = "250";

        din.set_request(request_id);
        din.set_function_id(function_id);
        din.set_event_def(event_def.c_str(), event_def.length());

        std::string act_def = "HelloWorld Action Definition";

        din.add_action(1,1,(void*)act_def.c_str(), act_def.length(), 0);

        res = sub_req.encode_e2ap_subscription(&buf[0], &buf_size, din);

        xapp_rmr_header rmr_header;
        rmr_header.message_type = RIC_SUB_REQ;
        rmr_header.payload_length = buf_size; //data_size
        strcpy((char*)rmr_header.meid,gnblist[i].c_str());

        mdclog_write(MDCLOG_INFO,"Sending subscription in file= %s, line=%d for MEID %s",__FILE__,__LINE__, meid);
        auto transmitter = std::bind(&XappRmr::xapp_rmr_send,rmr_ref, &rmr_header, (void*)buf );//(void*)data);

        int result = subhandler_ref->manage_subscription_request(meid, transmitter);
        if(result){
            mdclog_write(MDCLOG_INFO,"Subscription SUCCESSFUL in file= %s, line=%d for MEID %s",__FILE__,__LINE__, meid);
        }
    }
}

void Xapp::startup_get_policies(void){

    int policy_id = HELLOWORLD_POLICY_ID;

    std::string policy_query = "{\"policy_type_id\":" + std::to_string(policy_id) + "}";
    unsigned char * message = (unsigned char *)calloc(policy_query.length(), sizeof(unsigned char));
    memcpy(message, policy_query.c_str(),  policy_query.length());
    xapp_rmr_header header;
    header.payload_length = policy_query.length();
    header.message_type = A1_POLICY_QUERY;
    mdclog_write(MDCLOG_INFO, "Sending request for policy id %d\n", policy_id);
    rmr_ref->xapp_rmr_send(&header, (void *)message);
    free(message);

}

void Xapp::set_rnib_gnblist(void) {

    openSdl();

    void *result = getListGnbIds();
    if(strlen((char*)result) < 1){
        mdclog_write(MDCLOG_ERR, "ERROR: no data from getListGnbIds\n");
        return;
    }

    mdclog_write(MDCLOG_INFO, "GNB List in R-NIB %s\n", (char*)result);

    // remove non-unicode characters that make rapodjson fail the parsing
    std::string result_clean((char*) result);
    while (result_clean.back() != '}') {
        result_clean.pop_back();
    }

    Document doc;
    ParseResult parseJson = doc.Parse(result_clean.c_str());
    if (!parseJson) {
    	//std::cerr << "JSON parse error: %s (%u)\n", GetParseErrorFunc(parseJson.Code());
        std::cerr << "JSON parse error: " << GetParseErrorFunc(parseJson.Code()) << std::endl;
    	return;
    }

    if(!doc.HasMember("gnb_list")){
    	mdclog_write(MDCLOG_INFO, "JSON Has No GNB List Object");
    	return;
    }
    assert(doc.HasMember("gnb_list"));

    const Value& gnblist = doc["gnb_list"];
    if (gnblist.IsNull())
        return;

    if(!gnblist.IsArray()){
    	mdclog_write(MDCLOG_INFO, "GNB List is not an array");
    	return;
    }


    assert(gnblist.IsArray());
    for (SizeType i = 0; i < gnblist.Size(); i++) { // Uses SizeType instead of size_t
    	assert(gnblist[i].IsObject());
    	const Value& gnbobj = gnblist[i];
    	assert(gnbobj.HasMember("inventory_name"));
    	assert(gnbobj["inventory_name"].IsString());
    	std::string name = gnbobj["inventory_name"].GetString();
    	rnib_gnblist.push_back(name);
    }

    closeSdl();
    return;
}
