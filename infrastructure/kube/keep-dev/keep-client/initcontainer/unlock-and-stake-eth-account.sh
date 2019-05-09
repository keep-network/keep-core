#!/bin/sh

function unlock_eth_account() {
  geth --exec "personal.unlockAccount(\"$ETHEREUM_ACCOUNT\", \"$ETHEREUM_PASSWORD\", 700)" attach http://eth-tx-node.default.svc.cluster.local:8545
}

function stake_eth_account() {
  geth --exec 'loadScript("./stake-eth-account.js");' attach http://eth-tx-node.default.svc.cluster.local:8545
}

unlock_eth_account
