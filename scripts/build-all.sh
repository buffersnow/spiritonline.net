#!/usr/bin/env bash
set -euo pipefail

# Iterate over each folder in src/cmd
for dir in src/cmd/*/; do
  # Get just the folder name (service name)
  service_name=$(basename "$dir")

  echo ">>> Running service: $service_name"
  ./run-service "$service_name" "$@"
done