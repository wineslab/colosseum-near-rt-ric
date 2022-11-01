#!/bin/sh
# Script to import the base images to create the RIC containers from Wineslab Docker Hub

# Pull Wines base images
docker pull wineslab/o-ran-sc-bldr-ubuntu18-c-go:9-u18.04
docker pull wineslab/o-ran-sc-bldr-alpine3-go:6-a3.11-rmr3

# Pull Wines RIC images
docker pull wineslab/colo-ran-e2term:bronze
docker pull wineslab/colo-ran-e2mgr:bronze
docker pull wineslab/colo-ran-e2rtmansim:bronze
docker pull wineslab/colo-ran-dbaas:bronze

# Tag images to be used with the setup-ric script
docker tag wineslab/colo-ran-e2term:bronze e2term:bronze
docker tag wineslab/colo-ran-e2mgr:bronze e2mgr:bronze
docker tag wineslab/colo-ran-e2rtmansim:bronze e2rtmansim:bronze
docker tag wineslab/colo-ran-dbaas:bronze dbaas:bronze
