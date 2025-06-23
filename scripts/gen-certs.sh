#!/bin/bash

set -e

# Usage:
# ./gen-certs.sh <service-name>
# Example:
#   ./gen-certs.sh userservice

if [ $# -lt 1 ]; then
  echo "Usage: $0 <service-name>"
  exit 1
fi

CERT_NAME="$1"
CERT_PU_PATH="certs/${CERT_NAME}_public_key.pem"
CERT_PK_PATH="certs/${CERT_NAME}_private_key.pem"

echo "Generating certificates:"
echo " > Public Key:  $CERT_PU_PATH"
echo " > Private Key: $CERT_PK_PATH"

openssl ecparam -genkey -name prime256v1 -noout -out "$CERT_PU_PATH"
openssl ec -in "$CERT_PU_PATH" -pubout -out "$CERT_PK_PATH"
echo "Generation complete"