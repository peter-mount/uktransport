#!/bin/sh
#
# Script to run a build for a specific microservice and platform.
#
# SYNTAX
#
# docker.sh imagename microservice arch version
#
# Where arch is one of the following: amd64 arm32v6 arm32v7 arm64v8
#
# image should be the full name, e.g. area51/nre-feeds:latest or area51/nre-feeds:0.2
# This script will append -{microservice}-{arch} to that name
#

IMAGE=$1
ARCH=$2
VERSION=$3

# Resolve the architecture
case $ARCH in
  amd64)
    GOARCH=amd64
    ;;
  arm32v6)
    GOARCH=arm
    GOARM=6
    ;;
  arm32v7)
    GOARCH=arm
    GOARM=7
    ;;
  arm64v8)
    GOARCH=arm64
    ;;
  *)
    echo "Unsupported architecture $ARCH"
    exit 1
    ;;
esac

# For now just support Linux
GOOS=linux

# The actual image being built
TAG=${IMAGE}:${ARCH}-${VERSION}

echo "Building $SERVICE image $IMAGE on $ARCH"

CMD="docker build --force-rm=true"
CMD="$CMD -t ${IMAGE}"

CMD="$CMD --build-arg arch=${ARCH}"
CMD="$CMD --build-arg goos=${GOOS}"
CMD="$CMD --build-arg goarch=${GOARCH}"
CMD="$CMD --build-arg goarm=${GOARM}"

# Upload a tar file as part of the build
if [ -n "${UPLOAD_CRED}" -a -n "${UPLOAD_PATH}" ]
then
  CMD="$CMD --build-arg uploadCred=${UPLOAD_CRED}"
  CMD="$CMD --build-arg uploadPath=${UPLOAD_PATH}"
fi

CMD="$CMD ."

$CMD
