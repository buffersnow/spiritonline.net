#!/usr/bin/env bash
set -euo pipefail

# Iterate over each folder in src/cmd
for dir in src/cmd/*/; do
    # Get just the folder name (service name)
    service_name=$(basename "$dir")

    echo ""
    echo ">>> Building service: $service_name"
    ./scripts/run-service.sh "$service_name" "$@"
done