##-----------------------------------
# BS Connector
#------------------------------------
ARG SCHEMA_PATH=schemas
ARG XAPP_DIR=sample-xapp
ARG CONNECTOR_DIR=xapp-bs-connector
ARG STAGE_DIR=/home/${CONNECTOR_DIR}
ARG DBAAS_SERVICE_HOST=10.0.2.12
ARG DBAAS_SERVICE_PORT="6379"

#==================================================================================
FROM ubuntu:18.04 as ricbuild

# to override repo base, pass in repo argument when running docker build:
# docker build --build-arg REPOBASE=http://abc.def.org . ....
ARG REPOBASE=https://gerrit.oran-osc.org/r
ARG SCHEMA_FILE
ARG SCHEMA_PATH
ARG XAPP_DIR
ARG CONNECTOR_DIR
ARG STAGE_DIR

ARG MDC_VER=0.0.4-1
ARG RMR_VER=4.0.5
ARG RNIB_VER=1.0.0

ENV RMR_RTG_SVC="9999" \
  RMR_SEED_RT="/home/${CONNECTOR_DIR}/init/routes.txt" \
  LD_LIBRARY_PATH="/usr/local/lib:/usr/local/libexec" \
  VERBOSE=0 \
  CONFIG_FILE=/opt/ric/config/config-file.json \
  DBAAS_SERVICE_HOST=${DBAAS_SERVICE_HOST} \
  DBAAS_SERVICE_PORT=${DBAAS_SERVICE_PORT}

COPY ${CONNECTOR_DIR} ${STAGE_DIR}
COPY ${XAPP_DIR} /home/${XAPP_DIR}

WORKDIR ${STAGE_DIR}

# Install necessary packages
RUN apt-get update  \
  && apt-get install -y \
    cmake \
    git \
    build-essential \
    automake \
    autoconf-archive \
    autoconf \
    pkg-config \
    gawk \
    libtool \
    wget \
    zlib1g-dev \
    libffi-dev \
    libcurl4-openssl-dev \
    vim \
    cpputest \
    libboost-all-dev \
    libhiredis-dev \
  && apt-get remove -y --purge python* \
  && apt-get install -y \
    python3 \
    python3-pip \
# Install mdclog using debian package hosted at packagecloud.io
  && wget -nv --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/mdclog_${MDC_VER}_amd64.deb/download.deb \
  && wget -nv --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/mdclog-dev_${MDC_VER}_amd64.deb/download.deb \
  && dpkg -i mdclog_${MDC_VER}_amd64.deb \
  && dpkg -i mdclog-dev_${MDC_VER}_amd64.deb \
  && rm mdclog_${MDC_VER}_amd64.deb mdclog-dev_${MDC_VER}_amd64.deb \
# Install RMR using debian package hosted at packagecloud.io
  && wget -nv --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr_${RMR_VER}_amd64.deb/download.deb \
  && wget -nv --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr-dev_${RMR_VER}_amd64.deb/download.deb \
  && dpkg -i rmr_${RMR_VER}_amd64.deb \
  && dpkg -i rmr-dev_${RMR_VER}_amd64.deb \
  && rm rmr_${RMR_VER}_amd64.deb rmr-dev_${RMR_VER}_amd64.deb \
# Install RNIB libraries
  && wget -nv --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rnib_${RNIB_VER}_all.deb/download.deb \
  && dpkg -i rnib_${RNIB_VER}_all.deb \
  && rm rnib_${RNIB_VER}_all.deb \
# Install dbaas
  && cd ${STAGE_DIR} \
  && git clone https://gerrit.o-ran-sc.org/r/ric-plt/dbaas \
  && cd dbaas/redismodule \
  && ./autogen.sh \
  && ./configure \
  && make -j ${nproc} all \
  && make install \
  && rm -Rf ${STAGE_DIR}/dbaas \
# Install sdl
  && cd ${STAGE_DIR} \
  && git clone https://gerrit.o-ran-sc.org/r/ric-plt/sdl \
  && cd sdl \
  && ./autogen.sh \
  && ./configure \
  && make -j ${nproc} all \
  && make install \
  && rm -Rf ${STAGE_DIR}/sdl \
# Install rapidjson
  && cd ${STAGE_DIR} \
  && git clone https://github.com/Tencent/rapidjson \
  && cd rapidjson \
  && mkdir build \
  && cd build \
  && cmake -DCMAKE_INSTALL_PREFIX=/usr/local .. \
  && make -j ${nproc} \
  && make install \
  && cd ${STAGE_DIR} \
  && rm -rf rapidjson \
##-----------------------------------
# Now install the program
#------------------------------------
# build connector
  && export CPATH=$CPATH:/usr/local/include \
  && cd src \
  && make clean \
  && make -j ${nproc} \
  && make install \
  && ldconfig \
  && sysctl -w net.ipv6.conf.lo.disable_ipv6=1 \
###-----------------------------------
## DRL Agent
##------------------------------------
## install packages and pip requirements
#  && pip3 --no-cache-dir install -U pip \
#  && pip3 --no-cache-dir install -r /home/${XAPP_DIR}/requirements.txt \
  && apt-get autoremove -y \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*
