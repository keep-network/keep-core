#!/bin/bash
set -eou pipefail

ROOT_DIR="$(realpath "$(dirname $0)/../")"

ENVIRONMENT_DEFAULT=goerli
read -p "Ethereum Network [$ENVIRONMENT_DEFAULT]: " ENVIRONMENT
ENVIRONMENT=${ENVIRONMENT:-$ENVIRONMENT_DEFAULT}

VERSION_DEFAULT=$(git describe --tags --match "v[0-9]*" HEAD)
read -p "Version [$VERSION_DEFAULT]: " VERSION
VERSION=${ENVIRONMENT:-$VERSION_DEFAULT}

REVISION_DEFAULT=$(git rev-parse --short HEAD)
read -p "Revision [$REVISION_DEFAULT]: " REVISION
VERSION=${REVISION:-$REVISION_DEFAULT}

APP_NAME_DEFAULT=keep-client-${ENVIRONMENT}
read -p "Application Name (prefix of the output file name) [$APP_NAME_DEFAULT]: " APP_NAME
APP_NAME=${APP_NAME:-$APP_NAME_DEFAULT}

BIN_OUTPUT_DIR=${ROOT_DIR}/bin/

docker buildx build \
    --platform linux/amd64 \
    --output type=local,dest=${BIN_OUTPUT_DIR} \
    --target=output-bins \
    --build-arg ENVIRONMENT=${ENVIRONMENT} \
    --build-arg APP_NAME=${APP_NAME} \
    --build-arg VERSION=${VERSION} \
    --build-arg REVISION=${REVISION} \
    .
