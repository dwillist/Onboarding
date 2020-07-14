#!/usr/bin/env bash

set -e
set -u
set x
set -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

BUILDPACK_ROOT=$DIR/..

WORK_DIR="$(mktemp -d)"


function main {
  local version

  while [[ "${#}" != 0 ]]; do
    case "${1}" in
      --version|-v)
        version="${2}"
        shift 2
        ;;

      "")
        # skip if the argument is empty
        shift 1
        ;;

      *)
        echo "unknown argument \"${1}\""
    esac
  done

  if [[ "${version:-}" == "" ]]; then
    echo "$version"
    echo "--version is required"
  fi

  cp "$BUILDPACK_ROOT/buildpack.toml" "$WORK_DIR"
  mkdir "$WORK_DIR/bin"

  pushd $BUILDPACK_ROOT/detect
    go build -o "$WORK_DIR/bin/detect" .
  popd

  pushd $BUILDPACK_ROOT/build
    go build -o "$WORK_DIR/bin/build" .
  popd

  tar -czvf -C "$BUILDPACK_ROOT/buildpack.tgz" "$WORK_DIR"
}

function cleanup {
    echo "cleaning up $WORK_DIR"
    #rm -r "$WORK_DIR"
}

main
