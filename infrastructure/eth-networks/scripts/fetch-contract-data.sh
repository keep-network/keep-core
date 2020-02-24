#!/bin/bash
set -e

CONTRACTS=( KeepToken.json TokenStaking.json KeepRandomBeaconOperator.json )

SOURCE_PROJECT_ID=keep-test-f3e0
SOURCE_CONTRACT_BUCKET=keep-test-contract-data/keep-core
CURRENT_PROJECT=$(gcloud config get-value project)

TARGET_DIR=../keep-test/ropsten


  if [ $CURRENT_PROJECT != $SOURCE_PROJECT_ID ]
  then
    echo "--Current gcloud project: ${CURRENT_PROJECT}"
    echo "--Setting your gcloud context to the keep-test project!"
    gcloud config set project $KEEP_TEST_PROJECT_ID
    for contract in "${CONTRACTS[@]}"; do
      gsutil cp gs://${SOURCE_CONTRACT_BUCKET}/$contract $TARGET_DIR
    done
    echo "--Returning to original glcoud project: ${CURRENT_PROJECT}"
    gcloud config set project $CURRENT_PROJECT
  else
    for contract in "${CONTRACTS[@]}"; do
      gsutil cp gs://${SOURCE_CONTRACT_BUCKET}/$contract $TARGET_DIR
    done
  fi
