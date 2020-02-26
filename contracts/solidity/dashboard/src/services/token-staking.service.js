import { contractService } from './contracts.service'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import web3Utils from 'web3-utils'

const fetchDelegatedTokensData = async (web3Context) => {
  const { yourAddress } = web3Context
  const [
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
  ] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'balanceOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'ownerOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'magpieOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'authorizerOf', yourAddress),
  ])

  const { undelegationStatus } = await fetchPendingUndelegation(web3Context)

  return {
    stakedBalance,
    ownerAddress,
    beneficiaryAddress,
    authorizerAddress,
    undelegationStatus,
  }
}

const fetchPendingUndelegation = async (web3Context) => {
  const { yourAddress, eth } = web3Context
  const [undelegation, undelegationPeriod] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'getDelegationInfo', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'undelegationPeriod'),
  ])

  const undelegationCompletedAtInBN = web3Utils.toBN(undelegation.undelegatedAt).add(web3Utils.toBN(undelegationPeriod))
  const isUndelegation = undelegation.undelegatedAt !== '0'
  const pendingUnstakeBalance = isUndelegation ? undelegation.amount : 0
  const undelegationComplete = isUndelegation ? undelegationCompletedAtInBN.toString() : null
  let undelegationStatus
  if (isUndelegation) {
    undelegationStatus = web3Utils.toBN(await eth.getBlockNumber()).gt(undelegationCompletedAtInBN) ? 'COMPLETE' : 'PENDING'
  }

  return {
    pendingUnstakeBalance,
    undelegationComplete,
    undelegationPeriod,
    undelegationStatus,
  }
}

export const operatorService = {
  fetchDelegatedTokensData,
  fetchPendingUndelegation,
}
