#!/bin/sh

mkdir -p /data/keystore
cp /celo-init/keystore/* /mnt/data/keystore
cp /celo-init/password.txt /mnt/data/password.txt
[ ! -d "/mnt/data/celo" ] && geth --nousb --datadir=/mnt/data init /celo-init/genesis.json

 geth --port 3000 --networkid 1101 --identity "somerandomidentity" \
    --ws --wsaddr "0.0.0.0" --wsport "8546" --wsorigins "*" \
    --rpc --rpcport "8545" --rpcaddr "0.0.0.0" --rpccorsdomain "" \
    --rpcapi "db,ssh,miner,admin,eth,net,web3,personal" \
    --wsapi "db,ssh,miner,admin,eth,net,web3,personal" \
    --datadir /mnt/data --syncmode "fast" \
    --mine --miner.threads 1 --nousb \
    --unlock 0x2b2976824682233807a197081119da511af12f7a --password /mnt/data/password.txt \
    --allow-insecure-unlock 