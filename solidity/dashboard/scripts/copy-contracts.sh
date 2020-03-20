#!/bin/bash
# This script copies contracts artifacts from provided local storage destination.
# It expects an argument with a path to the source directory holding the contracts.

CONTRACTS_NAMES=("KeepToken.json" "TokenStaking.json" "TokenGrant.json" "KeepRandomBeaconOperator.json" "Registry.json")

SOURCE_DIR=$(realpath $1)
DESTINATION_DIR=$(realpath $(dirname $0)/../src/contracts)

function create_destination_dir() {
  mkdir -p $DESTINATION_DIR
}

function copy_contracts() {
  for CONTRACT_NAME in ${CONTRACTS_NAMES[@]}
  do
    cp $(realpath $SOURCE_DIR/$CONTRACT_NAME) $DESTINATION_DIR/
  done
}

echo "Copy contracts artifacts"
echo "Source: ${SOURCE_DIR}"
echo "Destination: ${DESTINATION_DIR}"
create_destination_dir
copy_contracts
