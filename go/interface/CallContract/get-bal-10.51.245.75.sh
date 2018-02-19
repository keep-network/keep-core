
# get the balance from an account on the "testnet"
# VPN version

curl  \
-H "Content-Type: application/json" \
-X POST \
--data '{"jsonrpc":"2.0", "method":"eth_getBalance", "params":["0x6ffba2d0f4c8fd7961f516af43c55fe2d56f6044", "latest"], "id":1}' \
http://10.51.245.75:8545

# OLD: --data '{"jsonrpc":"2.0", "method":"eth_getBalance", "params":["0x7460a97b337214e319c7366520ee804236a848c8", "latest"], "id":1}' \

# From: https://ethereum.stackexchange.com/questions/6888/getting-contract-balance-over-json-rpc
# also has web example of how to do this.
# 
# Here is an example using the Linux utility curl to retrieve the balance of the contract over JSON-RPC:
# 
# ```bash
# 	curl -s -X POST --data '{"jsonrpc":"2.0", "method":"eth_getBalance", "params":["0xc5910bcb2442e84845aa98b20ca51e8f5d2bee23", "latest"], "id":1}' http://localhost:8545
# ```
#
# You will have to start geth with the --rpc parameter to run the curl program above, e.g.,
# 
# geth --testnet --rpc console
# The results you receive will be the number of weis in hexadecimal format:
# 
# {"jsonrpc":"2.0","id":1,"result":"0x58d15e17628000f"}
# Converting this number using a Hexadecimal to Decimal Converter results in a decimal number of 400000000000000015 weis.
# 
# To calculate the number of ethers, divide the decimal number by 1e18 and your result is 0.400000000000000015 ethers.
# 
# And checking this balance using the geth JavaScript console:
# 
# > web3.fromWei(eth.getBalance("0xc5910bcb2442e84845aa98b20ca51e8f5d2bee23"), "ether")
# > 0.400000000000000015
# Here is my Perl code to convert the hex wei number to a double ether number:
# 
# sub hexToDouble{
#   my $param = shift;
#   $param =~ s/^0x//;
#   my $fullnum = 0.0;
#   while ($param =~ /(.)/g) {
#     my $num = hex($1);
#     $fullnum = $fullnum * 16 + $num;
#   }
#   $fullnum *= 1e-18;
#   return $fullnum;
# }
# 
# 
#
# node.js version
#
# Use the official web3 JavaScript library's method web3.eth.getBalance
# 
# Example Node.js code:
# 
# Web3 = require('web3');
# 
# var web3 = new Web3(new Web3.providers.HttpProvider("http://localhost:8545"));
# 
# var balance = web3.eth.getBalance("0xc5910bcb2442e84845aa98b20ca51e8f5d2bee23");
# 
# console.log("Balance = " + web3.fromWei(balance, "ether"));
# You should change http://localhost:8545 to point to your own geth server.
# 
# Also the documentation linked above has some best practices for not overriding previously defined web3 objects, but I skipped them
# for brevity. Consult the docs yourself to see what all you can do with that lib.
# 
# If you cannot use Node.JS then you can read the Official Ethereum JSON RPC spec which tells you all of the methods available
# and how to use them. 
#
