#!/bin/bash

set -e

DOCKER_IMAGE_NAME="terraform-provider-sakuracloud-docs"
DOCKER_CONTAINER_NAME="terraform-provider-sakuracloud-docs-container"

if [[ $(docker ps -a | grep $DOCKER_CONTAINER_NAME) != "" ]]; then
  docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
fi

docker build -f scripts/Dockerfile.docs -t $DOCKER_IMAGE_NAME .

docker run --name $DOCKER_CONTAINER_NAME \
       $DOCKER_IMAGE_NAME

rm -rf docs/
docker cp $DOCKER_CONTAINER_NAME:/terraform-provider-sakuracloud/build_docs/site docs
docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null

