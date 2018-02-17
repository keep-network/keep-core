#!/usr/bin/env bash
# Copyright 2018 The Keep Authors.  See LICENSE.md for details.

# Examples (also see usage below)
# To run only 1 node (node #3):  $ ./run-relay.sh 1 3
# To run 4 nodes :  $ ./run-relay.sh 4

set -o pipefail

NUM_NODES=$1
NODE_NUM=$2
SUBNET="192.168.0"

echo "NUM_NODES: $NUM_NODES"
echo "NODE_NUM: $NODE_NUM"
echo "$(grep = Dockerfile)" | sed -e 's/ENV//' | awk '{$1=$1}1' | sed -e 's/=/: /' | sed -e 's/\\//'

function usage() {
   echo "Usage: $(basename $0) <NUM_NODES> [NODE_NUM]"
   echo "Notes: "
   echo "· NUM_NODES must be > 3 and < 256 --OR-- NODE_NUMS == 1.  Use this when starting a swarm of nodes."
   echo "· If NUM_NODES == 1 then NODE_NUM must be > 1 and < 255."
}


function validate_params() {
	if [ -z $NODE_NUM ] && [ ! -z $NUM_NODES ]; then
		if [ $NUM_NODES -lt 4 ] || [ $NUM_NODES -gt 255 ]; then
		   usage
		   exit 2
		 fi
	elif [ ! -z $NODE_NUM ] && [ "$NUM_NODES" != "1" ]; then
	   usage
	   exit 2
	elif [ "$NUM_NODES"  == "1" ] && [ -z $NODE_NUM ]; then
	   usage
	   exit 2
	elif [ "$NUM_NODES"  == "1" ]; then
		if [ $NODE_NUM -lt 2 ] || [ $NODE_NUM -gt 255 ]; then
		   usage
		   exit 2
		fi
	elif [ -z $NODE_NUM ] && [ -z $NUM_NODES ]; then
	   usage
	   exit 2
	fi
}
validate_params

command -v go >/dev/null 2>&1 || { echo >&2 "Missing go.  For details, see: https://golang.org/doc/install"; exit 2; }

IMG="keep-client"
DOCKERFILE="./Dockerfile"
NETWORK_NAME="relay"
CMD="keep-client"

function build_docker_image() {
   echo ">> Build docker image $IMG"

   unset BUILD_ARG
   if [ "$NUM_NODES"  == "1" ] && [ ! -z $NODE_NUM ]; then
		GROUP_PORT="$(bc -l <<< $NODE_NUM+7000)"
		BUILD_ARG="--build-arg GROUP_PORT=$GROUP_PORT"
   fi

   set -x
   docker build -t "$IMG" -f "$DOCKERFILE" $BUILD_ARG  .

   { set +x;} 2>/dev/null
   if [ $? != 0 ]; then
        echo "Failed to build docker image."
        exit 1
   fi
}

function create_network_and_run_nodes() {
   echo ">> Creating $NETWORK_NAME with subnet ${SUBNET}.0/24"
   docker network rm $NETWORK_NAME &>/dev/null
   set -x
   docker network create --subnet "${SUBNET}.0/24" --driver bridge $NETWORK_NAME&>/dev/null
   { set +x;} 2>/dev/null

   for i in $(seq 1 +1 $NUM_NODES); do

	   if [ "$i" == "1" ] && [ ! -z $NODE_NUM ]; then
	          node_num=$NODE_NUM
	   else
	          node_num=$(bc -l <<< $i+1)
	   fi
       my_ip="${SUBNET}.$node_num"
       echo "My IP: $my_ip"

       # Example call: keep-client --idGenerationSeed 3 -p2pListenPort 7003 -p2pEncryption
       params=" -idGenerationSeed $node_num -p2pListenPort 700$node_num --p2pEncryption"
       echo "> Running relay node $node_num: ${SUBNET}.$node_num"

#	   echo "Running $CMD$params"
#       set -x
#       docker run \
#           --name node${node_num} \
#           --net $NETWORK_NAME \
#           --ip ${SUBNET}.${node_num} \
#           --publish 700${node_num}:700${node_num} \
#           --volume \
#           --detach \
#           --rm \
#           "$IMG" &
#
#	   sleep 2
#	   docker exec -it node${node_num} ${CMD}${params}

# For testing - Don't forget to comment out the ENTRYPOINT line in Dockerfile
  (and return --rm to docker run cmd above)
       docker run \
           --name node$i \
           --net $NETWORK_NAME \
           --ip ${SUBNET}.$node_num \
           --volume \
           --detach \
           -it \
           "$IMG"
exit

       { set +x; } 2>/dev/null

       sleep .1
   done
}


build_docker_image
create_network_and_run_nodes

