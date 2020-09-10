import { Web3Loaded } from "../contracts"
import { isSameEthAddress } from "./general.utils"

export const getEventFromLogs = async (logs, web3Contract, eventName) => {
  const web3 = await Web3Loaded

  const eventInterface = web3Contract.options.jsonInterface.find(
    (entry) => entry.type === "event" && entry.name === eventName
  )

  const eventSignature = eventInterface.signature

  const log = logs.find(
    (log) =>
      isSameEthAddress(web3Contract.options.address, log.address) &&
      log.topics[0] === eventSignature
  )

  if (!log) {
    return
  }

  return web3.eth.abi.decodeLog(
    eventInterface.inputs,
    log.data,
    log.topics.slice(1)
  )
}

export const getEventsFromTransaction = async (contractToEventName, txHash) => {
  const web3 = await Web3Loaded

  const receipt = await web3.eth.getTransactionReceipt(txHash)

  const events = {}

  for (const [contract, eventName] of contractToEventName) {
    const eventData = await getEventFromLogs(receipt.logs, contract, eventName)
    if (eventData) {
      events[eventName] = eventData
    }
  }

  return events
}
