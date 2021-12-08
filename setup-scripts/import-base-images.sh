#!/bin/bash

images=('bldr-alpine3-go_6-a3.11-rmr3' 'bldr-ubuntu18-c-go_9-u18.04')
images_dir=/root/o-ran_images/

# load Docker images in parallel
for i in "${images[@]}"; do
    docker load --input ${images_dir}${i}.tar.gz &
done

# wait for above jobs to finish
wait
