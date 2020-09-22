import Web3ProviderEngine from "web3-provider-engine"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket"
import CacheSubprovider from "web3-provider-engine/subproviders/cache"
import { getWsUrl } from "./utils"

import clone from "clone"

const DEFAULT_NUM_ADDRESSES_TO_FETCH = 15

export class AbstractHardwareWalletConnector extends Web3ProviderEngine {
  provider
  defaultAccount = ""

  constructor(provider) {
    super()
    this.provider = provider
    this.addProvider(this.provider)
    const cacheSubprovider = new CacheSubprovider()
    this.addProvider(cacheSubprovider) // initializes internal middleware

    // HACK ALERT Intercept middleware to always clone results. The cache
    // HACK ALERT subprovider caches results, but the cached values are mutable,
    // HACK ALERT and sure enough, the downstream handlers can and do at times
    // HACK ALERT mangle the results in non-idempotent ways. This means that
    // HACK ALERT when they receive cached values that they've already mangled
    // HACK ALERT later, everything blows up. This mini-middleware clones
    // HACK ALERT the results at the two exit points that the cache subprovider
    // HACK ALERT can use, ensuring that any downstream handlers are mutating
    // HACK ALERT a request-specific version of the value, without mangling the
    // HACK ALERT cached version.
    const originalMiddleware = cacheSubprovider.middleware.bind(cacheSubprovider)
    cacheSubprovider.middleware = (request, response, nextMiddleware, end) => {
      originalMiddleware(request, response, (handler) => {
        nextMiddleware((nextHandler) => {
          handler(nextHandler)
          // If the handler filled in a result, make sure to clone it so the
          // cache value is independent of downstream changes.
          response.result = clone(response.result)
        })
      }, (error) => {
        // If the handler filled in a result, make sure to clone it so the
        // cache value is independent of downstream changes.
        response.result = clone(response.result)
        end(error)
      })
    }

    this.addProvider(
      new WebsocketSubprovider({ rpcUrl: getWsUrl(), debug: false })
    )
  }

  enable = async () => {
    this.start()

    return await this.getAccount()
  }

  getAccounts = async (
    numberOfAccounts = DEFAULT_NUM_ADDRESSES_TO_FETCH,
    accountsOffSet = 0
  ) => {
    return await this.provider.getAccountsAsync(
      numberOfAccounts,
      accountsOffSet
    )
  }

  setProvider = (provider) => {
    this.provider = provider
  }

  getProvider = () => {
    return this.provider
  }

  getAccount = async () => {
    return [this.defaultAccount]
  }
}
