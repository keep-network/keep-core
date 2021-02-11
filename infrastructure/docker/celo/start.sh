#!/bin/sh

mkdir -p $CELO_DATA_DIR/keystore
cp $CELO_INIT_DIR/keystore/* $CELO_DATA_DIR/keystore
cp $CELO_INIT_DIR/password.txt $CELO_DATA_DIR/password.txt
[ ! -d "$CELO_DATA_DIR/celo" ] && geth --nousb --datadir=$CELO_DATA_DIR init $CELO_INIT_DIR/genesis.json

 geth --port 3000 --networkid 1101 --identity "somerandomidentity" \
    --ws --wsaddr "0.0.0.0" --wsport "8546" --wsorigins "*" \
    --rpc --rpcport "8545" --rpcaddr "0.0.0.0" --rpccorsdomain "" \
    --rpcapi "db,ssh,miner,admin,eth,net,web3,personal" \
    --wsapi "db,ssh,miner,admin,eth,net,web3,personal" \
    --datadir $CELO_DATA_DIR --syncmode "fast" \
    --mine --miner.threads 1 --nousb \
    --unlock 0x2b2976824682233807a197081119da511af12f7a --password $CELO_DATA_DIR/password.txt \
    --allow-insecure-unlock 