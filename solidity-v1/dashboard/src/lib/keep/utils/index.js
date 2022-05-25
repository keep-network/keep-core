import { getFirstNetworkIdFromArtifact } from "../contracts"

export const getChainId = () => {
  if (
    process.env.NODE_ENV === "development" ||
    process.env.NODE_ENV === "test"
  ) {
    // private chains (default), change if you use a different one
    return 1337
  }
  // For KEEP internal testnet, ropsten and mainnet `chainId == networkId`
  return Number(getFirstNetworkIdFromArtifact())
}