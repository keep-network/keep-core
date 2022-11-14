# Threshold Network rewards Merkle distribution

Solidity contract and scripts for Threshold Network rewards' distribution.

In Cumulative Merkle Drop contract, each new token distribution replaces previous one and should
contain the cumulative balances of all the participants. Cumulative claimed amount is used as
invalidation for every participant.

## Structure

- `contracts`: Source code for contract
- `test`: Hardhat contract tests.
- `src/scripts`:
  - `gen_rewards_dist.js`: generate new Merkle distribution for the Threshold Network rewards earned
    in a specific period.
  - `verify_proof.js`: verify Merkle proof of a distribution.
  - `stake_history.js`: fetch the information of a particular staker, including staking history.
  - `claimed_rewards.js`: calculate the Threshold rewards that has been already claimed.

## Installation

```bash
npm install
```

In order to run the scripts, it's needed to have a `.env` file that includes:

```
ETHERSCAN_TOKEN=<your Etherscan API token>
```

## Run scripts

> **NOTE:** Scripts must be run from the repo root, and not from the folder that contains them.

### gen_rewards_dist script

This script calculates the Threshold Network rewards earned during a specific period, adds them to
the previous distributions, and generates a new distribution that contains the cumulative rewards.

Note that some parameters (rewards weights, start time, end time, last distribution path...) must be
changed in the script before running it.

```bash
node src/scripts/gen_rewards_dist.js
```

### stake_history script

This script fetch the information of a particular staker, including staking history.

```bash
node src/scripts/stake_history <0x-prefixed staking provider address>
```

## Run Hardhat tests

```bash
npx hardhat test
```

## Deploy

To deploy to the Goerli test network you will need a `.env` that looks similar to:

```
GOERLI_RPC_URL="https://goerli.infura.io/v3/bd76xxxxxxxxxxxxxxxxxxxxxxxxxff0"
GOERLI_PRIVATE_KEY="3d3ad2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx87b"
MAINNET_ETHERSCAN_KEY="M5xxxxxxxxxxxxxxxxxxxxxxxxxxxxxSMV"
```

You can then run

```bash
npx hardhat --network goerli deploy
```

The contract will be deployed and the source code will be verified on etherscan.

## Test Deployment

In order to run a test deployment:

```bash
npx hardhat --network mainnet_test deploy
```

This will use the deployment script in `deploy-test`.
The difference is that it also deploys a mock Token contract, which makes testing on mainnet possible.
