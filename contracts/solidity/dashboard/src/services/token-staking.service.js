import { contractService } from './contracts.service'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
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

const fetchPendingUndelegation = async (web3Context, ...args) => {
  const { yourAddress } = web3Context
  const [withdrawals, stakeWithdrawalDelayInSec, events] = await Promise.all([
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'withdrawals', yourAddress),
    contractService.makeCall(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'stakeWithdrawalDelay'),
    contractService.getPastEvents(web3Context, TOKEN_STAKING_CONTRACT_NAME, 'InitiatedUnstake', { fromBlock: '0', filter: { operator: yourAddress } }),
  ])

  const undelegatedOn = withdrawals.amount === '0' ?
    null :
    moment.unix(withdrawals.createdAt).add(stakeWithdrawalDelayInSec, 'seconds')
  const stakeWithdrawalDelay = moment().add(stakeWithdrawalDelayInSec, 'seconds').fromNow(stakeWithdrawalDelayInSec)
  const pendinUndelegations = events
    .map(({ id, returnValues: { value, createdAt } }) => ({ eventId: id, amount: value, createdAt: moment.unix(createdAt) }))
    .reverse()

  return {
    pendingUnstakeBalance: withdrawals.amount,
    undelegatedOn,
    stakeWithdrawalDelay,
    stakeWithdrawalDelayInSec,
    pendinUndelegations,
  }
}

export const operatorService = {
  fetchDelegatedTokensData,
  fetchPendingUndelegation,
}
