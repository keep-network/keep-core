#!/bin/bash
set -e

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color

# Dafault inputs.
KEEP_ACCOUNT_PASSWORD_DEFAULT="password"
CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY=""
NETWORK_DEFAULT="local"
KEEP_CORE_PATH=$PWD
CONFIG_DIR_PATH_DEFAULT="$KEEP_CORE_PATH/configs"

help()
{
   echo ""
   echo "Usage: $0"\
        "--config-dir <path>"\
        "--account-password <password>"\
        "--private-key <private key>"\
        "--network <network>"
   echo -e "\t--config-dir: Configuration path for keep-core client"
   echo -e "\t--account-password: Account password"
   echo -e "\t--private-key: Contract owner's account private key"
   echo -e "\t--network: Host chain network for keep-core client"
   exit 1 # Exit script after printing help
}

if [ "$1" == "-help" ]; then
  help
fi

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--config-dir")        set -- "$@" "-c" ;;
    "--account-password")  set -- "$@" "-p" ;;
    "--private-key")       set -- "$@" "-k" ;;
    "--network")           set -- "$@" "-n" ;;
    *)                     set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "c:p:k:n:" opt
do
   case "$opt" in
      c ) config_dir_path="$OPTARG" ;;
      p ) account_password="$OPTARG" ;;
      k ) private_key="$OPTARG" ;;
      n ) network="$OPTARG" ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

CONFIG_DIR_PATH=${config_dir_path:-$CONFIG_DIR_PATH_DEFAULT}
KEEP_ACCOUNT_PASSWORD=${account_password:-$KEEP_ACCOUNT_PASSWORD_DEFAULT}
ACCOUNT_PRIVATE_KEY=${private_key:-$CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY}
NETWORK=${network:-$NETWORK_DEFAULT}

# Run script
printf "${LOG_START}Starting installation...${LOG_END}"

printf "Config dir path: $CONFIG_DIR_PATH\n"
printf "Network: $NETWORK"

KEEP_CORE_CONFIG_DIR_PATH=$CONFIG_DIR_PATH
KEEP_CORE_SOL_PATH="$KEEP_CORE_PATH/solidity"

cd $KEEP_CORE_SOL_PATH

printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
npm install

if [ "$NETWORK" != "alfajores" ]; then
    printf "${LOG_START}Unlocking ethereum accounts...${LOG_END}"
    KEEP_ETHEREUM_PASSWORD=$KEEP_ACCOUNT_PASSWORD \
    CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY=$ACCOUNT_PRIVATE_KEY \
        npx truffle exec scripts/unlock-eth-accounts.js --network $NETWORK
fi

printf "${LOG_START}Migrating contracts...${LOG_END}"
rm -rf build/

CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY=$ACCOUNT_PRIVATE_KEY \
    npx truffle migrate --reset --network $NETWORK

KEEP_CORE_SOL_ARTIFACTS_PATH="$KEEP_CORE_SOL_PATH/build/contracts"

printf "${LOG_START}Initializing contracts...${LOG_END}"

CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY=$ACCOUNT_PRIVATE_KEY \
    npx truffle exec scripts/delegate-tokens.js --network $NETWORK

printf "${LOG_START}Updating keep-core client configs...${LOG_END}"
for CONFIG_FILE in $KEEP_CORE_CONFIG_DIR_PATH/*.toml
do
    KEEP_CORE_CONFIG_FILE_PATH=$CONFIG_FILE \
    CONTRACT_OWNER_ACCOUNT_PRIVATE_KEY=$ACCOUNT_PRIVATE_KEY \
        npx truffle exec scripts/lcl-client-config.js --network $NETWORK
done

printf "${LOG_START}Building keep-core client...${LOG_END}"
cd $KEEP_CORE_PATH
go generate ./...
go build -a -o keep-core .
