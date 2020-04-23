# Introduction

Welcome to the Keep Network Smart contracts documentation.

- **[Keep Token](./docs/KeepToken/)** A standard ERC20Burnable token based on OpenZeppelin [ERC20Burnable.sol](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/ERC20Burnable.sol). Includes implementation for `approveAndCall` pattern.

- **[Token Staking](./docs/TokenStaking/)** A generic token staking contract for a specified standard ERC20Burnable token. A holder of the specified token can delegate its tokens to this contract and recover the stake after specified undelegation period is over.

- **[Token Grant](./docs/TokenGrant/)** A generic token grant contract for a specified standard ERC20Burnable token. Has additional functionality to stake delegate/undelegate token grants. Tokens are granted to the grantee via unlocking scheme and can be released gradually based on the unlocking schedule cliff and unlocking duration. Optionally grant can be revoked by the token grant manager.

# Buidler

To deploy contracts use [`buidler`](https://buidler.dev). Configuration can be
found in [buidler.config.js](../buidler.config.js).

## Deployment

Sample command for deployment on local network:

```
npx buidler run --network local scripts/buidler-deploy.js
```

`CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY` variable has to be set to a private key
that will be used for migrations.

Sample command for deployment on Ropsten:

```
CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY="0x...." \
    npx buidler run --network ropsten scripts/buidler-deploy.js
```

## Etherscan Verify Contracts

Update builder config and provide `apiKey` obtained from https://etherscan.io/myapikey.

Sample command to verify contracts with etherscan:

```sh
CONTRACT_OWNER_ETH_ACCOUNT_PRIVATE_KEY="0x...." \
    npx buidler verify-contract \
        --contract-name TokenGrant \
        --address 0xF9126Cd6554D98482788bE1808283C01B5a475D0 \
        "0x5c0995F988E7A4Ae8B8E300Fe592B5630B511566" # arguments
```
