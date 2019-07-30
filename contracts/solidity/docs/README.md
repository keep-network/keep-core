# Introduction

Welcome to the Keep Network Smart contracts documentation.

- **[Keep Token](./docs/KeepToken/)** A standard ERC20 token based on OpenZeppelin [StandardToken.sol](https://github.com/OpenZeppelin/zeppelin-solidity/blob/master/contracts/token/ERC20/StandardToken.sol). Includes implementation for `approveAndCall` pattern.

- **[Token Staking](./docs/TokenStaking/)** A generic token staking contract for a specified standard ERC20 token. A holder of the specified token can stake its tokens to this contract and unstake after specified withdrawal delay is over.

- **[Token Grant](./docs/TokenGrant/)** A generic token grant contract for a specified standard ERC20 token. Has additional functionality to stake/unstake token grants. Tokens are granted to the grantee via vesting scheme and can be released gradually based on the vesting schedule cliff and vesting duration. Optionally grant can be revoked by the token grant manager.
