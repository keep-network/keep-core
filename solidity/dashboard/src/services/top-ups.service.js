import { contracts, CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { isEmptyArray } from "../utils/array.utils"
import { add } from "../utils/arithmetics.utils"

export const commitTopUp = async (operator, onTransactionHashCallback) => {
  await contracts.stakingContract.methods
    .commitTopUp(operator)
    .send()
    .on("transactionHash", onTransactionHashCallback)
}

export const fetchAvailableTopUps = async (operators) => {
  const availableTopUps = []

  const toupUpsInitiatedByOperator = (
    await contracts.stakingContract.getPastEvents("TopUpInitiated", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
      filter: { operator: [operators] },
    })
  ).reduce(reduceByOperator, {})
  const toupUpsCompletedByOperator = (
    await contracts.stakingContract.getPastEvents("TopUpCompleted", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
      filter: { operator: [operators] },
    })
  ).reduce(reduceByOperator, {})

  for (const operator of operators) {
    const topUpInitiated = toupUpsInitiatedByOperator[operator]
    const topUpCompleted = toupUpsCompletedByOperator[operator]

    const latestTopUpCompletedEvent = !isEmptyArray(topUpCompleted)
      ? topUpCompleted.pop()
      : undefined

    if (latestTopUpCompletedEvent && !isEmptyArray(topUpInitiated)) {
      const availableTopUpAmount = topUpInitiated
        .filter(filterByAfterLatestCompletedTopUp(latestTopUpCompletedEvent))
        .reduce(reduceAmount, 0)

      availableTopUps.push({ operatorAddress: operator, availableTopUpAmount })
    }
  }

  return availableTopUps
}

const reduceByOperator = (result, event) => {
  const {
    returnValues: { operator },
  } = event

  result[operator] = (result[operator] || []).push(event)

  return result
}

const reduceAmount = (result, { returnValues: { topUp } }) => {
  return add(result, topUp)
}

const filterByAfterLatestCompletedTopUp = (latestTopUpCompletedEvent) => (
  initiatedEvent
) => {
  const isAfterLatestCompleteedTopUpBlock =
    initiatedEvent.blockNumber > latestTopUpCompletedEvent.blockNumber
  const isAfterLatestCompletedTopUpTransactionInBlock =
    latestTopUpCompletedEvent.blockNumber === initiatedEvent.blockNumber &&
    initiatedEvent.transactionIndex > latestTopUpCompletedEvent.transactionIndex

  return (
    isAfterLatestCompleteedTopUpBlock ||
    isAfterLatestCompletedTopUpTransactionInBlock
  )
}
