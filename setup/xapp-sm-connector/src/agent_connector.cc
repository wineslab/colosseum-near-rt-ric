#include "agent_connector.hpp"


// open client of control socket with agent
int open_control_socket_agent(const char* dest_ip, const int dest_port) {

  std::cout << "Opening control socket with host " << dest_ip << ":" << dest_port << std::endl;

  int control_sckfd = socket(AF_INET, SOCK_STREAM, 0);
  if (control_sckfd < 0) {
	  std::cout << "ERROR: OPEN SOCKET" << std::endl;
      close(control_sckfd);
      return -1;
  }

  // SET SOCKET OPTIONS TO RELEASE THE SOCKET ADDRESS IMMEDIATELY AFTER
  // THE SOCKET IS CLOSED
  int option(1);
  setsockopt(control_sckfd, SOL_SOCKET, SO_REUSEADDR, (char*)&option, sizeof(option));

  struct sockaddr_in dest_addr = {0};
  dest_addr.sin_family = AF_INET;
  dest_addr.sin_port = htons(dest_port);

  // convert dest_ip from char* to network address
  if (inet_pton(AF_INET, dest_ip, &dest_addr.sin_addr) <= 0) {
      std::cout << "ERROR CONVERTING IP TO INTERNET ADDR" << std::endl;
      close(control_sckfd); // if conversion fail, close the socket and return error -2
      return -2;
  }

  if (connect(control_sckfd, (struct sockaddr *) &dest_addr, sizeof(dest_addr)) < 0) {
      std::cout << "ERROR: CONNECT" << std::endl;
      close(control_sckfd);
      return -3;
  }

  // update map
  std::string agent_ip;
  agent_ip.assign(dest_ip);
  // std::cout << "Agent IP " << agent_ip << std::endl;
  agentIp_socket[agent_ip] = control_sckfd;

  return 0;
}


// close control sockets
void close_control_socket_agent(void) {

  std::cout << "Closing control sockets with agent(s)" << std::endl;

  std::map<std::string, int>::iterator it;
  for (it = agentIp_socket.begin(); it != agentIp_socket.end(); ++it) {
    std::string agent_ip = it->first;
    int control_sckfd = it->second;

    close(control_sckfd);
  }

  // clear maps
  std::cout << "Clearing maps" << std::endl;
  agentIp_socket.clear();
  agentIp_gnbId.clear();
}


// find agent IP of socket to use with gNB id
std::string find_agent_ip_from_gnb(unsigned char* gnb_id_trans) {

  std::map<std::string, int>::iterator it_sck;
  std::map<std::string, std::string>::iterator it_gnb;
  std::string agent_ip;

  // convert transaction identifier (unsigned char*) to string
  std::string gnb_id(reinterpret_cast<char*>(gnb_id_trans));

  // check if gnb_id is already in agentIp_gnbId map
  bool found = false;
  for (it_gnb = agentIp_gnbId.begin(); it_gnb != agentIp_gnbId.end(); ++it_gnb) {
    if (gnb_id.compare(it_gnb->second) == 0) {
      agent_ip = it_gnb->first;
      found = true;
      break;
    }
  }

  if (!found) {
    // check if agent_ip is already in agentIp_gnbId map
    for (it_sck = agentIp_socket.begin(); it_sck != agentIp_socket.end(); ++it_sck) {
      agent_ip = it_sck->first;

      it_gnb = agentIp_gnbId.find(agent_ip);
      if (it_gnb == agentIp_gnbId.end()) {
        // insert into agentIp_gnbId map
        agentIp_gnbId[agent_ip] = gnb_id;
        break;
      }
    }
  }

  return agent_ip;
}


// send through socket
int send_socket(char* buf, size_t payload_size, std::string dest_ip) {

  int control_sckfd = -1;

  // get socket file descriptor
  std::map<std::string, int>::iterator it;
  for (it = agentIp_socket.begin(); it != agentIp_socket.end(); ++it) {
    std::string agent_ip = it->first;

    if (dest_ip.compare(agent_ip) == 0) {
      control_sckfd = it->second;
      break;
    }
  }

  if (control_sckfd == -1) {
    std::cout << "ERROR: Could not find socket for destination " << dest_ip << std::endl;
    return -1;
  }

  // const size_t max_size = 512;
  // char buf[max_size] = "Hello, Server!";  // store the data in a buffer
  int sent_size = send(control_sckfd, buf, payload_size, 0);

  if(sent_size < 0) { // the send returns a size of -1 in case of errors
      std::cout <<  "ERROR: SEND to agent " << dest_ip << std::endl;
      return -2;
  }
  else {
      std::cout << "Message sent" << std::endl;
  }

  return 0;
}
