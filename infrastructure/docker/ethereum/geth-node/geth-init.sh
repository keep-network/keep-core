#!/bin/bash

DATADIR_DEFAULT=/root/.geth

# feed keep ETH accounts into genesis
for keyfile in ${DATADIR_DEFAULT}/keystore/*;
do
  ACCOUNT=`cat $keyfile | jq .address | tr -d '"'`
  echo "0x${ACCOUNT}" >> /root/keep_accounts
done
echo "-- Keep client accounts populated:"
cat /root/keep_accounts

# Generate genesis.json and issue tokens to Keep peers.
# We are setting mining difficulty to zero.
# Start with the preamble
# Generate genesis.json and issue tokens to Keep peers.
# We are setting mining difficulty to zero.
GENESIS=/root/genesis.json
GENESIS_TEMPLATE=/root/genesis-template.json
RESULT=`cat $GENESIS_TEMPLATE`
# add Keep client accounts and fund them
while read -r ACCOUNT; do
  LINE=".alloc += {\"${ACCOUNT}\": {"balance": \"1000000000000000000000\"}}"
  RESULT=`echo $RESULT | jq "$LINE"`
done < /root/keep_accounts
echo $RESULT | jq . > $GENESIS

## genesis.json generation done -------

## set miner account
echo "-- Setting miner account:"
head -1 /root/keep_accounts > /root/mining_account
cat /root/mining_account

# dump genesis file
echo "-- Dump genesis.json:"
cat $GENESIS
echo ""

# List the keystore directory
echo "-- KEYSTORE directory"
ls -la ${DATADIR_DEFAULT}/keystore

# List the .geth directory
echo "-- .geth directory"
ls -la $DATADIR_DEFAULT

# List the /root directory
echo "-- /root directory"
ls -la /root
