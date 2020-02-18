import { TOKEN_GRANT_CONTRACT_NAME } from '../constants/constants'
import { contractService } from './contracts.service'
import moment from 'moment'
import { displayAmount } from '../utils'
import web3Utils from 'web3-utils'

const fetchGrants = async (web3Context) => {
  const { yourAddress } = web3Context
  const grantIds = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrants', yourAddress)
  const grants = []

  for (let i = 0; i < grantIds.length; i++) {
    const grantDetails = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrant', grantIds[i])
    const vestingSchedule = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrantVestingSchedule', grantIds[i])

    const vested = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'grantedAmount', grantIds[i])
    const released = grantDetails.withdrawn

    grants.push({ id: grantIds[i], vested, released, ...vestingSchedule, ...grantDetails })
  }

  return grants
}

const fetchGrantVestingSchedule = async (web3Context, grantId) => {
  const { eth } = web3Context
  const { cliff, start, duration } = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrantVestingSchedule', grantId)
  const { amount } = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrant', grantId)

  const vested = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'grantedAmount', grantId)
  const cliffDuration = web3Utils.toBN(cliff).sub(web3Utils.toBN(start))
  const vestedAmountAfterCliff = web3Utils.toBN(amount).mul(cliffDuration).div(web3Utils.toBN(duration)).div(web3Utils.toBN(12))

  const cliffBreakpoint = {
    dotColorClassName: 'grey ring',
    label: 'Cliff Start',
    date: moment.unix(start),
  }
  const afterCliffBreakpoint = {
    dotColorClassName: 'grey ring',
    label: `Cliff End ${displayAmount(vestedAmountAfterCliff, 18, 2)} KEEP Will Be Vested`,
    date: moment.unix(cliff),
  }
  const breakpoints = [cliffBreakpoint, afterCliffBreakpoint]
  if (moment.unix(cliff).isBefore(moment())) {
    breakpoints.push({
      dotColorClassName: 'grey',
      label: `${displayAmount(vested, 18, 2)} KEEP Vested`,
      date: moment(),
    })
  }

  const filterObj = { fromBlock: '0', filter: { id: grantId } }
  const withdrawEvents = await contractService.getPastEvents(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'WithdrawnTokenGrant', filterObj)
  const stakedEvents = await contractService.getPastEvents(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'StakedGrant', filterObj)

  for (let i = 0; i < withdrawEvents.length; i++) {
    const { blockNumber, returnValues: { amount } } = withdrawEvents[i]

    const withdrawnAt = (await eth.getBlock(blockNumber)).timestamp
    breakpoints.push({
      dotColorClassName: 'primary',
      label: `${displayAmount(amount)} KEEP Released `,
      date: moment.unix(withdrawnAt),
    })
  }

  for (let i = 0; i < stakedEvents.length; i++) {
    const { blockNumber, returnValues: { value } } = stakedEvents[i]

    const withdrawnAt = (await eth.getBlock(blockNumber)).timestamp
    breakpoints.push({
      dotColorClassName: 'brown',
      label: `${displayAmount(value)} KEEP Staked`,
      date: moment.unix(withdrawnAt),
    })
  }

  return breakpoints.sort((a, b) => a.date.diff(b.date))
}

export const tokenGrantsService = {
  fetchGrants,
  fetchGrantVestingSchedule,
}
