#!/bin/sh
# Script to setup the RIC containers. Call as ./setup-ric.sh [network interface]

set -x

# get flags
for ARGUMENT in "$@"
do
    KEY=$(echo $ARGUMENT | cut -f1 -d=)
    case "$KEY" in
	        arena)              arena=true;;
        import)        import=true;;
            *)   
    esac    
done

# get RIC interface from cli arguments
if [ $# -eq 0 ] || [ "$import" = false ] ; then
    RIC_INTERFACE="can0"
else
    if [ "$arena" = true ]; then
        RIC_INTERFACE="brric"
    else
        RIC_INTERFACE=$1
    fi
fi

export SRC=`dirname $0`
. $SRC/setup-lib.sh

OURDIR=../setup

# import base RIC images
if [ "$import" = true ] || [ $(docker image ls -q | wc -l) -eq "0" ]; then
    echo "Importing base Docker images"
    cd $SRC
    ./import-base-images.sh
fi

cd $OURDIR
tagvers=`git log --pretty=format:"%h" -n 1`

# build e2term
$SUDO docker image inspect e2term:latest >/dev/null 2>&1
if [ ! $? -eq 0 ]; then
    cd e2/RIC-E2-TERMINATION
    $SUDO docker image inspect e2term:$tagvers >/dev/null 2>&1
    if [ ! $? -eq 0 ]; then
        $SUDO docker build -f Dockerfile -t e2term:$tagvers .
    fi
    $SUDO docker tag e2term:$tagvers e2term:latest
    $SUDO docker rmi e2term:$tagvers
    cd ../..
fi

# build e2mgr
$SUDO docker image inspect e2mgr:latest >/dev/null 2>&1
if [ ! $? -eq 0 ]; then
    cd e2mgr/E2Manager
    $SUDO docker image inspect e2mgr:$tagvers >/dev/null 2>&1
    if [ ! $? -eq 0 ]; then 
        $SUDO docker build -f Dockerfile -t e2mgr:$tagvers .
    fi
    $SUDO docker tag e2mgr:$tagvers e2mgr:latest
    $SUDO docker rmi e2mgr:$tagvers
    cd ../..
fi

# build e2rtmansim
$SUDO docker image inspect e2rtmansim:latest >/dev/null 2>&1
if [ ! $? -eq 0 ]; then
    cd e2mgr/tools/RoutingManagerSimulator
    $SUDO docker image inspect e2rtmansim:$tagvers >/dev/null 2>&1
    if [ ! $? -eq 0 ]; then
        $SUDO docker build -f Dockerfile -t e2rtmansim:$tagvers .
    fi
    $SUDO docker tag e2rtmansim:$tagvers e2rtmansim:latest
    $SUDO docker rmi e2rtmansim:$tagvers
    cd ../../..
fi

# build dbaas
$SUDO docker image inspect dbaas:latest >/dev/null 2>&1
if [ ! $? -eq 0 ]; then
    cd dbaas
    $SUDO docker build -f docker/Dockerfile.redis -t dbaas:latest .
    cd ..
fi

# remove dangling images
docker rmi $(docker images --filter "dangling=true" -q --no-trunc) 2> /dev/null

# create a private network for near-real-time RIC
$SUDO docker network inspect ric >/dev/null 2>&1
if [ ! $? -eq 0 ]; then
    $SUDO brctl addbr brric
    $SUDO docker network create --subnet=$RIC_SUBNET -d bridge --attachable -o com.docker.network.bridge.name=brric ric
fi

# Create a route info file to tell the containers where to send various
# messages. This will be mounted on the containers
ROUTERFILE=`pwd`/router.txt
cat << EOF > $ROUTERFILE
newrt|start
rte|10020|$E2MGR_IP:3801
rte|10060|$E2TERM_IP:38000
rte|10061|$E2MGR_IP:3801
rte|10062|$E2MGR_IP:3801
rte|10070|$E2MGR_IP:3801
rte|10071|$E2MGR_IP:3801
rte|10080|$E2MGR_IP:3801
rte|10081|$E2TERM_IP:38000
rte|10082|$E2TERM_IP:38000
rte|10360|$E2TERM_IP:38000
rte|10361|$E2MGR_IP:3801
rte|10362|$E2MGR_IP:3801
rte|10370|$E2MGR_IP:3801
rte|10371|$E2TERM_IP:38000
rte|10372|$E2TERM_IP:38000
rte|1080|$E2MGR_IP:3801
rte|1090|$E2TERM_IP:38000
rte|1100|$E2MGR_IP:3801
rte|12010|$E2MGR_IP:38010
rte|1101|$E2TERM_IP:38000
rte|12002|$E2TERM_IP:38000
rte|12003|$E2TERM_IP:38000
rte|10091|$E2MGR_IP:4801
rte|10092|$E2MGR_IP:4801
rte|1101|$E2TERM_IP:38000
rte|1102|$E2MGR_IP:3801
rte|12001|$E2MGR_IP:3801
mse|12050|$(echo $XAPP_IP | cut -d "." -f 4)|$XAPP_IP:4560
newrt|end
EOF

remove_container() {
    $SUDO docker inspect $1 >/dev/null 2>&1
    if [ $? -eq 0 ]; then
        $SUDO docker kill $1
        $SUDO docker rm $1
    fi
}

# create RIC various containers. Kill and remove them if they exist.
remove_container db
$SUDO docker run -d --network ric --ip $DBAAS_IP --name db dbaas:latest

remove_container e2rtmansim
$SUDO docker run -d -it --network ric --ip $E2RTMANSIM_IP --name e2rtmansim e2rtmansim:latest

remove_container e2mgr
$SUDO docker run -d -it --network ric --ip $E2MGR_IP -e RIC_ID=7b0000-000000/18 \
    -e DBAAS_PORT_6379_TCP_ADDR=$DBAAS_IP -e DBAAS_PORT_6379_TCP_PORT="6379" \
    -e DBAAS_SERVICE_HOST=$DBAAS_IP -e DBAAS_SERCE_PORT="6379" \
    --mount type=bind,source=$ROUTERFILE,destination=/opt/E2Manager/router.txt,ro \
    --name e2mgr e2mgr:latest

remove_container e2term
E2TERMCONFFILE=`pwd`/e2term_config.conf
if [ ! -e $E2TERMCONFFILE ]; then
cat <<EOF >$E2TERMCONFFILE
nano=38000
loglevel=debug
volume=log
#the key name of the environment holds the local ip address
#ip address of the E2T in the RMR
local-ip=$E2TERM_IP
#trace is start, stop
trace=start
external-fqdn=e2t.com
#put pointer to the key that point to pod name
pod_name=E2TERM_POD_NAME
sctp-port=$E2TERM_SCTP_PORT
EOF
fi
 E2TERM_CONFIG_BIND="--mount type=bind,source=$E2TERMCONFFILE,destination=/opt/e2/config/config.conf,ro"

export RIC_IP=`ifconfig ${RIC_INTERFACE} | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*'`


if [ "$arena" = true ]; then
    echo 'Starting local setup'
    # if both RIC and DU are executed on the same machine, do not set Docker NAT rules
    $SUDO docker run -d -it --network=ric --ip $E2TERM_IP --name e2term \
        --mount type=bind,source=$ROUTERFILE,destination=/opt/e2/dockerRouter.txt,ro \
        $E2TERM_CONFIG_BIND \
        e2term:latest
else
    $SUDO docker run -d -it --network=ric --ip $E2TERM_IP --name e2term -p ${RIC_IP}:${E2TERM_SCTP_PORT}:${E2TERM_SCTP_PORT}/sctp\
        --mount type=bind,source=$ROUTERFILE,destination=/opt/e2/dockerRouter.txt,ro \
        $E2TERM_CONFIG_BIND e2term:latest
fi

exit 0
