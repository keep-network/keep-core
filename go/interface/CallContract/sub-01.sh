
# get the balance from an account on the "testnet"
# VPN version

curl  \
-H "Content-Type: application/json" \
-X POST \
--data '{"method":"eth_subscribe","params":["newHeads",{"fromBlock":"latest","toBlock":"latest"}],"id":1,"jsonrpc":"2.0"}'  \
http://10.51.245.75:8545

#--data '{"jsonrpc":"2.0", "method":"eth_getBalance", "params":["0x6ffba2d0f4c8fd7961f516af43c55fe2d56f6044", "latest"], "id":1}' \
