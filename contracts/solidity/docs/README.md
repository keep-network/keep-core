# Introduction

Welcome to the Keep Network Smart contracts documentation.

- **[Keep Token](./docs/KeepToken/)** A standard ERC20Burnable token based on OpenZeppelin [ERC20Burnable.sol](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/ERC20Burnable.sol). Includes implementation for `approveAndCall` pattern.

- **[Token Staking](./docs/TokenStaking/)** A generic token staking contract for a specified standard ERC20Burnable token. A holder of the specified token can stake its tokens to this contract and unstake after specified undelegation period is over.

- **[Token Grant](./docs/TokenGrant/)** A generic token grant contract for a specified standard ERC20Burnable token. Has additional functionality to stake/unstake token grants. Tokens are granted to the grantee via vesting scheme and can be released gradually based on the vesting schedule cliff and vesting duration. Optionally grant can be revoked by the token grant manager.
