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

FROM nexus3.o-ran-sc.org:10004/bldr-ubuntu16-c-go:2-u16.04-nng as ubuntu

WORKDIR /opt/xappmock
COPY . . 
ENV PATH=$PATH:/usr/local/go/bin:/usr/lib/go-1.12/bin
# Install RMr library and dev files
RUN wget --content-disposition  https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr_1.10.0_amd64.deb/download.deb
RUN dpkg -i rmr_1.10.0_amd64.deb
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr-dev_1.10.0_amd64.deb/download.deb
RUN dpkg -i rmr-dev_1.10.0_amd64.deb

RUN go build main/xappmock.go
    


FROM ubuntu:16.04
COPY --from=ubuntu /opt/xappmock/xappmock /opt/xappmock/xappmock
COPY --from=ubuntu /opt/xappmock/resources /opt/xappmock/resources
COPY --from=ubuntu /usr/local/lib/librmr_nng.so.1 /usr/local/lib/librmr_nng.so.1
COPY --from=ubuntu /usr/local/lib/libnng.so.1 /usr/local/lib/libnng.so.1
WORKDIR /opt/xappmock
ENV LD_LIBRARY_PATH=/usr/local/lib 
ENV RMR_SEED_RT=resources/router.txt
ENV RMR_PORT=5001
#CMD mkdir -p resources/conf  exec ./xappmock
CMD mkdir -p resources/conf &&  exec /bin/bash 
