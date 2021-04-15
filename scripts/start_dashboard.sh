#!/bin/bash
set -e

# Dafault inputs.
KEEP_ECDSA_SOL_PATH="$PWD/../keep-ecdsa/solidity"
KEEP_ECDSA_DISTRIBUTOR_PATH="$PWD/../keep-ecdsa/staker-rewards/distributor"
MERKLE_DISTRIBUTOR_INPUT_PATH="$KEEP_ECDSA_DISTRIBUTOR_PATH/staker-reward-allocation.json"
MERKLE_DISTRIBUTOR_OUTPUT_PATH="$KEEP_ECDSA_DISTRIBUTOR_PATH/output-merkle-objects.json"

TBTC_SOL_PATH="$PWD/../tbtc/solidity"
TBTC_SOL_ARTIFACTS_PATH="$TBTC_SOL_PATH/build/contracts"

KEEP_CORE_PATH=$PWD
KEEP_CORE_SOL_PATH="$KEEP_CORE_PATH/solidity"
DASHBOARD_DIR_PATH="$KEEP_CORE_SOL_PATH/dashboard"
KEEP_CORE_SOL_ARTIFACTS_PATH="$KEEP_CORE_SOL_PATH/build/contracts"

# Run script.
LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color

cd $KEEP_CORE_SOL_PATH

## uncomment when version of a dependency in package.json has changed.
# printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
# rm -f package-lock.json
# rm -rf node_modules/
npm install

printf "${LOG_START}Migrating contracts for Keep-Core...${LOG_END}"
rm -rf build/
truffle migrate --reset --network sov
printf "${LOG_START}Delegating tokens...${LOG_END}"
truffle exec ./scripts/delegate-tokens.js --network sov

cd $TBTC_SOL_PATH

printf "${LOG_START}Migrating contracts for tBTC...${LOG_END}"
npm run clean
truffle migrate --reset --network development

printf "${LOG_START}Creating symlinks for tBTC...${LOG_END}"
rm -f artifacts
ln -s build/contracts artifacts
npm link

cd $KEEP_ECDSA_SOL_PATH

output=$(truffle exec ./scripts/get-network-id.js --network sov)
NETWORKID=$(echo "$output" | tail -1)
printf "Current network ID: ${NETWORKID}\n"

printf "${LOG_START}Provisioning Keep-Ecdsa...${LOG_END}"
KEEP_CORE_SOL_ARTIFACTS_PATH=$KEEP_CORE_SOL_ARTIFACTS_PATH \
NETWORKID=$NETWORKID \
    ./scripts/lcl-provision-external-contracts.sh

printf "${LOG_START}Provisioning TBTC...${LOG_END}"
TBTC_SOL_ARTIFACTS_PATH=$TBTC_SOL_ARTIFACTS_PATH \
NETWORKID=$NETWORKID \
    ./scripts/lcl-provision-tbtc.sh

printf "${LOG_START}Migrating contracts for Keep-Ecdsa...${LOG_END}"
npm run clean
truffle migrate --reset --network sov

printf "${LOG_START}Creating symlinks for Keep-Ecdsa...${LOG_END}"
rm -f artifacts
ln -s build/contracts artifacts
npm link

printf "${LOG_START}Initializing Keep-Ecdsa...${LOG_END}"
truffle exec scripts/lcl-initialize.js --network sov

cd $DASHBOARD_DIR_PATH

## uncomment when version of a dependency in package.json has changed.
# printf "${LOG_START}Installing NPM dependencies in dashboard...${LOG_END}"
# rm -rf node_modules/
# rm package-lock.json
npm install

cd $KEEP_CORE_SOL_PATH

printf "${LOG_START}Creating symlinks for Keep-Core...${LOG_END}"
rm -f artifacts
ln -s build/contracts artifacts
npm link

cd $DASHBOARD_DIR_PATH
npm link @keep-network/keep-core
npm link @keep-network/keep-ecdsa
npm link @keep-network/tbtc

# printf "${LOG_START}Starting dashboard...${LOG_END}"
# npm start

# Make sure files below exists in keep-ecdsa repository. Otherwise comment out.
printf "${LOG_START}Generating mock input data for ecdsa merkle distributor${LOG_END}"
cd $KEEP_ECDSA_SOL_PATH
truffle exec ./scripts/generate-staker-rewards-input.js --network sov

printf "${LOG_START}Generating mock merkle objects${LOG_END}"
cd $KEEP_ECDSA_DISTRIBUTOR_PATH
npm i
npm run generate-merkle-root -- --input="$MERKLE_DISTRIBUTOR_INPUT_PATH"

printf "${LOG_START}Copying the mock merkle objects to dashboard${LOG_END}"
cp $MERKLE_DISTRIBUTOR_OUTPUT_PATH "$DASHBOARD_DIR_PATH/src/rewards-allocation/rewards.json"

printf "${LOG_START}Initializing ECDSARewardsDistributor contract${LOG_END}"
cd $KEEP_ECDSA_SOL_PATH
truffle exec ./scripts/initialize-ecdsa-rewards-distributor.js --network sov
