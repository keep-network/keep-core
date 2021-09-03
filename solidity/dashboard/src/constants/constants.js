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
export const TBTC_TOKEN_CONTRACT_NAME = "tbtcTokenContract"
export const TBTC_SYSTEM_CONTRACT_NAME = "tbtcSystemContract"
export const TOKEN_STAKING_ESCROW_CONTRACT_NAME = "tokenStakingEscrow"
export const OLD_TOKEN_STAKING_CONTRACT_NAME = "oldTokenStakingContract"
export const STAKING_PORT_BACKER_CONTRACT_NAME = "stakingPortBackerContract"
export const LP_REWARDS_TBTC_SADDLE_CONTRACT_NAME = "LPRewardsTBTCSaddle"
export const LP_REWARDS_KEEP_ETH_CONTRACT_NAME = "LPRewardsKEEPETHContract"
export const LP_REWARDS_TBTC_ETH_CONTRACT_NAME = "LPRewardsTBTCETHContract"
export const LP_REWARDS_KEEP_TBTC_CONTRACT_NAME = "LPRewardsKEEPTBTCContract"
export const KEEP_TOKEN_GEYSER_CONTRACT_NAME = "keepTokenGeyserContract"
export const ECDSA_REWARDS_DISTRRIBUTOR_CONTRACT_NAME =
  "ECDSARewardsDistributorContract"

export const ASSET_POOL_CONTRACT_NAME = "assetPoolContract"

export const PENDING_STATUS = "PENDING"
export const COMPLETE_STATUS = "COMPLETE"

export const LINK = {
  discord: "https://discordapp.com/invite/wYezN7v",
  keepWebsite: "https://keep.network/",
  stakingDocumentation:
    "https://keep-network.gitbook.io/staking-documentation/",
}

export const WALLETS = {
  METAMASK: { label: "MetaMask", name: "METAMASK" },
  TREZOR: { label: "Trezor", name: "TREZOR" },
  LEDGER: { label: "Ledger", name: "LEDGER" },
  WALLET_CONNECT: { label: "WalletConnect", name: "WALLET_CONNECT" },
  EXPLORER_MODE: {
    label: "Explorer Mode",
    name: "EXPLORER_MODE",
  },
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
  TBTC_SADDLE: {
    contractName: LP_REWARDS_TBTC_SADDLE_CONTRACT_NAME,
    label: "TBTC + SADDLE",
    viewPoolLink: "https://saddle.exchange/#/deposit",
    pool: "SADDLE",
    lpTokens: [],
  },
  KEEP_ETH: {
    contractName: LP_REWARDS_KEEP_ETH_CONTRACT_NAME,
    label: "KEEP + ETH",
    viewPoolLink:
      "https://v2.info.uniswap.org/pair/0xe6f19dab7d43317344282f803f8e8d240708174a",
    address: "0xe6f19dab7d43317344282f803f8e8d240708174a",
    pool: "UNISWAP",
    lpTokens: [
      {
        tokenName: "KEEP",
        iconName: "KeepBlackGreen",
      },
      {
        tokenName: "ETH",
        iconName: "EthToken",
      },
    ],
  },
  KEEP_TBTC: {
    contractName: LP_REWARDS_KEEP_TBTC_CONTRACT_NAME,
    label: "KEEP + TBTC",
    viewPoolLink:
      "https://v2.info.uniswap.org/pair/0x38c8ffee49f286f25d25bad919ff7552e5daf081",
    address: "0x38c8ffee49f286f25d25bad919ff7552e5daf081",
    pool: "UNISWAP",
    lpTokens: [
      {
        tokenName: "KEEP",
        iconName: "KeepBlackGreen",
      },
      {
        tokenName: "TBTC",
        iconName: "TBTC",
      },
    ],
  },
  TBTC_ETH: {
    contractName: LP_REWARDS_TBTC_ETH_CONTRACT_NAME,
    label: "TBTC + ETH",
    viewPoolLink:
      "https://v2.info.uniswap.org/pair/0x854056fd40c1b52037166285b2e54fee774d33f6",
    address: "0x854056fd40c1b52037166285b2e54fee774d33f6",
    pool: "UNISWAP",
    lpTokens: [
      {
        tokenName: "TBTC",
        iconName: "TBTC",
      },
      {
        tokenName: "ETH",
        iconName: "EthToken",
      },
    ],
  },
  KEEP_ONLY: {
    contractName: KEEP_TOKEN_GEYSER_CONTRACT_NAME,
    label: "KEEP",
    pool: "TOKEN_GEYSER",
  },
}

export const AUTH_CONTRACTS_LABEL = {
  TBTC_SYSTEM: "TBTCSystem",
  BONDED_ECDSA_KEEP_FACTORY: "BondedECDSAKeepFactory",
  RANDOM_BEACON: "Keep Random Beacon Operator Contract",
}

export const TBTC_TOKEN_VERSION = {
  v1: "v1",
  v2: "v2",
}
