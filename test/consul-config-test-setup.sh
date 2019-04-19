#!/bin/sh

# to make the code in config_test.go work you need to run the following commands
# on your Mac:
# brew install consul
# brew services start consul
#
# Alternatively you can run Consul inside a pod on your Kubernetes cluster
#
consul kv put ethereum '{"URL": "ws://192.168.0.150:8546", "URLRPC": "http://192.168.0.151:8545"}'
consul kv get ethereum

consul kv put ethereum/account '{"Address": "0xc2a56884538778bacd91aa5bf343bf882c5fb18c", "KeyFile": "/tmp/UTC--2018-03-11T01-37-33.202765887Z--c2a56884538778bacd91aa5bf343bf882c5fb18c"}'
consul kv get ethereum/account

consul kv put ethereum/contractaddresses '{"KeepRandomBeacon": "0x639deb0dd975af8e4cc91fe9053a37e4faf37648", "KeepGroup": "0xcf64c2a367341170cb4e09cf8c0ed137d8473cec", "StakingProxy": "0xCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCD"}'
consul kv get ethereum/contractaddresses

consul kv put libp2p '{"Port": 27002, "Peers": ["/ip4/127.0.0.1/tcp/27001/ipfs/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy49tpPV95XdKu1BcvMA"], "Seed": 0}'
consul kv get libp2p
