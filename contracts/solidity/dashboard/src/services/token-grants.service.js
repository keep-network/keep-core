import { TOKEN_GRANT_CONTRACT_NAME } from '../constants/constants'
import { contractService } from './contracts.service'
import moment from 'moment'
import { displayAmount } from '../utils'

const ONE_WEEK_IN_SEC= 604800

const fetchGrantVestingSchedule = async (web3Context, grantId) => {
  const { eth } = web3Context

  const grantDetails = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrant', grantId)
  const { cliff, duration, start } = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrantVestingSchedule', grantId)

  const cliffBreakpoint = { dotColorClassName: 'grey ring', label: 'Cliff', date: moment.unix(start) }
  const afterCliffBreakpoint = { dotColorClassName: 'grey', label: 'After cliff vested', date: moment.unix(cliff) }
  const breakpoints = [cliffBreakpoint, afterCliffBreakpoint]
  if (moment.unix(cliff).isBefore(moment())) {
    breakpoints.push({ dotColorClassName: 'grey', label: 'current breakpoint', date: moment() })
  }

  const withdrawEvents = await contractService.getPastEvents(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'WithdrawnTokenGrant', { fromBlock: '0', filter: { id: grantId } })
  const stakedEvents = await contractService.getPastEvents(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'StakedGrant', { fromBlock: '0', filter: { id: grantId } })

  for (let i = 0; i < withdrawEvents.length; i++) {
    const { blockNumber, returnValues: { amount } } = withdrawEvents[i]

    const withdrawnAt = (await eth.getBlock(blockNumber)).timestamp
    breakpoints.push({ dotColorClassName: 'primary', label: `Released ${displayAmount(amount)}`, date: moment.unix(withdrawnAt) })
  }

  for (let i = 0; i < stakedEvents.length; i++) {
    const { blockNumber, returnValues: { value } } = stakedEvents[i]

    const withdrawnAt = (await eth.getBlock(blockNumber)).timestamp
    breakpoints.push({ dotColorClassName: 'brown', label: `Staked ${displayAmount(value)}`, date: moment.unix(withdrawnAt) })
  }

  return breakpoints.sort((a, b) => a.date.diff(b.date))
}

export const tokenGrantsService = {
  fetchGrantVestingSchedule,
}
