import {
  TOKEN_GRANT_CONTRACT_NAME,
  MANAGED_GRANT_FACTORY_CONTRACT_NAME,
} from "../constants/constants"
import { contractService } from "./contracts.service"
import { isSameEthAddress } from "../utils/general.utils"
import web3Utils from "web3-utils"
import {
  getGuaranteedMinimumStakingPolicyContractAddress,
  getPermissiveStakingPolicyContractAddress,
  createManagedGrantContractInstance,
  CONTRACT_DEPLOY_BLOCK_NUMBER,
} from "../contracts"

const fetchGrants = async (web3Context) => {
  const { yourAddress } = web3Context
  const grantIds = new Set(
    await contractService.makeCall(
      web3Context,
      TOKEN_GRANT_CONTRACT_NAME,
      "getGrants",
      yourAddress
    )
  )
  const managedGrants = await fetchManagedGrants(web3Context)
  const grants = []
  for (const grantId of grantIds) {
    let grantDetails = {}
    try {
      grantDetails = await getGrantDetails(grantId, web3Context)
    } catch {
      continue
    }
    grants.push({ ...grantDetails })
  }

  for (const managedGrant of managedGrants) {
    const { grantId, managedGrantContractInstance } = managedGrant
    const grantDetails = await getGrantDetails(grantId, web3Context, true)
    grants.push({
      ...grantDetails,
      isManagedGrant: true,
      managedGrantContractInstance,
    })
  }
  return grants
}

const getGrantDetails = async (
  grantId,
  web3Context,
  isManagedGrant = false
) => {
  const { yourAddress } = web3Context
  const grantDetails = await contractService.makeCall(
    web3Context,
    TOKEN_GRANT_CONTRACT_NAME,
    "getGrant",
    grantId
  )
  if (!isManagedGrant && !isSameEthAddress(yourAddress, grantDetails.grantee)) {
    throw new Error(
      `${yourAddress} does not match a grantee address for the grantId ${grantId}`
    )
  }
  const unlockingSchedule = await contractService.makeCall(
    web3Context,
    TOKEN_GRANT_CONTRACT_NAME,
    "getGrantUnlockingSchedule",
    grantId
  )

  const unlocked = await contractService.makeCall(
    web3Context,
    TOKEN_GRANT_CONTRACT_NAME,
    "unlockedAmount",
    grantId
  )
  let readyToRelease = "0"
  try {
    readyToRelease = await contractService.makeCall(
      web3Context,
      TOKEN_GRANT_CONTRACT_NAME,
      "withdrawable",
      grantId
    )
  } catch (error) {
    readyToRelease = "0"
  }
  const released = grantDetails.withdrawn
  const availableToStake = await contractService.makeCall(
    web3Context,
    TOKEN_GRANT_CONTRACT_NAME,
    "availableToStake",
    grantId
  )

  return {
    id: grantId,
    unlocked,
    released,
    readyToRelease,
    availableToStake,
    ...unlockingSchedule,
    ...grantDetails,
  }
}

const createGrant = async (web3Context, data, onTransationHashCallback) => {
  const { yourAddress, token, grantContract } = web3Context
  const tokenGrantContractAddress = grantContract.options.address
  const { grantee, amount, duration, start, cliff, revocable } = data

  /**
   * Extra data contains the following values:
   * from Address of the grant manager.
   * grantee Address of the grantee.
   * cliff Duration in seconds of the cliff after which tokens will begin to unlock.
   * start Timestamp at which unlocking will start.
   * revocable Whether the token grant is revocable or not (1 or 0).
   * stakingPolicyAddress The staking policy as an address
   */
  const stakingPolicyAddress = revocable
    ? getGuaranteedMinimumStakingPolicyContractAddress()
    : getPermissiveStakingPolicyContractAddress()

  const extraData = web3Context.eth.abi.encodeParameters(
    ["address", "address", "uint256", "uint256", "uint256", "bool", "address"],
    [
      yourAddress,
      grantee,
      duration,
      start,
      cliff,
      revocable,
      stakingPolicyAddress,
    ]
  )

  const formattedAmount = web3Utils
    .toBN(amount)
    .mul(web3Utils.toBN(10).pow(web3Utils.toBN(18)))
    .toString()

  await token.methods
    .approveAndCall(tokenGrantContractAddress, formattedAmount, extraData)
    .send({ from: yourAddress })
    .on("transactionHash", onTransationHashCallback)
}

const fetchManagedGrants = async (web3Context) => {
  const { managedGrantFactoryContract, yourAddress, web3 } = web3Context

  const managedGrantCreatedEvents = await managedGrantFactoryContract.getPastEvents(
    "ManagedGrantCreated",
    {
      fromBlock:
        CONTRACT_DEPLOY_BLOCK_NUMBER[MANAGED_GRANT_FACTORY_CONTRACT_NAME],
    }
  )
  const grants = []

  for (const event of managedGrantCreatedEvents) {
    const {
      returnValues: { grantAddress },
    } = event
    const managedGrantContractInstance = createManagedGrantContractInstance(
      web3,
      grantAddress
    )
    const grantee = await managedGrantContractInstance.methods.grantee().call()
    if (!isSameEthAddress(yourAddress, grantee)) {
      continue
    }
    const grantId = await managedGrantContractInstance.methods.grantId().call()
    grants.push({ grantId, managedGrantContractInstance })
  }

  return grants
}

export const stake = async (
  web3Context,
  data,
  onTransactionHashCallback = () => {}
) => {
  const { grantContract, stakingContract, yourAddress } = web3Context
  const { amount, delegation, grant } = data
  const { isManagedGrant, managedGrantContractInstance, id } = grant

  if (isManagedGrant) {
    await managedGrantContractInstance.methods
      .stake(stakingContract.options.address, amount, delegation)
      .send({ from: yourAddress })
      .on("transactionHash", onTransactionHashCallback)
  } else {
    await grantContract.methods
      .stake(id, stakingContract.options.address, amount, delegation)
      .send({ from: yourAddress })
      .on("transactionHash", onTransactionHashCallback)
  }
}

const getOperatorsFromManagedGrants = async (web3Context) => {
  const { grantContract } = web3Context
  const manageGrants = await fetchManagedGrants(web3Context)
  const operators = new Set()

  for (const managedGrant of manageGrants) {
    const { managedGrantContractInstance } = managedGrant
    const granteeAddress = managedGrantContractInstance.options.address
    const grenteeOperators = await grantContract.methods
      .getGranteeOperators(granteeAddress)
      .call()
    grenteeOperators.forEach(operators.add, operators)
  }

  return operators
}

export const tokenGrantsService = {
  fetchGrants,
  createGrant,
  fetchManagedGrants,
  stake,
  getOperatorsFromManagedGrants,
}
