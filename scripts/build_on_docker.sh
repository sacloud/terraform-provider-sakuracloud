#!/bin/bash

set -e

DOCKER_IMAGE_NAME="terraform-for-sakuracloud-build"
DOCKER_CONTAINER_NAME="terraform-for-sakuracloud-build-container"

if [[ $(docker ps -a | grep $DOCKER_CONTAINER_NAME) != "" ]]; then
  docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
fi

docker build -t $DOCKER_IMAGE_NAME .

docker run --name $DOCKER_CONTAINER_NAME $DOCKER_IMAGE_NAME make "$@"
docker cp $DOCKER_CONTAINER_NAME:/go/src/github.com/yamamoto-febc/terraform-provider-sakuracloud/bin ./
docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
