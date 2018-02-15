#!/usr/bin/env bash

# Run both the build and the test script like this:  
#   ./build-keep-client-docker-img.sh && ./run-keep-client.sh

IMG=keep-client
DOCKERFILE=Dockerfile
CMD=keep-client

docker run --rm --name keep-client_$(date +%Y%m%d-%H%M%S) \
	--detach \
	"$IMG" "$CMD"
