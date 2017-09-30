#!/bin/sh

DOCKER_SOCKET=/var/run/docker.sock
DOCKER_GROUP=docker
USER=sackci

# make the docker socket usable by the unprivileged user
if [ -S ${DOCKER_SOCKET} ]; then
    DOCKER_GID=$(stat -c '%g' ${DOCKER_SOCKET})
    addgroup -g ${DOCKER_GID} ${DOCKER_GROUP}
    addgroup ${USER} ${DOCKER_GROUP}
fi

# run sackci as unprivileged user
gosu ${USER}:${DOCKER_GROUP} /usr/sbin/sackci $@