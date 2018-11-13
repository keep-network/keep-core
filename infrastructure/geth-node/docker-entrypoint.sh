#!/bin/sh
set -e

export GETH_ETH_ACCOUNT=`cat /root/account1`
echo "GETH_ETH_ACCOUNT: $GETH_ETH_ACCOUNT"

export RANDOM_ID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
echo "RANDOM_ID: $RANDOM_ID"

exec "/geth" --port 30303 --networkid 1101 \
    --ws --wsaddr "0.0.0.0" --wsport 8546 --wsorigins "*" \
    --rpc --rpcport 8545 --rpcaddr 0.0.0.0 --rpccorsdomain "" \
    --rpcapi "db,ssh,miner,admin,eth,net,web3,personal" \
    --syncmode "fast" \
    --mine --miner.threads=1 \
    --identity $RANDOM_ID \
    --miner.etherbase=$GETH_ETH_ACCOUNT
