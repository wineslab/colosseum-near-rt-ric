#!/usr/bin/env python


import sys, os, traceback
import docker
# Getting the Arguments. if argument are missing exit the script with exit 1
try:
    ms = sys.argv[1].lower()
    action = sys.argv[2].lower()
    script = os.path.basename(sys.argv[0])
except:
    print("Usage: %s <microservice> <action> for now only stop action is allowd" % \
            (os.path.basename(sys.argv[0])))
    sys.exit(1)

ms=sys.argv[1].lower()
action=sys.argv[2].lower()
docker_host_ip=os.environ.get('DOCKER_HOST_IP', False)
cms=[]
if not docker_host_ip:
    print('The DOCKER_HOST_IP env varibale is not defined, exiting!')
    sys.exit(1)

def get_ms():
    try: 
        client = docker.DockerClient(base_url='tcp://%s:2376' % docker_host_ip)

        for ms in client.containers.list():
            if ms.name == sys.argv[1]:
                cms.append(ms)
                return cms[0]
    except:
        print('Can\'t  connect to docker API, Exiting!')
        print(traceback.format_exc()) 
        sys.exit(1)
         
if action == 'stop':
    print('Stop the %s pod'  % ms )
    get_ms().stop()
    sys.exit(0)
else:
    print ('Only stop commnad is allowed!, exiting!')
    sys.exit(1)
