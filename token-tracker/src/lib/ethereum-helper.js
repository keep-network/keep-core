import { backoffRetrier } from "@keep-network/tbtc.js/src/lib/backoff.js"
import EthereumHelpers from "@keep-network/tbtc.js/src/EthereumHelpers.js"

/**
 * Ethereum address.
 * @typedef {String} Address
 */

const GET_PAST_EVENTS_BLOCK_INTERVAL = 50000

export async function callWithRetry(
  contractMethod,
  sendParams,
  totalAttempts,
  block
) {
  return EthereumHelpers.callWithRetry(
    contractMethod,
    sendParams,
    totalAttempts,
    block
  )
}

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
    } catch (error) {
      console.warn(
        `switching to partial events pulls; failed to get events in one request: [${error.message}]`
      )

      try {
        const endBlock = toBlock || (await web3.eth.getBlockNumber())

        while (fromBlock <= endBlock) {
          let toBlock = fromBlock + GET_PAST_EVENTS_BLOCK_INTERVAL
          if (toBlock > endBlock) {
            toBlock = endBlock
          }
          const foundEvents = await backoffRetrier(3)(async () => {
            return await contract.getPastEvents(eventName, {
              fromBlock: fromBlock,
              toBlock: toBlock,
            })
          })

          resultEvents = resultEvents.concat(foundEvents)

          fromBlock = toBlock + 1
        }
      } catch (error) {
        return reject(error)
      }
    }

    return resolve(resultEvents)
  })
}
