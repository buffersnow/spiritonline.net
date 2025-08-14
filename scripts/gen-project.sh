#!/bin/bash

set -e

# Usage:
# ./gen-project.sh <service-name>
# Example:
#   ./gen-project.sh userservice

if [ $# -lt 1 ]; then
  echo "Usage: $0 <service-name> <service-short-name>"
  exit 1
fi

# Determine file extension 
EXT=".lxb"
if uname | grep -iq 'mingw\|msys\|cygwin'; then
  EXT=".exe"
fi

SERVICE_NAME="$1"
GEN_SRC_DIR="src/tools/pd"
GEN_SVC_PATH="tools/pd"
GEN_BIN="$(basename "$GEN_SVC_PATH")$EXT"
GEN_BIN_PATH="bin/$GEN_BIN"

# Build
echo "Building $GEN_SVC_PATH..."
go build -C "$GEN_SRC_DIR" -o "../../../$GEN_BIN_PATH" main.go
echo "Build complete: $GEN_SVC_PATH"

echo "Generating project:"
echo " > Service Name: $SERVICE_NAME"

"$GEN_BIN_PATH" "$SERVICE_NAME"