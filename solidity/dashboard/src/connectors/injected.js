import { AbstractConnector } from "./abstract-connector"
import { WALLETS } from "../constants/constants"
import BigNumber from "bignumber.js"

class InjectedConnector extends AbstractConnector {
  constructor() {
    super(WALLETS.METAMASK.name)
    this.provider = window.ethereum
  }

  enable = async () => {
    if (!this.provider) {
      throw new Error("window.ethereum provider not found")
    }

    // https://docs.metamask.io/guide/ethereum-provider.html#ethereum-autorefreshonnetworkchange
    this.provider.autoRefreshOnNetworkChange = false

    try {
      return await this.provider.request({ method: "eth_requestAccounts" })
    } catch (error) {
      if (error.code === 4001) {
        // EIP-1193 userRejectedRequest error
        // If this happens, the user rejected the connection request.
        console.error("The user rejected the connection request.")
        throw new Error("User rejected request")
      }
      throw error
    }
  }

  disconnect = async () => {
    // window.ethereum injected by MetaMask does not provide a method to
    // disconnect a wallet.
  }

  getChainId = async () => {
    const chainId = await this.provider.request({ method: "eth_chainId" })
    // In case the provider returns chainId in hex.
    return new BigNumber(chainId).toString()
  }

  getNetworkId = async () => {
    return await this.provider.request({ method: "net_version" })
  }
}

export { InjectedConnector }
