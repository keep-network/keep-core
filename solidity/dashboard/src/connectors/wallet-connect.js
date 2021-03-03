import WalletConnect from "@walletconnect/web3-provider"
import CacheSubprovider from "web3-provider-engine/subproviders/cache"
import BigNumber from "bignumber.js"
import { AbstractConnector } from "./abstract-connector"
import { WALLETS } from "../constants/constants"
import { getRPCRequestPayload, overrideCacheMiddleware } from "./utils"

export class WalletConnectConnector extends AbstractConnector {
  constructor(
    options = {
      // TODO add chains
      rpc: {
        1337: "http://localhost:8545",
      },
    }
  ) {
    super(WALLETS.WALLET_CONNECT.name)
    this.provider = new WalletConnect(options)

    const cacheSubprovider = this.provider._providers.find(
      (provider) => provider.constructor.name === CacheSubprovider.name
    )
    overrideCacheMiddleware(cacheSubprovider)
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

  getChainId = async () => {
    try {
      const response = await this.provider.handleReadRequests(
        getRPCRequestPayload("eth_chainId")
      )
      // In case the provider returns chainId in hex.
      return new BigNumber(response.result).toString()
    } catch (error) {
      throw error
    }
  }

  getNetworkId = async () => {
    try {
      const response = await this.provider.handleReadRequests(
        getRPCRequestPayload("net_version")
      )
      return response.result
    } catch (error) {
      throw error
    }
  }
}
