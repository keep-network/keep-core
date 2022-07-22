#!/bin/bash
set -eo pipefail

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color
DONE_START='\n\e[1;32m' # new line + bold + green
DONE_END='\n\n\e[0m'    # new line + reset

KEEP_CORE_PATH=$PWD
TMP_TBTC="$KEEP_CORE_PATH/tmp-tbtc"
KEEP_BEACON_SOL_PATH="$KEEP_CORE_PATH/solidity/random-beacon"
KEEP_ECDSA_SOL_PATH="$KEEP_CORE_PATH/solidity/ecdsa"
KEEP_TBTC_SOL_PATH="$TMP_TBTC/tbtc-v2/solidity"

# Defaults, can be overwritten by env variables/input parameters
NETWORK_DEFAULT="development"

help()
{
   echo -e "\nUsage: $0"\
           "--network <network>"\
           "--stake-owner <stake owner address>"\
           "--staking-provider <staking provider address>"\
           "--operator <operator address>"\
           "--beneficiary <beneficiary address>"\
           "--authorizer <authorizer address>"\
           "--staking-amount <staking amount>"\
           "--authorization-amount <authorization amount>"
   echo -e "\nMandatory line arguments:\n"
   echo -e "\t--stake-owner: Stake owner address"
   echo -e "\nOptional line arguments:\n"
   echo -e "\t--network: Ethereum network for keep-core client."\
                        "Available networks and settings are specified in the 'hardhat.config.ts'"
   echo -e "\t--staking-provider: Staking provider address"
   echo -e "\t--operator: Operator address"
   echo -e "\t--beneficiary: Staking beneficiary address"
   echo -e "\t--authorizer: Staking authorizer address"
   echo -e "\t--stake-amount: Staking amount"
   echo -e "\t--authorization-amount: Authorization amount"
   echo -e "\t--tbtc-path: 'Local' tbtc project's path. A temporary folder with tbtc is created and removed"\
                           "upon installation if the path is not provided"
   echo -e "\t--skip-ecdsa-deployment: This option skips ecdsa and tbtc deployment. Default is false"
   echo -e "\t--skip-tbtc-deployment: This option skips tbtc deployment. Default is false\n"
   exit 1 # Exit script after printing help
}

# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--network")              set -- "$@" "-n" ;;
    "--stake-owner")          set -- "$@" "-o" ;;
    "--staking-provider")     set -- "$@" "-p" ;;
    "--operator")             set -- "$@" "-d" ;;
    "--beneficiary")          set -- "$@" "-b" ;;
    "--authorizer")           set -- "$@" "-a" ;;
    "--stake-amount")         set -- "$@" "-s" ;;
    "--authorization-amount") set -- "$@" "-k" ;;
    "--tbtc-path")              set -- "$@" "-p" ;;
    "--skip-ecdsa-deployment")  set -- "$@" "-e" ;;
    "--skip-tbtc-deployment")   set -- "$@" "-t" ;;
    "--help")                 set -- "$@" "-h" ;;
    *)                        set -- "$@" "$arg"
  esac
done

# Parse short options
OPTIND=1
while getopts "n:o:p:d:b:a:s:k:p:eth" opt
do
   case "$opt" in
      n ) network="$OPTARG" ;;
      o ) stake_owner="$OPTARG" ;;
      p ) staking_provider="$OPTARG" ;;
      d ) operator="$OPTARG" ;;
      b ) beneficiary="$OPTARG" ;;
      a ) authorizer="$OPTARG" ;;
      s ) stake_amount="$OPTARG" ;;
      k ) authorization_amount="$OPTARG" ;;
      p ) tbtc_path="$OPTARG" ;;
      e ) skip_ecdsa_deployment=${OPTARG:-true} ;;
      t ) skip_tbtc_deployment=${OPTARG:-true} ;;
      h ) help ;;
      ? ) help ;; # Print help in case parameter is non-existent
   esac
done
shift $(expr $OPTIND - 1) # remove options from positional parameters

if [ -z "$stake_owner" ]; then
   echo 'Stake owner address must be provided. See --help'
   exit 1
fi

# Overwrite default properties
NETWORK=${network:-$NETWORK_DEFAULT}
TBTC_PATH=${tbtc_path:-""}
SKIP_ECDSA_DEPLOYMENT=${skip_ecdsa_deployment:-false}
SKIP_TBTC_DEPLOYMENT=${skip_tbtc_deployment:-false}

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

printf "${LOG_START}Setting up staking...${LOG_END}"

mint="npx hardhat mint --network $NETWORK --owner ${stake_owner}"
stake="npx hardhat stake --network $NETWORK --owner ${stake_owner} --provider ${staking_provider} --beneficiary ${beneficiary} --authorizer ${authorizer} "
increase_authorization="npx hardhat increase-authorization --network $NETWORK --owner ${stake_owner} --provider ${staking_provider} --authorizer ${authorizer}"
register_operator="npx hardhat register-operator --network $NETWORK --owner ${stake_owner} --provider ${staking_provider} --operator ${operator}"

application="--application RandomBeacon"

if [ "$SKIP_ECDSA_DEPLOYMENT" = true ]; then
   cd $KEEP_BEACON_SOL_PATH
   # go to beacon
elif [ "$SKIP_TBTC_DEPLOYMENT" = true ]; then
   # go to ecdsa (includes beacon contracts)
   cd $KEEP_ECDSA_SOL_PATH
else
   # go to tbtc (includes beacon and ecdsa contracts)
   if [ "$TBTC_PATH" = "" ]; then
      cd "$KEEP_TBTC_SOL_PATH"  
   else
      cd "$TBTC_PATH/solidity"
   fi
fi

eval ${mint} ${stake_amount_opt}
eval ${stake} ${stake_amount_opt}

application="--application RandomBeacon"
eval ${increase_authorization} ${application} ${authorization_amount_opt}
eval ${register_operator} ${application} 

if [ "$SKIP_ECDSA_DEPLOYMENT" = false ]; then
   application="--application WalletRegistry"
   eval ${increase_authorization} ${application} ${authorization_amount_opt}
   eval ${register_operator} ${application}
fi

printf "${DONE_START}Initialization completed!${DONE_END}"