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

FROM nexus3.o-ran-sc.org:10002/o-ran-sc/bldr-ubuntu18-c-go:1.9.0 as ubuntu

WORKDIR /opt/E2Manager
COPY . .

# Install dependencies, compile and test the module
RUN apt-get update \
  && apt-get install -y sudo \
  && apt-get --allow-releaseinfo-change update \
  && apt upgrade -y
RUN bash build-e2mgr-ubuntu.sh

# Build deployable container
FROM ubuntu:18.04

RUN apt-get update && apt-get install -y \
  net-tools \
  iputils-ping \
  curl \
  tcpdump \
  jq \
  && rm -rf /var/lib/apt/lists/*

COPY --from=ubuntu /opt/E2Manager/router.txt /opt/E2Manager/main /opt/E2Manager/
COPY --from=ubuntu /opt/E2Manager/resources/configuration.yaml /opt/E2Manager/resources/configuration.yaml
COPY --from=ubuntu /opt/E2Manager/resources/rmr.verbose /tmp/rmr.verbose
COPY --from=ubuntu /usr/local/lib/librmr_si.so.4 /usr/local/lib/librmr_si.so.4

WORKDIR /opt/E2Manager
ENV LD_LIBRARY_PATH=/usr/local/lib \
    port=3800
ENV RMR_SEED_RT=router.txt
ENV RMR_VCTL_FILE=/tmp/rmr.verbose
EXPOSE 3800
CMD ["sh", "-c", "./main  -port=$port"]
