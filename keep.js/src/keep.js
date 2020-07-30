import KeepToken from "@keep-network/keep-core/artifacts/KeepToken.json"
import TokenStaking from "@keep-network/keep-core/artifacts/TokenStaking.json"
import TokenGrant from "@keep-network/keep-core/artifacts/TokenGrant.json"
import KeepRandomBeaconOperator from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperator.json"
import BondedECDSAKeepFactory from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeepFactory.json"
import TBTCSystem from "@keep-network/tbtc/artifacts/TBTCSystem.json"
import KeepBonding from "@keep-network/keep-ecdsa/artifacts/KeepBonding.json"
import KeepRegistry from "@keep-network/keep-core/artifacts/KeepRegistry.json"
import GuaranteedMinimumStakingPolicy from "@keep-network/keep-core/artifacts/GuaranteedMinimumStakingPolicy.json"
import PermissiveStakingPolicy from "@keep-network/keep-core/artifacts/PermissiveStakingPolicy.json"
import KeepRandomBeaconOperatorStatistics from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperatorStatistics.json"
import ManagedGrant from "@keep-network/keep-core/artifacts/ManagedGrant.json"
import ManagedGrantFactory from "@keep-network/keep-core/artifacts/ManagedGrantFactory.json"
import TBTCToken from "@keep-network/tbtc/artifacts/TBTCToken.json"
import Deposit from "@keep-network/tbtc/artifacts/Deposit.json"
import BondedECDSAKeep from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeep.json"
import ContractFactory from "./contract-wrapper.js"
import { TokenStakingConstants } from "./constants.js"
import {
  add,
  gt,
  lte,
  isSameEthAddress,
  lookupArtifactAddress,
} from "./utils.js"
/** @typedef {import("./contract-wrapper").ContractWrapper} ContractWrapper */

export const contracts = new Map([
  [KeepToken, "keepTokenContract"],
  [TokenStaking, "tokenStakingContract"],
  [TokenGrant, "tokenGrantContract"],
  [KeepRandomBeaconOperator, "keepRandomBeaconOperatorContract"],
  [
    KeepRandomBeaconOperatorStatistics,
    "keepRandomBeaconOperatorStatisticsContract",
  ],
  [KeepRegistry, "keepRegirstyContract"],
  [BondedECDSAKeepFactory, "bondedECDSAKeepFactoryContract"],
  [KeepBonding, "keepBondingContract"],
  [TBTCSystem, "tbtcSystemContract"],
  [TBTCToken, "tbtcTokenContract"],
  [BondedECDSAKeep, "bondedECDSAKeepContract"],
  [ManagedGrantFactory, "managedGrantFactoryContract"],
])

export default class KEEP {
  static async initialize(config) {
    const keep = new KEEP(config)
    await keep.initializeContracts()

    return keep
  }

  constructor(config) {
    this.config = config
  }

  async initializeContracts() {
    for (const [artifact, propertyName] of contracts) {
      this[propertyName] = await ContractFactory.createContractInstance(
        artifact,
        this.config
      )
    }

    /**
     * @type TokenStakingConstants
     */
    this.tokenStakingConstants = await TokenStakingConstants.initialize(
      this.tokenStakingContract
    )

    /**
     * @type ContractWrapper
     */
    this.keepTokenContract

    /**
     * @type ContractWrapper
     */
    this.tokenStakingContract

    /**
     * @type ContractWrapper
     */
    this.tokenGrantContract

    /**
     * @type ContractWrapper
     */
    this.keepRandomBeaconOperatorContract

    /**
     * @type ContractWrapper
     */
    this.keepRandomBeaconOperatorStatisticsContract

    /**
     * @type ContractWrapper
     */
    this.keepRegirstyContract

    /**
     * @type ContractWrapper
     */
    this.bondedECDSAKeepFactoryContract

    /**
     * @type ContractWrapper
     */
    this.keepBondingContract

    /**
     * @type ContractWrapper
     */
    this.tbtcSystemContract

    /**
     * @type ContractWrapper
     */
    this.managedGrantFactoryContract

    /**
     * @type ContractWrapper
     */
    this.tbtcTokenContract

    /**
     * @type ContractWrapper
     */
    this.bondedECDSAKeepContract

    this.tbtcSortitionPoolAddress = ""
  }

  /**
   * Returns the authorizer for the given operator address.
   *
   * @param {string} operatorAddress
   * @return {Promise<string>} Authrorizer address.
   */
  async authorizerOf(operatorAddress) {
    return await this.tokenStakingContract.makeCall(
      "authorizerOf",
      operatorAddress
    )
  }

  /**
   * Returns the beneficiary for the given operator address.
   *
   * @param {string} operatorAddress
   * @return {Promise<string>} Beneficiary address.
   */
  async beneficiaryOf(operatorAddress) {
    return await this.tokenStakingContract.makeCall(
      "beneficiaryOf",
      operatorAddress
    )
  }

  /**
   * Returns the stake owner for the specified operator address.
   *
   * @param {string} operatorAddress
   * @return {Promise<string>} Stake owner address.
   */
  async ownerOf(operatorAddress) {
    return await this.tokenStakingContract.makeCall("ownerOf", operatorAddress)
  }

  /**
   * Returns the list of operators of the given owner address.
   *
   * @param {string} ownerAddress
   * @return {Promise<string[]>} An array of addresses.
   */
  async operatorsOf(ownerAddress) {
    return await this.tokenStakingContract.makeCall("operatorsOf", ownerAddress)
  }

  /**
   * Returns the list of operators of the provided beneficiary address.
   *
   * @param {string} beneficiary Beneficiary address.
   * @return {Primise<string[]>} An array of addresses.
   */
  async operatorsOfBeneficiary(beneficiary) {
    return (
      await this.tokenStakingContract.getPastEvents("Staked", {
        beneficiary,
      })
    ).map((_) => _.returnValues.operator)
  }

  /**
   * Returns the list of operators of the provided authorizer address.
   *
   * @param {string} authorizer Authorizer address.
   * @return {Primise<string[]>} An array of addresses.
   */
  async operatorsOfAuthorizer(authorizer) {
    return (
      await this.tokenStakingContract.getPastEvents("Staked", {
        authorizer,
      })
    ).map((_) => _.returnValues.operator)
  }

  /**
   * @typedef {Object} DelegationInfo
   * @property {string} amount The amount of tokens the given operator delegated.
   * @property {string} createdAt The time when the stake has been delegated.
   * @property {string} undelegatedAt The time when undelegation has been requested.
   */
  /**
   * Returns stake delegation info for the given operator address.
   *
   * @param {string} operatorAddress
   * @return {Promise<DelegationInfo>} Stake delegation info.
   */
  async getDelegationInfo(operatorAddress) {
    return await this.tokenStakingContract.makeCall(
      "getDelegationInfo",
      operatorAddress
    )
  }

  /**
   * @typedef {Object} DelegationAddresses
   * @property {string} authroizer
   * @property {string} beneficiary
   * @property {string} operator
   *
   * @typedef {DelegationInfo | DelegationAddresses} FullDelegationInfo
   */
  /**
   * Returns delegations for given operators.
   * @param {string[]} operatorAddresses An array of operator addresses.
   * @return {Promise<FullDelegationInfo[]>} Array of delegations
   */
  async getDelegations(operatorAddresses) {
    const delegations = []
    for (const operator of operatorAddresses) {
      const delegationInfo = await this.getDelegationInfo(operator)
      const beneficiary = await this.beneficiaryOf(operator)
      const authorizer = await this.authorizerOf(operator)
      delegations.push({ ...delegationInfo, beneficiary, authorizer, operator })
    }

    return delegations
  }

  /**
   * Authorizes operator contract to access staked token balance of the provided operator.
   * Can only be executed by stake operator authorizer.
   *
   * @param {string} operatorAddress
   * @return {*}
   */
  authorizeKeepRandomBeaconOperatorContract(operatorAddress) {
    const keepRandomBeaconOperatorContractAddress = this
      .keepRandomBeaconOperatorContract.address
    return this.tokenStakingContract.sendTransaction(
      "authorizeOperatorContract",
      operatorAddress,
      keepRandomBeaconOperatorContractAddress
    )
  }

  /**
   * Checks if operator contract has access to the staked token balance of the provided operator.
   *
   * @param {string} operatorAddress
   * @return {Promise<boolean>}
   */
  async isAuthorizedForKeepRandomBeacon(operatorAddress) {
    return await this.tokenStakingContract.makeCall(
      "isAuthorizedForOperator",
      operatorAddress,
      this.keepRandomBeaconOperatorContract.address
    )
  }

  /**
   * @typedef {Object} GroupMemberRewardsWithdrawnEventValues
   * @property {Object} returnValues
   * @property {string} returnValues.beneficiary
   * @property {string} returnValues.operator
   * @property {string} returnValues.amount
   * @property {string} returnValues.groupIndex
   *
   * @typedef {import("./contract-wrapper").EventData & GroupMemberRewardsWithdrawnEventValues} GroupMemberRewardsWithdrawnEvent
   */

  /**
   * Returns withdrawn rewards for a given beneficiary address.
   *
   * @param {string} beneficiaryAddress
   *
   * @return {Promise<Array<GroupMemberRewardsWithdrawnEvent>>} Withdrawal Events
   */
  async getWithdrawnRewardsForBeneficiary(beneficiaryAddress) {
    return await this.keepRandomBeaconOperatorContract.getPastEvents(
      "GroupMemberRewardsWithdrawn",
      { beneficiary: beneficiaryAddress }
    )
  }

  /**
   *  Withdraws accumulated group member rewards for operator using the provided group index.
   *
   * @param {string} memberAddress
   * @param {string | number} groupIndex
   *
   * @return {*}
   */
  withdrawGroupMemberRewards(memberAddress, groupIndex) {
    return this.keepRandomBeaconOperatorContract.sendTransaction(
      "withdrawGroupMemberRewards",
      memberAddress,
      groupIndex
    )
  }

  /**
   * @typedef {Object} DkgResultSubmittedEventValues
   * @property {Object} returnValues
   * @property {string} returnValues.memberIndex
   * @property {string} returnValues.groupPubKey
   * @property {*} returnValues.misbehaved
   *
   * @typedef {import("./contract-wrapper").EventData & DkgResultSubmittedEventValues} DkgResultSubmittedEvent
   */
  /**
   *
   * @return {Promise<Array<DkgResultSubmittedEvent>>}
   */
  async getAllCreatedGroups() {
    return await this.keepRandomBeaconOperatorContract.getPastEvents(
      "DkgResultSubmittedEvent"
    )
  }

  /**
   * Returns available rewards for a provided beneficiary address
   * @param {*} beneficiaryAddress
   *
   * @typedef {Object} Reward
   * @property {string} groupIndex
   * @property {string} groupPublicKey
   * @property {boolean} isStale
   * @property {boolean} isTerminated
   * @property {string} operatorAddress
   * @property {string} reward
   *
   * @return {Promise<Array<Reward>>} Available rewards
   */
  async findKeepRandomBeaconRewardsForBeneficiary(beneficiaryAddress) {
    const groupPublicKeys = (await this.getAllCreatedGroups()).map(
      (event) => event.returnValues.groupPubKey
    )
    const beneficiaryOperators = await this.operatorsOfBeneficiary(
      beneficiaryAddress
    )

    const groupsInfo = {}
    const rewards = []

    for (
      let groupIndex = 0;
      groupIndex < groupPublicKeys.length;
      groupIndex++
    ) {
      const groupPublicKey = groupPublicKeys[groupIndex]
      for (const memberAddress of beneficiaryOperators) {
        const awaitingRewards = await this.keepRandomBeaconOperatorStatisticsContract.makeCall(
          "awaitingRewards",
          memberAddress,
          groupIndex
        )

        if (!gt(awaitingRewards, 0)) {
          continue
        }

        let groupInfo = {}
        if (groupsInfo.hasOwnProperty(groupIndex)) {
          groupInfo = { ...groupsInfo[groupIndex] }
        } else {
          const isStale = await this.keepRandomBeaconOperatorContract.makeCall(
            "isStaleGroup",
            groupPublicKey
          )

          const isTerminated =
            !isStale &&
            (await this.keepRandomBeaconOperatorContract.makeCall(
              "isGroupTerminated",
              groupIndex
            ))

          groupInfo = {
            groupPublicKey,
            isStale,
            isTerminated,
          }

          groupsInfo[groupIndex] = groupInfo
        }

        rewards.push({
          groupIndex: groupIndex.toString(),
          ...groupInfo,
          operatorAddress: memberAddress,
          reward: awaitingRewards,
        })
      }
    }

    return rewards
  }

  /**
   * Returns slashed tokens for a provided operator address
   *
   * @param {string} operatorAddress
   * @return {Promise<any[]>} Slashed tokens data.
   */
  async getSlashedTokens(operatorAddress) {
    const data = []

    const slashedTokensEvents = await this.tokenStakingContract.getPastEvents(
      "TokensSlashed",
      { operator: operatorAddress }
    )
    const seizedTokensEvents = await this.tokenStakingContract.getPastEvents(
      "TokensSeized",
      { operator: operatorAddress }
    )

    if (slashedTokensEvents.length === 0 && seizedTokensEvents.length === 0) {
      return data
    }

    const unauthorizedSigningEvents = await this.keepRandomBeaconOperatorContract.getPastEvents(
      "UnauthorizedSigningReported"
    )

    const relayEntryTimeoutEvents = await this.keepRandomBeaconOperatorContract.getPastEvents(
      "RelayEntryTimeoutReported"
    )

    const punishmentEvents = [
      ...unauthorizedSigningEvents,
      ...relayEntryTimeoutEvents,
    ]

    const groupByTransactionHash = (events) => {
      const groupedByTransactionHash = {}

      events.forEach((event) => {
        const { transactionHash, returnValues } = event
        if (groupedByTransactionHash.hasOwnProperty(transactionHash)) {
          const prevData = groupedByTransactionHash[transactionHash]
          groupedByTransactionHash[transactionHash] = {
            ...returnValues,
            amount: add(returnValues.amount, prevData.amount),
          }
        } else {
          groupedByTransactionHash[transactionHash] = { ...returnValues }
        }
      })

      return groupedByTransactionHash
    }

    const slashedTokensGroupedByTxtHash = groupByTransactionHash(
      slashedTokensEvents
    )
    const seizedTokensGroupedByTxtHash = groupByTransactionHash(
      seizedTokensEvents
    )

    for (let i = 0; i < punishmentEvents.length; i++) {
      const {
        transactionHash,
        returnValues: { groupIndex },
      } = punishmentEvents[i]
      let punishmentData = {}
      if (slashedTokensGroupedByTxtHash.hasOwnProperty(transactionHash)) {
        const { amount } = slashedTokensGroupedByTxtHash[transactionHash]
        punishmentData = {
          amount,
          groupIndex,
          ...punishmentEvents[i],
        }
      } else if (seizedTokensGroupedByTxtHash.hasOwnProperty(transactionHash)) {
        const { amount } = seizedTokensGroupedByTxtHash[transactionHash]
        punishmentData = {
          amount,
          groupIndex,
          ...punishmentEvents[i],
        }
      }

      if (lte(punishmentData.amount, 0)) continue

      data.push(punishmentData)
    }

    return data
  }

  /**
   * @typedef {Object} TbtcReward
   * @property {string} depositTokenId,
   * @property {string} amount
   * @property {string} transactionHash
   */
  /**
   * Returns withdrawn tbtc rewards for a provided beneficiary address
   *
   * @param {string} beneficiaryAddress
   * @return {Promise<Array<TbtcReward>>} An array of tbtc rewards.
   */
  async getWithdrawnTBTCRewards(beneficiaryAddress) {
    const transefEventsToBeneficiary = await this.tbtcTokenContract.getPastEvents(
      "Transfer",
      { to: beneficiaryAddress }
    )
    const fromAddresses = transefEventsToBeneficiary.map(
      (_) => _.returnValues.from
    )

    const depositCreatedEvents = await this.tbtcSystemContract.getPastEvents(
      "Created",
      {
        _depositContractAddress: fromAddresses,
      }
    )

    const tbtcRewards = transefEventsToBeneficiary
      .filter(({ returnValues: { from } }) =>
        depositCreatedEvents.some(
          ({ returnValues: { _depositContractAddress } }) =>
            isSameEthAddress(_depositContractAddress, from)
        )
      )
      .map(({ returnValues: { from, value }, ...eventData }) => ({
        depositTokenId: from,
        amount: value,
        ...eventData,
      }))

    return tbtcRewards
  }

  /**
   * Returns the KEEP operator addresses associated with the Deposit.
   *
   * @param {string} depositAddress
   * @return {Promise<string[]>}
   */
  async getOperatorsForDeposit(depositAddress) {
    const depositContract = ContractFactory.new(
      new this.config.web3.eth.Contract(Deposit.abi, depositAddress)
    )
    const keepAddress = await depositContract.makeCall("getKeepAddress")

    const bondedECDSAKeepContract = ContractFactory.new(
      new this.config.web3.eth.Contract(BondedECDSAKeep.abi, keepAddress)
    )

    return await bondedECDSAKeepContract.makeCall("getMembers")
  }

  /**
   * @typedef {Object} CreateGrantData
   * @property {string} grantManager Address of the grant manager.
   * @property {string} grantee Address of the grantee.
   * @property {string | number} start Timestamp at which unlocking will start.
   * @property {string | number} duration Duration in seconds of the unlocking period.
   * @property {string | number} cliffDuration in seconds of the cliff; no tokens will be unlocked until the time `start + cliff`.
   * @property {boolean} revocable Whether the token grant is revocable or not (1 or 0).
   * @property {string} stakingPolicy Address of the staking policy for the grant.
   * @property {string} amount Approved amount in wei for the transfer to create token grant.
   */
  /**
   * Creates a token grant with a unlocking schedule where balance withdrawn to
   * the grantee gradually in a linear fashion until start + duration.
   * By then all of the balance will have unlocked.
   *
   * @param {CreateGrantData} data
   */
  async createGrant(data) {
    const {
      grantManager,
      grantee,
      duration,
      start,
      cliffDuration,
      revocable,
      stakingPolicyAddress,
      amount,
    } = data
    const extraData = this.config.web3.eth.abi.encodeParameters(
      [
        "address",
        "address",
        "uint256",
        "uint256",
        "uint256",
        "bool",
        "address",
      ],
      [
        grantManager,
        grantee,
        duration,
        start,
        cliffDuration,
        revocable,
        stakingPolicyAddress,
      ]
    )

    return this.keepTokenContract.sendTransaction(
      "approveAndCall",
      this.tokenGrantContract.address,
      amount,
      extraData
    )
  }

  get guaranteedMinimumStakingPolicyAddress() {
    return lookupArtifactAddress(
      GuaranteedMinimumStakingPolicy,
      this.config.networkId
    )
  }

  get permissiveStakingPolicyAddress() {
    return lookupArtifactAddress(PermissiveStakingPolicy, this.config.networkId)
  }

  /**
   * Returns a managed grants addresses for a provided grantee address.
   *
   * @param {string} grantee Address of the grantee.
   * @return {Promis<string[]>}
   */
  async getGranteeManagedGrantAddresses(grantee) {
    const managedGrantCreatedEvents = await this.managedGrantFactoryContract.getPastEvents(
      "ManagedGrantCreated"
    )
    const managedGrantAddresses = []

    for (const event of managedGrantCreatedEvents) {
      const {
        returnValues: { grantAddress },
      } = event

      const managedGrantContractInstance = ContractFactory.new(
        new this.config.web3.eth.Contract(ManagedGrant.abi, grantAddress)
      )
      const managedGrantGrantee = await managedGrantContractInstance.makeCall(
        "grantee"
      )

      if (!isSameEthAddress(grantee, managedGrantGrantee)) {
        continue
      }

      managedGrantAddresses.push(grantAddress)
    }

    return managedGrantAddresses
  }

  /**
   *
   * @typedef {Object} GrantDetails
   * @property {string} id Grant ID.
   * @property {string} unlcoked Unlocked grant amount
   * @property {string} released The amount of tokens that have already been withdrawn from the grant.
   * @property {string} readyToRelease Withdrawable granted amount.
   * @property {string} availableToStake The amount of tokens available for staking from the grant
   * @property {string} amount The amount of tokens the grant provides.
   * @property {string} staked The amount of tokens that have been staked from the grant.
   * @property {string} revokedAmount The number of tokens revoked from the grantee.
   * @property {string} revokedAt Timestamp at which grant was revoked by the grant manager.
   * @property {string} grantee The grantee of grant.
   * @property {string} grantManager The address designated as the manager of the grant, which is the only address that can revoke this grant.
   * @property {string | number} duration The duration, in seconds, during which the tokens will unlocking linearly.
   * @property {string | number} start The start time, as a timestamp comparing to `now`.
   * @property {string | number} cliff The timestamp, before which none of the tokens
   * in the grant will be unlocked, and after which a linear amount based on
   * the time elapsed since the start will be unlocked.
   * @property {string} policy The address of the grant's staking policy.
   */
  /**
   * Returns a details of the provided grant id.
   *
   * @param {string} grantId
   * @return {GrantDetails} Grant details
   */
  async getGrantDetails(grantId) {
    const grantDetails = await this.tokenGrantContract.makeCall(
      "getGrant",
      grantId
    )

    const unlockingSchedule = await this.tokenGrantContract.makeCall(
      "getGrantUnlockingSchedule",
      grantId
    )

    const unlocked = await this.tokenGrantContract.makeCall(
      "unlockedAmount",
      grantId
    )

    let readyToRelease = "0"
    try {
      readyToRelease = await this.tokenGrantContract.makeCall(
        "withdrawable",
        grantId
      )
    } catch (error) {
      readyToRelease = "0"
    }

    const released = grantDetails.withdrawn
    const availableToStake = await this.tokenGrantContract.makeCall(
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

  /**
   * Authorizes BondedECDSAKeepFactory contract to access staked token balance of the provided operator.
   *
   * @param {string} operatorAddress Address of the operator.
   * @return {*}
   */
  authorizeBondedECDSAKeepFactory(operatorAddress) {
    return this.tokenStakingContract.sendTransaction(
      "authorizeOperatorContract",
      operatorAddress,
      this.bondedECDSAKeepFactoryContract.address
    )
  }

  /**
   * Authorizes TBTC sortition pool for the provided operator.
   *
   * @param {string} operatorAddress Address of the operator.
   * @return {*}
   */
  async authorizeTBTCSystem(operatorAddress) {
    const tbtcSortitionPoolAddress = await this.getTBTCSortitionPoolAddress()

    return this.keepBondingContract.sendTransaction(
      "authorizeSortitionPoolContract",
      operatorAddress,
      tbtcSortitionPoolAddress
    )
  }

  /**
   * Returns the sortition pool address of the TBTC application.
   *
   * @return {Promise<string>} The TBTC sortition pool address
   */
  async getTBTCSortitionPoolAddress() {
    if (this.tbtcSortitionPoolAddress) {
      return this.tbtcSortitionPoolAddress
    }

    this.tbtcSortitionPoolAddress = await this.bondedECDSAKeepFactoryContract.makeCall(
      "getSortitionPool",
      this.tbtcSystemContract.address
    )

    return this.tbtcSortitionPoolAddress
  }

  /**
   * Checks if the TBTC sortition pool has been authorized for the provided operator by its authorizer.
   *
   * @param {string} operatorAddress Address of the operator.
   * @return {Promise<boolean>}
   */
  async isTBTCSystemAuthorized(operatorAddress) {
    try {
      const tbtcSortitionPoolAddress = await this.getTBTCSortitionPoolAddress()

      return await this.keepBondingContract.makeCall(
        "hasSecondaryAuthorization",
        operatorAddress,
        tbtcSortitionPoolAddress
      )
    } catch (error) {
      return false
    }
  }

  /**
   * Add the provided value to operator's pool available for bonding.
   *
   * @param {string} operatorAddress Address of the operator.
   * @return {*}
   */
  depositEthForBondingOperator(operatorAddress) {
    return this.keepBondingContract.sendTransaction("deposit", operatorAddress)
  }

  /**
   * Withdraws amount from operator's value available for bonding.
   *
   * @param {string} value Value to withdraw in wei.
   * @param {string} operatorAddress Address of the operator.
   * @return {*}
   */
  withdrawUnbondedEthForOperator(value, operatorAddress) {
    return this.keepBondingContract.sendTransaction(
      "withdraw",
      value,
      operatorAddress
    )
  }

  /**
   * Withdraws amount from operator's value available for bonding as a grantee in managed grant.
   *
   * @param {string} value Value to withdraw in wei.
   * @param {string} operatorAddress Address of the operator.
   * @param {string} managedGrantAddress Address of the managed grant.
   * @return {*}
   */
  withdrawUnbondedEthForOperatorAsManagedGrantee(
    value,
    operatorAddress,
    managedGrantAddress
  ) {
    return this.keepBondingContract.sendTransaction(
      "withdrawAsManagedGrantee",
      value,
      operatorAddress,
      managedGrantAddress
    )
  }

  /**
   * Deauthorizes sortition pool for the provided operator.
   *
   * @param {string} operatorAddress Address of the operator.
   * @return {*}
   */
  async deauthorizeTBTCSystem(operatorAddress) {
    const tbtcSortitionPoolAddress = await this.getTBTCSortitionPoolAddress()

    return this.keepBondingContract.sendTransaction(
      "deauthorizeSortitionPoolContract",
      operatorAddress,
      tbtcSortitionPoolAddress
    )
  }

  /**
   * Returns unassigned value in wei deposited by operators.
   *
   * @param {string} operatorAddress Address of the operator.
   * @return {Promise<string>} Unbonded value in wei.
   */
  async getOperatorUnbondedValue(operatorAddress) {
    return await this.keepBondingContract.makeCall(
      "unbondedValue",
      operatorAddress
    )
  }

  /**
   * Returns the amount of wei the operator has made available for
   * bonding and that is still unbounded. If the operator doesn't exist or
   * bond creator is not authorized as an operator contract or it is not
   * authorized by the operator or there is no secondary authorization for
   * the provided sortition pool, function returns 0.
   *
   * @param {string} operatorAddress Address of the operator.
   * @param {string} holder Address of the holder of the bond.
   * @param {string} referenceID Reference ID used to track the bond by holder.
   * @return {Promise<string>} Amount of authorized wei deposit available for bonding.
   */
  async getOperatorBondAmount(operatorAddress, holder, referenceID) {
    return await this.keepBondingContract.makeCall(
      "bondAmount",
      operatorAddress,
      holder,
      referenceID
    )
  }

  /**
   * @typedef {Object} CreatedBond
   * @property {string} operator Address of the operator to bond.
   * @property {string} holder Address of the holder of the bond.
   * @property {string} referenceID Reference ID used to track the bond by holder.
   */
  /**
   * Returns created bonds for the TBTC sortition pool and provided operator addresses.
   *
   * @param {string[]} operatorAddresses Address of the operator.
   * @return {Promise<Array<CreatedBond>>} Created bonds.
   */
  async getCreatedBondsForTBTC(operatorAddresses) {
    const tbtcSortitionPoolAddress = await this.getTBTCSortitionPoolAddress()

    return (
      await this.keepBondingContract.getPastEvents("BondCreated", {
        operator: operatorAddresses,
        sortitionPool: tbtcSortitionPoolAddress,
      })
    ).map((_) => {
      return {
        operator: _.returnValues.operator,
        holder: _.returnValues.holder,
        referenceID: _.returnValues.referenceID,
      }
    })
  }
}
