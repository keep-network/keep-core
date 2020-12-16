import { contractService } from "./contracts.service"
import {
  TOKEN_STAKING_CONTRACT_NAME
} from "../constants/constants"
import moment from "moment"
import {
  isCodeValid,
  createManagedGrantContractInstance,
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  ContractsLoaded,
  Web3Loaded,
} from "../contracts"
import { isSameEthAddress } from "../utils/general.utils"
import { isEmptyArray } from "../utils/array.utils"
import { getEventsFromTransaction } from "../utils/ethereum.utils"

const delegationInfoFromStakedEvents = async (address) => {
  const { stakingContract } = await ContractsLoaded
  let operatorStakedEvents
  try {
    operatorStakedEvents = await stakingContract.getPastEvents(
      "OperatorStaked",
      {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
        filter: { operator: address },
      }
    )
  } catch (err) {}

  if (operatorStakedEvents.length === 0)
    throw new Error("No OperatorStaked events found for address " + address)

  const {
    transactionHash: stakingTransactionHash,
    returnValues: {
      beneficiary: beneficiaryAddress,
      authorizer: authorizerAddress,
    },
  } = operatorStakedEvents[0]

  return {
    stakingTransactionHash,
    beneficiaryAddress,
    authorizerAddress,
  }
}

const fetchDelegatedTokensData = async () => {
  const { grantContract, stakingContract } = await ContractsLoaded
  const web3 = await Web3Loaded
  const {
    eth,
    eth: { defaultAccount: yourAddress },
  } = web3
  let ownerAddress

  const {
    stakingTransactionHash,
    beneficiaryAddress,
    authorizerAddress,
  } = await delegationInfoFromStakedEvents(yourAddress)

  const [stakedBalance, initializationPeriod] = await Promise.all([
    stakingContract.methods.balanceOf(yourAddress).call(),
    stakingContract.methods.initializationPeriod().call(),
  ])

  const eventsToCheck = [[grantContract, "TokenGrantStaked"]]
  let events
  try {
    events = await getEventsFromTransaction(
      eventsToCheck,
      stakingTransactionHash
    )
  } catch (err) {}

  let isDelegationFromGrant = true
  let isManagedGrant = false
  let managedGrantContractInstance
  let grantId
  if (events && events.TokenGrantStaked) {
    isDelegationFromGrant = true
    grantId = events.TokenGrantStaked.grantId
    const { grantee } = await grantContract.methods.getGrant(grantId).call()
    ownerAddress = grantee
    // check if grantee is a contract
    const code = await eth.getCode(grantee)
    if (isCodeValid(code)) {
      managedGrantContractInstance = createManagedGrantContractInstance(
        web3,
        grantee
      )
      isManagedGrant = true
      ownerAddress = await managedGrantContractInstance.methods.grantee().call()
    }
  } else {
    ownerAddress = await stakingContract.methods.ownerOf(yourAddress).call()
  }

  const {
    undelegationStatus,
    undelegation,
    undelegationPeriod,
    delegationStatus,
    undelegationCompletedAt,
  } = await fetchPendingUndelegation()
  const { createdAt } = undelegation
  const initializationOverAt = moment
    .unix(createdAt)
    .add(initializationPeriod, "seconds")
  const isInInitializationPeriod = moment().isSameOrBefore(initializationOverAt)

  return {
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
    undelegationStatus,
    isDelegationFromGrant,
    isInInitializationPeriod,
    undelegationPeriod,
    isManagedGrant,
    managedGrantContractInstance,
    delegationStatus,
    undelegationCompletedAt,
  }
}

const fetchPendingUndelegation = async () => {
  const contractsLoaded = await ContractsLoaded
  const { stakingContract } = contractsLoaded
  const web3 = await Web3Loaded
  const { defaultAccount: yourAddress } = web3.eth

  const [delegation, undelegationPeriod] = await Promise.all([
    stakingContract.methods.getDelegationInfo(yourAddress).call(),
    stakingContract.methods.undelegationPeriod().call(),
  ])

  const { undelegatedAt, createdAt, amount } = delegation

  const isUndelegation = delegation.undelegatedAt !== "0"
  const pendingUnstakeBalance = isUndelegation ? delegation.amount : 0
  const undelegationCompletedAt = isUndelegation
    ? moment.unix(undelegatedAt).add(undelegationPeriod, "seconds")
    : null

  let delegationStatus
  if (amount !== "0" && createdAt !== "0" && undelegatedAt !== "0") {
    // delegation undelegated
    delegationStatus = "UNDELEGATED"
  } else if (amount === "0" && createdAt !== "0" && undelegatedAt === "0") {
    // delegation canceled
    delegationStatus = "CANCELED"
  } else if (amount === "0" && createdAt !== "0" && undelegatedAt !== "0") {
    // delegation recovered
    delegationStatus = "RECOVERED"
  }

  return {
    pendingUnstakeBalance,
    undelegationCompletedAt,
    undelegationPeriod,
    delegationStatus,
    undelegation: delegation,
  }
}

export const operatorService = {
  fetchDelegatedTokensData,
  fetchPendingUndelegation,
}

export const getOperatorsOfAuthorizer = async (web3Context, authorizer) => {
  return (
    await contractService.getPastEvents(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "OperatorStaked",
      {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME],
        filter: { authorizer },
      }
    )
  ).map((_) => _.returnValues.operator)
}

export const getOperatorsOfBeneficiary = async (beneficiary) => {
  const { stakingContract } = await ContractsLoaded

  return (
    await stakingContract.getPastEvents("OperatorStaked", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME],
      filter: { beneficiary },
    })
  ).map((_) => _.returnValues.operator)
}

export const getOperatorsOfOwner = async (owner, operatorsFilterParam) => {
  const { stakingContract } = await ContractsLoaded
  const filterParam = operatorsFilterParam
    ? { operator: operatorsFilterParam }
    : {}

  const ownerDelegations = await stakingContract.getPastEvents(
    "StakeDelegated",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
      filter: { owner, ...filterParam },
    }
  )

  const transferEventsByOwner = await stakingContract.getPastEvents(
    "StakeOwnershipTransferred",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
      filter: { newOwner: owner, ...filterParam },
    }
  )

  const operators = Array.from(
    new Set(
      [...ownerDelegations, ...transferEventsByOwner].map(
        (_) => _.returnValues.operator
      )
    )
  )

  // Fetch `StakeOwnershipTransferred` by operator field. We need to check more recent event
  // to make sure the delegation ownership has not been transferred.
  let transferEventsByOperators = {}
  if (!isEmptyArray(operators)) {
    transferEventsByOperators = (
      await stakingContract.getPastEvents("StakeOwnershipTransferred", {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
        filter: { operator: operators },
      })
    ).reduce(reduceByOperator, {})
  }

  return operators.filter((operator) => {
    if (!transferEventsByOperators.hasOwnProperty(operator)) {
      return true
    }

    const transferEventsByOperator = transferEventsByOperators[operator]
    const latestTransfer =
      transferEventsByOperator[transferEventsByOperator.length - 1]
    if (
      latestTransfer &&
      isSameEthAddress(latestTransfer.returnValues.newOwner, owner)
    ) {
      return true
    }

    return false
  })
}

const reduceByOperator = (result, event) => {
  const {
    returnValues: { operator },
  } = event

  ;(result[operator] = result[operator] || []).push(event)

  return result
}
