import {
  TOKEN_GRANT_CONTRACT_NAME,
  MANAGED_GRANT_FACTORY_CONTRACT_NAME,
} from "../constants/constants"
import { contractService } from "./contracts.service"
import { isSameEthAddress } from "../utils/general.utils"
import { add, gt } from "../utils/arithmetics.utils"
import {
  getGuaranteedMinimumStakingPolicyContractAddress,
  getPermissiveStakingPolicyContractAddress,
  createManagedGrantContractInstance,
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  Web3Loaded,
} from "../contracts"
import BigNumber from "bignumber.js"
import {
  fetchEscrowDepositsByGrantId,
  fetchWithdrawableAmountForDeposit,
  fetchDepositWithdrawnAmount,
} from "./token-staking-escrow.service"

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
  const escrowDepositsEvents = await fetchEscrowDepositsByGrantId(grantId)
  const escrowOperatorsToWithdraw = []
  let escrowWithdrawableAmount = 0
  let escrowWithdrawTotalAmount = 0

  for (const event of escrowDepositsEvents) {
    const {
      returnValues: { operator },
    } = event
    const withdrawable = await fetchWithdrawableAmountForDeposit(operator)
    const withdraw = await fetchDepositWithdrawnAmount(operator)

    escrowWithdrawTotalAmount = add(escrowWithdrawTotalAmount, withdraw)

    if (gt(withdrawable, 0)) {
      escrowOperatorsToWithdraw.push(operator)
      escrowWithdrawableAmount = add(escrowWithdrawableAmount, withdrawable)
    }
  }

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
  const withdrawableAmountGrantOnly = await contractService.makeCall(
    web3Context,
    TOKEN_GRANT_CONTRACT_NAME,
    "withdrawable",
    grantId
  )

  const readyToRelease = add(
    withdrawableAmountGrantOnly,
    escrowWithdrawableAmount
  )
  const released = add(grantDetails.withdrawn, escrowWithdrawTotalAmount)
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
    escrowOperatorsToWithdraw,
    withdrawableAmountGrantOnly,
    ...unlockingSchedule,
    ...grantDetails,
  }
}

const getCreateTokenGrantExtraData = async (data) => {
  const web3Context = await Web3Loaded
  const { grantee, duration, start, cliff, revocable } = data

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
      web3Context.eth.defaultAccount,
      grantee,
      duration,
      start,
      cliff,
      revocable,
      stakingPolicyAddress,
    ]
  )

  return extraData
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

const fetchGrantById = async (web3Context, grantId) => {
  const id = new BigNumber(grantId)

  if (!id.isInteger() || id.isNegative()) {
    throw new Error("Invalid grant ID")
  }

  try {
    return await getGrantDetails(id.toString(), web3Context, true)
  } catch (error) {
    throw new Error("Grant ID not found")
  }
}

export const tokenGrantsService = {
  fetchGrants,
  getCreateTokenGrantExtraData,
  fetchManagedGrants,
  stake,
  getOperatorsFromManagedGrants,
  fetchGrantById,
}
