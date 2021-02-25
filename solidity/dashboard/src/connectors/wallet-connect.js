import WalletConnect from "@walletconnect/web3-provider"
import { AbstractConnector } from "./abstract-connector"

export class WalletConnectConnector extends AbstractConnector {
  constructor(
    options = {
      // TODO add chains
      rpc: {
        1337: "http://localhost:8545",
      },
    }
  ) {
    super()
    this.provider = new WalletConnect(options)
  }

  enable = async () => {
    try {
      const accounts = await this.provider.enable()
      return accounts
    } catch (error) {
      if (error.message === "User closed modal") {
        throw new Error("The user rejected the request")
      }

      throw error
    }
  }

  disconnect = async () => {
    await this.provider.disconnect()
  }

  sendAsync = async (payload, callback) => {
    return await this.provider.sendAsync(payload, callback)
  }

  getChainId = async () => {
    return this.provider.chainId
  }

  getNetworkId = async () => {
    try {
      const response = await this.provider.handleReadRequests({
        jsonrpc: "2.0",
        method: "net_version",
        params: [],
        id: new Date().getTime(),
      })
      return response.result.toString(16)
    } catch (error) {
      throw error
    }
  }
}
