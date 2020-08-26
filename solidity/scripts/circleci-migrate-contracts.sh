#!/bin/bash

set -e

if [[ -z $GOOGLE_PROJECT_NAME || -z $BUILD_TAG || -z $TRUFFLE_NETWORK || -z $TENDERLY_TOKEN || -z $ETH_NETWORK_ID ]]; then
  echo "one or more required variables are undefined"
  exit 1
fi

mkdir -p /tmp/keep-client/contracts
cd ./solidity
npm ci

echo "<<<<<<START Contract Migration START<<<<<<"
./node_modules/.bin/truffle migrate --reset --network $TRUFFLE_NETWORK
cp ./build/contracts/* /tmp/keep-client/contracts
echo ">>>>>>FINISH Contract Migration FINISH>>>>>>"

echo "<<<<<<START Tenderly Push START<<<<<<"
tenderly login --authentication-method token --token $TENDERLY_TOKEN
tenderly push --networks $ETH_NETWORK_ID --tag keep-core \ 
  --tag $GOOGLE_PROJECT_NAME --tag $BUILD_TAG || echo "tendery push failed :("
echo "<<<<<<FINISH Tenderly Push FINISH<<<<<<"
