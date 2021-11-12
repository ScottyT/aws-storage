#!/usr/bin/env bash
set -eo pipefail
echo "Downloading files..."
requireEnv() {
  test "${!1}" || (echo "server: '$1' not found" >&2 && exit 1)
}
requireEnv BUCKET_NAME

aws s3 sync s3://${BUCKET_NAME}/$1 $2

echo "wrote to $2"