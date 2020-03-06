import { TOKEN_GRANT_CONTRACT_NAME } from '../constants/constants'
import { contractService } from './contracts.service'

const fetchGrants = async (web3Context) => {
  const { yourAddress } = web3Context
  const grantIds = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrants', yourAddress)
  const grants = []

  for (let i = 0; i < grantIds.length; i++) {
    const grantDetails = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrant', grantIds[i])
    const vestingSchedule = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrantVestingSchedule', grantIds[i])

    const vested = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'grantedAmount', grantIds[i])
    let readyToRelease = '0'
    readyToRelease = await contractService
      .makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'withdrawable', grantIds[i])
      .catch(() => readyToRelease = '0')
    const released = grantDetails.withdrawn

    grants.push({ id: grantIds[i], vested, released, readyToRelease, ...vestingSchedule, ...grantDetails })
  }

  return grants
}

export const tokenGrantsService = {
  fetchGrants,
}
