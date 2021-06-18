import KeepToken from "@keep-network/keep-core/artifacts/KeepToken.json"
import TokenStaking from "@keep-network/keep-core/artifacts/TokenStaking.json"
import TokenGrant from "@keep-network/keep-core/artifacts/TokenGrant.json"
import KeepRandomBeaconOperator from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperator.json"
import BondedECDSAKeepFactory from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeepFactory.json"
import TBTCSystem from "@keep-network/tbtc/artifacts/TBTCSystem.json"
import KeepBonding from "@keep-network/keep-ecdsa/artifacts/KeepBonding.json"
// import GuaranteedMinimumStakingPolicy from "@keep-network/keep-core/artifacts/GuaranteedMinimumStakingPolicy.json"
// import PermissiveStakingPolicy from "@keep-network/keep-core/artifacts/PermissiveStakingPolicy.json"
import KeepRandomBeaconOperatorStatistics from "@keep-network/keep-core/artifacts/KeepRandomBeaconOperatorStatistics.json"
// import ManagedGrant from "@keep-network/keep-core/artifacts/ManagedGrant.json"
import ManagedGrantFactory from "@keep-network/keep-core/artifacts/ManagedGrantFactory.json"
import TBTCToken from "@keep-network/tbtc/artifacts/TBTCToken.json"
// import Deposit from "@keep-network/tbtc/artifacts/Deposit.json"
// import BondedECDSAKeep from "@keep-network/keep-ecdsa/artifacts/BondedECDSAKeep.json"
import TokenStakingEscrow from "@keep-network/keep-core/artifacts/TokenStakingEscrow.json"
import StakingPortBacker from "@keep-network/keep-core/artifacts/StakingPortBacker.json"
import BeaconRewards from "@keep-network/keep-core/artifacts/BeaconRewards.json"
import ECDSARewardsDistributor from "@keep-network/keep-ecdsa/artifacts/ECDSARewardsDistributor.json"
import LPRewardsKEEPETH from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPETH.json"
import LPRewardsTBTCETH from "@keep-network/keep-ecdsa/artifacts/LPRewardsTBTCETH.json"
import LPRewardsKEEPTBTC from "@keep-network/keep-ecdsa/artifacts/LPRewardsKEEPTBTC.json"
import LPRewardsTBTCSaddle from "@keep-network/keep-ecdsa/artifacts/LPRewardsTBTCSaddle.json"
import KeepOnlyPool from "@keep-network/keep-core/artifacts/KeepVault.json"
// import IERC20 from "@keep-network/keep-core/artifacts/IERC20.json"
// import SaddleSwap from "../../contracts-artifacts/SaddleSwap.json"

import {
  KEEP_TOKEN_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  OPERATOR_CONTRACT_NAME,
  KEEP_OPERATOR_STATISTICS_CONTRACT_NAME,
  MANAGED_GRANT_FACTORY_CONTRACT_NAME,
  KEEP_BONDING_CONTRACT_NAME,
  TBTC_TOKEN_CONTRACT_NAME,
  TBTC_SYSTEM_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
  BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
  STAKING_PORT_BACKER_CONTRACT_NAME,
  //   OLD_TOKEN_STAKING_CONTRACT_NAME,
  LP_REWARDS_KEEP_ETH_CONTRACT_NAME,
  LP_REWARDS_TBTC_ETH_CONTRACT_NAME,
  LP_REWARDS_KEEP_TBTC_CONTRACT_NAME,
  LP_REWARDS_TBTC_SADDLE_CONTRACT_NAME,
  KEEP_TOKEN_GEYSER_CONTRACT_NAME,
} from "../../../constants/constants"

const contracts = {
  [KEEP_TOKEN_CONTRACT_NAME]: { artifact: KeepToken },
  [TOKEN_GRANT_CONTRACT_NAME]: { artifact: TokenGrant },
  [OPERATOR_CONTRACT_NAME]: {
    artifact: KeepRandomBeaconOperator,
  },
  [TOKEN_STAKING_CONTRACT_NAME]: {
    artifact: TokenStaking,
  },
  [KEEP_OPERATOR_STATISTICS_CONTRACT_NAME]: {
    artifact: KeepRandomBeaconOperatorStatistics,
  },
  [MANAGED_GRANT_FACTORY_CONTRACT_NAME]: {
    artifact: ManagedGrantFactory,
  },
  [KEEP_BONDING_CONTRACT_NAME]: {
    artifact: KeepBonding,
  },
  [TBTC_TOKEN_CONTRACT_NAME]: {
    artifact: TBTCToken,
  },
  [TBTC_SYSTEM_CONTRACT_NAME]: {
    artifact: TBTCSystem,
  },
  [TOKEN_STAKING_ESCROW_CONTRACT_NAME]: {
    artifact: TokenStakingEscrow,
  },
  [BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME]: {
    artifact: BondedECDSAKeepFactory,
  },
  [STAKING_PORT_BACKER_CONTRACT_NAME]: {
    artifact: StakingPortBacker,
  },
  beaconRewardsContract: {
    artifact: BeaconRewards,
  },
  ECDSARewardsDistributorContract: {
    artifact: ECDSARewardsDistributor,
  },
  [LP_REWARDS_KEEP_ETH_CONTRACT_NAME]: {
    artifact: LPRewardsKEEPETH,
  },
  [LP_REWARDS_TBTC_ETH_CONTRACT_NAME]: {
    artifact: LPRewardsTBTCETH,
  },
  [LP_REWARDS_KEEP_TBTC_CONTRACT_NAME]: {
    artifact: LPRewardsKEEPTBTC,
  },
  [LP_REWARDS_TBTC_SADDLE_CONTRACT_NAME]: {
    artifact: LPRewardsTBTCSaddle,
  },
  [KEEP_TOKEN_GEYSER_CONTRACT_NAME]: {
    artifact: KeepOnlyPool,
  },
}

export default contracts
