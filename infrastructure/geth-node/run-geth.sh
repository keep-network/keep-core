#!/bin/bash

DATADIR=/root/.geth

# fetch accounts
export GETH_ETH_ACCOUNT0=`cat /root/account0`
echo "-- GETH_ETH_ACCOUNT0: $GETH_ETH_ACCOUNT0"
export GETH_ETH_ACCOUNT1=`cat /root/account1`
echo "-- GETH_ETH_ACCOUNT1: $GETH_ETH_ACCOUNT1"

# dump genesis file
echo "-- Dump genesis.json:"
cat /root/genesis.json
echo ""

GENESIS=/root/genesis.json
RPCPORT=8545
RPCHOST=0.0.0.0
RPCAPI="db,ssh,miner,admin,eth,net,web3,personal"
WSPORT=8546
WSHOST=0.0.0.0
WSORIGINS="*"
GETHPORT=30303
GETHARGS=
BOOTNODE_URL="$BOOTNODE_URL/staticenodes?network=$BOOTNODE_NETWORK"
BOOTNODES=$(curl --connect-timeout 1 --retry 10  --retry-max-time 10 -f -s $BOOTNODE_URL)
STATSARGS="--ethstats \"$NODE_NAME:$WS_SECRET@$WS_SERVER\""

if [ -z "$NETWORKID" ]; then
  echo "No NETWORKID was supplied"
  exit 1
fi

if [ -z "$GENESIS" ]; then
  echo "No GENESIS  was supplied"
  exit 1
fi

if [ -z "$NODE_NAME" ]; then
  echo "No NODE_NAME was supplied"
  exit 1
fi

if [ "$ENABLE_MINER" ]; then
  MINER_ADDRESS=$GETH_ETH_ACCOUNT0
  echo "-- MINER_ADDRESS: $MINER_ADDRESS"

  while [ -z "$BOOTNODES" ]
  do
    BOOTNODES=$(curl --connect-timeout 1 --retry 10  --retry-max-time 10 -f -s $BOOTNODE_URL)
  done

  GETHARGS="--mine --miner.etherbase=$MINER_ADDRESS"

  if [ "$MINER_THREADS" ]; then
    GETHARGS="$GETHARGS --minerthreads $MINER_THREADS"
  fi
else
  GETHARGS=""
fi


if [ "$BOOTNODES" ]; then
  echo "-- Adding bootnodes:"
  mkdir -p $DATADIR
  echo $BOOTNODES > $DATADIR/static-nodes.json
  cat $DATADIR/static-nodes.json
fi

if [ ! -d "$DATADIR/chaindata" ]; then
  echo "-- Initialize. Write genesis block"
  /geth --datadir $DATADIR init $GENESIS
fi

echo "-- BOOTNODES: $BOOTNODES"
echo "-- GETHARGS:  $GETHARGS"

/geth --datadir $DATADIR \
    --nodiscover \
    --port 30303 --networkid $NETWORKID \
    --ws --wsaddr "0.0.0.0" --wsport 8546 --wsorigins "*" \
    --rpc --rpcport 8545 --rpcaddr 0.0.0.0 --rpccorsdomain "" \
    --rpcapi "db,ssh,miner,admin,eth,net,web3,personal" \
    --identity $NODE_NAME \
    --syncmode "fast" \
    $GETHARGS
