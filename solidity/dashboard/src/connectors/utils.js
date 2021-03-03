import Common from "ethereumjs-common"
import { Transaction as EthereumTx } from "ethereumjs-tx"
import clone from "clone"
import config from "../config/config.json"
import { getFirstNetworkIdFromArtifact } from "../contracts"

export const getEthereumTxObj = (txData, chainId) => {
  const customCommon = Common.forCustomChain("mainnet", {
    name: "keep-dev",
    chainId,
  })
  const common = new Common(customCommon._chainParams, "petersburg", [
    "petersburg",
  ])
  return new EthereumTx(txData, { common })
}

// EIP-155 https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
// v = CHAIN_ID * 2 + 35 => CHAIN_ID = (v - 35) / 2
export const getChainIdFromV = (vInHex) => {
  const vIntValue = parseInt(vInHex, 16)
  const chainId = Math.floor((vIntValue - 35) / 2)
  return chainId < 0 ? 0 : chainId
}

export const getChainId = () => {
  if (process.env.NODE_ENV === "development") {
    // private chains (default), change if you use a different one
    return 1337
  }
  // For KEEP internal testnet, ropsten and mainnet `chainId == networkId`
  return Number(getFirstNetworkIdFromArtifact())
}

export const getWsUrl = () => {
  if (process.env.NODE_ENV === "development") {
    // Ganache web socket url, change if you use a different one
    return "ws://localhost:8545"
  }
  return config.networks[getChainId()].wsURL
}

export const getRPCRequestPayload = (method, params = []) => {
  return {
    jsonrpc: "2.0",
    method,
    params,
    id: new Date().getTime(),
  }
}

export const overrideCacheMiddleware = (cacheSubprovider) => {
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
    originalMiddleware(
      request,
      response,
      (handler) => {
        nextMiddleware((nextHandler) => {
          handler(nextHandler)
          // If the handler filled in a result, make sure to clone it so the
          // cache value is independent of downstream changes.
          response.result = clone(response.result)
        })
      },
      (error) => {
        // If the handler filled in a result, make sure to clone it so the
        // cache value is independent of downstream changes.
        response.result = clone(response.result)
        end(error)
      }
    )
  }
}
