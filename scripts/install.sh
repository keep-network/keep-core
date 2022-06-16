#!/bin/bash
set -euo pipefail

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

KEEP_CORE_PATH=$PWD
KEEP_BEACON_SOL_PATH="$KEEP_CORE_PATH/solidity/random-beacon"

# Defaults, can be overwritten by env variables/input parameters
CONFIG_DIR_PATH_DEFAULT="$KEEP_CORE_PATH/configs"
NETWORK_DEFAULT="development"
KEEP_ETHEREUM_PASSWORD=${KEEP_ETHEREUM_PASSWORD:-"password"}
CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=${CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY:-""}

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"\
           "--config-dir <path>"\
           "--network <network>"\
           "--contracts-only"
   echo -e "\nEnvironment variables:\n"
   echo -e "\tKEEP_ETHEREUM_PASSWORD: The password to unlock local Ethereum accounts to set up delegations."\
           "Required only for 'local' network. Default value is 'password'"
   echo -e "\tCONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY: Contracts owner private key on Ethereum. Required for non-local network only"
   echo -e "\nCommand line arguments:\n"
   echo -e "\t--config-dir: Path to keep-core client configuration file(s)"
   echo -e "\t--network: Ethereum network for keep-core client."\
                        "Available networks and settings are specified in the 'truffle-config.js'"
   echo -e "\t--contracts-only: Should execute contracts part only."\
                        "Client installation will not be executed.\n"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--config-dir")     set -- "$@" "-c" ;;
    "--network")        set -- "$@" "-n" ;;
    "--contracts-only") set -- "$@" "-m" ;;
    "--help")           set -- "$@" "-h" ;;
    *)                  set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "c:n:mh" opt
do
   case "$opt" in
      c ) config_dir_path="$OPTARG" ;;
      n ) network="$OPTARG" ;;
      m ) contracts_only=true ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Overwrite default properties
CONFIG_DIR_PATH=${config_dir_path:-$CONFIG_DIR_PATH_DEFAULT}
NETWORK=${network:-$NETWORK_DEFAULT}
CONTRACTS_ONLY=${contracts_only:-false}

# Run script
printf "${LOG_START}Starting installation...${LOG_END}"

printf "Config dir path: $CONFIG_DIR_PATH\n"
printf "Network: $NETWORK"

cd $KEEP_BEACON_SOL_PATH

printf "${LOG_START}Installing YARN dependencies...${LOG_END}"
yarn install

if [ "$NETWORK" == "development" ]; then
    printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
    KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
        npx hardhat run scripts/unlock-eth-accounts.ts --network $NETWORK
fi

printf "${LOG_START}Deploying contracts...${LOG_END}"
rm -rf build/

CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
    npx hardhat deploy --reset --export export.json --network $NETWORK
    
# TODO: T staking
#       Increasing authorization on RandomBeacon
#       Build the client

# printf "${LOG_START}Copying contract artifacts...${LOG_END}"
# rm -rf artifacts
# cp -r build/contracts artifacts
# npm link

# printf "${LOG_START}Initializing contracts...${LOG_END}"

# CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
#     npx truffle exec scripts/delegate-tokens.js --network $NETWORK

# if [ "$CONTRACTS_ONLY" = false ] ; then
#   printf "${LOG_START}Updating keep-core client configs...${LOG_END}"
#   for CONFIG_FILE in $CONFIG_DIR_PATH/*.toml
#   do
#       KEEP_CORE_CONFIG_FILE_PATH=$CONFIG_FILE \
#       CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
#           npx truffle exec scripts/lcl-client-config.js --network $NETWORK
#   done

#   printf "${LOG_START}Building keep-core client...${LOG_END}"
#   cd $KEEP_CORE_PATH
#   go generate ./...
#   go build -a -o keep-core .
# fi

printf "${DONE_START}Installation completed!${DONE_END}"
