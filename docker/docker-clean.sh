#!/bin/bash

# Do not run if removal already in progress.
pgrep "docker rm" && exit 0

# Remove Dead and Exited containers.
for CONTAINER in $(docker ps -a | grep "Dead\|Exited" | awk '{print $1}')
do
    docker rm -f ${CONTAINER}
done

# It will fail to remove images currently in use.
for IMAGE in $(docker images -qf dangling=true)
do
   docker rmi -f ${IMAGE}
done

# Clean up unused docker volumes
for VOLUME in $(docker volume ls -qf dangling=true)
do
    docker volume rm ${VOLUME}
done
