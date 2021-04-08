import {
  OPERATOR_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
} from "../constants/constants"
import {
  getContractDeploymentBlockNumber,
  Web3Loaded,
  ContractsLoaded,
} from "../contracts"
import { isEmptyArray } from "../utils/array.utils"
import { add, lte } from "../utils/arithmetics.utils"
import moment from "moment"

const fetchSlashedTokens = async (address) => {
  const { eth } = await Web3Loaded
  const {
    stakingContract,
    keepRandomBeaconOperatorContract,
  } = await ContractsLoaded
  const operatorEventsSearchFilters = {
    fromBlock: await getContractDeploymentBlockNumber(OPERATOR_CONTRACT_NAME),
  }

  const eventsSearchFilters = {
    fromBlock: await getContractDeploymentBlockNumber(
      TOKEN_STAKING_CONTRACT_NAME
    ),
    filter: { operator: address },
  }
  const data = []

  const slashedTokensEvents = await stakingContract.getPastEvents(
    "TokensSlashed",
    eventsSearchFilters
  )
  const seizedTokensEvents = await stakingContract.getPastEvents(
    "TokensSeized",
    eventsSearchFilters
  )

  if (isEmptyArray(slashedTokensEvents) && isEmptyArray(seizedTokensEvents)) {
    return data
  }

  const unauthorizedSigningEvents = await keepRandomBeaconOperatorContract.getPastEvents(
    "UnauthorizedSigningReported",
    operatorEventsSearchFilters
  )

  const relayEntryTimeoutEvents = await keepRandomBeaconOperatorContract.getPastEvents(
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
    let punishmentData
    if (slashedTokensGroupedByTxtHash.hasOwnProperty(transactionHash)) {
      const { amount } = slashedTokensGroupedByTxtHash[transactionHash]
      punishmentData = { amount, type: "SLASHED", event }
    } else if (seizedTokensGroupedByTxtHash.hasOwnProperty(transactionHash)) {
      const { amount } = seizedTokensGroupedByTxtHash[transactionHash]
      punishmentData = { amount, type: "SEIZED", event }
    } else {
      continue
    }

    if (punishmentData && lte(punishmentData.amount, 0)) continue

    punishmentData.date = moment.unix(
      (await eth.getBlock(blockNumber)).timestamp
    )
    punishmentData.groupPublicKey = await keepRandomBeaconOperatorContract.methods
      .getGroupPublicKey(groupIndex)
      .call()

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
      groupedByTransactionHash[transactionHash] = {
        ...returnValues,
        amount: returnValues.amount,
      }
    }
  })

  return groupedByTransactionHash
}
