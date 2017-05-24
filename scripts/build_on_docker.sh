#!/bin/bash

set -e

DOCKER_IMAGE_NAME="terraform-for-sakuracloud-build"
DOCKER_CONTAINER_NAME="terraform-for-sakuracloud-build-container"

if [[ $(docker ps -a | grep $DOCKER_CONTAINER_NAME) != "" ]]; then
  docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
fi

docker build -t $DOCKER_IMAGE_NAME .

docker run --name $DOCKER_CONTAINER_NAME \
       -e SAKURACLOUD_ACCESS_TOKEN \
       -e SAKURACLOUD_ACCESS_TOKEN_SECRET \
       -e SAKURACLOUD_ZONE \
       -e SAKURACLOUD_TRACE_MODE \
       -e TF_LOG \
       -e TESTARGS \
       $DOCKER_IMAGE_NAME make "$@"
if [[ "$@" == *"build"* ]]; then
  docker cp $DOCKER_CONTAINER_NAME:/go/src/github.com/sacloud/terraform-provider-sakuracloud/bin ./
fi
docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
