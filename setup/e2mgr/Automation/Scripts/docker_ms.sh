#!/bin/bash


MS=$1
ACTION=$2

do_stop(){
     MS=$1
     if ! docker ps --filter "name=^/${MS}" | grep -q "${MS}"; then
           echo "${MS} is already stopped, ignore the action."
     else
           echo "Executing 'docker stop ${MS}'"
           docker stop ${MS}
     fi
}

do_start(){
     MS=$1
     if docker ps --filter "name=^/${MS}" | grep -q "${MS}"; then
        echo "${MS} is running, performing restart."
        echo "Executing \'\docker stop ${MS}'"
        docker stop ${MS} && sleep 2
        echo "Executing 'docker start ${MS}'"
        docker start ${MS} 
     else
        echo "Executing 'docker start ${MS}'"
        docker start ${MS}
     fi
}


do_status(){
     MS=$1
     out=$(docker ps --filter "name=^/${MS}" | grep "${MS}")
     res=$?
     if [ "$res" == "0" ]; then
        echo $out
        echo "The ${MS} is currnetly up & running!"
     else
        echo "The ${MS} is currnetly not running!"
     fi
}


case $ACTION in
   start)
       do_start ${MS}
   ;;
   stop)
       do_stop ${MS}
   ;;
   status)
       do_status ${MS}
   ;;
   restart)
       do_stop ${MS}
       do_start ${MS}
   ;;
   *)
       do_status ${MS}
   ;;
esac