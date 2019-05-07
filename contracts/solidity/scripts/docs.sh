#!/bin/sh
shopt -s expand_aliases

alias flatten="solidity_flattener --solc-paths openzeppelin-solidity=$(pwd)/node_modules/openzeppelin-solidity"

mkdir -p docs/contracts docs/output docs/doxity/pages/docs

source ~/venv/bin/activate
flatten contracts/KeepToken.sol --output docs/contracts/KeepToken.sol
flatten contracts/Staking.sol --output docs/contracts/Staking.sol
flatten contracts/TokenGrant.sol --output docs/contracts/TokenGrant.sol

cd docs && node ../node_modules/@digix/doxity/lib/bin/doxity.js build
