#!/bin/bash

DATADIR=/root/.geth

# create new account for this Keep client. We are prepending "0x".
/geth --datadir $DATADIR account new  --password /root/passphrase | \
    cut -d "{" -f2 | cut -d "}" -f1 | \
    while read line; do echo "0x$line"; done >> /root/account0
echo "---Generating account0"

# create new account for Keep peers. We are prepending "0x".
/geth --datadir $DATADIR account new  --password /root/passphrase | \
    cut -d "{" -f2 | cut -d "}" -f1 | \
    while read line; do echo "0x$line"; done >> /root/account1
echo "---Generating account1"

export GETH_ETH_ACCOUNT0=`cat /root/account0`
echo "-- GETH_ETH_ACCOUNT0: $GETH_ETH_ACCOUNT0"
export GETH_ETH_ACCOUNT1=`cat /root/account1`
echo "-- GETH_ETH_ACCOUNT1: $GETH_ETH_ACCOUNT1"

# Generate genesis.json and issue tokens to Keep peers account1.
# We are setting mining difficulty to zero.
cat <<EOF >> /root/genesis.json
{
  "config": {
    "chainId": 1101,
    "homesteadBlock": 0,
    "byzantiumBlock": 0,
    "eip155Block": 0,
    "eip158Block": 0
  },
  "difficulty" : "0x0",
  "gasLimit"   : "0x493E00",
  "alloc": {
EOF

echo "   \"$GETH_ETH_ACCOUNT1\": {" >> /root/genesis.json
echo "    \"balance\": \"1000000000000000000000\"" >> /root/genesis.json
cat <<EOF >> /root/genesis.json
   }
  }
}
EOF

# dump genesis file
echo "-- Dump genesis.json:"
cat /root/genesis.json
echo ""
