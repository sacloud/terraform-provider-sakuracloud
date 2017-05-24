#!/bin/bash

set -e

DOCKER_IMAGE_NAME="terraform-provider-sakuracloud-docs"
DOCKER_CONTAINER_NAME="terraform-provider-sakuracloud-docs-container"

if [[ $(docker ps -a | grep $DOCKER_CONTAINER_NAME) != "" ]]; then
  docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
fi

docker build -t $DOCKER_IMAGE_NAME -f scripts/Dockerfile.docs .

docker run --name $DOCKER_CONTAINER_NAME \
       -v $PWD/build_docs/:/go/src/github.com/sacloud/terraform-provider-sakuracloud/build_docs \
       -p 80:80 \
       $DOCKER_IMAGE_NAME serve --dev-addr=0.0.0.0:80

docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
