#!/bin/sh
#
# Builds all binaries
#
DEST=$1

BIN_DIR=${DEST}/bin/

mkdir -p ${BIN_DIR}

for bin in \
  nptgimport
do
  echo "Building ${bin}"
  OUT=
  go build \
    -o ${BIN_DIR}/${bin} \
    github.com/peter-mount/uktransport/${bin}/bin
done
