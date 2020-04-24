import { contractService } from './contracts.service'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import moment from 'moment'
import { COMPLETE_STATUS, PENDING_STATUS } from '../constants/constants'

const fetchDelegatedTokensData = async (web3Context) => {
  const { yourAddress, grantContract } = web3Context
  const [
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
    initializationPeriod,
  ] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'balanceOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'ownerOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'beneficiaryOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'authorizerOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'initializationPeriod'),
  ])

  let isUndelegationFromGrant = true
  try {
    await grantContract.methods.getGrantStakeDetails(yourAddress).call()
  } catch (error) {
    isUndelegationFromGrant = false
  }

  const { undelegationStatus, undelegation, undelegationPeriod } = await fetchPendingUndelegation(web3Context)
  const { createdAt } = undelegation
  const initializationOverAt = moment.unix(createdAt).add(initializationPeriod, 'seconds')
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
  }
}

const fetchPendingUndelegation = async (web3Context) => {
  const { yourAddress } = web3Context
  const [delegation, undelegationPeriod] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'getDelegationInfo', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'undelegationPeriod'),
  ])

  const { undelegatedAt } = delegation

  const isUndelegation = delegation.undelegatedAt !== '0'
  const pendingUnstakeBalance = isUndelegation ? delegation.amount : 0
  const undelegationCompletedAt = isUndelegation ? moment.unix(undelegatedAt).add(undelegationPeriod, 'seconds') : null
  let undelegationStatus
  if (isUndelegation) {
    undelegationStatus = undelegationCompletedAt.isBefore(moment()) ? COMPLETE_STATUS : PENDING_STATUS
  } else if (delegation.undelegatedAt === '0' && delegation.createdAt !== '0' && delegation.amount === '0') {
    undelegationStatus = COMPLETE_STATUS
  }

  return {
    pendingUnstakeBalance,
    undelegationCompletedAt,
    undelegationPeriod,
    undelegationStatus,
    undelegation: delegation,
  }
}

export const operatorService = {
  fetchDelegatedTokensData,
  fetchPendingUndelegation,
}
