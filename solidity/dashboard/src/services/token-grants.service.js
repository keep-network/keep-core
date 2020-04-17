import { TOKEN_GRANT_CONTRACT_NAME } from '../constants/constants'
import { contractService } from './contracts.service'
import { isSameEthAddress } from '../utils/general.utils'
import web3Utils from 'web3-utils'
import { getGuaranteedMinimumStakingPolicyContractAddress, getPermissiveStakingPolicyContractAddress } from '../contracts'

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
    const availableToStake = await contractService.makeCall(web3Context, TOKEN_GRANT_CONTRACT_NAME, 'availableToStake', grantIds[i])

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
   * from Address of the grant manager.
   * grantee Address of the grantee.
   * cliff Duration in seconds of the cliff after which tokens will begin to unlock.
   * start Timestamp at which unlocking will start.
   * revocable Whether the token grant is revocable or not (1 or 0).
   * stakingPolicyAddress The staking policy as an address
   */
  const stakingPolicyAddress = revocable ?
    getGuaranteedMinimumStakingPolicyContractAddress() :
    getPermissiveStakingPolicyContractAddress()

  const extraData = web3Context.eth.abi.encodeParameters(
    ['address', 'address', 'uint256', 'uint256', 'uint256', 'bool', 'address'],
    [yourAddress, grantee, duration, start, cliff, revocable, stakingPolicyAddress]
  );

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
