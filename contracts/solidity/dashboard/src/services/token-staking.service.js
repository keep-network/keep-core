import { contractService } from './contracts.service'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import web3Utils from 'web3-utils'

const fetchDelegatedTokensData = async (web3Context) => {
  const { yourAddress } = web3Context
  const [stakedBalance, ownerAddress, beneficiaryAddress] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'balanceOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'ownerOf', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'magpieOf', yourAddress),
  ])

  return { stakedBalance, ownerAddress, beneficiaryAddress }
}

const fetchPendingUndelegation = async (web3Context) => {
  const { yourAddress } = web3Context
  const [undelegation, undelegationPeriod] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'getUndelegation', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'undelegationPeriod'),
  ])

  const undelegationAtInBN = web3Utils.toBN(undelegation.undelegatedAt).add(web3Utils.toBN(undelegationPeriod))
  const isUndelegation = undelegation.undelegatedAt === '0'
  const pendingUnstakeBalance = isUndelegation ? 0 : undelegation.amount
  const undelegationComplete = isUndelegation ? null : undelegationAtInBN.toString()

  return {
    pendingUnstakeBalance,
    undelegationComplete,
    undelegationPeriod,
  }
}

export const operatorService = {
  fetchDelegatedTokensData,
  fetchPendingUndelegation,
}
