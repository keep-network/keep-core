import Web3ProviderEngine from "web3-provider-engine"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket"
import CacheSubprovider from "web3-provider-engine/subproviders/cache"
import BigNumber from "bignumber.js"
import {
  getWsUrl,
  getRPCRequestPayload,
  overrideCacheMiddleware,
} from "./utils"

import { Deferred } from "../contracts"

const DEFAULT_NUM_ADDRESSES_TO_FETCH = 15

export class AbstractConnector {
  name

  constructor(name) {
    this.name = name
  }

  enable = async () => {
    throw Error("Implement first")
  }

  getAccounts = async () => {
    throw Error("Implement first")
  }

  disconnect = async () => {
    throw Error("Implement first")
  }

  getChainId = async () => {
    const chainIdDeferred = new Deferred()
    this.provider.sendAsync(
      getRPCRequestPayload("eth_chainId"),
      (error, response) => {
        if (error) {
          chainIdDeferred.reject(error)
        } else {
          // Convert response to `BigNumber` in case the provider returns chainId in hex.
          chainIdDeferred.resolve(new BigNumber(response.result).toString())
        }
      }
    )
    return await chainIdDeferred.promise
  }

  getNetworkId = async () => {
    const netIdDeferred = new Deferred()
    this.provider.sendAsync(
      getRPCRequestPayload("net_version"),
      (error, response) => {
        if (error) {
          netIdDeferred.reject(error)
        } else {
          netIdDeferred.resolve(response.result)
        }
      }
    )
    return await netIdDeferred.promise
  }

  getProvider = () => {
    return this.provider
  }
}

export class AbstractHardwareWalletConnector extends Web3ProviderEngine {
  hardwareWalletProvider
  defaultAccount = ""

  constructor(hardwareWalletProvider, name) {
    super(name)
    this.hardwareWalletProvider = hardwareWalletProvider
  }

  enable = async () => {
    const web3Engine = new Web3ProviderEngine()

    web3Engine.addProvider(this.hardwareWalletProvider)
    const cacheSubprovider = new CacheSubprovider()
    web3Engine.addProvider(cacheSubprovider) // initializes internal middleware
    overrideCacheMiddleware(cacheSubprovider)

    web3Engine.addProvider(
      new WebsocketSubprovider({ rpcUrl: getWsUrl(), debug: true })
    )
    this.provider = web3Engine
    this.provider.start()

    return this.defaultAccount
      ? [this.defaultAccount]
      : await this.getAccounts()
  }

  getAccounts = async (
    numberOfAccounts = DEFAULT_NUM_ADDRESSES_TO_FETCH,
    accountsOffSet = 0
  ) => {
    return await this.hardwareWalletSubprovider.getAccountsAsync(
      numberOfAccounts,
      accountsOffSet
    )
  }
}
