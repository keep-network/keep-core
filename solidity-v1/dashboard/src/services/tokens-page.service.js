import web3Utils from "web3-utils"
import moment from "moment"
import {
  createManagedGrantContractInstance,
  getContractDeploymentBlockNumber,
  Keep,
} from "../contracts"
import { ContractsLoaded, Web3Loaded } from "../contracts"
import {
  getOperatorsOfOwner,
  getOperatorsOfCopiedDelegations,
  getOperatorsOfGrantee,
  getOperatorsOfManagedGrantee,
} from "./token-staking.service"
import { sub, gt } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
} from "../constants/constants"

export const fetchTokensPageData = async () => {
  const web3 = await Web3Loaded
  const yourAddress = web3.eth.defaultAccount

  const { stakingContract, token } = await ContractsLoaded

  const keepTokenBalance = await token.methods.balanceOf(yourAddress).call()

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

  const [granteeDelegations, granteeUndelegations, granteeGrantsIds] =
    await getGranteeDelegations(
      yourAddress,
      initializationPeriod,
      undelegationPeriod
    )

  const [
    managedGrantsDelegations,
    managedGrantsUndelegations,
    managedGrantsIds,
  ] = await getManagedGranteeDelegations(
    yourAddress,
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

    const { createdAt, undelegatedAt, amount } = await stakingContract.methods
      .getDelegationInfo(operatorAddress)
      .call()

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

    const thresholdTokenStakingContractAddress =
      Keep.thresholdStakingContract.address

    const isTStakingContractAuthorized = await Keep.stakingContract.methods
      .isAuthorizedForOperator(
        operatorAddress,
        thresholdTokenStakingContractAddress
      )
      .call()

    const hasKeepTokensStakedInTNetwork =
      await Keep.keepToTStaking.hasKeepTokensStakedInTNetwork(operatorAddress)

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
      isTStakingContractAuthorized,
      isStakedToT: hasKeepTokensStakedInTNetwork,
    }
    const balance = web3Utils.toBN(amount)

    if (!balance.isZero() && operatorData.undelegatedAt === "0") {
      const initializationOverAt = moment
        .unix(createdAt)
        .add(initializationPeriod, "seconds")
      operatorData.isInInitializationPeriod =
        moment().isSameOrBefore(initializationOverAt)
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
      operatorData.canRecoverStake =
        operatorData.undelegationCompleteAt.isBefore(moment())
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
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_STAKING_CONTRACT_NAME
      ),
      filter: { operator: operators },
    })
  ).reduce(toOperator, {})

  return await getDelegations(
    operatorToDetails,
    initializationPeriod,
    undelegationPeriod
  )
}

const getManagedGranteeDelegations = async (
  address,
  initializationPeriod,
  undelegationPeriod
) => {
  const { operatorToGrantDetailsMap, granteeGrants } =
    await getOperatorsOfManagedGrantee(address)

  const [delegations, undelegations] = await getDelegations(
    operatorToGrantDetailsMap,
    initializationPeriod,
    undelegationPeriod
  )

  return [delegations, undelegations, granteeGrants]
}

const getGranteeDelegations = async (
  address,
  initializationPeriod,
  undelegationPeriod
) => {
  const { granteeGrants, operatorToGrantDetailsMap } =
    await getOperatorsOfGrantee(address)

  const [delegations, undelegations] = await getDelegations(
    operatorToGrantDetailsMap,
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
  const { grantContract, stakingContract } = await ContractsLoaded

  const operators = await getOperatorsOfCopiedDelegations(ownerOrGrantee)

  // No delegations.
  if (isEmptyArray(operators)) {
    return await getDelegations({}, initializationPeriod, undelegationPeriod)
  }

  // Scan `OperatorStaked` event by operator(indexed param) to get authorizer
  // and beneficiary.
  const operatorToDetails = (
    isEmptyArray(operators)
      ? []
      : await stakingContract.getPastEvents("OperatorStaked", {
          fromBlock: await getContractDeploymentBlockNumber(
            TOKEN_STAKING_CONTRACT_NAME
          ),
          filter: { operator: operators },
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
        fromBlock: await getContractDeploymentBlockNumber(
          TOKEN_GRANT_CONTRACT_NAME
        ),
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
