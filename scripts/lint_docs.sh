#!/bin/bash

set -e

DOCKER_IMAGE_NAME="sacloud/textlint:latest"

docker run -ti --rm \
       -v $PWD/build_docs:/workdir \
       $DOCKER_IMAGE_NAME .
