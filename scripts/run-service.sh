#!/bin/bash

set -e

usage() {
  echo "Usage: $0 <service-name> [...] [extra-run-args...]"
  echo "  Build Mode:                   [-d|-t|-p]"
  echo "  Run after Build:              [-r]"
  echo "  Cross-Compile Arch:           [-a <arch>]"
  echo "  Cross-Compile Platform:       [-p <platform>]"
  echo "  Executable Extension:         [-e <extension>]"
  echo "  Skip DevFingerprint Creation: [-s]"
  echo ""
  echo "Specify no Arguments or -h to display the usage help"
  exit 1
}

# Defaults

MODE="prod"            
ACTION="build_only"      
ARCH="native"              
PLATFORM="native"       
DO_FINGERPRINT="true"
EXT_OVERRIDE=""

POSITIONAL=()
while [[ $# -gt 0 ]]; do
  case "$1" in
    -r)
      ACTION="build_and_run"
      shift
      ;;
    -d)
      MODE="dev"
      shift
      ;;
    -t)
      MODE="test"
      shift
      ;;
    -p)
      MODE="prod"
      shift
      ;;
    -s)
      DO_FINGERPRINT="false"
      shift
      ;;
    -a)
      ARCH="$2"
      shift 2
      ;;
    -p)
      PLATFORM="$2"
      shift 2
      ;;
    -e)
      EXT_OVERRIDE="$2"
      shift 2
      ;;
    -h)
      usage
      ;;
    *)
      POSITIONAL+=("$1")
      shift
      ;;
  esac
done

# Restore positional args
set -- "${POSITIONAL[@]}"

if [ $# -lt 1 ]; then
  usage
fi

REL_PATH="$1"
shift
SRC_DIR="src/cmd/$REL_PATH"

if [ ! -d "$SRC_DIR" ]; then
  echo "Error: Directory $SRC_DIR does not exist."
  exit 1
fi

# Determine file extension based on OS
EXT=".lxb"
if uname | grep -iq 'mingw\|msys\|cygwin|windows_nt'; then
  EXT=".exe"
fi

# Manual Cross-Compilation override for linux
if [ "$PLATFORM" == "linux" ]; then
  EXT=".lxb" # manual override for cross-compilation
fi

if [ -n "$EXT_OVERRIDE" ]; then
  EXT="$EXT_OVERRIDE"
fi

# Output binary info
BUILD_DIR="bin"
mkdir -p "$BUILD_DIR"

BIN_NAME="$REL_PATH$EXT"
BIN_PATH="$BUILD_DIR/$BIN_NAME"

# Build Tag
DNT_BUILD_SERVICE="$REL_PATH"
DNT_BUILD_VERSION_MAJOR="0"
DNT_BUILD_VERSION_MINOR="0"
DNT_BUILD_VERSION_HOTFIX="0"
DNT_BUILD_VERSION_COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")"
DNT_BUILD_CFG_OSARCH="$(go env GOARCH)"
DNT_BUILD_CFG_BLDMODE="$MODE"
DNT_BUILD_CFG_OSNAME="$(go env GOOS)"
DNT_BUILD_LAB_NAME="$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown-branch")"
DNT_BUILD_LAB_HOST="$(hostname)"
DNT_BUILD_LAB_USERNAME="$(whoami)"
DNT_BUILD_IDEN_TIMESTAMP="$(date +"%y%m%d-%H%M")"

# Determine Git Tag
TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
if [[ "$TAG" =~ ^v([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
  DNT_BUILD_VERSION_MAJOR="${BASH_REMATCH[1]}"
  DNT_BUILD_VERSION_MINOR="${BASH_REMATCH[2]}"
  DNT_BUILD_VERSION_HOTFIX="${BASH_REMATCH[3]}"
fi

# Dev Fingerprint
DNT_BUILD_IDEN_DEVFINGERPRINT="0"
if [ "$DO_FINGERPRINT" == "true" ]; then
  FINGERPRINT_SRC="${DNT_BUILD_SERVICE}-${DNT_BUILD_VERSION_MAJOR}-"
  FINGERPRINT_SRC+="${DNT_BUILD_VERSION_MINOR}-${DNT_BUILD_VERSION_HOTFIX}-"
  FINGERPRINT_SRC+="${DNT_BUILD_VERSION_COMMIT}-${DNT_BUILD_CFG_OSARCH}-"
  FINGERPRINT_SRC+="${DNT_BUILD_CFG_BLDMODE}-${DNT_BUILD_CFG_OSNAME}-"
  FINGERPRINT_SRC+="${DNT_BUILD_LAB_NAME}-${DNT_BUILD_LAB_HOST}-"
  FINGERPRINT_SRC+="${DNT_BUILD_LAB_USERNAME}-${DNT_BUILD_IDEN_TIMESTAMP}"
  DNT_BUILD_IDEN_DEVFINGERPRINT="$(echo "$FINGERPRINT_SRC" | sha1sum | cut -c1-16)"
fi

MODULE_PATH="buffersnow.com/spiritonline/pkg/version"
LDFLAGS=" -X '${MODULE_PATH}.DoNotTouch_Build_Service=${DNT_BUILD_SERVICE}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Version_Major=${DNT_BUILD_VERSION_MAJOR}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Version_Minor=${DNT_BUILD_VERSION_MINOR}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Version_HotFix=${DNT_BUILD_VERSION_HOTFIX}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Version_Commit=${DNT_BUILD_VERSION_COMMIT}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Cfg_OSArch=${DNT_BUILD_CFG_OSARCH}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Cfg_BldMode=${DNT_BUILD_CFG_BLDMODE}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Cfg_OSName=${DNT_BUILD_CFG_OSNAME}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Lab_Name=${DNT_BUILD_LAB_NAME}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Lab_Host=${DNT_BUILD_LAB_HOST}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Lab_Username=${DNT_BUILD_LAB_USERNAME}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Iden_Timestamp=${DNT_BUILD_IDEN_TIMESTAMP}'"
LDFLAGS+=" -X '${MODULE_PATH}.DoNotTouch_Build_Iden_DevFingerprint=${DNT_BUILD_IDEN_DEVFINGERPRINT}'"

if [ "$MODE" == "prod" ]; then
  LDFLAGS="-w -s ${LDFLAGS}"
fi

# Enable cross compliation support
ENV_VARS=""
if [ "$ARCH" != "native" ]; then
  ENV_VARS+="GOARCH=$ARCH "
fi
if [ "$PLATFORM" != "native" ]; then
  ENV_VARS+="GOOS=$PLATFORM "
fi

# Build
echo "Building $DNT_BUILD_SERVICE..."
echo " > Version: $DNT_BUILD_VERSION_MAJOR.$DNT_BUILD_VERSION_MINOR.$DNT_BUILD_VERSION_HOTFIX.$DNT_BUILD_VERSION_COMMIT"
echo " > Config:  ${DNT_BUILD_CFG_OSARCH}${DNT_BUILD_CFG_BLDMODE}.${DNT_BUILD_CFG_OSNAME}"
echo " > Lab:     $DNT_BUILD_LAB_NAME"
eval "${ENV_VARS} go build -C \"$SRC_DIR\" -ldflags=\"$LDFLAGS\" -o \"../../../$BIN_PATH\" main.go" 
echo "Build complete: $BIN_PATH"

# If running, build run args
if [ "$ACTION" == "build_and_run" ]; then
  RUN_ARGS=()

  # If dev run (-d), add debug and nologcmpr flags
  if [ "$MODE" == "dev" ]; then
    RUN_ARGS+=("--devel" "--no-archive" "-v")
  elif [ "$MODE" == "test" ]; then 
    RUN_ARGS+=("--no-archive" "-v")
  fi

  # Append any extra user-provided args
  RUN_ARGS+=("$@")

  echo -e "Running $DNT_BUILD_SERVICE with args: ${RUN_ARGS[*]}\n"
  "$BIN_PATH" "${RUN_ARGS[@]}"
fi