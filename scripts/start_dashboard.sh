#!/bin/bash
set -e

LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color

KEEP_CORE_PATH=$PWD
KEEP_CORE_SOL_PATH="$PWD/solidity-v1"
DASHBOARD_DIR_PATH="$KEEP_CORE_PATH/solidity-v1/dashboard"
KEEP_CORE_ARTIFACTS_PATH="$KEEP_CORE_PATH/solidity-v1/artifacts"

KEEP_ECDSA_PATH="$PWD/../keep-ecdsa"
KEEP_ECDSA_SOL_PATH="$PWD/../keep-ecdsa/solidity"

KEEP_ECDSA_DISTRIBUTOR_PATH="$PWD/../keep-ecdsa/staker-rewards/distributor"
MERKLE_DISTRIBUTOR_INPUT_PATH="$KEEP_ECDSA_DISTRIBUTOR_PATH/staker-reward-allocation.json"
MERKLE_DISTRIBUTOR_OUTPUT_PATH="$KEEP_ECDSA_DISTRIBUTOR_PATH/output-merkle-objects.json"

TBTC_PATH="$PWD/../tbtc"
TBTC_SOL_ARTIFACTS_PATH="$TBTC_PATH/solidity/artifacts"

COV_POOLS_PATH="$PWD/../coverage-pools"

THRESHOLD_CONTRACTS_PATH="$PWD/../solidity-contracts"

printf "${LOG_START}Migrating contracts for keep-core...${LOG_END}"
cd "$KEEP_CORE_PATH"
./scripts/install.sh --network local --contracts-only
cd "$KEEP_CORE_SOL_PATH"
# Link keep-core contracts via `yarn`- the threshold solidity contracts repo
# uses the `yarn` so we need to link keep-core package with `yarn` as well.
yarn link

printf "${LOG_START}Migrating contracts for keep-ecdsa...${LOG_END}"
cd "$KEEP_ECDSA_PATH"
./scripts/install-v1.sh --network local

printf "${LOG_START}Migrating contracts for tBTC...${LOG_END}"
cd "$TBTC_PATH"
./scripts/install.sh

printf "${LOG_START}Initialize contracts for keep-ecdsa...${LOG_END}"
cd $KEEP_ECDSA_SOL_PATH

# Get network ID.
NETWORK_ID_OUTPUT=$(npx truffle exec ./scripts/get-network-id.js --network local)
NETWORK_ID=$(echo "$NETWORK_ID_OUTPUT" | tail -1)
GET_NETWORK_JSON_QUERY=".networks.\"${NETWORK_ID}\".address"

# Extract TBTCSystem contract address.
TBTC_SYSTEM_CONTRACT="${TBTC_SOL_ARTIFACTS_PATH}/TBTCSystem.json"
TBTC_SYSTEM_CONTRACT_ADDRESS=$(cat ${TBTC_SYSTEM_CONTRACT} | jq "${GET_NETWORK_JSON_QUERY}" | tr -d '"')

printf "${LOG_START}TBTCSystem contract address is: ${TBTC_SYSTEM_CONTRACT_ADDRESS}${LOG_END}"

# Run keep-ecdsa initialization script.
cd $KEEP_ECDSA_PATH
./scripts/initialize.sh --contracts-only --network local --application-address $TBTC_SYSTEM_CONTRACT_ADDRESS

printf "${LOG_START}KEEP token contract address is: ${KEEP_TOKEN_CONTRACT_ADDRESS}${LOG_END}"

cd $COV_POOLS_PATH
printf "${LOG_START}Creating links for cvovrage pools...${LOG_END}"
./scripts/install.sh

# In the Keep Token Dashboard we use `npm` instead of `yarn` so we need to link
# the `keep-network/coverage-pool` package manually via npm.
npm link

cd "$THRESHOLD_CONTRACTS_PATH"
printf "${LOG_START}Installing Threshold network solidity contracts dependencies...${LOG_START}"
yarn
yarn link @keep-network/keep-core
./scripts/prepare-dependencies.sh
printf "${LOG_START}Deploying contracts for threshold solidity contracts...${LOG_START}"
yarn deploy --network development --reset
./scripts/prepare-artifacts.sh --network development
npm link

printf "${LOG_START}Installing Keep Token Dashboard...${LOG_END}"

cd $DASHBOARD_DIR_PATH

# uncomment when version of a dependency in package.json has changed.
printf "${LOG_START}Installing NPM dependencies in dashboard...${LOG_END}"
# rm -rf node_modules/
# rm package-lock.json
npm install

cd $DASHBOARD_DIR_PATH
npm link @keep-network/keep-core \
    @keep-network/keep-ecdsa \
    @keep-network/tbtc \
    @keep-network/coverage-pools \
    @threshold-network/solidity-contracts

# Make sure files below exists in keep-ecdsa repository. Otherwise comment out.
# printf "${LOG_START}Generating mock input data for ecdsa merkle distributor${LOG_END}"
# cd $KEEP_ECDSA_SOL_PATH
# truffle exec ./scripts/generate-staker-rewards-input.js --network local

# printf "${LOG_START}Generating mock merkle objects${LOG_END}"
# cd $KEEP_ECDSA_DISTRIBUTOR_PATH
# npm i
# npm run generate-merkle-root -- --input="$MERKLE_DISTRIBUTOR_INPUT_PATH"

# printf "${LOG_START}Copying the mock merkle objects to dashboard${LOG_END}"
# cp $MERKLE_DISTRIBUTOR_OUTPUT_PATH "$DASHBOARD_DIR_PATH/src/rewards-allocation/rewards.json"

# printf "${LOG_START}Initializing ECDSARewardsDistributor contract${LOG_END}"
# cd $KEEP_ECDSA_SOL_PATH
# truffle exec ./scripts/initialize-ecdsa-rewards-distributor.js --network local

# printf "${LOG_START}Starting dashboard...${LOG_END}"
# cd $DASHBOARD_DIR_PATH
# npm start
