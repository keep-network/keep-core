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

/**
 * Representing an abstract connector.
 */
export class AbstractConnector {
  /** @type {string} */
  name
  /** @type {Web3ProviderEngine} */
  provider

  /**
   * Create a connector.
   *
   * @param {string} name - Name of the connector.
   */
  constructor(name) {
    this.name = name
  }

  /**
   * Enable the connector and return available accounts.
   *
   * @return {Promise<Array<string>>} Available accounts.
   */
  enable = async () => {
    throw Error("Implement first")
  }

  /**
   * Disconnect a connector.
   */
  disconnect = async () => {
    throw Error("Implement first")
  }

  /**
   * @return {Promise<string>} The chain id.
   */
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

  /**
   * @return {Promise<string>} The network id.
   */
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

  /**
   * @return {Web3ProviderEngine} The web3 proivder engine.
   */
  getProvider = () => {
    return this.provider
  }
}

const DEFAULT_NUM_ADDRESSES_TO_FETCH = 15
/**
 * Representing an abstract connector for hardware wallets.
 * @extends AbstractConnector
 */
export class AbstractHardwareWalletConnector extends AbstractConnector {
  hardwareWalletProvider
  /** @type {string} To store selected account by user. */
  defaultAccount = ""

  /**
   * Create an abstract hardware wallet connector.
   *
   * @param {any} hardwareWalletProvider - The hardware wallet subprovider.
   * @param {string} name - Name of the connector.
   */
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

  /**
   * Get available accounts from the hardware wallet subprovider.
   *
   * @param {number} numberOfAccounts - The number of accounts to return.
   * @param {number} accountsOffSet - The start index.
   *
   * @return {Primise<Array<string>>} Accounts.
   */
  getAccounts = async (
    numberOfAccounts = DEFAULT_NUM_ADDRESSES_TO_FETCH,
    accountsOffSet = 0
  ) => {
    return await this.hardwareWalletProvider.getAccountsAsync(
      numberOfAccounts,
      accountsOffSet
    )
  }

  disconnect = () => {
    this.provider.stop()
  }
}
