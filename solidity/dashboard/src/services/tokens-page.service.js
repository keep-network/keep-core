import web3Utils from "web3-utils"
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

export const fetchTokensPageData = async () => {
  const web3 = await Web3Loaded
  const yourAddress = web3.eth.defaultAccount

  const { stakingContract, token, grantContract } = await ContractsLoaded

  const keepTokenBalance = await token.methods.balanceOf(yourAddress).call()

  const grantTokenBalance = await grantContract.methods
    .balanceOf(yourAddress)
    .call()

  const minimumStake = await stakingContract.methods.minimumStake().call()

  const undelegationPeriod = await stakingContract.methods
    .undelegationPeriod()
    .call()

  const initializationPeriod = await stakingContract.methods
    .initializationPeriod()
    .call()

  const [
    ownedDelegations,
    ownedUndelegations,
    tokenStakingBalance,
    pendingUndelegationBalance,
  ] = await getOwnedDelegations(initializationPeriod, undelegationPeriod)

  const [
    granteeDelegations,
    granteeUndelegations,
    granteeGrantsIds,
  ] = await getGranteeDelegations(initializationPeriod, undelegationPeriod)

  const [
    managedGrantsDelegations,
    managedGrantsUndelegations,
    managedGrantsIds,
  ] = await getManagedGranteeDelegations(
    initializationPeriod,
    undelegationPeriod
  )

  const [copiedDelegations, copiedUndelegations] = await getCopiedDelegations(
    yourAddress,
    [...granteeGrantsIds, ...managedGrantsIds],
    initializationPeriod,
    undelegationPeriod
  )

  const delegations = [
    ...ownedDelegations,
    ...granteeDelegations,
    ...managedGrantsDelegations,
    ...copiedDelegations,
  ].sort((a, b) => sub(b.createdAt, a.createdAt))
  const undelegations = [
    ...ownedUndelegations,
    ...granteeUndelegations,
    ...managedGrantsUndelegations,
    ...copiedUndelegations,
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
      isCopiedStake,
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
      isCopiedStake,
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

  // No delegations
  if (isEmptyArray(operators)) {
    return await getDelegations({}, initializationPeriod, undelegationPeriod)
  }

  // Scan `OperatorStaked` event by operator(indexed param) to get authorizer and beneeficiary.
  const operatorToDetails = (
    await stakingContract.getPastEvents("OperatorStaked", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
      filter: { operator: operators },
    })
  ).reduce(toOperator, {})

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
  const {
    stakingContract,
    tokenStakingEscrow,
    stakingPortBackerContract,
  } = await ContractsLoaded

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

  const newOperators = escrowRedelegation.map((_) => _.returnValues.newOperator)
  const obsoleteOperators = escrowRedelegation.map(
    (_) => _.returnValues.previousOperator
  )

  let activeOperators = granteeOperators
    .filter((operator) => !obsoleteOperators.includes(operator))
    .concat(newOperators)

  const operatorsOfPortBacker = isEmptyArray(activeOperators)
    ? []
    : await getOperatorsOfOwner(
        stakingPortBackerContract.options.address,
        activeOperators
      )

  // We want to skip copied delegations
  activeOperators = activeOperators.filter(
    (operator) => !operatorsOfPortBacker.includes(operator)
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

  for (const operator of Object.keys(operatorsDetailsMap)) {
    const grantId = newOperatorToGrantId.hasOwnProperty(operator)
      ? newOperatorToGrantId[operator]
      : null
    operatorsDetailsMap[operator] = {
      ...operatorsDetailsMap[operator],
      isFromGrant: true,
      isManagedGrant,
      grantId,
    }
  }

  return operatorsDetailsMap
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
    grantIds,
    true
  )

  const [delegations, undelegations] = await getDelegations(
    allOperators,
    initializationPeriod,
    undelegationPeriod
  )

  return [delegations, undelegations, grantIds]
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
    const { grantee } = await grantContract.methods.getGrant(grantId).call()
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

  const [delegations, undelegations] = await getDelegations(
    allOperators,
    initializationPeriod,
    undelegationPeriod
  )

  return [delegations, undelegations, granteeGrants]
}

export const getCopiedDelegations = async (
  ownerOrGrantee,
  grantIds,
  initializationPeriod,
  undelegationPeriod
) => {
  const {
    stakingPortBackerContract,
    grantContract,
    stakingContract,
  } = await ContractsLoaded

  const operatorsToCheck = (
    await stakingPortBackerContract.getPastEvents("StakeCopied", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingPortBackerContract,
      filter: { owner: ownerOrGrantee },
    })
  ).map((_) => _.returnValues.operator)

  // No delegations
  if (isEmptyArray(operatorsToCheck)) {
    return await getDelegations({}, initializationPeriod, undelegationPeriod)
  }

  // We only want operators for whom the delegation has not been paid back.
  const operatorsOfPortBacker = await getOperatorsOfOwner(
    stakingPortBackerContract.options.address,
    operatorsToCheck
  )

  // Scan `OperatorStaked` event by operator(indexed param) to get authorizer and beneficiary.
  const operatorToDetails = (isEmptyArray(operatorsOfPortBacker)
    ? []
    : await stakingContract.getPastEvents("OperatorStaked", {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
        filter: { operator: operatorsOfPortBacker },
      })
  ).reduce((reducer, _) => {
    reducer[_.returnValues.operator] = _.returnValues
    reducer[_.returnValues.operator].isCopiedStake = true

    return reducer
  }, {})

  // Fill operator's grant details for delegations created from a grant.
  if (!isEmptyArray(grantIds)) {
    const tokenGrantStakingEvents = await grantContract.getPastEvents(
      "TokenGrantStaked",
      {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.grantContract,
        filter: { grantId: grantIds },
      }
    )

    for (const grantStakedEvent of tokenGrantStakingEvents) {
      const operator = grantStakedEvent.returnValues.operator
      const grantId = grantStakedEvent.returnValues.grantId
      if (operatorToDetails.hasOwnProperty(operator)) {
        operatorToDetails[operator].grantId = grantId
        operatorToDetails[operator].isFromGrant = true
      }
    }
  }

  const [delegations, undelegations] = await getDelegations(
    operatorToDetails,
    initializationPeriod,
    undelegationPeriod
  )

  return [delegations, undelegations]
}

const toOperator = (reducer, _) => {
  reducer[_.returnValues.operator] = _.returnValues
  return reducer
}

export const tokensPageService = {
  fetchTokensPageData,
}
