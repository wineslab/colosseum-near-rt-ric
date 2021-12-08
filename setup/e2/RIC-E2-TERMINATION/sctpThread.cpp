// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

// TODO: High-level file comment.



#include <3rdparty/oranE2/RANfunctions-List.h>
#include "sctpThread.h"
#include "BuildRunName.h"

#include "3rdparty/oranE2SM/E2SM-gNB-NRT-RANfunction-Definition.h"
#include "BuildXml.h"
#include "pugixml/src/pugixml.hpp"

// #include <string.h>
// #include <netinet/tcp.h>
// #include <arpa/inet.h>
// #include <sys/socket.h>
// #include <netinet/in.h>
// #include <unistd.h>
// #include <iostream>

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <libgen.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <netinet/sctp.h>
#include <arpa/inet.h>

#include <rapidjson/document.h>
#include <rapidjson/writer.h>
#include <rapidjson/stringbuffer.h>
#include <iostream>

using namespace std;
//using namespace std::placeholders;
using namespace boost::filesystem;

using namespace rapidjson;

//#ifdef __cplusplus
//extern "C"
//{
//#endif

int sendMessageSocket(const int dest_port);
int parseMessageData(ReportingMessages_t &message, char *data, char *host, uint16_t &port);

// need to expose without the include of gcov
extern "C" void __gcov_flush(void);

static void catch_function(int signal) {
    __gcov_flush();
    exit(signal);
}


BOOST_LOG_INLINE_GLOBAL_LOGGER_DEFAULT(my_logger, src::logger_mt)

boost::shared_ptr<sinks::synchronous_sink<sinks::text_file_backend>> boostLogger;
double cpuClock = 0.0;
bool jsonTrace = true;

void init_log() {
    mdclog_attr_t *attr;
    mdclog_attr_init(&attr);
    mdclog_attr_set_ident(attr, "E2Terminator");
    mdclog_init(attr);
    mdclog_attr_destroy(attr);
}
auto start_time = std::chrono::high_resolution_clock::now();
typedef std::chrono::duration<double, std::ratio<1,1>> seconds_t;

double age() {
    return seconds_t(std::chrono::high_resolution_clock::now() - start_time).count();
}

double approx_CPU_MHz(unsigned sleeptime) {
    using namespace std::chrono_literals;
    uint32_t aux = 0;
    uint64_t cycles_start = rdtscp(aux);
    double time_start = age();
    std::this_thread::sleep_for(sleeptime * 1ms);
    uint64_t elapsed_cycles = rdtscp(aux) - cycles_start;
    double elapsed_time = age() - time_start;
    return elapsed_cycles / elapsed_time;
}

//std::atomic<int64_t> rmrCounter{0};
std::atomic<int64_t> num_of_messages{0};
std::atomic<int64_t> num_of_XAPP_messages{0};
static long transactionCounter = 0;

int buildListeningPort(sctp_params_t &sctpParams) {
    sctpParams.listenFD = socket (AF_INET6, SOCK_STREAM, IPPROTO_SCTP);
    struct sockaddr_in6 servaddr {};
    servaddr.sin6_family = AF_INET6;
    servaddr.sin6_addr   = in6addr_any;
    servaddr.sin6_port = htons(sctpParams.sctpPort);
    if (bind(sctpParams.listenFD, (SA *)&servaddr, sizeof(servaddr)) < 0 ) {
        mdclog_write(MDCLOG_ERR, "Error binding. %s\n", strerror(errno));
        return -1;
    }
    if (setSocketNoBlocking(sctpParams.listenFD) == -1) {
        //mdclog_write(MDCLOG_ERR, "Error binding. %s", strerror(errno));
        return -1;
    }
    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        struct sockaddr_in6 cliaddr {};
        socklen_t len = sizeof(cliaddr);
        getsockname(sctpParams.listenFD, (SA *)&cliaddr, &len);
        char buff[1024] {};
        inet_ntop(AF_INET6, &cliaddr.sin6_addr, buff, sizeof(buff));
        mdclog_write(MDCLOG_DEBUG, "My address: %s, port %d\n", buff, htons(cliaddr.sin6_port));
    }

    if (listen(sctpParams.listenFD, SOMAXCONN) < 0) {
        mdclog_write(MDCLOG_ERR, "Error listening. %s\n", strerror(errno));
        return -1;
    }
    struct epoll_event event {};
    event.events = EPOLLIN | EPOLLET;
    event.data.fd = sctpParams.listenFD;

    // add listening port to epoll
    if (epoll_ctl(sctpParams.epoll_fd, EPOLL_CTL_ADD, sctpParams.listenFD, &event)) {
        printf("Failed to add descriptor to epoll\n");
        mdclog_write(MDCLOG_ERR, "Failed to add descriptor to epoll. %s\n", strerror(errno));
        return -1;
    }

    return 0;
}

int buildConfiguration(sctp_params_t &sctpParams) {
    path p = (sctpParams.configFilePath + "/" + sctpParams.configFileName).c_str();
    if (exists(p)) {
        const int size = 2048;
        auto fileSize = file_size(p);
        if (fileSize > size) {
            mdclog_write(MDCLOG_ERR, "File %s larger than %d", p.string().c_str(), size);
            return -1;
        }
    } else {
        mdclog_write(MDCLOG_ERR, "Configuration File %s not exists", p.string().c_str());
        return -1;
    }

    ReadConfigFile conf;
    if (conf.openConfigFile(p.string()) == -1) {
        mdclog_write(MDCLOG_ERR, "Filed to open config file %s, %s",
                     p.string().c_str(), strerror(errno));
        return -1;
    }
    int rmrPort = conf.getIntValue("nano");
    if (rmrPort == -1) {
        mdclog_write(MDCLOG_ERR, "illigal RMR port ");
        return -1;
    }
    sctpParams.rmrPort = (uint16_t)rmrPort;
    snprintf(sctpParams.rmrAddress, sizeof(sctpParams.rmrAddress), "%d", (int) (sctpParams.rmrPort));

    auto tmpStr = conf.getStringValue("loglevel");
    if (tmpStr.length() == 0) {
        mdclog_write(MDCLOG_ERR, "illigal loglevel. Set loglevel to MDCLOG_INFO");
        tmpStr = "info";
    }
    transform(tmpStr.begin(), tmpStr.end(), tmpStr.begin(), ::tolower);

    if ((tmpStr.compare("debug")) == 0) {
        sctpParams.logLevel = MDCLOG_DEBUG;
    } else if ((tmpStr.compare("info")) == 0) {
        sctpParams.logLevel = MDCLOG_INFO;
    } else if ((tmpStr.compare("warning")) == 0) {
        sctpParams.logLevel = MDCLOG_WARN;
    } else if ((tmpStr.compare("error")) == 0) {
        sctpParams.logLevel = MDCLOG_ERR;
    } else {
        mdclog_write(MDCLOG_ERR, "illigal loglevel = %s. Set loglevel to MDCLOG_INFO", tmpStr.c_str());
        sctpParams.logLevel = MDCLOG_INFO;
    }
    mdclog_level_set(sctpParams.logLevel);

    tmpStr = conf.getStringValue("volume");
    if (tmpStr.length() == 0) {
        mdclog_write(MDCLOG_ERR, "illigal volume.");
        return -1;
    }

    char tmpLogFilespec[VOLUME_URL_SIZE];
    tmpLogFilespec[0] = 0;
    sctpParams.volume[0] = 0;
    snprintf(sctpParams.volume, VOLUME_URL_SIZE, "%s", tmpStr.c_str());
    // copy the name to temp file as well
    snprintf(tmpLogFilespec, VOLUME_URL_SIZE, "%s", tmpStr.c_str());


    // define the file name in the tmp directory under the volume
    strcat(tmpLogFilespec,"/tmp/E2Term_%Y-%m-%d_%H-%M-%S.%N.tmpStr");

    sctpParams.myIP = conf.getStringValue("local-ip");
    if (sctpParams.myIP.length() == 0) {
        mdclog_write(MDCLOG_ERR, "illigal local-ip.");
        return -1;
    }

    int sctpPort = conf.getIntValue("sctp-port");
    if (sctpPort == -1) {
        mdclog_write(MDCLOG_ERR, "illigal SCTP port ");
        return -1;
    }
    sctpParams.sctpPort = (uint16_t)sctpPort;

    sctpParams.fqdn = conf.getStringValue("external-fqdn");
    if (sctpParams.fqdn.length() == 0) {
        mdclog_write(MDCLOG_ERR, "illigal external-fqdn");
        return -1;
    }

    std::string pod = conf.getStringValue("pod_name");
    if (pod.length() == 0) {
        mdclog_write(MDCLOG_ERR, "illigal pod_name in config file");
        return -1;
    }
    auto *podName = getenv(pod.c_str());
    if (podName == nullptr) {
        mdclog_write(MDCLOG_ERR, "illigal pod_name or environment varible not exists : %s", pod.c_str());
        return -1;

    } else {
        sctpParams.podName.assign(podName);
        if (sctpParams.podName.length() == 0) {
            mdclog_write(MDCLOG_ERR, "illigal pod_name");
            return -1;
        }
    }

    tmpStr = conf.getStringValue("trace");
    transform(tmpStr.begin(), tmpStr.end(), tmpStr.begin(), ::tolower);
    if ((tmpStr.compare("start")) == 0) {
        mdclog_write(MDCLOG_INFO, "Trace set to: start");
        sctpParams.trace = true;
    } else if ((tmpStr.compare("stop")) == 0) {
        mdclog_write(MDCLOG_INFO, "Trace set to: stop");
        sctpParams.trace = false;
    }
    jsonTrace = sctpParams.trace;

    sctpParams.ka_message_length = snprintf(sctpParams.ka_message, KA_MESSAGE_SIZE, "{\"address\": \"%s:%d\","
                                                                                    "\"fqdn\": \"%s\","
                                                                                    "\"pod_name\": \"%s\"}",
                                            (const char *)sctpParams.myIP.c_str(),
                                            sctpParams.rmrPort,
                                            sctpParams.fqdn.c_str(),
                                            sctpParams.podName.c_str());

    if (mdclog_level_get() >= MDCLOG_INFO) {
        mdclog_mdc_add("RMR Port", to_string(sctpParams.rmrPort).c_str());
        mdclog_mdc_add("LogLevel", to_string(sctpParams.logLevel).c_str());
        mdclog_mdc_add("volume", sctpParams.volume);
        mdclog_mdc_add("tmpLogFilespec", tmpLogFilespec);
        mdclog_mdc_add("my ip", sctpParams.myIP.c_str());
        mdclog_mdc_add("pod name", sctpParams.podName.c_str());

        mdclog_write(MDCLOG_INFO, "running parameters for instance : %s", sctpParams.ka_message);
    }
    mdclog_mdc_clean();

    // Files written to the current working directory
    boostLogger = logging::add_file_log(
            keywords::file_name = tmpLogFilespec, // to temp directory
            keywords::rotation_size = 10 * 1024 * 1024,
            keywords::time_based_rotation = sinks::file::rotation_at_time_interval(posix_time::hours(1)),
            keywords::format = "%Message%"
            //keywords::format = "[%TimeStamp%]: %Message%" // use each tmpStr with time stamp
    );

    // Setup a destination folder for collecting rotated (closed) files --since the same volumn can use rename()
    boostLogger->locked_backend()->set_file_collector(sinks::file::make_collector(
            keywords::target = sctpParams.volume
    ));

    // Upon restart, scan the directory for files matching the file_name pattern
    boostLogger->locked_backend()->scan_for_files();

    // Enable auto-flushing after each tmpStr record written
    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        boostLogger->locked_backend()->auto_flush(true);
    }

    return 0;
}



int main(const int argc, const char **argv) {
    sctp_params_t sctpParams;

    {
        std::random_device device{};
        std::mt19937 generator(device());
        std::uniform_int_distribution<long> distribution(1, (long) 1e12);
        transactionCounter = distribution(generator);
    }

//    uint64_t st = 0;
//    uint32_t aux1 = 0;
//   st = rdtscp(aux1);

    unsigned num_cpus = std::thread::hardware_concurrency();
    init_log();
    mdclog_level_set(MDCLOG_INFO);

    if (std::signal(SIGINT, catch_function) == SIG_ERR) {
        mdclog_write(MDCLOG_ERR, "Error initializing SIGINT");
        exit(1);
    }
    if (std::signal(SIGABRT, catch_function)== SIG_ERR) {
        mdclog_write(MDCLOG_ERR, "Error initializing SIGABRT");
        exit(1);
    }
    if (std::signal(SIGTERM, catch_function)== SIG_ERR) {
        mdclog_write(MDCLOG_ERR, "Error initializing SIGTERM");
        exit(1);
    }

    cpuClock = approx_CPU_MHz(100);

    mdclog_write(MDCLOG_DEBUG, "CPU speed %11.11f", cpuClock);

    auto result = parse(argc, argv, sctpParams);

    if (buildConfiguration(sctpParams) != 0) {
        exit(-1);
    }

    // start epoll
    sctpParams.epoll_fd = epoll_create1(0);
    if (sctpParams.epoll_fd == -1) {
        mdclog_write(MDCLOG_ERR, "failed to open epoll descriptor");
        exit(-1);
    }

    getRmrContext(sctpParams);
    if (sctpParams.rmrCtx == nullptr) {
        close(sctpParams.epoll_fd);
        exit(-1);
    }

    if (buildInotify(sctpParams) == -1) {
        close(sctpParams.rmrListenFd);
        rmr_close(sctpParams.rmrCtx);
        close(sctpParams.epoll_fd);
        exit(-1);
    }

    if (buildListeningPort(sctpParams) != 0) {
        close(sctpParams.rmrListenFd);
        rmr_close(sctpParams.rmrCtx);
        close(sctpParams.epoll_fd);
        exit(-1);
    }

    sctpParams.sctpMap = new mapWrapper();

    std::vector<std::thread> threads(num_cpus);
//    std::vector<std::thread> threads;

    num_cpus = 1;
    for (unsigned int i = 0; i < num_cpus; i++) {
        threads[i] = std::thread(listener, &sctpParams);

        cpu_set_t cpuset;
        CPU_ZERO(&cpuset);
        CPU_SET(i, &cpuset);
        int rc = pthread_setaffinity_np(threads[i].native_handle(), sizeof(cpu_set_t), &cpuset);
        if (rc != 0) {
            mdclog_write(MDCLOG_ERR, "Error calling pthread_setaffinity_np: %d", rc);
        }
    }

    auto statFlag = false;
    auto statThread = std::thread(statColectorThread, (void *)&statFlag);

    //loop over term_init until first message from xApp
    handleTermInit(sctpParams);

    for (auto &t : threads) {
        t.join();
    }

    statFlag = true;
    statThread.join();

    return 0;
}

void handleTermInit(sctp_params_t &sctpParams) {
    sendTermInit(sctpParams);
    //send to e2 manager init of e2 term
    //E2_TERM_INIT

    int count = 0;
    while (true) {
        auto xappMessages = num_of_XAPP_messages.load(std::memory_order_acquire);
        if (xappMessages > 0) {
            if (mdclog_level_get() >=  MDCLOG_INFO) {
                mdclog_write(MDCLOG_INFO, "Got a message from some appliction, stop sending E2_TERM_INIT");
            }
            return;
        }
        usleep(100000);
        count++;
        if (count % 1000 == 0) {
            mdclog_write(MDCLOG_ERR, "GOT No messages from any xApp");
            sendTermInit(sctpParams);
        }
    }
}

void sendTermInit(sctp_params_t &sctpParams) {
    rmr_mbuf_t *msg = rmr_alloc_msg(sctpParams.rmrCtx, sctpParams.ka_message_length);
    auto count = 0;
    while (true) {
        msg->mtype = E2_TERM_INIT;
        msg->state = 0;
        rmr_bytes2payload(msg, (unsigned char *)sctpParams.ka_message, sctpParams.ka_message_length);
        static unsigned char tx[32];
        auto txLen = snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
        rmr_bytes2xact(msg, tx, txLen);
        msg = rmr_send_msg(sctpParams.rmrCtx, msg);
        if (msg == nullptr) {
            msg = rmr_alloc_msg(sctpParams.rmrCtx, sctpParams.ka_message_length);
        } else if (msg->state == 0) {
            rmr_free_msg(msg);
            if (mdclog_level_get() >=  MDCLOG_INFO) {
                mdclog_write(MDCLOG_INFO, "E2_TERM_INIT succsesfuly sent ");
            }
            return;
        } else {
            if (count % 100 == 0) {
                mdclog_write(MDCLOG_ERR, "Error sending E2_TERM_INIT cause : %s ", translateRmrErrorMessages(msg->state).c_str());
            }
            sleep(1);
        }
        count++;
    }
}

/**
 *
 * @param argc
 * @param argv
 * @param sctpParams
 * @return
 */
cxxopts::ParseResult parse(int argc, const char *argv[], sctp_params_t &sctpParams) {
    cxxopts::Options options(argv[0], "e2 term help");
    options.positional_help("[optional args]").show_positional_help();
    options.allow_unrecognised_options().add_options()
            ("p,path", "config file path", cxxopts::value<std::string>(sctpParams.configFilePath)->default_value("config"))
            ("f,file", "config file name", cxxopts::value<std::string>(sctpParams.configFileName)->default_value("config.conf"))
            ("h,help", "Print help");

    auto result = options.parse(argc, argv);

    if (result.count("help")) {
        std::cout << options.help({""}) << std::endl;
        exit(0);
    }
    return result;
}

/**
 *
 * @param sctpParams
 * @return -1 failed 0 success
 */
int buildInotify(sctp_params_t &sctpParams) {
    sctpParams.inotifyFD = inotify_init1(IN_NONBLOCK);
    if (sctpParams.inotifyFD == -1) {
        mdclog_write(MDCLOG_ERR, "Failed to init inotify (inotify_init1) %s", strerror(errno));
        close(sctpParams.rmrListenFd);
        rmr_close(sctpParams.rmrCtx);
        close(sctpParams.epoll_fd);
        return -1;
    }

    sctpParams.inotifyWD = inotify_add_watch(sctpParams.inotifyFD,
                                             (const char *)sctpParams.configFilePath.c_str(),
                                             (unsigned)IN_OPEN | (unsigned)IN_CLOSE_WRITE | (unsigned)IN_CLOSE_NOWRITE); //IN_CLOSE = (IN_CLOSE_WRITE | IN_CLOSE_NOWRITE)
    if (sctpParams.inotifyWD == -1) {
        mdclog_write(MDCLOG_ERR, "Failed to add directory : %s to  inotify (inotify_add_watch) %s",
                     sctpParams.configFilePath.c_str(),
                     strerror(errno));
        close(sctpParams.inotifyFD);
        return -1;
    }

    struct epoll_event event{};
    event.events = (EPOLLIN);
    event.data.fd = sctpParams.inotifyFD;
    // add listening RMR FD to epoll
    if (epoll_ctl(sctpParams.epoll_fd, EPOLL_CTL_ADD, sctpParams.inotifyFD, &event)) {
        mdclog_write(MDCLOG_ERR, "Failed to add inotify FD to epoll");
        close(sctpParams.inotifyFD);
        return -1;
    }
    return 0;
}

/**
 *
 * @param args
 * @return
 */
void listener(sctp_params_t *params) {
    int num_of_SCTP_messages = 0;
    auto totalTime = 0.0;
    mdclog_mdc_clean();
    mdclog_level_set(params->logLevel);

    std::thread::id this_id = std::this_thread::get_id();
    //save cout
    streambuf *oldCout = cout.rdbuf();
    ostringstream memCout;
    // create new cout
    cout.rdbuf(memCout.rdbuf());
    cout << this_id;
    //return to the normal cout
    cout.rdbuf(oldCout);

    char tid[32];
    memcpy(tid, memCout.str().c_str(), memCout.str().length() < 32 ? memCout.str().length() : 31);
    tid[memCout.str().length()] = 0;
    mdclog_mdc_add("thread id", tid);

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "started thread number %s", tid);
    }


    RmrMessagesBuffer_t rmrMessageBuffer{};
    //create and init RMR
    rmrMessageBuffer.rmrCtx = params->rmrCtx;

    auto *events = (struct epoll_event *) calloc(MAXEVENTS, sizeof(struct epoll_event));
    struct timespec end{0, 0};
    struct timespec start{0, 0};

    rmrMessageBuffer.rcvMessage = rmr_alloc_msg(rmrMessageBuffer.rmrCtx, RECEIVE_XAPP_BUFFER_SIZE);
    rmrMessageBuffer.sendMessage = rmr_alloc_msg(rmrMessageBuffer.rmrCtx, RECEIVE_XAPP_BUFFER_SIZE);

    memcpy(rmrMessageBuffer.ka_message, params->ka_message, params->ka_message_length);
    rmrMessageBuffer.ka_message_len = params->ka_message_length;
    rmrMessageBuffer.ka_message[rmrMessageBuffer.ka_message_len] = 0;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "keep alive message is : %s", rmrMessageBuffer.ka_message);
    }

    ReportingMessages_t message {};

//    for (int i = 0; i < MAX_RMR_BUFF_ARRY; i++) {
//        rmrMessageBuffer.rcvBufferedMessages[i] = rmr_alloc_msg(rmrMessageBuffer.rmrCtx, RECEIVE_XAPP_BUFFER_SIZE);
//        rmrMessageBuffer.sendBufferedMessages[i] = rmr_alloc_msg(rmrMessageBuffer.rmrCtx, RECEIVE_XAPP_BUFFER_SIZE);
//    }

    message.statCollector = StatCollector::GetInstance();

    while (true) {
        if (mdclog_level_get() >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG, "Start EPOLL Wait");
        }
        auto numOfEvents = epoll_wait(params->epoll_fd, events, MAXEVENTS, -1);
        if (numOfEvents < 0 && errno == EINTR) {
            if (mdclog_level_get() >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "got EINTR : %s", strerror(errno));
            }
            continue;
        }
        if (numOfEvents < 0) {
            mdclog_write(MDCLOG_ERR, "Epoll wait failed, errno = %s", strerror(errno));
            return;
        }
        for (auto i = 0; i < numOfEvents; i++) {
            if (mdclog_level_get() >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "handling epoll event %d out of %d", i + 1, numOfEvents);
            }
            clock_gettime(CLOCK_MONOTONIC, &message.message.time);
            start.tv_sec = message.message.time.tv_sec;
            start.tv_nsec = message.message.time.tv_nsec;


            if ((events[i].events & EPOLLERR) || (events[i].events & EPOLLHUP)) {
                handlepoll_error(events[i], message, rmrMessageBuffer, params);
            } else if (events[i].events & EPOLLOUT) {
                handleEinprogressMessages(events[i], message, rmrMessageBuffer, params);
            } else if (params->listenFD == events[i].data.fd) {
                if (mdclog_level_get() >= MDCLOG_INFO) {
                    mdclog_write(MDCLOG_INFO, "New connection request from sctp network\n");
                }
                // new connection is requested from RAN  start build connection
                while (true) {
                    struct sockaddr in_addr {};
                    socklen_t in_len;
                    char hostBuff[NI_MAXHOST];
                    char portBuff[NI_MAXSERV];

                    in_len = sizeof(in_addr);
                    auto *peerInfo = (ConnectedCU_t *)calloc(1, sizeof(ConnectedCU_t));
                    peerInfo->sctpParams = params;
                    peerInfo->fileDescriptor = accept(params->listenFD, &in_addr, &in_len);
                    if (peerInfo->fileDescriptor == -1) {
                        if ((errno == EAGAIN) || (errno == EWOULDBLOCK)) {
                            /* We have processed all incoming connections. */
                            break;
                        } else {
                            mdclog_write(MDCLOG_ERR, "Accept error, errno = %s", strerror(errno));
                            break;
                        }
                    }
                    if (setSocketNoBlocking(peerInfo->fileDescriptor) == -1) {
                        mdclog_write(MDCLOG_ERR, "setSocketNoBlocking failed to set new connection %s on port %s\n", hostBuff, portBuff);
                        close(peerInfo->fileDescriptor);
                        break;
                    }
                    auto  ans = getnameinfo(&in_addr, in_len,
                                            peerInfo->hostName, NI_MAXHOST,
                                            peerInfo->portNumber, NI_MAXSERV, (unsigned )((unsigned int)NI_NUMERICHOST | (unsigned int)NI_NUMERICSERV));
                    if (ans < 0) {
                        mdclog_write(MDCLOG_ERR, "Failed to get info on connection request. %s\n", strerror(errno));
                        close(peerInfo->fileDescriptor);
                        break;
                    }
                    if (mdclog_level_get() >= MDCLOG_DEBUG) {
                        mdclog_write(MDCLOG_DEBUG, "Accepted connection on descriptor %d (host=%s, port=%s)\n", peerInfo->fileDescriptor, peerInfo->hostName, peerInfo->portNumber);
                    }
                    peerInfo->isConnected = false;
                    peerInfo->gotSetup = false;
                    if (addToEpoll(params->epoll_fd,
                                   peerInfo,
                                   (EPOLLIN | EPOLLET),
                                   params->sctpMap, nullptr,
                                   0) != 0) {
                        break;
                    }
                    break;
                }
            } else if (params->rmrListenFd == events[i].data.fd) {
                // got message from XAPP
                num_of_XAPP_messages.fetch_add(1, std::memory_order_release);
                num_of_messages.fetch_add(1, std::memory_order_release);
                if (mdclog_level_get() >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "new message from RMR");
                }
                if (receiveXappMessages(params->sctpMap,
                                        rmrMessageBuffer,
                                        message.message.time) != 0) {
                    mdclog_write(MDCLOG_ERR, "Error handling Xapp message");
                }
            } else if (params->inotifyFD == events[i].data.fd) {
                mdclog_write(MDCLOG_INFO, "Got event from inotify (configuration update)");
                handleConfigChange(params);
            } else {
                /* We RMR_ERR_RETRY have data on the fd waiting to be read. Read and display it.
                 * We must read whatever data is available completely, as we are running
                 *  in edge-triggered mode and won't get a notification again for the same data. */
                num_of_messages.fetch_add(1, std::memory_order_release);
                if (mdclog_level_get() >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "new message from SCTP, epoll flags are : %0x", events[i].events);
                }
                receiveDataFromSctp(&events[i],
                                    params->sctpMap,
                                    num_of_SCTP_messages,
                                    rmrMessageBuffer,
                                    message.message.time);
            }

            clock_gettime(CLOCK_MONOTONIC, &end);
            if (mdclog_level_get() >= MDCLOG_INFO) {
                totalTime += ((end.tv_sec + 1.0e-9 * end.tv_nsec) -
                              ((double) start.tv_sec + 1.0e-9 * start.tv_nsec));
            }
            if (mdclog_level_get() >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "message handling is %ld seconds %ld nanoseconds",
                             end.tv_sec - start.tv_sec,
                             end.tv_nsec - start.tv_nsec);
            }
        }
    }
}

/**
 *
 * @param sctpParams
 */
void handleConfigChange(sctp_params_t *sctpParams) {
    char buf[4096] __attribute__ ((aligned(__alignof__(struct inotify_event))));
    const struct inotify_event *event;
    char *ptr;

    path p = (sctpParams->configFilePath + "/" + sctpParams->configFileName).c_str();
    auto endlessLoop = true;
    while (endlessLoop) {
        auto len = read(sctpParams->inotifyFD, buf, sizeof buf);
        if (len == -1) {
            if (errno != EAGAIN) {
                mdclog_write(MDCLOG_ERR, "read %s ", strerror(errno));
                endlessLoop = false;
                continue;
            }
            else {
                endlessLoop = false;
                continue;
            }
        }

        for (ptr = buf; ptr < buf + len; ptr += sizeof(struct inotify_event) + event->len) {
            event = (const struct inotify_event *)ptr;
            if (event->mask & (uint32_t)IN_ISDIR) {
                continue;
            }

            // the directory name
            if (sctpParams->inotifyWD == event->wd) {
                // not the directory
            }
            if (event->len) {
                auto  retVal = strcmp(sctpParams->configFileName.c_str(), event->name);
                if (retVal != 0) {
                    continue;
                }
            }
            // only the file we want
            if (event->mask & (uint32_t)IN_CLOSE_WRITE) {
                if (mdclog_level_get() >= MDCLOG_INFO) {
                    mdclog_write(MDCLOG_INFO, "Configuration file changed");
                }
                if (exists(p)) {
                    const int size = 2048;
                    auto fileSize = file_size(p);
                    if (fileSize > size) {
                        mdclog_write(MDCLOG_ERR, "File %s larger than %d", p.string().c_str(), size);
                        return;
                    }
                } else {
                    mdclog_write(MDCLOG_ERR, "Configuration File %s not exists", p.string().c_str());
                    return;
                }

                ReadConfigFile conf;
                if (conf.openConfigFile(p.string()) == -1) {
                    mdclog_write(MDCLOG_ERR, "Filed to open config file %s, %s",
                                 p.string().c_str(), strerror(errno));
                    return;
                }

                auto tmpStr = conf.getStringValue("loglevel");
                if (tmpStr.length() == 0) {
                    mdclog_write(MDCLOG_ERR, "illigal loglevel. Set loglevel to MDCLOG_INFO");
                    tmpStr = "info";
                }
                transform(tmpStr.begin(), tmpStr.end(), tmpStr.begin(), ::tolower);

                if ((tmpStr.compare("debug")) == 0) {
                    mdclog_write(MDCLOG_INFO, "Log level set to MDCLOG_DEBUG");
                    sctpParams->logLevel = MDCLOG_DEBUG;
                } else if ((tmpStr.compare("info")) == 0) {
                    mdclog_write(MDCLOG_INFO, "Log level set to MDCLOG_INFO");
                    sctpParams->logLevel = MDCLOG_INFO;
                } else if ((tmpStr.compare("warning")) == 0) {
                    mdclog_write(MDCLOG_INFO, "Log level set to MDCLOG_WARN");
                    sctpParams->logLevel = MDCLOG_WARN;
                } else if ((tmpStr.compare("error")) == 0) {
                    mdclog_write(MDCLOG_INFO, "Log level set to MDCLOG_ERR");
                    sctpParams->logLevel = MDCLOG_ERR;
                } else {
                    mdclog_write(MDCLOG_ERR, "illigal loglevel = %s. Set loglevel to MDCLOG_INFO", tmpStr.c_str());
                    sctpParams->logLevel = MDCLOG_INFO;
                }
                mdclog_level_set(sctpParams->logLevel);


                tmpStr = conf.getStringValue("trace");
                if (tmpStr.length() == 0) {
                    mdclog_write(MDCLOG_ERR, "illigal trace. Set trace to stop");
                    tmpStr = "stop";
                }

                transform(tmpStr.begin(), tmpStr.end(), tmpStr.begin(), ::tolower);
                if ((tmpStr.compare("start")) == 0) {
                    mdclog_write(MDCLOG_INFO, "Trace set to: start");
                    sctpParams->trace = true;
                } else if ((tmpStr.compare("stop")) == 0) {
                    mdclog_write(MDCLOG_INFO, "Trace set to: stop");
                    sctpParams->trace = false;
                } else {
                    mdclog_write(MDCLOG_ERR, "Trace was set to wrong value %s, set to stop", tmpStr.c_str());
                    sctpParams->trace = false;
                }
                jsonTrace = sctpParams->trace;
                endlessLoop = false;
            }
        }
    }
}

/**
 *
 * @param event
 * @param message
 * @param rmrMessageBuffer
 * @param params
 */
void handleEinprogressMessages(struct epoll_event &event,
                               ReportingMessages_t &message,
                               RmrMessagesBuffer_t &rmrMessageBuffer,
                               sctp_params_t *params) {
    auto *peerInfo = (ConnectedCU_t *)event.data.ptr;
    memcpy(message.message.enodbName, peerInfo->enodbName, sizeof(peerInfo->enodbName));

    mdclog_write(MDCLOG_INFO, "file descriptor %d got EPOLLOUT", peerInfo->fileDescriptor);
    auto retVal = 0;
    socklen_t retValLen = 0;
    auto rc = getsockopt(peerInfo->fileDescriptor, SOL_SOCKET, SO_ERROR, &retVal, &retValLen);
    if (rc != 0 || retVal != 0) {
        if (rc != 0) {
            rmrMessageBuffer.sendMessage->len = snprintf((char *)rmrMessageBuffer.sendMessage->payload, 256,
                                                         "%s|Failed SCTP Connection, after EINPROGRESS the getsockopt%s",
                                                         peerInfo->enodbName, strerror(errno));
        } else if (retVal != 0) {
            rmrMessageBuffer.sendMessage->len = snprintf((char *)rmrMessageBuffer.sendMessage->payload, 256,
                                                         "%s|Failed SCTP Connection after EINPROGRESS, SO_ERROR",
                                                         peerInfo->enodbName);
        }

        message.message.asndata = rmrMessageBuffer.sendMessage->payload;
        message.message.asnLength = rmrMessageBuffer.sendMessage->len;
        mdclog_write(MDCLOG_ERR, "%s", rmrMessageBuffer.sendMessage->payload);
        message.message.direction = 'N';
        if (sendRequestToXapp(message, RIC_SCTP_CONNECTION_FAILURE, rmrMessageBuffer) != 0) {
            mdclog_write(MDCLOG_ERR, "SCTP_CONNECTION_FAIL message failed to send to xAPP");
        }
        memset(peerInfo->asnData, 0, peerInfo->asnLength);
        peerInfo->asnLength = 0;
        peerInfo->mtype = 0;
        return;
    }

    peerInfo->isConnected = true;

    if (modifyToEpoll(params->epoll_fd, peerInfo, (EPOLLIN | EPOLLET), params->sctpMap, peerInfo->enodbName,
                      peerInfo->mtype) != 0) {
        mdclog_write(MDCLOG_ERR, "epoll_ctl EPOLL_CTL_MOD");
        return;
    }

    message.message.asndata = (unsigned char *)peerInfo->asnData;
    message.message.asnLength = peerInfo->asnLength;
    message.message.messageType = peerInfo->mtype;
    memcpy(message.message.enodbName, peerInfo->enodbName, sizeof(peerInfo->enodbName));
    num_of_messages.fetch_add(1, std::memory_order_release);
    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "send the delayed SETUP/ENDC SETUP to sctp for %s",
                     message.message.enodbName);
    }
    if (sendSctpMsg(peerInfo, message, params->sctpMap) != 0) {
        if (mdclog_level_get() >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG, "Error write to SCTP  %s %d", __func__, __LINE__);
        }
        return;
    }

    memset(peerInfo->asnData, 0, peerInfo->asnLength);
    peerInfo->asnLength = 0;
    peerInfo->mtype = 0;
}


void handlepoll_error(struct epoll_event &event,
                      ReportingMessages_t &message,
                      RmrMessagesBuffer_t &rmrMessageBuffer,
                      sctp_params_t *params) {
    if (event.data.fd != params->rmrListenFd) {
        auto *peerInfo = (ConnectedCU_t *)event.data.ptr;
        mdclog_write(MDCLOG_ERR, "epoll error, events %0x on fd %d, RAN NAME : %s",
                     event.events, peerInfo->fileDescriptor, peerInfo->enodbName);

        rmrMessageBuffer.sendMessage->len = snprintf((char *)rmrMessageBuffer.sendMessage->payload, 256,
                                                     "%s|Failed SCTP Connection",
                                                     peerInfo->enodbName);
        message.message.asndata = rmrMessageBuffer.sendMessage->payload;
        message.message.asnLength = rmrMessageBuffer.sendMessage->len;

        memcpy(message.message.enodbName, peerInfo->enodbName, sizeof(peerInfo->enodbName));
        message.message.direction = 'N';
        if (sendRequestToXapp(message, RIC_SCTP_CONNECTION_FAILURE, rmrMessageBuffer) != 0) {
            mdclog_write(MDCLOG_ERR, "SCTP_CONNECTION_FAIL message failed to send to xAPP");
        }

        close(peerInfo->fileDescriptor);
        params->sctpMap->erase(peerInfo->enodbName);
        cleanHashEntry((ConnectedCU_t *) event.data.ptr, params->sctpMap);
    } else {
        mdclog_write(MDCLOG_ERR, "epoll error, events %0x on RMR FD", event.events);
    }
}
/**
 *
 * @param socket
 * @return
 */
int setSocketNoBlocking(int socket) {
    auto flags = fcntl(socket, F_GETFL, 0);

    if (flags == -1) {
        mdclog_mdc_add("func", "fcntl");
        mdclog_write(MDCLOG_ERR, "%s, %s", __FUNCTION__, strerror(errno));
        mdclog_mdc_clean();
        return -1;
    }

    flags = (unsigned) flags | (unsigned) O_NONBLOCK;
    if (fcntl(socket, F_SETFL, flags) == -1) {
        mdclog_mdc_add("func", "fcntl");
        mdclog_write(MDCLOG_ERR, "%s, %s", __FUNCTION__, strerror(errno));
        mdclog_mdc_clean();
        return -1;
    }

    return 0;
}

/**
 *
 * @param val
 * @param m
 */
void cleanHashEntry(ConnectedCU_t *val, Sctp_Map_t *m) {
    char *dummy;
    auto port = (uint16_t) strtol(val->portNumber, &dummy, 10);
    char searchBuff[2048]{};

    snprintf(searchBuff, sizeof searchBuff, "host:%s:%d", val->hostName, port);
    m->erase(searchBuff);

    m->erase(val->enodbName);
    free(val);
}

/**
 *
 * @param fd file discriptor
 * @param data the asn data to send
 * @param len  length of the data
 * @param enodbName the enodbName as in the map for printing purpose
 * @param m map host information
 * @param mtype message number
 * @return 0 success, anegative number on fail
 */
int sendSctpMsg(ConnectedCU_t *peerInfo, ReportingMessages_t &message, Sctp_Map_t *m) {
    auto loglevel = mdclog_level_get();
    int fd = peerInfo->fileDescriptor;
    if (loglevel >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "Send SCTP message for CU %s, %s",
                     message.message.enodbName, __FUNCTION__);
    }

    while (true) {
        if (send(fd,message.message.asndata, message.message.asnLength,MSG_NOSIGNAL) < 0) {
            if (errno == EINTR) {
                continue;
            }
            mdclog_write(MDCLOG_ERR, "error writing to CU a message, %s ", strerror(errno));
            if (!peerInfo->isConnected) {
                mdclog_write(MDCLOG_ERR, "connection to CU %s is still in progress.", message.message.enodbName);
                return -1;
            }
            cleanHashEntry(peerInfo, m);
            close(fd);
            char key[MAX_ENODB_NAME_SIZE * 2];
            snprintf(key, MAX_ENODB_NAME_SIZE * 2, "msg:%s|%d", message.message.enodbName,
                     message.message.messageType);
            if (loglevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "remove key = %s from %s at line %d", key, __FUNCTION__, __LINE__);
            }
            auto tmp = m->find(key);
            if (tmp) {
                free(tmp);
            }
            m->erase(key);
            return -1;
        }
        // TODO remove stat update
        //message.statCollector->incSentMessage(string(message.message.enodbName));
        message.message.direction = 'D';
        // send report.buffer of size
        buildJsonMessage(message);

        if (loglevel >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG,
                         "SCTP message for CU %s sent from %s",
                         message.message.enodbName,
                         __FUNCTION__);
        }
        return 0;
    }
}

/**
 *
 * @param message
 * @param rmrMessageBuffer
 */
void getRequestMetaData(ReportingMessages_t &message, RmrMessagesBuffer_t &rmrMessageBuffer) {
    rmr_get_meid(rmrMessageBuffer.rcvMessage, (unsigned char *) (message.message.enodbName));

    message.message.asndata = rmrMessageBuffer.rcvMessage->payload;
    message.message.asnLength = rmrMessageBuffer.rcvMessage->len;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "Message from Xapp RAN name = %s message length = %ld",
                     message.message.enodbName, (unsigned long) message.message.asnLength);
    }
}


/**
 * 09/11/2020 - Leo: parse json message
 *
 * @param metaData all the data strip to structure
 * @param data the data recived from xAPP
 * @return 0 success all other values are fault
 */
int parseMessageData(ReportingMessages_t &message, char *data, char *host, uint16_t &port) {

    auto loglevel = mdclog_level_get();

    Document document;
    document.Parse(data);

    const char* ran_ip;
    if (document.HasMember("ranIp")) {
        ran_ip = document["ranIp"].GetString();
        std::cout << "ranIP: " << ran_ip << ", length: " << strlen(ran_ip) << std::endl;
        memcpy(host, ran_ip, strlen(ran_ip));
    }
    else {
        mdclog_write(MDCLOG_ERR, "ranIp not provided");
        return -1;
    }

    int ran_port;
    if (document.HasMember("ranPort")) {
        ran_port = document["ranPort"].GetInt();
        std::cout << "ranPort: " << ran_port << std::endl;
        port = (uint16_t) ran_port;
    }
    else {
        mdclog_write(MDCLOG_ERR, "ranPort not provided");
        return -2;
    }

    const char* ran_name;
    if (document.HasMember("ranName")) {
        ran_name = document["ranName"].GetString();
        std::cout << "ranName: " << ran_name << ", length: " << strlen(ran_name) << std::endl;
        memcpy(message.message.enodbName, ran_name, strlen(ran_name));
    }
    else {
        mdclog_write(MDCLOG_ERR, "ranName not provided");
        return -3;
    }

    const char* message_payload;
    if (document.HasMember("payload")) {
        message_payload = document["payload"].GetString();
        std::cout << "payload: " << message_payload << ", length: " << strlen(message_payload) << std::endl;
        message.message.asndata = (unsigned char *) message_payload;
        std::cout << "message.message.asndata: " << message.message.asndata << std::endl;
        message.message.asnLength = (uint16_t) strlen(message_payload);
        std::cout << "message.message.asnLength: " << message.message.asnLength << std::endl;
    }
    else {
        mdclog_write(MDCLOG_ERR, "Payload not provided");
        return -4;
    }

    std::cout << "Step 1" << std::endl;

    // val = strtok_r(nullptr, delimiter, &tmp);
    // if (val != nullptr) {
    //     mdclog_write(MDCLOG_DEBUG, "ASN length parameter from message = %s", val);
    //     if (mdclog_level_get() >= MDCLOG_DEBUG) {
    //         mdclog_write(MDCLOG_DEBUG, "ASN length parameter from message = %s", val);
    //     }
    //     char *dummy;
    //     message.message.asnLength = (uint16_t) strtol(val, &dummy, 10);
    // } else {
    //     mdclog_write(MDCLOG_ERR, "wrong ASN length for setup request %s", data);
    //     return -4;
    // }
    //
    // message.message.asndata = (unsigned char *)tmp;  // tmp is local but point to the location in data

    if (loglevel >= MDCLOG_INFO) {
        mdclog_write(MDCLOG_INFO, "Message from Xapp RAN name = %s host address = %s port = %d",
                     message.message.enodbName, host, port);
    }

    std::cout << "Step 2" << std::endl;

    return 0;
}


/**
 *
 * @param events
 * @param sctpMap
 * @param numOfMessages
 * @param rmrMessageBuffer
 * @param ts
 * @return
 */
int receiveDataFromSctp(struct epoll_event *events,
                        Sctp_Map_t *sctpMap,
                        int &numOfMessages,
                        RmrMessagesBuffer_t &rmrMessageBuffer,
                        struct timespec &ts) {
    /* We have data on the fd waiting to be read. Read and display it.
 * We must read whatever data is available completely, as we are running
 *  in edge-triggered mode and won't get a notification again for the same data. */
    ReportingMessages_t message {};
    auto done = 0;
    auto loglevel = mdclog_level_get();

    // get the identity of the interface
    message.peerInfo = (ConnectedCU_t *)events->data.ptr;

    message.statCollector = StatCollector::GetInstance();
    struct timespec start{0, 0};
    struct timespec decodestart{0, 0};
    struct timespec end{0, 0};

    E2AP_PDU_t *pdu = nullptr;


    while (true) {
        if (loglevel >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG, "Start Read from SCTP %d fd", message.peerInfo->fileDescriptor);
            clock_gettime(CLOCK_MONOTONIC, &start);
        }
        // read the buffer directly to rmr payload
        message.message.asndata = rmrMessageBuffer.sendMessage->payload;
        message.message.asnLength = rmrMessageBuffer.sendMessage->len =
                read(message.peerInfo->fileDescriptor, rmrMessageBuffer.sendMessage->payload, RECEIVE_SCTP_BUFFER_SIZE);

        if (loglevel >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG, "Finish Read from SCTP %d fd message length = %ld",
                         message.peerInfo->fileDescriptor, message.message.asnLength);
        }

        memcpy(message.message.enodbName, message.peerInfo->enodbName, sizeof(message.peerInfo->enodbName));
        message.statCollector->incRecvMessage(string(message.message.enodbName));
        message.message.direction = 'U';
        message.message.time.tv_nsec = ts.tv_nsec;
        message.message.time.tv_sec = ts.tv_sec;

        if (message.message.asnLength < 0) {
            if (errno == EINTR) {
                continue;
            }
            /* If errno == EAGAIN, that means we have read all
               data. So goReportingMessages_t back to the main loop. */
            if (errno != EAGAIN) {
                mdclog_write(MDCLOG_ERR, "Read error, %s ", strerror(errno));
                done = 1;
            } else if (loglevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "EAGAIN - descriptor = %d", message.peerInfo->fileDescriptor);
            }
            break;
        } else if (message.message.asnLength == 0) {
            /* End of file. The remote has closed the connection. */
            if (loglevel >= MDCLOG_INFO) {
                mdclog_write(MDCLOG_INFO, "END of File Closed connection - descriptor = %d",
                             message.peerInfo->fileDescriptor);
            }
            done = 1;
            break;
        }

        if (loglevel >= MDCLOG_DEBUG) {
            char printBuffer[4096]{};
            char *tmp = printBuffer;
            for (size_t i = 0; i < (size_t)message.message.asnLength; ++i) {
                snprintf(tmp, 3, "%02x", message.message.asndata[i]);
                tmp += 2;
            }
            printBuffer[message.message.asnLength] = 0;
            clock_gettime(CLOCK_MONOTONIC, &end);
            mdclog_write(MDCLOG_DEBUG, "Before Encoding E2AP PDU for : %s, Read time is : %ld seconds, %ld nanoseconds",
                         message.peerInfo->enodbName, end.tv_sec - start.tv_sec, end.tv_nsec - start.tv_nsec);
            mdclog_write(MDCLOG_DEBUG, "PDU buffer length = %ld, data =  : %s", message.message.asnLength,
                         printBuffer);
            clock_gettime(CLOCK_MONOTONIC, &decodestart);
        }

        auto rval = asn_decode(nullptr, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2AP_PDU, (void **) &pdu,
                          message.message.asndata, message.message.asnLength);
        if (rval.code != RC_OK) {
            mdclog_write(MDCLOG_ERR, "Error %d Decoding (unpack) E2AP PDU from RAN : %s", rval.code,
                         message.peerInfo->enodbName);
            //todo may need reset to pdu
            break;
        }

        if (loglevel >= MDCLOG_DEBUG) {
            clock_gettime(CLOCK_MONOTONIC, &end);
            mdclog_write(MDCLOG_DEBUG, "After Encoding E2AP PDU for : %s, Read time is : %ld seconds, %ld nanoseconds",
                         message.peerInfo->enodbName, end.tv_sec - decodestart.tv_sec, end.tv_nsec - decodestart.tv_nsec);
            char *printBuffer;
            size_t size;
            FILE *stream = open_memstream(&printBuffer, &size);
            asn_fprint(stream, &asn_DEF_E2AP_PDU, pdu);
            mdclog_write(MDCLOG_DEBUG, "Encoding E2AP PDU past : %s", printBuffer);
            clock_gettime(CLOCK_MONOTONIC, &decodestart);
        }

        switch (pdu->present) {
            case E2AP_PDU_PR_initiatingMessage: {//initiating message
                asnInitiatingRequest(pdu, sctpMap,message, rmrMessageBuffer);
                break;
            }
            case E2AP_PDU_PR_successfulOutcome: { //successful outcome
                asnSuccsesfulMsg(pdu, sctpMap, message,  rmrMessageBuffer);
                break;
            }
            case E2AP_PDU_PR_unsuccessfulOutcome: { //Unsuccessful Outcome
                asnUnSuccsesfulMsg(pdu, sctpMap, message, rmrMessageBuffer);
                break;
            }
            default:
                mdclog_write(MDCLOG_ERR, "Unknown index %d in E2AP PDU", pdu->present);
                break;
        }
        if (loglevel >= MDCLOG_DEBUG) {
            clock_gettime(CLOCK_MONOTONIC, &end);
            mdclog_write(MDCLOG_DEBUG,
                         "After processing message and sent to rmr for : %s, Read time is : %ld seconds, %ld nanoseconds",
                         message.peerInfo->enodbName, end.tv_sec - decodestart.tv_sec, end.tv_nsec - decodestart.tv_nsec);
        }
        numOfMessages++;
        if (pdu != nullptr) {
            ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, pdu);
            //ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, pdu);
            //pdu = nullptr;
        }
    }

    if (done) {
        if (loglevel >= MDCLOG_INFO) {
            mdclog_write(MDCLOG_INFO, "Closed connection - descriptor = %d", message.peerInfo->fileDescriptor);
        }
        message.message.asnLength = rmrMessageBuffer.sendMessage->len =
                snprintf((char *)rmrMessageBuffer.sendMessage->payload,
                         256,
                         "%s|CU disconnected unexpectedly",
                         message.peerInfo->enodbName);
        message.message.asndata = rmrMessageBuffer.sendMessage->payload;

        if (sendRequestToXapp(message,
                              RIC_SCTP_CONNECTION_FAILURE,
                              rmrMessageBuffer) != 0) {
            mdclog_write(MDCLOG_ERR, "SCTP_CONNECTION_FAIL message failed to send to xAPP");
        }

        /* Closing descriptor make epoll remove it from the set of descriptors which are monitored. */
        close(message.peerInfo->fileDescriptor);
        cleanHashEntry((ConnectedCU_t *) events->data.ptr, sctpMap);
    }
    if (loglevel >= MDCLOG_DEBUG) {
        clock_gettime(CLOCK_MONOTONIC, &end);
        mdclog_write(MDCLOG_DEBUG, "from receive SCTP to send RMR time is %ld seconds and %ld nanoseconds",
                     end.tv_sec - start.tv_sec, end.tv_nsec - start.tv_nsec);

    }
    return 0;
}

static void buildAndsendSetupRequest(ReportingMessages_t &message,
                                     RmrMessagesBuffer_t &rmrMessageBuffer,
                                     E2AP_PDU_t *pdu,
                                     vector<string> &repValues) {
    auto logLevel = mdclog_level_get();
    // now we can send the data to e2Mgr
    auto buffer_size = RECEIVE_SCTP_BUFFER_SIZE * 2;
    unsigned char buffer[RECEIVE_SCTP_BUFFER_SIZE * 2];
    auto *rmrMsg = rmr_alloc_msg(rmrMessageBuffer.rmrCtx, buffer_size);
    // encode to xml
    auto er = asn_encode_to_buffer(nullptr, ATS_BASIC_XER, &asn_DEF_E2AP_PDU, pdu, buffer, buffer_size);
    if (er.encoded == -1) {
        mdclog_write(MDCLOG_ERR, "encoding of %s failed, %s", asn_DEF_E2AP_PDU.name, strerror(errno));
    } else if (er.encoded > (ssize_t) buffer_size) {
        mdclog_write(MDCLOG_ERR, "Buffer of size %d is to small for %s, at %s line %d",
                     (int) buffer_size,
                     asn_DEF_E2AP_PDU.name, __func__, __LINE__);
    } else {
        string messageType("E2setupRequest");
        string ieName("E2setupRequestIEs");
        buffer[er.encoded] = '\0';
        buildXmlData(messageType, ieName, repValues, buffer, (size_t)er.encoded);

//        string xmlStr = (char *)buffer;
//        auto removeSpaces = [] (string str) -> string {
//            str.erase(remove(str.begin(), str.end(), ' '), str.end());
//            str.erase(remove(str.begin(), str.end(), '\t'), str.end());
//            return str;
//        };
//
//        xmlStr = removeSpaces(xmlStr);
//        // we have the XML
//        rmrMsg->len = snprintf((char *)rmrMsg->payload, RECEIVE_SCTP_BUFFER_SIZE * 2, "%s:%d|%s",
//                               message.peerInfo->sctpParams->myIP.c_str(),
//                               message.peerInfo->sctpParams->rmrPort,
//                               xmlStr.c_str());
        rmrMsg->len = snprintf((char *)rmrMsg->payload, RECEIVE_SCTP_BUFFER_SIZE * 2, "%s:%d|%s",
                               message.peerInfo->sctpParams->myIP.c_str(),
                               message.peerInfo->sctpParams->rmrPort,
                               buffer);
        if (logLevel >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG, "Setup request of size %d :\n %s\n", rmrMsg->len, rmrMsg->payload);
        }
        // send to RMR
        message.message.messageType = rmrMsg->mtype = RIC_E2_SETUP_REQ;
        rmrMsg->state = 0;
        rmr_bytes2meid(rmrMsg, (unsigned char *) message.message.enodbName, strlen(message.message.enodbName));

        static unsigned char tx[32];
        snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
        rmr_bytes2xact(rmrMsg, tx, strlen((const char *) tx));

        rmrMsg = rmr_send_msg(rmrMessageBuffer.rmrCtx, rmrMsg);
        if (rmrMsg == nullptr) {
            mdclog_write(MDCLOG_ERR, "RMR failed to send returned nullptr");
        } else if (rmrMsg->state != 0) {
            char meid[RMR_MAX_MEID]{};
            if (rmrMsg->state == RMR_ERR_RETRY) {
                usleep(5);
                rmrMsg->state = 0;
                mdclog_write(MDCLOG_INFO, "RETRY sending Message %d to Xapp from %s",
                             rmrMsg->mtype, rmr_get_meid(rmrMsg, (unsigned char *) meid));
                rmrMsg = rmr_send_msg(rmrMessageBuffer.rmrCtx, rmrMsg);
                if (rmrMsg == nullptr) {
                    mdclog_write(MDCLOG_ERR, "RMR failed send returned nullptr");
                } else if (rmrMsg->state != 0) {
                    mdclog_write(MDCLOG_ERR,
                                 "RMR Retry failed %s sending request %d to Xapp from %s",
                                 translateRmrErrorMessages(rmrMsg->state).c_str(),
                                 rmrMsg->mtype,
                                 rmr_get_meid(rmrMsg, (unsigned char *) meid));
                }
            } else {
                mdclog_write(MDCLOG_ERR, "RMR failed: %s. sending request %d to Xapp from %s",
                             translateRmrErrorMessages(rmrMsg->state).c_str(),
                             rmrMsg->mtype,
                             rmr_get_meid(rmrMsg, (unsigned char *) meid));
            }
        }
        message.peerInfo->gotSetup = true;
        buildJsonMessage(message);
        if (rmrMsg != nullptr) {
            rmr_free_msg(rmrMsg);
        }
    }

}

int RAN_Function_list_To_Vector(RANfunctions_List_t& list, vector <string> &runFunXML_v) {
    auto index = 0;
    runFunXML_v.clear();
    for (auto j = 0; j < list.list.count; j++) {
        auto *raNfunctionItemIEs = (RANfunction_ItemIEs_t *)list.list.array[j];
        if (raNfunctionItemIEs->id == ProtocolIE_ID_id_RANfunction_Item &&
            (raNfunctionItemIEs->value.present == RANfunction_ItemIEs__value_PR_RANfunction_Item)) {
            // encode to xml
            E2SM_gNB_NRT_RANfunction_Definition_t *ranFunDef = nullptr;
            auto rval = asn_decode(nullptr, ATS_ALIGNED_BASIC_PER,
                                   &asn_DEF_E2SM_gNB_NRT_RANfunction_Definition,
                                   (void **)&ranFunDef,
                                   raNfunctionItemIEs->value.choice.RANfunction_Item.ranFunctionDefinition.buf,
                                   raNfunctionItemIEs->value.choice.RANfunction_Item.ranFunctionDefinition.size);
            if (rval.code != RC_OK) {
                mdclog_write(MDCLOG_ERR, "Error %d Decoding (unpack) E2SM message from : %s",
                             rval.code,
                             asn_DEF_E2SM_gNB_NRT_RANfunction_Definition.name);
                return -1;
            }

//                        if (mdclog_level_get() >= MDCLOG_DEBUG) {
//                            char *printBuffer;
//                            size_t size;
//                            FILE *stream = open_memstream(&printBuffer, &size);
//                            asn_fprint(stream, &asn_DEF_E2SM_gNB_NRT_RANfunction_Definition, ranFunDef);
//                            mdclog_write(MDCLOG_DEBUG, "Encoding E2SM %s PDU past : %s",
//                                         asn_DEF_E2SM_gNB_NRT_RANfunction_Definition.name,
//                                         printBuffer);
//                        }
            auto xml_buffer_size = RECEIVE_SCTP_BUFFER_SIZE * 2;
            unsigned char xml_buffer[RECEIVE_SCTP_BUFFER_SIZE * 2];
            // encode to xml
            auto er = asn_encode_to_buffer(nullptr,
                                           ATS_BASIC_XER,
                                           &asn_DEF_E2SM_gNB_NRT_RANfunction_Definition,
                                           ranFunDef,
                                           xml_buffer,
                                           xml_buffer_size);
            if (er.encoded == -1) {
                mdclog_write(MDCLOG_ERR, "encoding of %s failed, %s",
                             asn_DEF_E2SM_gNB_NRT_RANfunction_Definition.name,
                             strerror(errno));
            } else if (er.encoded > (ssize_t)xml_buffer_size) {
                mdclog_write(MDCLOG_ERR, "Buffer of size %d is to small for %s, at %s line %d",
                             (int) xml_buffer_size,
                             asn_DEF_E2SM_gNB_NRT_RANfunction_Definition.name, __func__, __LINE__);
            } else {
                if (mdclog_level_get() >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "Encoding E2SM %s PDU number %d : %s",
                                 asn_DEF_E2SM_gNB_NRT_RANfunction_Definition.name,
                                 index++,
                                 xml_buffer);
                }
                string runFuncs = (char *)(xml_buffer);
                runFunXML_v.emplace_back(runFuncs);
            }
        }
    }
    return 0;
}



int collectSetupAndServiceUpdate_RequestData(E2AP_PDU_t *pdu,
                                             Sctp_Map_t *sctpMap,
                                             ReportingMessages_t &message,
                                             vector <string> &RANfunctionsAdded_v,
                                             vector <string> &RANfunctionsModified_v) {
    memset(message.peerInfo->enodbName, 0 , MAX_ENODB_NAME_SIZE);
    for (auto i = 0; i < pdu->choice.initiatingMessage->value.choice.E2setupRequest.protocolIEs.list.count; i++) {
        auto *ie = pdu->choice.initiatingMessage->value.choice.E2setupRequest.protocolIEs.list.array[i];
        if (ie->id == ProtocolIE_ID_id_GlobalE2node_ID) {
            // get the ran name for meid
            if (ie->value.present == E2setupRequestIEs__value_PR_GlobalE2node_ID) {
                if (buildRanName(message.peerInfo->enodbName, ie) < 0) {
                    mdclog_write(MDCLOG_ERR, "Bad param in E2setupRequestIEs GlobalE2node_ID.\n");
                    // no mesage will be sent
                    return -1;
                }
                memcpy(message.message.enodbName, message.peerInfo->enodbName, strlen(message.peerInfo->enodbName));
                sctpMap->setkey(message.message.enodbName, message.peerInfo);
            }
        } else if (ie->id == ProtocolIE_ID_id_RANfunctionsAdded) {
            if (ie->value.present == E2setupRequestIEs__value_PR_RANfunctions_List) {
                if (mdclog_level_get() >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "Run function list have %d entries",
                                 ie->value.choice.RANfunctions_List.list.count);
                }
                if (RAN_Function_list_To_Vector(ie->value.choice.RANfunctions_List, RANfunctionsAdded_v) != 0 ) {
                    return -1;
                }
            }
        } else if (ie->id == ProtocolIE_ID_id_RANfunctionsModified) {
            if (ie->value.present == E2setupRequestIEs__value_PR_RANfunctions_List) {
                if (mdclog_level_get() >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "Run function list have %d entries",
                                 ie->value.choice.RANfunctions_List.list.count);
                }
                if (RAN_Function_list_To_Vector(ie->value.choice.RANfunctions_List, RANfunctionsModified_v) != 0 ) {
                    return -1;
                }
            }
        }
    }
    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "Run function vector have %ld entries",
                     RANfunctionsAdded_v.size());
    }
    return 0;
}
/**
 *
 * @param pdu
 * @param message
 * @param rmrMessageBuffer
 */
void asnInitiatingRequest(E2AP_PDU_t *pdu,
                          Sctp_Map_t *sctpMap,
                          ReportingMessages_t &message,
                          RmrMessagesBuffer_t &rmrMessageBuffer) {
    auto logLevel = mdclog_level_get();
    auto procedureCode = ((InitiatingMessage_t *) pdu->choice.initiatingMessage)->procedureCode;
    if (logLevel >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "Initiating message %ld\n", procedureCode);
    }
    switch (procedureCode) {
        case ProcedureCode_id_E2setup: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got E2setup");
            }

            // first get the message as XML buffer
            auto setup_xml_buffer_size = RECEIVE_SCTP_BUFFER_SIZE * 2;
            unsigned char setup_xml_buffer[RECEIVE_SCTP_BUFFER_SIZE * 2];

            auto er = asn_encode_to_buffer(nullptr, ATS_BASIC_XER, &asn_DEF_E2AP_PDU, pdu, setup_xml_buffer, setup_xml_buffer_size);
            if (er.encoded == -1) {
                mdclog_write(MDCLOG_ERR, "encoding of %s failed, %s", asn_DEF_E2AP_PDU.name, strerror(errno));
                break;
            } else if (er.encoded > (ssize_t) setup_xml_buffer_size) {
                mdclog_write(MDCLOG_ERR, "Buffer of size %d is to small for %s, at %s line %d",
                             (int)setup_xml_buffer_size,
                             asn_DEF_E2AP_PDU.name, __func__, __LINE__);
                break;
            }
            std::string xmlString(setup_xml_buffer_size,  setup_xml_buffer_size + er.encoded);

            vector <string> RANfunctionsAdded_v;
            vector <string> RANfunctionsModified_v;
            RANfunctionsAdded_v.clear();
            RANfunctionsModified_v.clear();
            if (collectSetupAndServiceUpdate_RequestData(pdu, sctpMap, message,
                    RANfunctionsAdded_v, RANfunctionsModified_v) != 0) {
                break;
            }

            buildAndsendSetupRequest(message, rmrMessageBuffer, pdu, RANfunctionsAdded_v);
            break;
        }
        case ProcedureCode_id_RICserviceUpdate: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICserviceUpdate %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_SERVICE_UPDATE, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "RIC_SERVICE_UPDATE message failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_ErrorIndication: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got ErrorIndication %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_ERROR_INDICATION, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "RIC_ERROR_INDICATION failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_Reset: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got Reset %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_X2_RESET, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "RIC_X2_RESET message failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_RICcontrol: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICcontrol %s", message.message.enodbName);
            }
            break;
        }
        case ProcedureCode_id_RICindication: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICindication - initiating %s", message.message.enodbName);
            }
            for (auto i = 0; i < pdu->choice.initiatingMessage->value.choice.RICindication.protocolIEs.list.count; i++) {
                auto messageSent = false;
                RICindication_IEs_t *ie = pdu->choice.initiatingMessage->value.choice.RICindication.protocolIEs.list.array[i];
                if (logLevel >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "ie type (ProtocolIE_ID) = %ld", ie->id);
                }
                if (ie->id == ProtocolIE_ID_id_RICrequestID) {
                    if (logLevel >= MDCLOG_DEBUG) {
                        mdclog_write(MDCLOG_DEBUG, "Got RIC requestId entry, ie type (ProtocolIE_ID) = %ld", ie->id);
                    }
                    if (ie->value.present == RICindication_IEs__value_PR_RICrequestID) {
                        static unsigned char tx[32];
                        message.message.messageType = rmrMessageBuffer.sendMessage->mtype = RIC_INDICATION;
                        snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
                        rmr_bytes2xact(rmrMessageBuffer.sendMessage, tx, strlen((const char *) tx));
                        int exit_status = rmr_bytes2meid(rmrMessageBuffer.sendMessage,
                                       (unsigned char *)message.message.enodbName,
                                       strlen(message.message.enodbName));
                        rmrMessageBuffer.sendMessage->state = 0;
                        
                        // set sub_id to ricRequestorID for future lookup in rmr routing table.
                        // it was set to ricInstanceID before
                        rmrMessageBuffer.sendMessage->sub_id = ie->value.choice.RICrequestID.ricRequestorID;
                        // rmrMessageBuffer.sendMessage->sub_id = (int)ie->value.choice.RICrequestID.ricInstanceID;
                        
                        unsigned char *me_id;
			unsigned char* me_id_ptr = rmr_get_meid(rmrMessageBuffer.sendMessage, me_id);
			mdclog_write(MDCLOG_DEBUG, "Received MEID: %s, exit_status %d, ptr %s", me_id, exit_status, me_id_ptr);                        

                        //ie->value.choice.RICrequestID.ricInstanceID;
                        if (mdclog_level_get() >= MDCLOG_DEBUG) {
                            mdclog_write(MDCLOG_DEBUG, "sub id = %d, mtype = %d, ric instance id %ld, requestor id = %ld",
                                         rmrMessageBuffer.sendMessage->sub_id,
                                         rmrMessageBuffer.sendMessage->mtype,
                                         ie->value.choice.RICrequestID.ricInstanceID,
                                         ie->value.choice.RICrequestID.ricRequestorID);
                        }
                        sendRmrMessage(rmrMessageBuffer, message);
                        messageSent = true;
                    } else {
                        mdclog_write(MDCLOG_ERR, "RIC request id missing illigal request");
                    }
                }
                if (messageSent) {
                    break;
                }
            }
            break;
        }
        case ProcedureCode_id_RICserviceQuery: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICserviceQuery %s", message.message.enodbName);
            }
            break;
        }
        case ProcedureCode_id_RICsubscription: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICsubscription %s", message.message.enodbName);
            }
            break;
        }
        case ProcedureCode_id_RICsubscriptionDelete: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICsubscriptionDelete %s", message.message.enodbName);
            }
            break;
        }
        default: {
            mdclog_write(MDCLOG_ERR, "Undefined or not supported message = %ld", procedureCode);
            message.message.messageType = 0; // no RMR message type yet

            buildJsonMessage(message);

            break;
        }
    }
}

/**
 *
 * @param pdu
 * @param message
 * @param rmrMessageBuffer
 */
void asnSuccsesfulMsg(E2AP_PDU_t *pdu,
                      Sctp_Map_t *sctpMap,
                      ReportingMessages_t &message,
                      RmrMessagesBuffer_t &rmrMessageBuffer) {
    auto procedureCode = pdu->choice.successfulOutcome->procedureCode;
    auto logLevel = mdclog_level_get();
    if (logLevel >= MDCLOG_INFO) {
        mdclog_write(MDCLOG_INFO, "Successful Outcome %ld", procedureCode);
    }
    switch (procedureCode) {
        case ProcedureCode_id_E2setup: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got E2setup\n");
            }
            break;
        }
        case ProcedureCode_id_ErrorIndication: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got ErrorIndication %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_ERROR_INDICATION, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "RIC_ERROR_INDICATION failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_Reset: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got Reset %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_X2_RESET, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "RIC_X2_RESET message failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_RICcontrol: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICcontrol %s", message.message.enodbName);
            }
            for (auto i = 0;
                 i < pdu->choice.successfulOutcome->value.choice.RICcontrolAcknowledge.protocolIEs.list.count; i++) {
                auto messageSent = false;
                RICcontrolAcknowledge_IEs_t *ie = pdu->choice.successfulOutcome->value.choice.RICcontrolAcknowledge.protocolIEs.list.array[i];
                if (mdclog_level_get() >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "ie type (ProtocolIE_ID) = %ld", ie->id);
                }
                if (ie->id == ProtocolIE_ID_id_RICrequestID) {
                    if (mdclog_level_get() >= MDCLOG_DEBUG) {
                        mdclog_write(MDCLOG_DEBUG, "Got RIC requestId entry, ie type (ProtocolIE_ID) = %ld", ie->id);
                    }
                    if (ie->value.present == RICcontrolAcknowledge_IEs__value_PR_RICrequestID) {
                        message.message.messageType = rmrMessageBuffer.sendMessage->mtype = RIC_CONTROL_ACK;
                        rmrMessageBuffer.sendMessage->state = 0;
//                        rmrMessageBuffer.sendMessage->sub_id = (int) ie->value.choice.RICrequestID.ricRequestorID;
                        rmrMessageBuffer.sendMessage->sub_id = (int)ie->value.choice.RICrequestID.ricInstanceID;

                        static unsigned char tx[32];
                        snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
                        rmr_bytes2xact(rmrMessageBuffer.sendMessage, tx, strlen((const char *) tx));
                        rmr_bytes2meid(rmrMessageBuffer.sendMessage,
                                       (unsigned char *)message.message.enodbName,
                                       strlen(message.message.enodbName));

                        sendRmrMessage(rmrMessageBuffer, message);
                        messageSent = true;
                    } else {
                        mdclog_write(MDCLOG_ERR, "RIC request id missing illigal request");
                    }
                }
                if (messageSent) {
                    break;
                }
            }

            break;
        }
        case ProcedureCode_id_RICindication: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICindication %s", message.message.enodbName);
            }
            for (auto i = 0; i < pdu->choice.initiatingMessage->value.choice.RICindication.protocolIEs.list.count; i++) {
                auto messageSent = false;
                RICindication_IEs_t *ie = pdu->choice.initiatingMessage->value.choice.RICindication.protocolIEs.list.array[i];
                if (logLevel >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "ie type (ProtocolIE_ID) = %ld", ie->id);
                }
                if (ie->id == ProtocolIE_ID_id_RICrequestID) {
                    if (logLevel >= MDCLOG_DEBUG) {
                        mdclog_write(MDCLOG_DEBUG, "Got RIC requestId entry, ie type (ProtocolIE_ID) = %ld", ie->id);
                    }
                    if (ie->value.present == RICindication_IEs__value_PR_RICrequestID) {
                        static unsigned char tx[32];
                        message.message.messageType = rmrMessageBuffer.sendMessage->mtype = RIC_INDICATION;
                        snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
                        rmr_bytes2xact(rmrMessageBuffer.sendMessage, tx, strlen((const char *) tx));
                        rmr_bytes2meid(rmrMessageBuffer.sendMessage,
                                       (unsigned char *)message.message.enodbName,
                                       strlen(message.message.enodbName));
                        rmrMessageBuffer.sendMessage->state = 0;
                        rmrMessageBuffer.sendMessage->sub_id = (int)ie->value.choice.RICrequestID.ricInstanceID;
                        if (mdclog_level_get() >= MDCLOG_DEBUG) {
                            mdclog_write(MDCLOG_DEBUG, "RIC sub id = %d, message type = %d",
                                         rmrMessageBuffer.sendMessage->sub_id,
                                         rmrMessageBuffer.sendMessage->mtype);
                        }
                        sendRmrMessage(rmrMessageBuffer, message);
                        messageSent = true;
                    } else {
                        mdclog_write(MDCLOG_ERR, "RIC request id missing illigal request");
                    }
                }
                if (messageSent) {
                    break;
                }
            }
            break;
        }
        case ProcedureCode_id_RICserviceQuery: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICserviceQuery %s", message.message.enodbName);
            }
            break;
        }
        case ProcedureCode_id_RICserviceUpdate: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICserviceUpdate %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_SERVICE_UPDATE, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "RIC_SERVICE_UPDATE message failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_RICsubscription: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICsubscription %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_SUB_RESP, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "Subscription successful message failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_RICsubscriptionDelete: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICsubscriptionDelete %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_SUB_DEL_RESP, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "Subscription delete successful message failed to send to xAPP");
            }
            break;
        }
        default: {
            mdclog_write(MDCLOG_WARN, "Undefined or not supported message = %ld", procedureCode);
            message.message.messageType = 0; // no RMR message type yet
            buildJsonMessage(message);

            break;
        }
    }
}

/**
 *
 * @param pdu
 * @param message
 * @param rmrMessageBuffer
 */
void asnUnSuccsesfulMsg(E2AP_PDU_t *pdu,
                        Sctp_Map_t *sctpMap,
                        ReportingMessages_t &message,
                        RmrMessagesBuffer_t &rmrMessageBuffer) {
    auto procedureCode = pdu->choice.unsuccessfulOutcome->procedureCode;
    auto logLevel = mdclog_level_get();
    if (logLevel >= MDCLOG_INFO) {
        mdclog_write(MDCLOG_INFO, "Unsuccessful Outcome %ld", procedureCode);
    }
    switch (procedureCode) {
        case ProcedureCode_id_E2setup: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got E2setup\n");
            }
            break;
        }
        case ProcedureCode_id_ErrorIndication: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got ErrorIndication %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_ERROR_INDICATION, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "RIC_ERROR_INDICATION failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_Reset: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got Reset %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_X2_RESET, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "RIC_X2_RESET message failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_RICcontrol: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICcontrol %s", message.message.enodbName);
            }
            for (int i = 0;
                 i < pdu->choice.unsuccessfulOutcome->value.choice.RICcontrolFailure.protocolIEs.list.count; i++) {
                auto messageSent = false;
                RICcontrolFailure_IEs_t *ie = pdu->choice.unsuccessfulOutcome->value.choice.RICcontrolFailure.protocolIEs.list.array[i];
                if (logLevel >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "ie type (ProtocolIE_ID) = %ld", ie->id);
                }
                if (ie->id == ProtocolIE_ID_id_RICrequestID) {
                    if (logLevel >= MDCLOG_DEBUG) {
                        mdclog_write(MDCLOG_DEBUG, "Got RIC requestId entry, ie type (ProtocolIE_ID) = %ld", ie->id);
                    }
                    if (ie->value.present == RICcontrolFailure_IEs__value_PR_RICrequestID) {
                        message.message.messageType = rmrMessageBuffer.sendMessage->mtype = RIC_CONTROL_FAILURE;
                        rmrMessageBuffer.sendMessage->state = 0;
//                        rmrMessageBuffer.sendMessage->sub_id = (int)ie->value.choice.RICrequestID.ricRequestorID;
                        rmrMessageBuffer.sendMessage->sub_id = (int)ie->value.choice.RICrequestID.ricInstanceID;
                        static unsigned char tx[32];
                        snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
                        rmr_bytes2xact(rmrMessageBuffer.sendMessage, tx, strlen((const char *) tx));
                        rmr_bytes2meid(rmrMessageBuffer.sendMessage, (unsigned char *) message.message.enodbName,
                                       strlen(message.message.enodbName));
                        sendRmrMessage(rmrMessageBuffer, message);
                        messageSent = true;
                    } else {
                        mdclog_write(MDCLOG_ERR, "RIC request id missing illigal request");
                    }
                }
                if (messageSent) {
                    break;
                }
            }
            break;
        }
        case ProcedureCode_id_RICindication: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICindication %s", message.message.enodbName);
            }
            for (auto i = 0; i < pdu->choice.initiatingMessage->value.choice.RICindication.protocolIEs.list.count; i++) {
                auto messageSent = false;
                RICindication_IEs_t *ie = pdu->choice.initiatingMessage->value.choice.RICindication.protocolIEs.list.array[i];
                if (logLevel >= MDCLOG_DEBUG) {
                    mdclog_write(MDCLOG_DEBUG, "ie type (ProtocolIE_ID) = %ld", ie->id);
                }
                if (ie->id == ProtocolIE_ID_id_RICrequestID) {
                    if (logLevel >= MDCLOG_DEBUG) {
                        mdclog_write(MDCLOG_DEBUG, "Got RIC requestId entry, ie type (ProtocolIE_ID) = %ld", ie->id);
                    }
                    if (ie->value.present == RICindication_IEs__value_PR_RICrequestID) {
                        static unsigned char tx[32];
                        message.message.messageType = rmrMessageBuffer.sendMessage->mtype = RIC_INDICATION;
                        snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
                        rmr_bytes2xact(rmrMessageBuffer.sendMessage, tx, strlen((const char *) tx));
                        rmr_bytes2meid(rmrMessageBuffer.sendMessage,
                                       (unsigned char *)message.message.enodbName,
                                       strlen(message.message.enodbName));
                        rmrMessageBuffer.sendMessage->state = 0;
//                        rmrMessageBuffer.sendMessage->sub_id = (int)ie->value.choice.RICrequestID.ricRequestorID;
                        rmrMessageBuffer.sendMessage->sub_id = (int)ie->value.choice.RICrequestID.ricInstanceID;
                        if (mdclog_level_get() >= MDCLOG_DEBUG) {
                            mdclog_write(MDCLOG_DEBUG, "RIC sub id = %d, message type = %d",
                                         rmrMessageBuffer.sendMessage->sub_id,
                                         rmrMessageBuffer.sendMessage->mtype);
                        }
                        sendRmrMessage(rmrMessageBuffer, message);
                        messageSent = true;
                    } else {
                        mdclog_write(MDCLOG_ERR, "RIC request id missing illigal request");
                    }
                }
                if (messageSent) {
                    break;
                }
            }
            break;
        }
        case ProcedureCode_id_RICserviceQuery: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICserviceQuery %s", message.message.enodbName);
            }
            break;
        }
        case ProcedureCode_id_RICserviceUpdate: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICserviceUpdate %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_SERVICE_UPDATE, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "RIC_SERVICE_UPDATE message failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_RICsubscription: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICsubscription %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_SUB_FAILURE, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "Subscription unsuccessful message failed to send to xAPP");
            }
            break;
        }
        case ProcedureCode_id_RICsubscriptionDelete: {
            if (logLevel >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got RICsubscriptionDelete %s", message.message.enodbName);
            }
            if (sendRequestToXapp(message, RIC_SUB_DEL_FAILURE, rmrMessageBuffer) != 0) {
                mdclog_write(MDCLOG_ERR, "Subscription Delete unsuccessful message failed to send to xAPP");
            }
            break;
        }
        default: {
            mdclog_write(MDCLOG_WARN, "Undefined or not supported message = %ld", procedureCode);
            message.message.messageType = 0; // no RMR message type yet

            buildJsonMessage(message);

            break;
        }
    }
}

/**
 *
 * @param message
 * @param requestId
 * @param rmrMmessageBuffer
 * @return
 */
int sendRequestToXapp(ReportingMessages_t &message,
                      int requestId,
                      RmrMessagesBuffer_t &rmrMmessageBuffer) {
    rmr_bytes2meid(rmrMmessageBuffer.sendMessage,
                   (unsigned char *)message.message.enodbName,
                   strlen(message.message.enodbName));
    message.message.messageType = rmrMmessageBuffer.sendMessage->mtype = requestId;
    rmrMmessageBuffer.sendMessage->state = 0;
    static unsigned char tx[32];
    snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
    rmr_bytes2xact(rmrMmessageBuffer.sendMessage, tx, strlen((const char *) tx));

    auto rc = sendRmrMessage(rmrMmessageBuffer, message);
    return rc;
}


void getRmrContext(sctp_params_t &pSctpParams) {
    pSctpParams.rmrCtx = nullptr;
    pSctpParams.rmrCtx = rmr_init(pSctpParams.rmrAddress, RECEIVE_XAPP_BUFFER_SIZE, RMRFL_NONE);
    if (pSctpParams.rmrCtx == nullptr) {
        mdclog_write(MDCLOG_ERR, "Failed to initialize RMR");
        return;
    }

    rmr_set_stimeout(pSctpParams.rmrCtx, 0);    // disable retries for any send operation
    // we need to find that routing table exist and we can run
    if (mdclog_level_get() >= MDCLOG_INFO) {
        mdclog_write(MDCLOG_INFO, "We are after RMR INIT wait for RMR_Ready");
    }
    int rmrReady = 0;
    int count = 0;
    while (!rmrReady) {
        if ((rmrReady = rmr_ready(pSctpParams.rmrCtx)) == 0) {
            sleep(1);
        }
        count++;
        if (count % 60 == 0) {
            mdclog_write(MDCLOG_INFO, "waiting to RMR ready state for %d seconds", count);
        }
    }
    if (mdclog_level_get() >= MDCLOG_INFO) {
        mdclog_write(MDCLOG_INFO, "RMR running");
    }
    rmr_init_trace(pSctpParams.rmrCtx, 200);
    // get the RMR fd for the epoll
    pSctpParams.rmrListenFd = rmr_get_rcvfd(pSctpParams.rmrCtx);
    struct epoll_event event{};
    // add RMR fd to epoll
    event.events = (EPOLLIN);
    event.data.fd = pSctpParams.rmrListenFd;
    // add listening RMR FD to epoll
    if (epoll_ctl(pSctpParams.epoll_fd, EPOLL_CTL_ADD, pSctpParams.rmrListenFd, &event)) {
        mdclog_write(MDCLOG_ERR, "Failed to add RMR descriptor to epoll");
        close(pSctpParams.rmrListenFd);
        rmr_close(pSctpParams.rmrCtx);
        pSctpParams.rmrCtx = nullptr;
    }
}

int PER_FromXML(ReportingMessages_t &message, RmrMessagesBuffer_t &rmrMessageBuffer) {
    E2AP_PDU_t *pdu = nullptr;

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "got xml setup response of size %d is:%s",
                rmrMessageBuffer.rcvMessage->len, rmrMessageBuffer.rcvMessage->payload);
    }
    auto rval = asn_decode(nullptr, ATS_BASIC_XER, &asn_DEF_E2AP_PDU, (void **) &pdu,
                           rmrMessageBuffer.rcvMessage->payload, rmrMessageBuffer.rcvMessage->len);
    if (rval.code != RC_OK) {
        mdclog_write(MDCLOG_ERR, "Error %d Decoding (unpack) setup response  from E2MGR : %s",
                     rval.code,
                     message.message.enodbName);
        return -1;
    }

    int buff_size = RECEIVE_XAPP_BUFFER_SIZE;
    auto er = asn_encode_to_buffer(nullptr, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2AP_PDU, pdu,
                                   rmrMessageBuffer.rcvMessage->payload, buff_size);
    if (er.encoded == -1) {
        mdclog_write(MDCLOG_ERR, "encoding of %s failed, %s", asn_DEF_E2AP_PDU.name, strerror(errno));
        return -1;
    } else if (er.encoded > (ssize_t)buff_size) {
        mdclog_write(MDCLOG_ERR, "Buffer of size %d is to small for %s, at %s line %d",
                     (int)rmrMessageBuffer.rcvMessage->len,
                     asn_DEF_E2AP_PDU.name,
                     __func__,
                     __LINE__);
        return -1;
    }
    rmrMessageBuffer.rcvMessage->len = er.encoded;
    return 0;
}

/**
 *
 * @param sctpMap
 * @param rmrMessageBuffer
 * @param ts
 * @return
 */
int receiveXappMessages(Sctp_Map_t *sctpMap,
                        RmrMessagesBuffer_t &rmrMessageBuffer,
                        struct timespec &ts) {
    if (rmrMessageBuffer.rcvMessage == nullptr) {
        //we have error
        mdclog_write(MDCLOG_ERR, "RMR Allocation message, %s", strerror(errno));
        return -1;
    }

    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "Call to rmr_rcv_msg");
    }
    rmrMessageBuffer.rcvMessage = rmr_rcv_msg(rmrMessageBuffer.rmrCtx, rmrMessageBuffer.rcvMessage);
    if (rmrMessageBuffer.rcvMessage == nullptr) {
        mdclog_write(MDCLOG_ERR, "RMR Receving message with null pointer, Realloc rmr mesage buffer");
        rmrMessageBuffer.rcvMessage = rmr_alloc_msg(rmrMessageBuffer.rmrCtx, RECEIVE_XAPP_BUFFER_SIZE);
        return -2;
    }
    ReportingMessages_t message;
    message.message.direction = 'D';
    message.message.time.tv_nsec = ts.tv_nsec;
    message.message.time.tv_sec = ts.tv_sec;

    // get message payload
    //auto msgData = msg->payload;
    if (rmrMessageBuffer.rcvMessage->state != 0) {
        mdclog_write(MDCLOG_ERR, "RMR Receving message with stat = %d", rmrMessageBuffer.rcvMessage->state);
        return -1;
    }
    rmr_get_meid(rmrMessageBuffer.rcvMessage, (unsigned char *)message.message.enodbName);
    switch (rmrMessageBuffer.rcvMessage->mtype) {
        case RIC_E2_SETUP_RESP : {
            if (PER_FromXML(message, rmrMessageBuffer) != 0) {
                break;
            }

            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_E2_SETUP_RESP");
                return -6;
            }
            break;
        }
        case RIC_E2_SETUP_FAILURE : {
            if (PER_FromXML(message, rmrMessageBuffer) != 0) {
                break;
            }
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_E2_SETUP_FAILURE");
                return -6;
            }
            break;
        }
        case RIC_ERROR_INDICATION: {
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_ERROR_INDICATION");
                return -6;
            }
            break;
        }
        case RIC_SUB_REQ: {
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_SUB_REQ");
                return -6;
            }
            break;
        }
        case RIC_SUB_DEL_REQ: {
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_SUB_DEL_REQ");
                return -6;
            }
            break;
        }
        case RIC_CONTROL_REQ: {
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_CONTROL_REQ");
                return -6;
            }
            break;
        }
        case RIC_SERVICE_QUERY: {
            if (PER_FromXML(message, rmrMessageBuffer) != 0) {
                break;
            }
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_SERVICE_QUERY");
                return -6;
            }
            break;
        }
        case RIC_SERVICE_UPDATE_ACK: {
            if (PER_FromXML(message, rmrMessageBuffer) != 0) {
                break;
            }
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_SERVICE_UPDATE_ACK");
                return -6;
            }
            break;
        }
        case RIC_SERVICE_UPDATE_FAILURE: {
            if (PER_FromXML(message, rmrMessageBuffer) != 0) {
                break;
            }
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_SERVICE_UPDATE_FAILURE");
                return -6;
            }
            break;
        }
        case RIC_X2_RESET: {
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_X2_RESET");
                return -6;
            }
            break;
        }
        case RIC_X2_RESET_RESP: {
            if (sendDirectionalSctpMsg(rmrMessageBuffer, message, 0, sctpMap) != 0) {
                mdclog_write(MDCLOG_ERR, "Failed to send RIC_X2_RESET_RESP");
                return -6;
            }
            break;
        }
        case RIC_SCTP_CLEAR_ALL: {
            mdclog_write(MDCLOG_INFO, "RIC_SCTP_CLEAR_ALL");
            // loop on all keys and close socket and then erase all map.
            vector<char *> v;
            sctpMap->getKeys(v);
            for (auto const &iter : v) { //}; iter != sctpMap.end(); iter++) {
                if (!boost::starts_with((string) (iter), "host:") && !boost::starts_with((string) (iter), "msg:")) {
                    auto *peerInfo = (ConnectedCU_t *) sctpMap->find(iter);
                    if (peerInfo == nullptr) {
                        continue;
                    }
                    close(peerInfo->fileDescriptor);
                    memcpy(message.message.enodbName, peerInfo->enodbName, sizeof(peerInfo->enodbName));
                    message.message.direction = 'D';
                    message.message.time.tv_nsec = ts.tv_nsec;
                    message.message.time.tv_sec = ts.tv_sec;

                    message.message.asnLength = rmrMessageBuffer.sendMessage->len =
                            snprintf((char *)rmrMessageBuffer.sendMessage->payload,
                                     256,
                                     "%s|RIC_SCTP_CLEAR_ALL",
                                     peerInfo->enodbName);
                    message.message.asndata = rmrMessageBuffer.sendMessage->payload;
                    mdclog_write(MDCLOG_INFO, "%s", message.message.asndata);
                    if (sendRequestToXapp(message, RIC_SCTP_CONNECTION_FAILURE, rmrMessageBuffer) != 0) {
                        mdclog_write(MDCLOG_ERR, "SCTP_CONNECTION_FAIL message failed to send to xAPP");
                    }
                    free(peerInfo);
                }
            }

            sleep(1);
            sctpMap->clear();
            break;
        }
        case E2_TERM_KEEP_ALIVE_REQ: {
            // send message back
            rmr_bytes2payload(rmrMessageBuffer.sendMessage,
                              (unsigned char *)rmrMessageBuffer.ka_message,
                              rmrMessageBuffer.ka_message_len);
            rmrMessageBuffer.sendMessage->mtype = E2_TERM_KEEP_ALIVE_RESP;
            rmrMessageBuffer.sendMessage->state = 0;
            static unsigned char tx[32];
            auto txLen = snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
            rmr_bytes2xact(rmrMessageBuffer.sendMessage, tx, txLen);
            rmrMessageBuffer.sendMessage = rmr_send_msg(rmrMessageBuffer.rmrCtx, rmrMessageBuffer.sendMessage);
            if (rmrMessageBuffer.sendMessage == nullptr) {
                rmrMessageBuffer.sendMessage = rmr_alloc_msg(rmrMessageBuffer.rmrCtx, RECEIVE_XAPP_BUFFER_SIZE);
                mdclog_write(MDCLOG_ERR, "Failed to send E2_TERM_KEEP_ALIVE_RESP RMR message returned NULL");
            } else if (rmrMessageBuffer.sendMessage->state != 0)  {
                mdclog_write(MDCLOG_ERR, "Failed to send E2_TERM_KEEP_ALIVE_RESP, on RMR state = %d ( %s)",
                             rmrMessageBuffer.sendMessage->state, translateRmrErrorMessages(rmrMessageBuffer.sendMessage->state).c_str());
            } else if (mdclog_level_get() >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "Got Keep Alive Request send : %s", rmrMessageBuffer.ka_message);
            }

            break;
        }
        case 42000: {
            mdclog_write(MDCLOG_INFO, "Received message 42000");

            // sendMessageSocket(1234);
            break;
        }
        default:
            mdclog_write(MDCLOG_WARN, "Message Type : %d is not seported", rmrMessageBuffer.rcvMessage->mtype);
            message.message.asndata = rmrMessageBuffer.rcvMessage->payload;
            message.message.asnLength = rmrMessageBuffer.rcvMessage->len;
            message.message.time.tv_nsec = ts.tv_nsec;
            message.message.time.tv_sec = ts.tv_sec;
            message.message.messageType = rmrMessageBuffer.rcvMessage->mtype;

            buildJsonMessage(message);


            return -7;
    }
    if (mdclog_level_get() >= MDCLOG_DEBUG) {
        mdclog_write(MDCLOG_DEBUG, "EXIT OK from %s", __FUNCTION__);
    }
    return 0;
}

int sendMessageSocket(const int dest_port) {

    const char* dest_ip = "10.0.2.100";

    int sckfd = socket(AF_INET, SOCK_STREAM, IPPROTO_SCTP);
    if (sckfd < 0) {
        mdclog_write(MDCLOG_ERR, "ERROR: OPEN SOCKET");
        close(sckfd);
        return -1;
    }

    // SET SOCKET OPTIONS TO RELEASE THE SOCKET ADDRESS IMMEDIATELY AFTER
    // THE SOCKET IS CLOSED
    int option(1);
    setsockopt(sckfd, SOL_SOCKET, SO_REUSEADDR, (char*)&option, sizeof(option));

    struct sockaddr_in dest_addr = {0};
    dest_addr.sin_family = AF_INET;
    dest_addr.sin_port = htons(dest_port);

    // convert dest_ip from char* to network address
    if (inet_pton(AF_INET, dest_ip, &dest_addr.sin_addr) <= 0) {
        mdclog_write(MDCLOG_ERR, "ERROR CONVERTING IP TO INTERNET ADDR");
        close(sckfd); // if conversion fail, close the socket and return error -2
        return -2;
    }

    if (connect(sckfd, (struct sockaddr *) &dest_addr, sizeof(dest_addr)) < 0) {
        mdclog_write(MDCLOG_ERR, "ERROR: CONNECT");
        close(sckfd);
        return -3;
    }

    // TRANSMIT DATA
    const size_t max_size = 512;
    char buf[max_size] = "Hello, Server!";  // store the data in a buffer
    size_t data_size = 14;
    int sent_size = sctp_sendmsg(sckfd, buf, data_size, NULL, 0, 0, 0, 0, 0, 0 );
    // int sent_size = sctp_sendmsg(sckfd, (void *) msg, strlen(msg) + 1, NULL, 0, 0, 0, 0, 0, 0 );

    if(sent_size < 0) { // the send returns a size of -1 in case of errors
        mdclog_write(MDCLOG_ERR, "ERROR: SEND");
        close(sckfd); // if error close the socket and exit
        return -4;
    }
    else {
        mdclog_write(MDCLOG_INFO, "Message sent");
    }

    close(sckfd);

    return 0;

    // // open a SOCK_STREAM (TCP) socket
    // int sckfd = socket(AF_INET, SOCK_STREAM, 0);
    // if (sckfd < 0){
    //     mdclog_write(MDCLOG_ERR, "ERROR: OPEN SOCKET");
    //     close(sckfd);
    //     return -1;
    // }
    //
    // // SET SOCKET OPTIONS TO RELEASE THE SOCKET ADDRESS IMMEDIATELY AFTER
    // // THE SOCKET IS CLOSED
    // int option(1);
    // setsockopt(sckfd, SOL_SOCKET, SO_REUSEADDR, (char*)&option, sizeof(option));
    //
    // //SET SERVER DESTINATION ADDRESS
    // struct sockaddr_in dest_addr = {0}; // set all elements of the struct to 0
    // dest_addr.sin_family = AF_INET; // address family is AF_INET (IPV4)
    //
    // // convert dest_port to network number format
    // dest_addr.sin_port = htons(dest_port);
    // // convert dest_ip from char* to network address
    // if (inet_pton(AF_INET, dest_ip, &dest_addr.sin_addr) <= 0) {
    //     mdclog_write(MDCLOG_ERR, "ERROR CONVERTING IP TO INTERNET ADDR");
    //     close(sckfd); // if conversion fail, close the socket and return error -2
    //     return -2;
    // }
    //
    // //CONNECT THE SOCKET TO THE SERVER IP:PORT SPECIFIED INTO DEST_ADDR
    // if (connect(sckfd, (struct sockaddr*) &dest_addr, sizeof(dest_addr)) < 0) {
    //     mdclog_write(MDCLOG_ERR, "ERROR: CONNECT");
    //     close(sckfd); // if connection failed return
    //     return -3;
    // }
    //
    // //TRANSMIT DATA
    // const size_t max_size = 512;
    // char buf[max_size] = "Hello from e2term"; // store the data in a buffer
    // size_t data_size = 10;
    // int sent_size = send(sckfd,buf,data_size,0); // send the data through sckfd
    //
    // if(sent_size < 0) { // the send returns a size of -1 in case of errors
    //     mdclog_write(MDCLOG_ERR, "ERROR: SEND");
    //     close(sckfd); // if error close the socket and exit
    //     return -4;
    // }
    // else {
    //     mdclog_write(MDCLOG_INFO, "Message sent");
    // }
    //
    // memset(buf, 0, max_size); // set buffer to zero for next read
    //
    // //CLOSE THE SOCKET
    // close(sckfd);
}

/**
 * Send message to the CU that is not expecting for successful or unsuccessful results
 * @param messageBuffer
 * @param message
 * @param failedMsgId
 * @param sctpMap
 * @return
 */
int sendDirectionalSctpMsg(RmrMessagesBuffer_t &messageBuffer,
                           ReportingMessages_t &message,
                           int failedMsgId,
                           Sctp_Map_t *sctpMap) {

    getRequestMetaData(message, messageBuffer);
    if (mdclog_level_get() >= MDCLOG_INFO) {
        mdclog_write(MDCLOG_INFO, "send message to %s address", message.message.enodbName);
    }

    auto rc = sendMessagetoCu(sctpMap, messageBuffer, message, failedMsgId);
    return rc;
}

/**
 *
 * @param sctpMap
 * @param messageBuffer
 * @param message
 * @param failedMesgId
 * @return
 */
int sendMessagetoCu(Sctp_Map_t *sctpMap,
                    RmrMessagesBuffer_t &messageBuffer,
                    ReportingMessages_t &message,
                    int failedMesgId) {
    auto *peerInfo = (ConnectedCU_t *) sctpMap->find(message.message.enodbName);
    if (peerInfo == nullptr) {
        if (failedMesgId != 0) {
            sendFailedSendingMessagetoXapp(messageBuffer, message, failedMesgId);
        } else {
            mdclog_write(MDCLOG_ERR, "Failed to send message no CU entry %s", message.message.enodbName);
        }
        return -1;
    }

    // get the FD
    message.message.messageType = messageBuffer.rcvMessage->mtype;
    auto rc = sendSctpMsg(peerInfo, message, sctpMap);
    return rc;
}

/**
 *
 * @param rmrCtx the rmr context to send and receive
 * @param msg the msg we got fromxApp
 * @param metaData data from xApp in ordered struct
 * @param failedMesgId the return message type error
 */
void
sendFailedSendingMessagetoXapp(RmrMessagesBuffer_t &rmrMessageBuffer, ReportingMessages_t &message, int failedMesgId) {
    rmr_mbuf_t *msg = rmrMessageBuffer.sendMessage;
    msg->len = snprintf((char *) msg->payload, 200, "the gNb/eNode name %s not found",
                        message.message.enodbName);
    if (mdclog_level_get() >= MDCLOG_INFO) {
        mdclog_write(MDCLOG_INFO, "%s", msg->payload);
    }
    msg->mtype = failedMesgId;
    msg->state = 0;

    static unsigned char tx[32];
    snprintf((char *) tx, sizeof tx, "%15ld", transactionCounter++);
    rmr_bytes2xact(msg, tx, strlen((const char *) tx));

    sendRmrMessage(rmrMessageBuffer, message);
}



/**
 *
 * @param epoll_fd
 * @param peerInfo
 * @param events
 * @param sctpMap
 * @param enodbName
 * @param msgType
 * @return
 */
int addToEpoll(int epoll_fd,
               ConnectedCU_t *peerInfo,
               uint32_t events,
               Sctp_Map_t *sctpMap,
               char *enodbName,
               int msgType) {
    // Add to Epol
    struct epoll_event event{};
    event.data.ptr = peerInfo;
    event.events = events;
    if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, peerInfo->fileDescriptor, &event) < 0) {
        if (mdclog_level_get() >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG, "epoll_ctl EPOLL_CTL_ADD (may chack not to quit here), %s, %s %d",
                         strerror(errno), __func__, __LINE__);
        }
        close(peerInfo->fileDescriptor);
        if (enodbName != nullptr) {
            cleanHashEntry(peerInfo, sctpMap);
            char key[MAX_ENODB_NAME_SIZE * 2];
            snprintf(key, MAX_ENODB_NAME_SIZE * 2, "msg:%s|%d", enodbName, msgType);
            if (mdclog_level_get() >= MDCLOG_DEBUG) {
                mdclog_write(MDCLOG_DEBUG, "remove key = %s from %s at line %d", key, __FUNCTION__, __LINE__);
            }
            auto tmp = sctpMap->find(key);
            if (tmp) {
                free(tmp);
                sctpMap->erase(key);
            }
        } else {
            peerInfo->enodbName[0] = 0;
        }
        mdclog_write(MDCLOG_ERR, "epoll_ctl EPOLL_CTL_ADD (may chack not to quit here)");
        return -1;
    }
    return 0;
}

/**
 *
 * @param epoll_fd
 * @param peerInfo
 * @param events
 * @param sctpMap
 * @param enodbName
 * @param msgType
 * @return
 */
int modifyToEpoll(int epoll_fd,
                  ConnectedCU_t *peerInfo,
                  uint32_t events,
                  Sctp_Map_t *sctpMap,
                  char *enodbName,
                  int msgType) {
    // Add to Epol
    struct epoll_event event{};
    event.data.ptr = peerInfo;
    event.events = events;
    if (epoll_ctl(epoll_fd, EPOLL_CTL_MOD, peerInfo->fileDescriptor, &event) < 0) {
        if (mdclog_level_get() >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG, "epoll_ctl EPOLL_CTL_MOD (may chack not to quit here), %s, %s %d",
                         strerror(errno), __func__, __LINE__);
        }
        close(peerInfo->fileDescriptor);
        cleanHashEntry(peerInfo, sctpMap);
        char key[MAX_ENODB_NAME_SIZE * 2];
        snprintf(key, MAX_ENODB_NAME_SIZE * 2, "msg:%s|%d", enodbName, msgType);
        if (mdclog_level_get() >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG, "remove key = %s from %s at line %d", key, __FUNCTION__, __LINE__);
        }
        auto tmp = sctpMap->find(key);
        if (tmp) {
            free(tmp);
        }
        sctpMap->erase(key);
        mdclog_write(MDCLOG_ERR, "epoll_ctl EPOLL_CTL_ADD (may chack not to quit here)");
        return -1;
    }
    return 0;
}


int sendRmrMessage(RmrMessagesBuffer_t &rmrMessageBuffer, ReportingMessages_t &message) {
    buildJsonMessage(message);

    rmrMessageBuffer.sendMessage = rmr_send_msg(rmrMessageBuffer.rmrCtx, rmrMessageBuffer.sendMessage);

    if (rmrMessageBuffer.sendMessage == nullptr) {
        rmrMessageBuffer.sendMessage = rmr_alloc_msg(rmrMessageBuffer.rmrCtx, RECEIVE_XAPP_BUFFER_SIZE);
        mdclog_write(MDCLOG_ERR, "RMR failed send message returned with NULL pointer");
        return -1;
    }

    if (rmrMessageBuffer.sendMessage->state != 0) {
        char meid[RMR_MAX_MEID]{};
        if (rmrMessageBuffer.sendMessage->state == RMR_ERR_RETRY) {
            usleep(5);
            rmrMessageBuffer.sendMessage->state = 0;
            mdclog_write(MDCLOG_INFO, "RETRY sending Message type %d to Xapp from %s",
                         rmrMessageBuffer.sendMessage->mtype,
                         rmr_get_meid(rmrMessageBuffer.sendMessage, (unsigned char *)meid));
            rmrMessageBuffer.sendMessage = rmr_send_msg(rmrMessageBuffer.rmrCtx, rmrMessageBuffer.sendMessage);
            if (rmrMessageBuffer.sendMessage == nullptr) {
                mdclog_write(MDCLOG_ERR, "RMR failed send message returned with NULL pointer");
                rmrMessageBuffer.sendMessage = rmr_alloc_msg(rmrMessageBuffer.rmrCtx, RECEIVE_XAPP_BUFFER_SIZE);
                return -1;
            } else if (rmrMessageBuffer.sendMessage->state != 0) {
                mdclog_write(MDCLOG_ERR,
                             "Message state %s while sending request %d to Xapp from %s after retry of 10 microseconds",
                             translateRmrErrorMessages(rmrMessageBuffer.sendMessage->state).c_str(),
                             rmrMessageBuffer.sendMessage->mtype,
                             rmr_get_meid(rmrMessageBuffer.sendMessage, (unsigned char *)meid));
                auto rc = rmrMessageBuffer.sendMessage->state;
                return rc;
            }
        } else {
            mdclog_write(MDCLOG_ERR, "Message state %s while sending request %d to Xapp from %s",
                         translateRmrErrorMessages(rmrMessageBuffer.sendMessage->state).c_str(),
                         rmrMessageBuffer.sendMessage->mtype,
                         rmr_get_meid(rmrMessageBuffer.sendMessage, (unsigned char *)meid));
            return rmrMessageBuffer.sendMessage->state;
        }
    }
    return 0;
}

void buildJsonMessage(ReportingMessages_t &message) {
    if (jsonTrace) {
        message.outLen = sizeof(message.base64Data);
        base64::encode((const unsigned char *) message.message.asndata,
                       (const int) message.message.asnLength,
                       message.base64Data,
                       message.outLen);
        if (mdclog_level_get() >= MDCLOG_DEBUG) {
            mdclog_write(MDCLOG_DEBUG, "Tracing: ASN length = %d, base64 message length = %d ",
                         (int) message.message.asnLength,
                         (int) message.outLen);
        }

        snprintf(message.buffer, sizeof(message.buffer),
                 "{\"header\": {\"ts\": \"%ld.%09ld\","
                 "\"ranName\": \"%s\","
                 "\"messageType\": %d,"
                 "\"direction\": \"%c\"},"
                 "\"base64Length\": %d,"
                 "\"asnBase64\": \"%s\"}",
                 message.message.time.tv_sec,
                 message.message.time.tv_nsec,
                 message.message.enodbName,
                 message.message.messageType,
                 message.message.direction,
                 (int) message.outLen,
                 message.base64Data);
        static src::logger_mt &lg = my_logger::get();

        BOOST_LOG(lg) << message.buffer;
    }
}


/**
 * take RMR error code to string
 * @param state
 * @return
 */
string translateRmrErrorMessages(int state) {
    string str = {};
    switch (state) {
        case RMR_OK:
            str = "RMR_OK - state is good";
            break;
        case RMR_ERR_BADARG:
            str = "RMR_ERR_BADARG - argument passd to function was unusable";
            break;
        case RMR_ERR_NOENDPT:
            str = "RMR_ERR_NOENDPT - send//call could not find an endpoint based on msg type";
            break;
        case RMR_ERR_EMPTY:
            str = "RMR_ERR_EMPTY - msg received had no payload; attempt to send an empty message";
            break;
        case RMR_ERR_NOHDR:
            str = "RMR_ERR_NOHDR - message didn't contain a valid header";
            break;
        case RMR_ERR_SENDFAILED:
            str = "RMR_ERR_SENDFAILED - send failed; errno has nano reason";
            break;
        case RMR_ERR_CALLFAILED:
            str = "RMR_ERR_CALLFAILED - unable to send call() message";
            break;
        case RMR_ERR_NOWHOPEN:
            str = "RMR_ERR_NOWHOPEN - no wormholes are open";
            break;
        case RMR_ERR_WHID:
            str = "RMR_ERR_WHID - wormhole id was invalid";
            break;
        case RMR_ERR_OVERFLOW:
            str = "RMR_ERR_OVERFLOW - operation would have busted through a buffer/field size";
            break;
        case RMR_ERR_RETRY:
            str = "RMR_ERR_RETRY - request (send/call/rts) failed, but caller should retry (EAGAIN for wrappers)";
            break;
        case RMR_ERR_RCVFAILED:
            str = "RMR_ERR_RCVFAILED - receive failed (hard error)";
            break;
        case RMR_ERR_TIMEOUT:
            str = "RMR_ERR_TIMEOUT - message processing call timed out";
            break;
        case RMR_ERR_UNSET:
            str = "RMR_ERR_UNSET - the message hasn't been populated with a transport buffer";
            break;
        case RMR_ERR_TRUNC:
            str = "RMR_ERR_TRUNC - received message likely truncated";
            break;
        case RMR_ERR_INITFAILED:
            str = "RMR_ERR_INITFAILED - initialisation of something (probably message) failed";
            break;
        case RMR_ERR_NOTSUPP:
            str = "RMR_ERR_NOTSUPP - the request is not supported, or RMr was not initialised for the request";
            break;
        default:
            char buf[128]{};
            snprintf(buf, sizeof buf, "UNDOCUMENTED RMR_ERR : %d", state);
            str = buf;
            break;
    }
    return str;
}
