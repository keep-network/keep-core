#!/bin/bash

# Copyright 2018 The Keep Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# This script creates a docker network for and runs relay nodes, each with
# its own IP address.
# - The first node will be the leader node; The remaining nodes are worker
#  nodes.
# - This will allow from 4 to 245 nodes per relay group

set -oe pipefail

NUM_NODES=$1
if [ -z $NUM_NODES ] || [ $NUM_NODES -lt 4 ] || [ $NUM_NODES -gt 255 ]; then
    echo "Usage: $(basename $0) <NUM_NODES>"
    echo "Notes: NUM_NODES must be > 4 and < 254"
    exit 2
fi

IMG="l3xx/relay"
DOCKERFILE="$GOPATH/src/github.com/keep-network/keep-core/go/group/Dockerfile"
DOCKER_NETWORK="relay"
SUBNET="192.168.0."
TMPDIR=$(mktemp -d /tmp/relay_$(date +%Y%m%d-%H%M%S))
CMD="relay"

function build_docker_image() {
    echo ">> Build docker image $IMG"
    docker build -t "$IMG" -f "$DOCKERFILE" .
}

function create_network_and_run_nodes() {
    echo ">> Creating $DOCKER_NETWORK with subnet ${SUBNET}0/24"
    docker network rm $DOCKER_NETWORK &>/dev/null
    docker network create $DOCKER_NETWORK --subnet "${SUBNET}0/24" &>/dev/null

    NODE_CNTR=$(seq 1 +1 $NUM_NODES)
    for i in $NODE_CNTR; do
        if [ "$i" -eq 1 ]; then LEADER=' --leader';else LEADER='';fi
        NEXT_NUM=$(bc -l <<< $i+1)
        echo "> Running relay node $i: ${SUBNET}$NEXT_NUM"

         docker run \
            --name node$i \
            --net $DOCKER_NETWORK \
            --ip ${SUBNET}$NEXT_NUM \
            --volume \
            --detach \
            --rm \
            "$IMG" \
            "$CMD$LEADER"

        sleep .1
    done
}

function cleanup_docker() {
    echo ">> Cleanup docker containers and dangling images"
    docker rm $(docker ps -q -f status=exited) 2>/dev/null
    docker rmi $(docker images -q -f dangling=true) 2>/dev/null
}

build_docker_image
create_network_and_run_nodes
cleanup_docker
