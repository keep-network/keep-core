#!/bin/bash
set -e

# Dafault inputs.
KEEP_CORE_PATH=$PWD
KEEP_CORE_SOL_PATH="$KEEP_CORE_PATH/solidity"
DASHBOARD_DIR_PATH="$KEEP_CORE_PATH/solidity/dashboard"

# Run script.
LOG_START='\n\e[1;36m' # new line + bold + color
LOG_END='\n\e[0m' # new line + reset color

cd $KEEP_CORE_SOL_PATH

printf "${LOG_START}Installing NPM dependencies...${LOG_END}"
rm -rf node_modules/
npm install

printf "${LOG_START}Migrating contracts...${LOG_END}"
rm -rf build/
rm package-lock.json
truffle migrate --reset --network local
truffle exec ./scripts/delegate-tokens.js --network local

cd $DASHBOARD_DIR_PATH

printf "${LOG_START}Installing NPM dependencies in dashboard...${LOG_END}"
rm -rf node_modules/
rm package-lock.json
npm install

cd $KEEP_CORE_SOL_PATH

printf "${LOG_START}Creating symlinks...${LOG_END}"
rm -f artifacts
ln -s build/contracts artifacts
npm link

cd $DASHBOARD_DIR_PATH
npm link @keep-network/keep-core

printf "${LOG_START}Starting dashboard...${LOG_END}"
npm start
