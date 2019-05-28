# keep-core

The core code behind the [Keep network].

### Getting Started as a Developer

For now, start with the README in [the developer docs directory](docs/development/).

### [`docs/`](docs/)

Documentation related to the Keep network, Keep client, and Keep contracts.

#### [`docs/development/`](docs/development/)

Specifically developer documentation for the various parts of Keep.

### [`solidity/`](contracts/solidity/)

The smart contracts behind the [Keep network].

They handle creating and managing keeps, bridging off-chain secret storage
and the public blockchain.

### [`go/`](cmd/)

The Keep Go client.

It runs the Keep network's random beacon, hosts keep nodes, and participates in
keep computations.

  [Keep network]: https://keep.network
