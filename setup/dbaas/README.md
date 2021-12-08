# Database as a service repository

This repository containes all the needed elements to deploy database as a service to kubernetes

## Architecture

Redis is the chosen database technology and the final product will deploy autonomous
redis cluster. In R0 schedule, single, non-redundant, non-persistent redis server is
deployed

## Subsystem structure

**docker** Contains dockerfiles to produce dbaas / testapplication container images
**charts** Contais helm charts to deploy dbaas service / testapplication
**testapplication** Contains dbaas test applications with various languages such as go, ..

## Container image creation

The images must be built at subsystem root level

To produce dbaas service image:
```
docker build --file docker/Dockerfile.redis --tag redis-standalone .
```

To produce testapplication image:
```
docker build --file docker/Dockerfile.testapp --tag dbaas-test .
```

## Deployment

### DBaaS service

In R0, dbaas service is realized with single container running redis database.
The database is configured to be non-persistent and non-redundant. The container
exposes single port which is hardcoded to be 6379.

After dbaas service is installed, environment variables **DBAAS_SERVICE_HOST**
and **DBAAS_SERVICE_PORT** are exposed to application containers. SDL library
will automatically use these environment variables.

The service is installed via helm by using dbaas-service chart. Modify the
values accordingly before installation (repository location, image name, ..)

```
helm install ./dbaas-service
```

### DBaaS test application

Test application is installed via helm by using dbaas-test chart. Modify the
values accordingly before installation (repository location, image name, ..)

```
helm install ./dbaas-test
```

## Testing

Make sure that dbaas-service and dbaas-test application are deployed:
```
>>helm ls
NAME            REVISION  UPDATED                    STATUS     CHART                   APP VERSION	NAMESPACE
angry-greyhound	1         Thu Mar 21 11:36:23 2019   DEPLOYED	dbaas-test-0.1.0	1.0             default
loitering-toad  1         Thu Mar 21 11:35:21 2019   DEPLOYED	dbaas-0.1.0             1.0             default
```

Check the deployed pods
```
>>kubectl get pods
NAME                                READY   STATUS    RESTARTS   AGE
dbaas-test-app-7695dbb9ff-qn8c2     1/1     Running   0          5s
redis-standalone-78978f4c6f-54b2s   1/1     Running   0          66s
```

Connect to the test application container:
```
kubectl exec -it dbaas-test-app-7695dbb9ff-qn8c2 -- /bin/bash
```

In test application container:
```
The environment variables for database backend should be set:

>>printenv
DBAAS_SERVICE_HOST=10.108.103.51
DBAAS_SERVICE_PORT=6379


Go test application using preliminary go SDL-API should be able to perform reads and writes:

>>./testapp
key1:data1
key3:%!s(<nil>)
key2:data2
num1:1
num2:2
-------------
mix2:2
num1:1
num2:2
pair1:data1
array1:adata1
mix1:data1
mix3:data3
mix4:4
arr1:
key1:data1
key2:data2
pair2:data2
array2:adata2


Redis server can be pinged with redis-cli:

>>redis-cli -h $DBAAS_SERVICE_HOST -p $DBAAS_SERVICE_PORT ping
PONG
```

## License
This project is licensed under the Apache License 2.0 - see the [LICENSE.md](LICENSE.md) file for details
