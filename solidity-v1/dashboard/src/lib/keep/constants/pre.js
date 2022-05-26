import {
  ADDRESS_ZERO,
  PRE_CONTRACT_ADDRESS_MAINNET,
  PRE_CONTRACT_ADDRESS_ROPSTEN,
} from "./addresses"

export const MAINNET_PRE_DEPLOYMENT_BLOCK = 14141140
export const MAINNET_PRE_DEPLOYMENT_TX_HASH =
  "0x9a6db81f7db30bf80aa36c444729b3e4800c0dc77351f2fecbf8e2c0515defb0"
export const PRE_ADDRESSESS = {
  // https://etherscan.io/address/0x7E01c9c03FD3737294dbD7630a34845B0F70E5Dd
  1: PRE_CONTRACT_ADDRESS_MAINNET, // MAINNET
  // https://ropsten.etherscan.io/address/0xb6f98dA65174CE8F50eA0ea4350D96B2d3eFde9a
  3: PRE_CONTRACT_ADDRESS_ROPSTEN, // ROPSTEN
  // Set the correct `SimplePREApplication` contract address. If you deployed
  // the `@threshold-network/solidity-contracts` to your local chain and linked
  // package using `yarn link @threshold-network/solidity-contracts` you can
  // find the contract address at
  // `node_modules/@threshold-network/solidity-contracts/artifacts/SimplePREApplication.json`.
  1337: ADDRESS_ZERO, // LOCALHOST
}
