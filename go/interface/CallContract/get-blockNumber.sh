
# get the balance from an account on the "testnet"
# VPN version

curl  \
-H "Content-Type: application/json" \
-X POST \
--data '{"jsonrpc":"2.0", "method":"eth_blockNumber", "params":[], "id":992}' \
http://10.51.245.75:8545

