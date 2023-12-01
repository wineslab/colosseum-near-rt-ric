##############################################################################
#
#   Copyright (c) 2019 AT&T Intellectual Property.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
##############################################################################

#
#   This source code is part of the near-RT RIC (RAN Intelligent Controller)
#   platform project (RICP).
#
FROM wineslab/o-ran-sc-bldr-ubuntu18-c-go:9-u18.04 as ubuntu

WORKDIR /opt/e2/

ARG BUILD_TYPE="Release"

RUN mkdir -p /opt/e2/RIC-E2-TERMINATION/ \
#    && mkdir -p /opt/e2/RIC-E2-TERMINATION/TEST/T1 \
    && mkdir -p /opt/e2/RIC-E2-TERMINATION/TEST/T2 \
    && mkdir -p /opt/e2/RIC-E2-TERMINATION/3rdparty

COPY . /opt/e2/RIC-E2-TERMINATION/

RUN apt-get update && apt-get install -y libcurl4-gnutls-dev gawk libtbb-dev libtbb-doc libtbb2 libtbb2-dbg rapidjson-dev && rm -rf /var/lib/apt/lists/* \
    && mv /opt/e2/RIC-E2-TERMINATION/CMakeLists.txt /opt/e2/ && mv /opt/e2/RIC-E2-TERMINATION/Makefile /opt/e2/ && cat /opt/e2/RIC-E2-TERMINATION/config/config.conf \
    && wget --content-disposition https://github.com/cgreen-devs/cgreen/releases/download/1.2.0/cgreen_1.2.0_amd64.deb \
    && dpkg -i cgreen_1.2.0_amd64.deb \
    && wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr_4.0.2_amd64.deb/download.deb \
    && dpkg -i rmr_4.0.2_amd64.deb \
    && wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr-dev_4.0.2_amd64.deb/download.deb \
    && dpkg -i rmr-dev_4.0.2_amd64.deb \
    && wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/mdclog_0.0.4-1_amd64.deb/download.deb \
    && dpkg -i mdclog_0.0.4-1_amd64.deb \
    && wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/mdclog-dev_0.0.4-1_amd64.deb/download.deb \
    && dpkg -i mdclog-dev_0.0.4-1_amd64.deb \
    && cd /opt/e2/RIC-E2-TERMINATION/3rdparty && git clone https://github.com/oktal/pistache.git && cd pistache && git checkout 270bbefeb25a402153a55053f845e9c7674ab713 \
    && git submodule update --init && mkdir build && cd build \
    && cmake -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release -DPISTACHE_BUILD_DOCS=false \
       -DPISTACHE_USE_SSL=false -DPISTACHE_BUILD_TESTS=false -DCMAKE_INSTALL_PREFIX=/usr/local ../ \
    && make -j ${nproc} && make install && ldconfig \
#    && git clone http://gerrit.o-ran-sc.org/r/com/log \
#    && cd log && ./autogen.sh && ./configure && make && make install && ldconfig \
    && cd /opt/e2/RIC-E2-TERMINATION/3rdparty && git clone https://github.com/jarro2783/cxxopts.git \
    && cd /opt/e2/RIC-E2-TERMINATION/3rdparty && git clone https://github.com/Tencent/rapidjson.git \
    && cd /opt/e2/RIC-E2-TERMINATION/3rdparty && git clone https://github.com/zeux/pugixml.git \
    && cd /opt/e2/ && git clone https://github.com/bilke/cmake-modules.git \
    && cd /opt/e2/ && /usr/local/bin/cmake -D CMAKE_BUILD_TYPE=$BUILD_TYPE . && make -j ${nproc} \
    && echo "3" > /opt/e2/rmr.verbose



RUN if [$BUILD_TYPE == "Debug"] ; then make e2_coverage ; fi

#    && git clone http://gerrit.o-ran-sc.org/r/ric-plt/tracelibcpp \
#    && cd tracelibcpp && mkdir build && cd build  \
#    &&  sed -i '19iset\(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -O3"\)' ../CMakeLists.txt \
#    &&  sed -i '19iset\(CMAKE_CXX_STANDARD 17\)' ../CMakeLists.txt\
#    && /usr/local/bin/cmake .. && make && cd .. && cp include/tracelibcpp/tracelibcpp.hpp . && cd .. \
#    && mkdir nlohmann && cd nlohmann && find / -type f -name "json.hpp" -exec cp {} json.hpp \; && cd .. \
#    && find / -type d -name "opentracing" -exec cp -r {} . \; \
#    && cd /usr/local/lib/ && find / -type f -name "libyaml-cpp.a" -exec cp {} libyaml-cpp.a \; \
#    && find / -type f -name "libopentracing.a" -exec cp {} libopentracing.a \; && cd /opt/e2/RIC-E2-TERMINATION && ls nlohmann  \

FROM ubuntu:18.04
RUN apt-get update && apt-get install -y net-tools iputils-ping curl tcpdump libsctp-dev vim build-essential && mkdir -p /opt/e2/config

# COPY --from=ubuntu /opt/e2/e2 /opt/e2/e2
COPY --from=ubuntu /opt/e2/ /opt/e2/
COPY --from=ubuntu /opt/e2/rmr.verbose /tmp/rmr.verbose
COPY --from=ubuntu /opt/e2/RIC-E2-TERMINATION/dockerRouter.txt /opt/e2/RIC-E2-TERMINATION/startup.sh /opt/e2/
COPY --from=ubuntu /opt/e2/RIC-E2-TERMINATION/config/config.conf /opt/e2/config/config.conf
#COPY --from=ubuntu /opt/e2/RIC-E2-TERMINATION/tracelibcpp/build/libtracelibcpp.so /usr/local/lib/libtracelibcpp.so
#COPY --from=ubuntu /opt/e2/RIC-E2-TERMINATION/tracelibcpp/build/libtracelibcpp.so.0 /usr/local/lib/libtracelibcpp.so.0
#COPY --from=ubuntu /opt/e2/RIC-E2-TERMINATION/tracelibcpp/build/libtracelibcpp.so.0.0.2 /usr/local/lib/libtracelibcpp.so.0.0.2

COPY --from=ubuntu /usr/local/lib/librmr_si.so.4.0.2 \
                   /usr/local/lib/librmr_si.so.4 \
                   /usr/local/lib/librmr_si.so \
                   /usr/lib/x86_64-linux-gnu/libmdclog.so \
                   /usr/lib/x86_64-linux-gnu/libmdclog.so.0 \
                   /usr/lib/x86_64-linux-gnu/libmdclog.so.0.0.4 \
                   /usr/local/lib/

COPY --from=ubuntu /usr/lib/x86_64-linux-gnu/libcgreen.so \
                   /usr/lib/x86_64-linux-gnu/libcgreen.so.1 \
                   /usr/lib/x86_64-linux-gnu/libcgreen.so.1.2.0 \
                   /usr/lib/x86_64-linux-gnu/libtbb.so \
                   /usr/lib/x86_64-linux-gnu/libtbb.so.2 \
                   /usr/lib/x86_64-linux-gnu/libtbbmalloc.so \
                   /usr/lib/x86_64-linux-gnu/libtbbmalloc.so.2 \
                   /usr/lib/x86_64-linux-gnu/libtbbmalloc_proxy.so \
                   /usr/lib/x86_64-linux-gnu/libtbbmalloc_proxy.so.2 \
                   /usr/local/

COPY --from=ubuntu /usr/lib/libboost_log_setup.so \
                   /usr/lib/libboost_log_setup.so.1.69.0 \
                   /usr/lib/libboost_log.so \
                   /usr/lib/libboost_log.so.1.69.0 \
                   /usr/lib/libboost_thread.so \
                   /usr/lib/libboost_thread.so.1.69.0 \
                   /usr/lib/libboost_filesystem.so \
                   /usr/lib/libboost_filesystem.so.1.69.0 \
                   /usr/lib/libboost_regex.so \
                   /usr/lib/libboost_regex.so.1.69.0 \
                   /usr/lib/x86_64-linux-gnu/libicui18n.so \
                   /usr/lib/x86_64-linux-gnu/libicui18n.so.60.2 \
                   /usr/lib/x86_64-linux-gnu/libicui18n.so.60 \
                   /usr/lib/x86_64-linux-gnu/libicuuc.so \
                   /usr/lib/x86_64-linux-gnu/libicuuc.so.60.2 \
                   /usr/lib/x86_64-linux-gnu/libicuuc.so.60 \
                   /usr/lib/x86_64-linux-gnu/libicudata.so \
                   /usr/lib/x86_64-linux-gnu/libicudata.so.60.2 \
                   /usr/lib/x86_64-linux-gnu/libicudata.so.60 \
                   /usr/lib/x86_64-linux-gnu/libstdc++.so.6 \
                   /usr/lib/x86_64-linux-gnu/libsctp.so \
                   /usr/lib/x86_64-linux-gnu/libsctp.so.1 \
                   /usr/lib/x86_64-linux-gnu/libsctp.so.1.0.17 \
                   /usr/lib/

COPY --from=ubuntu /usr/local/bin/cmake /usr/local/bin/cmake
COPY --from=ubuntu /usr/local/share/cmake-3.14 /usr/local/share/cmake-3.14

# the following are needed to compile e2term inside the Docker container
RUN apt-get install -y libboost-all-dev libtbb-dev rapidjson-dev && rm -rf /var/lib/apt/lists/*
RUN cd /opt/e2/RIC-E2-TERMINATION/3rdparty/pistache/build && make install
RUN dpkg -i /opt/e2/cgreen_1.2.0_amd64.deb && dpkg -i /opt/e2/rmr_4.0.2_amd64.deb && \
    dpkg -i /opt/e2/rmr-dev_4.0.2_amd64.deb && dpkg -i /opt/e2/mdclog_0.0.4-1_amd64.deb && \
    dpkg -i /opt/e2/mdclog-dev_0.0.4-1_amd64.deb

WORKDIR /opt/e2/
ENV LD_LIBRARY_PATH=/usr/local/lib
ENV RMR_SEED_RT=dockerRouter.txt
ENV E2TERM_POD_NAME=e2term
RUN chmod +x /opt/e2/startup.sh
EXPOSE 38000
CMD ["sh", "-c", "./startup.sh"]
