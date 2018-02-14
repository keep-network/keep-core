#!/usr/bin/env bash

IMG=keep-client
DOCKERFILE=Dockerfile

docker run --rm --name keep-client_$(date +%Y%m%d-%H%M%S) \
	--detach \
	"$IMG" "$cmd"
