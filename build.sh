#!/bin/bash
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

# The actual image being built
TAG=${IMAGE}:${ARCH}-${VERSION}

. functions.sh

CMD="docker build --force-rm=true"
CMD="$CMD -t ${TAG}"

CMD="$CMD --build-arg arch=${ARCH}"

# For now just support linux
CMD="$CMD --build-arg goos=linux"

CMD="$CMD --build-arg goarch=$(goarch $ARCH)"
CMD="$CMD --build-arg goarm=$(goarm $ARCH)"

# Upload a tar file as part of the build
if [ -n "${UPLOAD_CRED}" -a -n "${JOB_NAME}" ]
then
  CMD="$CMD --build-arg uploadCred=${UPLOAD_CRED}"

  if [ -z "${UPLOAD_PATH}" ]
  then
    UPLOAD_PATH=https://nexus.area51.onl/repository/snapshots/
  fi

  repoPath=$JOB_NAME
  if [ "$(basename $repoPath)" = "${BRANCH_NAME}" ]
  then
    repoPath=$(dirname $repoPath)
  fi

  # add version to the path
  repoPath="${repoPath}/${VERSION}"

  repoName="$(basename $repoPath)-${ARCH}-${VERSION}"
  if [ -n "${BUILD_NUMBER}" ]
  then
    repoName="${repoName}.${BUILD_NUMBER}"
  fi

  CMD="$CMD --build-arg uploadPath=${UPLOAD_PATH}/${repoPath}"
  CMD="$CMD --build-arg uploadName=${repoName}"
fi

CMD="$CMD ."

execute $CMD
