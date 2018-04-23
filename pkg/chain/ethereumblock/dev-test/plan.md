
1. Shows how to wait for blocks [using a filter in JS / web3][https://ethereum.stackexchange.com/questions/9636/whats-the-proper-way-to-wait-for-a-transaction-to-be-mined-and-get-the-results]

2. Go example to [filter][https://ethereum.stackexchange.com/questions/22954/mistake-when-using-web3-eth-filter-and-filter-get]

3. JSON RPC Post to do [get_logs][https://github.com/ethereum/go-ethereum/issues/15091]

https://ethereum.stackexchange.com/questions/21694/using-web3-eth-filter


{"id": 1, "method": "eth_subscribe", "params": ["syncing"]}
	can also be "pending" and "latest" instead of "syncing"
	https://github.com/ethereum/go-ethereum/wiki/RPC-PUB-SUB


https://github.com/ethereum/go-ethereum/blob/master/ethclient/ethclient.go
	- has get block by number
	
// +build ignore

	if false {
		var latestBlock Block
		if err := client.CallContext(ctx, &latestBlock, "eth_getBlockByNumber", "latest"); err != nil {
			fmt.Println("Cannot get the latest Block => ", err)
			return
		}
		subch <- latestBlock
	}
