#!/bin/bash

# Requires geth https://github.com/ethereum/go-ethereum/wiki/Installation-Instructions-for-Mac
# Requres bip39-cli https://www.npmjs.com/package/bip39-cli

# Please don't use this for mainnet accounts, it's for testing only.

HELP="Usage: ./$(basename $0) -n <number of accounts> \ne.g. ./$(basename $0) -n 5"

while getopts ":n:" opt; do
  case $opt in
    n  ) NUMBER_OF_ACCOUNTS=$OPTARG;;
    \? ) echo "Unknown option: -$OPTARG"; echo -e $HELP; exit 1
  esac
done

if [ $# -eq 0 ]
then
  echo -e $HELP
  exit 1
fi

PASSWORD_FILE="./account-password.txt"
ACCOUNT_INFO_FILE="./account-info.txt"

for ((ACCOUNT_ORDINAL=0; $ACCOUNT_ORDINAL<$NUMBER_OF_ACCOUNTS; ACCOUNT_ORDINAL++))
do
  echo "=====Account $ACCOUNT_ORDINAL====="

  ACCOUNT_PASSWORD=$(bip39-cli generate)

  echo $ACCOUNT_PASSWORD > $PASSWORD_FILE

  ACCOUNT=$(geth account new --keystore ./ --password $PASSWORD_FILE)

  echo "Account $ACCOUNT_ORDINAL: $ACCOUNT / $ACCOUNT_PASSWORD" | tee -a $ACCOUNT_INFO_FILE
done

rm $PASSWORD_FILE