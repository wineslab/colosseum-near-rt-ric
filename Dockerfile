#==================================================================================
#	Copyright (c) 2022 Northeastern University
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#	   http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#==================================================================================

FROM wineslab/o-ran-sc-bldr-ubuntu18-c-go:9-u18.04 as buildenv
ARG log_level_e2sim=2
# log_level_e2sim = 0 ->  LOG_LEVEL_UNCOND   0
# log_level_e2sim = 1 -> LOG_LEVEL_ERROR     1
# log_level_e2sim = 2 -> LOG_LEVEL_INFO      2
# log_level_e2sim = 3 -> LOG_LEVEL_DEBUG     3

# Install E2sim
RUN mkdir -p /workspace/e2sim
RUN apt-get update && apt-get install -y build-essential git cmake libsctp-dev autoconf automake libtool bison flex libboost-all-dev

WORKDIR /workspace/e2sim

COPY ./e2sim/e2sim /workspace/e2sim

RUN mkdir /workspace/e2sim/build
WORKDIR /workspace/e2sim/build

RUN cmake .. -DDEV_PKG=1 -DLOG_LEVEL=${log_level_e2sim}
RUN make package
RUN echo "Going to install e2sim-dev"
RUN dpkg --install ./e2sim-dev_1.0.0_amd64.deb
RUN ldconfig

WORKDIR /workspace

# Install ns-3
RUN apt-get install -y g++ python3 qtbase5-dev qtchooser qt5-qmake qtbase5-dev-tools

COPY ./ns3-mmwave-oran /workspace/ns3-mmwave-oran
COPY ./ns-o-ran /workspace/ns3-mmwave-oran/contrib/oran-interface

WORKDIR /workspace/ns3-mmwave-oran

RUN ./waf configure --enable-tests --enable-examples
# RUN ./waf build

WORKDIR /workspace

CMD [ "/bin/sh" ]


