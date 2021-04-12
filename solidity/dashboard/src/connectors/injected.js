import { AbstractConnector } from "./abstract-connector"
import { WALLETS } from "../constants/constants"
import BigNumber from "bignumber.js"
import { UserRejectedConnectionRequestError } from "./utils"

class InjectedConnector extends AbstractConnector {
  constructor() {
    super(WALLETS.METAMASK.name)
    this.provider = window.ethereum
  }

  _onAccountChanged = ([address]) => {
    this.emit("accountsChanged", address)
  }

  _onDisconnect = () => {
    this.emit("disconnect")
  }

  _onChainChanged = (chainId) => {
    this.emit("chainChanged", chainId)
  }

  enable = async () => {
    if (!this.provider) {
      throw new Error("window.ethereum provider not found")
    }

    try {
      const accounts = await this.provider.request({
        method: "eth_requestAccounts",
      })
      // https://docs.metamask.io/guide/ethereum-provider.html#ethereum-autorefreshonnetworkchange
      this.provider.autoRefreshOnNetworkChange = false
      if (this.provider && this.provider.on) {
        this.provider.on("accountsChanged", this._onAccountChanged)
        this.provider.on("disconnect", this._onDisconnect)
        this.provider.on("chainChanged", this._onChainChanged)
      }
      return accounts
    } catch (error) {
      if (error.code === 4001) {
        // EIP-1193 userRejectedRequest error
        // If this happens, the user rejected the connection request.
        console.error("The user rejected the connection request.")
        throw new UserRejectedConnectionRequestError()
      }
      throw error
    }
  }

  disconnect = async () => {
    // window.ethereum injected by MetaMask does not provide a method to
    // disconnect a wallet.
    this._onDisconnect()
    this.provider.removeListener("accountsChanged", this._onAccountChanged)
    this.provider.removeListener("disconnect", this._onDisconnect)
    this.provider.removeListener("chainChanged", this._onChainChanged)
  }

  getChainId = async () => {
    const chainId = await this.provider.request({ method: "eth_chainId" })
    // In case the provider returns chainId in hex.
    return new BigNumber(chainId).toString()
  }

  getNetworkId = async () => {
    return await this.provider.request({ method: "net_version" })
  }

  getAccounts = async () => {
    return this.provider && this.provider.request
      ? await this.provider.request({ method: "eth_accounts" })
      : []
  }
}

export { InjectedConnector }
