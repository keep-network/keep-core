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
import ContractFactory from "./contract-wrapper"
import { TokenStakingConstants } from "./constants"

export default class KEEP {
  static async initialize(config) {
    const keep = new KEEP(config)
    await keep.initializeContracts()

    return keep
  }

  constructor(config) {
    this.config = config
  }

  keepTokenContract
  tokenStakingContract
  tokenGrantContract
  keepRandomBeaconOperatorContract
  keepRandomBeaconOperatorStatisticsContract
  keepRegirstyContract
  bondedECDSAKeepFactoryContract
  keepBondingContract
  tbtcSystemContract
  guaranteedMinimumStakingPolicyContract
  managedGrantContract
  managedGrantFactoryContract
  tbtcTokenContract
  depositContract
  bondedECDSAKeepContract

  async initializeContracts() {
    const contracts = new Map([
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
      [BondedECDSAKeep, bondedECDSAKeepContract],
      [
        GuaranteedMinimumStakingPolicy,
        "guaranteedMinimumStakingPolicyContract",
      ],
      [PermissiveStakingPolicy, "permissiveStakingPolicyContract"],
      [ManagedGrant, "managedGrantContract"],
      [ManagedGrantFactory, "managedGrantFactoryContract"],
    ])

    for (const [artifact, propertyName] of contracts) {
      this[propertyName] = await ContractFactory.createContractInstance(
        artifact,
        this.config
      )
    }

    this.tokenStakingConstants = await TokenStakingConstants.initialize(
      this.tokenStakingContract
    )
  }
}
