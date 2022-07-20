#!/bin/bash
set -e

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

KEEP_CORE_PATH=$PWD

KEEP_BEACON_SOL_PATH="$KEEP_CORE_PATH/solidity/random-beacon"
TBTC_PATH="$KEEP_CORE_PATH/tbtc-v2"
TMP_TBTC="$KEEP_CORE_PATH/tmp-tbtc"

# Defaults, can be overwritten by env variables/input parameters
NETWORK_DEFAULT="development"
KEEP_ETHEREUM_PASSWORD=${KEEP_ETHEREUM_PASSWORD:-"password"}

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"\
           "--network <network>"\
           "--tbtc-path <tbtc-path>"\
           "--skip-beacon-deployment"\
           "--skip-tbtc-deployment"\
           "--skip-client-build"
   echo -e "\nEnvironment variables:\n"
   echo -e "\tKEEP_ETHEREUM_PASSWORD: The password to unlock local Ethereum accounts to set up delegations."\
           "Required only for 'local' network. Default value is 'password'"
   echo -e "\nCommand line arguments:\n"
   echo -e "\t--network: Ethereum network for keep-core client(s)."\
                        "Available networks and settings are specified in the 'hardhat.config.ts'"
   echo -e "\t--tbtc-path: 'Local' tbtc project's path. A temporary folder with tbtc is created and removed "\
                           "upon installation if the path is not provided"
   echo -e "\t--skip-beacon-deployment: When set to true the old artifacts from the '/deployments' dir are used. Default is false"
   echo -e "\t--skip-tbtc-deployment: When set to true the old artifacts from the '/deployments' dir are used. Default is false"
   echo -e "\t--skip-client-build: Should execute contracts part only. Client installation will not be executed\n"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--network")                    set -- "$@" "-n" ;;
    "--tbtc-path")                  set -- "$@" "-t" ;;
    "--skip-beacon-deployment")     set -- "$@" "-d" ;;
    "--skip-tbtc-deployment")       set -- "$@" "-e" ;;
    "--skip-client-build")          set -- "$@" "-b" ;;
    "--help")                       set -- "$@" "-h" ;;
    *)                              set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "n:t:debh" opt
do
   case "$opt" in
      n ) network="$OPTARG" ;;
      t ) tbtc_path="$OPTARG" ;;
      d ) skip_beacon_deployment=${OPTARG:-true} ;;
      e ) skip_tbtc_deployment=${OPTARG:-true} ;;
      b ) skip_client_build=${OPTARG:-true} ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Overwrite default properties
NETWORK=${network:-$NETWORK_DEFAULT}
TBTC_PATH=${tbtc_path:-""}
SKIP_BEACON_DEPLOYMENT=${skip_beacon_deployment:-false}
SKIP_TBTC_DEPLOYMENT=${skip_tbtc_deployment:-false}
SKIP_CLIENT_BUILD=${skip_client_build:-false}

# Run script
printf "${LOG_START}Starting installation...${LOG_END}"

printf "Network: $NETWORK\n"

cd $KEEP_BEACON_SOL_PATH

printf "${LOG_START}Installing beacon YARN dependencies...${LOG_END}"
yarn install

if [ "$NETWORK" == "development" ]; then
    printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
    KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
        npx hardhat unlock-accounts --network $NETWORK
fi

# Skip the beacon deployment if none of the skip-* options are passed and
# deploy TBTC which include the beacon contracts.
if [ "$SKIP_BEACON_DEPLOYMENT" = false ] && [ "$SKIP_TBTC_DEPLOYMENT" = false ]; then
  SKIP_BEACON_DEPLOYMENT=true
fi

if [ "$SKIP_BEACON_DEPLOYMENT" = false ] ; then
  printf "${LOG_START}Building beacon contracts...${LOG_END}"
  yarn clean
  yarn build

  printf "${LOG_START}Deploying contracts for beacon...${LOG_END}"

  npx hardhat deploy --reset --export export.json --network $NETWORK
fi

cd $KEEP_CORE_PATH

printf "${LOG_START}Installing tbtc YARN dependencies...${LOG_END}"
yarn install

if [ "$SKIP_TBTC_DEPLOYMENT" = false ] ; then
  if [ "$TBTC_PATH" = "" ] ; then
    printf "${LOG_START}Cloning tbtc...${LOG_END}"
    # recreate a temporary tbtc dir for fresh installation
    rm -rf $TMP_TBTC && mkdir $TMP_TBTC && cd $TMP_TBTC
    # clone project from the repository
    git clone https://github.com/keep-network/tbtc-v2.git
    cd "tbtc-v2/solidity"
  else
    printf "${LOG_START}Installing tbtc from the local directory...${LOG_END}"
    cd "$TBTC_PATH/solidity"
  fi

  printf "${LOG_START}Building tbtc contracts...${LOG_END}"
  yarn install
  yarn build

  printf "${LOG_START}Deploying contracts for tbtc...${LOG_END}"

  npx hardhat deploy --reset --export export.json --network $NETWORK

  rm -rf $TMP_TBTC
fi
    
if [ "$SKIP_CLIENT_BUILD" = false ] ; then
   printf "${LOG_START}Building beacon client...${LOG_END}"

   cd $KEEP_CORE_PATH
   go generate ./...
   go build -a -o keep-core .
fi

printf "${DONE_START}Installation completed!${DONE_END}"