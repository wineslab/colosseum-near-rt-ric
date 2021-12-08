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

FROM snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/e2_base:1.0.0 as ubuntu

WORKDIR /opt/e2/
RUN mkdir -p /opt/e2/RIC-E2-TERMINATION/ \
    && mkdir -p /opt/e2/RIC-E2-TERMINATION/TEST/T1 \
    && mkdir -p /opt/e2/RIC-E2-TERMINATION/TEST/T2
COPY . /opt/e2/RIC-E2-TERMINATION/
RUN mv /opt/e2/RIC-E2-TERMINATION/CMakeLists.txt /opt/e2/

RUN echo "137.135.91.204 gerrit.ranco-dev-tools.eastus.cloudapp.azure.com" >> /etc/hosts \
    && apt-get install -y autoconf gawk libtool automake pkg-config autoconf-archive \
    && git clone http://gerrit.ranco-dev-tools.eastus.cloudapp.azure.com/com/log \
    && cd log \
    && ./autogen.sh && ./configure && make && make install && ldconfig

RUN  git clone https://gerrit.o-ran-sc.org/r/ric-plt/lib/rmr \
     && cd rmr/; mkdir build; cd build; /opt/bin/cmake -DDEV_PKG=1 ..; make install \
     && cd /opt/e2/ && /opt/bin/cmake . && make

FROM ubuntu:16.04
COPY --from=ubuntu /opt/e2/setUpTest /opt/e2/setUpTest
COPY --from=ubuntu /opt/e2/RIC-E2-TERMINATION/TEST/T1/dockerRouter.txt /opt/e2/dockerRouter.txt
COPY --from=ubuntu /usr/local/lib/librmr_nng.so.1 /usr/local/lib/librmr_nng.so.1
COPY --from=ubuntu /usr/local/lib/libnng.so.1 /usr/local/lib/libnng.so.1
COPY --from=ubuntu /usr/local/lib/libmdclog.so.0 /usr/local/lib/libmdclog.so.0

WORKDIR /opt/e2/
ENV LD_LIBRARY_PATH=/usr/local/lib
#ENV RMR_RTG_SVC=127.0.0.1
ENV RMR_SEED_RT=dockerRouter.txt
ENV host=127.0.0.1
ENV port=5566
ENV ran=ranname
ENV rmr=38012
ENV loglevel=info
EXPOSE 5566
EXPOSE 38012
CMD ["sh", "-c", "./setUpTest host $host port $port ran $ran rmr $rmr loglevel $loglevel" ]
