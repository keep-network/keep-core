#!/bin/bash
set -eou pipefail

docker buildx build \
    --platform=linux/amd64 \
    --target build-docker \
    --tag keep-client-go-build-env \
    --load \
    .

docker run \
    --workdir /go/src/github.com/keep-network/keep-core \
    keep-client-go-build-env \
    gotestsum -- -timeout 20m -tags=integration ./...
