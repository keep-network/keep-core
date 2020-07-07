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
// import ManagedGrant from "@keep-network/keep-core/artifacts/ManagedGrant.json"
import ManagedGrantFactory from "@keep-network/keep-core/artifacts/ManagedGrantFactory.json"
import TBTCToken from "@keep-network/tbtc/artifacts/TBTCToken.json"
import Deposit from "@keep-network/tbtc/artifacts/Deposit.json"
import BondedECDSAKeep from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeep.json"
import ContractFactory, { ContractWrapper } from "./contract-wrapper.js"
import { TokenStakingConstants } from "./constants.js"
import { gt, lte } from "./utils.js"

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
  [Deposit, "depositContract"],
  [BondedECDSAKeep, "bondedECDSAKeepContract"],
  [GuaranteedMinimumStakingPolicy, "guaranteedMinimumStakingPolicyContract"],
  [PermissiveStakingPolicy, "permissiveStakingPolicyContract"],
  // TODO create managed grant instance for a given address
  // [ManagedGrant, "managedGrantContract"],
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

    this.tokenStakingConstants = await TokenStakingConstants.initialize(
      this.tokenStakingContract
    )

    this.keepTokenContract
    this.tokenStakingContract
    this.tokenGrantContract
    this.keepRandomBeaconOperatorContract
    this.keepRandomBeaconOperatorStatisticsContract
    this.keepRegirstyContract
    this.bondedECDSAKeepFactoryContract
    this.keepBondingContract
    this.tbtcSystemContract
    this.guaranteedMinimumStakingPolicyContract
    // this.managedGrantContract
    this.managedGrantFactoryContract
    this.tbtcTokenContract
    this.depositContract
    this.bondedECDSAKeepContract
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
  async oepratorsOf(ownerAddress) {
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

          groups[groupIndex] = groupInfo
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
}
