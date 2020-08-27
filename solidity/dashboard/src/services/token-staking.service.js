import { contractService } from "./contracts.service"
import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
} from "../constants/constants"
import moment from "moment"
import {
  isCodeValid,
  createManagedGrantContractInstance,
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  ContractsLoaded,
} from "../contracts"
import { isSameEthAddress } from "../utils/general.utils"

const fetchDelegatedTokensData = async (web3Context) => {
  const { yourAddress, grantContract, eth, web3 } = web3Context
  const [
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
    initializationPeriod,
  ] = await Promise.all([
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "balanceOf",
      yourAddress
    ),
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "ownerOf",
      yourAddress
    ),
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "beneficiaryOf",
      yourAddress
    ),
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "authorizerOf",
      yourAddress
    ),
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "initializationPeriod"
    ),
  ])

  let isUndelegationFromGrant = true
  let grantStakeDetails
  try {
    grantStakeDetails = await grantContract.methods
      .getGrantStakeDetails(yourAddress)
      .call()
  } catch (error) {
    isUndelegationFromGrant = false
  }

  let isManagedGrant = false
  let managedGrantContractInstance
  if (isUndelegationFromGrant) {
    const { grantee } = await contractService.makeCall(
      web3Context,
      TOKEN_GRANT_CONTRACT_NAME,
      "getGrant",
      grantStakeDetails.grantId
    )
    // check if grantee is a contract
    const code = await eth.getCode(grantee)
    if (isCodeValid(code)) {
      managedGrantContractInstance = createManagedGrantContractInstance(
        web3,
        grantee
      )
      isManagedGrant = true
    }
  }

  const {
    undelegationStatus,
    undelegation,
    undelegationPeriod,
    delegationStatus,
    undelegationCompletedAt,
  } = await fetchPendingUndelegation(web3Context)
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
    isUndelegationFromGrant,
    isInInitializationPeriod,
    undelegationPeriod,
    isManagedGrant,
    managedGrantContractInstance,
    delegationStatus,
    undelegationCompletedAt,
  }
}

const fetchPendingUndelegation = async (web3Context) => {
  const { yourAddress } = web3Context
  const [delegation, undelegationPeriod] = await Promise.all([
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "getDelegationInfo",
      yourAddress
    ),
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "undelegationPeriod"
    ),
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

export const getOperatorsOfBeneficiary = async (web3Context, beneficiary) => {
  return (
    await contractService.getPastEvents(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "OperatorStaked",
      {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME],
        filter: { beneficiary },
      }
    )
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
  const transferEventsByOperators = (
    await stakingContract.getPastEvents("StakeOwnershipTransferred", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
      filter: { operator: operators },
    })
  ).reduce(reduceByOperator, {})

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
