#!/bin/bash
set -e

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color

WORKDIR=$PWD

# default file for calculated amounts of stakedrop for KEEP holders / stakers
KEEP_AMOUNTS_BY_HOLDERS_PATH="$WORKDIR/output/result.json"

NODE_CURRENT_VER="$(node --version)"
NODE_REQUIRED_VER="v14.3.0"

if [ "$(printf '%s\n' "$NODE_REQUIRED_VER" "$NODE_CURRENT_VER" | sort -V | head -n1)" != "$NODE_REQUIRED_VER" ]; 
then
      echo "Required node version must be at least ${NODE_REQUIRED_VER}" 
      exit 1
fi

help()
{
   echo -e "\nUsage: ENV_VAR(S) $0"\
           "--target-block <block_number>"
   echo -e "\nEnvironment variables:\n"
   echo -e "\tETH_ACCOUNT_PRIVATE_KEY: Ethereum account private key\n"\
           "\tETH_HOSTNAME: Ethereum endpoint hostname"
   echo -e "\nCommand line arguments:\n"
   echo -e "\t--target-block: Block height when the stakedrop happens"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--target-block") set -- "$@" "-t" ;;
    "--help")         set -- "$@" "-h" ;;
    *)                set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "t:h" opt
do
   case "$opt" in
      t ) target_block="$OPTARG" ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

# enter ETH_ACCOUNT_PRIVATE_KEY if not provided
if [[ -z "$ETH_ACCOUNT_PRIVATE_KEY" ]]; then
  read -p "Enter Ethereum account private key: " ETH_ACCOUNT_PRIVATE_KEY
fi

# enter ETH_HOSTNAME if not provided
if [[ -z "$ETH_HOSTNAME" ]]; then
  read -p "Enter Ethereum hostname: " ETH_HOSTNAME
fi

# enter target_block if not provided
if [[ -z "$target_block" ]]; then
  read -p "Enter target block number: " target_block
fi

printf "${LOG_START}Initializing merkle-distributor submodule...${LOG_END}"

git submodule update --init --recursive --remote --rebase --force

printf "${LOG_START}Installing dependencies for merkle-distributor...${LOG_END}"

cd "$WORKDIR/merkle-distributor"
npm i

printf "${LOG_START}Installing dependencies for token-stakedrop...${LOG_END}"

cd "$WORKDIR"
npm i

printf "${LOG_START}Tracking KEEP holders...${LOG_END}"

node --experimental-json-modules ./bin/inspect-token-ownership.js --target-block $target_block

printf "${LOG_START}Installing dependencies for merkle-generator...${LOG_END}"

cd "$WORKDIR/merkle-generator"
npm i

printf "${LOG_START}Generating merkle output object...${LOG_END}"

npm run generate-merkle-root -- --input="$KEEP_AMOUNTS_BY_HOLDERS_PATH"

printf "${LOG_START}Script finished successfully!${LOG_END}"