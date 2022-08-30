#!/bin/bash
set -eou pipefail

ROOT_DIR="$(realpath "$(dirname $0)/../bundles")"

if [ -z "${CHAIN_API_URL+x}" ]; then
    read -p "Provide Ethereum API URL: " CHAIN_API_URL
fi

if [ -z "${PURSE_PRIVATE_KEY+x}" ]; then
    read -p "Provide ETH Purse Private Key: " PURSE_PRIVATE_KEY
fi

if [ -z "${GOERLI_DEPLOYER_PRIVATE_KEY+x}" ]; then
    read -p "Provide GOERLI_DEPLOYER_PRIVATE_KEY: " GOERLI_DEPLOYER_PRIVATE_KEY
fi

STAKING_PROVIDER=${1-}
if [ -z "$STAKING_PROVIDER" ]; then
    read -p "Provide Staking Provider name: " STAKING_PROVIDER
fi

STAKING_PROVIDER_DIR="$(realpath "$ROOT_DIR/$STAKING_PROVIDER")"

if [ ! -d "$STAKING_PROVIDER_DIR" ]; then
    echo "Directory for $STAKING_PROVIDER does not exists."
    exit 1
fi

CONFIG_DIR="$STAKING_PROVIDER_DIR/config"
SECRETS_DIR="$STAKING_PROVIDER_DIR/secret"

KEY_FILE_PATH="$CONFIG_DIR/staking-provider-eth-account-key-file.json"
KEY_FILE_PASSWORD_PATH="$SECRETS_DIR/staking-provider-eth-account-password"
PRIVATE_KEY_FILE_PATH="$SECRETS_DIR/staking-provider-eth-account-private-key"

ACCOUNT_ADDRESS=$(jq -jr .address $KEY_FILE_PATH)
ACCOUNT_PRIVATE_KEY=$(cat $PRIVATE_KEY_FILE_PATH)

[[ $ACCOUNT_ADDRESS == 0x* ]] || ACCOUNT_ADDRESS="0x$ACCOUNT_ADDRESS"

printf "Staking Provider Account Address: $ACCOUNT_ADDRESS\n"

printf "Pull the latest images...\n"

docker pull gcr.io/keep-test-f3e0/keep-random-beacon-hardhat:latest
docker pull gcr.io/keep-test-f3e0/keep-ecdsa-hardhat:latest

printf "Fund Staking Provider Account Address with ether from purse...\n"

docker run \
    --rm \
    --env "CHAIN_API_URL=$CHAIN_API_URL" \
    --env "ACCOUNTS_PRIVATE_KEYS=$PURSE_PRIVATE_KEY" \
    gcr.io/keep-test-f3e0/keep-random-beacon-hardhat:latest \
    ensure-eth-balance \
    --network goerli \
    --target-balance "0.1 ether" \
    $ACCOUNT_ADDRESS

printf "Initialize staking...\n"

docker run \
    --rm \
    --env "CHAIN_API_URL=$CHAIN_API_URL" \
    --env "ACCOUNTS_PRIVATE_KEYS=$GOERLI_DEPLOYER_PRIVATE_KEY,$ACCOUNT_PRIVATE_KEY" \
    gcr.io/keep-test-f3e0/keep-random-beacon-hardhat:latest \
    initialize:staking \
    --network goerli \
    --owner $ACCOUNT_ADDRESS \
    --provider $ACCOUNT_ADDRESS

printf "Authorize the Random Beacon...\n"

docker run \
    --rm \
    --env "CHAIN_API_URL=$CHAIN_API_URL" \
    --env "ACCOUNTS_PRIVATE_KEYS=$GOERLI_DEPLOYER_PRIVATE_KEY,$ACCOUNT_PRIVATE_KEY" \
    gcr.io/keep-test-f3e0/keep-random-beacon-hardhat:latest \
    authorize:beacon \
    --network goerli \
    --owner $ACCOUNT_ADDRESS \
    --provider $ACCOUNT_ADDRESS

printf "Authorize the ECDSA...\n"

docker run \
    --rm \
    --env "CHAIN_API_URL=$CHAIN_API_URL" \
    --env "ACCOUNTS_PRIVATE_KEYS=$GOERLI_DEPLOYER_PRIVATE_KEY,$ACCOUNT_PRIVATE_KEY" \
    gcr.io/keep-test-f3e0/keep-ecdsa-hardhat:latest \
    authorize:ecdsa \
    --network goerli \
    --owner $ACCOUNT_ADDRESS \
    --provider $ACCOUNT_ADDRESS

printf "\n\e[1;32mDONE!\n\n\e[0m"
