#!/bin/bash

set -e

DOCKER_IMAGE_NAME="sacloud/mkdocs:latest"
DOCKER_CONTAINER_NAME="terraform-provider-sakuracloud-docs-container"

if [[ $(docker ps -a | grep $DOCKER_CONTAINER_NAME) != "" ]]; then
  docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
fi

docker run --name $DOCKER_CONTAINER_NAME \
       -v $PWD/build_docs:/workdir \
       $DOCKER_IMAGE_NAME

rm -rf docs/
docker cp $DOCKER_CONTAINER_NAME:/workdir/site docs
docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
