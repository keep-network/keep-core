import { contractService } from './contracts.service'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import web3Utils from 'web3-utils'
import moment from 'moment'

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
  const [undelegation, undelegationPeriod, events] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'getUndelegation', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'undelegationPeriod'),
    contractService.getPastEvents(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'Undelegated', { fromBlock: '0', filter: { operator: yourAddress } }),
  ])

  const undelegationAtInBN = web3Utils.toBN(undelegation.undelegatedAt).add(web3Utils.toBN(undelegationPeriod))
  const pendingUnstakeBalance = undelegation.undelegatedAt === '0' ? 0 : undelegation.balance
  const undelegationComplete = undelegation.undelegatedAt === '0' ? null : undelegationAtInBN.toString()

  const pendinUndelegations = events
    .map(({ id, returnValues: { value, createdAt } }) => ({ eventId: id, amount: value, createdAt: moment.unix(createdAt) }))
    .reverse()

  return {
    pendingUnstakeBalance,
    undelegationComplete,
    undelegationPeriod,
    pendinUndelegations,
  }
}

export const operatorService = {
  fetchDelegatedTokensData,
  fetchPendingUndelegation,
}
