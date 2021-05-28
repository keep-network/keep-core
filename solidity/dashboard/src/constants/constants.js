export const KEEP_TOKEN_CONTRACT_NAME = "token"
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
export const BOND_ERC20_CONTRACT_NAME = 'bondERC20'
export const TBTC_TOKEN_CONTRACT_NAME = "tbtcTokenContract"
export const TBTC_SYSTEM_CONTRACT_NAME = "tbtcSystemContract"
export const TOKEN_STAKING_ESCROW_CONTRACT_NAME = "tokenStakingEscrow"
export const OLD_TOKEN_STAKING_CONTRACT_NAME = "oldTokenStakingContract"
export const STAKING_PORT_BACKER_CONTRACT_NAME = "stakingPortBackerContract"
export const LP_REWARDS_KEEP_ETH_CONTRACT_NAME = "LPRewardsKEEPETHContract"
export const LP_REWARDS_TBTC_ETH_CONTRACT_NAME = "LPRewardsTBTCETHContract"
export const LP_REWARDS_KEEP_TBTC_CONTRACT_NAME = "LPRewardsKEEPTBTCContract"

export const PENDING_STATUS = "PENDING"
export const COMPLETE_STATUS = "COMPLETE"

export const WALLETS = {
  METAMASK: { label: "MetaMask" },
  TREZOR: { label: "Trezor" },
  LEDGER: { label: "Ledger" },
  COINBASE: { label: "Coinbase" },
}

export const REWARD_STATUS = {
  AVAILABLE: "AVAILABLE",
  WITHDRAWN: "WITHDRAWN",
  ACCUMULATING: "ACCUMULATING",
}

export const SIGNING_GROUP_STATUS = {
  COMPLETED: "Completed",
  TERMINATED: "Terminated",
  ACTIVE: "Active work",
}

export const LIQUIDITY_REWARD_PAIRS = {
  KEEP_ETH: {
    contractName: LP_REWARDS_KEEP_ETH_CONTRACT_NAME,
    label: "KEEP + ETH",
    viewPoolLink:
      "https://info.uniswap.org/pair/0xe6f19dab7d43317344282f803f8e8d240708174a",
    rewardPoolPerWeek: 150000,
    address: "0xe6f19dab7d43317344282f803f8e8d240708174a",
  },
  KEEP_TBTC: {
    contractName: LP_REWARDS_KEEP_TBTC_CONTRACT_NAME,
    label: "KEEP + TBTC",
    viewPoolLink:
      "https://info.uniswap.org/pair/0x38c8ffee49f286f25d25bad919ff7552e5daf081",
    rewardPoolPerWeek: 200000,
    address: "0x38c8ffee49f286f25d25bad919ff7552e5daf081",
  },
  TBTC_ETH: {
    contractName: LP_REWARDS_TBTC_ETH_CONTRACT_NAME,
    label: "TBTC + ETH",
    viewPoolLink:
      "https://info.uniswap.org/pair/0x854056fd40c1b52037166285b2e54fee774d33f6",
    rewardPoolPerWeek: 50000,
    address: "0x854056fd40c1b52037166285b2e54fee774d33f6",
  },
}
