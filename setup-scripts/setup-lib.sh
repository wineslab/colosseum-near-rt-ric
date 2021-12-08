#!/bin/sh

# get sudo if needed
if [ -z "$EUID" ]; then
    EUID=`id -u`
fi
SUDO=
if [ ! $EUID -eq 0 ] ; then
    SUDO=sudo
fi

# default IPs and ports
RIC_SUBNET=10.0.2.0/24
RIC_IP=10.0.2.1
E2TERM_IP=10.0.2.10
E2TERM_SCTP_PORT=36422
E2MGR_IP=10.0.2.11
DBAAS_IP=10.0.2.12
DBAAS_PORT=6379
E2RTMANSIM_IP=10.0.2.15
XAPP_IP=10.0.2.24  # generic xApp IP

