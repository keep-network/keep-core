# Keep Maintainer Staging Mainnet

Bitcoin Relay used by the tTBTC system on testnet environment is an optimized version
of the Light Relay implemented by [SepoliaLightRelay](https://github.com/keep-network/tbtc-v2/blob/main/solidity/contracts/test/SepoliaLightRelay.sol)
contract. The reason the Light Relay had to be modified for testnet is that on
Bitcoin testnet difficulty often drops to `1`, which makes the blocks validation
on such change impossible to the regular Light Relay contract.

The `SepoliaLightRelay` version doesn't require a maintainer bot to submit block
headers on difficulty change. It accepts ad-hoc difficulty alignment according to
the tests needs.

The setup defined in this directory is meant for testing the bitcoin difficulty module 
of the maintainer bot with the Bitcoin mainnet blockchain state.
