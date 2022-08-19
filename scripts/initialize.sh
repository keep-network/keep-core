#!/bin/bash
set -eo pipefail

LOG_START='\n\e[1;36m'  # new line + bold + color
LOG_END='\n\e[0m'       # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

KEEP_CORE_PATH=$PWD
CONFIG_DIR_PATH_DEFAULT="$KEEP_CORE_PATH/configs"
KEEP_BEACON_SOL_PATH="$KEEP_CORE_PATH/solidity/random-beacon"
KEEP_ECDSA_SOL_PATH="$KEEP_CORE_PATH/solidity/ecdsa"

# Defaults, can be overwritten by env variables/input parameters
NETWORK_DEFAULT="development"

help() {
   echo -e "\nUsage: $0" \
      "--network <network>" \
      "--stake-owner <stake owner address>" \
      "--staking-provider <staking provider address>" \
      "--operator <operator address>" \
      "--beneficiary <beneficiary address>" \
      "--authorizer <authorizer address>" \
      "--stake-amount <stake amount>" \
      "--authorization-amount <authorization amount>"
   echo -e "\nMandatory line arguments:\n"
   echo -e "\t--stake-owner: Stake owner address"
   echo -e "\nOptional line arguments:\n"
   echo -e "\t--network: Ethereum network for keep-core client." \
      "Available networks and settings are specified in the 'hardhat.config.ts'"
   echo -e "\t--staking-provider: Staking provider address"
   echo -e "\t--operator: Operator address"
   echo -e "\t--beneficiary: Staking beneficiary address"
   echo -e "\t--authorizer: Staking authorizer address"
   echo -e "\t--stake-amount: Stake amount"
   echo -e "\t--authorization-amount: Authorization amount"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--network")              set -- "$@" "-n" ;;
    "--config-dir-path")      set -- "$@" "-c" ;;
    "--stake-owner")          set -- "$@" "-o" ;;
    "--staking-provider")     set -- "$@" "-p" ;;
    "--operator")             set -- "$@" "-d" ;;
    "--beneficiary")          set -- "$@" "-b" ;;
    "--authorizer")           set -- "$@" "-a" ;;
    "--stake-amount")         set -- "$@" "-s" ;;
    "--authorization-amount") set -- "$@" "-k" ;;
    "--help")                 set -- "$@" "-h" ;;
    *)                        set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "n:c:o:p:d:b:a:s:k:h" opt
do
   case "$opt" in
      n ) network="$OPTARG" ;;
      c ) config_dir_path="$OPTARG" ;;
      o ) stake_owner="$OPTARG" ;;
      p ) staking_provider="$OPTARG" ;;
      d ) operator="$OPTARG" ;;
      b ) beneficiary="$OPTARG" ;;
      a ) authorizer="$OPTARG" ;;
      s ) stake_amount="$OPTARG" ;;
      k ) authorization_amount="$OPTARG" ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

CONFIG_DIR_PATH=${config_dir_path:-$CONFIG_DIR_PATH_DEFAULT}

# read from the config file if a stake_owner is not provided as parameter
if [ -z "$stake_owner" ]; then
   printf "\nReading stake owner address from the config file..."

   # read from the config file when the stake owner is not provided
   config_files=($CONFIG_DIR_PATH/*.toml)
   config_files_count=${#config_files[@]}
   while :; do
      printf "\nSelect client config file: \n"
      i=1
      for o in "${config_files[@]}"; do
         echo "$i) ${o##*/}"
         let i++
      done

      read reply
      if [ "$reply" -ge 1 ] && [ "$reply" -le $config_files_count ]; then
         CONFIG_FILE_PATH=${config_files["$reply" - 1]}
         break
      else
         printf "\nInvalid choice. Please choose an existing option number.\n"
      fi
   done

   sed -n -e '/^keyFile -/p' $CONFIG_FILE_PATH

   key_file_str=$(grep "^keyFile" $CONFIG_FILE_PATH)
   # internal field separator, creates array from key_file separated by '"' sign
   # array[1] is the key file path
   IFS='"' read -r -a array <<< "$key_file_str"

   # find address:<address> in the key file
   address_str=$(grep -Eo '"address":.*?[^\\]"' "${array[1]}")
   # address_array[3] is the ethereum address
   IFS='"' read -r -a address_array <<< "$address_str"
   printf "\nStake owner address: ${address_array[3]} \n"

   stake_owner="${address_array[3]}"
fi

# Overwrite default properties
NETWORK=${network:-$NETWORK_DEFAULT}

if [ -z "$staking_provider" ]; then
   staking_provider=${stake_owner}
fi

if [ -z "$operator" ]; then
   operator=${stake_owner}
fi

if [ -z "$beneficiary" ]; then
   beneficiary=${stake_owner}
fi

if [ -z "$authorizer" ]; then
   authorizer=${stake_owner}
fi

stake_amount_opt=""
if [ ! -z "$stake_amount" ]; then
   stake_amount_opt="--amount ${stake_amount}"
fi

authorization_amount_opt=""
if [ ! -z "$authorization_amount" ]; then
   authorization_amount_opt="--authorization ${authorization_amount}"
fi

initialize="npx hardhat initialize 
   --network $NETWORK \
   --owner ${stake_owner} \
   --provider ${staking_provider} \
   --operator ${operator} \
   --beneficiary ${beneficiary} \
   --authorizer ${authorizer}"

printf "${LOG_START}Initializing beacon...${LOG_END}"
cd $KEEP_BEACON_SOL_PATH
eval ${initialize} ${stake_amount_opt} ${authorization_amount_opt}

printf "${LOG_START}Initializing ecdsa...${LOG_END}"
cd $KEEP_ECDSA_SOL_PATH
eval ${initialize} ${stake_amount_opt} ${authorization_amount_opt}

printf "${DONE_START}Initialization completed!${DONE_END}"
