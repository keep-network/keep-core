#!/bin/bash

DATADIR_DEFAULT=/root/.geth

usage()
{
    echo "usage: geth-init.sh <KEEP_ACCOUNTS>"
}

generate_Keep_account()
{
  ACCOUNT=$1
  # create a new Ethereum account. We are prepending "0x".
  echo "-- Generating account $ACCOUNT"
  GETH_ETH_ACCOUNT=`/geth --datadir $DATADIR_DEFAULT account new  \
    --password /root/passphrase | \
    cut -d "{" -f2 | cut -d "}" -f1 | \
    while read line; do echo "0x$line"; done`
  echo "$GETH_ETH_ACCOUNT" >> /root/keep_accounts
  echo "   GETH_ETH_ACCOUNT: $GETH_ETH_ACCOUNT"
}


# parse command line parameter
if [ "$1" != "" ]; then
  re='^[0-9]+$'
  if ! [[ $1 =~ $re ]] ; then
     echo "error: Not a number" >&2; exit 1
  fi
  KEEP_ACCOUNTS=$1
  echo "Will generate $KEEP_ACCOUNTS Keep client accounts."
else
    echo "error: KEEP_ACCOUNTS undefined!"
    usage
    exit 1
fi

# create mining account. We are prepending "0x".
echo "-- Generating mining_account"
GETH_ETH_MINING_ACCOUNT=`/geth --datadir $DATADIR_DEFAULT account new  \
    --password /root/passphrase | \
    cut -d "{" -f2 | cut -d "}" -f1 | \
    while read line; do echo "0x$line"; done`
echo "$GETH_ETH_MINING_ACCOUNT"  >> /root/mining_account
echo "-- GETH_ETH_MINING_ACCOUNT: $GETH_ETH_MINING_ACCOUNT"

# create Keep client accounts
for i in `seq 1 $KEEP_ACCOUNTS`;
do
  generate_Keep_account $i
done
echo "-- Keep client accounts generated:"
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
  LINE=".alloc += {\"$ACCOUNT\": {"balance": \"1000000000000000000000\"}}"
  RESULT=`echo $RESULT | jq "$LINE"`
done < /root/keep_accounts
echo $RESULT | jq . > $GENESIS

## genesis.json generation done -------

# dump genesis file
echo "-- Dump genesis.json:"
cat $GENESIS
echo ""

# List the keystore directory
echo "-- KEYSTORE directory"
ls -la $DATADIR_DEFAULT/keystore

# List the .geth directory
echo "-- .geth directory"
ls -la $DATADIR_DEFAULT

# List the /root directory
echo "-- /root directory"
ls -la /root
