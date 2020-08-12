import { Web3Loaded } from "../contracts"
import { isSameEthAddress } from "./general.utils"

export const getEventFromTransactionHash = async (
  web3Contract,
  eventName,
  txHash
) => {
  const web3 = await Web3Loaded

  const eventInterface = web3Contract.options.jsonInterface.find(
    (entry) => entry.type === "event" && entry.name === eventName
  )

  const eventSignature = eventInterface.signature

  const receipt = await web3.eth.getTransactionReceipt(txHash)

  const log = receipt.logs.find(
    (log) =>
      isSameEthAddress(web3Contract.options.address, log.address) &&
      log.topics[0] === eventSignature
  )

  if (!log) {
    return
  }

  return web3.eth.abi.decodeLog(eventInterface.inputs, log.data, log.topics)
}
