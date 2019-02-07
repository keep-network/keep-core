#!/bin/bash

DATADIR_DEFAULT=/root/.geth
ETH_IPC_PATH_DEFAULT=/root/.geth/geth.ipc

RPCPORT=8545
RPCHOST=0.0.0.0
RPCAPI=db,ssh,miner,admin,eth,net,web3,personal
WSPORT=8546
#WSHOST=0.0.0.0
#WSORIGINS="*"
GETHPORT=30303
GETHARGS=
BOOTNODE_URL="$BOOTNODE_URL/staticenodes?network=$BOOTNODE_NETWORK"
BOOTNODES=$(curl --connect-timeout 1 --retry 10  --retry-max-time 10 -f -s $BOOTNODE_URL)

# fetch accounts
export GETH_ETH_MINING_ACCOUNT=`cat /root/mining_account`
echo "-- GETH_ETH_MINING_ACCOUNT: $GETH_ETH_MINING_ACCOUNT"

# dump genesis file
echo "-- Dump genesis.json:"
GENESIS=/root/genesis.json
cat $GENESIS
echo ""

if [ -z "$HOSTVOLUME" ]; then
  DATADIR="$DATADIR_DEFAULT"
  echo "-- No HOSTVOLUME was supplied. Using default DATADIR: $DATADIR"
else
  DATADIR="$HOSTVOLUME" # GCP: each pod has a private volume attached
  echo "-- Setting DATADIR to: $DATADIR"
  # check if we need to create the directory
  if [ ! -d "$DATADIR" ]; then
    echo "-- Creating $DATADIR"
    mkdir -p $DATADIR
  fi
  echo "-- Copying keystore to DATADIR"
  cp -rv $DATADIR_DEFAULT/keystore $DATADIR
  echo "-- List DATADIR/keystore:"
  ls -la $DATADIR/keystore
fi

if [ -z "$ETH_IPC_PATH" ]; then
  ETH_IPC_PATH="$ETH_IPC_PATH_DEFAULT"
  echo "-- No ETH_IPC_PATH was supplied. Using default ETH_IPC_PATH: $ETH_IPC_PATH"
fi

if [ -z "$NETWORKID" ]; then
  echo "-- No NETWORKID was supplied"
  exit 1
fi

if [ -z "$GENESIS" ]; then
  echo "-- No GENESIS  was supplied"
  exit 1
fi

if [ -z "$NODE_NAME" ]; then
  echo "-- No NODE_NAME was supplied"
  exit 1
fi

if [ "$ENABLE_MINER" ]; then
  MINER_ADDRESS=$GETH_ETH_MINING_ACCOUNT
  echo "-- MINER_ADDRESS: $MINER_ADDRESS"

  while [ -z "$BOOTNODES" ]
  do
    BOOTNODES=$(curl --connect-timeout 1 --retry 10  --retry-delay 0 --retry-max-time 10 -f -s $BOOTNODE_URL)
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

# TODO: only initialize if DATADIR has no chain data
if [ ! -d "$DATADIR/geth/chaindata" ]; then
  echo "-- No chaindata directory. Neet to Initialize. Writing genesis block..."
  /geth --datadir $DATADIR init $GENESIS
fi

echo "-- BOOTNODES: $BOOTNODES"
echo "-- GETHARGS: $GETHARGS"

echo "-- Starting geth..."

/geth --datadir $DATADIR --ethash.dagdir $DATADIR --ipcpath $ETH_IPC_PATH \
      --nodiscover \
      --port $GETHPORT --networkid $NETWORKID \
      --ws --wsaddr "0.0.0.0" --wsport $WSPORT --wsorigins "*" \
      --rpc --rpcport $RPCPORT --rpcaddr $RPCHOST --rpccorsdomain "" \
      --rpcapi $RPCAPI \
      --identity $NODE_NAME \
      --syncmode "fast" \
      $GETHARGS
