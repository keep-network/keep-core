import { getChainId } from "../connectors/utils"

const chainID = getChainId()
export const useEtherscanUrl = () => {
  return chainID === 1 // Mainnet network ID.
    ? "https://etherscan.io"
    : "https://ropsten.etherscan.io"
}
