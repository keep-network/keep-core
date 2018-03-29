Contract Call / Capture Event Test
=========================================

You will need to have  a version of geth running locally to support IPC or
running with websocket enabled.  The regular testnet lacks web sockets.  Therefore
you will not be able to directly connect to it.

You will have to build from source because of a breaking change in the Application Binary Interface
(ABI) in geth.  You need version 1.8.1 or newer (1.8.1 was release Feb 16, 2018).

For example, init geth:

```bash
geth --datadir=/Users/pschlump/Projects/eth/data init /Users/pschlump/Projects/eth/genesis.json
```

See @Nik for the genesis.json file. You may want to put your data in a different location.

Then to run a local `geth` with websockets:

```bash
geth \
	--port 3000 --networkid 1101 \
	--identity "Philip.Schlump.1.8.1" --ws --wsaddr 192.168.0.157 --wsport 8546  \
	--rpc --rpcport 8545 --rpcaddr 192.168.0.157 --rpccorsdomain "*" --rpcapi "db,ssh,miner,admin,eth,net,web3,personal" \
	--datadir=/Users/pschlump/Projects/eth/data \
	--fast \
	--bootnodes=enode://1fab8525218222b26b7a72997fdccca2a79fa292172de1edda417061b1bd831dcfd257392c713f9b42d8853ad6a22d703affb8dbe118da48ff626943226bc64f@10.48.2.217:30301 \
	--etherbase="0xf08000b67882b699e475130b9e5faefc377f3c3f" \
	--mine --minerthreads=1
```

You probably need to use a different IP address, replace both 192.168.0.157's with localhost(127.0.0.1) or the IP address
of your system.

Note that I have turned on every JSON RPC interface.

The following assumes that you are going to call into the 'testnet' and that you are calling the existing 
contract.   The source for the contract is in ../contract/Threshold

First start the event watcher:

```
cd ../WatchForEvent
go build
./WatchForEvent
```

This will by default watch for RequestRelayEvent.

Now to call the contract:

```
cd ../CallContract
go build
./CallContract
```

The program has a bunch of options.   By default it will call KStart.RequestRelay.  This produces an event.




