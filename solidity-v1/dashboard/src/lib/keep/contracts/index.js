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
import SaddleSwap from "../../../contracts-artifacts/SaddleSwap.json"
import AssetPool from "@keep-network/coverage-pools/artifacts/AssetPool.json"
import UnderwriterToken from "@keep-network/coverage-pools/artifacts/UnderwriterToken.json"
import RewardsPool from "@keep-network/coverage-pools/artifacts/RewardsPool.json"
import TBTCV2Token from "@keep-network/tbtc-v2/artifacts/TBTC.json"
import TBTCV2VendingMachine from "@keep-network/tbtc-v2/artifacts/VendingMachine.json"
import RiskManagerV1 from "@keep-network/coverage-pools/artifacts/RiskManagerV1.json"
import ThresholdTokenStaking from "@threshold-network/solidity-contracts/artifacts/TokenStaking.json"
import KeepStake from "@threshold-network/solidity-contracts/artifacts/KeepStake.json"

export const KEEP_TOKEN_CONTRACT_NAME = "keepTokenContract"
export const TOKEN_STAKING_CONTRACT_NAME = "stakingContract"
export const TOKEN_GRANT_CONTRACT_NAME = "grantContract"
export const OPERATOR_CONTRACT_NAME = "keepRandomBeaconOperatorContract"
export const REGISTRY_CONTRACT_NAME = "registryContract"
export const KEEP_OPERATOR_STATISTICS_CONTRACT_NAME =
  "keepRandomBeaconOperatorStatistics"
export const MANAGED_GRANT_FACTORY_CONTRACT_NAME = "managedGrantFactoryContract"
export const BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME =
  "bondedEcdsaKeepFactoryContract"
export const KEEP_BONDING_CONTRACT_NAME = "keepBondingContract"
export const TBTC_TOKEN_CONTRACT_NAME = "tbtcTokenContract"
export const TBTC_SYSTEM_CONTRACT_NAME = "tbtcSystemContract"
export const TOKEN_STAKING_ESCROW_CONTRACT_NAME = "tokenStakingEscrowContract"
export const OLD_TOKEN_STAKING_CONTRACT_NAME = "oldTokenStakingContract"
export const STAKING_PORT_BACKER_CONTRACT_NAME = "stakingPortBackerContract"
export const LP_REWARDS_TBTC_SADDLE_CONTRACT_NAME =
  "LPRewardsTBTCSaddleContract"
export const LP_REWARDS_KEEP_ETH_CONTRACT_NAME = "LPRewardsKEEPETHContract"
export const LP_REWARDS_TBTC_ETH_CONTRACT_NAME = "LPRewardsTBTCETHContract"
export const LP_REWARDS_KEEP_TBTC_CONTRACT_NAME = "LPRewardsKEEPTBTCContract"
export const KEEP_TOKEN_GEYSER_CONTRACT_NAME = "keepTokenGeyserContract"
export const ECDSA_REWARDS_DISTRRIBUTOR_CONTRACT_NAME =
  "ECDSARewardsDistributorContract"

export const SADDLE_SWAP_CONTRACT_NAME = "saddleSwapContract"
export const SaddleSwapArtifact = SaddleSwap

export const ASSET_POOL_CONTRACT_NAME = "assetPoolContract"
export const COV_TOKEN_CONTRACT_NAME = "covTokenContract"
export const RISK_MANAGER_V1_CONTRACT_NAME = "riskManagerV1Contract"
export const REWARDS_POOL_CONTRACT_NAME = "rewardsPoolContract"

export const RewardsPoolArtifact = RewardsPool

export const TBTCV2_TOKEN_CONTRACT_NAME = "tbtcV2Contract"
export const TBTCV2_VENDING_MACHINE_CONTRACT_NAME = "vendingMachineContract"

export const THRESHOLD_STAKING_CONTRACT_NAME = "thresholdStakingContract"
export const THRESHOLD_KEEP_STAKE_CONTRACT_NAME = "thresholdKeepStakeContract"

export const SIMPLE_PRE_APPLICATION_CONTRACT_NAME =
  "simplePREApplicationContract"

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
  [ASSET_POOL_CONTRACT_NAME]: {
    artifact: AssetPool,
  },
  [COV_TOKEN_CONTRACT_NAME]: {
    artifact: UnderwriterToken,
  },
  [TBTCV2_TOKEN_CONTRACT_NAME]: {
    artifact: TBTCV2Token,
  },
  [TBTCV2_VENDING_MACHINE_CONTRACT_NAME]: {
    artifact: TBTCV2VendingMachine,
  },
  [RISK_MANAGER_V1_CONTRACT_NAME]: {
    artifact: RiskManagerV1,
  },
  [REWARDS_POOL_CONTRACT_NAME]: {
    artifact: RewardsPool,
  },
  [THRESHOLD_STAKING_CONTRACT_NAME]: {
    artifact: ThresholdTokenStaking,
  },
  [THRESHOLD_KEEP_STAKE_CONTRACT_NAME]: {
    artifact: KeepStake,
  },
}

export default contracts

// The artifacts from @keep-network/keep-core for a given build only support a single network id
export function getFirstNetworkIdFromArtifact() {
  return Object.keys(KeepToken.networks)[0]
}
