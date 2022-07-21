#!/bin/bash
set -euo pipefail

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

KEEP_CORE_PATH=$PWD
KEEP_CORE_SOL_PATH="$KEEP_CORE_PATH/solidity-v1"

# Defaults, can be overwritten by env variables/input parameters
NETWORK_DEFAULT="local"
KEEP_ETHEREUM_PASSWORD=${KEEP_ETHEREUM_PASSWORD:-"password"}
CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=${CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY:-""}

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"\
           "--network <network>"\
   echo -e "\nEnvironment variables:\n"
   echo -e "\tKEEP_ETHEREUM_PASSWORD: The password to unlock local Ethereum accounts to set up delegations."\
           "Required only for 'local' network. Default value is 'password'"
   echo -e "\tCONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY: Contracts owner private key on Ethereum. Required for non-local network only"
   echo -e "\nCommand line arguments:\n"
   echo -e "\t--network: Ethereum network for keep-core client."\
                        "Available networks and settings are specified in the 'truffle-config.js'"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--network")        set -- "$@" "-n" ;;
    "--help")           set -- "$@" "-h" ;;
    *)                  set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "n:h" opt
do
   case "$opt" in
      n ) network="$OPTARG" ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Overwrite default properties
NETWORK=${network:-$NETWORK_DEFAULT}

# Run script
printf "${LOG_START}Starting installation...${LOG_END}"

printf "Network: $NETWORK"

cd $KEEP_CORE_SOL_PATH

printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
npm install

if [ "$NETWORK" == "local" ]; then
    printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
    KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
        npx truffle exec scripts/unlock-eth-accounts.js --network $NETWORK
fi

printf "${LOG_START}Migrating contracts...${LOG_END}"
rm -rf build/

CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
    npx truffle migrate --reset --network $NETWORK

printf "${LOG_START}Copying contract artifacts...${LOG_END}"
rm -rf artifacts
cp -r build/contracts artifacts
npm link

printf "${LOG_START}Initializing contracts...${LOG_END}"

printf "${DONE_START}Installation completed!${DONE_END}"
