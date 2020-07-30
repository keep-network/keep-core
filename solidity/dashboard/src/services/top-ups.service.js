import { contracts, CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"
import { isEmptyArray } from "../utils/array.utils"
import { add } from "../utils/arithmetics.utils"

export const commitTopUp = async (operator, onTransactionHashCallback) => {
  await contracts.stakingContract.methods
    .commitTopUp(operator)
    .send()
    .on("transactionHash", onTransactionHashCallback)
}

export const fetchAvailableTopUps = async (_, operators) => {
  const availableTopUps = []

  if (isEmptyArray(operators)) {
    return availableTopUps
  }

  const toupUpsInitiatedByOperator = (
    await contracts.stakingContract.getPastEvents("TopUpInitiated", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
      filter: { operator: operators },
    })
  ).reduce(reduceByOperator, {})

  const topUpsCompletedByOperator = (
    await contracts.stakingContract.getPastEvents("TopUpCompleted", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
      filter: { operator: operators },
    })
  ).reduce(reduceByOperator, {})

  for (const operator of operators) {
    const topUpInitiated = toupUpsInitiatedByOperator[operator]
    const topUpCompleted = topUpsCompletedByOperator[operator]

    const latestTopUpCompletedEvent = !isEmptyArray(topUpCompleted)
      ? [...topUpCompleted].pop()
      : undefined

    if (!isEmptyArray(topUpInitiated)) {
      const availableOperatorTopUps = latestTopUpCompletedEvent
        ? topUpInitiated.filter(
            filterByAfterLatestCompletedTopUp(latestTopUpCompletedEvent)
          )
        : topUpInitiated
      const availableTopUpAmount = availableOperatorTopUps.reduce(
        reduceAmount,
        0
      )
      if (availableTopUpAmount > 0)
        availableTopUps.push({
          operatorAddress: operator,
          availableTopUpAmount,
        })
    }
  }

  return availableTopUps
}

const reduceByOperator = (result, event) => {
  const {
    returnValues: { operator },
  } = event

  ;(result[operator] = result[operator] || []).push(event)

  return result
}

const reduceAmount = (result, { returnValues: { topUp } }) => {
  return add(result, topUp)
}

const filterByAfterLatestCompletedTopUp = (latestTopUpCompletedEvent) => (
  initiatedEvent
) => {
  const isAfterLatestCompletedTopUpBlock =
    initiatedEvent.blockNumber > latestTopUpCompletedEvent.blockNumber
  const isAfterLatestCompletedTopUpTransactionInBlock =
    latestTopUpCompletedEvent.blockNumber === initiatedEvent.blockNumber &&
    initiatedEvent.transactionIndex > latestTopUpCompletedEvent.transactionIndex

  return (
    isAfterLatestCompletedTopUpBlock ||
    isAfterLatestCompletedTopUpTransactionInBlock
  )
}
