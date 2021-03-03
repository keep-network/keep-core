/** @typedef { import("@keep-network/tbtc.js/src/EthereumHelpers.js").Contract } EthereumContract */

import Web3 from "web3"
import ProviderEngine from "web3-provider-engine"
import FilterSubprovider from "web3-provider-engine/subproviders/filters.js"
import NonceSubprovider from "web3-provider-engine/subproviders/nonce-tracker.js"
import WebsocketSubprovider from "web3-provider-engine/subproviders/websocket.js"
import { backoffRetrier } from "@keep-network/tbtc.js/src/lib/backoff.js"
import { logger } from "./winston.js"

/** @typedef {String} Address Ethereum address. */

/** @type {Number} Number of blocks for iterative events fetch. */
const GET_PAST_EVENTS_BLOCK_INTERVAL = 50000

/**
 * Initializes Web3 instance with a provider for given RPC URL and account private
 * key.
 * @param {String} url RPC URL
 * @param {String} ethPrivateKey Default account's private key.
 * @return {Web3} Web3 instance initialized to interact with the chain.
 */
export async function initWeb3(url, ethPrivateKey) {
  const engine = new ProviderEngine()
  const web3 = new Web3(engine)

  // filters
  engine.addProvider(new FilterSubprovider())

  // pending nonce
  engine.addProvider(new NonceSubprovider())

  // ws
  engine.addProvider(
    new WebsocketSubprovider({
      rpcUrl: url,
    })
  )

  engine.on("error", logger.warn) // report network connectivity errors

  engine.start()

  const account = web3.eth.accounts.privateKeyToAccount(ethPrivateKey)
  web3.eth.accounts.wallet.add(account)

  web3.eth.defaultAccount = account.address

  return web3
}

/**
 * Gets chain ID from the currently initialized Web3 instance.
 * @param {Web3} web3
 * @return {Number} Chain ID (1 - mainnet, 3 - ropsten, etc.)
 */
export async function getChainID(web3) {
  return await web3.eth.getChainId()
}

/**
 * Gets events from a contract. It tried to fetch events with one call for the whole
 * range. Some Ethereum API providers may reject request if the range is too wide
 * as a fallback the function will slice the range and gather events in chunks.
 *
 * @param {Web3} web3 Web3 instance.
 * @param {EthereumContract} contract Contract instance.
 * @param {String} eventName Event name.
 * @param {Number} fromBlock Starting block.
 * @param {Number} toBlock End block.
 * @return {Promise<[*]>} Found events.
 */
export async function getPastEvents(
  web3,
  contract,
  eventName,
  fromBlock = 0,
  toBlock
) {
  if (fromBlock < 0) {
    throw new Error(
      `fromBlock cannot be less than 0, current value: ${fromBlock}`
    )
  }

  return new Promise(async (resolve, reject) => {
    let resultEvents = []
    try {
      resultEvents = await backoffRetrier(3)(async () => {
        return await contract.getPastEvents(eventName, {
          fromBlock: fromBlock,
          toBlock: toBlock || "latest",
        })
      })
    } catch (err) {
      logger.warn(
        `switching to partial events pulls;` +
          `failed to get events in one request for event [${eventName}], ` +
          `fromBlock: [${fromBlock}], toBlock: [${toBlock}]: [${err.message}]`
      )

      try {
        const endBlock = toBlock || (await web3.eth.getBlockNumber())

        while (fromBlock <= endBlock) {
          let toBlock = fromBlock + GET_PAST_EVENTS_BLOCK_INTERVAL
          if (toBlock > endBlock) {
            toBlock = endBlock
          }
          logger.debug(
            `executing partial events pull for event [${eventName}], ` +
              `fromBlock: [${fromBlock}], toBlock: [${toBlock}]`
          )
          const foundEvents = await backoffRetrier(3)(async () => {
            return await contract.getPastEvents(eventName, {
              fromBlock: fromBlock,
              toBlock: toBlock,
            })
          })

          resultEvents = resultEvents.concat(foundEvents)
          logger.debug(
            `fetched [${foundEvents.length}] events, has ` +
              `[${resultEvents.length}] total`
          )

          fromBlock = toBlock + 1
        }
      } catch (error) {
        return reject(error)
      }
    }

    return resolve(resultEvents)
  })
}
