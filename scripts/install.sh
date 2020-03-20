#!/bin/bash
set -e

# Dafault inputs.
KEEP_ETHEREUM_PASSWORD_DEFAULT="password"
CONFIG_FILE_PATH_DEFAULT=$(realpath -m $(dirname $0)/../config.toml)

# Read user inputs.
read -p "Enter ethereum accounts password [$KEEP_ETHEREUM_PASSWORD_DEFAULT]: " ethereum_password
KEEP_ETHEREUM_PASSWORD=${ethereum_password:-$KEEP_ETHEREUM_PASSWORD_DEFAULT}

read -p "Enter path to keep-core client config [$CONFIG_FILE_PATH_DEFAULT]: " config_file_path
CONFIG_FILE_PATH=${config_file_path:-$CONFIG_FILE_PATH_DEFAULT}

# Run script.
LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color

printf "${LOG_START}Starting installation...${LOG_END}"
KEEP_CORE_PATH=$(realpath $(dirname $0)/../)
KEEP_CORE_CONFIG_FILE_PATH=$(realpath $CONFIG_FILE_PATH)
KEEP_CORE_SOL_PATH=$(realpath $KEEP_CORE_PATH/contracts/solidity)

cd $KEEP_CORE_SOL_PATH

printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
npm install

printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
    truffle exec scripts/unlock-eth-accounts.js --network local

printf "${LOG_START}Migrating contracts...${LOG_END}"
rm -rf build/
truffle migrate --reset --network local

KEEP_CORE_SOL_ARTIFACTS_PATH=$(realpath $KEEP_CORE_SOL_PATH/build/contracts)

printf "${LOG_START}Initializing contracts...${LOG_END}"
truffle exec scripts/delegate-tokens.js --network local

printf "${LOG_START}Updating keep-core client config...${LOG_END}"
KEEP_CORE_CONFIG_FILE_PATH=$KEEP_CORE_CONFIG_FILE_PATH \
    truffle exec scripts/lcl-client-config.js --network local

printf "${LOG_START}Building keep-core client...${LOG_END}"
cd $KEEP_CORE_PATH
go generate ./...
go build -a -o keep-core .
