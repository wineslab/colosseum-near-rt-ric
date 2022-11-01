#!/bin/sh
# Script to import the base images to create the RIC containers from Wineslab Docker Hub


# This shall execute first wineslab images and ns-o-ran after
# import-wines-images.sh
# setup-ric-bronze.sh

# Build image for ns-o-ran

$SUDO docker build -t ns-o-ran -f Dockerfile .

remove_container() {
    $SUDO docker inspect $1 >/dev/null 2>&1
    if [ $? -eq 0 ]; then
        $SUDO docker kill $1
        $SUDO docker rm $1
    fi
}

remove_container ns-o-ran
$SUDO docker run -d -it --network=ric --name ns-o-ran ns-o-ran