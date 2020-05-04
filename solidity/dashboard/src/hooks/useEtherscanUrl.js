import { useContext } from "react"
import { Web3Context } from "../components/WithWeb3Context"

export const useEtherscanUrl = () => {
  const { networkType } = useContext(Web3Context)

  switch (networkType) {
    case "main":
      return "https://etherscan.io"
    case "ropsten":
      return "https://ropsten.etherscan.io"
    default:
      return "https://ropsten.etherscan.io"
  }
}
