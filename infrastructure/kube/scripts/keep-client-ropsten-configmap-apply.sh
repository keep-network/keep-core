#!/bin/bash
set -e

HELP="Usage: ./$(basename $0) -e environment"

while getopts ":e:" opt; do
  case $opt in
    e  ) ENVIRONMENT=$OPTARG;;
    \? ) echo "Unknown option: -$OPTARG"; echo -e $HELP; exit 1
  esac
done

if [ $# -eq 0 ]
then
  echo -e $HELP
  exit 1
fi

kubectl create configmap keep-client-0 \
  --from-file=eth_account_keyfile=../${ENVIRONMENT}/files/keep-client-0-eth-keyfile.json \
  --from-file=eth_account_password=../../eth-networks/${ENVIRONMENT}/ropsten/eth-account-password.txt \
  --from-file=keep-client-config.toml=../${ENVIRONMENT}/files/keep-client-0-config.toml \
  --namespace ropsten \
  --dry-run \
  --save-config \
  -o yaml | kubectl apply -f -


kubectl create configmap keep-client-1 \
  --from-file=eth_account_keyfile=../${ENVIRONMENT}/files/keep-client-1-eth-keyfile.json \
  --from-file=eth_account_password=../../eth-networks/${ENVIRONMENT}/ropsten/eth-account-password.txt \
  --from-file=keep-client-config.toml=../${ENVIRONMENT}/files/keep-client-1-config.toml \
  --namespace ropsten \
  --dry-run \
  --save-config \
  -o yaml | kubectl apply -f -

kubectl create configmap keep-client-2 \
  --from-file=eth_account_keyfile=../${ENVIRONMENT}/files/keep-client-2-eth-keyfile.json \
  --from-file=eth_account_password=../../eth-networks/${ENVIRONMENT}/ropsten/eth-account-password.txt \
  --from-file=keep-client-config.toml=../${ENVIRONMENT}/files/keep-client-2-config.toml \
  --namespace ropsten \
  --dry-run \
  --save-config \
  -o yaml | kubectl apply -f -

kubectl create configmap keep-client-3 \
  --from-file=eth_account_keyfile=../${ENVIRONMENT}/files/keep-client-3-eth-keyfile.json \
  --from-file=eth_account_password=../../eth-networks/${ENVIRONMENT}/ropsten/eth-account-password.txt \
  --from-file=keep-client-config.toml=../${ENVIRONMENT}/files/keep-client-3-config.toml \
  --namespace ropsten \
  --dry-run \
  --save-config \
  -o yaml | kubectl apply -f -

kubectl create configmap keep-client-4 \
  --from-file=eth_account_keyfile=../${ENVIRONMENT}/files/keep-client-4-eth-keyfile.json \
  --from-file=eth_account_password=../../eth-networks/${ENVIRONMENT}/ropsten/eth-account-password.txt \
  --from-file=keep-client-config.toml=../${ENVIRONMENT}/files/keep-client-4-config.toml \
  --namespace ropsten \
  --dry-run \
  --save-config \
  -o yaml | kubectl apply -f -

kubectl create configmap keep-client-5 \
  --from-file=eth_account_keyfile=../${ENVIRONMENT}/files/keep-client-5-eth-keyfile.json \
  --from-file=eth_account_password=../../eth-networks/${ENVIRONMENT}/ropsten/eth-account-password.txt \
  --from-file=keep-client-config.toml=../${ENVIRONMENT}/files/keep-client-5-config.toml \
  --namespace ropsten \
  --dry-run \
  --save-config \
  -o yaml | kubectl apply -f -
