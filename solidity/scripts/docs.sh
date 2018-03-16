#!/bin/sh
shopt -s expand_aliases
alias flatten="solidity_flattener --solc-paths zeppelin-solidity=$(pwd)/node_modules/zeppelin-solidity"

mkdir -p docs/contracts docs/output docs/doxity/pages/docs

flatten contracts/KeepToken.sol --output docs/contracts/KeepToken.sol
flatten contracts/TokenStaking.sol --output docs/contracts/TokenStaking.sol
flatten contracts/TokenGrant.sol --output docs/contracts/TokenGrant.sol

cd docs && node ../node_modules/@digix/doxity/lib/bin/doxity.js build
