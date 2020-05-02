import Common from "ethereumjs-common"
import { Transaction as EthereumTx } from "ethereumjs-tx"
import config from "../config/config.json"
import { getFirstNetworkIdFromArtifact } from "../contracts"

export const getEthereumTxObj = (txData, chainId) => {
  const customCommon = Common.forCustomChain("mainnet", {
    name: "keep-dev",
    chainId,
  })
  const common = new Common(customCommon._chainParams, "petersburg", [
    "petersburg",
  ])
  return new EthereumTx(txData, { common })
}

// EIP-155 https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
// v = CHAIN_ID * 2 + 35 => CHAIN_ID = (v - 35) / 2
export const getChainIdFromV = (vInHex) => {
  const vIntValue = parseInt(vInHex, 16)
  const chainId = Math.floor((vIntValue - 35) / 2)
  return chainId < 0 ? 0 : chainId
}

export const getChainId = () => {
  if (process.env.NODE_ENV === "development") {
    // private chains (default), change if you use a different one
    return 1337
  }
  // For KEEP internal testnet, ropsten and mainnet `chainId == networkId`
  return getFirstNetworkIdFromArtifact()
}

export const getWsUrl = () => {
  if (process.env.NODE_ENV === "development") {
    // Ganache web socket url, change if you use a different one
    return "ws://localhost:8545"
  }
  return config.networks[getChainId()].wsURL
}
