#!/bin/bash

DATADIR=/root/.geth
DAGDIR=/root/.ethash
BLOCKNUMBER=1
GENESIS=/root/genesis.json

# initialize using the genesis block created by geth-init.sh
echo "-- Initialize. Write genesis block"
/geth --datadir $DATADIR init $GENESIS

# prebake the DAG for faster geth startup times
# details are available here:
# https://github.com/ethereum/wiki/wiki/Mining#ethash-dag
echo "-- Prebake the DAG for faster geth startup times"
/geth makedag $BLOCKNUMBER $DAGDIR
echo "-- Done"

echo "-- Check DATADIR contents"
ls -lah $DATADIR

echo "-- Check DAGDIR contents"
ls -lah $DAGDIR
