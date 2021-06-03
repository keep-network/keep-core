import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
  STAKING_PORT_BACKER_CONTRACT_NAME,
} from "../constants/constants"
import moment from "moment"
import {
  isCodeValid,
  createManagedGrantContractInstance,
  getContractDeploymentBlockNumber,
  ContractsLoaded,
  Web3Loaded,
} from "../contracts"
import { isSameEthAddress } from "../utils/general.utils"
import { isEmptyArray } from "../utils/array.utils"
import { getEventsFromTransaction, ZERO_ADDRESS } from "../utils/ethereum.utils"
import { tokenGrantsService } from "./token-grants.service"

const delegationInfoFromStakedEvents = async (address) => {
  const { stakingContract } = await ContractsLoaded
  let operatorStakedEvents
  try {
    operatorStakedEvents = await stakingContract.getPastEvents(
      "OperatorStaked",
      {
        fromBlock: await getContractDeploymentBlockNumber(
          TOKEN_STAKING_CONTRACT_NAME
        ),
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

const fetchDelegatedTokensData = async (address) => {
  const { grantContract, stakingContract } = await ContractsLoaded
  const web3 = await Web3Loaded
  const { eth } = web3
  let ownerAddress

  const {
    stakingTransactionHash,
    beneficiaryAddress,
    authorizerAddress,
  } = await delegationInfoFromStakedEvents(address)

  const [stakedBalance, initializationPeriod] = await Promise.all([
    stakingContract.methods.balanceOf(address).call(),
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
    ownerAddress = await stakingContract.methods.ownerOf(address).call()
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

export const getOperatorsOfAuthorizer = async (authorizer) => {
  const { stakingContract } = await ContractsLoaded
  return (
    await stakingContract.getPastEvents("OperatorStaked", {
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_STAKING_CONTRACT_NAME
      ),
      filter: { authorizer },
    })
  ).map((_) => _.returnValues.operator)
}

export const getOperatorsOfBeneficiary = async (beneficiary) => {
  const { stakingContract } = await ContractsLoaded

  return (
    await stakingContract.getPastEvents("OperatorStaked", {
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_STAKING_CONTRACT_NAME
      ),
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
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_STAKING_CONTRACT_NAME
      ),
      filter: { owner, ...filterParam },
    }
  )

  const transferEventsByOwner = await stakingContract.getPastEvents(
    "StakeOwnershipTransferred",
    {
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_STAKING_CONTRACT_NAME
      ),
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

  // Fetch `StakeOwnershipTransferred` by operator field. We need to check more
  // recent event to make sure the delegation ownership has not been
  // transferred.
  let transferEventsByOperators = {}
  if (!isEmptyArray(operators)) {
    transferEventsByOperators = (
      await stakingContract.getPastEvents("StakeOwnershipTransferred", {
        fromBlock: await getContractDeploymentBlockNumber(
          TOKEN_STAKING_CONTRACT_NAME
        ),
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

export const getOperatorsOfGrantee = async (address) => {
  const { grantContract } = await ContractsLoaded

  // The `getGrants` function returns grants for a grant manager or grantee. So
  // it's possible that the provided address is a grantee in grant A and a grant
  // manager in grant B. In that case a `getGrants` function returns [A, B].
  const grantIds = new Set(
    await grantContract.methods.getGrants(address).call()
  )

  // Filter out grants. We just want grants from the grantee's perspective
  const granteeGrants = []
  for (const grantId of grantIds) {
    const { grantee } = await grantContract.methods.getGrant(grantId).call()
    if (isSameEthAddress(grantee, address)) {
      granteeGrants.push(grantId)
    }
  }

  const granteeOperators = new Set(
    await grantContract.methods.getGranteeOperators(address).call()
  )

  const {
    allOperators,
    operatorToGrantDetailsMap,
  } = await getAllGranteeOperators(
    Array.from(granteeOperators),
    granteeGrants,
    address
  )

  return { allOperators, granteeGrants, operatorToGrantDetailsMap }
}

export const getOperatorsOfManagedGrantee = async (address) => {
  const { grantContract } = await ContractsLoaded

  const managedGrants = await tokenGrantsService.fetchManagedGrants(address)

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

  const {
    allOperators,
    operatorToGrantDetailsMap,
  } = await getAllGranteeOperators(
    Array.from(operators),
    grantIds,
    address,
    true
  )

  return { allOperators, granteeGrants: grantIds, operatorToGrantDetailsMap }
}

/**
 * The `getGranteeOperators` function returns only operators stored in the
 * `TokenGrant` contract. If the grant delegation will be canceled/revoked,
 * tokens go to the `TokenStakingEscrow` contract. The grantee can redelagte
 * tokens via `TokenStakingEscrow` and  only the `TokenStakingEscrow` knows
 * about the new operators so we need to take into account redelegations from
 * `TokenStakingEscrow`.
 *
 * @param {string[]} granteeOperators The result of the
 * `TokenGrant::getGranteeOperators`.
 * @param {string[]} grantIds Array of all grantee grants ids.
 * @param {string} grantee Grantee address.
 * @param {boolean} isManagedGrant The flag informs that grants in `grantIds`
 * param are managed grants.
 * @return {Promise<string[]>} Array of all grantee operators.
 */
const getAllGranteeOperators = async (
  granteeOperators,
  grantIds,
  grantee,
  isManagedGrant = false
) => {
  const {
    tokenStakingEscrow,
    stakingPortBackerContract,
    stakingContract,
  } = await ContractsLoaded

  let escrowRedelegation = []
  if (!isEmptyArray(grantIds)) {
    escrowRedelegation = await tokenStakingEscrow.getPastEvents(
      "DepositRedelegated",
      {
        fromBlock: await getContractDeploymentBlockNumber(
          TOKEN_STAKING_ESCROW_CONTRACT_NAME
        ),
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

  // Copied delegations but not yet paid back. The owner of the not paid back
  // delegation is `StakingPortBacker` contract.
  const operatorsOfPortBacker = isEmptyArray(activeOperators)
    ? []
    : await getOperatorsOfOwner(
        stakingPortBackerContract.options.address,
        activeOperators
      )

  // Fetching paid back delegations. Paid-back delegations are as the
  // delegations from liquid tokens in the new `TokenStaking` contract. So we
  // want to display them as delegations from liquid tokens ,not from a grant.
  const paidBackDelegations = isEmptyArray(activeOperators)
    ? []
    : (
        await stakingPortBackerContract.getPastEvents("StakePaidBack", {
          fromBlock: await getContractDeploymentBlockNumber(
            STAKING_PORT_BACKER_CONTRACT_NAME
          ),
          filter: { owner: grantee, operator: activeOperators },
        })
      ).map((_) => _.returnValues.operator)

  const operatorsToFilterOut = [
    ...operatorsOfPortBacker,
    ...paidBackDelegations,
  ]

  // We want to skip copied delegations.
  activeOperators = activeOperators.filter(
    (operator) => !operatorsToFilterOut.includes(operator)
  )

  let operatorToGrantDetailsMap = {}
  if (!isEmptyArray(activeOperators)) {
    operatorToGrantDetailsMap = (
      await stakingContract.getPastEvents("OperatorStaked", {
        fromBlock: await getContractDeploymentBlockNumber(
          TOKEN_STAKING_CONTRACT_NAME
        ),
        filter: { operator: activeOperators },
      })
    ).reduce((reducer, _) => {
      reducer[_.returnValues.operator] = _.returnValues
      return reducer
    }, {})
  }

  for (const operator of Object.keys(operatorToGrantDetailsMap)) {
    const grantId = newOperatorToGrantId.hasOwnProperty(operator)
      ? newOperatorToGrantId[operator]
      : null
    operatorToGrantDetailsMap[operator] = {
      ...operatorToGrantDetailsMap[operator],
      isFromGrant: true,
      isManagedGrant,
      grantId,
    }
  }

  return { allOperators: activeOperators, operatorToGrantDetailsMap }
}

export const getOperatorsOfCopiedDelegations = async (ownerOrGrantee) => {
  const { stakingPortBackerContract } = await ContractsLoaded
  const operatorsToCheck = (
    await stakingPortBackerContract.getPastEvents("StakeCopied", {
      fromBlock: await getContractDeploymentBlockNumber(
        STAKING_PORT_BACKER_CONTRACT_NAME
      ),
      filter: { owner: ownerOrGrantee },
    })
  ).map((_) => _.returnValues.operator)

  // No copied delegations.
  if (isEmptyArray(operatorsToCheck)) {
    return []
  }

  // We only want operators for whom the delegation has not been paid back. The
  // owner of the not paid back delegation is `StakingPortBacker` contract.
  return await getOperatorsOfOwner(
    stakingPortBackerContract.options.address,
    operatorsToCheck
  )
}

const reduceByOperator = (result, event) => {
  const {
    returnValues: { operator },
  } = event

  ;(result[operator] = result[operator] || []).push(event)

  return result
}

export const isDelegationExists = async (operator) => {
  const { stakingContract } = await ContractsLoaded

  const owner = await stakingContract.methods.ownerOf(operator).call()

  return owner !== ZERO_ADDRESS
}
