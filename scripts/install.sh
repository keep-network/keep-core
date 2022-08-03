#!/bin/bash
set -eou pipefail

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

KEEP_CORE_PATH=$PWD

KEEP_BEACON_SOL_PATH="$KEEP_CORE_PATH/solidity/random-beacon"
KEEP_ECDSA_SOL_PATH="$KEEP_CORE_PATH/solidity/ecdsa"
TMP_TBTC="$KEEP_CORE_PATH/tmp-tbtc"

# Defaults, can be overwritten by env variables/input parameters
NETWORK_DEFAULT="development"
KEEP_ETHEREUM_PASSWORD=${KEEP_ETHEREUM_PASSWORD:-"password"}

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"\
           "--network <network>"\
           "--tbtc-path <tbtc-path>"\
           "--skip-deployment"\
           "--skip-client-build"
   echo -e "\nEnvironment variables:\n"
   echo -e "\tKEEP_ETHEREUM_PASSWORD: The password to unlock local Ethereum accounts to set up delegations."\
           "Required only for 'local' network. Default value is 'password'"
   echo -e "\nCommand line arguments:\n"
   echo -e "\t--network: Ethereum network for keep-core client(s)."\
                        "Available networks and settings are specified in the 'hardhat.config.ts'"
   echo -e "\t--tbtc-path: 'Local' tbtc project's path. 'tbtc' is cloned to a temporary directory"\
                           "upon installation if the path is not provided"
   echo -e "\t--skip-deployment: This option skips all the contracts deployment. Default is false"
   echo -e "\t--skip-client-build: Should execute contracts part only. Client installation will not be executed\n"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--network")           set -- "$@" "-n" ;;
    "--tbtc-path")         set -- "$@" "-t" ;;
    "--skip-deployment")   set -- "$@" "-e" ;;
    "--skip-client-build") set -- "$@" "-b" ;;
    "--help")              set -- "$@" "-h" ;;
    *)                     set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "n:t:ebh" opt
do
   case "$opt" in
      n ) network="$OPTARG" ;;
      t ) tbtc_path="$OPTARG" ;;
      e ) skip_deployment=${OPTARG:-true} ;;
      b ) skip_client_build=${OPTARG:-true} ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# Overwrite default properties
NETWORK=${network:-$NETWORK_DEFAULT}
TBTC_PATH=${tbtc_path:-""}
SKIP_DEPLOYMENT=${skip_deployment:-false}
SKIP_CLIENT_BUILD=${skip_client_build:-false}

# Run script
printf "${LOG_START}Starting installation...${LOG_END}"

printf "Network: $NETWORK\n"

cd $KEEP_BEACON_SOL_PATH

printf "${LOG_START}Installing beacon YARN dependencies...${LOG_END}"
yarn

if [ "$NETWORK" == "development" ]; then
    printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
    KEEP_ETHEREUM_PASSWORD=$KEEP_ETHEREUM_PASSWORD \
        npx hardhat unlock-accounts --network $NETWORK
fi

if [ "$SKIP_DEPLOYMENT" != true ]; then
  printf "${LOG_START}Building random-beacon...${LOG_END}"
  yarn clean && yarn build

  # deploy beacon
  printf "${LOG_START}Deploying random-beacon contracts...${LOG_END}"
  USE_EXTERNAL_DEPLOY=true npx hardhat deploy --reset --export export.json --network $NETWORK

  printf "${LOG_START}Creating random-beacon link...${LOG_END}"
  yarn link
  # create export folder
  yarn prepack

  cd $KEEP_ECDSA_SOL_PATH

  printf "${LOG_START}Linking random-beacon...${LOG_END}"
  yarn link @keep-network/random-beacon

  printf "${LOG_START}Building ecdsa...${LOG_END}"
  yarn && yarn clean && yarn build

  # deploy ecdsa
  printf "${LOG_START}Deploying ecdsa contracts...${LOG_END}"
  npx hardhat deploy --reset --export export.json --network $NETWORK
  
  printf "${LOG_START}Creating ecdsa link...${LOG_END}"
  yarn link
  # create export folder
  yarn prepack

  cd $KEEP_CORE_PATH

  if [ "$TBTC_PATH" = "" ]; then
    printf "${LOG_START}Cloning tbtc...${LOG_END}"
    # create a temporary tbtc dir for fresh installation
    rm -rf $TMP_TBTC && mkdir $TMP_TBTC && cd $TMP_TBTC
    # clone project from the repository
    git clone https://github.com/keep-network/tbtc-v2.git
    
    printf "${LOG_START}Building tbtc contracts...${LOG_END}"
    cd "tbtc-v2/solidity"
    yarn && yarn build && yarn prepack
  else
    printf "${LOG_START}Installing tbtc from the local directory...${LOG_END}"
    cd "$TBTC_PATH/solidity"
  fi

  printf "${LOG_START}Linking random-beacon...${LOG_END}"
  yarn link @keep-network/random-beacon

  printf "${LOG_START}Linking ecdsa...${LOG_END}"
  yarn link @keep-network/ecdsa

  # deploy tbtc
  printf "${LOG_START}Deploying tbtc contracts...${LOG_END}"
  npx hardhat deploy --reset --export export.json --network $NETWORK
fi

if [ "$SKIP_CLIENT_BUILD" = false ]; then
   printf "${LOG_START}Building client...${LOG_END}"

   cd $KEEP_CORE_PATH
   go generate ./...
   go build -a -o keep-core .
fi

printf "${DONE_START}Installation completed!${DONE_END}"