import { useContext } from "react"
import { Web3Context } from "../components/WithWeb3Context"

export const useEtherscanUrl = () => {
  const { networkType } = useContext(Web3Context)

  return networkType === "main"
    ? "https://etherscan.io"
    : "https://ropsten.etherscan.io"
}
