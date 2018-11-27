#!/bin/bash
#
# Script to generate a multi-architecture docker image from the
# individual images
#
# SYNTAX
#
# multiarch.sh imagename version architectures...
#
# Where arch is one of the following: amd64 arm32v6 arm32v7 arm64v8
#
# image should be the full name, e.g. area51/nre-feeds:latest or area51/nre-feeds:0.2
# This script will append -{microservice}-{arch} to that name
#

IMAGE=$1
shift

VERSION=$2
shift

# The final multiarch image
MULTIIMAGE=${IMAGE}:${VERSION}

. functions.sh

CMD="docker manifest create -a ${MULTIIMAGE}"
for arch in $@
do
  CMD="$CMD $(dockerImage $arch)"
done
execute $CMD

for arch in $@
do
  CMD="docker manifest annotate"
  CMD="$CMD --os linux"
  CMD="$CMD --arch $(goarch $arch)"
  CMD="$CMD $MULTIIMAGE"
  CMD="$CMD $(dockerImage $arch)"
  execute $CMD
done

execute docker manifest push -p $MULTIIMAGE