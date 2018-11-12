#!/bin/sh
set -e

export GETH_ETH_ACCOUNT=`cat /root/account1`
echo "GETH_ETH_ACCOUNT: $GETH_ETH_ACCOUNT"

export RANDOM_ID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
echo "RANDOM_ID: $RANDOM_ID"

export GETH_PARAMETERS='--port 3000 --networkid 1101 \
    --ws --wsport "8546" --wsorigins "*" \
    --rpc --rpcport "8545" --rpccorsdomain "" \
    --rpcapi "db,ssh,miner,admin,eth,net,web3,personal" \
    --syncmode "fast" \
    --mine --miner.threads=1'
export GETH_PARAMETERS="$GETH_PARAMETERS \\
    --identity $RANDOM_ID \\
    --miner.etherbase=$GETH_ETH_ACCOUNT"
echo "GETH_PARAMETERS: $GETH_PARAMETERS"

exec "/geth" $GETH_STATIC
