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

FROM nexus3.o-ran-sc.org:10002/o-ran-sc/bldr-ubuntu18-c-go:1.9.0 as ubuntu

WORKDIR /opt/rmsimulator
COPY . . 
ENV PATH=$PATH:/usr/local/go/bin:/usr/lib/go-1.12/bin

RUN apt-get update \
  && apt-get --allow-releaseinfo-change update \
  && apt upgrade -y
RUN go build main.go


FROM ubuntu:18.04
COPY --from=ubuntu /opt/rmsimulator/main /opt/rmsimulator/main
COPY --from=ubuntu /opt/rmsimulator/resources /opt/rmsimulator/resources

WORKDIR /opt/rmsimulator
#CMD mkdir -p resources/conf  exec ./xappmock
CMD exec ./main
