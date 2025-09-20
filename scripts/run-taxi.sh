#!/bin/bash

set -e

# Determine file extension 
EXT=".lxb"
if uname | grep -iq 'mingw\|msys\|cygwin'; then
  EXT=".exe"
fi

GEN_SRC_DIR="src/tools/taxi"
GEN_SVC_PATH="tools/taxi"
GEN_BIN="$(basename "$GEN_SVC_PATH")$EXT"
GEN_BIN_PATH="bin/$GEN_BIN"

# Build
echo "Building $GEN_SVC_PATH..."
go build -C "$GEN_SRC_DIR" -o "../../../$GEN_BIN_PATH" main.go
echo "Build complete: $GEN_SVC_PATH"

"$GEN_BIN_PATH" "$@"