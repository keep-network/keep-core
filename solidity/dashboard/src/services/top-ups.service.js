import moment from "moment"
import {
  ContractsLoaded,
  Web3Loaded,
  getContractDeploymentBlockNumber,
} from "../contracts"
import { isEmptyArray } from "../utils/array.utils"
import { add } from "../utils/arithmetics.utils"
import { TOKEN_STAKING_CONTRACT_NAME } from "../constants/constants"

export const commitTopUp = async (operator, onTransactionHashCallback) => {
  const { stakingContract } = await ContractsLoaded

  await stakingContract.methods
    .commitTopUp(operator)
    .send()
    .on("transactionHash", onTransactionHashCallback)
}

export const fetchAvailableTopUps = async (operators) => {
  const web3Context = await Web3Loaded
  const { stakingContract } = await ContractsLoaded

  const availableTopUps = []

  if (isEmptyArray(operators)) {
    return availableTopUps
  }

  const toupUpsInitiatedByOperator = (
    await stakingContract.getPastEvents("TopUpInitiated", {
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_STAKING_CONTRACT_NAME
      ),
      filter: { operator: operators },
    })
  ).reduce(reduceByOperator, {})

  const topUpsCompletedByOperator = (
    await stakingContract.getPastEvents("TopUpCompleted", {
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_STAKING_CONTRACT_NAME
      ),
      filter: { operator: operators },
    })
  ).reduce(reduceByOperator, {})

  for (const operator of operators) {
    const topUpsInitiated = toupUpsInitiatedByOperator[operator]
    const topUpsCompleted = topUpsCompletedByOperator[operator]

    const latestTopUpCompletedEvent = !isEmptyArray(topUpsCompleted)
      ? [...topUpsCompleted].pop()
      : undefined

    if (!isEmptyArray(topUpsInitiated)) {
      const availableOperatorTopUps = latestTopUpCompletedEvent
        ? topUpsInitiated.filter(
            filterByAfterLatestCompletedTopUp(latestTopUpCompletedEvent)
          )
        : topUpsInitiated
      const availableTopUpAmount = availableOperatorTopUps.reduce(
        reduceAmount,
        0
      )
      if (availableTopUpAmount > 0) {
        const createdAt = (
          await web3Context.eth.getBlock(
            topUpsInitiated[topUpsInitiated.length - 1].blockNumber
          )
        ).timestamp

        availableTopUps.push({
          operatorAddress: operator,
          availableTopUpAmount,
          createdAt,
        })
      }
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

export const isTopUpReadyToBeCommitted = (
  topUp,
  initializationPeriodInSeconds
) => {
  return moment
    .unix(topUp.createdAt)
    .add(initializationPeriodInSeconds, "seconds")
    .isBefore(moment())
}
