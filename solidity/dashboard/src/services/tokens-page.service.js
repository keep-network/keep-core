import { contractService } from "./contracts.service"
import web3Utils from "web3-utils"
import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  KEEP_TOKEN_CONTRACT_NAME,
} from "../constants/constants"
import { sub, gt } from "../utils/arithmetics.utils"
import moment from "moment"
import { tokenGrantsService } from "../services/token-grants.service"
import {
  createManagedGrantContractInstance,
  CONTRACT_DEPLOY_BLOCK_NUMBER,
} from "../contracts"
import { ContractsLoaded, Web3Loaded } from "../contracts"
import { isEmptyArray } from "../utils/array.utils"
import { getOperatorsOfOwner } from "./token-staking.service"
import { isSameEthAddress } from "../utils/general.utils"

export const fetchTokensPageData = async (web3Context) => {
  const { yourAddress } = web3Context

  const [
    keepTokenBalance,
    grantTokenBalance,
    minimumStake,
    undelegationPeriod,
    initializationPeriod,
  ] = await Promise.all([
    contractService.makeCall(
      web3Context,
      KEEP_TOKEN_CONTRACT_NAME,
      "balanceOf",
      yourAddress
    ),
    contractService.makeCall(
      web3Context,
      TOKEN_GRANT_CONTRACT_NAME,
      "balanceOf",
      yourAddress
    ),
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "minimumStake"
    ),
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "undelegationPeriod"
    ),
    contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "initializationPeriod"
    ),
  ])

  const [
    ownedDelegations,
    ownedUndelegations,
    tokenStakingBalance,
    pendingUndelegationBalance,
  ] = await getOwnedDelegations(initializationPeriod, undelegationPeriod)

  const [
    granteeDelegations,
    granteeUndelegations,
  ] = await getGranteeDelegations(initializationPeriod, undelegationPeriod)
  const [
    managedGrantsDelegations,
    managedGrantsUndelegations,
  ] = await getManagedGranteeDelegations(
    initializationPeriod,
    undelegationPeriod
  )

  const delegations = [
    ...ownedDelegations,
    ...granteeDelegations,
    ...managedGrantsDelegations,
  ].sort((a, b) => sub(b.createdAt, a.createdAt))
  const undelegations = [
    ...ownedUndelegations,
    ...granteeUndelegations,
    ...managedGrantsUndelegations,
  ].sort((a, b) => sub(b.undelegatedAt, a.undelegatedAt))

  return {
    delegations,
    undelegations,
    keepTokenBalance,
    grantTokenBalance,
    ownedTokensDelegationsBalance: tokenStakingBalance.toString(),
    ownedTokensUndelegationsBalance: pendingUndelegationBalance.toString(),
    minimumStake,
    initializationPeriod,
    undelegationPeriod,
  }
}

const getDelegations = async (
  operatorToDetailsMap,
  initializationPeriod,
  undelegationPeriod
) => {
  const web3 = await Web3Loaded
  const { stakingContract, grantContract } = await ContractsLoaded

  let tokenStakingBalance = web3Utils.toBN(0)
  let pendingUndelegationBalance = web3Utils.toBN(0)
  const delegations = []
  const undelegations = []

  for (const [operatorAddress, details] of Object.entries(
    operatorToDetailsMap
  )) {
    let {
      beneficiary,
      authorizer: authorizerAddress,
      isFromGrant,
      isManagedGrant,
      grantId,
    } = details

    const {
      createdAt,
      undelegatedAt,
      amount,
    } = await stakingContract.methods.getDelegationInfo(operatorAddress).call()

    let managedGrantContractInstance = null
    if (isFromGrant && !grantId) {
      try {
        const grantStakeDetails = await grantContract.methods
          .getGrantStakeDetails(operatorAddress)
          .call()
        grantId = grantStakeDetails.grantId
      } catch (error) {
        grantId = null
      }
    }

    if (isManagedGrant && grantId) {
      const { grantee } = await grantContract.methods.getGrant(grantId).call()
      managedGrantContractInstance = createManagedGrantContractInstance(
        web3,
        grantee
      )
    }

    const operatorData = {
      undelegatedAt,
      amount,
      beneficiary,
      operatorAddress,
      createdAt,
      authorizerAddress,
      isFromGrant,
      grantId,
      isManagedGrant,
      managedGrantContractInstance,
    }
    const balance = web3Utils.toBN(amount)

    if (!balance.isZero() && operatorData.undelegatedAt === "0") {
      const initializationOverAt = moment
        .unix(createdAt)
        .add(initializationPeriod, "seconds")
      operatorData.isInInitializationPeriod = moment().isSameOrBefore(
        initializationOverAt
      )
      operatorData.initializationOverAt = initializationOverAt
      delegations.push(operatorData)
      if (!isFromGrant) {
        tokenStakingBalance = tokenStakingBalance.add(balance)
      }
    }
    if (operatorData.undelegatedAt !== "0" && gt(amount, 0)) {
      operatorData.undelegationCompleteAt = moment
        .unix(undelegatedAt)
        .add(undelegationPeriod, "seconds")
      operatorData.canRecoverStake = operatorData.undelegationCompleteAt.isBefore(
        moment()
      )
      if (!isFromGrant) {
        pendingUndelegationBalance = pendingUndelegationBalance.add(balance)
      }
      undelegations.push(operatorData)
    }
  }

  return [
    delegations,
    undelegations,
    tokenStakingBalance,
    pendingUndelegationBalance,
  ]
}

const getOwnedDelegations = async (
  initializationPeriod,
  undelegationPeriod
) => {
  const web3 = await Web3Loaded
  const yourAddress = web3.eth.defaultAccount
  const { stakingContract } = await ContractsLoaded

  // Get operators
  const operators = await getOperatorsOfOwner(yourAddress)

  // Scan `OperatorStaked` event by operator(indexed param) to get authorizer and beneeficiary.
  let operatorToDetails = {}
  if (!isEmptyArray(operators)) {
    operatorToDetails = (
      await stakingContract.getPastEvents("OperatorStaked", {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
        filter: { operator: operators },
      })
    ).reduce(toOperator, {})
  }

  return await getDelegations(
    operatorToDetails,
    initializationPeriod,
    undelegationPeriod
  )
}

const getAllGranteeOperators = async (
  granteeOperators,
  grantIds,
  isManagedGrant = false
) => {
  const { stakingContract, tokenStakingEscrow } = await ContractsLoaded

  let escrowRedelegation = []
  if (!isEmptyArray(grantIds)) {
    escrowRedelegation = await tokenStakingEscrow.getPastEvents(
      "DepositRedelegated",
      {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.tokenStakingEscrow,
        filter: {
          grantId: grantIds,
        },
      }
    )
  }

  const newOperatorToGrantId = escrowRedelegation.reduce((reducer, event) => {
    const {
      returnValues: { newOperator, grantId },
    } = event
    reducer[newOperator] = grantId
    return reducer
  }, {})

  let newOperatorsStakeDetails = {}
  if (!isEmptyArray(Object.keys(newOperatorToGrantId))) {
    newOperatorsStakeDetails = (
      await stakingContract.getPastEvents("OperatorStaked", {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
        filter: { operator: Object.keys(newOperatorToGrantId) },
      })
    ).reduce(toOperator, {})
  }

  const oldOperators = escrowRedelegation
    .filter((_) => {
      return newOperatorsStakeDetails.hasOwnProperty(_.returnValues.newOperator)
    })
    .map((_) => _.returnValues.previousOperator)

  const activeOperators = granteeOperators.filter(
    (operator) => !oldOperators.includes(operator)
  )

  let operatorsDetailsMap = {}
  if (!isEmptyArray(activeOperators)) {
    operatorsDetailsMap = (operatorsDetailsMap = await stakingContract.getPastEvents(
      "OperatorStaked",
      {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
        filter: { operator: activeOperators },
      }
    )).reduce(toOperator, {})
  }

  const allOperatorsDetails = {
    ...operatorsDetailsMap,
    ...newOperatorsStakeDetails,
  }

  for (const operator of Object.keys(allOperatorsDetails)) {
    const grantId = newOperatorToGrantId.hasOwnProperty(operator)
      ? newOperatorToGrantId[operator]
      : null
    allOperatorsDetails[operator] = {
      ...allOperatorsDetails[operator],
      isFromGrant: true,
      isManagedGrant,
      grantId,
    }
  }

  return allOperatorsDetails
}

const getManagedGranteeDelegations = async (
  initializationPeriod,
  undelegationPeriod
) => {
  const { grantContract } = await ContractsLoaded

  const managedGrants = await tokenGrantsService.fetchManagedGrants()

  const grantIds = managedGrants.map(({ grantId }) => grantId)

  const operators = new Set()

  for (const managedGrant of managedGrants) {
    const { managedGrantContractInstance } = managedGrant
    const granteeAddress = managedGrantContractInstance.options.address
    const grenteeOperators = await grantContract.methods
      .getGranteeOperators(granteeAddress)
      .call()
    grenteeOperators.forEach(operators.add, operators)
  }

  const allOperators = await getAllGranteeOperators(
    Array.from(operators),
    Array.from(grantIds),
    true
  )

  return await getDelegations(
    allOperators,
    initializationPeriod,
    undelegationPeriod
  )
}

const getGranteeDelegations = async (
  initializationPeriod,
  undelegationPeriod
) => {
  const web3 = await Web3Loaded
  const yourAddress = web3.eth.defaultAccount
  const { grantContract } = await ContractsLoaded

  // `getGrants` function returns grants for a grant manager or grantee.
  // So it's possible that the provided address is a grantee in grant A and a grant manager in grant B.
  // In that case a `getGrants` function returns [A, B].
  const grantIds = new Set(
    await grantContract.methods.getGrants(yourAddress).call()
  )

  // Filter out grants. We just want grants from the grantee's perspective
  const granteeGrants = []
  for (const grantId of grantIds) {
    const { grantee } = await grantContract.methods.getGrant(grantId)
    if (isSameEthAddress(grantee, yourAddress)) {
      granteeGrants.push(grantId)
    }
  }

  const granteeOperators = new Set(
    await grantContract.methods.getGranteeOperators(yourAddress).call()
  )

  const allOperators = await getAllGranteeOperators(
    Array.from(granteeOperators),
    granteeGrants,
    false
  )

  return await getDelegations(
    allOperators,
    initializationPeriod,
    undelegationPeriod
  )
}

const toOperator = (reducer, _) => {
  reducer[_.returnValues.operator] = _.returnValues
  return reducer
}

export const tokensPageService = {
  fetchTokensPageData,
}
