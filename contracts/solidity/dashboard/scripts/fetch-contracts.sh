#!/bin/bash
# This script fetches contracts artifacts published to the GCP bucket. It is expected
# that `CONTRACT_DATA_BUCKET` environment variable is set to the name of a bucket
# from which contracts should be downloaded.

CONTRACTS_NAMES=("KeepToken.json" "TokenStaking.json" "TokenGrant.json" "KeepRandomBeaconOperator.json")

DESTINATION_DIR=$(realpath $(dirname $0)/../src/contracts)

function create_destination_dir() {
  mkdir -p $DESTINATION_DIR
}

function fetch_contracts() {
  for CONTRACT_NAME in ${CONTRACTS_NAMES[@]}
  do
    gsutil -q cp gs://${CONTRACT_DATA_BUCKET}/keep-core/${CONTRACT_NAME} $DESTINATION_DIR
  done
}

echo "Fetch contracts artifacts to: ${DESTINATION_DIR}"
create_destination_dir
fetch_contracts
