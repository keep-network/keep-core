import { TOKEN_GRANT_CONTRACT_NAME } from '../constants/constants'
import { contractService } from './contracts.service'
import { isSameEthAddress } from '../utils/general.utils'
import web3Utils from 'web3-utils'
import { sub } from '../utils/arithmetics.utils'

const fetchGrants = async (web3Context) => {
  const { yourAddress } = web3Context
  const grantIds = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrants', yourAddress)
  const grants = []

  for (let i = 0; i < grantIds.length; i++) {
    const grantDetails = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrant', grantIds[i])
    if (!isSameEthAddress(yourAddress, grantDetails.grantee)) {
      continue
    }
    const unlockingSchedule = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'getGrantUnlockingSchedule', grantIds[i])

    const unlocked = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'unlockedAmount', grantIds[i])
    let readyToRelease = '0'
    try {
      readyToRelease = await contractService
        .makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'withdrawable', grantIds[i])
    } catch (error) {
      readyToRelease = '0'
    }
    const released = grantDetails.withdrawn
    const availableToStake = sub(sub(grantDetails.amount, grantDetails.withdrawn), grantDetails.staked)

    grants.push({ id: grantIds[i], unlocked, released, readyToRelease, availableToStake, ...unlockingSchedule, ...grantDetails })
  }

  return grants
}

const createGrant = async (web3Context, data, onTransationHashCallback) => {
  const { yourAddress, token, grantContract } = web3Context
  const tokenGrantContractAddress = grantContract.options.address
  const {
    grantee,
    amount,
    duration,
    start,
    cliff,
    revocable,
  } = data

  /**
   * Extra data contains the following values:
   * grantee (20 bytes) Address of the grantee.
   * cliff (32 bytes) Duration in seconds of the cliff after which tokens will begin to unlock.
   * start (32 bytes) Timestamp at which unlocking will start.
   * revocable (1 byte) Whether the token grant is revocable or not (1 or 0).
   */
  const extraData = Buffer.concat([
    Buffer.from(grantee.substr(2), 'hex'),
    web3Utils.toBN(duration).toBuffer('be', 32),
    web3Utils.toBN(start).toBuffer('be', 32),
    web3Utils.toBN(cliff).toBuffer('be', 32),
    Buffer.from(revocable ? '01' : '00', 'hex'),
  ])

  const formattedAmount = web3Utils.toBN(amount).mul(web3Utils.toBN(10).pow(web3Utils.toBN(18))).toString()

  await token.methods
    .approveAndCall(
      tokenGrantContractAddress,
      formattedAmount,
      extraData
    )
    .send({ from: yourAddress })
    .on('transactionHash', onTransationHashCallback)
}

export const tokenGrantsService = {
  fetchGrants,
  createGrant,
}
