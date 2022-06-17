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
           "--reuse-contracts"
   echo -e "\nEnvironment variables:\n"
   echo -e "\tKEEP_ETHEREUM_PASSWORD: The password to unlock local Ethereum accounts to set up delegations."\
           "Required only for 'local' network. Default value is 'password'"
   echo -e "\tCONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY: Contracts owner private key on Ethereum. Required for non-local network only"
   echo -e "\nCommand line arguments:\n"
   echo -e "\t--config-dir: Path to keep-core client configuration file(s)"
   echo -e "\t--network: Ethereum network for keep-core client."\
                        "Available networks and settings are specified in the 'hardhat.config.ts'"
   echo -e "\t--reuse-contracts: The already deployed contracts can be reused to speed up the development. Default is false.\n"\
   "Client installation will not be executed.\n"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--config-dir")      set -- "$@" "-c" ;;
    "--network")         set -- "$@" "-n" ;;
    "--reuse-contracts") set -- "$@" "-d" ;;
    "--help")            set -- "$@" "-h" ;;
    *)                   set -- "$@" "$arg"
  esac
done



# Parse short options
OPTIND=1
while getopts "c:n:dh" opt
do
   case "$opt" in
      c ) config_dir_path="$OPTARG" ;;
      n ) network="$OPTARG" ;;
      d ) reuse_contracts=true ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Overwrite default properties
CONFIG_DIR_PATH=${config_dir_path:-$CONFIG_DIR_PATH_DEFAULT}
NETWORK=${network:-$NETWORK_DEFAULT}
REUSE_CONTRACTS=${reuse_contracts:-false}

echo "reuse contracts::" ${REUSE_CONTRACTS}

# Run script
printf "${LOG_START}Starting installation...${LOG_END}"

printf "Config dir path: $CONFIG_DIR_PATH\n"
printf "Network: $NETWORK"

cd $KEEP_BEACON_SOL_PATH

printf "${LOG_START}Installing YARN dependencies...${LOG_END}"
# rm -rf node_modules/
yarn install

if [ "$NETWORK" == "development" ]; then
    printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
    KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
        npx hardhat run scripts/unlock-eth-accounts.ts --network $NETWORK
fi

printf "${LOG_START}Building contracts...${LOG_END}"
# rm -rf build && rm -rf cache && rm -rf typechain

yarn build

if [ "$REUSE_CONTRACTS" = false ] ; then
   printf "${LOG_START}Deploying contracts...${LOG_END}"
   rm -rf deployments/development/

   CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
      npx hardhat deploy --reset --export export.json --network $NETWORK
fi
    
printf "${LOG_START}Initializing contracts...${LOG_END}"

CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
    npx hardhat run scripts/init-contracts.ts --network $NETWORK

printf "${LOG_START}Building keep-core client...${LOG_END}"

cd $KEEP_CORE_PATH
go generate ./...
go build -a -o keep-core .

printf "${DONE_START}Installation completed!${DONE_END}"
