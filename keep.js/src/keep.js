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
import ContractFactory from "./contract-wrapper.js"
import { TokenStakingConstants } from "./constants.js"
import { isSameEthAddress, gt, lte } from "./utils.js"

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
}
