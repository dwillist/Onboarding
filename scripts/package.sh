#!/usr/bin/env bash

set -e
set -u
set x
set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

BUILDPACK_ROOT=$DIR/..

WORK_DIR="$(mktemp -d)"


function main {
  cp "$BUILDPACK_ROOT/buildpack.toml" "$WORK_DIR"
  mkdir "$WORK_DIR/bin"

  pushd $BUILDPACK_ROOT/detect
    GOOS=linux go build -o "$WORK_DIR/bin/detect" .
  popd

  pushd $BUILDPACK_ROOT/build
    GOOS=linux go build -o "$WORK_DIR/bin/build" .
  popd

  pushd "$WORK_DIR"
    tar -czvf "$BUILDPACK_ROOT/buildpack.tgz" .
  popd
}

function cleanup {
    echo "cleaning up $WORK_DIR"
    rm -r "$WORK_DIR"
}

main
