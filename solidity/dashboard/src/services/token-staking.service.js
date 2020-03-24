import { contractService } from './contracts.service'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import web3Utils from 'web3-utils'
import { COMPLETE_STATUS, PENDING_STATUS } from '../constants/constants'
import { gt, add } from '../utils/arithmetics.utils'

const fetchDelegatedTokensData = async (web3Context) => {
  const { yourAddress, grantContract, eth } = web3Context
  const [
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
    initializationPeriod,
  ] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'balanceOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'ownerOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'magpieOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'authorizerOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'initializationPeriod'),
  ])

  let isUndelegationFromGrant = true
  try {
    await grantContract.methods.getGrantStakeDetails(yourAddress).call()
  } catch (error) {
    isUndelegationFromGrant = false
  }

  const { undelegationStatus, undelegation } = await fetchPendingUndelegation(web3Context)
  const { createdAt } = undelegation
  const initializationOverAt = add(createdAt || 0, initializationPeriod)
  const isInInitializationPeriod = gt(initializationOverAt, await eth.getBlockNumber())

  return {
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
    undelegationStatus,
    isUndelegationFromGrant,
    isInInitializationPeriod,
  }
}

const fetchPendingUndelegation = async (web3Context) => {
  const { yourAddress, eth } = web3Context
  const [delegation, undelegationPeriod] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'getDelegationInfo', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'undelegationPeriod'),
  ])

  const undelegationCompletedAtInBN = web3Utils.toBN(delegation.undelegatedAt).add(web3Utils.toBN(undelegationPeriod))
  const isUndelegation = delegation.undelegatedAt !== '0'
  const pendingUnstakeBalance = isUndelegation ? delegation.amount : 0
  const undelegationComplete = isUndelegation ? undelegationCompletedAtInBN.toString() : null
  let undelegationStatus
  if (isUndelegation) {
    undelegationStatus = gt(await eth.getBlockNumber(), undelegationCompletedAtInBN) ? COMPLETE_STATUS : PENDING_STATUS
  } else if (delegation.undelegatedAt === '0' && delegation.createdAt !== '0' && delegation.amount === '0') {
    undelegationStatus = COMPLETE_STATUS
  }

  return {
    pendingUnstakeBalance,
    undelegationComplete,
    undelegationPeriod,
    undelegationStatus,
    undelegation: delegation,
  }
}

export const operatorService = {
  fetchDelegatedTokensData,
  fetchPendingUndelegation,
}
