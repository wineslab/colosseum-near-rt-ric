#!/bin/bash

export DEBUG=0
export CONNECTOR_DIR=/home/xapp-sm-connector

# these are replaced through the dockerfile
export GNB_ID
export XAPP_ID

# get build clean from cli arguments
if [ $# -ne 0 ]; then
    BUILD_CLEAN=1
fi

# setup debug field
sed -i "s/^#define DEBUG.*/#define DEBUG ${DEBUG}/g" ${CONNECTOR_DIR}/src/xapp.hpp

# setup parameters
if [ -n "${GNB_ID}" ]; then
    sed -i "s/^#define GNB_ID.*/#define GNB_ID \"${GNB_ID}\"/g" ${CONNECTOR_DIR}/src/xapp.hpp
fi

if [ -n "${XAPP_ID}" ]; then
    sed -i "s/^#define XAPP_REQ_ID.*/#define XAPP_REQ_ID ${XAPP_ID}/g" ${CONNECTOR_DIR}/src/xapp.hpp
fi

# build
if [ ${BUILD_CLEAN} ]; then
    cd ${CONNECTOR_DIR}/src && make clean && make -j ${nproc} && make install && ldconfig
else
    cd ${CONNECTOR_DIR}/src && make -j ${nproc} && make install && ldconfig
fi
