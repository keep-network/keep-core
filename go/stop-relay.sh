#!/usr/bin/env bash
# Copyright 2018 The Keep Authors.  See LICENSE.md for details.

set -o pipefail

NUM_NODES=$1
NODE_NUM=$2
NETWORK_NAME=relay

function usage() {
   echo "Usage: $(basename $0) <NUM_NODES> [NODE_NUM]"
   echo "Notes: "
   echo "· NUM_NODES must be > 3 and < 256 --OR-- NODE_NUMS == 1.  Use this when starting a swarm of nodes."
   echo "· If NUM_NODES == 1 then NODE_NUM must be > 1 and < 255."
   echo "· If --purge is passed this script will remove all related docker containers"
}

function validate_params() {
	if [ ! -z $NODE_NUM ] && [ ! -z $NUM_NODES ]; then
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
	fi
}

function cleanup_docker() {
	if [ "$NUM_NODES" == "--purge" ]; then
		echo ">> Cleanup docker containers and dangling images"
		docker rm $(docker ps -q -f status=exited) 2>/dev/null
		docker rmi $(docker images -q -f dangling=true) 2>/dev/null
		docker network rm "$NETWORK_NAME"
	elif [ "$NUM_NODES" == "1" ] && [ ! -z $NODE_NUM ]; then
		set -x
		docker stop "node$NODE_NUM"
		{ set +x; } 2>/dev/null
	else
		# stop all containers
		for i in $(seq 1 +1 $NUM_NODES); do
			set -x
			docker stop "node$i"
			{ set +x; } 2>/dev/null
		done
	fi
}

validate_params
cleanup_docker
