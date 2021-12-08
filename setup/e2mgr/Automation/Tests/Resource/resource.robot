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
#
#   This source code is part of the near-RT RIC (RAN Intelligent Controller)
#   platform project (RICP).
#

*** Settings ***
Documentation    Resource file


*** Variables ***
${docker_number}    5
${docker_number-1}    4
${url}   http://localhost:3800
${ranName}  gnb:208-092-303030
${getNodeb}  /v1/nodeb/${ranName}
${update_gnb_url}   /v1/nodeb/${ranName}/update
${update_gnb_body}  {"servedNrCells":[{"servedNrCellInformation":{"cellId":"abcd","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1,"servedPlmns":["whatever"]},"nrNeighbourInfos":[{"nrCgi":"one","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1}]}]}
${update_gnb_body_notvalid}  {"servedNrCells":[{"servedNrCellInformation":{"choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1,"servedPlmns":["whatever"]},"nrNeighbourInfos":[{"nrCgi":"whatever","choiceNrMode":{"fdd":{}},"nrMode":1,"nrPci":1}]}]}
${E2tInstanceAddress}   10.0.2.15:38000
${header}  {"Content-Type": "application/json"}
${docker_command}  docker ps | grep Up | wc --lines
${stop_simu}  docker stop gnbe2_oran_simu
${run_simu_regular}  docker run -d --name gnbe2_oran_simu --net host --env gNBipv4=10.0.2.15 --env gNBport=5577 --env ricIpv4=10.0.2.15 --env ricPort=36422 --env nbue=0  snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/gnbe2_oran_simu:3.2-32
${docker_Remove}    docker rm gnbe2_oran_simu
${docker_restart}   docker restart e2mgr
${restart_simu}  docker restart gnbe2_oran_simu
${start_e2}  docker start e2
${stop_e2}      docker stop e2
${dbass_start}   docker run -d --name dbass -p 6379:6379 --env DBAAS_SERVICE_HOST=10.0.2.15  snapshot.docker.ranco-dev-tools.eastus.cloudapp.azure.com:10001/dbass:1.0.0
${dbass_remove}    docker rm dbass
${dbass_stop}      docker stop dbass
${restart_simu}  docker restart gnbe2_oran_simu
${stop_docker_e2}      docker stop e2
${stop_routingmanager_sim}  docker stop rm_sim
${start_routingmanager_sim}  docker start rm_sim



