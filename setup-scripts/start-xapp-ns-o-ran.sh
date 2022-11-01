#!/bin/bash
docker kill sample-xapp-24
docker rm sample-xapp-24
docker rmi sample-xapp:latest
./setup-sample-xapp.sh ns-o-ran

docker exec -it sample-xapp-24 bash