#!/bin/sh

#dockerIp=$(ifconfig eth0 | awk '{ print $2}' | grep -E -o "([0-9]{1,3}[\.]){3}[0-9]{1,3}")
#echo "docker ip: $dockerIp"
serviceName=$(env| grep RMR_SRC_ID | tr -d 'RMR_SRC_ID=' | awk '{split($0,srv,"."); print srv[1]}'| sed -r 's/-/_/g' | tr '[:lower:]' '[:upper:]')
fullServiceName=${serviceName}_SERVICE_HOST
echo "environments service name is $fullServiceName"
serviceIp=$(env | grep $fullServiceName | awk '{split($0,ip,"="); print ip[2]}')
echo "service ip is $serviceIp"
sed -i "s/local-ip=127.0.0.1/local-ip=$serviceIp/g" "/opt/e2/config/config.conf"
sed -i "s/external-fqdn=e2t.com/external-fqdn=$serviceIp/g" "/opt/e2/config/config.conf"
cat "/opt/e2/config/config.conf"
./e2 -p config -f config.conf
