import web3Utils from "web3-utils"
import { sub, gt } from "../utils/arithmetics.utils"
import moment from "moment"
import { tokenGrantsService } from "../services/token-grants.service"
import {
  createManagedGrantContractInstance,
  Web3Loaded,
  ContractsLoaded,
} from "../contracts"

export const fetchOldDelegations = async (web3Context) => {
  const web3 = await Web3Loaded
  const yourAddress = web3.eth.defaultAccount
  const { oldTokenStakingContract, grantContract } = await ContractsLoaded

  const operatorsAddresses = await oldTokenStakingContract.methods
    .operatorsOf(yourAddress)
    .call()

  const undelegationPeriod = await oldTokenStakingContract.methods
    .undelegationPeriod()
    .call()

  const initializationPeriod = await oldTokenStakingContract.methods
    .initializationPeriod()
    .call()

  const operatorsAddressesSet = new Set(operatorsAddresses)
  const granteeOperators = await grantContract.methods
    .getGranteeOperators(yourAddress)
    .call()

  const granteeOperatorsSet = new Set(granteeOperators)
  const managedGrantOperators = await tokenGrantsService.getOperatorsFromManagedGrants()

  const ownedDelegations = await getDelegations(
    operatorsAddressesSet,
    initializationPeriod,
    undelegationPeriod
  )

  const granteeDelegations = await getDelegations(
    granteeOperatorsSet,
    initializationPeriod,
    undelegationPeriod,
    true
  )
  const managedGrantsDelegations = await getDelegations(
    managedGrantOperators,
    initializationPeriod,
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
    initializationPeriod,
    undelegationPeriod,
  }
}

const getDelegations = async (
  operators,
  initializationPeriod,
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

    const beneficiary = await oldTokenStakingContract.methods
      .beneficiaryOf(operatorAddress)
      .call()
    const authorizerAddress = await oldTokenStakingContract.methods
      .authorizerOf(operatorAddress)
      .call()

    let grantId
    let managedGrantContractInstance
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
      operatorData.isUndelegation = false
      delegations.push(operatorData)
    }

    if (operatorData.undelegatedAt !== "0" && gt(amount, 0)) {
      operatorData.undelegationCompleteAt = moment
        .unix(undelegatedAt)
        .add(undelegationPeriod, "seconds")
      operatorData.canRecoverStake = operatorData.undelegationCompleteAt.isBefore(
        moment()
      )
      operatorData.isUndelegation = true
      delegations.push(operatorData)
    }
  }

  return delegations
}
