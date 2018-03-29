

local geth truffle version

curl  \
-H "Content-Type: application/json" \
-X POST \
--data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' \
http://127.0.0.1:9545
