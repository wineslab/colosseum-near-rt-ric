#!/bin/sh
# Script to setup the E2 termination.

set -x

export SRC=`dirname $0`
. $SRC/setup-lib.sh

OURDIR=../setup

cd $OURDIR
tagvers=`git log --pretty=format:"%h" -n 1`

docker kill e2term
docker rm e2term
docker rmi e2term:bronze

# build e2term
$SUDO docker image inspect e2term:bronze >/dev/null 2>&1
if [ ! $? -eq 0 ]; then
    cd e2/RIC-E2-TERMINATION
    $SUDO docker image inspect e2term:$tagvers >/dev/null 2>&1
    if [ ! $? -eq 0 ]; then
        $SUDO docker build -f Dockerfile -t e2term:$tagvers .
    fi
    $SUDO docker tag e2term:$tagvers e2term:bronze
    $SUDO docker rmi e2term:$tagvers
    cd ../..
fi

# remove dangling images
docker rmi $(docker images --filter "dangling=true" -q --no-trunc) 2> /dev/null

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

$SUDO docker run -d -it --network=ric --ip $E2TERM_IP --name e2term \
        --mount type=bind,source=$ROUTERFILE,destination=/opt/e2/dockerRouter.txt,ro \
        $E2TERM_CONFIG_BIND \
        e2term:bronze

exit 0
