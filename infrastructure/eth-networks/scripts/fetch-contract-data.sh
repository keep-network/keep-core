#!/bin/bash

# For ropsten contract data fetch this script assumes the last migration run from your local
# machine was against ropsten
set -e

HELP="Usage: ./$(basename $0) -n <ETH_NETWORK> -e <ENVIRONMENT>
      \n\nAvailable ETH_NETWORK: ropsten, internal
      \nAvailable ENVIRONMENT: keep-dev,keep-test"

while getopts ":n:e:" opt; do
  case $opt in
    n ) ETH_NETWORK=$OPTARG;;
    e ) ENVIRONMENT=$OPTARG;;

    \?)
      echo -e $HELP
      exit 1
  esac
done

if [ $# -eq 0 ]
then
  echo -e $HELP
  exit 1
fi

TRUFFLE_MIGRATED_CONTRACT_PATH=../../../contracts/solidity/build/contracts

KEEP_DEV_ROPSTEN_DIR=../keep-dev/ropsten

KEEP_TEST_PROJECT_ID=keep-test-f3e0
KEEP_TEST_CONTRACT_BUCKET=keep-test-contract-data/keep-core
KEEP_TEST_INTERNAL_DIR=../keep-test/internal
KEEP_TEST_ROPSTEN_DIR=../keep-test/ropsten


if [ $ETH_NETWORK = ropsten ] && [ $ENVIRONMENT = keep-dev ]
then
  cp ${TRUFFLE_MIGRATED_CONTRACT_PATH}/KeepToken.json $KEEP_DEV_ROPSTEN_DIR
  cp ${TRUFFLE_MIGRATED_CONTRACT_PATH}/TokenStaking.json $KEEP_DEV_ROPSTEN_DIR

elif [ $ETH_NETWORK = ropsten ] && [ $ENVIRONMENT = keep-test ]
then
  cp ${TRUFFLE_MIGRATED_CONTRACT_PATH}/KeepToken.json $KEEP_TEST_ROPSTEN_DIR
  cp ${TRUFFLE_MIGRATED_CONTRACT_PATH}/TokenStaking.json $KEEP_TEST_ROPSTEN_DIR

elif [ $ETH_NETWORK = internal ] && [ $ENVIRONMENT = keep-test ]
then
  CURRENT_PROJECT=$(gcloud config get-value project)

  if [ $CURRENT_PROJECT != $KEEP_TEST_PROJECT_ID ]
  then
    echo "--Current gcloud project: ${CURRENT_PROJECT}"
    echo "--Setting your gcloud context to the keep-test project!"
    gcloud config set project $KEEP_TEST_PROJECT_ID
    gsutil cp gs://${KEEP_TEST_CONTRACT_BUCKET}/KeepToken.json $KEEP_TEST_INTERNAL_DIR
    gsutil cp gs://${KEEP_TEST_CONTRACT_BUCKET}/TokenStaking.json $KEEP_TEST_INTERNAL_DIR
    echo "--Returning to original glcoud project: ${CURRENT_PROJECT}"
    gcloud config set project $CURRENT_PROJECT

  else
    gsutil cp gs://${KEEP_TEST_CONTRACT_BUCKET}/KeepToken.json $KEEP_TEST_INTERNAL_DIR
    gsutil cp gs://${KEEP_TEST_CONTRACT_BUCKET}/TokenStaking.json $KEEP_TEST_INTERNAL_DIR
  fi

else
  echo -e $HELP
fi


