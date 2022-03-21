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
