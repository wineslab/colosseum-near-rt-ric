

#ifndef AGENT_CONNECTOR_HPP_
#define AGENT_CONNECTOR_HPP_

#include <iostream>
#include <string.h>
#include <netinet/tcp.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <unistd.h>
#include <map>
#include <vector>

#define AGENT_0 "127.0.0.1"

// vector of agent IPs
extern std::vector<std::string> drl_agent_ip;

// key: DRL agent IP, value: socket file descriptor
extern std::map<std::string, int> agentIp_socket;

// key: DRL agent IP, value: gNB id
extern std::map<std::string, std::string> agentIp_gnbId;

int open_control_socket_agent(const char* dest_ip, const int dest_port);
void close_control_socket_agent(void);
std::string find_agent_ip_from_gnb(unsigned char* gnb_id);
int send_socket(char* buf, std::string dest_ip);

#endif
