import WalletConnect from "@walletconnect/web3-provider"
import CacheSubprovider from "web3-provider-engine/subproviders/cache"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket"
import BigNumber from "bignumber.js"
import { AbstractConnector } from "./abstract-connector"
import { WALLETS } from "../constants/constants"
import {
  getRPCRequestPayload,
  overrideCacheMiddleware,
  getWsUrl,
} from "./utils"

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
      // Override `handleRequest` in order to support event subscriptions via
      // `WebsocketSubprovider`
      const requestProvider = this.provider._providers[
        this.provider._providers.length - 1
      ]
      const originalHandleRequest = requestProvider.handleRequest.bind(
        requestProvider
      )

      requestProvider.handleRequest = async (payload, next, end) => {
        switch (payload.method) {
          case "eth_getLogs":
          case "eth_subscribe": {
            // Pass this request to the next subprovider.
            next()
            return
          }
          default: {
            await originalHandleRequest(payload, next, end)
            return
          }
        }
      }

      this.provider.addProvider(
        new WebsocketSubprovider({
          rpcUrl: getWsUrl(),
          debug: true,
        })
      )
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
