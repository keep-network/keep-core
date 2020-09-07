import web3Utils from "web3-utils"
import { sub } from "../utils/arithmetics.utils"
import moment from "moment"
import { tokenGrantsService } from "./token-grants.service"
import {
  createManagedGrantContractInstance,
  Web3Loaded,
  ContractsLoaded,
  CONTRACT_DEPLOY_BLOCK_NUMBER,
} from "../contracts"
import { isEmptyArray } from "../utils/array.utils"

const filterOutByOperator = (toFilterOut) => (operator) =>
  !toFilterOut.includes(operator)

export const fetchOldDelegations = async () => {
  const web3 = await Web3Loaded
  const yourAddress = web3.eth.defaultAccount
  const {
    oldTokenStakingContract,
    grantContract,
    stakingPortBackerContract,
  } = await ContractsLoaded

  // We want to skip the already copied stakes. To get copied stakes we should scan
  // the `StakedCopied` event from the `StakingPortBacker` contract.
  const copiedStakesOperator = (
    await stakingPortBackerContract.getPastEvents("StakeCopied", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingPortBackerContract,
      filter: { owner: yourAddress },
    })
  ).map((_) => _.returnValues.operator)

  const operatorsAddresses = (
    await oldTokenStakingContract.methods.operatorsOf(yourAddress).call()
  ).filter(filterOutByOperator(copiedStakesOperator))
  const operatorsAddressesSet = new Set(operatorsAddresses)

  const undelegationPeriod = await oldTokenStakingContract.methods
    .undelegationPeriod()
    .call()

  const granteeOperators = (
    await grantContract.methods.getGranteeOperators(yourAddress).call()
  ).filter(filterOutByOperator(copiedStakesOperator))
  const granteeOperatorsSet = new Set(granteeOperators)

  const managedGrantOperators = (
    await tokenGrantsService.getOperatorsFromManagedGrants()
  ).filter(filterOutByOperator(copiedStakesOperator))

  // We want to skip delegations that were undelegated after `StakingPortBacker` deploy.
  let operatorsToSkip = []
  if (
    operatorsAddressesSet.size > 0 ||
    granteeOperatorsSet.size > 0 ||
    !isEmptyArray(managedGrantOperators)
  ) {
    operatorsToSkip = (
      await oldTokenStakingContract.getPastEvents("Undelegated", {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingPortBackerContract,
        filter: {
          operator: [
            ...operatorsAddressesSet,
            ...granteeOperatorsSet,
            ...managedGrantOperators,
          ],
        },
      })
    ).map((_) => _.returnValues.operator)
  }

  const filteredOwnedOperators = Array.from(operatorsAddressesSet).filter(
    filterOutByOperator(operatorsToSkip)
  )
  const filteredGranteeOperators = Array.from(granteeOperatorsSet).filter(
    filterOutByOperator(operatorsToSkip)
  )
  const filteredManagedGrantOperators = managedGrantOperators.filter(
    filterOutByOperator(operatorsToSkip)
  )

  const ownedDelegations = await getDelegations(
    filteredOwnedOperators,
    undelegationPeriod
  )

  const granteeDelegations = await getDelegations(
    filteredGranteeOperators,
    undelegationPeriod,
    true
  )

  const managedGrantsDelegations = await getDelegations(
    filteredManagedGrantOperators,
    undelegationPeriod,
    true,
    true
  )

  const delegations = [
    ...ownedDelegations,
    ...granteeDelegations,
    ...managedGrantsDelegations,
  ].sort((a, b) => sub(b.createdAt, a.createdAt))

  return {
    delegations,
    undelegationPeriod,
  }
}

const getDelegations = async (
  operators,
  undelegationPeriod,
  isFromGrant,
  isManagedGrant
) => {
  const web3 = await Web3Loaded
  const { oldTokenStakingContract, grantContract } = await ContractsLoaded

  const delegations = []

  for (const operatorAddress of operators) {
    const {
      createdAt,
      undelegatedAt,
      amount,
    } = await oldTokenStakingContract.methods
      .getDelegationInfo(operatorAddress)
      .call()

    const beneficiaryAddress = await oldTokenStakingContract.methods
      .beneficiaryOf(operatorAddress)
      .call()
    const authorizerAddress = await oldTokenStakingContract.methods
      .authorizerOf(operatorAddress)
      .call()

    let grantId = null
    let managedGrantContractInstance = null
    if (isFromGrant) {
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
      beneficiaryAddress,
      operatorAddress,
      createdAt,
      authorizerAddress,
      isFromGrant,
      grantId,
      isManagedGrant,
      managedGrantContractInstance,
    }

    const delegationAmount = web3Utils.toBN(amount)
    // StakingPortBacker contract requires the delegation amount to be greater than zero.
    if (delegationAmount.gt(web3Utils.toBN(0))) {
      // Check if the stake is currently undelegating.
      if (operatorData.undelegatedAt === "0") {
        operatorData.isUndelegating = false
      } else {
        operatorData.undelegationCompleteAt = moment
          .unix(undelegatedAt)
          .add(undelegationPeriod, "seconds")
        operatorData.canRecoverStake = operatorData.undelegationCompleteAt.isBefore(
          moment()
        )
        operatorData.isUndelegating = true
      }

      delegations.push(operatorData)
    }
  }

  return delegations
}
