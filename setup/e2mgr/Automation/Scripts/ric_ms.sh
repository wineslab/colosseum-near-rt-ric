#!/bin/bash 

stringContain() { [ -z "${2##*$1*}" ]; }

MS=${1}
ACTION=${2}
SCRIPT=`basename "$0"`
USAGE="\nUsage: ./$SCRIPT <MS NAME> <ACTION>\nValid Options: ./$SCRIPT <simu,dbass,e2mgr,e2> <stop,start,restart,status>\nE.g ./$SCRIPT dbass stop"
OPTIONS="simu,dbass,e2mgr,e2,stop,start,restart,status"
DOC_SCRIPT="/opt/docker_ms.sh"
K8S_SCRIPT="/opt/k8s_ms.py"
# Check if script got the reqiured arguments 
[ -z ${SYS_TYPE} ] && echo -e "\nThe SYS_TYPE environemnt variable is not set!" && echo -e "${USAGE}" && exit 1
[ -z ${MS} ] && echo -e "\nThe MS argument is reqiured!" && echo -e "${USAGE}" && exit 2
[ -z ${ACTION} ] && echo -e "\nThe ACTION argument is reqiured!" && echo -e "${USAGE}" && exit 2
! grep -q $MS <<<"$OPTIONS" && echo -e "\nThe microservice '${MS}' is not a valid value!" &&  echo -e "${USAGE}" && exit 3
! grep -q $ACTION <<<"$OPTIONS" && echo -e "\nThe action '${ACTION}' is not a valid value!" &&  echo -e "${USAGE}" && exit 3

if [ "${SYS_TYPE}" == "docker" ]; then
   echo "SYS_TYPE=docker, Docker mode is set"
   [ ! -f ${DOC_SCRIPT} ] && echo "reqiured file '${DOC_SCRIPT}' is missing, exit" && exit 4
   echo "Executing the '${DOC_SCRIPT}' script!"
   ${DOC_SCRIPT} ${MS} ${ACTION}
elif [ "${SYS_TYPE}" == "k8s" ]; then
   echo "SYS_TYPE=k8s, K8S mode is set"
   [ ! -f ${K8S_SCRIPT} ] && echo "reqiured file '${K8S_SCRIPT}' is missing, exit" && exit 4
   echo "Executing the '${K8S_SCRIPT}' script!"
   ${K8S_SCRIPT} ${MS} ${ACTION}
fi