#!/bin/bash
# This script fetches contracts artifacts published to the keep-dev bucket.
# These contracts are expected to be the ones deployed on keep-dev environment.

CONTRACTS_NAMES=("KeepToken.json" "TokenStaking.json" "TokenGrant.json")

DESTINATION_DIR=$(realpath $(dirname $0)/../src/contracts)

function create_destination_dir() {
  mkdir -p $DESTINATION_DIR
}

function fetch_contracts() {
  for CONTRACT_NAME in ${CONTRACTS_NAMES[@]}
  do
    gsutil -q cp gs://keep-dev-contract-data/keep-core/${CONTRACT_NAME} $DESTINATION_DIR
  done
}

echo "Fetch contracts artifacts to: ${DESTINATION_DIR}"
create_destination_dir
fetch_contracts
