# T Merkle Distributor

Solidity contract for Threshold rewards' distribution.

In Cumulative Merkle Drop contract each new token distribution replaces previous one and should contain cumulative balances of all the participants. Cumulative claimed amount is used as invalidation for every participant.

## Structure

This is a Hardhat project:

* `contracts`: Source code for contract
* `test`: Contract tests
* `scripts`: Hardhat scripts:
  * `gen_merkle_dist.js`: generate new Merkle distribution
  * `verify_proof.js`: verify Merkle proof
  * `example_run_claim.js`: example about how to claim tokens

## Installation

```bash
npm install
```

## Run scripts

```bash
npx hardhat run scripts/<script.js>
```

## Run tests

```bash
npx hardhat test
```

## Deploy
To deploy to the Ropsten test network you will need a `.env` that looks similar to:

```
ROPSTEN_RPC_URL="https://ropsten.infura.io/v3/bd76xxxxxxxxxxxxxxxxxxxxxxxxxff0"
ROPSTEN_PRIVATE_KEY="3d3ad2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx87b"
MAINNET_ETHERSCAN_KEY="M5xxxxxxxxxxxxxxxxxxxxxxxxxxxxxSMV"
```

You can then run
```bash
npx hardhat --network ropsten deploy
```

The contract will be deployed and the source code will be verified on etherscan.