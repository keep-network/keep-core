import {
  createSaddleSwapContract,
  createSaddleTBTCMetaPool,
  createSaddleTBTCMetaPoolV2,
} from "../contracts"

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
export const LP_REWARDS_TBTCV2_SADDLE_CONTRACT_NAME = "LPRewardsTBTCv2Saddle"
export const LP_REWARDS_TBTCV2_SADDLEV2_CONTRACT_NAME =
  "LPRewardsTBTCv2SaddleV2"

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
  coveragePools: {
    docs: "https://github.com/keep-network/coverage-pools/blob/main/docs/design.adoc",
    etherscanCovKeepToken:
      "https://etherscan.io/token/0x2aa24dac5e494e7b028ed43023530e5769df5d8b",
  },
  tbtcMigration: {
    docs: "https://coda.io/@keep-network/how-to-mint-tbtc-v2-with-etherscan",
  },
  pools: {
    saddle: {
      tbtc: "https://saddle.exchange/#/pools/tbtc/deposit",
      tbtcV2: "https://saddle.exchange/#/pools/tbtcv2/deposit",
    },
    uniswap: {
      tbtcETH: `https://app.uniswap.org/#/add/v2/0x8daebade922df735c38c80c7ebd708af50815faa/ETH`,
    },
  },
  proposals: {
    shiftingIncentivesToV2Metapool:
      "https://forum.keep.network/t/move-saddle-tbtc-liquidity-incentives-to-v2-metapool/380",
    shiftingIncentivesToCoveragePools:
      "https://forum.keep.network/t/shifting-incentives-towards-tbtc-v2-and-coverage-pool-version-2/322",
    removeIncentivesForKEEPTBTCpool:
      "https://forum.keep.network/t/proposal-remove-incentives-for-the-keep-tbtc-pool/56",
    removeIncentivesForTBTCETHpool:
      "https://forum.keep.network/t/proposal-to-remove-incentives-for-tbtc-eth-pool/341",
    repurposingKEEPETHLiquidityIncentives:
      "https://forum.keep.network/t/repurposing-keep-eth-liquidity-incentives/387",
  },
  tbtcDapp: "https://dapp.tbtc.network",
  thresholdDapp: "https://dashboard.threshold.network/",
  setUpPRE: "https://docs.nucypher.com/en/latest/",
}

export const WALLETS = {
  TALLY: { label: "Tally", name: "TALLY" },
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

export const POOL_TYPE = {
  SADDLE: "SADDLE",
  MSTABLE: "MSTABLE",
  UNISWAP: "UNISWAP",
  TOKEN_GEYSER: "TOKEN_GEYSER", // KEEP_ONLY
}

export const LIQUIDITY_REWARD_PAIRS = {
  TBTCV2_SADDLE_META_V2: {
    contractName: LP_REWARDS_TBTCV2_SADDLEV2_CONTRACT_NAME,
    label: "TBTC V2 + SADDLE Meta V2",
    viewPoolLink: LINK.pools.saddle.tbtcV2,
    pool: POOL_TYPE.SADDLE,
    lpTokens: [],
    options: {
      createSwapContract: createSaddleTBTCMetaPoolV2,
      poolTokens: [
        { name: "TBTC-V2", decimals: 18 },
        { name: "saddleBTC-V2", decimals: 18 },
      ],
    },
  },
  TBTCV2_SADDLE: {
    contractName: LP_REWARDS_TBTCV2_SADDLE_CONTRACT_NAME,
    label: "TBTC V2 + SADDLE",
    viewPoolLink: LINK.pools.saddle.tbtc,
    pool: POOL_TYPE.SADDLE,
    lpTokens: [],
    options: {
      createSwapContract: createSaddleTBTCMetaPool,
      poolTokens: [
        { name: "TBTC", decimals: 18 },
        { name: "saddleBTC-V2", decimals: 18 },
      ],
    },
  },
  TBTC_SADDLE: {
    contractName: LP_REWARDS_TBTC_SADDLE_CONTRACT_NAME,
    label: "TBTC + SADDLE",
    viewPoolLink: "https://saddle.exchange/#/deposit",
    pool: POOL_TYPE.SADDLE,
    lpTokens: [],
    options: {
      createSwapContract: createSaddleSwapContract,
      poolTokens: [
        { name: "TBTC", decimals: 18 },
        { name: "WBTC", decimals: 8 },
        { name: "RENBTC", decimals: 8 },
        { name: "SBTC", decimals: 18 },
      ],
    },
  },
  KEEP_ETH: {
    contractName: LP_REWARDS_KEEP_ETH_CONTRACT_NAME,
    label: "KEEP + ETH",
    viewPoolLink:
      "https://v2.info.uniswap.org/pair/0xe6f19dab7d43317344282f803f8e8d240708174a",
    address: "0xe6f19dab7d43317344282f803f8e8d240708174a",
    pool: POOL_TYPE.UNISWAP,
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
  TBTCV2_MBTC: {
    label: "TBTCv2 + mBTC",
    viewPoolLink:
      "https://mstable.app/#/mbtc/pools/0xc3280306b6218031e61752d060b091278d45c329",
    pool: POOL_TYPE.MSTABLE,
  },
  KEEP_TBTC: {
    contractName: LP_REWARDS_KEEP_TBTC_CONTRACT_NAME,
    label: "KEEP + TBTC",
    viewPoolLink:
      "https://v2.info.uniswap.org/pair/0x38c8ffee49f286f25d25bad919ff7552e5daf081",
    address: "0x38c8ffee49f286f25d25bad919ff7552e5daf081",
    pool: POOL_TYPE.UNISWAP,
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
    pool: POOL_TYPE.UNISWAP,
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
    pool: POOL_TYPE.TOKEN_GEYSER,
  },
}

export const AUTH_CONTRACTS_LABEL = {
  TBTC_SYSTEM: "TBTCSystem",
  BONDED_ECDSA_KEEP_FACTORY: "BondedECDSAKeepFactory",
  RANDOM_BEACON: "Keep Random Beacon Operator Contract",
  THRESHOLD_TOKEN_STAKING: "Threshold Staking",
}

export const TBTC_TOKEN_VERSION = {
  v1: "v1",
  v2: "v2",
}

/**
 * Enum defines a supported types of Modals
 * @readonly
 * @enum {string}
 */
export const MODAL_TYPES = {
  MobileUsers: "MobileUsers",
  BondingAddETH: "BondingAddEth",
  BondingWithdrawETH: "BondingWithdrawETH",
  MetaMask: "MetaMask",
  Tally: "Tally",
  ExplorerMode: "ExplorerMode",
  Ledger: "Ledger",
  Trezor: "Trezor",
  WalletConnect: "WalletConnect",
  WalletSelection: "WalletSelection",
  DelegationAlreadyExists: "DelegationAlreadyExists",
  TopUpInitiatedConfirmation: "TopUpInitiatedConfirmation",
  TopUpInitialization: "TopUpInitialization",
  ConfirmTopUpInitialization: "ConfirmTopUpInitialization",
  KeepOnlyPoolAddKeep: "KeepOnlyPoolAddKeep",
  KeepOnlyPoolWithdrawKeep: "KeepOnlyPoolWithdrawKeep",
  ConfirmDelegation: "ConfirmDelegation",
  ConfirmRecovering: "ConfirmRecovering",
  ConfirmCancelDelegationFromGrant: "ConfirmCancelDelegationFromGrant",
  UndelegateStake: "UndelegateStake",
  UndelegationInitiated: "UndelegationInitiated",
  ClaimStakingTokens: "ClaimStakingTokens",
  StakingTokensClaimed: "StakingTokensClaimed",
  GrantTokensWithdrawn: "GrantTokensWithdrawn",
  CopyStake: "CopyStake",
  ConfirmTBTCMigration: "ConfirmTBTCMigration",
  TBTCMigrationCompleted: "TBTCMigrationCompleted",
  ConfirmReleaseTokensFromGrant: "ConfirmReleaseTokensFromGrant",
  WarningBeforeCovPoolDeposit: "WarningBeforeCovPoolDeposit",
  InitiateCovPoolDeposit: "InitiateCovPoolDeposit",
  InitiateCovPoolWithdraw: "InitiateCovPoolWithdraw",
  CovPoolWithdrawInitialized: "CovPoolWithdrawInitialized",
  CovPoolClaimTokens: "CovPoolClaimTokens",
  ReInitiateCovPoolWithdraw: "ReInitiateCovPoolWithdraw",
  ConfirmCovPoolIncreaseWithdrawal: "ConfirmCovPoolIncreaseWithdrawal",
  IncreaseCovPoolWithdrawal: "IncreaseCovPoolWithdrawal",
  WithdrawGrantedTokens: "WithdrawGrantedTokens",
  AuthorizeAndStakeOnThreshold: "AuthorizeAndStakeOnThreshold",
  StakeOnThresholdWithoutAuthorization: "StakeOnThresholdWithoutAuthorization",
  StakeOnThresholdConfirmed: "StakeOnThresholdConfirmed",
  ThresholdAuthorizationLoadingModal: "ThresholdAuthorizationLoadingModal",
  ThresholdStakeConfirmationLoadingModal:
    "ThresholdStakeConfirmationLoadingModal",
  AuthorizedButNotStakedToTWarningModal:
    "AuthorizedButNotStakedToTWarningModal",
  ContactYourGrantManagerWarning: "ContactYourGrantManagerWarning",
}

export const COV_POOL_TIMELINE_STEPS = {
  DEPOSITED_TOKENS: 1,
  WITHDRAW_DEPOSIT: 2,
  COOLDOWN: 3,
  CLAIM_TOKENS: 4,
}

export const STAKE_ON_THRESHOLD_TIMELINE_STEPS = {
  NONE: 0,
  AUTHORIZE_CONTRACT: 1,
  CONFIRM_STAKE: 2,
  SET_UP_PRE: 3,
}

export const COV_POOLS_FORMS_MAX_DECIMAL_PLACES = 6

/**
 * Enum defines cov pools withdrawal status
 * @readonly
 * @enum {string}
 */
export const PENDING_WITHDRAWAL_STATUS = {
  NONE: "none",
  PENDING: "pending",
  COMPLETED: "completed",
  EXPIRED: "expired",
  NEW: "new",
}

export const ADD_TO_CALENDAR_OPTIONS = {
  GOOGLE_CALENDER: "google-calendar",
  APPLE_CALENDAR: "apple-calendar",
}

export const COVERAGE_POOL_CLAIM_TOKENS_CALENDAR_EVENT = {
  name: "Coverage Pools - Tokens Ready To Claim",
  details: "You have 48 hours to claim your tokens!",
  // eslint-disable-next-line no-template-curly-in-string
  location: "https://dashboard.keep.network/${address}/coverage-pools/deposit",
}

export const UNDELEGATE_STAKE_CALENDAR_EVENT = {
  name: "Stake Undelegation - Tokens Ready To Claim",
  details:
    "The stake has been undelegated! The tokens are ready to be claimed!",
  // eslint-disable-next-line no-template-curly-in-string
  location: "https://dashboard.keep.network/${address}/overview",
}

export const GRANT_MANAGER_EMAIL = "grantmanager@keep.network"
