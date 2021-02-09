#!/bin/bash
set -euo pipefail

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color

KEEP_CORE_PATH=$PWD
KEEP_CORE_SOL_PATH="$KEEP_CORE_PATH/solidity"

# Defaults, can be overwritten by env variables/input parameters
CONFIG_DIR_PATH_DEFAULT="$KEEP_CORE_PATH/configs"
NETWORK_DEFAULT="local"
KEEP_ETHEREUM_PASSWORD=${KEEP_ETH_ACCOUNT_PASSWORD:-"password"}
CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=${CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY:-""}

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"\
           "--config-dir <path>"\
           "--network <network>"
   echo -e "\nEnvironment variables:\n"
   echo -e "\tKEEP_ETH_ACCOUNT_PASSWORD: Unlock an account with a password. Default password is 'password'"
   echo -e "\tCONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY: Contracts owner private key on Ethereum"
   echo -e "\nCommand line arguments:\n"
   echo -e "\t--config-dir: Path to keep-core client configuration file(s)"
   echo -e "\t--network: Host chain network for keep-core client\n"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--config-dir")  set -- "$@" "-c" ;;
    "--network")     set -- "$@" "-n" ;;
    "--help")        set -- "$@" "-h" ;;
    *)               set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "c:n:h" opt
do
   case "$opt" in
      c ) config_dir_path="$OPTARG" ;;
      n ) network="$OPTARG" ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Overwrite default properties
CONFIG_DIR_PATH=${config_dir_path:-$CONFIG_DIR_PATH_DEFAULT}
NETWORK=${network:-$NETWORK_DEFAULT}

# Run script
printf "${LOG_START}Starting installation...${LOG_END}"

printf "Config dir path: $CONFIG_DIR_PATH\n"
printf "Network: $NETWORK"

cd $KEEP_CORE_SOL_PATH

printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
npm install

printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
    npx truffle exec scripts/unlock-eth-accounts.js --network $NETWORK

printf "${LOG_START}Migrating contracts...${LOG_END}"
rm -rf build/

CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
    npx truffle migrate --reset --network $NETWORK

KEEP_CORE_SOL_ARTIFACTS_PATH="$KEEP_CORE_SOL_PATH/build/contracts"

printf "${LOG_START}Initializing contracts...${LOG_END}"

CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
    npx truffle exec scripts/delegate-tokens.js --network $NETWORK

printf "${LOG_START}Updating keep-core client configs...${LOG_END}"
for CONFIG_FILE in $CONFIG_DIR_PATH/*.toml
do
    KEEP_CORE_CONFIG_FILE_PATH=$CONFIG_FILE \
    CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
        npx truffle exec scripts/lcl-client-config.js --network $NETWORK
done

printf "${LOG_START}Building keep-core client...${LOG_END}"
cd $KEEP_CORE_PATH
go generate ./...
go build -a -o keep-core .
