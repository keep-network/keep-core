import { contractService } from './contracts.service'
import { TOKEN_GRANT_CONTRACT_NAME } from '../constants/constants'

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

export const tokenGrantsService = {
  fetchGrants,
}
