import { contractService } from "./contracts.service"
import {
  OPERATOR_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
} from "../constants/constants"
import { CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { isEmptyArray } from "../utils/array.utils"
import { add } from "../utils/arithmetics.utils"
import moment from "moment"

const fetchSlashedTokens = async (web3Context) => {
  const { yourAddress, eth } = web3Context
  const operatorEventsSearchFilters = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[OPERATOR_CONTRACT_NAME],
  }

  const eventsSearchFilters = {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME],
    filter: { operator: yourAddress },
  }
  const data = []

  const slashedTokensEvents = await contractService.getPastEvents(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "TokensSlashed",
    eventsSearchFilters
  )
  const seizedTokensEvents = await contractService.getPastEvents(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "TokensSeized",
    eventsSearchFilters
  )

  if (isEmptyArray(slashedTokensEvents) && isEmptyArray(seizedTokensEvents)) {
    return data
  }

  const unauthorizedSigningEvents = await contractService.getPastEvents(
    web3Context,
    OPERATOR_CONTRACT_NAME,
    "UnauthorizedSigningReported",
    operatorEventsSearchFilters
  )

  const relayEntryTimeoutEvents = await contractService.getPastEvents(
    web3Context,
    OPERATOR_CONTRACT_NAME,
    "RelayEntryTimeoutReported",
    operatorEventsSearchFilters
  )

  const punishmentEvents = [
    ...unauthorizedSigningEvents,
    ...relayEntryTimeoutEvents,
  ]
  const slashedTokensGroupedByTxtHash = groupByTransactionHash(
    slashedTokensEvents
  )
  const seizedTokensGroupedByTxtHash = groupByTransactionHash(
    seizedTokensEvents
  )

  for (let i = 0; i < punishmentEvents.length; i++) {
    const {
      event,
      transactionHash,
      blockNumber,
      returnValues: { groupIndex },
    } = punishmentEvents[i]
    let punishmentData = {}
    if (slashedTokensGroupedByTxtHash.hasOwnProperty(transactionHash)) {
      const { amount } = slashedTokensGroupedByTxtHash[transactionHash]
      punishmentData = { amount, type: "SLASHED", event }
    } else if (seizedTokensGroupedByTxtHash.hasOwnProperty(transactionHash)) {
      const { amount } = seizedTokensGroupedByTxtHash[transactionHash]
      punishmentData = { amount, type: "SEIZED", event }
    }

    punishmentData.date = moment.unix(
      (await eth.getBlock(blockNumber)).timestamp
    )
    punishmentData.groupPublicKey = await contractService.makeCall(
      web3Context,
      OPERATOR_CONTRACT_NAME,
      "getGroupPublicKey",
      groupIndex
    )

    data.push(punishmentData)
  }
  return data
}

export const slashedTokensService = {
  fetchSlashedTokens,
}

const groupByTransactionHash = (events) => {
  const groupedByTransactionHash = {}

  events.forEach((event) => {
    const { transactionHash, returnValues } = event
    if (groupedByTransactionHash.hasOwnProperty(transactionHash)) {
      const prevData = groupedByTransactionHash[transactionHash]
      groupedByTransactionHash[transactionHash] = {
        ...returnValues,
        amount: add(returnValues.amount, prevData.amount),
      }
    } else {
      groupedByTransactionHash[transactionHash] = { ...returnValues }
    }
  })

  return groupedByTransactionHash
}
