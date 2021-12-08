#!/bin/sh
# call as setup-sample-xapp.sh gnb_id

set -x

IMAGE_NAME=sample-xapp
MODEL_DIR=sample-xapp
CONNECTOR_DIR=xapp-bs-connector
DOCKER_FILE=Dockerfile
SETUP_DIR=../setup

export SRC=`dirname $0`
cd $SRC
. $SRC/setup-lib.sh

ENTRYPOINT=/bin/bash
GNB_ID=$1

# if changing xApp IP or ID, you need to define new RMR routes
# in the setup-ric.sh/setup-lib.sh scripts and restart the RIC
XAPP_IP=$XAPP_IP
XAPP_ID=$(echo $XAPP_IP | cut -d "." -f 4)

CONTAINER_NAME=${IMAGE_NAME}-${XAPP_ID}

# Build docker image
$SUDO docker image inspect ${IMAGE_NAME}:latest >/dev/null 2>&1
if [ ! $? -eq 0 ]; then
    tagvers=`git log --pretty=format:"%h" -n 1`
    $SUDO docker image inspect ${IMAGE_NAME}:$tagvers >/dev/null 2>&1
    if [ ! $? -eq 0 ]; then
            # copy Dockerfile out
            cd ${SETUP_DIR}
            cp ${MODEL_DIR}/${DOCKER_FILE} ./${DOCKER_FILE}_${IMAGE_NAME}

	    $SUDO docker build  \
            --build-arg DBAAS_SERVICE_HOST=$DBAAS_IP \
            --build-arg DBAAS_SERVICE_PORT=$DBAAS_PORT \
            -f ${DOCKER_FILE}_${IMAGE_NAME} -t ${IMAGE_NAME}:$tagvers .

            # remove copied Dockerfile
            rm ${DOCKER_FILE}_${IMAGE_NAME}

    fi
    $SUDO docker tag ${IMAGE_NAME}:$tagvers ${IMAGE_NAME}:latest
    $SUDO docker rmi ${IMAGE_NAME}:$tagvers
fi

remove_container() {
    $SUDO docker inspect $1 >/dev/null 2>&1
    if [ $? -eq 0 ]; then
	$SUDO docker kill $1
	$SUDO docker rm $1
    fi
}

# run containers
remove_container ${CONTAINER_NAME}

# replace parameters, recompile code and restart container
$SUDO docker run -d -it --entrypoint ${ENTRYPOINT} --network ric --ip ${XAPP_IP} \
    -e DBAAS_SERVICE_HOST=$DBAAS_IP -e DBAAS_SERVICE_PORT=$DBAAS_PORT --name ${CONTAINER_NAME} ${IMAGE_NAME}:latest

if [ -n "${GNB_ID}" ]; then
    docker exec ${CONTAINER_NAME} sed -i "s/^export GNB_ID.*/export GNB_ID=${GNB_ID}/g" /home/xapp-bs-connector/build_xapp.sh
fi

docker exec ${CONTAINER_NAME} sed -i "s/^export XAPP_ID.*/export XAPP_ID=${XAPP_ID}/g" /home/xapp-bs-connector/build_xapp.sh
docker exec ${CONTAINER_NAME} /home/xapp-bs-connector/build_xapp.sh clean
docker container restart ${CONTAINER_NAME}

