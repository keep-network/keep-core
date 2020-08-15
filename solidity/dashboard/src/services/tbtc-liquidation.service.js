import { contractService } from "./contracts.service"
import {
  TBTC_SYSTEM_CONTRACT_NAME,
  DEPOSIT_STATES,
} from "../constants/constants"
import moment from "moment"
import TBTCSystem from "@keep-network/tbtc/artifacts/TBTCSystem.json"
import {
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  createDepositContractInstance,
} from "../contracts"

const getPastRedemptionRequestedEvents = async (web3Context) => {
  const pastRedemptionRequestedEvents = await contractService.getPastEvents(
    web3Context,
    TBTC_SYSTEM_CONTRACT_NAME,
    "RedemptionRequested",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TBTC_SYSTEM_CONTRACT_NAME],
      // filter: { attr : paramVal },
    }
  )
  const augmentedRedemRequestEvents = await Promise.all(
    pastRedemptionRequestedEvents.map(async (event) => {
      let evDepStatus = null
      try {
        evDepStatus = await getDepositState(
          web3Context,
          event.returnValues._depositContractAddress
        )
      } catch (error) {
        console.log(error)
        throw error
      }
      event.depositStatusObj = evDepStatus
      event.isAlignedToDeposit = isEventAlignedToDeposit(event, evDepStatus)
      event.timestamp = await eventTimestamp(web3Context, event)

      return event
    })
  )
  return augmentedRedemRequestEvents
}

const depositNotifySignatureTimeout = async (
  web3Context,
  depositContractAddress
) => {
  const depositContractInstance = createDepositContractInstance(
    web3Context,
    depositContractAddress
  )
  depositContractInstance.methods
    .notifyRedemptionSignatureTimedOut()
    .send({
      from: web3Context.yourAddress,
    })
    .on("receipt", function (receipt) {
      console.log(
        `notifySignatureTimeout on Deposit ${depositContractAddress} was successful.`
      )
      return receipt
    })
    .on("error", function (error, receipt) {
      console.log(error.message)
      console.log(
        `notifySignatureTimeout on Deposit ${depositContractAddress} failed.`
      )
      throw error
    })
}

const depositNotifyRedemptionProofTimedOut = async (
  web3Context,
  depositContractAddress
) => {
  const depositContractInstance = createDepositContractInstance(
    web3Context,
    depositContractAddress
  )
  depositContractInstance.methods
    .notifyRedemptionProofTimedOut()
    .send({
      from: web3Context.yourAddress,
    })
    .on("receipt", function (receipt) {
      console.log(
        `notifyRedemptionProofTimedOut on Deposit ${depositContractAddress} was successful.`
      )
      return receipt
    })
    .on("error", function (error, receipt) {
      console.log(error.message)
      console.log(
        `notifyRedemptionProofTimedOut on Deposit ${depositContractAddress} failed.`
      )
      throw error
    })
}

const eventTimestamp = async (web3Context, event) => {
  const { eth } = web3Context
  const block = await eth.getBlock(event.blockNumber)
  return moment.unix(block.timestamp)
}

const getDepositState = async (web3Context, depositContractAddress) => {
  const depositContractInstance = createDepositContractInstance(
    web3Context,
    depositContractAddress
  )
  const depStateCode = await depositContractInstance.methods
    .currentState()
    .call()
  const depStateObj = DEPOSIT_STATES.find(
    (obj) => obj.depositStatusId == depStateCode
  )
  return depStateObj
}

const isEventAlignedToDeposit = (tbtcSystemContractEvent, evDepStatus) => {
  const contract = TBTCSystem
  let isAligned = false
  switch (evDepStatus.name) {
    case DEPOSIT_STATES.find((dso) => dso.depositStatusId == 5).name: // AWAITING_WITHDRAWAL_SIGNATURE
      if (tbtcSystemContractEvent.event == contract.abi[20].name)
        // RedemptionRequested
        isAligned = true
      break
    case DEPOSIT_STATES.find((dso) => dso.depositStatusId == 6).name: // AWAITING_WITHDRAWAL_PROOF
      if (tbtcSystemContractEvent.event == contract.abi[12].name)
        // GotRedemptionSignature
        isAligned = true
      break
    default:
      isAligned = true
  }
  return isAligned
}

export const liquidationService = {
  getPastRedemptionRequestedEvents,
  depositNotifySignatureTimeout,
  depositNotifyRedemptionProofTimedOut,
}
