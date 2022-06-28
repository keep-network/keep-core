#!/bin/bash

set -e

PROJECT_ROOT=$(realpath $(dirname "$0")/..)

packages=(ecdsa random-beacon tbtc-v2 threshold-network)

for package in "${packages[@]}"; do
    (cd ${PROJECT_ROOT}/pkg/chain/${package}/gen && make download_artifacts)
done
