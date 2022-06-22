#!/bin/bash
set -euo pipefail

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

KEEP_CORE_PATH=$PWD
KEEP_BEACON_SOL_PATH="$KEEP_CORE_PATH/solidity/random-beacon"

# Defaults, can be overwritten by env variables/input parameters
NETWORK_DEFAULT="development"
KEEP_ETHEREUM_PASSWORD=${KEEP_ETHEREUM_PASSWORD:-"password"}
CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=${CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY:-""}

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"\
           "--network <network>"\
           "--skip-deployment"\
           "--skip-client-build"
   echo -e "\nEnvironment variables:\n"
   echo -e "\tKEEP_ETHEREUM_PASSWORD: The password to unlock local Ethereum accounts to set up delegations."\
           "Required only for 'local' network. Default value is 'password'"
   echo -e "\tCONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY: Contracts owner private key on Ethereum. Required for non-development network only"
   echo -e "\nCommand line arguments:\n"
   echo -e "\t--network: Ethereum network for keep-core client."\
                        "Available networks and settings are specified in the 'hardhat.config.ts'"
   echo -e "\t--skip-deployment: When set to true the old artifacts from the '/deployments' dir are used. Default is false."
   echo -e "\t--skip-client-build: Should execute contracts part only. Client installation will not be executed.\n"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--network")           set -- "$@" "-n" ;;
    "--skip-deployment")   set -- "$@" "-d" ;;
    "--skip-client-build") set -- "$@" "-b" ;;
    "--help")              set -- "$@" "-h" ;;
    *)                     set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "n:dbh" opt
do
   case "$opt" in
      n ) network="$OPTARG" ;;
      d ) skip_deployment=${OPTARG:-true} ;;
      b ) skip_client_build=${OPTARG:-true} ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Overwrite default properties
NETWORK=${network:-$NETWORK_DEFAULT}
SKIP_DEPLOYMENT=${skip_deployment:-false}
SKIP_CLIENT_BUILD=${skip_client_build:-false}

# Run script
printf "${LOG_START}Starting installation...${LOG_END}"

printf "Network: $NETWORK\n"

cd $KEEP_BEACON_SOL_PATH

printf "${LOG_START}Installing YARN dependencies...${LOG_END}"
# rm -rf node_modules/
yarn install

if [ "$NETWORK" == "development" ]; then
    printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
    KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
        npx hardhat unlock-accounts --network $NETWORK
fi

printf "${LOG_START}Building contracts...${LOG_END}"

rm -rf build && rm -rf cache && rm -rf typechain
yarn build

if [ "$SKIP_DEPLOYMENT" = false ] ; then
   printf "${LOG_START}Deploying contracts...${LOG_END}"
   # rm -rf deployments/development/

   CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
      npx hardhat deploy --reset --export export.json --network $NETWORK
fi
    
printf "${LOG_START}Setting up staking...${LOG_END}"

CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY=$CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY \
    npx hardhat staking --network $NETWORK


if [ "$SKIP_CLIENT_BUILD" = false ] ; then
   printf "${LOG_START}Building keep-core client...${LOG_END}"

   cd $KEEP_CORE_PATH
   go generate ./...
   go build -a -o keep-core .
fi

printf "${DONE_START}Installation completed!${DONE_END}"
