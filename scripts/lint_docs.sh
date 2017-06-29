#!/bin/bash

set -e

DOCKER_IMAGE_NAME="terraform-provider-sakuracloud-lint"
DOCKER_CONTAINER_NAME="terraform-provider-sakuracloud-lint-container"

if [[ $(docker ps -a | grep $DOCKER_CONTAINER_NAME) != "" ]]; then
  docker rm -f $DOCKER_CONTAINER_NAME 2>/dev/null
fi

docker build -f scripts/Dockerfile.textlint -t $DOCKER_IMAGE_NAME .

docker run -ti --rm \
       $DOCKER_IMAGE_NAME .
